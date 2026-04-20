import { CampaignRepository } from "../domain/campaign-repository.js";
import {
  CreateTouchpointInput,
  ensurePublicId,
  ensureTouchpointText
} from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export class CreateTouchpoint {
  constructor(
    private readonly campaigns: CampaignRepository,
    private readonly repository: TouchpointRepository
  ) {}

  async execute(input: CreateTouchpointInput) {
    const tenantSlug = ensureTouchpointText(input.tenantSlug, "tenant_slug_required");
    const campaignPublicId = ensurePublicId(input.campaignPublicId, "campaign_public_id_invalid");
    const campaign = await this.campaigns.getByPublicId(campaignPublicId);

    if (campaign === null || campaign.tenantSlug !== tenantSlug) {
      throw new Error("campaign_not_found");
    }

    return this.repository.create({
      tenantSlug,
      campaignPublicId,
      campaignKey: campaign.key,
      channel: campaign.channel,
      workflowDefinitionKey: campaign.workflowDefinitionKey,
      leadPublicId: ensurePublicId(input.leadPublicId, "lead_public_id_invalid"),
      contactValue: ensureTouchpointText(input.contactValue, "touchpoint_contact_value_required"),
      source: ensureTouchpointText(input.source, "touchpoint_source_required"),
      createdBy: ensureTouchpointText(input.createdBy, "touchpoint_created_by_required"),
      notes: input.notes.trim(),
    });
  }
}
