import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinition, WorkflowDefinitionStatus, createWorkflowDefinition } from "../domain/workflow-definition.js";

export class InMemoryWorkflowDefinitionRepository implements WorkflowDefinitionRepository {
  private readonly definitions: WorkflowDefinition[];

  public constructor() {
    this.definitions = [
      createWorkflowDefinition({
        id: 1,
        key: "lead-follow-up",
        name: "Lead Follow-Up",
        description: "Orquestra o acompanhamento inicial de novos leads do CRM.",
        status: "active",
        trigger: "lead.created"
      })
    ];
  }

  public async list(): Promise<WorkflowDefinition[]> {
    return [...this.definitions];
  }

  public async findByKey(key: string): Promise<WorkflowDefinition | null> {
    const normalizedKey = key.trim().toLowerCase();

    return this.definitions.find((definition) => definition.key === normalizedKey) ?? null;
  }

  public async add(definition: WorkflowDefinition): Promise<WorkflowDefinition> {
    this.definitions.push(definition);
    return definition;
  }

  public async updateDefinition(key: string, definition: WorkflowDefinition): Promise<WorkflowDefinition> {
    const normalizedKey = key.trim().toLowerCase();
    const currentDefinition = this.definitions.find((candidate) => candidate.key === normalizedKey);

    if (!currentDefinition) {
      throw new Error("workflow_definition_not_found");
    }

    currentDefinition.name = definition.name;
    currentDefinition.description = definition.description;
    currentDefinition.trigger = definition.trigger;

    return currentDefinition;
  }

  public async updateStatus(key: string, status: WorkflowDefinitionStatus): Promise<WorkflowDefinition> {
    const normalizedKey = key.trim().toLowerCase();
    const definition = this.definitions.find((candidate) => candidate.key === normalizedKey);

    if (!definition) {
      throw new Error("workflow_definition_not_found");
    }

    definition.status = status;
    return definition;
  }

  public async nextId(): Promise<number> {
    return this.definitions.reduce((max, definition) => Math.max(max, definition.id), 0) + 1;
  }
}
