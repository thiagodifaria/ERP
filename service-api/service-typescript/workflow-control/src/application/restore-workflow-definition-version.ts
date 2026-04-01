import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { WorkflowDefinition, createWorkflowDefinition } from "../domain/workflow-definition.js";

export class RestoreWorkflowDefinitionVersion {
  public constructor(
    private readonly definitionRepository: WorkflowDefinitionRepository,
    private readonly versionRepository: WorkflowDefinitionVersionRepository
  ) {}

  public async execute(definitionKey: string, versionNumber: number): Promise<WorkflowDefinition> {
    const definition = await this.definitionRepository.findByKey(definitionKey);

    if (definition === null) {
      throw new Error("workflow_definition_not_found");
    }

    const version = await this.versionRepository.findByWorkflowDefinitionIdAndVersionNumber(definition.id, versionNumber);

    if (version === null) {
      throw new Error("workflow_definition_version_not_found");
    }

    const restoredDefinition = createWorkflowDefinition({
      id: definition.id,
      key: definition.key,
      name: version.snapshotName,
      description: version.snapshotDescription,
      status: version.snapshotStatus,
      trigger: version.snapshotTrigger
    });

    await this.definitionRepository.updateDefinition(definition.key, restoredDefinition);
    await this.definitionRepository.updateStatus(definition.key, restoredDefinition.status);

    return (await this.definitionRepository.findByKey(definition.key)) ?? restoredDefinition;
  }
}
