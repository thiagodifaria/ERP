import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";
import { WorkflowRun } from "../domain/workflow-run.js";

export class ListWorkflowRuns {
  public constructor(
    private readonly repository: WorkflowRunRepository
  ) {}

  public execute(): Promise<WorkflowRun[]> {
    return this.repository.list();
  }
}
