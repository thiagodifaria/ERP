import { DeliveryRepository } from "../domain/delivery-repository.js";
import { UpdateDeliveryStatusInput, ensureDeliveryStatusInput } from "../domain/delivery.js";
import { ensurePublicId } from "../domain/touchpoint.js";

export class UpdateTouchpointDeliveryStatus {
  constructor(private readonly repository: DeliveryRepository) {}

  async execute(publicId: string, input: UpdateDeliveryStatusInput) {
    return this.repository.updateStatus(
      ensurePublicId(publicId, "delivery_public_id_invalid"),
      ensureDeliveryStatusInput(input)
    );
  }
}
