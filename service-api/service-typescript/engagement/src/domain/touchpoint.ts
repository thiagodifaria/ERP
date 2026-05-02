import { CampaignChannel, campaignChannels, ensureText } from "./campaign.js";

export const touchpointStatuses = ["queued", "sent", "delivered", "responded", "converted", "failed"] as const;
export type TouchpointStatus = (typeof touchpointStatuses)[number];

export type Touchpoint = {
  id: number;
  publicId: string;
  tenantSlug: string;
  campaignPublicId: string;
  campaignKey: string;
  leadPublicId: string;
  businessEntityType: string | null;
  businessEntityPublicId: string | null;
  channel: CampaignChannel;
  contactValue: string;
  source: string;
  status: TouchpointStatus;
  workflowDefinitionKey: string | null;
  lastWorkflowRunPublicId: string | null;
  createdBy: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
};

export type TouchpointFilters = {
  tenantSlug?: string;
  campaignPublicId?: string;
  status?: TouchpointStatus;
  channel?: CampaignChannel;
  leadPublicId?: string;
  businessEntityType?: string;
  businessEntityPublicId?: string;
};

export type CreateTouchpointInput = {
  tenantSlug: string;
  campaignPublicId: string;
  leadPublicId: string;
  businessEntityType?: string | null;
  businessEntityPublicId?: string | null;
  contactValue: string;
  source: string;
  createdBy: string;
  notes: string;
};

export type UpdateTouchpointStatusInput = {
  status: TouchpointStatus;
  lastWorkflowRunPublicId?: string | null;
};

export type TouchpointSummary = {
  tenantSlug: string;
  generatedAt: string;
  totals: {
    campaigns: number;
    touchpoints: number;
    workflowConfigured: number;
    workflowDispatched: number;
    businessLinked: number;
  };
  byStatus: Record<TouchpointStatus, number>;
  byChannel: Record<CampaignChannel, number>;
  outcomes: {
    responded: number;
    converted: number;
    failed: number;
    responseRate: number;
    conversionRate: number;
  };
};

export function buildTouchpointSummary(
  touchpoints: Touchpoint[],
  tenantSlug?: string
): TouchpointSummary {
  const campaignIds = new Set(touchpoints.map((touchpoint) => touchpoint.campaignPublicId));
  const byStatus = Object.fromEntries(touchpointStatuses.map((status) => [status, 0])) as Record<TouchpointStatus, number>;
  const byChannel = Object.fromEntries(campaignChannels.map((channel) => [channel, 0])) as Record<CampaignChannel, number>;

  for (const touchpoint of touchpoints) {
    byStatus[touchpoint.status] += 1;
    byChannel[touchpoint.channel] += 1;
  }

  const responded = byStatus.responded + byStatus.converted;
  const converted = byStatus.converted;
  const failed = byStatus.failed;
  const sentBase = byStatus.sent + byStatus.delivered + byStatus.responded + byStatus.converted;
  const responseBase = byStatus.responded + byStatus.converted + byStatus.failed;

  return {
    tenantSlug: tenantSlug ?? "global",
    generatedAt: new Date().toISOString(),
    totals: {
      campaigns: campaignIds.size,
      touchpoints: touchpoints.length,
      workflowConfigured: touchpoints.filter((touchpoint) => touchpoint.workflowDefinitionKey !== null).length,
      workflowDispatched: touchpoints.filter((touchpoint) => touchpoint.lastWorkflowRunPublicId !== null).length,
      businessLinked: touchpoints.filter((touchpoint) => touchpoint.businessEntityPublicId !== null).length
    },
    byStatus,
    byChannel,
    outcomes: {
      responded,
      converted,
      failed,
      responseRate: sentBase > 0 ? Number((responded / sentBase).toFixed(4)) : 0,
      conversionRate: responseBase > 0 ? Number((converted / responseBase).toFixed(4)) : 0
    }
  };
}

function ensureIncluded<T extends string>(value: string, items: readonly T[], errorCode: string): T {
  const normalizedValue = value.trim().toLowerCase();

  if (!items.includes(normalizedValue as T)) {
    throw new Error(errorCode);
  }

  return normalizedValue as T;
}

export function ensureTouchpointStatus(value: string): TouchpointStatus {
  return ensureIncluded(value, touchpointStatuses, "touchpoint_status_invalid");
}

export function ensureTouchpointChannel(value: string): CampaignChannel {
  return ensureIncluded(value, campaignChannels, "touchpoint_channel_invalid");
}

export function ensurePublicId(value: string, errorCode: string): string {
  const normalizedValue = value.trim();

  if (!/^[0-9a-fA-F-]{36}$/.test(normalizedValue)) {
    throw new Error(errorCode);
  }

  return normalizedValue.toLowerCase();
}

export function ensureOptionalPublicId(value: string | null | undefined): string | null {
  if (value == null || value.trim().length === 0) {
    return null;
  }

  return ensurePublicId(value, "touchpoint_workflow_run_public_id_invalid");
}

export function ensureTouchpointText(value: string, errorCode: string): string {
  return ensureText(value, errorCode);
}

export function ensureBusinessEntityType(value: string | null | undefined): string | null {
  const normalizedValue = (value ?? "").trim().toLowerCase();

  if (normalizedValue.length === 0) {
    return null;
  }

  if (!/^[a-z0-9_.-]{3,80}$/.test(normalizedValue)) {
    throw new Error("business_entity_type_invalid");
  }

  return normalizedValue;
}

export function ensureBusinessEntityPublicId(value: string | null | undefined): string | null {
  if (value == null || value.trim().length === 0) {
    return null;
  }

  return ensurePublicId(value, "business_entity_public_id_invalid");
}
