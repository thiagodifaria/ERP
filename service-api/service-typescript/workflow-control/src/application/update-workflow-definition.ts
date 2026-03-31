import { WorkflowDefinition, WorkflowDefinitionUpdateInput, applyWorkflowDefinitionUpdate } from "../domain/workflow-definition.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class UpdateWorkflowDefinition {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository
  ) {}

  public async execute(key: string, input: WorkflowDefinitionUpdateInput): Promise<WorkflowDefinition> {
    const currentDefinition = await this.repository.findByKey(key);

    if (currentDefinition === null) {
      throw new Error("workflow_definition_not_found");
    }

    const updatedDefinition = applyWorkflowDefinitionUpdate(currentDefinition, input);
    return this.repository.updateDefinition(currentDefinition.key, updatedDefinition);
  }
}
