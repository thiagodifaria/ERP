import { WorkflowDefinition } from "../domain/workflow-definition.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class ListWorkflowDefinitions {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository
  ) {}

  public execute(): Promise<WorkflowDefinition[]> {
    return this.repository.list();
  }
}
