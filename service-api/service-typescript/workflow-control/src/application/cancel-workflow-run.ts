import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";
import { WorkflowRun } from "../domain/workflow-run.js";

export class CancelWorkflowRun {
  public constructor(
    private readonly repository: WorkflowRunRepository
  ) {}

  public async execute(publicId: string): Promise<WorkflowRun> {
    const workflowRun = await this.repository.findByPublicId(publicId);

    if (workflowRun === null) {
      throw new Error("workflow_run_not_found");
    }

    if (workflowRun.status !== "pending" && workflowRun.status !== "running") {
      throw new Error("workflow_run_transition_invalid");
    }

    return this.repository.updateStatus(publicId, "cancelled", {
      cancelledAt: new Date().toISOString()
    });
  }
}
