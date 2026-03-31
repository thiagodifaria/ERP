import { ListWorkflowDefinitions } from "../application/list-workflow-definitions.js";
import { InMemoryWorkflowDefinitionRepository } from "../infrastructure/in-memory-workflow-definition-repository.js";

const repository = new InMemoryWorkflowDefinitionRepository();

export const services = {
  listWorkflowDefinitions: new ListWorkflowDefinitions(repository)
};
