import pg from "pg";
import { CampaignRepository } from "../domain/campaign-repository.js";
import { Campaign, CampaignFilters, CampaignStatus, CreateCampaignInput } from "../domain/campaign.js";

const { Pool } = pg;

type CampaignRow = {
  id: number;
  public_id: string;
  tenant_slug: string;
  key: string;
  name: string;
  description: string;
  channel: Campaign["channel"];
  status: CampaignStatus;
  touchpoint_goal: string;
  workflow_definition_key: string | null;
  budget_cents: number;
  created_at: Date;
  updated_at: Date;
};

export class PostgresCampaignRepository implements CampaignRepository {
  private readonly cachedTenantIds = new Map<string, number>();

  public constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  public async list(filters: CampaignFilters = {}): Promise<Campaign[]> {
    const tenantId = await this.resolveTenantId(filters.tenantSlug);
    const conditions = ["campaign.tenant_id = $1"];
    const params: Array<number | string> = [tenantId];

    if (filters.status) {
      params.push(filters.status);
      conditions.push(`campaign.status = $${params.length}`);
    }

    if (filters.channel) {
      params.push(filters.channel);
      conditions.push(`campaign.channel = $${params.length}`);
    }

    if (filters.q) {
      params.push(`%${filters.q.trim().toLowerCase()}%`);
      conditions.push(`LOWER(campaign.key || ' ' || campaign.name || ' ' || campaign.description) LIKE $${params.length}`);
    }

    const result = await this.pool.query<CampaignRow>(
      `
        SELECT
          campaign.id,
          campaign.public_id::text,
          tenant.slug AS tenant_slug,
          campaign.key,
          campaign.name,
          campaign.description,
          campaign.channel,
          campaign.status,
          campaign.touchpoint_goal,
          campaign.workflow_definition_key,
          campaign.budget_cents,
          campaign.created_at,
          campaign.updated_at
        FROM engagement.campaigns AS campaign
        INNER JOIN identity.tenants AS tenant ON tenant.id = campaign.tenant_id
        WHERE ${conditions.join(" AND ")}
        ORDER BY campaign.id
      `,
      params
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  public async getByPublicId(publicId: string): Promise<Campaign | null> {
    const result = await this.pool.query<CampaignRow>(
      `
        SELECT
          campaign.id,
          campaign.public_id::text,
          tenant.slug AS tenant_slug,
          campaign.key,
          campaign.name,
          campaign.description,
          campaign.channel,
          campaign.status,
          campaign.touchpoint_goal,
          campaign.workflow_definition_key,
          campaign.budget_cents,
          campaign.created_at,
          campaign.updated_at
        FROM engagement.campaigns AS campaign
        INNER JOIN identity.tenants AS tenant ON tenant.id = campaign.tenant_id
        WHERE campaign.public_id = $1::uuid
        LIMIT 1
      `,
      [publicId]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  public async create(input: CreateCampaignInput): Promise<Campaign> {
    const tenantId = await this.resolveTenantId(input.tenantSlug);
    const result = await this.pool
      .query<CampaignRow>(
        `
          INSERT INTO engagement.campaigns (
            tenant_id,
            public_id,
            key,
            name,
            description,
            channel,
            status,
            touchpoint_goal,
            workflow_definition_key,
            budget_cents
          )
          VALUES (
            $1,
            gen_random_uuid(),
            $2,
            $3,
            $4,
            $5,
            $6,
            $7,
            $8,
            $9
          )
          RETURNING
            id,
            public_id::text,
            $10::text AS tenant_slug,
            key,
            name,
            description,
            channel,
            status,
            touchpoint_goal,
            workflow_definition_key,
            budget_cents,
            created_at,
            updated_at
        `,
        [
          tenantId,
          input.key,
          input.name,
          input.description,
          input.channel,
          input.status ?? "draft",
          input.touchpointGoal,
          input.workflowDefinitionKey,
          input.budgetCents,
          input.tenantSlug
        ]
      )
      .catch((error: { code?: string }) => {
        if (error.code === "23505") {
          throw new Error("campaign_key_conflict");
        }

        throw error;
      });

    return this.mapRow(result.rows[0]);
  }

  public async updateStatus(publicId: string, status: CampaignStatus): Promise<Campaign | null> {
    const result = await this.pool.query<CampaignRow>(
      `
        UPDATE engagement.campaigns AS campaign
        SET status = $2
        FROM identity.tenants AS tenant
        WHERE tenant.id = campaign.tenant_id
          AND campaign.public_id = $1::uuid
        RETURNING
          campaign.id,
          campaign.public_id::text,
          tenant.slug AS tenant_slug,
          campaign.key,
          campaign.name,
          campaign.description,
          campaign.channel,
          campaign.status,
          campaign.touchpoint_goal,
          campaign.workflow_definition_key,
          campaign.budget_cents,
          campaign.created_at,
          campaign.updated_at
      `,
      [publicId, status]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
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
      throw new Error("campaign_tenant_not_found");
    }

    const tenantId = Number(result.rows[0].id);
    this.cachedTenantIds.set(slug, tenantId);
    return tenantId;
  }

  private mapRow(row: CampaignRow): Campaign {
    return {
      id: Number(row.id),
      publicId: row.public_id,
      tenantSlug: row.tenant_slug,
      key: row.key,
      name: row.name,
      description: row.description,
      channel: row.channel,
      status: row.status,
      touchpointGoal: row.touchpoint_goal,
      workflowDefinitionKey: row.workflow_definition_key,
      budgetCents: Number(row.budget_cents),
      createdAt: row.created_at.toISOString(),
      updatedAt: row.updated_at.toISOString()
    };
  }
}
