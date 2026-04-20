import { TouchpointStatus } from "../../domain/touchpoint.js";

export type UpdateTouchpointStatusRequest = {
  status: TouchpointStatus;
  lastWorkflowRunPublicId?: string | null;
};
