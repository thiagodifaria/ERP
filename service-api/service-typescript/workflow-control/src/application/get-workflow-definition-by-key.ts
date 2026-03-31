import { WorkflowDefinition } from "../domain/workflow-definition.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";

export class GetWorkflowDefinitionByKey {
  public constructor(
    private readonly repository: InMemoryWorkflowDefinitionRepository
  ) {}

  public execute(key: string): WorkflowDefinition | null {
    return this.repository.findByKey(key);
  }
}
