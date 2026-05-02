export type DispatchWorkflowTouchpointRequest = {
  tenantSlug: string;
  campaignPublicId: string;
  templatePublicId?: string | null;
  leadPublicId: string;
  contactValue: string;
  provider: string;
  workflowRunPublicId: string;
  providerMessageId?: string | null;
  createdBy: string;
  notes?: string;
};
