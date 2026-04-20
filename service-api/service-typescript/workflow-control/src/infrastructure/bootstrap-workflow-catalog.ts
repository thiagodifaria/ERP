import {
  WorkflowActionCatalogItem,
  WorkflowCatalogRepository,
  WorkflowTriggerCatalogItem
} from "../domain/workflow-catalog.js";

const triggerCatalog: WorkflowTriggerCatalogItem[] = [
  {
    key: "lead.created",
    name: "Lead Created",
    description: "Starts automation when a new CRM lead is captured.",
    category: "crm",
    subjectTypes: ["crm.lead"]
  },
  {
    key: "lead.qualified",
    name: "Lead Qualified",
    description: "Starts follow-up orchestration after commercial qualification.",
    category: "crm",
    subjectTypes: ["crm.lead"]
  },
  {
    key: "lead.contacted",
    name: "Lead Contacted",
    description: "Continues automation after the first sales touchpoint.",
    category: "crm",
    subjectTypes: ["crm.lead"]
  },
  {
    key: "opportunity.stage.changed",
    name: "Opportunity Stage Changed",
    description: "Reacts to pipeline movement in the sales flow.",
    category: "sales",
    subjectTypes: ["sales.opportunity"]
  },
  {
    key: "invoice.overdue",
    name: "Invoice Overdue",
    description: "Launches follow-up actions for delayed receivables.",
    category: "billing",
    subjectTypes: ["sales.invoice"]
  },
  {
    key: "touchpoint.responded",
    name: "Touchpoint Responded",
    description: "Responds to engagement activity and qualification signals.",
    category: "engagement",
    subjectTypes: ["engagement.touchpoint"]
  },
  {
    key: "manual.dispatch",
    name: "Manual Dispatch",
    description: "Lets operators launch controlled workflow runs on demand.",
    category: "manual",
    subjectTypes: ["crm.lead", "sales.opportunity", "sales.invoice", "engagement.touchpoint"]
  }
];

const actionCatalog: WorkflowActionCatalogItem[] = [
  {
    key: "task.create",
    name: "Create Task",
    description: "Queues an operational task for a downstream team.",
    kind: "task",
    supportsCompensation: true,
    requiresRuntime: true
  },
  {
    key: "delay.wait",
    name: "Wait Delay",
    description: "Pauses execution until a relative delay elapses.",
    kind: "delay",
    supportsCompensation: false,
    requiresRuntime: true
  },
  {
    key: "integration.webhook",
    name: "Emit Webhook",
    description: "Calls an external webhook as part of an automation step.",
    kind: "integration",
    supportsCompensation: true,
    requiresRuntime: true
  },
  {
    key: "decision.condition",
    name: "Condition Branch",
    description: "Branches execution based on a resolved business condition.",
    kind: "decision",
    supportsCompensation: false,
    requiresRuntime: true
  },
  {
    key: "sales.stage.advance",
    name: "Advance Sales Stage",
    description: "Moves a commercial artifact to the next operational stage.",
    kind: "task",
    supportsCompensation: true,
    requiresRuntime: true
  }
];

export class BootstrapWorkflowCatalogRepository implements WorkflowCatalogRepository {
  public async listTriggers(): Promise<WorkflowTriggerCatalogItem[]> {
    return triggerCatalog.map((trigger) => ({
      ...trigger,
      subjectTypes: [...trigger.subjectTypes]
    }));
  }

  public async listActions(): Promise<WorkflowActionCatalogItem[]> {
    return actionCatalog.map((action) => ({ ...action }));
  }

  public async hasTrigger(triggerKey: string): Promise<boolean> {
    const normalizedTriggerKey = triggerKey.trim().toLowerCase();
    return triggerCatalog.some((trigger) => trigger.key === normalizedTriggerKey);
  }
}
