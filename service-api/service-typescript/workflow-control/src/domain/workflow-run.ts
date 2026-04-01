export type WorkflowRunStatus = "pending" | "running" | "completed" | "failed" | "cancelled";

export type WorkflowRun = {
  id: number;
  publicId: string;
  workflowDefinitionId: number;
  workflowDefinitionVersionId: number;
  status: WorkflowRunStatus;
  triggerEvent: string;
  subjectType: string;
  subjectPublicId: string;
  initiatedBy: string;
  startedAt: string | null;
  completedAt: string | null;
  failedAt: string | null;
  cancelledAt: string | null;
};

export function ensureWorkflowRunStatus(value: string): WorkflowRunStatus {
  const normalizedValue = value.trim().toLowerCase();

  if (
    normalizedValue !== "pending" &&
    normalizedValue !== "running" &&
    normalizedValue !== "completed" &&
    normalizedValue !== "failed" &&
    normalizedValue !== "cancelled"
  ) {
    throw new Error("workflow_run_status_invalid");
  }

  return normalizedValue;
}

export function createWorkflowRun(input: {
  id: number;
  publicId: string;
  workflowDefinitionId: number;
  workflowDefinitionVersionId: number;
  status?: WorkflowRunStatus;
  triggerEvent: string;
  subjectType: string;
  subjectPublicId: string;
  initiatedBy: string;
  startedAt?: string | null;
  completedAt?: string | null;
  failedAt?: string | null;
  cancelledAt?: string | null;
}): WorkflowRun {
  const publicId = input.publicId.trim().toLowerCase();
  const triggerEvent = input.triggerEvent.trim().toLowerCase();
  const subjectType = input.subjectType.trim().toLowerCase();
  const subjectPublicId = input.subjectPublicId.trim().toLowerCase();
  const initiatedBy = input.initiatedBy.trim().toLowerCase();

  if (publicId.length === 0) {
    throw new Error("workflow_run_public_id_required");
  }

  if (triggerEvent.length === 0) {
    throw new Error("workflow_run_trigger_event_required");
  }

  if (subjectType.length === 0) {
    throw new Error("workflow_run_subject_type_required");
  }

  if (subjectPublicId.length === 0) {
    throw new Error("workflow_run_subject_public_id_required");
  }

  if (initiatedBy.length === 0) {
    throw new Error("workflow_run_initiated_by_required");
  }

  return {
    id: input.id,
    publicId,
    workflowDefinitionId: input.workflowDefinitionId,
    workflowDefinitionVersionId: input.workflowDefinitionVersionId,
    status: ensureWorkflowRunStatus(input.status ?? "pending"),
    triggerEvent,
    subjectType,
    subjectPublicId,
    initiatedBy,
    startedAt: input.startedAt ?? null,
    completedAt: input.completedAt ?? null,
    failedAt: input.failedAt ?? null,
    cancelledAt: input.cancelledAt ?? null
  };
}
