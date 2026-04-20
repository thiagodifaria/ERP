export type WorkflowDefinitionStatus = "draft" | "active" | "archived";

export type WorkflowActionDefinition = {
  stepId: string;
  actionKey: string;
  label: string;
  delaySeconds: number | null;
  compensationActionKey: string | null;
};

export type WorkflowDefinition = {
  id: number;
  key: string;
  name: string;
  description: string | null;
  status: WorkflowDefinitionStatus;
  trigger: string;
  actions: WorkflowActionDefinition[];
};

export type WorkflowDefinitionUpdateInput = {
  name?: string;
  description?: string | null;
  trigger?: string;
  actions?: WorkflowActionDefinition[];
};

export function ensureWorkflowDefinitionStatus(value: string): WorkflowDefinitionStatus {
  const normalizedValue = value.trim().toLowerCase();

  if (normalizedValue !== "draft" && normalizedValue !== "active" && normalizedValue !== "archived") {
    throw new Error("workflow_definition_status_invalid");
  }

  return normalizedValue;
}

export function createWorkflowDefinition(input: {
  id: number;
  key: string;
  name: string;
  description?: string | null;
  status?: WorkflowDefinitionStatus;
  trigger: string;
  actions?: WorkflowActionDefinition[];
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
    status: ensureWorkflowDefinitionStatus(input.status ?? "draft"),
    trigger,
    actions: normalizeWorkflowActions(input.actions)
  };
}

export function applyWorkflowDefinitionUpdate(
  currentDefinition: WorkflowDefinition,
  input: WorkflowDefinitionUpdateInput
): WorkflowDefinition {
  const hasName = Object.prototype.hasOwnProperty.call(input, "name");
  const hasDescription = Object.prototype.hasOwnProperty.call(input, "description");
  const hasTrigger = Object.prototype.hasOwnProperty.call(input, "trigger");
  const hasActions = Object.prototype.hasOwnProperty.call(input, "actions");

  if (!hasName && !hasDescription && !hasTrigger && !hasActions) {
    throw new Error("workflow_definition_update_required");
  }

  return createWorkflowDefinition({
    id: currentDefinition.id,
    key: currentDefinition.key,
    name: hasName ? input.name ?? "" : currentDefinition.name,
    description: hasDescription ? input.description : currentDefinition.description,
    status: currentDefinition.status,
    trigger: hasTrigger ? input.trigger ?? "" : currentDefinition.trigger,
    actions: Object.prototype.hasOwnProperty.call(input, "actions") ? input.actions : currentDefinition.actions
  });
}

export function normalizeWorkflowActions(actions?: WorkflowActionDefinition[]): WorkflowActionDefinition[] {
  if (actions === undefined) {
    return [];
  }

  return actions.map((action, index) => {
    const stepId = action.stepId.trim().toLowerCase();
    const actionKey = action.actionKey.trim().toLowerCase();
    const label = action.label.trim();
    const compensationActionKey = action.compensationActionKey?.trim().toLowerCase() || null;
    const delaySeconds = action.actionKey.trim().toLowerCase() === "delay.wait"
      ? Number(action.delaySeconds ?? 0)
      : null;

    if (stepId.length === 0) {
      throw new Error("workflow_definition_action_step_id_required");
    }

    if (actionKey.length === 0) {
      throw new Error("workflow_definition_action_key_required");
    }

    if (label.length === 0) {
      throw new Error("workflow_definition_action_label_required");
    }

    if (actionKey === "delay.wait") {
      if (delaySeconds === null || !Number.isInteger(delaySeconds) || delaySeconds <= 0) {
        throw new Error("workflow_definition_action_delay_invalid");
      }
    }

    return {
      stepId,
      actionKey,
      label,
      delaySeconds,
      compensationActionKey
    };
  }).filter((action, index, array) => {
    if (array.findIndex((candidate) => candidate.stepId === action.stepId) !== index) {
      throw new Error("workflow_definition_action_step_id_conflict");
    }

    return true;
  });
}
