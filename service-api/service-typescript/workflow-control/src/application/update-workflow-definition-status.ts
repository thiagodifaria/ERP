import { WorkflowDefinition, ensureWorkflowDefinitionStatus } from "../domain/workflow-definition.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class UpdateWorkflowDefinitionStatus {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository
  ) {}

  public execute(key: string, status: string): WorkflowDefinition {
    const definition = this.repository.findByKey(key);

    if (definition === null) {
      throw new Error("workflow_definition_not_found");
    }

    return this.repository.updateStatus(definition.key, ensureWorkflowDefinitionStatus(status));
  }
}
