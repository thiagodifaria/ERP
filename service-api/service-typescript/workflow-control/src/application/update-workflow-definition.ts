import { WorkflowDefinition, WorkflowDefinitionUpdateInput, applyWorkflowDefinitionUpdate } from "../domain/workflow-definition.js";
import { WorkflowCatalogRepository } from "../domain/workflow-catalog.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class UpdateWorkflowDefinition {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository,
    private readonly catalogRepository: WorkflowCatalogRepository
  ) {}

  public async execute(key: string, input: WorkflowDefinitionUpdateInput): Promise<WorkflowDefinition> {
    const currentDefinition = await this.repository.findByKey(key);

    if (currentDefinition === null) {
      throw new Error("workflow_definition_not_found");
    }

    if (input.trigger !== undefined && !(await this.catalogRepository.hasTrigger(input.trigger))) {
      throw new Error("workflow_definition_trigger_unknown");
    }

    const updatedDefinition = applyWorkflowDefinitionUpdate(currentDefinition, input);
    return this.repository.updateDefinition(currentDefinition.key, updatedDefinition);
  }
}
