export type CreateTouchpointRequest = {
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
