import pg from "pg";
import { WorkflowDefinitionRepository } from "../domain/workflow-definition-repository.js";
import { WorkflowDefinition, WorkflowDefinitionStatus, createWorkflowDefinition } from "../domain/workflow-definition.js";

const { Pool } = pg;

type WorkflowDefinitionRow = {
  id: number;
  key: string;
  name: string;
  description: string | null;
  status: WorkflowDefinitionStatus;
  trigger: string;
};

export class PostgresWorkflowDefinitionRepository implements WorkflowDefinitionRepository {
  private cachedTenantId: number | null = null;

  public constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  public async list(): Promise<WorkflowDefinition[]> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowDefinitionRow>(
      `
        SELECT id, key, name, description, status, trigger
        FROM workflow_control.workflow_definitions
        WHERE tenant_id = $1
        ORDER BY id
      `,
      [tenantId]
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  public async findByKey(key: string): Promise<WorkflowDefinition | null> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowDefinitionRow>(
      `
        SELECT id, key, name, description, status, trigger
        FROM workflow_control.workflow_definitions
        WHERE tenant_id = $1
          AND key = $2
        LIMIT 1
      `,
      [tenantId, key.trim().toLowerCase()]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  public async add(definition: WorkflowDefinition): Promise<WorkflowDefinition> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowDefinitionRow>(
      `
        INSERT INTO workflow_control.workflow_definitions (
          id,
          tenant_id,
          public_id,
          key,
          name,
          description,
          status,
          trigger
        )
        VALUES (
          $1,
          $2,
          gen_random_uuid(),
          $3,
          $4,
          $5,
          $6,
          $7
        )
        RETURNING id, key, name, description, status, trigger
      `,
      [
        definition.id,
        tenantId,
        definition.key,
        definition.name,
        definition.description,
        definition.status,
        definition.trigger
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  public async updateStatus(key: string, status: WorkflowDefinitionStatus): Promise<WorkflowDefinition> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowDefinitionRow>(
      `
        UPDATE workflow_control.workflow_definitions
        SET status = $3
        WHERE tenant_id = $1
          AND key = $2
        RETURNING id, key, name, description, status, trigger
      `,
      [tenantId, key.trim().toLowerCase(), status]
    );

    if (result.rowCount === 0) {
      throw new Error("workflow_definition_not_found");
    }

    return this.mapRow(result.rows[0]);
  }

  public async nextId(): Promise<number> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<{ next_id: number }>(
      `
        SELECT COALESCE(MAX(id), 0) + 1 AS next_id
        FROM workflow_control.workflow_definitions
        WHERE tenant_id = $1
      `,
      [tenantId]
    );

    return Number(result.rows[0]?.next_id ?? 1);
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

  private mapRow(row: WorkflowDefinitionRow): WorkflowDefinition {
    return createWorkflowDefinition({
      id: Number(row.id),
      key: row.key,
      name: row.name,
      description: row.description,
      status: row.status,
      trigger: row.trigger
    });
  }
}
