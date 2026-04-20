import { CampaignRepository } from "../domain/campaign-repository.js";

export class GetCampaignByPublicId {
  constructor(private readonly repository: CampaignRepository) {}

  async execute(publicId: string) {
    return this.repository.getByPublicId(publicId);
  }
}
