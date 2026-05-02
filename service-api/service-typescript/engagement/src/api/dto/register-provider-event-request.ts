export type RegisterProviderEventRequest = {
  tenantSlug: string;
  provider: string;
  eventType: string;
  externalEventId?: string;
  touchpointPublicId?: string | null;
  deliveryPublicId?: string | null;
  workflowRunPublicId?: string | null;
  leadPublicId?: string | null;
  providerMessageId?: string | null;
  payloadSummary?: string;
  responseSummary?: string;
};
