import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRunEvent } from "../domain/workflow-run-event.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";

export class ListWorkflowRunEvents {
  public constructor(
    private readonly runRepository: WorkflowRunRepository,
    private readonly eventRepository: WorkflowRunEventRepository
  ) {}

  public async execute(workflowRunPublicId: string): Promise<WorkflowRunEvent[]> {
    const workflowRun = await this.runRepository.findByPublicId(workflowRunPublicId);

    if (workflowRun === null) {
      throw new Error("workflow_run_not_found");
    }

    return this.eventRepository.listByWorkflowRunPublicId(workflowRunPublicId);
  }
}
