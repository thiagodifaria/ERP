import { DeliveryRepository } from "../domain/delivery-repository.js";
import { ProviderEventRepository } from "../domain/provider-event-repository.js";
import { ensureProvider } from "../domain/provider-event.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

export type RegisterProviderEventInput = {
  tenantSlug: string;
  provider: string;
  eventType: string;
  externalEventId?: string;
  touchpointPublicId?: string | null;
  deliveryPublicId?: string | null;
  workflowRunPublicId?: string | null;
  leadPublicId?: string | null;
  providerMessageId?: string | null;
  payloadSummary?: string;
  responseSummary?: string;
};

export class RegisterProviderEvent {
  constructor(
    private readonly touchpoints: TouchpointRepository,
    private readonly deliveries: DeliveryRepository,
    private readonly providerEvents: ProviderEventRepository
  ) {}

  async execute(input: RegisterProviderEventInput) {
    const tenantSlug = input.tenantSlug.trim();
    if (tenantSlug.length === 0) {
      throw new Error("tenant_slug_required");
    }

    const provider = ensureProvider(input.provider);
    const eventType = input.eventType.trim().toLowerCase();
    if (eventType.length < 3) {
      throw new Error("provider_event_type_invalid");
    }

    const externalEventId = (input.externalEventId ?? "").trim();
    if (externalEventId.length > 0) {
      const existing = await this.providerEvents.findByProviderAndExternalEventId(tenantSlug, provider, externalEventId);
      if (existing !== null) {
        throw new Error("provider_event_conflict");
      }
    }

    let touchpointPublicId = (input.touchpointPublicId ?? "").trim() || null;
    let deliveryPublicId = (input.deliveryPublicId ?? "").trim() || null;
    let status: "processed" | "failed" = "processed";
    let responseSummary = (input.responseSummary ?? "").trim();

    if (deliveryPublicId) {
      const delivery = await this.deliveries.getByPublicId(deliveryPublicId);
      if (delivery === null || delivery.tenantSlug !== tenantSlug) {
        throw new Error("delivery_not_found");
      }

      deliveryPublicId = delivery.publicId;
      touchpointPublicId = touchpointPublicId ?? delivery.touchpointPublicId;

      if (eventType === "delivery.delivered") {
        await this.deliveries.updateStatus(deliveryPublicId, {
          status: "delivered",
          providerMessageId: input.providerMessageId ?? null,
          notes: responseSummary || "Provider confirmed delivery."
        });
      } else if (eventType === "delivery.failed") {
        await this.deliveries.updateStatus(deliveryPublicId, {
          status: "failed",
          providerMessageId: input.providerMessageId ?? null,
          errorCode: "provider_delivery_failed",
          notes: responseSummary || "Provider returned delivery failure."
        });
        status = "failed";
      }
    }

    if (touchpointPublicId) {
      const touchpoint = await this.touchpoints.getByPublicId(touchpointPublicId);
      if (touchpoint === null || touchpoint.tenantSlug !== tenantSlug) {
        throw new Error("touchpoint_not_found");
      }

      if (eventType === "delivery.delivered" && touchpoint.status === "queued") {
        await this.touchpoints.updateStatus(touchpointPublicId, "delivered", input.workflowRunPublicId ?? null);
      } else if (eventType === "delivery.responded") {
        await this.touchpoints.updateStatus(touchpointPublicId, "responded", input.workflowRunPublicId ?? null);
      } else if (eventType === "delivery.converted") {
        await this.touchpoints.updateStatus(touchpointPublicId, "converted", input.workflowRunPublicId ?? null);
      } else if (eventType === "delivery.failed") {
        await this.touchpoints.updateStatus(touchpointPublicId, "failed", input.workflowRunPublicId ?? null);
      }
    }

    if (responseSummary.length === 0) {
      responseSummary = eventType === "delivery.responded"
        ? "Provider response linked to touchpoint."
        : eventType === "delivery.converted"
          ? "Provider conversion linked to touchpoint."
          : status === "failed"
            ? "Provider returned a failed delivery event."
            : "Provider event processed successfully.";
    }

    return this.providerEvents.create({
      tenantSlug,
      provider,
      eventType,
      direction: "inbound",
      externalEventId,
      leadPublicId: input.leadPublicId ?? null,
      touchpointPublicId,
      deliveryPublicId,
      workflowRunPublicId: input.workflowRunPublicId ?? null,
      status,
      payloadSummary: (input.payloadSummary ?? "").trim() || `${provider} provider event ${eventType} received.`,
      responseSummary,
      processedAt: new Date().toISOString()
    });
  }
}
