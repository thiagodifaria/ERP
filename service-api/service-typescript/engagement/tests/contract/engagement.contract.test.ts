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

test("template routes expose catalog and status transition", async () => {
  await withServer(async (baseUrl) => {
    const listResponse = await fetch(`${baseUrl}/api/engagement/templates?tenantSlug=bootstrap-ops`);
    const listPayload = (await listResponse.json()) as Array<{ publicId: string; key: string }>;

    assert.equal(listResponse.status, 200);
    assert.equal(listPayload.length, 1);

    const createResponse = await fetch(`${baseUrl}/api/engagement/templates`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        key: "reactivation-email-template",
        name: "Reactivation Email Template",
        channel: "email",
        provider: "resend",
        subject: "Ainda faz sentido retomarmos?",
        body: "Ola {{firstName}}, retomamos seu contexto para seguir daqui."
      })
    });
    const createdTemplate = (await createResponse.json()) as { publicId: string; status: string };

    assert.equal(createResponse.status, 201);
    assert.equal(createdTemplate.status, "draft");

    const statusResponse = await fetch(`${baseUrl}/api/engagement/templates/${createdTemplate.publicId}/status`, {
      method: "PATCH",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ status: "active" })
    });
    const updatedTemplate = (await statusResponse.json()) as { status: string };

    assert.equal(statusResponse.status, 200);
    assert.equal(updatedTemplate.status, "active");
  });
});

test("touchpoint routes expose list, creation and summary", async () => {
  await withServer(async (baseUrl) => {
    const createResponse = await fetch(`${baseUrl}/api/engagement/touchpoints`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        campaignPublicId: "00000000-0000-0000-0000-00000000c102",
        leadPublicId: "00000000-0000-0000-0000-000000008851",
        contactValue: "+5531777777777",
        source: "crm",
        createdBy: "contract-test",
        notes: "Teste publico do contrato."
      })
    });
    const createdTouchpoint = (await createResponse.json()) as {
      publicId: string;
      status: string;
      businessEntityType: string;
      businessEntityPublicId: string;
    };

    assert.equal(createResponse.status, 201);
    assert.equal(createdTouchpoint.status, "queued");
    assert.equal(createdTouchpoint.businessEntityType, "crm.lead");
    assert.equal(createdTouchpoint.businessEntityPublicId, "00000000-0000-0000-0000-000000008851");

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
      totals: { touchpoints: number; workflowDispatched: number; businessLinked: number };
      byStatus: { responded: number };
    };

    assert.equal(summaryResponse.status, 200);
    assert.equal(summaryPayload.totals.touchpoints, 2);
    assert.equal(summaryPayload.totals.workflowDispatched, 2);
    assert.equal(summaryPayload.totals.businessLinked, 2);
    assert.equal(summaryPayload.byStatus.responded, 2);
  });
});

test("delivery routes expose creation, status transition and summary", async () => {
  await withServer(async (baseUrl) => {
    const createTouchpointResponse = await fetch(`${baseUrl}/api/engagement/touchpoints`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        campaignPublicId: "00000000-0000-0000-0000-00000000c101",
        leadPublicId: "00000000-0000-0000-0000-000000008823",
        contactValue: "+5531777777777",
        source: "crm",
        createdBy: "contract-test",
        notes: "Teste publico de entregas."
      })
    });
    const createdTouchpoint = (await createTouchpointResponse.json()) as { publicId: string };
    assert.equal(createTouchpointResponse.status, 201);

    const createDeliveryResponse = await fetch(
      `${baseUrl}/api/engagement/touchpoints/${createdTouchpoint.publicId}/deliveries`,
      {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          tenantSlug: "bootstrap-ops",
          provider: "resend",
          providerMessageId: "msg-contract-001",
          sentBy: "contract-test",
          notes: "Entrega inicial do contrato."
        })
      }
    );
    const createdDelivery = (await createDeliveryResponse.json()) as { publicId: string; status: string };
    assert.equal(createDeliveryResponse.status, 201);
    assert.equal(createdDelivery.status, "sent");

    const listDeliveriesResponse = await fetch(
      `${baseUrl}/api/engagement/touchpoints/${createdTouchpoint.publicId}/deliveries?tenantSlug=bootstrap-ops`
    );
    const listDeliveriesPayload = (await listDeliveriesResponse.json()) as Array<{ publicId: string }>;
    assert.equal(listDeliveriesResponse.status, 200);
    assert.equal(listDeliveriesPayload.length, 1);

    const statusResponse = await fetch(
      `${baseUrl}/api/engagement/touchpoints/${createdTouchpoint.publicId}/deliveries/${createdDelivery.publicId}/status`,
      {
        method: "PATCH",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ status: "delivered", providerMessageId: "msg-contract-001" })
      }
    );
    const updatedDelivery = (await statusResponse.json()) as { status: string };
    assert.equal(statusResponse.status, 200);
    assert.equal(updatedDelivery.status, "delivered");

    const summaryResponse = await fetch(`${baseUrl}/api/engagement/deliveries/summary?tenantSlug=bootstrap-ops`);
    const summaryPayload = (await summaryResponse.json()) as {
      totals: { deliveries: number; templates: number };
      byStatus: { delivered: number };
    };

    assert.equal(summaryResponse.status, 200);
    assert.equal(summaryPayload.totals.templates, 2);
    assert.equal(summaryPayload.totals.deliveries, 2);
    assert.equal(summaryPayload.byStatus.delivered, 2);
  });
});

