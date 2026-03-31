export type CreateWorkflowDefinitionRequest = {
  key: string;
  name: string;
  description?: string | null;
  trigger: string;
};
