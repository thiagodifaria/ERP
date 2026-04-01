import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { WorkflowDefinitionVersion } from "../domain/workflow-definition-version.js";

export class ListWorkflowDefinitionVersions {
  public constructor(
    private readonly definitionRepository: WorkflowDefinitionRepository,
    private readonly versionRepository: WorkflowDefinitionVersionRepository
  ) {}

  public async execute(definitionKey: string): Promise<WorkflowDefinitionVersion[]> {
    const definition = await this.definitionRepository.findByKey(definitionKey);

    if (definition === null) {
      throw new Error("workflow_definition_not_found");
    }

    return this.versionRepository.listByWorkflowDefinitionId(definition.id);
  }
}
