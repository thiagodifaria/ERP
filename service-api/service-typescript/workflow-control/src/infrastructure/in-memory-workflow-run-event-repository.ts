import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRunEvent, createWorkflowRunEvent } from "../domain/workflow-run-event.js";

export class InMemoryWorkflowRunEventRepository implements WorkflowRunEventRepository {
  private readonly events: WorkflowRunEvent[];

  public constructor() {
    this.events = [
      createWorkflowRunEvent({
        id: 1,
        publicId: "00000000-0000-0000-0000-000000000501",
        workflowRunPublicId: "00000000-0000-0000-0000-000000000301",
        category: "note",
        body: "Execucao bootstrap criada para acompanhamento inicial do fluxo.",
        createdBy: "bootstrap-seed",
        createdAt: "2026-04-01T08:02:00.000Z"
      })
    ];
  }

  public async listByWorkflowRunPublicId(workflowRunPublicId: string): Promise<WorkflowRunEvent[]> {
    const normalizedWorkflowRunPublicId = workflowRunPublicId.trim().toLowerCase();

    return this.events
      .filter((event) => event.workflowRunPublicId === normalizedWorkflowRunPublicId)
      .sort((left, right) => left.id - right.id);
  }

  public async add(event: WorkflowRunEvent): Promise<WorkflowRunEvent> {
    this.events.push(event);
    return event;
  }

  public async nextId(): Promise<number> {
    return this.events.reduce((max, event) => Math.max(max, event.id), 0) + 1;
  }
}
