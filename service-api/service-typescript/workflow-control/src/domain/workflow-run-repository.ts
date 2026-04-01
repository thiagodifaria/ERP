import { WorkflowRun, WorkflowRunStatus } from "./workflow-run.js";

export type WorkflowRunListFilters = {
  workflowDefinitionId?: number;
  status?: WorkflowRunStatus;
  triggerEvent?: string;
  initiatedBy?: string;
  subjectType?: string;
  subjectPublicId?: string;
};

export interface WorkflowRunRepository {
  list(filters?: WorkflowRunListFilters): Promise<WorkflowRun[]>;
  findByPublicId(publicId: string): Promise<WorkflowRun | null>;
  add(run: WorkflowRun): Promise<WorkflowRun>;
  updateStatus(
    publicId: string,
    status: WorkflowRunStatus,
    timestamps?: Partial<Pick<WorkflowRun, "startedAt" | "completedAt" | "failedAt" | "cancelledAt">>
  ): Promise<WorkflowRun>;
  nextId(): Promise<number>;
}
