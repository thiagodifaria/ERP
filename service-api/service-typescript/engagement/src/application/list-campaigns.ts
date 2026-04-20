import { CampaignFilters } from "../domain/campaign.js";
import { CampaignRepository } from "../domain/campaign-repository.js";

export class ListCampaigns {
  constructor(private readonly repository: CampaignRepository) {}

  async execute(filters?: CampaignFilters) {
    return this.repository.list(filters);
  }
}
