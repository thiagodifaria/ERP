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

type WorkflowDefinitionVersionResponse = {
  id: number;
  workflowDefinitionId: number;
  versionNumber: number;
  snapshotName: string;
  snapshotDescription: string | null;
  snapshotStatus: string;
  snapshotTrigger: string;
};

type WorkflowRunResponse = {
  id: number;
  publicId: string;
  workflowDefinitionId: number;
  workflowDefinitionVersionId: number;
  status: string;
  triggerEvent: string;
  subjectType: string;
  subjectPublicId: string;
  initiatedBy: string;
  startedAt: string | null;
  completedAt: string | null;
  failedAt: string | null;
  cancelledAt: string | null;
};

type WorkflowRunEventResponse = {
  id: number;
  publicId: string;
  workflowRunPublicId: string;
  category: string;
  body: string;
  createdBy: string;
  createdAt: string;
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

test("workflow runs contract should expose public fields on list", async () => {
  const response = await request("/api/workflow-control/runs");
  const payload = await response.json() as WorkflowRunResponse[];

  assert.equal(response.status, 200);
  assert.ok(payload.length > 0);

  for (const workflowRun of payload) {
    assert.ok(workflowRun.id > 0);
    assert.match(workflowRun.publicId, /^[0-9a-f-]{36}$/);
    assert.ok(workflowRun.workflowDefinitionId > 0);
    assert.ok(workflowRun.workflowDefinitionVersionId > 0);
    assert.ok(workflowRun.status.trim().length > 0);
    assert.ok(workflowRun.triggerEvent.trim().length > 0);
    assert.ok(workflowRun.subjectType.trim().length > 0);
    assert.match(workflowRun.subjectPublicId, /^[0-9a-f-]{36}$/);
    assert.ok(workflowRun.initiatedBy.trim().length > 0);
  }
});

test("workflow run contract should expose execution detail by public id", async () => {
  const response = await request("/api/workflow-control/runs/00000000-0000-0000-0000-000000000301");
  const payload = await response.json() as WorkflowRunResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.publicId, "00000000-0000-0000-0000-000000000301");
  assert.equal(payload.status, "running");
  assert.equal(payload.triggerEvent, "lead.created");
  assert.equal(payload.subjectType, "crm.lead");
});

test("workflow run events contract should expose event history by run", async () => {
  const response = await request("/api/workflow-control/runs/00000000-0000-0000-0000-000000000301/events");
  const payload = await response.json() as WorkflowRunEventResponse[];

  assert.equal(response.status, 200);
  assert.ok(payload.length > 0);
  assert.equal(payload[0].workflowRunPublicId, "00000000-0000-0000-0000-000000000301");
  assert.equal(payload[0].category, "note");
  assert.ok(payload[0].body.trim().length > 0);
  assert.ok(payload[0].createdBy.trim().length > 0);
  assert.ok(payload[0].createdAt.trim().length > 0);
});

test("workflow run events contract should return created note resource", async () => {
  const createRunResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-notes"
    })
  });
  const createdRun = await createRunResponse.json() as WorkflowRunResponse;

  assert.equal(createRunResponse.status, 201);

  const response = await request(`/api/workflow-control/runs/${createdRun.publicId}/events`, {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      body: "Observacao contratual registrada no ledger da execucao.",
      createdBy: "contract-notes"
    })
  });
  const payload = await response.json() as WorkflowRunEventResponse;

  assert.equal(response.status, 201);
  assert.match(payload.publicId, /^[0-9a-f-]{36}$/);
  assert.equal(payload.workflowRunPublicId, createdRun.publicId);
  assert.equal(payload.category, "note");
  assert.equal(payload.body, "Observacao contratual registrada no ledger da execucao.");
  assert.equal(payload.createdBy, "contract-notes");
  assert.ok(payload.createdAt.trim().length > 0);
});

test("workflow run contract should return created execution resource", async () => {
  const response = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-suite"
    })
  });
  const payload = await response.json() as WorkflowRunResponse;

  assert.equal(response.status, 201);
  assert.match(payload.publicId, /^[0-9a-f-]{36}$/);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.workflowDefinitionVersionId, 1);
  assert.equal(payload.status, "pending");
  assert.equal(payload.triggerEvent, "lead.created");
  assert.equal(payload.subjectType, "crm.lead");
  assert.equal(payload.initiatedBy, "contract-suite");
});

