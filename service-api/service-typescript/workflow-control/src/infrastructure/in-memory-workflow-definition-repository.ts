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
}
