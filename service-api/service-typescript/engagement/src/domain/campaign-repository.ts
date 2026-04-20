import { Campaign, CampaignFilters, CampaignStatus, CreateCampaignInput } from "./campaign.js";

export interface CampaignRepository {
  list(filters?: CampaignFilters): Promise<Campaign[]>;
  getByPublicId(publicId: string): Promise<Campaign | null>;
  create(input: CreateCampaignInput): Promise<Campaign>;
  updateStatus(publicId: string, status: CampaignStatus): Promise<Campaign | null>;
}
