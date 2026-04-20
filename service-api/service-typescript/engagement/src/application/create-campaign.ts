import {
  CampaignRepository
} from "../domain/campaign-repository.js";
import {
  CreateCampaignInput,
  ensureBudgetCents,
  ensureCampaignChannel,
  ensureCampaignKey,
  ensureCampaignStatus,
  ensureText
} from "../domain/campaign.js";

export class CreateCampaign {
  constructor(private readonly repository: CampaignRepository) {}

  async execute(input: CreateCampaignInput) {
    return this.repository.create({
      tenantSlug: ensureText(input.tenantSlug, "tenant_slug_required"),
      key: ensureCampaignKey(input.key),
      name: ensureText(input.name, "campaign_name_required"),
      description: ensureText(input.description, "campaign_description_required"),
      channel: ensureCampaignChannel(input.channel),
      status: input.status ? ensureCampaignStatus(input.status) : "draft",
      touchpointGoal: ensureText(input.touchpointGoal, "campaign_touchpoint_goal_required"),
      workflowDefinitionKey:
        input.workflowDefinitionKey && input.workflowDefinitionKey.trim().length > 0
          ? input.workflowDefinitionKey.trim()
          : null,
      budgetCents: ensureBudgetCents(input.budgetCents)
    });
  }
}
