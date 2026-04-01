import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRunEvent, ensureWorkflowRunEventCategory } from "../domain/workflow-run-event.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";

export class ListWorkflowRunEvents {
  public constructor(
    private readonly runRepository: WorkflowRunRepository,
    private readonly eventRepository: WorkflowRunEventRepository
  ) {}

  public async execute(
    workflowRunPublicId: string,
    filters?: {
      category?: string;
      createdBy?: string;
    }
  ): Promise<WorkflowRunEvent[]> {
    const workflowRun = await this.runRepository.findByPublicId(workflowRunPublicId);

    if (workflowRun === null) {
      throw new Error("workflow_run_not_found");
    }

    const workflowRunEvents = await this.eventRepository.listByWorkflowRunPublicId(workflowRunPublicId);
    const category = filters?.category ? ensureWorkflowRunEventCategory(filters.category) : undefined;
    const createdBy = filters?.createdBy?.trim().toLowerCase();

    return workflowRunEvents.filter((workflowRunEvent) => {
      if (category && workflowRunEvent.category !== category) {
        return false;
      }

      if (createdBy && workflowRunEvent.createdBy !== createdBy) {
        return false;
      }

      return true;
    });
  }
}
