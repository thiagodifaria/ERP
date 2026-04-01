export type CreateWorkflowRunRequest = {
  workflowDefinitionKey: string;
  subjectType: string;
  subjectPublicId: string;
  initiatedBy: string;
};
