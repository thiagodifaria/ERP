import { ProviderEventRepository } from "../domain/provider-event-repository.js";
import { ensurePublicId } from "../domain/touchpoint.js";

export class GetProviderEventByPublicId {
  constructor(private readonly repository: ProviderEventRepository) {}

  async execute(publicId: string) {
    return this.repository.getByPublicId(ensurePublicId(publicId, "provider_event_public_id_invalid"));
  }
}
