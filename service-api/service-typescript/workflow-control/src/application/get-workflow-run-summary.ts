import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";

export type WorkflowRunSummary = {
  total: number;
  pending: number;
  running: number;
  completed: number;
  failed: number;
  cancelled: number;
};

export class GetWorkflowRunSummary {
  public constructor(
    private readonly repository: WorkflowRunRepository
  ) {}

  public async execute(): Promise<WorkflowRunSummary> {
    const runs = await this.repository.list();

    return {
      total: runs.length,
      pending: runs.filter((run) => run.status === "pending").length,
      running: runs.filter((run) => run.status === "running").length,
      completed: runs.filter((run) => run.status === "completed").length,
      failed: runs.filter((run) => run.status === "failed").length,
      cancelled: runs.filter((run) => run.status === "cancelled").length
    };
  }
}
