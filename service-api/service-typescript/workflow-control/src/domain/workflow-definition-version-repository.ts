import { WorkflowDefinition } from "./workflow-definition.js";
import { WorkflowDefinitionVersion } from "./workflow-definition-version.js";

export interface WorkflowDefinitionVersionRepository {
  listByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion[]>;
  findByWorkflowDefinitionIdAndVersionNumber(workflowDefinitionId: number, versionNumber: number): Promise<WorkflowDefinitionVersion | null>;
  findCurrentByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion | null>;
  publish(definition: WorkflowDefinition): Promise<WorkflowDefinitionVersion>;
}
