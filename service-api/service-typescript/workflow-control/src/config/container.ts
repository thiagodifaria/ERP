import { CreateWorkflowDefinition } from "../application/create-workflow-definition.js";
import { GetWorkflowDefinitionByKey } from "../application/get-workflow-definition-by-key.js";
import { ListWorkflowDefinitions } from "../application/list-workflow-definitions.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";

const repository = new InMemoryWorkflowDefinitionRepository();

export const services = {
  createWorkflowDefinition: new CreateWorkflowDefinition(repository),
  getWorkflowDefinitionByKey: new GetWorkflowDefinitionByKey(repository),
  listWorkflowDefinitions: new ListWorkflowDefinitions(repository)
};
