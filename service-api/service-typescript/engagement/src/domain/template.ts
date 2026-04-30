import { CampaignChannel, campaignChannels, ensureText } from "./campaign.js";

export const templateStatuses = ["draft", "active", "archived"] as const;
export const templateProviders = ["resend", "whatsapp_cloud", "telegram_bot", "manual"] as const;

export type TemplateStatus = (typeof templateStatuses)[number];
export type TemplateProvider = (typeof templateProviders)[number];

export type EngagementTemplate = {
  id: number;
  publicId: string;
  tenantSlug: string;
  key: string;
  name: string;
  channel: CampaignChannel;
  status: TemplateStatus;
  provider: TemplateProvider;
  subject: string | null;
  body: string;
  createdAt: string;
  updatedAt: string;
};

export type TemplateFilters = {
  tenantSlug?: string;
  channel?: CampaignChannel;
  status?: TemplateStatus;
  provider?: TemplateProvider;
  q?: string;
};

export type CreateTemplateInput = {
  tenantSlug: string;
  key: string;
  name: string;
  channel: CampaignChannel;
  status?: TemplateStatus;
  provider: TemplateProvider;
  subject?: string | null;
  body: string;
};

function ensureIncluded<T extends string>(value: string, items: readonly T[], errorCode: string): T {
  const normalizedValue = value.trim().toLowerCase();

  if (!items.includes(normalizedValue as T)) {
    throw new Error(errorCode);
  }

  return normalizedValue as T;
}

export function ensureTemplateStatus(value: string): TemplateStatus {
  return ensureIncluded(value, templateStatuses, "template_status_invalid");
}

export function ensureTemplateProvider(value: string): TemplateProvider {
  return ensureIncluded(value, templateProviders, "template_provider_invalid");
}

export function ensureTemplateKey(value: string): string {
  const normalizedValue = value.trim().toLowerCase();

  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(normalizedValue)) {
    throw new Error("template_key_invalid");
  }

  return normalizedValue;
}

export function ensureTemplateSubject(value: string | null | undefined): string | null {
  if (value == null) {
    return null;
  }

  const normalizedValue = value.trim();
  return normalizedValue.length > 0 ? normalizedValue : null;
}

export function ensureTemplateBody(value: string): string {
  return ensureText(value, "template_body_required");
}

export function ensureTemplateName(value: string): string {
  return ensureText(value, "template_name_required");
}

export function ensureTemplateChannel(value: string): CampaignChannel {
  return ensureIncluded(value, campaignChannels, "template_channel_invalid");
}
