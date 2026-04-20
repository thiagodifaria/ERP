export type WorkflowTriggerCatalogItem = {
  key: string;
  name: string;
  description: string;
  category: "crm" | "sales" | "billing" | "engagement" | "manual";
  subjectTypes: string[];
};

export type WorkflowActionCatalogItem = {
  key: string;
  name: string;
  description: string;
  kind: "task" | "delay" | "integration" | "decision";
  supportsCompensation: boolean;
  requiresRuntime: boolean;
};

export interface WorkflowCatalogRepository {
  listTriggers(): Promise<WorkflowTriggerCatalogItem[]>;
  listActions(): Promise<WorkflowActionCatalogItem[]>;
  hasTrigger(triggerKey: string): Promise<boolean>;
}
