import assert from "node:assert/strict";
import test from "node:test";
import { createServer } from "node:http";
import { route } from "../../src/api/router.js";

test("GET /health/details exposes bootstrap dependencies", async () => {
  const server = createServer((request, response) => {
    void route(request, response);
  });

  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const address = server.address();

  if (address === null || typeof address === "string") {
    throw new Error("server_address_invalid");
  }

  const response = await fetch(`http://127.0.0.1:${address.port}/health/details`);
  const payload = (await response.json()) as {
    service: string;
    status: string;
    dependencies: Array<{ name: string; status: string }>;
  };

  assert.equal(response.status, 200);
  assert.equal(payload.service, "engagement");
  assert.equal(payload.status, "ready");
  assert.ok(payload.dependencies.some((dependency) => dependency.name === "campaign-catalog"));
  assert.ok(payload.dependencies.some((dependency) => dependency.name === "touchpoints"));

  await new Promise<void>((resolve, reject) => server.close((error) => (error ? reject(error) : resolve())));
});
