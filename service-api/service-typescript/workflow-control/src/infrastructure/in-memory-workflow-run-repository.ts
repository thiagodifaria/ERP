import { WorkflowRunRepository, WorkflowRunListFilters } from "../domain/workflow-run-repository.js";
import { WorkflowRun, WorkflowRunStatus, createWorkflowRun } from "../domain/workflow-run.js";

export class InMemoryWorkflowRunRepository implements WorkflowRunRepository {
  private readonly runs: WorkflowRun[];

  public constructor() {
    this.runs = [
      createWorkflowRun({
        id: 1,
        publicId: "00000000-0000-0000-0000-000000000301",
        workflowDefinitionId: 1,
        workflowDefinitionVersionId: 1,
        status: "running",
        triggerEvent: "lead.created",
        subjectType: "crm.lead",
        subjectPublicId: "00000000-0000-0000-0000-000000000401",
        initiatedBy: "bootstrap-seed",
        startedAt: "2026-04-01T08:10:00.000Z"
      })
    ];
  }

  public async list(filters?: WorkflowRunListFilters): Promise<WorkflowRun[]> {
    return this.runs
      .filter((run) => this.matchesFilters(run, filters))
      .sort((left, right) => right.id - left.id);
  }

  public async findByPublicId(publicId: string): Promise<WorkflowRun | null> {
    const normalizedPublicId = publicId.trim().toLowerCase();
    return this.runs.find((run) => run.publicId === normalizedPublicId) ?? null;
  }

  public async add(run: WorkflowRun): Promise<WorkflowRun> {
    this.runs.push(run);
    return run;
  }

  public async updateStatus(
    publicId: string,
    status: WorkflowRunStatus,
    timestamps?: Partial<Pick<WorkflowRun, "startedAt" | "completedAt" | "failedAt" | "cancelledAt">>
  ): Promise<WorkflowRun> {
    const run = await this.findByPublicId(publicId);

    if (run === null) {
      throw new Error("workflow_run_not_found");
    }

    run.status = status;
    run.startedAt = timestamps?.startedAt ?? run.startedAt;
    run.completedAt = timestamps?.completedAt ?? run.completedAt;
    run.failedAt = timestamps?.failedAt ?? run.failedAt;
    run.cancelledAt = timestamps?.cancelledAt ?? run.cancelledAt;

    return run;
  }

  public async nextId(): Promise<number> {
    return this.runs.reduce((max, run) => Math.max(max, run.id), 0) + 1;
  }

  private matchesFilters(run: WorkflowRun, filters?: WorkflowRunListFilters): boolean {
    if (!filters) {
      return true;
    }

    if (filters.workflowDefinitionId !== undefined && run.workflowDefinitionId !== filters.workflowDefinitionId) {
      return false;
    }

    if (filters.status !== undefined && run.status !== filters.status) {
      return false;
    }

    if (filters.triggerEvent !== undefined && run.triggerEvent !== filters.triggerEvent) {
      return false;
    }

    if (filters.initiatedBy !== undefined && run.initiatedBy !== filters.initiatedBy) {
      return false;
    }

    if (filters.subjectType !== undefined && run.subjectType !== filters.subjectType) {
      return false;
    }

    if (filters.subjectPublicId !== undefined && run.subjectPublicId !== filters.subjectPublicId) {
      return false;
    }

    return true;
  }
}
