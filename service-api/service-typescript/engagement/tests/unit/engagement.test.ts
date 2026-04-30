import assert from "node:assert/strict";
import test from "node:test";
import { services } from "../../src/config/container.js";

test("campaign catalog exposes bootstrap campaign and supports creation", async () => {
  const campaigns = await services.listCampaigns.execute({ tenantSlug: "bootstrap-ops" });

  assert.equal(campaigns.length, 2);
  assert.equal(campaigns[0]?.key, "lead-follow-up-campaign");

  const created = await services.createCampaign.execute({
    tenantSlug: "bootstrap-ops",
    key: "reactivation-wave",
    name: "Reactivation Wave",
    description: "Reengaja leads mornos com sequencia curta.",
    channel: "email",
    status: "draft",
    touchpointGoal: "reactivate-lead",
    workflowDefinitionKey: "lead-follow-up",
    budgetCents: 42000
  });

  assert.equal(created.key, "reactivation-wave");
  assert.equal(created.channel, "email");
});

test("template catalog exposes bootstrap template and supports creation", async () => {
  const templates = await services.listTemplates.execute({ tenantSlug: "bootstrap-ops" });

  assert.equal(templates.length, 1);
  assert.equal(templates[0]?.key, "lead-follow-up-whatsapp");

  const created = await services.createTemplate.execute({
    tenantSlug: "bootstrap-ops",
    key: "proposal-reminder-email",
    name: "Proposal Reminder Email",
    channel: "email",
    provider: "resend",
    status: "draft",
    subject: "Sua proposta segue aberta",
    body: "Ola {{firstName}}, seguimos disponiveis para concluir sua proposta."
  });

  assert.equal(created.provider, "resend");
  assert.equal(created.status, "draft");
});

test("touchpoint stream supports creation, status update and summary", async () => {
  const created = await services.createTouchpoint.execute({
    tenantSlug: "bootstrap-ops",
    campaignPublicId: "00000000-0000-0000-0000-00000000c101",
    leadPublicId: "00000000-0000-0000-0000-000000008851",
    contactValue: "+5531888888888",
    source: "manual",
    createdBy: "unit-test",
    notes: "Primeiro contato do teste."
  });

  assert.equal(created.status, "queued");
  assert.equal(created.workflowDefinitionKey, "lead-follow-up");

  const updated = await services.updateTouchpointStatus.execute(
    created.publicId,
    "converted",
    "00000000-0000-0000-0000-000000000399"
  );

  assert.equal(updated?.status, "converted");
  assert.equal(updated?.lastWorkflowRunPublicId, "00000000-0000-0000-0000-000000000399");

  const summary = await services.getTouchpointSummary.execute({ tenantSlug: "bootstrap-ops" });

  assert.equal(summary.totals.touchpoints, 2);
  assert.equal(summary.totals.workflowDispatched, 2);
  assert.equal(summary.byStatus.converted, 1);
  assert.equal(summary.byChannel.whatsapp, 2);
});

test("delivery stream supports creation, status update and summary", async () => {
  const created = await services.createTouchpoint.execute({
    tenantSlug: "bootstrap-ops",
    campaignPublicId: "00000000-0000-0000-0000-00000000c102",
    leadPublicId: "00000000-0000-0000-0000-000000008812",
    contactValue: "lead@example.com",
    source: "crm",
    createdBy: "unit-test",
    notes: "Touchpoint para entrega do teste."
  });

  const delivery = await services.createTouchpointDelivery.execute(created.publicId, {
    tenantSlug: "bootstrap-ops",
    templatePublicId: null,
    provider: "resend",
    providerMessageId: "msg-unit-001",
    sentBy: "unit-test",
    notes: "Entrega inicial do teste."
  });

  assert.equal(delivery.status, "sent");
  assert.equal(delivery.provider, "resend");

  const updated = await services.updateTouchpointDeliveryStatus.execute(delivery.publicId, {
    status: "delivered",
    providerMessageId: "msg-unit-001",
    notes: "Confirmado pelo provider."
  });

  assert.equal(updated?.status, "delivered");

  const summary = await services.getDeliverySummary.execute({ tenantSlug: "bootstrap-ops" });

  assert.equal(summary.totals.templates, 2);
  assert.equal(summary.totals.deliveries, 2);
  assert.equal(summary.byProvider.manual, 1);
  assert.equal(summary.byProvider.resend, 1);
  assert.equal(summary.byStatus.delivered, 2);
});
