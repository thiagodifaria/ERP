import { WorkflowDefinition } from "../domain/workflow-definition.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";

export class ListWorkflowDefinitions {
  public constructor(
    private readonly repository: InMemoryWorkflowDefinitionRepository
  ) {}

  public execute(): WorkflowDefinition[] {
    return this.repository.list();
  }
}
