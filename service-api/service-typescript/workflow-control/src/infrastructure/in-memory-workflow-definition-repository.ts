import { WorkflowDefinition, createWorkflowDefinition } from "../domain/workflow-definition.js";

export class InMemoryWorkflowDefinitionRepository {
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

  public list(): WorkflowDefinition[] {
    return [...this.definitions];
  }

  public findByKey(key: string): WorkflowDefinition | null {
    const normalizedKey = key.trim().toLowerCase();

    return this.definitions.find((definition) => definition.key === normalizedKey) ?? null;
  }

  public add(definition: WorkflowDefinition): WorkflowDefinition {
    this.definitions.push(definition);
    return definition;
  }

  public nextId(): number {
    return this.definitions.reduce((max, definition) => Math.max(max, definition.id), 0) + 1;
  }
}
