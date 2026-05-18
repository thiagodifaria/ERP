import { IncomingMessage, ServerResponse } from "node:http";
import { z } from "zod";

const jsonObjectSchema = z.record(z.unknown());

export function json(response: ServerResponse, statusCode: number, body: unknown): void {
  response.writeHead(statusCode, { "content-type": "application/json" });
  response.end(JSON.stringify(body));
}

export async function readJson<T>(request: IncomingMessage): Promise<T> {
  const chunks: Buffer[] = [];
  let totalBytes = 0;

  for await (const chunk of request) {
    const buffer = Buffer.from(chunk);
    totalBytes += buffer.length;
    if (totalBytes > 1_048_576) {
      throw new Error("payload_too_large");
    }
    chunks.push(buffer);
  }

  const rawBody = Buffer.concat(chunks).toString("utf8");

  if (rawBody.length === 0) {
    throw new Error("invalid_json");
  }

  const parsed = JSON.parse(rawBody) as unknown;
  assertSafeJson(parsed);
  return parsed as T;
}

export async function readJsonSchema<T>(request: IncomingMessage): Promise<T> {
  return jsonObjectSchema.parse(await readJson<unknown>(request)) as T;
}

function assertSafeJson(value: unknown): void {
  if (Array.isArray(value)) {
    value.forEach(assertSafeJson);
    return;
  }

  if (value && typeof value === "object") {
    for (const [key, child] of Object.entries(value as Record<string, unknown>)) {
      if (key === "__proto__" || key === "prototype" || key === "constructor") {
        throw new Error("invalid_json_key");
      }
      assertSafeJson(child);
    }
  }
}
