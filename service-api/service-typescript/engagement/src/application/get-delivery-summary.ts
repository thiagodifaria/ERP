import { DeliveryRepository } from "../domain/delivery-repository.js";
import { DeliveryFilters, ensureDeliveryProvider, ensureDeliveryStatus } from "../domain/delivery.js";
import { ensureTemplateChannel } from "../domain/template.js";
import { ensureTouchpointText } from "../domain/touchpoint.js";

export class GetDeliverySummary {
  constructor(private readonly repository: DeliveryRepository) {}

  async execute(filters: DeliveryFilters = {}) {
    return this.repository.getSummary({
      tenantSlug: filters.tenantSlug ? ensureTouchpointText(filters.tenantSlug, "tenant_slug_required") : undefined,
      channel: filters.channel ? ensureTemplateChannel(filters.channel) : undefined,
      provider: filters.provider ? ensureDeliveryProvider(filters.provider) : undefined,
      status: filters.status ? ensureDeliveryStatus(filters.status) : undefined
    });
  }
}
