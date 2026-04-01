import { randomUUID } from "node:crypto";
import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";
import { WorkflowRunEvent, createWorkflowRunEvent } from "../domain/workflow-run-event.js";

export class CreateWorkflowRunNote {
  public constructor(
    private readonly runRepository: WorkflowRunRepository,
    private readonly runEventRepository: WorkflowRunEventRepository
  ) {}

  public async execute(input: {
    workflowRunPublicId: string;
    body: string;
    createdBy: string;
  }): Promise<WorkflowRunEvent> {
    const workflowRunPublicId = input.workflowRunPublicId.trim().toLowerCase();

    if (workflowRunPublicId.length === 0) {
      throw new Error("workflow_run_public_id_required");
    }

    const workflowRun = await this.runRepository.findByPublicId(workflowRunPublicId);

    if (workflowRun === null) {
      throw new Error("workflow_run_not_found");
    }

    const workflowRunEvent = createWorkflowRunEvent({
      id: await this.runEventRepository.nextId(),
      publicId: randomUUID(),
      workflowRunPublicId,
      category: "note",
      body: input.body,
      createdBy: input.createdBy,
      createdAt: new Date().toISOString()
    });

    return this.runEventRepository.add(workflowRunEvent);
  }
}
