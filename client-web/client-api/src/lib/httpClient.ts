import type { EndpointContract } from "../generated/apiCatalog";
import type { RequestResult, RuntimeEnvironment } from "../data/runtime";

export type SendRequestOptions = {
  endpoint: EndpointContract;
  environment: RuntimeEnvironment;
  bodyText: string;
  pathValues: Record<string, string>;
  queryText: string;
};

const retryableStatuses = new Set([408, 425, 429, 500, 502, 503, 504]);
const maxAttempts = 3;
const requestTimeoutMs = 15000;
const bearerTokenPattern = /^[A-Za-z0-9._~+/=-]+$/;

function applyPathValues(path: string, values: Record<string, string>): string {
  return path.replace(/\{([^}]+)\}/g, (_match, key: string) => {
    return encodeURIComponent(values[key] || `{${key}}`);
  });
}

function buildUrl(endpoint: EndpointContract, environment: RuntimeEnvironment, path: string, queryText: string): string {
  const query = queryText.trim();
  const suffix = query ? `${path}${path.includes("?") ? "&" : "?"}${query.replace(/^\?/, "")}` : path;

  if (environment.mode === "proxy") {
    return `/__erp/${endpoint.service}${suffix}`;
  }

  const baseUrl = environment.baseUrls[endpoint.service] ?? "";
  return `${baseUrl.replace(/\/$/, "")}${suffix}`;
}

async function parseBody(response: Response): Promise<unknown> {
  const text = await response.text();
  if (!text) return null;

  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function buildRequestId(prefix: string): string {
  const randomPart = typeof crypto !== "undefined" && "randomUUID" in crypto
    ? crypto.randomUUID()
    : `${Date.now()}-${Math.random().toString(16).slice(2)}`;

  return `${prefix || "console"}-${randomPart}`;
}

function sanitizeBearerToken(token: string): string {
  const trimmed = token.trim();

  if (!trimmed) return "";
  if (!bearerTokenPattern.test(trimmed)) {
    throw new Error("Bearer token contem caracteres invalidos para um header HTTP.");
  }

  return trimmed;
}

function shouldRetry(method: string, response?: Response): boolean {
  if (response && !retryableStatuses.has(response.status)) return false;
  return ["GET", "HEAD", "OPTIONS", "PUT", "PATCH", "DELETE"].includes(method);
}

function wait(ms: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

async function fetchWithTimeout(url: string, init: RequestInit): Promise<Response> {
  const controller = new AbortController();
  const timeout = window.setTimeout(() => controller.abort(), requestTimeoutMs);

  try {
    return await fetch(url, { ...init, signal: controller.signal });
  } finally {
    window.clearTimeout(timeout);
  }
}

export async function sendEndpointRequest(options: SendRequestOptions): Promise<RequestResult> {
  const started = performance.now();
  const path = applyPathValues(options.endpoint.path, options.pathValues);
  const url = buildUrl(options.endpoint, options.environment, path, options.queryText);
  const method = options.endpoint.method;
  const headers: Record<string, string> = {
    "content-type": "application/json",
    "x-correlation-id": options.environment.correlationId || buildRequestId("console")
  };

  try {
    const bearerToken = sanitizeBearerToken(options.environment.bearerToken);
    if (bearerToken) {
      headers.authorization = `Bearer ${bearerToken}`;
    }
  } catch (error) {
    return {
      status: 0,
      statusText: "Invalid Request",
      durationMs: Math.round(performance.now() - started),
      headers: {},
      body: {
        error: error instanceof Error ? error.message : "Token invalido",
        hint: "Remova espacos, quebras de linha e caracteres de controle do token bearer."
      },
      url,
      error: error instanceof Error ? error.message : "Token invalido"
    };
  }

  if (!["GET", "HEAD", "OPTIONS"].includes(method)) {
    headers["idempotency-key"] = buildRequestId(options.environment.idempotencyKey || "console");
  }

  try {
    let response: Response | undefined;
    let lastError: unknown;

    for (let attempt = 1; attempt <= maxAttempts; attempt += 1) {
      try {
        response = await fetchWithTimeout(url, {
          method,
          headers,
          body: ["GET", "HEAD", "OPTIONS"].includes(method)
            ? undefined
            : options.bodyText.trim() || "{}"
        });

        if (!shouldRetry(method, response) || attempt === maxAttempts) {
          break;
        }
      } catch (error) {
        lastError = error;
        if (!shouldRetry(method) || attempt === maxAttempts) {
          throw error;
        }
      }

      await wait(150 * 2 ** (attempt - 1));
    }

    if (!response) {
      throw lastError instanceof Error ? lastError : new Error("Falha desconhecida");
    }

    const responseHeaders: Record<string, string> = {};
    response.headers.forEach((value, key) => {
      responseHeaders[key] = value;
    });

    return {
      status: response.status,
      statusText: response.statusText,
      durationMs: Math.round(performance.now() - started),
      headers: responseHeaders,
      body: await parseBody(response),
      url
    };
  } catch (error) {
    return {
      status: 0,
      statusText: "Network Error",
      durationMs: Math.round(performance.now() - started),
      headers: {},
      body: {
        error: error instanceof Error ? error.message : "Falha desconhecida",
        hint: "Confirme se o stack do ERP esta rodando e se o modo proxy/direct esta adequado."
      },
      url,
      error: error instanceof Error ? error.message : "Falha desconhecida"
    };
  }
}

export function curlFor(options: SendRequestOptions): string {
  const path = applyPathValues(options.endpoint.path, options.pathValues);
  const url = buildUrl(options.endpoint, options.environment, path, options.queryText);
  const idempotencyKey = buildRequestId(options.environment.idempotencyKey || "console");
  const lines = [
    `curl -X ${options.endpoint.method} '${url}'`,
    `  -H 'content-type: application/json'`,
    `  -H 'x-correlation-id: ${options.environment.correlationId}'`
  ];

  if (options.environment.bearerToken.trim()) {
    lines.push(`  -H 'authorization: Bearer ${options.environment.bearerToken.trim()}'`);
  }

  if (!["GET", "HEAD", "OPTIONS"].includes(options.endpoint.method)) {
    lines.push(`  -H 'idempotency-key: ${idempotencyKey}'`);
    lines.push(`  -d '${options.bodyText.trim() || "{}"}'`);
  }

  return lines.join(" \\\n");
}
