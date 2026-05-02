import { CreateTouchpointInput, Touchpoint, TouchpointFilters, TouchpointStatus, TouchpointSummary } from "./touchpoint.js";

export interface TouchpointRepository {
  list(filters?: TouchpointFilters): Promise<Touchpoint[]>;
  getByPublicId(publicId: string): Promise<Touchpoint | null>;
  create(
    input: CreateTouchpointInput & {
      campaignKey: string;
      channel: Touchpoint["channel"];
      workflowDefinitionKey: string | null;
      businessEntityType: string | null;
      businessEntityPublicId: string | null;
    }
  ): Promise<Touchpoint>;
  updateStatus(
    publicId: string,
    status: TouchpointStatus,
    lastWorkflowRunPublicId?: string | null
  ): Promise<Touchpoint | null>;
  getSummary(filters?: TouchpointFilters): Promise<TouchpointSummary>;
}
