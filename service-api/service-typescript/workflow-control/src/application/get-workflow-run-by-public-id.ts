import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";
import { WorkflowRun } from "../domain/workflow-run.js";

export class GetWorkflowRunByPublicId {
  public constructor(
    private readonly repository: WorkflowRunRepository
  ) {}

  public execute(publicId: string): Promise<WorkflowRun | null> {
    return this.repository.findByPublicId(publicId);
  }
}
