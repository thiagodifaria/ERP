import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowRun, ensureWorkflowRunStatus } from "../domain/workflow-run.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";

export class ListWorkflowRuns {
  public constructor(
    private readonly definitionRepository: WorkflowDefinitionRepository,
    private readonly repository: WorkflowRunRepository
  ) {}

  public async execute(filters?: {
    workflowDefinitionKey?: string;
    status?: string;
    subjectType?: string;
    initiatedBy?: string;
  }): Promise<WorkflowRun[]> {
    let workflowDefinitionId: number | undefined;

    if (filters?.workflowDefinitionKey) {
      const definition = await this.definitionRepository.findByKey(filters.workflowDefinitionKey);

      if (definition === null) {
        return [];
      }

      workflowDefinitionId = definition.id;
    }

    return this.repository.list({
      workflowDefinitionId,
      status: filters?.status ? ensureWorkflowRunStatus(filters.status) : undefined,
      subjectType: filters?.subjectType?.trim().toLowerCase() || undefined,
      initiatedBy: filters?.initiatedBy?.trim().toLowerCase() || undefined
    });
  }
}
