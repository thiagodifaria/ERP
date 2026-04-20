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
        snapshotTrigger: "lead.created",
        snapshotActions: [
          {
            stepId: "create-task",
            actionKey: "task.create",
            label: "Criar tarefa comercial inicial",
            delaySeconds: null,
            compensationActionKey: "task.create"
          },
          {
            stepId: "notify-webhook",
            actionKey: "integration.webhook",
            label: "Emitir webhook operacional",
            delaySeconds: null,
            compensationActionKey: "integration.webhook"
          }
        ]
      })
    ];
  }

  public async listByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion[]> {
    return this.versions
      .filter((version) => version.workflowDefinitionId === workflowDefinitionId)
      .sort((left, right) => right.versionNumber - left.versionNumber);
  }

  public async findByWorkflowDefinitionIdAndVersionNumber(
    workflowDefinitionId: number,
    versionNumber: number
  ): Promise<WorkflowDefinitionVersion | null> {
    return this.versions.find((version) => (
      version.workflowDefinitionId === workflowDefinitionId &&
      version.versionNumber === versionNumber
    )) ?? null;
  }

  public async findCurrentByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion | null> {
    const versions = await this.listByWorkflowDefinitionId(workflowDefinitionId);
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
      snapshotTrigger: definition.trigger,
      snapshotActions: definition.actions
    });

    this.versions.push(createdVersion);
    return createdVersion;
  }
}
