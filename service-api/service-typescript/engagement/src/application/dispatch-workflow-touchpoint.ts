import { CampaignRepository } from "../domain/campaign-repository.js";
import { DeliveryRepository } from "../domain/delivery-repository.js";
import { ensureDeliveryInput, ensureDeliveryProvider } from "../domain/delivery.js";
import { ProviderEventRepository } from "../domain/provider-event-repository.js";
import { ensureProvider } from "../domain/provider-event.js";
import { TemplateRepository } from "../domain/template-repository.js";
import { ensureOptionalPublicId, ensurePublicId, ensureTouchpointText } from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export type DispatchWorkflowTouchpointInput = {
  tenantSlug: string;
  campaignPublicId: string;
  templatePublicId?: string | null;
  leadPublicId: string;
  contactValue: string;
  provider: string;
  workflowRunPublicId: string;
  providerMessageId?: string | null;
  createdBy: string;
  notes?: string;
};

export class DispatchWorkflowTouchpoint {
  constructor(
    private readonly campaigns: CampaignRepository,
    private readonly templates: TemplateRepository,
    private readonly touchpoints: TouchpointRepository,
    private readonly deliveries: DeliveryRepository,
    private readonly providerEvents: ProviderEventRepository
  ) {}

  async execute(input: DispatchWorkflowTouchpointInput) {
    const tenantSlug = ensureTouchpointText(input.tenantSlug, "tenant_slug_required");
    const campaignPublicId = ensurePublicId(input.campaignPublicId, "campaign_public_id_invalid");
    const leadPublicId = ensurePublicId(input.leadPublicId, "lead_public_id_invalid");
    const workflowRunPublicId = ensurePublicId(input.workflowRunPublicId, "touchpoint_workflow_run_public_id_invalid");
    const contactValue = ensureTouchpointText(input.contactValue, "touchpoint_contact_value_required");
    const createdBy = ensureTouchpointText(input.createdBy, "touchpoint_created_by_required");
    ensureProvider(input.provider);
    const provider = ensureDeliveryProvider(input.provider);
    const campaign = await this.campaigns.getByPublicId(campaignPublicId);

    if (campaign === null || campaign.tenantSlug !== tenantSlug) {
      throw new Error("campaign_not_found");
    }

    let templatePublicId: string | null = null;
    if (input.templatePublicId) {
      const normalizedTemplatePublicId = ensurePublicId(input.templatePublicId, "template_public_id_invalid");
      const template = await this.templates.getByPublicId(normalizedTemplatePublicId);
      if (template === null || template.tenantSlug !== tenantSlug) {
        throw new Error("template_not_found");
      }

      templatePublicId = normalizedTemplatePublicId;
    }

    const touchpoint = await this.touchpoints.create({
      tenantSlug,
      campaignPublicId,
      campaignKey: campaign.key,
      channel: campaign.channel,
      workflowDefinitionKey: campaign.workflowDefinitionKey,
      leadPublicId,
      businessEntityType: "crm.lead",
      businessEntityPublicId: leadPublicId,
      contactValue,
      source: "workflow",
      createdBy,
      notes: (input.notes ?? "").trim() || `Workflow dispatch created from ${workflowRunPublicId}.`
    });

    const delivery = await this.deliveries.create(
      ensureDeliveryInput({
        tenantSlug,
        touchpointPublicId: touchpoint.publicId,
        templatePublicId,
        provider,
        providerMessageId: input.providerMessageId ?? null,
        sentBy: createdBy,
        notes: (input.notes ?? "").trim() || "Workflow dispatch created delivery."
      })
    );

    const updatedTouchpoint = await this.touchpoints.updateStatus(touchpoint.publicId, "sent", ensureOptionalPublicId(workflowRunPublicId));
    const providerEvent = await this.providerEvents.create({
      tenantSlug,
      provider,
      eventType: "workflow.dispatched",
      direction: "outbound",
      leadPublicId,
      businessEntityType: "crm.lead",
      businessEntityPublicId: leadPublicId,
      touchpointPublicId: touchpoint.publicId,
      deliveryPublicId: delivery.publicId,
      workflowRunPublicId,
      status: "processed",
      payloadSummary: `Workflow dispatch created touchpoint ${touchpoint.publicId}.`,
      responseSummary: `Delivery ${delivery.publicId} sent by workflow dispatch.`,
      processedAt: new Date().toISOString()
    });

    return {
      touchpoint: updatedTouchpoint ?? touchpoint,
      delivery,
      providerEvent
    };
  }
}
