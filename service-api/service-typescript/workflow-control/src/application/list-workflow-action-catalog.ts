import { WorkflowActionCatalogItem, WorkflowCatalogRepository } from "../domain/workflow-catalog.js";

export class ListWorkflowActionCatalog {
  public constructor(
    private readonly catalogRepository: WorkflowCatalogRepository
  ) {}

  public async execute(): Promise<WorkflowActionCatalogItem[]> {
    return this.catalogRepository.listActions();
  }
}
