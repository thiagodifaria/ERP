import assert from "node:assert/strict";
import { afterEach, test } from "node:test";
import { createServer, Server } from "node:http";
import { AddressInfo } from "node:net";
import { route } from "../../src/api/router.js";

const activeServers: Server[] = [];

afterEach(async () => {
  while (activeServers.length > 0) {
    const server = activeServers.pop();

    if (server) {
      await new Promise<void>((resolve, reject) => {
        server.close((error) => {
          if (error) {
            reject(error);
            return;
          }

          resolve();
        });
      });
    }
  }
});

async function request(pathname: string): Promise<Response> {
  const server = createServer(route);

  await new Promise<void>((resolve) => {
    server.listen(0, "127.0.0.1", () => resolve());
  });

  activeServers.push(server);

  const address = server.address() as AddressInfo;

  return fetch(`http://127.0.0.1:${address.port}${pathname}`);
}

test("health live should return live status", async () => {
  const response = await request("/health/live");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.service, "workflow-control");
  assert.equal(payload.status, "live");
});

test("health ready should return ready status", async () => {
  const response = await request("/health/ready");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.service, "workflow-control");
  assert.equal(payload.status, "ready");
});

test("health details should expose readiness dependencies", async () => {
  const response = await request("/health/details");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.service, "workflow-control");
  assert.equal(payload.status, "ready");
  assert.deepEqual(payload.dependencies, [
    { name: "router", status: "ready" },
    { name: "definitions-catalog", status: "ready" },
    { name: "workflow-runs", status: "ready" }
  ]);
});
