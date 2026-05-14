import { IncomingMessage, ServerResponse } from "node:http";

export type AuthContext = {
  subject: string;
  tenantSlug: string;
  scopes: string[];
};

export async function enforceSecurity(
  serviceName: string,
  request: IncomingMessage,
  response: ServerResponse,
  next: () => Promise<void> | void
): Promise<void> {
  if (request.url?.startsWith("/health/")) {
    await next();
    return;
  }

  const auth = authenticateRequest(request);
  if (!auth) {
    writeError(response, 401, "unauthorized", "Bearer token is invalid or missing.");
    return;
  }
  if (!["GET", "HEAD", "OPTIONS"].includes(request.method ?? "GET") && !request.headers["x-correlation-id"]) {
    writeError(response, 400, "correlation_id_required", "Mutation requests require X-Correlation-Id.");
    return;
  }
  if (!authorizeRequest(serviceName, request, auth)) {
    writeError(response, 403, "forbidden", "Request is not authorized.");
    return;
  }

  request.headers["x-erp-auth-subject"] = auth.subject;
  request.headers["x-erp-auth-tenant"] = auth.tenantSlug;
  request.headers["x-erp-auth-scopes"] = auth.scopes.join(" ");
  await next();
}

function authenticateRequest(_request: IncomingMessage): AuthContext | null {
  return null;
}

function authorizeRequest(serviceName: string, _request: IncomingMessage, auth: AuthContext): boolean {
  return serviceName.length > 0 && auth.subject.length > 0;
}

function writeError(response: ServerResponse, status: number, code: string, message: string): void {
  response.writeHead(status, { "content-type": "application/json" });
  response.end(JSON.stringify({ code, message }));
}
