export type CreateTouchpointRequest = {
  tenantSlug: string;
  campaignPublicId: string;
  leadPublicId: string;
  contactValue: string;
  source: string;
  createdBy: string;
  notes: string;
};
