import assert from "node:assert/strict";
import { afterEach, test } from "node:test";
import { createServer, Server } from "node:http";
import { randomUUID } from "node:crypto";
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

test("workflow definitions list should expose bootstrap catalog", async () => {
  const response = await request("/api/workflow-control/definitions");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.ok(Array.isArray(payload));
  assert.ok(payload.some((definition) => definition.key === "lead-follow-up"));
});

test("workflow runs list should expose bootstrap execution ledger", async () => {
  const response = await request("/api/workflow-control/runs");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.ok(Array.isArray(payload));
  assert.equal(payload.length, 1);
  assert.equal(payload[0].publicId, "00000000-0000-0000-0000-000000000301");
  assert.equal(payload[0].status, "running");
  assert.equal(payload[0].triggerEvent, "lead.created");
  assert.equal(payload[0].subjectType, "crm.lead");
});

test("workflow run detail should expose bootstrap execution by public id", async () => {
  const response = await request("/api/workflow-control/runs/00000000-0000-0000-0000-000000000301");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.publicId, "00000000-0000-0000-0000-000000000301");
  assert.equal(payload.status, "running");
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.workflowDefinitionVersionId, 1);
});

test("workflow run create should append a pending execution linked to current version", async () => {
  const response = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: "00000000-0000-0000-0000-000000009999",
      initiatedBy: "ops-console"
    })
  });
  const payload = await response.json();

  assert.equal(response.status, 201);
  assert.match(payload.publicId, /^[0-9a-f-]{36}$/);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.workflowDefinitionVersionId, 1);
  assert.equal(payload.status, "pending");
  assert.equal(payload.triggerEvent, "lead.created");
  assert.equal(payload.initiatedBy, "ops-console");
});

test("workflow run summary should expose operational buckets", async () => {
  const listResponse = await request("/api/workflow-control/runs");
  const runs = await listResponse.json();
  const response = await request("/api/workflow-control/runs/summary");
  const payload = await response.json();

  assert.equal(listResponse.status, 200);
  assert.equal(response.status, 200);
  assert.equal(payload.total, runs.length);
  assert.equal(payload.pending, runs.filter((run) => run.status === "pending").length);
  assert.equal(payload.running, runs.filter((run) => run.status === "running").length);
  assert.equal(payload.completed, runs.filter((run) => run.status === "completed").length);
  assert.equal(payload.failed, runs.filter((run) => run.status === "failed").length);
  assert.equal(payload.cancelled, runs.filter((run) => run.status === "cancelled").length);
});

test("workflow run start should transition a pending execution to running", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: "00000000-0000-0000-0000-000000008888",
      initiatedBy: "ops-console"
    })
  });
  const created = await createResponse.json();

  assert.equal(createResponse.status, 201);

  const startResponse = await request(`/api/workflow-control/runs/${created.publicId}/start`, {
    method: "POST"
  });
  const payload = await startResponse.json();

  assert.equal(startResponse.status, 200);
  assert.equal(payload.publicId, created.publicId);
  assert.equal(payload.status, "running");
  assert.match(payload.startedAt, /.+/);
});

test("workflow run complete should close a running execution", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: "00000000-0000-0000-0000-000000007777",
      initiatedBy: "ops-console"
    })
  });
  const created = await createResponse.json();

  assert.equal(createResponse.status, 201);

  const startResponse = await request(`/api/workflow-control/runs/${created.publicId}/start`, {
    method: "POST"
  });
  assert.equal(startResponse.status, 200);

  const completeResponse = await request(`/api/workflow-control/runs/${created.publicId}/complete`, {
    method: "POST"
  });
  const payload = await completeResponse.json();

  assert.equal(completeResponse.status, 200);
  assert.equal(payload.publicId, created.publicId);
  assert.equal(payload.status, "completed");
  assert.match(payload.completedAt, /.+/);
});

test("workflow definition detail should expose created resource by key", async () => {
  const key = `detail-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Detail Flow",
      description: "Fluxo criado para teste unitario de leitura por chave.",
      trigger: "lead.qualified"
    })
  });

  assert.equal(createResponse.status, 201);

  const detailResponse = await request(`/api/workflow-control/definitions/${key}`);
  const detailPayload = await detailResponse.json();

  assert.equal(detailResponse.status, 200);
  assert.equal(detailPayload.key, key);
  assert.equal(detailPayload.name, "Detail Flow");
  assert.equal(detailPayload.trigger, "lead.qualified");
});

test("workflow definition versions should expose bootstrap publication history", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.ok(Array.isArray(payload));
  assert.equal(payload.length, 1);
  assert.equal(payload[0].workflowDefinitionId, 1);
  assert.equal(payload[0].versionNumber, 1);
  assert.equal(payload[0].snapshotTrigger, "lead.created");
});

test("workflow definition publish should create a new version snapshot", async () => {
  const key = `publish-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Publish Flow",
      description: "Fluxo para testar publicacao manual de versao.",
      trigger: "lead.created"
    })
  });

  assert.equal(createResponse.status, 201);

  const publishResponse = await request(`/api/workflow-control/definitions/${key}/versions`, {
    method: "POST"
  });
  const publishPayload = await publishResponse.json();

  assert.equal(publishResponse.status, 201);
  assert.ok(publishPayload.workflowDefinitionId > 1);
  assert.equal(publishPayload.versionNumber, 1);
  assert.equal(publishPayload.snapshotName, "Publish Flow");
  assert.equal(publishPayload.snapshotTrigger, "lead.created");
});

