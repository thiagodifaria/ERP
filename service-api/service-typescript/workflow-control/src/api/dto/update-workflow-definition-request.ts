import { WorkflowActionDefinition } from "../../domain/workflow-definition.js";

export type UpdateWorkflowDefinitionRequest = {
  name?: string;
  description?: string | null;
  trigger?: string;
  actions?: WorkflowActionDefinition[];
};
