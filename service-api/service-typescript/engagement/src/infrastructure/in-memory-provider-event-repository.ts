import { randomUUID } from "node:crypto";
import { ProviderEventRepository } from "../domain/provider-event-repository.js";
import {
  buildProviderEventSummary,
  CreateProviderEventInput,
  ensureProviderEventInput,
  ProviderEvent,
  ProviderEventFilters,
  ProviderEventSummary
} from "../domain/provider-event.js";

export class InMemoryProviderEventRepository implements ProviderEventRepository {
  private nextId = 1;
  private readonly events: ProviderEvent[] = [];

  constructor(private readonly bootstrapTenantSlug: string) {}

  async list(filters: ProviderEventFilters = {}): Promise<ProviderEvent[]> {
    return this.events.filter((event) => {
      if (filters.tenantSlug && event.tenantSlug !== filters.tenantSlug) {
        return false;
      }
      if (filters.provider && event.provider !== filters.provider) {
        return false;
      }
      if (filters.eventType && event.eventType !== filters.eventType) {
        return false;
      }
      if (filters.direction && event.direction !== filters.direction) {
        return false;
      }
      if (filters.status && event.status !== filters.status) {
        return false;
      }
      if (filters.businessEntityType && event.businessEntityType !== filters.businessEntityType) {
        return false;
      }
      if (filters.businessEntityPublicId && event.businessEntityPublicId !== filters.businessEntityPublicId) {
        return false;
      }

      return true;
    });
  }

  async getByPublicId(publicId: string): Promise<ProviderEvent | null> {
    return this.events.find((event) => event.publicId === publicId) ?? null;
  }

  async findByProviderAndExternalEventId(tenantSlug: string, provider: string, externalEventId: string): Promise<ProviderEvent | null> {
    const normalizedExternalEventId = externalEventId.trim();
    if (normalizedExternalEventId.length === 0) {
      return null;
    }

    return this.events.find((event) => event.tenantSlug === tenantSlug && event.provider === provider && event.externalEventId === normalizedExternalEventId) ?? null;
  }

  async create(input: CreateProviderEventInput): Promise<ProviderEvent> {
    const normalized = ensureProviderEventInput(input);

    if (normalized.externalEventId && (await this.findByProviderAndExternalEventId(normalized.tenantSlug, normalized.provider, normalized.externalEventId)) !== null) {
      throw new Error("provider_event_conflict");
    }

    const createdAt = new Date().toISOString();
    const event: ProviderEvent = {
      id: this.nextId++,
      publicId: randomUUID(),
      tenantSlug: normalized.tenantSlug,
      provider: normalized.provider,
      eventType: normalized.eventType,
      direction: normalized.direction,
      externalEventId: normalized.externalEventId || null,
      leadPublicId: normalized.leadPublicId ?? null,
      businessEntityType: normalized.businessEntityType ?? null,
      businessEntityPublicId: normalized.businessEntityPublicId ?? null,
      touchpointPublicId: normalized.touchpointPublicId ?? null,
      deliveryPublicId: normalized.deliveryPublicId ?? null,
      workflowRunPublicId: normalized.workflowRunPublicId ?? null,
      status: normalized.status,
      payloadSummary: normalized.payloadSummary ?? "",
      responseSummary: normalized.responseSummary ?? "",
      createdAt,
      processedAt: normalized.processedAt ?? (normalized.status === "processed" ? createdAt : null)
    };

    this.events.push(event);
    return event;
  }

  async getSummary(filters: ProviderEventFilters = {}): Promise<ProviderEventSummary> {
    return buildProviderEventSummary(await this.list(filters), filters.tenantSlug ?? this.bootstrapTenantSlug);
  }
}
