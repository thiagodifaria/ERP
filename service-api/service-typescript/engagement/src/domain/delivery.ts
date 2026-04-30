import { CampaignChannel, campaignChannels } from "./campaign.js";
import { ensurePublicId, ensureTouchpointText } from "./touchpoint.js";
import { EngagementTemplate, TemplateProvider, templateProviders } from "./template.js";

export const deliveryStatuses = ["queued", "sent", "delivered", "failed"] as const;
export type DeliveryStatus = (typeof deliveryStatuses)[number];

export type TouchpointDelivery = {
  id: number;
  publicId: string;
  tenantSlug: string;
  touchpointPublicId: string;
  templatePublicId: string | null;
  templateKey: string | null;
  channel: CampaignChannel;
  provider: TemplateProvider;
  providerMessageId: string | null;
  status: DeliveryStatus;
  sentBy: string;
  errorCode: string | null;
  notes: string;
  attemptedAt: string;
  createdAt: string;
  updatedAt: string;
};

export type DeliveryFilters = {
  tenantSlug?: string;
  touchpointPublicId?: string;
  channel?: CampaignChannel;
  provider?: TemplateProvider;
  status?: DeliveryStatus;
};

export type CreateDeliveryInput = {
  tenantSlug: string;
  touchpointPublicId: string;
  templatePublicId?: string | null;
  provider: TemplateProvider;
  providerMessageId?: string | null;
  sentBy: string;
  notes?: string;
};

export type UpdateDeliveryStatusInput = {
  status: DeliveryStatus;
  providerMessageId?: string | null;
  errorCode?: string | null;
  notes?: string;
};

export type DeliverySummary = {
  tenantSlug: string;
  generatedAt: string;
  totals: {
    templates: number;
    activeTemplates: number;
    touchpoints: number;
    deliveries: number;
    convertedTouchpoints: number;
    workflowDispatched: number;
  };
  byProvider: Record<TemplateProvider, number>;
  byStatus: Record<DeliveryStatus, number>;
  outcomes: {
    delivered: number;
    failed: number;
    deliveryRate: number;
    failureRate: number;
    templateLinked: number;
  };
};

function ensureIncluded<T extends string>(value: string, items: readonly T[], errorCode: string): T {
  const normalizedValue = value.trim().toLowerCase();

  if (!items.includes(normalizedValue as T)) {
    throw new Error(errorCode);
  }

  return normalizedValue as T;
}

export function ensureDeliveryStatus(value: string): DeliveryStatus {
  return ensureIncluded(value, deliveryStatuses, "delivery_status_invalid");
}

export function ensureDeliveryProvider(value: string): TemplateProvider {
  return ensureIncluded(value, templateProviders, "delivery_provider_invalid");
}

export function ensureDeliveryChannel(value: string): CampaignChannel {
  return ensureIncluded(value, campaignChannels, "delivery_channel_invalid");
}

export function ensureOptionalText(value: string | null | undefined): string | null {
  if (value == null) {
    return null;
  }

  const normalizedValue = value.trim();
  return normalizedValue.length > 0 ? normalizedValue : null;
}

export function ensureDeliveryInput(input: CreateDeliveryInput): CreateDeliveryInput {
  return {
    tenantSlug: ensureTouchpointText(input.tenantSlug, "tenant_slug_required"),
    touchpointPublicId: ensurePublicId(input.touchpointPublicId, "touchpoint_public_id_invalid"),
    templatePublicId: input.templatePublicId ? ensurePublicId(input.templatePublicId, "template_public_id_invalid") : null,
    provider: ensureDeliveryProvider(input.provider),
    providerMessageId: ensureOptionalText(input.providerMessageId),
    sentBy: ensureTouchpointText(input.sentBy, "delivery_sent_by_required"),
    notes: ensureOptionalText(input.notes) ?? ""
  };
}

export function ensureDeliveryStatusInput(input: UpdateDeliveryStatusInput): UpdateDeliveryStatusInput {
  return {
    status: ensureDeliveryStatus(input.status),
    providerMessageId: ensureOptionalText(input.providerMessageId),
    errorCode: ensureOptionalText(input.errorCode),
    notes: ensureOptionalText(input.notes) ?? ""
  };
}

export function buildDeliverySummary(
  tenantSlug: string,
  templates: EngagementTemplate[],
  touchpoints: Array<{ status: string; lastWorkflowRunPublicId: string | null }>,
  deliveries: TouchpointDelivery[]
): DeliverySummary {
  const byProvider = Object.fromEntries(templateProviders.map((provider) => [provider, 0])) as Record<TemplateProvider, number>;
  const byStatus = Object.fromEntries(deliveryStatuses.map((status) => [status, 0])) as Record<DeliveryStatus, number>;

  for (const delivery of deliveries) {
    byProvider[delivery.provider] += 1;
    byStatus[delivery.status] += 1;
  }

  const delivered = byStatus.delivered;
  const failed = byStatus.failed;
  const total = deliveries.length;
  const templateLinked = deliveries.filter((delivery) => delivery.templatePublicId !== null).length;

  return {
    tenantSlug,
    generatedAt: new Date().toISOString(),
    totals: {
      templates: templates.length,
      activeTemplates: templates.filter((template) => template.status === "active").length,
      touchpoints: touchpoints.length,
      deliveries: total,
      convertedTouchpoints: touchpoints.filter((touchpoint) => touchpoint.status === "converted").length,
      workflowDispatched: touchpoints.filter((touchpoint) => touchpoint.lastWorkflowRunPublicId !== null).length
    },
    byProvider,
    byStatus,
    outcomes: {
      delivered,
      failed,
      deliveryRate: total > 0 ? Number((delivered / total).toFixed(4)) : 0,
      failureRate: total > 0 ? Number((failed / total).toFixed(4)) : 0,
      templateLinked
    }
  };
}
