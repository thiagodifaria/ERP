import { CampaignRepository } from "../domain/campaign-repository.js";
import { ensureCampaignStatus } from "../domain/campaign.js";
import { ensurePublicId } from "../domain/touchpoint.js";

export class UpdateCampaignStatus {
  constructor(private readonly repository: CampaignRepository) {}

  async execute(publicId: string, status: string) {
    return this.repository.updateStatus(ensurePublicId(publicId, "campaign_public_id_invalid"), ensureCampaignStatus(status));
  }
}
