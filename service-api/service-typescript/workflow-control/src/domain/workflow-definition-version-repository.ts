import { WorkflowDefinition } from "./workflow-definition.js";
import { WorkflowDefinitionVersion } from "./workflow-definition-version.js";

export interface WorkflowDefinitionVersionRepository {
  listByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion[]>;
  findCurrentByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion | null>;
  publish(definition: WorkflowDefinition): Promise<WorkflowDefinitionVersion>;
}
