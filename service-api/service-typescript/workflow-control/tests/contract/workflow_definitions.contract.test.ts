import assert from "node:assert/strict";
import { afterEach, test } from "node:test";
import { createServer, Server } from "node:http";
import { randomUUID } from "node:crypto";
import { AddressInfo } from "node:net";
import { route } from "../../src/api/router.js";

type WorkflowDefinitionResponse = {
  id: number;
  key: string;
  name: string;
  description: string | null;
  status: string;
  trigger: string;
};

type ErrorResponse = {
  code: string;
  message: string;
};

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

async function request(pathname: string, init?: RequestInit): Promise<Response> {
  const server = createServer((incoming, outgoing) => {
    void route(incoming, outgoing);
  });

  await new Promise<void>((resolve) => {
    server.listen(0, "127.0.0.1", () => resolve());
  });

  activeServers.push(server);

  const address = server.address() as AddressInfo;

  return fetch(`http://127.0.0.1:${address.port}${pathname}`, init);
}

test("workflow definitions contract should expose public fields on list", async () => {
  const response = await request("/api/workflow-control/definitions");
  const payload = await response.json() as WorkflowDefinitionResponse[];

  assert.equal(response.status, 200);
  assert.ok(payload.length > 0);

  for (const definition of payload) {
    assert.ok(definition.id > 0);
    assert.ok(definition.key.trim().length > 0);
    assert.ok(definition.name.trim().length > 0);
    assert.ok(definition.status.trim().length > 0);
    assert.ok(definition.trigger.trim().length > 0);
  }
});

test("workflow definition contract should return created resource", async () => {
  const key = `contract-${randomUUID()}`;
  const response = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Contract Flow",
      description: "Fluxo criado por contract test.",
      trigger: "lead.created"
    })
  });
  const payload = await response.json() as WorkflowDefinitionResponse;

  assert.equal(response.status, 201);
  assert.equal(payload.key, key);
  assert.equal(payload.name, "Contract Flow");
  assert.equal(payload.description, "Fluxo criado por contract test.");
  assert.equal(payload.status, "draft");
  assert.equal(payload.trigger, "lead.created");
});

test("workflow definition contract should expose detail and status update", async () => {
  const key = `contract-status-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Contract Status Flow",
      trigger: "lead.contacted"
    })
  });

  assert.equal(createResponse.status, 201);

  const updateResponse = await request(`/api/workflow-control/definitions/${key}/status`, {
    method: "PATCH",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      status: "active"
    })
  });
  const updated = await updateResponse.json() as WorkflowDefinitionResponse;

  assert.equal(updateResponse.status, 200);
  assert.equal(updated.key, key);
  assert.equal(updated.status, "active");

  const detailResponse = await request(`/api/workflow-control/definitions/${key}`);
  const detail = await detailResponse.json() as WorkflowDefinitionResponse;

  assert.equal(detailResponse.status, 200);
  assert.equal(detail.key, key);
  assert.equal(detail.status, "active");
});

test("workflow definition contract should expose not found error shape", async () => {
  const response = await request("/api/workflow-control/definitions/missing-contract-flow");
  const payload = await response.json() as ErrorResponse;

  assert.equal(response.status, 404);
  assert.equal(payload.code, "workflow_definition_not_found");
  assert.equal(payload.message, "Workflow definition was not found.");
});
