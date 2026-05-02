export type IngestProviderLeadRequest = {
  tenantSlug: string;
  provider: string;
  campaignPublicId?: string;
  externalEventId?: string;
  name: string;
  email: string;
  contactValue?: string;
  notes?: string;
};
