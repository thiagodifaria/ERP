import { DeliveryRepository } from "../domain/delivery-repository.js";
import { CreateDeliveryInput, ensureDeliveryInput } from "../domain/delivery.js";

export class CreateTouchpointDelivery {
  constructor(private readonly repository: DeliveryRepository) {}

  async execute(touchpointPublicId: string, input: Omit<CreateDeliveryInput, "touchpointPublicId">) {
    return this.repository.create(
      ensureDeliveryInput({
        ...input,
        touchpointPublicId
      })
    );
  }
}