test("workflow definition current version should expose latest snapshot", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions/current");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.versionNumber, 1);
  assert.equal(payload.snapshotName, "Lead Follow-Up");
});

test("workflow definition version detail should expose requested snapshot", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions/1");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.versionNumber, 1);
  assert.equal(payload.snapshotTrigger, "lead.created");
});

test("workflow definition version summary should expose totals and current version", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions/summary");
  const payload = await response.json();

  assert.equal(response.status, 200);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.totalVersions, 1);
  assert.equal(payload.currentVersionNumber, 1);
  assert.equal(payload.currentSnapshotStatus, "active");
});

test("workflow definition restore should bring metadata back from a published version", async () => {
  const key = `restore-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Restore Flow",
      description: "Fluxo inicial para restore.",
      trigger: "lead.created"
    })
  });

  assert.equal(createResponse.status, 201);

  const publishV1Response = await request(`/api/workflow-control/definitions/${key}/versions`, {
    method: "POST"
  });
  assert.equal(publishV1Response.status, 201);

  const updateResponse = await request(`/api/workflow-control/definitions/${key}`, {
    method: "PATCH",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      name: "Restore Flow Prime",
      description: "Fluxo refinado antes do rollback.",
      trigger: "lead.qualified"
    })
  });
  assert.equal(updateResponse.status, 200);

  const statusResponse = await request(`/api/workflow-control/definitions/${key}/status`, {
    method: "PATCH",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      status: "active"
    })
  });
  assert.equal(statusResponse.status, 200);

  const publishV2Response = await request(`/api/workflow-control/definitions/${key}/versions`, {
    method: "POST"
  });
  assert.equal(publishV2Response.status, 201);

  const restoreResponse = await request(`/api/workflow-control/definitions/${key}/versions/1/restore`, {
    method: "POST"
  });
  const restorePayload = await restoreResponse.json();

  assert.equal(restoreResponse.status, 200);
  assert.equal(restorePayload.name, "Restore Flow");
  assert.equal(restorePayload.description, "Fluxo inicial para restore.");
  assert.equal(restorePayload.trigger, "lead.created");
  assert.equal(restorePayload.status, "draft");
});

test("workflow definition create should normalize payload and return draft status", async () => {
  const key = `create-${randomUUID()}`;
  const response = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key: `  ${key.toUpperCase()}  `,
      name: "  Create Flow  ",
      description: "  Fluxo criado para validar normalizacao.  ",
      trigger: "  LEAD.CREATED  "
    })
  });
  const payload = await response.json();

  assert.equal(response.status, 201);
  assert.equal(payload.key, key);
  assert.equal(payload.name, "Create Flow");
  assert.equal(payload.description, "Fluxo criado para validar normalizacao.");
  assert.equal(payload.trigger, "lead.created");
  assert.equal(payload.status, "draft");
});

test("workflow definition status should update created resource", async () => {
  const key = `status-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Status Flow",
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
  const updatePayload = await updateResponse.json();

  assert.equal(updateResponse.status, 200);
  assert.equal(updatePayload.key, key);
  assert.equal(updatePayload.status, "active");
});

test("workflow definition update should persist metadata changes", async () => {
  const key = `profile-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Profile Flow",
      description: "Descricao inicial do fluxo.",
      trigger: "lead.created"
    })
  });

  assert.equal(createResponse.status, 201);

  const updateResponse = await request(`/api/workflow-control/definitions/${key}`, {
    method: "PATCH",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      name: "Profile Flow Prime",
      description: "Descricao refinada para operacao comercial.",
      trigger: "lead.qualified"
    })
  });
  const updatePayload = await updateResponse.json();

  assert.equal(updateResponse.status, 200);
  assert.equal(updatePayload.key, key);
  assert.equal(updatePayload.name, "Profile Flow Prime");
  assert.equal(updatePayload.description, "Descricao refinada para operacao comercial.");
  assert.equal(updatePayload.trigger, "lead.qualified");
});

test("workflow definition update should reject empty payload", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up", {
    method: "PATCH",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({})
  });
  const payload = await response.json();

  assert.equal(response.status, 400);
  assert.equal(payload.code, "workflow_definition_update_required");
});

test("workflow definition status should reject invalid transitions payload", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/status", {
    method: "PATCH",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      status: "paused"
    })
  });
  const payload = await response.json();

  assert.equal(response.status, 400);
  assert.equal(payload.code, "workflow_definition_status_invalid");
});
