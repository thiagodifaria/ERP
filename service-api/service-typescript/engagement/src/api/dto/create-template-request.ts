import { CampaignChannel } from "../../domain/campaign.js";
import { TemplateProvider, TemplateStatus } from "../../domain/template.js";

export type CreateTemplateRequest = {
  tenantSlug: string;
  key: string;
  name: string;
  channel: CampaignChannel;
  status?: TemplateStatus;
  provider: TemplateProvider;
  subject?: string | null;
  body: string;
};
