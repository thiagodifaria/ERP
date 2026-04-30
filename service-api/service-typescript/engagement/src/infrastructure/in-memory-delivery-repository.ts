import { randomUUID } from "node:crypto";
import { DeliveryRepository } from "../domain/delivery-repository.js";
import { buildDeliverySummary, CreateDeliveryInput, DeliveryFilters, DeliverySummary, TouchpointDelivery, UpdateDeliveryStatusInput } from "../domain/delivery.js";
import { EngagementTemplate } from "../domain/template.js";
import { Touchpoint } from "../domain/touchpoint.js";

function bootstrapDeliveries(tenantSlug: string): TouchpointDelivery[] {
  return [
    {
      id: 1,
      publicId: "00000000-0000-0000-0000-00000000d101",
      tenantSlug,
      touchpointPublicId: "00000000-0000-0000-0000-00000000e101",
      templatePublicId: "00000000-0000-0000-0000-00000000f101",
      templateKey: "lead-follow-up-whatsapp",
      channel: "whatsapp",
      provider: "manual",
      providerMessageId: "bootstrap-manual-001",
      status: "delivered",
      sentBy: "bootstrap-owner",
      errorCode: null,
      notes: "Entrega bootstrap confirmada.",
      attemptedAt: "2026-04-20T00:10:30.000Z",
      createdAt: "2026-04-20T00:10:30.000Z",
      updatedAt: "2026-04-20T00:10:30.000Z"
    }
  ];
}

export class InMemoryDeliveryRepository implements DeliveryRepository {
  private readonly deliveries: TouchpointDelivery[];

  constructor(
    tenantSlug: string,
    private readonly templates: () => Promise<EngagementTemplate[]>,
    private readonly touchpoints: () => Promise<Touchpoint[]>
  ) {
    this.deliveries = bootstrapDeliveries(tenantSlug);
  }

  async list(filters: DeliveryFilters = {}): Promise<TouchpointDelivery[]> {
    return this.deliveries.filter((delivery) => {
      if (filters.tenantSlug && delivery.tenantSlug !== filters.tenantSlug) {
        return false;
      }

      if (filters.touchpointPublicId && delivery.touchpointPublicId !== filters.touchpointPublicId) {
        return false;
      }

      if (filters.channel && delivery.channel !== filters.channel) {
        return false;
      }

      if (filters.provider && delivery.provider !== filters.provider) {
        return false;
      }

      if (filters.status && delivery.status !== filters.status) {
        return false;
      }

      return true;
    });
  }

  async getByPublicId(publicId: string): Promise<TouchpointDelivery | null> {
    return this.deliveries.find((delivery) => delivery.publicId === publicId) ?? null;
  }

  async create(input: CreateDeliveryInput): Promise<TouchpointDelivery> {
    const touchpoint = (await this.touchpoints()).find(
      (candidate) => candidate.publicId === input.touchpointPublicId && candidate.tenantSlug === input.tenantSlug
    );

    if (!touchpoint) {
      throw new Error("touchpoint_not_found");
    }

    let templateKey: string | null = null;
    if (input.templatePublicId) {
      const template = (await this.templates()).find(
        (candidate) => candidate.publicId === input.templatePublicId && candidate.tenantSlug === input.tenantSlug
      );

      if (!template) {
        throw new Error("template_not_found");
      }

      if (template.channel !== touchpoint.channel) {
        throw new Error("delivery_channel_mismatch");
      }

      templateKey = template.key;
    }

    const now = new Date().toISOString();
    const delivery: TouchpointDelivery = {
      id: this.deliveries.length + 1,
      publicId: randomUUID(),
      tenantSlug: input.tenantSlug,
      touchpointPublicId: input.touchpointPublicId,
      templatePublicId: input.templatePublicId ?? null,
      templateKey,
      channel: touchpoint.channel,
      provider: input.provider,
      providerMessageId: input.providerMessageId ?? null,
      status: "sent",
      sentBy: input.sentBy,
      errorCode: null,
      notes: input.notes ?? "",
      attemptedAt: now,
      createdAt: now,
      updatedAt: now
    };

    this.deliveries.push(delivery);
    return delivery;
  }

  async updateStatus(publicId: string, input: UpdateDeliveryStatusInput): Promise<TouchpointDelivery | null> {
    const delivery = this.deliveries.find((item) => item.publicId === publicId);

    if (!delivery) {
      return null;
    }

    delivery.status = input.status;
    delivery.providerMessageId = input.providerMessageId ?? delivery.providerMessageId;
    delivery.errorCode = input.errorCode ?? null;
    if (input.notes && input.notes.length > 0) {
      delivery.notes = input.notes;
    }
    delivery.updatedAt = new Date().toISOString();
    return delivery;
  }

  async getSummary(filters: DeliveryFilters = {}): Promise<DeliverySummary> {
    const deliveries = await this.list(filters);
    const templates = await this.templates();
    const touchpoints = await this.touchpoints();
    const tenantSlug = filters.tenantSlug ?? deliveries[0]?.tenantSlug ?? touchpoints[0]?.tenantSlug ?? "global";

    return buildDeliverySummary(
      tenantSlug,
      templates.filter((template) => !filters.tenantSlug || template.tenantSlug === filters.tenantSlug),
      touchpoints.filter((touchpoint) => !filters.tenantSlug || touchpoint.tenantSlug === filters.tenantSlug),
      deliveries
    );
  }

  async listByTouchpointPublicId(touchpointPublicId: string, tenantSlug?: string): Promise<TouchpointDelivery[]> {
    return this.list({ tenantSlug, touchpointPublicId });
  }
}
