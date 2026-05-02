import { ProviderEventFilters } from "../domain/provider-event.js";
import { ProviderEventRepository } from "../domain/provider-event-repository.js";

export class GetProviderEventSummary {
  constructor(private readonly repository: ProviderEventRepository) {}

  async execute(filters: ProviderEventFilters = {}) {
    return this.repository.getSummary(filters);
  }
}
