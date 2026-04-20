import { WorkflowDefinition, createWorkflowDefinition } from "../domain/workflow-definition.js";
import { WorkflowCatalogRepository } from "../domain/workflow-catalog.js";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";

export class CreateWorkflowDefinition {
  public constructor(
    private readonly repository: WorkflowDefinitionRepository,
    private readonly catalogRepository: WorkflowCatalogRepository
  ) {}

  public async execute(input: {
    key: string;
    name: string;
    description?: string | null;
    trigger: string;
    actions?: WorkflowDefinition["actions"];
  }): Promise<WorkflowDefinition> {
    if (await this.repository.findByKey(input.key) !== null) {
      throw new Error("workflow_definition_key_conflict");
    }

    if (!(await this.catalogRepository.hasTrigger(input.trigger))) {
      throw new Error("workflow_definition_trigger_unknown");
    }

    if (input.actions !== undefined) {
      for (const action of input.actions) {
        if (!(await this.catalogRepository.hasAction(action.actionKey))) {
          throw new Error("workflow_definition_action_key_unknown");
        }

        if (action.compensationActionKey && !(await this.catalogRepository.hasAction(action.compensationActionKey))) {
          throw new Error("workflow_definition_action_key_unknown");
        }
      }
    }

    const definition = createWorkflowDefinition({
      id: await this.repository.nextId(),
      key: input.key,
      name: input.name,
      description: input.description,
      status: "draft",
      trigger: input.trigger,
      actions: input.actions
    });

    return this.repository.add(definition);
  }
}
