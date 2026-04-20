import { WorkflowCatalogRepository, WorkflowTriggerCatalogItem } from "../domain/workflow-catalog.js";

export class ListWorkflowTriggerCatalog {
  public constructor(
    private readonly catalogRepository: WorkflowCatalogRepository
  ) {}

  public async execute(): Promise<WorkflowTriggerCatalogItem[]> {
    return this.catalogRepository.listTriggers();
  }
}