test("workflow run summary contract should expose operational buckets", async () => {
  const response = await request("/api/workflow-control/runs/summary");
  const payload = await response.json() as {
    total: number;
    pending: number;
    running: number;
    completed: number;
    failed: number;
    cancelled: number;
  };

  assert.equal(response.status, 200);
  assert.ok(payload.total >= 1);
  assert.ok(payload.running >= 1);
  assert.ok(payload.pending >= 0);
  assert.ok(payload.completed >= 0);
  assert.ok(payload.failed >= 0);
  assert.ok(payload.cancelled >= 0);
});

test("workflow runs contract should support operational filters", async () => {
  const subjectPublicId = randomUUID();
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId,
      initiatedBy: "contract-filter"
    })
  });

  assert.equal(createResponse.status, 201);

  const response = await request("/api/workflow-control/runs?status=pending&workflowDefinitionKey=lead-follow-up&subjectType=crm.lead&initiatedBy=contract-filter");
  const payload = await response.json() as WorkflowRunResponse[];

  assert.equal(response.status, 200);
  assert.ok(payload.length >= 1);
  assert.ok(payload.every((workflowRun) => workflowRun.status === "pending"));
  assert.ok(payload.every((workflowRun) => workflowRun.subjectType === "crm.lead"));
  assert.ok(payload.every((workflowRun) => workflowRun.initiatedBy === "contract-filter"));
});

test("workflow run contract should expose start transition", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-start"
    })
  });
  const created = await createResponse.json() as WorkflowRunResponse;

  assert.equal(createResponse.status, 201);

  const response = await request(`/api/workflow-control/runs/${created.publicId}/start`, {
    method: "POST"
  });
  const payload = await response.json() as WorkflowRunResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.publicId, created.publicId);
  assert.equal(payload.status, "running");
  assert.ok(payload.startedAt);
});

test("workflow run contract should expose status events generated by lifecycle transitions", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-timeline"
    })
  });
  const created = await createResponse.json() as WorkflowRunResponse;

  assert.equal(createResponse.status, 201);

  const startResponse = await request(`/api/workflow-control/runs/${created.publicId}/start`, {
    method: "POST"
  });
  assert.equal(startResponse.status, 200);

  const completeResponse = await request(`/api/workflow-control/runs/${created.publicId}/complete`, {
    method: "POST"
  });
  assert.equal(completeResponse.status, 200);

  const eventsResponse = await request(`/api/workflow-control/runs/${created.publicId}/events`);
  const eventsPayload = await eventsResponse.json() as WorkflowRunEventResponse[];

  assert.equal(eventsResponse.status, 200);
  assert.equal(eventsPayload.length, 2);
  assert.equal(eventsPayload[0].category, "status");
  assert.equal(eventsPayload[0].body, "Workflow run moved to running.");
  assert.equal(eventsPayload[0].createdBy, "workflow-control");
  assert.equal(eventsPayload[1].body, "Workflow run moved to completed.");
});

test("workflow run contract should expose complete transition", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-complete"
    })
  });
  const created = await createResponse.json() as WorkflowRunResponse;

  assert.equal(createResponse.status, 201);

  const startResponse = await request(`/api/workflow-control/runs/${created.publicId}/start`, {
    method: "POST"
  });
  assert.equal(startResponse.status, 200);

  const response = await request(`/api/workflow-control/runs/${created.publicId}/complete`, {
    method: "POST"
  });
  const payload = await response.json() as WorkflowRunResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.publicId, created.publicId);
  assert.equal(payload.status, "completed");
  assert.ok(payload.completedAt);
});

test("workflow run contract should expose fail transition", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-fail"
    })
  });
  const created = await createResponse.json() as WorkflowRunResponse;

  assert.equal(createResponse.status, 201);

  const startResponse = await request(`/api/workflow-control/runs/${created.publicId}/start`, {
    method: "POST"
  });
  assert.equal(startResponse.status, 200);

  const response = await request(`/api/workflow-control/runs/${created.publicId}/fail`, {
    method: "POST"
  });
  const payload = await response.json() as WorkflowRunResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.publicId, created.publicId);
  assert.equal(payload.status, "failed");
  assert.ok(payload.failedAt);
});

