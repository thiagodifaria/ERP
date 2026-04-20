import { CampaignChannel, CampaignStatus } from "../../domain/campaign.js";

export type CreateCampaignRequest = {
  tenantSlug: string;
  key: string;
  name: string;
  description: string;
  channel: CampaignChannel;
  status?: CampaignStatus;
  touchpointGoal: string;
  workflowDefinitionKey?: string | null;
  budgetCents: number;
};
