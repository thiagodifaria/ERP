import { WorkflowDefinition } from "./workflow-definition.js";
import { WorkflowDefinitionVersion } from "./workflow-definition-version.js";

export interface WorkflowDefinitionVersionRepository {
  listByDefinitionKey(definitionKey: string): Promise<WorkflowDefinitionVersion[]>;
  findCurrentByDefinitionKey(definitionKey: string): Promise<WorkflowDefinitionVersion | null>;
  publish(definition: WorkflowDefinition): Promise<WorkflowDefinitionVersion>;
}
