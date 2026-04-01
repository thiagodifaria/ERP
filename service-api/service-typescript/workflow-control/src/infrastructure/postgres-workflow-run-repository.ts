import pg from "pg";
import { WorkflowRunRepository, WorkflowRunListFilters } from "../domain/workflow-run-repository.js";
import { WorkflowRun, WorkflowRunStatus, createWorkflowRun } from "../domain/workflow-run.js";

const { Pool } = pg;

type WorkflowRunRow = {
  id: number;
  public_id: string;
  workflow_definition_id: number;
  workflow_definition_version_id: number;
  status: WorkflowRunStatus;
  trigger_event: string;
  subject_type: string;
  subject_public_id: string;
  initiated_by: string;
  started_at: Date | null;
  completed_at: Date | null;
  failed_at: Date | null;
  cancelled_at: Date | null;
};

export class PostgresWorkflowRunRepository implements WorkflowRunRepository {
  private cachedTenantId: number | null = null;

  public constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  public async list(filters?: WorkflowRunListFilters): Promise<WorkflowRun[]> {
    const tenantId = await this.resolveTenantId();
    const conditions = ["tenant_id = $1"];
    const params: Array<number | string> = [tenantId];

    if (filters?.workflowDefinitionId !== undefined) {
      params.push(filters.workflowDefinitionId);
      conditions.push(`workflow_definition_id = $${params.length}`);
    }

    if (filters?.status !== undefined) {
      params.push(filters.status);
      conditions.push(`status = $${params.length}`);
    }

    if (filters?.triggerEvent !== undefined) {
      params.push(filters.triggerEvent);
      conditions.push(`trigger_event = $${params.length}`);
    }

    if (filters?.initiatedBy !== undefined) {
      params.push(filters.initiatedBy);
      conditions.push(`initiated_by = $${params.length}`);
    }

    if (filters?.subjectType !== undefined) {
      params.push(filters.subjectType);
      conditions.push(`subject_type = $${params.length}`);
    }

    if (filters?.subjectPublicId !== undefined) {
      params.push(filters.subjectPublicId);
      conditions.push(`subject_public_id = $${params.length}::uuid`);
    }

    const result = await this.pool.query<WorkflowRunRow>(
      `
        SELECT
          id,
          public_id::text,
          workflow_definition_id,
          workflow_definition_version_id,
          status,
          trigger_event,
          subject_type,
          subject_public_id::text,
          initiated_by,
          started_at,
          completed_at,
          failed_at,
          cancelled_at
        FROM workflow_control.workflow_runs
        WHERE ${conditions.join(" AND ")}
        ORDER BY id DESC
      `,
      params
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  public async findByPublicId(publicId: string): Promise<WorkflowRun | null> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowRunRow>(
      `
        SELECT
          id,
          public_id::text,
          workflow_definition_id,
          workflow_definition_version_id,
          status,
          trigger_event,
          subject_type,
          subject_public_id::text,
          initiated_by,
          started_at,
          completed_at,
          failed_at,
          cancelled_at
        FROM workflow_control.workflow_runs
        WHERE tenant_id = $1
          AND public_id = $2::uuid
        LIMIT 1
      `,
      [tenantId, publicId.trim().toLowerCase()]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  public async add(run: WorkflowRun): Promise<WorkflowRun> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowRunRow>(
      `
        INSERT INTO workflow_control.workflow_runs (
          id,
          tenant_id,
          public_id,
          workflow_definition_id,
          workflow_definition_version_id,
          status,
          trigger_event,
          subject_type,
          subject_public_id,
          initiated_by,
          started_at,
          completed_at,
          failed_at,
          cancelled_at
        )
        VALUES (
          $1,
          $2,
          $3::uuid,
          $4,
          $5,
          $6,
          $7,
          $8,
          $9::uuid,
          $10,
          $11,
          $12,
          $13,
          $14
        )
        RETURNING
          id,
          public_id::text,
          workflow_definition_id,
          workflow_definition_version_id,
          status,
          trigger_event,
          subject_type,
          subject_public_id::text,
          initiated_by,
          started_at,
          completed_at,
          failed_at,
          cancelled_at
      `,
      [
        run.id,
        tenantId,
        run.publicId,
        run.workflowDefinitionId,
        run.workflowDefinitionVersionId,
        run.status,
        run.triggerEvent,
        run.subjectType,
        run.subjectPublicId,
        run.initiatedBy,
        run.startedAt,
        run.completedAt,
        run.failedAt,
        run.cancelledAt
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  public async updateStatus(
    publicId: string,
    status: WorkflowRunStatus,
    timestamps?: Partial<Pick<WorkflowRun, "startedAt" | "completedAt" | "failedAt" | "cancelledAt">>
  ): Promise<WorkflowRun> {
    const tenantId = await this.resolveTenantId();
    const result = await this.pool.query<WorkflowRunRow>(
      `
        UPDATE workflow_control.workflow_runs
        SET status = $3,
            started_at = COALESCE($4, started_at),
            completed_at = COALESCE($5, completed_at),
            failed_at = COALESCE($6, failed_at),
            cancelled_at = COALESCE($7, cancelled_at)
        WHERE tenant_id = $1
          AND public_id = $2::uuid
        RETURNING
          id,
          public_id::text,
          workflow_definition_id,
          workflow_definition_version_id,
          status,
          trigger_event,
          subject_type,
          subject_public_id::text,
          initiated_by,
          started_at,
          completed_at,
          failed_at,
          cancelled_at
      `,
      [
        tenantId,
        publicId.trim().toLowerCase(),
        status,
        timestamps?.startedAt ?? null,
        timestamps?.completedAt ?? null,
        timestamps?.failedAt ?? null,
        timestamps?.cancelledAt ?? null
      ]
    );

    if (result.rowCount === 0) {
      throw new Error("workflow_run_not_found");
    }

    return this.mapRow(result.rows[0]);
  }

  public async nextId(): Promise<number> {
    const result = await this.pool.query<{ next_id: number }>(
      `
        SELECT COALESCE(MAX(id), 0) + 1 AS next_id
        FROM workflow_control.workflow_runs
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
      throw new Error("workflow_run_tenant_not_found");
    }

    this.cachedTenantId = Number(result.rows[0].id);
    return this.cachedTenantId;
  }

  private mapRow(row: WorkflowRunRow): WorkflowRun {
    return createWorkflowRun({
      id: Number(row.id),
      publicId: row.public_id,
      workflowDefinitionId: Number(row.workflow_definition_id),
      workflowDefinitionVersionId: Number(row.workflow_definition_version_id),
      status: row.status,
      triggerEvent: row.trigger_event,
      subjectType: row.subject_type,
      subjectPublicId: row.subject_public_id,
      initiatedBy: row.initiated_by,
      startedAt: row.started_at?.toISOString() ?? null,
      completedAt: row.completed_at?.toISOString() ?? null,
      failedAt: row.failed_at?.toISOString() ?? null,
      cancelledAt: row.cancelled_at?.toISOString() ?? null
    });
  }
}
