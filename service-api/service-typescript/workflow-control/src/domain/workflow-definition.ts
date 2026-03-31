export type WorkflowDefinitionStatus = "draft" | "active" | "archived";

export type WorkflowDefinition = {
  id: number;
  key: string;
  name: string;
  description: string | null;
  status: WorkflowDefinitionStatus;
  trigger: string;
};

export function createWorkflowDefinition(input: {
  id: number;
  key: string;
  name: string;
  description?: string | null;
  status?: WorkflowDefinitionStatus;
  trigger: string;
}): WorkflowDefinition {
  const key = input.key.trim().toLowerCase();
  const name = input.name.trim();
  const trigger = input.trigger.trim().toLowerCase();

  if (key.length === 0) {
    throw new Error("workflow_definition_key_required");
  }

  if (name.length === 0) {
    throw new Error("workflow_definition_name_required");
  }

  if (trigger.length === 0) {
    throw new Error("workflow_definition_trigger_required");
  }

  return {
    id: input.id,
    key,
    name,
    description: input.description?.trim() || null,
    status: input.status ?? "draft",
    trigger
  };
}
