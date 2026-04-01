import { WorkflowDefinition } from "../domain/workflow-definition.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { WorkflowDefinitionVersion, createWorkflowDefinitionVersion } from "../domain/workflow-definition-version.js";

export class InMemoryWorkflowDefinitionVersionRepository implements WorkflowDefinitionVersionRepository {
  private readonly versions: WorkflowDefinitionVersion[];

  public constructor() {
    this.versions = [
      createWorkflowDefinitionVersion({
        id: 1,
        workflowDefinitionId: 1,
        versionNumber: 1,
        snapshotName: "Lead Follow-Up",
        snapshotDescription: "Orquestra o acompanhamento inicial de novos leads do CRM.",
        snapshotStatus: "active",
        snapshotTrigger: "lead.created"
      })
    ];
  }

  public async listByDefinitionKey(definitionKey: string): Promise<WorkflowDefinitionVersion[]> {
    const workflowDefinitionId = this.resolveDefinitionId(definitionKey);

    if (workflowDefinitionId === null) {
      return [];
    }

    return this.versions
      .filter((version) => version.workflowDefinitionId === workflowDefinitionId)
      .sort((left, right) => right.versionNumber - left.versionNumber);
  }

  public async findCurrentByDefinitionKey(definitionKey: string): Promise<WorkflowDefinitionVersion | null> {
    const versions = await this.listByDefinitionKey(definitionKey);
    return versions[0] ?? null;
  }

  public async publish(definition: WorkflowDefinition): Promise<WorkflowDefinitionVersion> {
    const definitionVersions = this.versions.filter((version) => version.workflowDefinitionId === definition.id);
    const nextVersionNumber = definitionVersions.reduce((max, version) => Math.max(max, version.versionNumber), 0) + 1;
    const nextId = this.versions.reduce((max, version) => Math.max(max, version.id), 0) + 1;

    const createdVersion = createWorkflowDefinitionVersion({
      id: nextId,
      workflowDefinitionId: definition.id,
      versionNumber: nextVersionNumber,
      snapshotName: definition.name,
      snapshotDescription: definition.description,
      snapshotStatus: definition.status,
      snapshotTrigger: definition.trigger
    });

    this.versions.push(createdVersion);
    return createdVersion;
  }

  private resolveDefinitionId(definitionKey: string): number | null {
    const normalizedKey = definitionKey.trim().toLowerCase();

    if (normalizedKey === "lead-follow-up") {
      return 1;
    }

    const version = this.versions.find((candidate) => candidate.workflowDefinitionId > 1 && candidate.snapshotName.length > 0 && normalizedKey.length > 0);
    return version?.workflowDefinitionId ?? null;
  }
}
