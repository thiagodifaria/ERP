import { WorkflowDefinition, createWorkflowDefinition } from "../domain/workflow-definition.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class CreateWorkflowDefinition {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository
  ) {}

  public execute(input: {
    key: string;
    name: string;
    description?: string | null;
    trigger: string;
  }): WorkflowDefinition {
    if (this.repository.findByKey(input.key) !== null) {
      throw new Error("workflow_definition_key_conflict");
    }

    const definition = createWorkflowDefinition({
      id: this.repository.nextId(),
      key: input.key,
      name: input.name,
      description: input.description,
      status: "draft",
      trigger: input.trigger
    });

    return this.repository.add(definition);
  }
}
