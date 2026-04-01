import pg from "pg";
import { WorkflowRunEventRepository } from "../domain/workflow-run-event-repository.js";
import { WorkflowRunEvent, WorkflowRunEventCategory, createWorkflowRunEvent } from "../domain/workflow-run-event.js";

const { Pool } = pg;

type WorkflowRunEventRow = {
  id: number;
  public_id: string;
  workflow_run_public_id: string;
  category: WorkflowRunEventCategory;
  body: string;
  created_by: string;
  created_at: Date;
};

export class PostgresWorkflowRunEventRepository implements WorkflowRunEventRepository {
  private cachedTenantId: number | null = null;

  public constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  public async listByWorkflowRunPublicId(workflowRunPublicId: string): Promise<WorkflowRunEvent[]> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowRunEventRow>(
      `
        SELECT
          event.id,
          event.public_id::text,
          workflow_run.public_id::text AS workflow_run_public_id,
          event.category,
          event.body,
          event.created_by,
          event.created_at
        FROM workflow_control.workflow_run_events AS event
        INNER JOIN workflow_control.workflow_runs AS workflow_run
          ON workflow_run.id = event.workflow_run_id
        WHERE event.tenant_id = $1
          AND workflow_run.public_id = $2::uuid
        ORDER BY event.id
      `,
      [tenantId, workflowRunPublicId.trim().toLowerCase()]
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  public async add(event: WorkflowRunEvent): Promise<WorkflowRunEvent> {
    const tenantId = await this.resolveTenantId();
    const workflowRunIdResult = await this.pool.query<{ id: number }>(
      `
        SELECT id
        FROM workflow_control.workflow_runs
        WHERE tenant_id = $1
          AND public_id = $2::uuid
        LIMIT 1
      `,
      [tenantId, event.workflowRunPublicId]
    );

    if (workflowRunIdResult.rowCount === 0) {
      throw new Error("workflow_run_not_found");
    }

    const result = await this.pool.query<WorkflowRunEventRow>(
      `
        INSERT INTO workflow_control.workflow_run_events (
          id,
          tenant_id,
          public_id,
          workflow_run_id,
          category,
          body,
          created_by,
          created_at
        )
        VALUES (
          $1,
          $2,
          $3::uuid,
          $4,
          $5,
          $6,
          $7,
          $8
        )
        RETURNING
          id,
          public_id::text,
          $9::text AS workflow_run_public_id,
          category,
          body,
          created_by,
          created_at
      `,
      [
        event.id,
        tenantId,
        event.publicId,
        Number(workflowRunIdResult.rows[0].id),
        event.category,
        event.body,
        event.createdBy,
        event.createdAt,
        event.workflowRunPublicId
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  public async nextId(): Promise<number> {
    const result = await this.pool.query<{ next_id: number }>(
      `
        SELECT COALESCE(MAX(id), 0) + 1 AS next_id
        FROM workflow_control.workflow_run_events
      `
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
      throw new Error("workflow_run_event_tenant_not_found");
    }

    this.cachedTenantId = Number(result.rows[0].id);
    return this.cachedTenantId;
  }

  private mapRow(row: WorkflowRunEventRow): WorkflowRunEvent {
    return createWorkflowRunEvent({
      id: Number(row.id),
      publicId: row.public_id,
      workflowRunPublicId: row.workflow_run_public_id,
      category: row.category,
      body: row.body,
      createdBy: row.created_by,
      createdAt: row.created_at.toISOString()
    });
  }
}
