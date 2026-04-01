import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";

export type WorkflowRunEventSummary = {
  workflowRunPublicId: string;
  total: number;
  byCategory: {
    status: number;
    note: number;
  };
  latestEventPublicId: string | null;
  latestCategory: string | null;
  latestCreatedAt: string | null;
};

export class GetWorkflowRunEventSummary {
  public constructor(
    private readonly runRepository: WorkflowRunRepository,
    private readonly runEventRepository: WorkflowRunEventRepository
  ) {}

  public async execute(workflowRunPublicId: string): Promise<WorkflowRunEventSummary> {
    const workflowRun = await this.runRepository.findByPublicId(workflowRunPublicId);

    if (workflowRun === null) {
      throw new Error("workflow_run_not_found");
    }

    const workflowRunEvents = await this.runEventRepository.listByWorkflowRunPublicId(workflowRunPublicId);
    const latestEvent = workflowRunEvents.at(-1) ?? null;

    return {
      workflowRunPublicId: workflowRun.publicId,
      total: workflowRunEvents.length,
      byCategory: {
        status: workflowRunEvents.filter((workflowRunEvent) => workflowRunEvent.category === "status").length,
        note: workflowRunEvents.filter((workflowRunEvent) => workflowRunEvent.category === "note").length
      },
      latestEventPublicId: latestEvent?.publicId ?? null,
      latestCategory: latestEvent?.category ?? null,
      latestCreatedAt: latestEvent?.createdAt ?? null
    };
  }
}
