import { WorkflowDefinition, WorkflowDefinitionStatus } from "./workflow-definition.js";

export interface WorkflowDefinitionRepository {
  list(): WorkflowDefinition[];
  findByKey(key: string): WorkflowDefinition | null;
  add(definition: WorkflowDefinition): WorkflowDefinition;
  updateStatus(key: string, status: WorkflowDefinitionStatus): WorkflowDefinition;
  nextId(): number;
}
