import { DeliveryRepository } from "../domain/delivery-repository.js";
import { ensurePublicId, ensureTouchpointText } from "../domain/touchpoint.js";

export class ListTouchpointDeliveries {
  constructor(private readonly repository: DeliveryRepository) {}

  async execute(touchpointPublicId: string, tenantSlug?: string) {
    return this.repository.listByTouchpointPublicId(
      ensurePublicId(touchpointPublicId, "touchpoint_public_id_invalid"),
      tenantSlug ? ensureTouchpointText(tenantSlug, "tenant_slug_required") : undefined
    );
  }
}
