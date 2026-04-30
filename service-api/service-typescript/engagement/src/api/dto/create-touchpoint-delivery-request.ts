import { TemplateProvider } from "../../domain/template.js";

export type CreateTouchpointDeliveryRequest = {
  tenantSlug: string;
  templatePublicId?: string | null;
  provider: TemplateProvider;
  providerMessageId?: string | null;
  sentBy: string;
  notes?: string;
};
