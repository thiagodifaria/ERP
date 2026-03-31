import { WorkflowDefinition } from "../domain/workflow-definition.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class GetWorkflowDefinitionByKey {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository
  ) {}

  public execute(key: string): Promise<WorkflowDefinition | null> {
    return this.repository.findByKey(key);
  }
}
