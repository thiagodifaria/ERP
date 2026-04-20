import { WorkflowActionDefinition, normalizeWorkflowActions } from "./workflow-definition.js";
import { WorkflowDefinitionStatus } from "./workflow-definition.js";

export type WorkflowDefinitionVersion = {
  id: number;
  workflowDefinitionId: number;
  versionNumber: number;
  snapshotName: string;
  snapshotDescription: string | null;
  snapshotStatus: WorkflowDefinitionStatus;
  snapshotTrigger: string;
  snapshotActions: WorkflowActionDefinition[];
};

export function createWorkflowDefinitionVersion(input: {
  id: number;
  workflowDefinitionId: number;
  versionNumber: number;
  snapshotName: string;
  snapshotDescription?: string | null;
  snapshotStatus: WorkflowDefinitionStatus;
  snapshotTrigger: string;
  snapshotActions?: WorkflowActionDefinition[];
}): WorkflowDefinitionVersion {
  const snapshotName = input.snapshotName.trim();
  const snapshotTrigger = input.snapshotTrigger.trim().toLowerCase();

  if (snapshotName.length === 0) {
    throw new Error("workflow_definition_version_name_required");
  }

  if (snapshotTrigger.length === 0) {
    throw new Error("workflow_definition_version_trigger_required");
  }

  if (input.versionNumber <= 0) {
    throw new Error("workflow_definition_version_number_invalid");
  }

  return {
    id: input.id,
    workflowDefinitionId: input.workflowDefinitionId,
    versionNumber: input.versionNumber,
    snapshotName,
    snapshotDescription: input.snapshotDescription?.trim() || null,
    snapshotStatus: input.snapshotStatus,
    snapshotTrigger,
    snapshotActions: normalizeWorkflowActions(input.snapshotActions)
  };
}
