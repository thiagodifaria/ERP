import { randomUUID } from "node:crypto";
import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRun, WorkflowRunStatus } from "../domain/workflow-run.js";
import { createWorkflowRunEvent } from "../domain/workflow-run-event.js";

const statusEventBodies: Record<WorkflowRunStatus, string> = {
  pending: "Workflow run moved to pending.",
  running: "Workflow run moved to running.",
  completed: "Workflow run moved to completed.",
  failed: "Workflow run moved to failed.",
  cancelled: "Workflow run moved to cancelled."
};

export async function appendWorkflowRunStatusEvent(
  repository: WorkflowRunEventRepository,
  workflowRun: WorkflowRun
): Promise<void> {
  if (workflowRun.status === "pending") {
    return;
  }

  const workflowRunEvent = createWorkflowRunEvent({
    id: await repository.nextId(),
    publicId: randomUUID(),
    workflowRunPublicId: workflowRun.publicId,
    category: "status",
    body: statusEventBodies[workflowRun.status],
    createdBy: "workflow-control",
    createdAt: new Date().toISOString()
  });

  await repository.add(workflowRunEvent);
}