test("workflow run contract should expose cancel transition", async () => {
  const createResponse = await request("/api/workflow-control/runs", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      workflowDefinitionKey: "lead-follow-up",
      subjectType: "crm.lead",
      subjectPublicId: randomUUID(),
      initiatedBy: "contract-cancel"
    })
  });
  const created = await createResponse.json() as WorkflowRunResponse;

  assert.equal(createResponse.status, 201);

  const response = await request(`/api/workflow-control/runs/${created.publicId}/cancel`, {
    method: "POST"
  });
  const payload = await response.json() as WorkflowRunResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.publicId, created.publicId);
  assert.equal(payload.status, "cancelled");
  assert.ok(payload.cancelledAt);
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

test("workflow definition contract should return updated metadata resource", async () => {
  const key = `contract-update-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Contract Update Flow",
      description: "Fluxo inicial para contract de update.",
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
      name: "Contract Update Flow Prime",
      description: "Fluxo revisado pelo contract test.",
      trigger: "lead.qualified"
    })
  });
  const payload = await updateResponse.json() as WorkflowDefinitionResponse;

  assert.equal(updateResponse.status, 200);
  assert.equal(payload.key, key);
  assert.equal(payload.name, "Contract Update Flow Prime");
  assert.equal(payload.description, "Fluxo revisado pelo contract test.");
  assert.equal(payload.trigger, "lead.qualified");
});

test("workflow definition versions contract should expose version history", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions");
  const payload = await response.json() as WorkflowDefinitionVersionResponse[];

  assert.equal(response.status, 200);
  assert.ok(payload.length > 0);
  assert.equal(payload[0].versionNumber, 1);
});

test("workflow definition publish contract should return created version resource", async () => {
  const key = `contract-publish-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Contract Publish Flow",
      description: "Fluxo publicado via contract test.",
      trigger: "lead.created"
    })
  });

  assert.equal(createResponse.status, 201);

  const publishResponse = await request(`/api/workflow-control/definitions/${key}/versions`, {
    method: "POST"
  });
  const payload = await publishResponse.json() as WorkflowDefinitionVersionResponse;

  assert.equal(publishResponse.status, 201);
  assert.ok(payload.id > 0);
  assert.ok(payload.workflowDefinitionId > 0);
  assert.equal(payload.versionNumber, 1);
  assert.equal(payload.snapshotName, "Contract Publish Flow");
  assert.equal(payload.snapshotTrigger, "lead.created");
});

test("workflow definition current version contract should expose latest version snapshot", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions/current");
  const payload = await response.json() as WorkflowDefinitionVersionResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.versionNumber, 1);
  assert.equal(payload.snapshotName, "Lead Follow-Up");
});

test("workflow definition version detail contract should expose requested snapshot", async () => {
  const response = await request("/api/workflow-control/definitions/lead-follow-up/versions/1");
  const payload = await response.json() as WorkflowDefinitionVersionResponse;

  assert.equal(response.status, 200);
  assert.equal(payload.workflowDefinitionId, 1);
  assert.equal(payload.versionNumber, 1);
  assert.equal(payload.snapshotTrigger, "lead.created");
});

test("workflow definition restore contract should return restored definition resource", async () => {
  const key = `contract-restore-${randomUUID()}`;
  const createResponse = await request("/api/workflow-control/definitions", {
    method: "POST",
    headers: {
      "content-type": "application/json"
    },
    body: JSON.stringify({
      key,
      name: "Contract Restore Flow",
      description: "Fluxo inicial para restore contratual.",
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
      name: "Contract Restore Flow Prime",
      description: "Fluxo refinado antes do restore.",
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

  const restoreResponse = await request(`/api/workflow-control/definitions/${key}/versions/1/restore`, {
    method: "POST"
  });
  const payload = await restoreResponse.json() as WorkflowDefinitionResponse;

  assert.equal(restoreResponse.status, 200);
  assert.equal(payload.key, key);
  assert.equal(payload.name, "Contract Restore Flow");
  assert.equal(payload.description, "Fluxo inicial para restore contratual.");
  assert.equal(payload.trigger, "lead.created");
  assert.equal(payload.status, "draft");
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
