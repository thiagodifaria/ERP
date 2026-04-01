import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";

export type WorkflowDefinitionVersionSummary = {
  workflowDefinitionId: number;
  totalVersions: number;
  currentVersionNumber: number | null;
  currentSnapshotStatus: string | null;
};

export class GetWorkflowDefinitionVersionSummary {
  public constructor(
    private readonly definitionRepository: WorkflowDefinitionRepository,
    private readonly versionRepository: WorkflowDefinitionVersionRepository
  ) {}

  public async execute(definitionKey: string): Promise<WorkflowDefinitionVersionSummary> {
    const definition = await this.definitionRepository.findByKey(definitionKey);

    if (definition === null) {
      throw new Error("workflow_definition_not_found");
    }

    const versions = await this.versionRepository.listByWorkflowDefinitionId(definition.id);
    const currentVersion = await this.versionRepository.findCurrentByWorkflowDefinitionId(definition.id);

    return {
      workflowDefinitionId: definition.id,
      totalVersions: versions.length,
      currentVersionNumber: currentVersion?.versionNumber ?? null,
      currentSnapshotStatus: currentVersion?.snapshotStatus ?? null
    };
  }
}
