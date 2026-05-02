export const engagementProviders = ["resend", "whatsapp_cloud", "telegram_bot", "meta_ads", "manual"] as const;
export const providerEventDirections = ["inbound", "outbound"] as const;
export const providerEventStatuses = ["received", "processed", "failed"] as const;

export type EngagementProvider = (typeof engagementProviders)[number];
export type ProviderEventDirection = (typeof providerEventDirections)[number];
export type ProviderEventStatus = (typeof providerEventStatuses)[number];

export type ProviderEvent = {
  id: number;
  publicId: string;
  tenantSlug: string;
  provider: EngagementProvider;
  eventType: string;
  direction: ProviderEventDirection;
  externalEventId: string | null;
  leadPublicId: string | null;
  businessEntityType: string | null;
  businessEntityPublicId: string | null;
  touchpointPublicId: string | null;
  deliveryPublicId: string | null;
  workflowRunPublicId: string | null;
  status: ProviderEventStatus;
  payloadSummary: string;
  responseSummary: string;
  createdAt: string;
  processedAt: string | null;
};

export type ProviderEventFilters = {
  tenantSlug?: string;
  provider?: EngagementProvider;
  eventType?: string;
  direction?: ProviderEventDirection;
  status?: ProviderEventStatus;
  businessEntityType?: string;
  businessEntityPublicId?: string;
};

export type CreateProviderEventInput = {
  tenantSlug: string;
  provider: EngagementProvider;
  eventType: string;
  direction: ProviderEventDirection;
  externalEventId?: string | null;
  leadPublicId?: string | null;
  businessEntityType?: string | null;
  businessEntityPublicId?: string | null;
  touchpointPublicId?: string | null;
  deliveryPublicId?: string | null;
  workflowRunPublicId?: string | null;
  status: ProviderEventStatus;
  payloadSummary?: string;
  responseSummary?: string;
  processedAt?: string | null;
};

export type ProviderEventSummary = {
  tenantSlug: string;
  generatedAt: string;
  totals: {
    total: number;
    inbound: number;
    outbound: number;
    processed: number;
    failed: number;
  };
  byProvider: Record<EngagementProvider, number>;
  byStatus: Record<ProviderEventStatus, number>;
  byDirection: Record<ProviderEventDirection, number>;
};

export type ProviderCapability = {
  provider: EngagementProvider;
  scope: "messaging" | "ads" | "email" | "manual";
  configured: boolean;
  mode: "configured" | "fallback" | "manual";
  supportsInbound: boolean;
  supportsOutbound: boolean;
  supportsTracking: boolean;
};

function ensureIncluded<T extends string>(value: string, items: readonly T[], errorCode: string): T {
  const normalized = value.trim().toLowerCase();
  if (!items.includes(normalized as T)) {
    throw new Error(errorCode);
  }

  return normalized as T;
}

export function ensureProvider(value: string): EngagementProvider {
  return ensureIncluded(value, engagementProviders, "provider_invalid");
}

export function ensureProviderDirection(value: string): ProviderEventDirection {
  return ensureIncluded(value, providerEventDirections, "provider_event_direction_invalid");
}

export function ensureProviderEventStatus(value: string): ProviderEventStatus {
  return ensureIncluded(value, providerEventStatuses, "provider_event_status_invalid");
}

export function ensureProviderEventType(value: string): string {
  const normalized = value.trim().toLowerCase();
  if (normalized.length < 3) {
    throw new Error("provider_event_type_invalid");
  }

  return normalized;
}

function ensureText(value: string | null | undefined): string {
  return (value ?? "").trim();
}

function ensureOptionalPublicId(value: string | null | undefined, errorCode: string): string | null {
  const normalized = ensureText(value);
  if (normalized.length === 0) {
    return null;
  }

  if (!/^[0-9a-fA-F-]{36}$/.test(normalized)) {
    throw new Error(errorCode);
  }

  return normalized.toLowerCase();
}

function ensureOptionalBusinessEntityType(value: string | null | undefined): string | null {
  const normalized = ensureText(value).toLowerCase();
  if (normalized.length === 0) {
    return null;
  }

  if (!/^[a-z0-9_.-]{3,80}$/.test(normalized)) {
    throw new Error("business_entity_type_invalid");
  }

  return normalized;
}

export function ensureProviderEventInput(input: CreateProviderEventInput): CreateProviderEventInput {
  const tenantSlug = ensureText(input.tenantSlug);
  if (tenantSlug.length === 0) {
    throw new Error("tenant_slug_required");
  }

  return {
    tenantSlug,
    provider: ensureProvider(input.provider),
    eventType: ensureProviderEventType(input.eventType),
    direction: ensureProviderDirection(input.direction),
    externalEventId: ensureText(input.externalEventId),
    leadPublicId: ensureOptionalPublicId(input.leadPublicId, "lead_public_id_invalid"),
    businessEntityType: ensureOptionalBusinessEntityType(input.businessEntityType),
    businessEntityPublicId: ensureOptionalPublicId(input.businessEntityPublicId, "business_entity_public_id_invalid"),
    touchpointPublicId: ensureOptionalPublicId(input.touchpointPublicId, "touchpoint_public_id_invalid"),
    deliveryPublicId: ensureOptionalPublicId(input.deliveryPublicId, "delivery_public_id_invalid"),
    workflowRunPublicId: ensureOptionalPublicId(input.workflowRunPublicId, "touchpoint_workflow_run_public_id_invalid"),
    status: ensureProviderEventStatus(input.status),
    payloadSummary: ensureText(input.payloadSummary),
    responseSummary: ensureText(input.responseSummary),
    processedAt: ensureText(input.processedAt) || null
  };
}

export function buildProviderEventSummary(events: ProviderEvent[], tenantSlug?: string): ProviderEventSummary {
  const byProvider = Object.fromEntries(engagementProviders.map((provider) => [provider, 0])) as Record<EngagementProvider, number>;
  const byStatus = Object.fromEntries(providerEventStatuses.map((status) => [status, 0])) as Record<ProviderEventStatus, number>;
  const byDirection = Object.fromEntries(providerEventDirections.map((direction) => [direction, 0])) as Record<ProviderEventDirection, number>;

  for (const event of events) {
    byProvider[event.provider] += 1;
    byStatus[event.status] += 1;
    byDirection[event.direction] += 1;
  }

  return {
    tenantSlug: tenantSlug ?? "global",
    generatedAt: new Date().toISOString(),
    totals: {
      total: events.length,
      inbound: byDirection.inbound,
      outbound: byDirection.outbound,
      processed: byStatus.processed,
      failed: byStatus.failed
    },
    byProvider,
    byStatus,
    byDirection
  };
}
