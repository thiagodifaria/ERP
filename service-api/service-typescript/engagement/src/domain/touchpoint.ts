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
};

export type CreateTouchpointInput = {
  tenantSlug: string;
  campaignPublicId: string;
  leadPublicId: string;
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
