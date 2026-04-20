export const campaignChannels = ["whatsapp", "email", "telegram", "meta_ads", "manual"] as const;
export const campaignStatuses = ["draft", "active", "paused", "archived"] as const;

export type CampaignChannel = (typeof campaignChannels)[number];
export type CampaignStatus = (typeof campaignStatuses)[number];

export type Campaign = {
  id: number;
  publicId: string;
  tenantSlug: string;
  key: string;
  name: string;
  description: string;
  channel: CampaignChannel;
  status: CampaignStatus;
  touchpointGoal: string;
  workflowDefinitionKey: string | null;
  budgetCents: number;
  createdAt: string;
  updatedAt: string;
};

export type CampaignFilters = {
  tenantSlug?: string;
  status?: CampaignStatus;
  channel?: CampaignChannel;
  q?: string;
};

export type CreateCampaignInput = {
  tenantSlug: string;
  key: string;
  name: string;
  description: string;
  channel: CampaignChannel;
  status?: CampaignStatus;
  touchpointGoal: string;
  workflowDefinitionKey?: string | null;
  budgetCents: number;
};

function ensureIncluded<T extends string>(value: string, items: readonly T[], errorCode: string): T {
  const normalizedValue = value.trim().toLowerCase();

  if (!items.includes(normalizedValue as T)) {
    throw new Error(errorCode);
  }

  return normalizedValue as T;
}

export function ensureCampaignChannel(value: string): CampaignChannel {
  return ensureIncluded(value, campaignChannels, "campaign_channel_invalid");
}

export function ensureCampaignStatus(value: string): CampaignStatus {
  return ensureIncluded(value, campaignStatuses, "campaign_status_invalid");
}

export function ensureCampaignKey(value: string): string {
  const normalizedValue = value.trim().toLowerCase();

  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(normalizedValue)) {
    throw new Error("campaign_key_invalid");
  }

  return normalizedValue;
}

export function ensureText(value: string, errorCode: string): string {
  const normalizedValue = value.trim();

  if (normalizedValue.length === 0) {
    throw new Error(errorCode);
  }

  return normalizedValue;
}

export function ensureBudgetCents(value: number): number {
  if (!Number.isInteger(value) || value < 0) {
    throw new Error("campaign_budget_invalid");
  }

  return value;
}
