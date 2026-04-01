import pg from "pg";
import { WorkflowDefinition } from "../domain/workflow-definition.js";
import { WorkflowDefinitionVersionRepository } from "../domain/workflow-definition-version-repository.js";
import { WorkflowDefinitionVersion, createWorkflowDefinitionVersion } from "../domain/workflow-definition-version.js";

const { Pool } = pg;

type WorkflowDefinitionVersionRow = {
  id: number;
  workflow_definition_id: number;
  version_number: number;
  snapshot_name: string;
  snapshot_description: string | null;
  snapshot_status: "draft" | "active" | "archived";
  snapshot_trigger: string;
};

export class PostgresWorkflowDefinitionVersionRepository implements WorkflowDefinitionVersionRepository {
  private cachedTenantId: number | null = null;

  public constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  public async listByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion[]> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowDefinitionVersionRow>(
      `
        SELECT
          id,
          workflow_definition_id,
          version_number,
          snapshot_name,
          snapshot_description,
          snapshot_status,
          snapshot_trigger
        FROM workflow_control.workflow_definition_versions
        WHERE tenant_id = $1
          AND workflow_definition_id = $2
        ORDER BY version_number DESC
      `,
      [tenantId, workflowDefinitionId]
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  public async findByWorkflowDefinitionIdAndVersionNumber(
    workflowDefinitionId: number,
    versionNumber: number
  ): Promise<WorkflowDefinitionVersion | null> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowDefinitionVersionRow>(
      `
        SELECT
          id,
          workflow_definition_id,
          version_number,
          snapshot_name,
          snapshot_description,
          snapshot_status,
          snapshot_trigger
        FROM workflow_control.workflow_definition_versions
        WHERE tenant_id = $1
          AND workflow_definition_id = $2
          AND version_number = $3
        LIMIT 1
      `,
      [tenantId, workflowDefinitionId, versionNumber]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  public async findCurrentByWorkflowDefinitionId(workflowDefinitionId: number): Promise<WorkflowDefinitionVersion | null> {
    const versions = await this.listByWorkflowDefinitionId(workflowDefinitionId);
    return versions[0] ?? null;
  }

  public async publish(definition: WorkflowDefinition): Promise<WorkflowDefinitionVersion> {
    const tenantId = await this.resolveTenantId();
    const currentVersion = await this.findCurrentByWorkflowDefinitionId(definition.id);
    const nextVersionNumber = (currentVersion?.versionNumber ?? 0) + 1;

    const result = await this.pool.query<WorkflowDefinitionVersionRow>(
      `
        INSERT INTO workflow_control.workflow_definition_versions (
          tenant_id,
          workflow_definition_id,
          version_number,
          snapshot_name,
          snapshot_description,
          snapshot_status,
          snapshot_trigger
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING
          id,
          workflow_definition_id,
          version_number,
          snapshot_name,
          snapshot_description,
          snapshot_status,
          snapshot_trigger
      `,
      [
        tenantId,
        definition.id,
        nextVersionNumber,
        definition.name,
        definition.description,
        definition.status,
        definition.trigger
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  private async resolveTenantId(): Promise<number> {
    if (this.cachedTenantId !== null) {
      return this.cachedTenantId;
    }

    const result = await this.pool.query<{ id: number }>(
      `
        SELECT id
        FROM identity.tenants
        WHERE slug = $1
        LIMIT 1
      `,
      [this.bootstrapTenantSlug]
    );

    if (result.rowCount === 0) {
      throw new Error("workflow_definition_tenant_not_found");
    }

    this.cachedTenantId = Number(result.rows[0].id);
    return this.cachedTenantId;
  }

  private mapRow(row: WorkflowDefinitionVersionRow): WorkflowDefinitionVersion {
    return createWorkflowDefinitionVersion({
      id: Number(row.id),
      workflowDefinitionId: Number(row.workflow_definition_id),
      versionNumber: Number(row.version_number),
      snapshotName: row.snapshot_name,
      snapshotDescription: row.snapshot_description,
      snapshotStatus: row.snapshot_status,
      snapshotTrigger: row.snapshot_trigger
    });
  }
}
