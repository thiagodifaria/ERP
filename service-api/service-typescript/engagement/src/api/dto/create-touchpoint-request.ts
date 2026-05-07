export type CreateTouchpointRequest = {
  tenantSlug: string;
  campaignPublicId: string;
  leadPublicId: string;
  threadPublicId?: string | null;
  participantKind?: string | null;
  participantPublicId?: string | null;
  businessEntityType?: string | null;
  businessEntityPublicId?: string | null;
  contactValue: string;
  source: string;
  createdBy: string;
  notes: string;
};
