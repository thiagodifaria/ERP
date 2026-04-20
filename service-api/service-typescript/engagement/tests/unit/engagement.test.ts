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
