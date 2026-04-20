import { CampaignStatus } from "../../domain/campaign.js";

export type UpdateCampaignStatusRequest = {
  status: CampaignStatus;
};
