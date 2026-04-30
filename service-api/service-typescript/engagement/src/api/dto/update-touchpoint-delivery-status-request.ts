import { DeliveryStatus } from "../../domain/delivery.js";

export type UpdateTouchpointDeliveryStatusRequest = {
  status: DeliveryStatus;
  providerMessageId?: string | null;
  errorCode?: string | null;
  notes?: string;
};
