import pg from "pg";
import {
  buildTouchpointSummary,
  CreateTouchpointInput,
  Touchpoint,
  TouchpointFilters,
  TouchpointStatus,
  TouchpointSummary
} from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

const { Pool } = pg;

type TouchpointRow = {
  id: number;
  public_id: string;
  tenant_slug: string;
  campaign_public_id: string;
  campaign_key: string;
  lead_public_id: string;
  business_entity_type: string | null;
  business_entity_public_id: string | null;
  channel: Touchpoint["channel"];
  contact_value: string;
  source: string;
  status: TouchpointStatus;
  workflow_definition_key: string | null;
  last_workflow_run_public_id: string | null;
  created_by: string;
  notes: string;
  created_at: Date;
  updated_at: Date;
};

export class PostgresTouchpointRepository implements TouchpointRepository {
  private readonly cachedTenantIds = new Map<string, number>();

  public constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  public async list(filters: TouchpointFilters = {}): Promise<Touchpoint[]> {
    const tenantId = await this.resolveTenantId(filters.tenantSlug);
    const conditions = ["touchpoint.tenant_id = $1"];
    const params: Array<number | string> = [tenantId];

    if (filters.campaignPublicId) {
      params.push(filters.campaignPublicId);
      conditions.push(`campaign.public_id = $${params.length}::uuid`);
    }

    if (filters.status) {
      params.push(filters.status);
      conditions.push(`touchpoint.status = $${params.length}`);
    }

    if (filters.channel) {
      params.push(filters.channel);
      conditions.push(`touchpoint.channel = $${params.length}`);
    }

    if (filters.leadPublicId) {
      params.push(filters.leadPublicId);
      conditions.push(`touchpoint.lead_public_id = $${params.length}::uuid`);
    }

    if (filters.businessEntityType) {
      params.push(filters.businessEntityType);
      conditions.push(`touchpoint.business_entity_type = $${params.length}`);
    }

    if (filters.businessEntityPublicId) {
      params.push(filters.businessEntityPublicId);
      conditions.push(`touchpoint.business_entity_public_id = $${params.length}::uuid`);
    }

    const result = await this.pool.query<TouchpointRow>(
      `
        SELECT
          touchpoint.id,
          touchpoint.public_id::text,
          tenant.slug AS tenant_slug,
          campaign.public_id::text AS campaign_public_id,
          campaign.key AS campaign_key,
          touchpoint.lead_public_id::text,
          touchpoint.business_entity_type,
          touchpoint.business_entity_public_id::text,
          touchpoint.channel,
          touchpoint.contact_value,
          touchpoint.source,
          touchpoint.status,
          touchpoint.workflow_definition_key,
          touchpoint.last_workflow_run_public_id::text,
          touchpoint.created_by,
          touchpoint.notes,
          touchpoint.created_at,
          touchpoint.updated_at
        FROM engagement.touchpoints AS touchpoint
        INNER JOIN engagement.campaigns AS campaign ON campaign.id = touchpoint.campaign_id
        INNER JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id
        WHERE ${conditions.join(" AND ")}
        ORDER BY touchpoint.id
      `,
      params
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  public async getByPublicId(publicId: string): Promise<Touchpoint | null> {
    const result = await this.pool.query<TouchpointRow>(
      `
        SELECT
          touchpoint.id,
          touchpoint.public_id::text,
          tenant.slug AS tenant_slug,
          campaign.public_id::text AS campaign_public_id,
          campaign.key AS campaign_key,
          touchpoint.lead_public_id::text,
          touchpoint.business_entity_type,
          touchpoint.business_entity_public_id::text,
          touchpoint.channel,
          touchpoint.contact_value,
          touchpoint.source,
          touchpoint.status,
          touchpoint.workflow_definition_key,
          touchpoint.last_workflow_run_public_id::text,
          touchpoint.created_by,
          touchpoint.notes,
          touchpoint.created_at,
          touchpoint.updated_at
        FROM engagement.touchpoints AS touchpoint
        INNER JOIN engagement.campaigns AS campaign ON campaign.id = touchpoint.campaign_id
        INNER JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id
        WHERE touchpoint.public_id = $1::uuid
        LIMIT 1
      `,
      [publicId]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  public async create(
    input: CreateTouchpointInput & {
      campaignKey: string;
      channel: Touchpoint["channel"];
      workflowDefinitionKey: string | null;
    }
  ): Promise<Touchpoint> {
    const tenantId = await this.resolveTenantId(input.tenantSlug);
    const campaignResult = await this.pool.query<{ id: number }>(
      `
        SELECT id
        FROM engagement.campaigns
        WHERE tenant_id = $1
          AND public_id = $2::uuid
        LIMIT 1
      `,
      [tenantId, input.campaignPublicId]
    );

    if (campaignResult.rowCount === 0) {
      throw new Error("campaign_not_found");
    }

    const result = await this.pool.query<TouchpointRow>(
      `
        INSERT INTO engagement.touchpoints (
          tenant_id,
          campaign_id,
          public_id,
          lead_public_id,
          business_entity_type,
          business_entity_public_id,
          channel,
          contact_value,
          source,
          status,
          workflow_definition_key,
          created_by,
          notes
        )
        VALUES (
          $1,
          $2,
          gen_random_uuid(),
          $3::uuid,
          $4,
          $5::uuid,
          $6,
          $7,
          $8,
          'queued',
          $9,
          $10,
          $11
        )
        RETURNING
          id,
          public_id::text,
          $12::text AS tenant_slug,
          $13::text AS campaign_public_id,
          $14::text AS campaign_key,
          lead_public_id::text,
          business_entity_type,
          business_entity_public_id::text,
          channel,
          contact_value,
          source,
          status,
          workflow_definition_key,
          last_workflow_run_public_id::text,
          created_by,
          notes,
          created_at,
          updated_at
      `,
      [
        tenantId,
        campaignResult.rows[0].id,
        input.leadPublicId,
        input.businessEntityType,
        input.businessEntityPublicId,
        input.channel,
        input.contactValue,
        input.source,
        input.workflowDefinitionKey,
        input.createdBy,
        input.notes,
        input.tenantSlug,
        input.campaignPublicId,
        input.campaignKey
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  public async updateStatus(
    publicId: string,
    status: TouchpointStatus,
    lastWorkflowRunPublicId?: string | null
  ): Promise<Touchpoint | null> {
    const result = await this.pool.query<TouchpointRow>(
      `
        UPDATE engagement.touchpoints AS touchpoint
        SET status = $2,
            last_workflow_run_public_id = COALESCE($3::uuid, touchpoint.last_workflow_run_public_id)
        FROM engagement.campaigns AS campaign, identity.tenants AS tenant
        WHERE campaign.id = touchpoint.campaign_id
          AND tenant.id = touchpoint.tenant_id
          AND touchpoint.public_id = $1::uuid
        RETURNING
          touchpoint.id,
          touchpoint.public_id::text,
          tenant.slug AS tenant_slug,
          campaign.public_id::text AS campaign_public_id,
          campaign.key AS campaign_key,
          touchpoint.lead_public_id::text,
          touchpoint.business_entity_type,
          touchpoint.business_entity_public_id::text,
          touchpoint.channel,
          touchpoint.contact_value,
          touchpoint.source,
          touchpoint.status,
          touchpoint.workflow_definition_key,
          touchpoint.last_workflow_run_public_id::text,
          touchpoint.created_by,
          touchpoint.notes,
          touchpoint.created_at,
          touchpoint.updated_at
      `,
      [publicId, status, lastWorkflowRunPublicId ?? null]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  public async getSummary(filters: TouchpointFilters = {}): Promise<TouchpointSummary> {
    return buildTouchpointSummary(await this.list(filters), filters.tenantSlug ?? this.bootstrapTenantSlug);
  }

  private async resolveTenantId(tenantSlug?: string): Promise<number> {
    const slug = tenantSlug ?? this.bootstrapTenantSlug;

    if (this.cachedTenantIds.has(slug)) {
      return this.cachedTenantIds.get(slug)!;
    }

    const result = await this.pool.query<{ id: number }>(
      `
        SELECT id
        FROM identity.tenants
        WHERE slug = $1
        LIMIT 1
      `,
      [slug]
    );

    if (result.rowCount === 0) {
      throw new Error("touchpoint_tenant_not_found");
    }

    const tenantId = Number(result.rows[0].id);
    this.cachedTenantIds.set(slug, tenantId);
    return tenantId;
  }

  private mapRow(row: TouchpointRow): Touchpoint {
    return {
      id: Number(row.id),
      publicId: row.public_id,
      tenantSlug: row.tenant_slug,
      campaignPublicId: row.campaign_public_id,
      campaignKey: row.campaign_key,
      leadPublicId: row.lead_public_id,
      businessEntityType: row.business_entity_type,
      businessEntityPublicId: row.business_entity_public_id,
      channel: row.channel,
      contactValue: row.contact_value,
      source: row.source,
      status: row.status,
      workflowDefinitionKey: row.workflow_definition_key,
      lastWorkflowRunPublicId: row.last_workflow_run_public_id,
      createdBy: row.created_by,
      notes: row.notes,
      createdAt: row.created_at.toISOString(),
      updatedAt: row.updated_at.toISOString()
    };
  }
}
