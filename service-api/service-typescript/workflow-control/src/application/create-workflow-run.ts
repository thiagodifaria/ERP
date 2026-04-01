import { randomUUID } from "node:crypto";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { WorkflowRunRepository } from "../domain/workflow-run-repository.js";
import { WorkflowRun, createWorkflowRun } from "../domain/workflow-run.js";

export class CreateWorkflowRun {
  public constructor(
    private readonly definitionRepository: WorkflowDefinitionRepository,
    private readonly versionRepository: WorkflowDefinitionVersionRepository,
    private readonly runRepository: WorkflowRunRepository
  ) {}

  public async execute(input: {
    workflowDefinitionKey: string;
    subjectType: string;
    subjectPublicId: string;
    initiatedBy: string;
  }): Promise<WorkflowRun> {
    const definitionKey = input.workflowDefinitionKey.trim().toLowerCase();

    if (definitionKey.length === 0) {
      throw new Error("workflow_run_definition_key_required");
    }

    const definition = await this.definitionRepository.findByKey(definitionKey);

    if (definition === null) {
      throw new Error("workflow_definition_not_found");
    }

    const currentVersion = await this.versionRepository.findCurrentByWorkflowDefinitionId(definition.id);

    if (currentVersion === null) {
      throw new Error("workflow_definition_version_not_found");
    }

    const workflowRun = createWorkflowRun({
      id: await this.runRepository.nextId(),
      publicId: randomUUID(),
      workflowDefinitionId: definition.id,
      workflowDefinitionVersionId: currentVersion.id,
      status: "pending",
      triggerEvent: currentVersion.snapshotTrigger,
      subjectType: input.subjectType,
      subjectPublicId: input.subjectPublicId,
      initiatedBy: input.initiatedBy
    });

    return this.runRepository.add(workflowRun);
  }
}
