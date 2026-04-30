import pg from "pg";
import { TemplateRepository } from "../domain/template-repository.js";
import { CreateTemplateInput, EngagementTemplate, TemplateFilters, TemplateStatus } from "../domain/template.js";

const { Pool } = pg;

type TemplateRow = {
  id: number;
  public_id: string;
  tenant_slug: string;
  key: string;
  name: string;
  channel: EngagementTemplate["channel"];
  status: TemplateStatus;
  provider: EngagementTemplate["provider"];
  subject: string | null;
  body: string;
  created_at: Date;
  updated_at: Date;
};

export class PostgresTemplateRepository implements TemplateRepository {
  private readonly cachedTenantIds = new Map<string, number>();

  constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  async list(filters: TemplateFilters = {}): Promise<EngagementTemplate[]> {
    const tenantId = await this.resolveTenantId(filters.tenantSlug);
    const conditions = ["template.tenant_id = $1"];
    const params: Array<number | string> = [tenantId];

    if (filters.channel) {
      params.push(filters.channel);
      conditions.push(`template.channel = $${params.length}`);
    }
    if (filters.status) {
      params.push(filters.status);
      conditions.push(`template.status = $${params.length}`);
    }
    if (filters.provider) {
      params.push(filters.provider);
      conditions.push(`template.provider = $${params.length}`);
    }
    if (filters.q) {
      params.push(`%${filters.q.trim().toLowerCase()}%`);
      conditions.push(`LOWER(template.key || ' ' || template.name || ' ' || COALESCE(template.subject, '') || ' ' || template.body) LIKE $${params.length}`);
    }

    const result = await this.pool.query<TemplateRow>(
      `
        SELECT
          template.id,
          template.public_id::text,
          tenant.slug AS tenant_slug,
          template.key,
          template.name,
          template.channel,
          template.status,
          template.provider,
          template.subject,
          template.body,
          template.created_at,
          template.updated_at
        FROM engagement.templates AS template
        INNER JOIN identity.tenants AS tenant ON tenant.id = template.tenant_id
        WHERE ${conditions.join(" AND ")}
        ORDER BY template.id
      `,
      params
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  async getByPublicId(publicId: string): Promise<EngagementTemplate | null> {
    const result = await this.pool.query<TemplateRow>(
      `
        SELECT
          template.id,
          template.public_id::text,
          tenant.slug AS tenant_slug,
          template.key,
          template.name,
          template.channel,
          template.status,
          template.provider,
          template.subject,
          template.body,
          template.created_at,
          template.updated_at
        FROM engagement.templates AS template
        INNER JOIN identity.tenants AS tenant ON tenant.id = template.tenant_id
        WHERE template.public_id = $1::uuid
        LIMIT 1
      `,
      [publicId]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  async create(input: CreateTemplateInput): Promise<EngagementTemplate> {
    const tenantId = await this.resolveTenantId(input.tenantSlug);
    const result = await this.pool
      .query<TemplateRow>(
        `
          INSERT INTO engagement.templates (
            tenant_id,
            public_id,
            key,
            name,
            channel,
            status,
            provider,
            subject,
            body
          )
          VALUES ($1, gen_random_uuid(), $2, $3, $4, $5, $6, $7, $8)
          RETURNING
            id,
            public_id::text,
            $9::text AS tenant_slug,
            key,
            name,
            channel,
            status,
            provider,
            subject,
            body,
            created_at,
            updated_at
        `,
        [
          tenantId,
          input.key,
          input.name,
          input.channel,
          input.status ?? "draft",
          input.provider,
          input.subject ?? null,
          input.body,
          input.tenantSlug
        ]
      )
      .catch((error: { code?: string }) => {
        if (error.code === "23505") {
          throw new Error("template_key_conflict");
        }
        throw error;
      });

    return this.mapRow(result.rows[0]);
  }

  async updateStatus(publicId: string, status: TemplateStatus): Promise<EngagementTemplate | null> {
    const result = await this.pool.query<TemplateRow>(
      `
        UPDATE engagement.templates AS template
        SET status = $2
        FROM identity.tenants AS tenant
        WHERE tenant.id = template.tenant_id
          AND template.public_id = $1::uuid
        RETURNING
          template.id,
          template.public_id::text,
          tenant.slug AS tenant_slug,
          template.key,
          template.name,
          template.channel,
          template.status,
          template.provider,
          template.subject,
          template.body,
          template.created_at,
          template.updated_at
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
      `SELECT id FROM identity.tenants WHERE slug = $1 LIMIT 1`,
      [slug]
    );

    if (result.rowCount === 0) {
      throw new Error("template_tenant_not_found");
    }

    const tenantId = Number(result.rows[0].id);
    this.cachedTenantIds.set(slug, tenantId);
    return tenantId;
  }

  private mapRow(row: TemplateRow): EngagementTemplate {
    return {
      id: Number(row.id),
      publicId: row.public_id,
      tenantSlug: row.tenant_slug,
      key: row.key,
      name: row.name,
      channel: row.channel,
      status: row.status,
      provider: row.provider,
      subject: row.subject,
      body: row.body,
      createdAt: row.created_at.toISOString(),
      updatedAt: row.updated_at.toISOString()
    };
  }
}
