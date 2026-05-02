import { CampaignRepository } from "../domain/campaign-repository.js";
import { ensurePublicId, ensureTouchpointText } from "../domain/touchpoint.js";
import { ProviderEventRepository } from "../domain/provider-event-repository.js";
import { ensureProvider } from "../domain/provider-event.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";
import { CrmGateway } from "../infrastructure/crm-gateway.js";

export type IngestProviderLeadInput = {
  tenantSlug: string;
  provider: string;
  campaignPublicId?: string;
  externalEventId?: string;
  name: string;
  email: string;
  contactValue?: string;
  notes?: string;
};

export class IngestProviderLead {
  constructor(
    private readonly campaigns: CampaignRepository,
    private readonly touchpoints: TouchpointRepository,
    private readonly providerEvents: ProviderEventRepository,
    private readonly crmGateway: CrmGateway
  ) {}

  async execute(input: IngestProviderLeadInput) {
    const tenantSlug = ensureTouchpointText(input.tenantSlug, "tenant_slug_required");
    const provider = ensureProvider(input.provider);
    const name = ensureTouchpointText(input.name, "provider_lead_name_required");
    const email = ensureTouchpointText(input.email, "provider_lead_email_required");
    const contactValue = ensureTouchpointText(input.contactValue ?? input.email, "touchpoint_contact_value_required");
    const externalEventId = (input.externalEventId ?? "").trim();
    const notes = (input.notes ?? "").trim();

    if (externalEventId.length > 0) {
      const existing = await this.providerEvents.findByProviderAndExternalEventId(tenantSlug, provider, externalEventId);
      if (existing !== null) {
        throw new Error("provider_event_conflict");
      }
    }

    const lead = await this.crmGateway.createLead({
      tenantSlug,
      name,
      email,
      source: provider
    });

    let touchpoint = null;
    if (input.campaignPublicId) {
      const campaignPublicId = ensurePublicId(input.campaignPublicId, "campaign_public_id_invalid");
      const campaign = await this.campaigns.getByPublicId(campaignPublicId);
      if (campaign === null || campaign.tenantSlug !== tenantSlug) {
        throw new Error("campaign_not_found");
      }

      touchpoint = await this.touchpoints.create({
        tenantSlug,
        campaignPublicId,
        campaignKey: campaign.key,
        channel: campaign.channel,
        workflowDefinitionKey: campaign.workflowDefinitionKey,
        leadPublicId: lead.publicId,
        businessEntityType: "crm.lead",
        businessEntityPublicId: lead.publicId,
        contactValue,
        source: provider,
        createdBy: "engagement-provider",
        notes: notes.length > 0 ? notes : `Inbound lead captured from ${provider}.`
      });
    }

    const providerEvent = await this.providerEvents.create({
      tenantSlug,
      provider,
      eventType: "lead.ingested",
      direction: "inbound",
      externalEventId,
      leadPublicId: lead.publicId,
      businessEntityType: "crm.lead",
      businessEntityPublicId: lead.publicId,
      touchpointPublicId: touchpoint?.publicId ?? null,
      status: "processed",
      payloadSummary: `Inbound lead ${email} captured from ${provider}.`,
      responseSummary: touchpoint ? `Lead linked to touchpoint ${touchpoint.publicId}.` : "Lead captured without campaign linkage.",
      processedAt: new Date().toISOString()
    });

    return {
      lead,
      touchpoint,
      providerEvent
    };
  }
}
