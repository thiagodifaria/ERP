import { CreateDeliveryInput, DeliveryFilters, DeliveryStatus, DeliverySummary, TouchpointDelivery, UpdateDeliveryStatusInput } from "./delivery.js";

export interface DeliveryRepository {
  list(filters?: DeliveryFilters): Promise<TouchpointDelivery[]>;
  getByPublicId(publicId: string): Promise<TouchpointDelivery | null>;
  create(input: CreateDeliveryInput): Promise<TouchpointDelivery>;
  updateStatus(publicId: string, input: UpdateDeliveryStatusInput): Promise<TouchpointDelivery | null>;
  getSummary(filters?: DeliveryFilters): Promise<DeliverySummary>;
  listByTouchpointPublicId(touchpointPublicId: string, tenantSlug?: string): Promise<TouchpointDelivery[]>;
}