test("provider routes expose callback detail and traceability", async () => {
  await withServer(async (baseUrl) => {
    const providerResponse = await fetch(`${baseUrl}/api/engagement/providers/meta-ads`);
    const providerPayload = (await providerResponse.json()) as { provider: string; mode: string; credentialKey: string };
    assert.equal(providerResponse.status, 200);
    assert.equal(providerPayload.provider, "meta_ads");
    assert.equal(providerPayload.mode, "fallback");
    assert.equal(providerPayload.credentialKey, "ENGAGEMENT_META_ADS_ACCESS_TOKEN");

    const templateResponse = await fetch(`${baseUrl}/api/engagement/templates`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        key: "provider-trace-email",
        name: "Provider Trace Email",
        channel: "email",
        provider: "resend",
        status: "active",
        subject: "Trace contract route",
        body: "Tracking callback flow in contract coverage."
      })
    });
    const templatePayload = (await templateResponse.json()) as { publicId: string };
    assert.equal(templateResponse.status, 201);

    const touchpointResponse = await fetch(`${baseUrl}/api/engagement/touchpoints`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        campaignPublicId: "00000000-0000-0000-0000-00000000c102",
        leadPublicId: "00000000-0000-0000-0000-000000008851",
        contactValue: "contract.trace@example.com",
        source: "workflow",
        createdBy: "contract-test",
        notes: "Provider trace touchpoint."
      })
    });
    const touchpointPayload = (await touchpointResponse.json()) as { publicId: string };
    assert.equal(touchpointResponse.status, 201);

    const deliveryResponse = await fetch(
      `${baseUrl}/api/engagement/touchpoints/${touchpointPayload.publicId}/deliveries`,
      {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          tenantSlug: "bootstrap-ops",
          templatePublicId: templatePayload.publicId,
          provider: "resend",
          providerMessageId: "msg-contract-trace-001",
          sentBy: "contract-test",
          notes: "Trace callback delivery."
        })
      }
    );
    const deliveryPayload = (await deliveryResponse.json()) as { publicId: string };
    assert.equal(deliveryResponse.status, 201);

    const callbackResponse = await fetch(`${baseUrl}/api/engagement/providers/resend/events`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        tenantSlug: "bootstrap-ops",
        eventType: "delivery.responded",
        externalEventId: "contract-provider-event-001",
        touchpointPublicId: touchpointPayload.publicId,
        deliveryPublicId: deliveryPayload.publicId,
        workflowRunPublicId: "00000000-0000-0000-0000-000000000455",
        leadPublicId: "00000000-0000-0000-0000-000000008851"
      })
    });
    const callbackPayload = (await callbackResponse.json()) as { publicId: string; eventType: string; status: string };

    assert.equal(callbackResponse.status, 201);
    assert.equal(callbackPayload.eventType, "delivery.responded");
    assert.equal(callbackPayload.status, "processed");

    const detailResponse = await fetch(`${baseUrl}/api/engagement/provider-events/${callbackPayload.publicId}`);
    const detailPayload = (await detailResponse.json()) as {
      publicId: string;
      eventType: string;
      responseSummary: string;
      businessEntityType: string;
      businessEntityPublicId: string;
    };

    assert.equal(detailResponse.status, 200);
    assert.equal(detailPayload.publicId, callbackPayload.publicId);
    assert.equal(detailPayload.eventType, "delivery.responded");
    assert.equal(detailPayload.businessEntityType, "crm.lead");
    assert.equal(detailPayload.businessEntityPublicId, "00000000-0000-0000-0000-000000008851");
    assert.ok(detailPayload.responseSummary.length > 0);
  });
});
