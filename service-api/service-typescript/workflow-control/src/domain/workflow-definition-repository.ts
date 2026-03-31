import { WorkflowDefinition, WorkflowDefinitionStatus } from "./workflow-definition.js";

export interface WorkflowDefinitionRepository {
  list(): Promise<WorkflowDefinition[]>;
  findByKey(key: string): Promise<WorkflowDefinition | null>;
  add(definition: WorkflowDefinition): Promise<WorkflowDefinition>;
  updateStatus(key: string, status: WorkflowDefinitionStatus): Promise<WorkflowDefinition>;
  nextId(): Promise<number>;
}
