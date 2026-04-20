import { WorkflowActionDefinition } from "../../domain/workflow-definition.js";

export type CreateWorkflowDefinitionRequest = {
  key: string;
  name: string;
  description?: string | null;
  trigger: string;
  actions?: WorkflowActionDefinition[];
};
