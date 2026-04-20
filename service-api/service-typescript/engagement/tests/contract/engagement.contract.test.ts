import assert from "node:assert/strict";
import test from "node:test";
import { createServer } from "node:http";
import { route } from "../../src/api/router.js";

async function withServer(run: (baseUrl: string) => Promise<void>): Promise<void> {
  const server = createServer((request, response) => {
    void route(request, response);
  });

  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const address = server.address();

  if (address === null || typeof address === "string") {
    throw new Error("server_address_invalid");
  }

  try {
    await run(`http://127.0.0.1:${address.port}`);
  } finally {
    await new Promise<void>((resolve, reject) => server.close((error) => (error ? reject(error) : resolve())));
  }
}

test("campaign routes expose catalog and status transition", async () => {
  await withServer(async (baseUrl) => {
    const listResponse = await fetch(`${baseUrl}/api/engagement/campaigns?tenantSlug=bootstrap-ops`);
    const listPayload = (await listResponse.json()) as Array<{ publicId: string; key: string }>;

    assert.equal(listResponse.status, 200);
    assert.equal(listPayload.length, 2);

    const createResponse = await fetch(`${baseUrl}/api/engagement/campaigns`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        key: "whatsapp-reengagement",
        name: "WhatsApp Reengagement",
        description: "Retoma conversas paradas via WhatsApp.",
        channel: "whatsapp",
        touchpointGoal: "resume-thread",
        workflowDefinitionKey: "lead-follow-up",
        budgetCents: 28000
      })
    });
    const createdCampaign = (await createResponse.json()) as { publicId: string; status: string };

    assert.equal(createResponse.status, 201);
    assert.equal(createdCampaign.status, "draft");

    const statusResponse = await fetch(`${baseUrl}/api/engagement/campaigns/${createdCampaign.publicId}/status`, {
      method: "PATCH",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ status: "active" })
    });
    const updatedCampaign = (await statusResponse.json()) as { status: string };

    assert.equal(statusResponse.status, 200);
    assert.equal(updatedCampaign.status, "active");
  });
});

test("touchpoint routes expose list, creation and summary", async () => {
  await withServer(async (baseUrl) => {
    const createResponse = await fetch(`${baseUrl}/api/engagement/touchpoints`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        campaignPublicId: "00000000-0000-0000-0000-00000000c101",
        leadPublicId: "00000000-0000-0000-0000-000000008899",
        contactValue: "+5531777777777",
        source: "crm",
        createdBy: "contract-test",
        notes: "Teste publico do contrato."
      })
    });
    const createdTouchpoint = (await createResponse.json()) as { publicId: string; status: string };

    assert.equal(createResponse.status, 201);
    assert.equal(createdTouchpoint.status, "queued");

    const statusResponse = await fetch(`${baseUrl}/api/engagement/touchpoints/${createdTouchpoint.publicId}/status`, {
      method: "PATCH",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        status: "responded",
        lastWorkflowRunPublicId: "00000000-0000-0000-0000-000000000355"
      })
    });
    const updatedTouchpoint = (await statusResponse.json()) as { status: string; lastWorkflowRunPublicId: string };

    assert.equal(statusResponse.status, 200);
    assert.equal(updatedTouchpoint.status, "responded");
    assert.equal(updatedTouchpoint.lastWorkflowRunPublicId, "00000000-0000-0000-0000-000000000355");

    const summaryResponse = await fetch(`${baseUrl}/api/engagement/touchpoints/summary?tenantSlug=bootstrap-ops`);
    const summaryPayload = (await summaryResponse.json()) as {
      totals: { touchpoints: number; workflowDispatched: number };
      byStatus: { responded: number };
    };

    assert.equal(summaryResponse.status, 200);
    assert.equal(summaryPayload.totals.touchpoints, 2);
    assert.equal(summaryPayload.totals.workflowDispatched, 2);
    assert.equal(summaryPayload.byStatus.responded, 2);
  });
});
