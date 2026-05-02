import { CreateProviderEventInput, ProviderEvent, ProviderEventFilters, ProviderEventSummary } from "./provider-event.js";

export interface ProviderEventRepository {
  list(filters?: ProviderEventFilters): Promise<ProviderEvent[]>;
  getByPublicId(publicId: string): Promise<ProviderEvent | null>;
  findByProviderAndExternalEventId(tenantSlug: string, provider: string, externalEventId: string): Promise<ProviderEvent | null>;
  create(input: CreateProviderEventInput): Promise<ProviderEvent>;
  getSummary(filters?: ProviderEventFilters): Promise<ProviderEventSummary>;
}
