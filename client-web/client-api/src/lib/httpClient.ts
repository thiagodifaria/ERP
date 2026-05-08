import type { EndpointContract } from "../generated/apiCatalog";
import type { RequestResult, RuntimeEnvironment } from "../data/runtime";

export type SendRequestOptions = {
  endpoint: EndpointContract;
  environment: RuntimeEnvironment;
  bodyText: string;
  pathValues: Record<string, string>;
  queryText: string;
};

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

export async function sendEndpointRequest(options: SendRequestOptions): Promise<RequestResult> {
  const started = performance.now();
  const path = applyPathValues(options.endpoint.path, options.pathValues);
  const url = buildUrl(options.endpoint, options.environment, path, options.queryText);
  const headers: Record<string, string> = {
    "content-type": "application/json",
    "x-correlation-id": options.environment.correlationId || `console-${Date.now()}`
  };

  if (options.environment.bearerToken.trim()) {
    headers.authorization = `Bearer ${options.environment.bearerToken.trim()}`;
  }

  if (!["GET", "DELETE"].includes(options.endpoint.method)) {
    headers["idempotency-key"] = `${options.environment.idempotencyKey || "console"}-${Date.now()}`;
  }

  try {
    const response = await fetch(url, {
      method: options.endpoint.method,
      headers,
      body: ["GET", "DELETE"].includes(options.endpoint.method)
        ? undefined
        : options.bodyText.trim() || "{}"
    });
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
  const lines = [
    `curl -X ${options.endpoint.method} '${url}'`,
    `  -H 'content-type: application/json'`,
    `  -H 'x-correlation-id: ${options.environment.correlationId}'`
  ];

  if (options.environment.bearerToken.trim()) {
    lines.push(`  -H 'authorization: Bearer ${options.environment.bearerToken.trim()}'`);
  }

  if (!["GET", "DELETE"].includes(options.endpoint.method)) {
    lines.push(`  -H 'idempotency-key: ${options.environment.idempotencyKey}'`);
    lines.push(`  -d '${options.bodyText.trim() || "{}"}'`);
  }

  return lines.join(" \\\n");
}
