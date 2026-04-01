export type WorkflowRunEventCategory = "status" | "note";

export type WorkflowRunEvent = {
  id: number;
  publicId: string;
  workflowRunPublicId: string;
  category: WorkflowRunEventCategory;
  body: string;
  createdBy: string;
  createdAt: string;
};

export function ensureWorkflowRunEventCategory(value: string): WorkflowRunEventCategory {
  const normalizedValue = value.trim().toLowerCase();

  if (normalizedValue !== "status" && normalizedValue !== "note") {
    throw new Error("workflow_run_event_category_invalid");
  }

  return normalizedValue;
}

export function createWorkflowRunEvent(input: {
  id: number;
  publicId: string;
  workflowRunPublicId: string;
  category: WorkflowRunEventCategory;
  body: string;
  createdBy: string;
  createdAt: string;
}): WorkflowRunEvent {
  const publicId = input.publicId.trim().toLowerCase();
  const workflowRunPublicId = input.workflowRunPublicId.trim().toLowerCase();
  const body = input.body.trim();
  const createdBy = input.createdBy.trim().toLowerCase();
  const createdAt = input.createdAt.trim();

  if (publicId.length === 0) {
    throw new Error("workflow_run_event_public_id_required");
  }

  if (workflowRunPublicId.length === 0) {
    throw new Error("workflow_run_event_workflow_run_public_id_required");
  }

  if (body.length === 0) {
    throw new Error("workflow_run_event_body_required");
  }

  if (createdBy.length === 0) {
    throw new Error("workflow_run_event_created_by_required");
  }

  if (createdAt.length === 0) {
    throw new Error("workflow_run_event_created_at_required");
  }

  return {
    id: input.id,
    publicId,
    workflowRunPublicId,
    category: ensureWorkflowRunEventCategory(input.category),
    body,
    createdBy,
    createdAt
  };
}
