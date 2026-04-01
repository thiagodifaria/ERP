import { WorkflowRunEvent } from "./workflow-run-event.js";

export interface WorkflowRunEventRepository {
  listByWorkflowRunPublicId(workflowRunPublicId: string): Promise<WorkflowRunEvent[]>;
  add(event: WorkflowRunEvent): Promise<WorkflowRunEvent>;
  nextId(): Promise<number>;
}
