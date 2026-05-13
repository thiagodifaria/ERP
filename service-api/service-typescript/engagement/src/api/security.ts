import { createHmac, timingSafeEqual } from "node:crypto";
import type { IncomingMessage, ServerResponse } from "node:http";

type Claims = {
  sub?: string;
  user_public_id?: string;
  tenant_slug?: string;
  tenant?: string;
  scope?: string | string[];
  exp?: number;
};

type AuthContext = {
  subject: string;
  tenantSlug: string;
  scopes: string[];
};

export async function enforceSecurity(
  serviceName: string,
  request: IncomingMessage,
  response: ServerResponse,
  next: () => Promise<void> | void,
): Promise<void> {
  if (!securityEnforced() || (request.url ?? "").startsWith("/health/")) {
    await next();
    return;
  }

  if (requiresCorrelation(request.method) && !request.headers["x-correlation-id"]) {
    writeSecurityError(response, 400, "correlation_id_required", "Mutation requests require X-Correlation-Id.");
    return;
  }

  const auth = authenticate(request);
  if (!auth) {
    writeSecurityError(response, 401, "unauthorized", "Bearer token is invalid or missing.");
    return;
  }

  request.headers["x-erp-auth-subject"] = auth.subject;
  request.headers["x-erp-auth-tenant"] = auth.tenantSlug;
  request.headers["x-erp-auth-scopes"] = auth.scopes.join(" ");

  const allowed = await authorizeOpenFga(serviceName, request, auth);
  if (!allowed) {
    writeSecurityError(response, 403, "openfga_denied", "OpenFGA denied the request.");
    return;
  }

  await next();
}

function securityEnforced(): boolean {
  const mode = (process.env.ERP_AUTH_ENFORCEMENT ?? "").trim().toLowerCase();
  if (["disabled", "off", "false"].includes(mode)) return false;
  if (["enforced", "strict", "true"].includes(mode)) return true;
  const environment = (process.env.ERP_ENV ?? "local").trim().toLowerCase();
  return !["", "local", "dev", "development", "test", "testing"].includes(environment);
}

function authenticate(request: IncomingMessage): AuthContext | null {
  const authorization = String(request.headers.authorization ?? "");
  if (!authorization.toLowerCase().startsWith("bearer ")) return null;
  const token = authorization.slice("Bearer ".length).trim();
  const internalToken = (process.env.ERP_INTERNAL_SERVICE_TOKEN ?? "").trim();
  if (internalToken && fixedTimeEquals(token, internalToken)) {
    return { subject: "service:internal", tenantSlug: resolveTenant(request), scopes: ["service"] };
  }

  const claims = verifyJwt(token);
  if (!claims) return null;
  const subject = claims.sub ?? claims.user_public_id ?? "";
  const tenantSlug = claims.tenant_slug ?? claims.tenant ?? resolveTenant(request);
  const scopes = Array.isArray(claims.scope) ? claims.scope : String(claims.scope ?? "").split(/\s+/).filter(Boolean);
  return subject ? { subject, tenantSlug, scopes } : null;
}

function verifyJwt(token: string): Claims | null {
  const secret = process.env.ERP_JWT_HS256_SECRET ?? "";
  const parts = token.split(".");
  if (!secret || parts.length !== 3) return null;
  try {
    const header = JSON.parse(base64UrlDecode(parts[0]).toString("utf8")) as { alg?: string };
    if (header.alg !== "HS256") return null;
    const expected = createHmac("sha256", secret).update(`${parts[0]}.${parts[1]}`).digest();
    if (!fixedTimeEquals(parts[2], base64UrlEncode(expected))) return null;
    const claims = JSON.parse(base64UrlDecode(parts[1]).toString("utf8")) as Claims;
    if (typeof claims.exp === "number" && claims.exp <= Math.floor(Date.now() / 1000)) return null;
    return claims;
  } catch {
    return null;
  }
}

async function authorizeOpenFga(serviceName: string, request: IncomingMessage, auth: AuthContext): Promise<boolean> {
  if ((process.env.ERP_OPENFGA_ENFORCEMENT ?? "").toLowerCase() !== "true") return true;
  const baseUrl = (process.env.OPENFGA_BASE_URL ?? "").replace(/\/$/, "");
  const storeId = process.env.OPENFGA_STORE_ID ?? "";
  if (!baseUrl || !storeId) return false;
  const relation = requiresCorrelation(request.method) ? "write" : "read";
  const targetObject = auth.tenantSlug ? `tenant:${normalize(auth.tenantSlug)}` : `service:${normalize(serviceName)}`;
  const payload: Record<string, unknown> = {
    tuple_key: {
      user: auth.subject.startsWith("service:") ? auth.subject : `user:${auth.subject}`,
      relation,
      object: targetObject,
    },
  };
  if (process.env.OPENFGA_AUTHORIZATION_MODEL_ID) {
    payload.authorization_model_id = process.env.OPENFGA_AUTHORIZATION_MODEL_ID;
  }
  try {
    const result = await fetch(`${baseUrl}/stores/${storeId}/check`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(2000),
    });
    if (!result.ok) return false;
    const body = (await result.json()) as { allowed?: boolean };
    return body.allowed === true;
  } catch {
    return false;
  }
}

function resolveTenant(request: IncomingMessage): string {
  const headerTenant = request.headers["x-tenant-slug"] ?? request.headers["x-erp-tenant-slug"];
  if (typeof headerTenant === "string" && headerTenant.trim()) return headerTenant.trim();
  const url = new URL(request.url ?? "/", "http://local");
  return url.searchParams.get("tenant_slug") ?? "";
}

function requiresCorrelation(method = "GET"): boolean {
  return !["GET", "HEAD", "OPTIONS"].includes(method.toUpperCase());
}

function base64UrlDecode(value: string): Buffer {
  return Buffer.from(value.replace(/-/g, "+").replace(/_/g, "/"), "base64");
}

function base64UrlEncode(value: Buffer): string {
  return value.toString("base64").replace(/=/g, "").replace(/\+/g, "-").replace(/\//g, "_");
}

function fixedTimeEquals(left: string, right: string): boolean {
  const leftBuffer = Buffer.from(left);
  const rightBuffer = Buffer.from(right);
  return leftBuffer.length === rightBuffer.length && timingSafeEqual(leftBuffer, rightBuffer);
}

function normalize(value: string): string {
  return value.trim().toLowerCase().replace(/\s+/g, "-");
}

function writeSecurityError(response: ServerResponse, statusCode: number, code: string, message: string): void {
  response.writeHead(statusCode, { "content-type": "application/json" });
  response.end(JSON.stringify({ code, message }));
}
