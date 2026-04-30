import pg from "pg";
import { DeliveryRepository } from "../domain/delivery-repository.js";
import { buildDeliverySummary, CreateDeliveryInput, DeliveryFilters, DeliverySummary, TouchpointDelivery, UpdateDeliveryStatusInput } from "../domain/delivery.js";
import { EngagementTemplate } from "../domain/template.js";
import { Touchpoint } from "../domain/touchpoint.js";

const { Pool } = pg;

type DeliveryRow = {
  id: number;
  public_id: string;
  tenant_slug: string;
  touchpoint_public_id: string;
  template_public_id: string | null;
  template_key: string | null;
  channel: TouchpointDelivery["channel"];
  provider: TouchpointDelivery["provider"];
  provider_message_id: string | null;
  status: TouchpointDelivery["status"];
  sent_by: string;
  error_code: string | null;
  notes: string;
  attempted_at: Date;
  created_at: Date;
  updated_at: Date;
};

export class PostgresDeliveryRepository implements DeliveryRepository {
  private readonly cachedTenantIds = new Map<string, number>();

  constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string,
    private readonly listTemplatesForTenant: (tenantSlug?: string) => Promise<EngagementTemplate[]>,
    private readonly listTouchpointsForTenant: (tenantSlug?: string) => Promise<Touchpoint[]>
  ) {}

  async list(filters: DeliveryFilters = {}): Promise<TouchpointDelivery[]> {
    const tenantId = await this.resolveTenantId(filters.tenantSlug);
    const conditions = ["delivery.tenant_id = $1"];
    const params: Array<number | string> = [tenantId];

    if (filters.touchpointPublicId) {
      params.push(filters.touchpointPublicId);
      conditions.push(`touchpoint.public_id = $${params.length}::uuid`);
    }
    if (filters.channel) {
      params.push(filters.channel);
      conditions.push(`delivery.channel = $${params.length}`);
    }
    if (filters.provider) {
      params.push(filters.provider);
      conditions.push(`delivery.provider = $${params.length}`);
    }
    if (filters.status) {
      params.push(filters.status);
      conditions.push(`delivery.status = $${params.length}`);
    }

    const result = await this.pool.query<DeliveryRow>(
      `
        SELECT
          delivery.id,
          delivery.public_id::text,
          tenant.slug AS tenant_slug,
          touchpoint.public_id::text AS touchpoint_public_id,
          template.public_id::text AS template_public_id,
          template.key AS template_key,
          delivery.channel,
          delivery.provider,
          delivery.provider_message_id,
          delivery.status,
          delivery.sent_by,
          delivery.error_code,
          delivery.notes,
          delivery.attempted_at,
          delivery.created_at,
          delivery.updated_at
        FROM engagement.touchpoint_deliveries AS delivery
        INNER JOIN engagement.touchpoints AS touchpoint ON touchpoint.id = delivery.touchpoint_id
        LEFT JOIN engagement.templates AS template ON template.id = delivery.template_id
        INNER JOIN identity.tenants AS tenant ON tenant.id = delivery.tenant_id
        WHERE ${conditions.join(" AND ")}
        ORDER BY delivery.id
      `,
      params
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  async getByPublicId(publicId: string): Promise<TouchpointDelivery | null> {
    const result = await this.pool.query<DeliveryRow>(
      `
        SELECT
          delivery.id,
          delivery.public_id::text,
          tenant.slug AS tenant_slug,
          touchpoint.public_id::text AS touchpoint_public_id,
          template.public_id::text AS template_public_id,
          template.key AS template_key,
          delivery.channel,
          delivery.provider,
          delivery.provider_message_id,
          delivery.status,
          delivery.sent_by,
          delivery.error_code,
          delivery.notes,
          delivery.attempted_at,
          delivery.created_at,
          delivery.updated_at
        FROM engagement.touchpoint_deliveries AS delivery
        INNER JOIN engagement.touchpoints AS touchpoint ON touchpoint.id = delivery.touchpoint_id
        LEFT JOIN engagement.templates AS template ON template.id = delivery.template_id
        INNER JOIN identity.tenants AS tenant ON tenant.id = delivery.tenant_id
        WHERE delivery.public_id = $1::uuid
        LIMIT 1
      `,
      [publicId]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  async create(input: CreateDeliveryInput): Promise<TouchpointDelivery> {
    const tenantId = await this.resolveTenantId(input.tenantSlug);
    const touchpointResult = await this.pool.query<{ id: number; channel: TouchpointDelivery["channel"] }>(
      `
        SELECT id, channel
        FROM engagement.touchpoints
        WHERE tenant_id = $1
          AND public_id = $2::uuid
        LIMIT 1
      `,
      [tenantId, input.touchpointPublicId]
    );

    if (touchpointResult.rowCount === 0) {
      throw new Error("touchpoint_not_found");
    }

    let templateId: number | null = null;
    let templatePublicId: string | null = null;
    let templateKey: string | null = null;

    if (input.templatePublicId) {
      const templateResult = await this.pool.query<{ id: number; public_id: string; key: string; channel: TouchpointDelivery["channel"] }>(
        `
          SELECT id, public_id::text, key, channel
          FROM engagement.templates
          WHERE tenant_id = $1
            AND public_id = $2::uuid
          LIMIT 1
        `,
        [tenantId, input.templatePublicId]
      );

      if (templateResult.rowCount === 0) {
        throw new Error("template_not_found");
      }

      if (templateResult.rows[0].channel !== touchpointResult.rows[0].channel) {
        throw new Error("delivery_channel_mismatch");
      }

      templateId = Number(templateResult.rows[0].id);
      templatePublicId = templateResult.rows[0].public_id;
      templateKey = templateResult.rows[0].key;
    }

    const result = await this.pool.query<DeliveryRow>(
      `
        INSERT INTO engagement.touchpoint_deliveries (
          tenant_id,
          touchpoint_id,
          template_id,
          public_id,
          channel,
          provider,
          provider_message_id,
          status,
          sent_by,
          error_code,
          notes,
          attempted_at
        )
        VALUES ($1, $2, $3, gen_random_uuid(), $4, $5, $6, 'sent', $7, NULL, $8, timezone('utc', now()))
        RETURNING
          id,
          public_id::text,
          $9::text AS tenant_slug,
          $10::text AS touchpoint_public_id,
          $11::text AS template_public_id,
          $12::text AS template_key,
          channel,
          provider,
          provider_message_id,
          status,
          sent_by,
          error_code,
          notes,
          attempted_at,
          created_at,
          updated_at
      `,
      [
        tenantId,
        Number(touchpointResult.rows[0].id),
        templateId,
        touchpointResult.rows[0].channel,
        input.provider,
        input.providerMessageId ?? null,
        input.sentBy,
        input.notes ?? "",
        input.tenantSlug,
        input.touchpointPublicId,
        templatePublicId,
        templateKey
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  async updateStatus(publicId: string, input: UpdateDeliveryStatusInput): Promise<TouchpointDelivery | null> {
    const result = await this.pool.query<{ public_id: string }>(
      `
        UPDATE engagement.touchpoint_deliveries AS delivery
        SET status = $2::text,
            provider_message_id = COALESCE($3::text, delivery.provider_message_id),
            error_code = CASE WHEN $2::text = 'failed' THEN $4::text ELSE NULL END,
            notes = CASE WHEN length($5::text) > 0 THEN $5::text ELSE delivery.notes END
        WHERE delivery.public_id = $1::uuid
        RETURNING delivery.public_id::text
      `,
      [publicId, input.status, input.providerMessageId ?? null, input.errorCode ?? null, input.notes ?? ""]
    );

    return result.rowCount === 0 ? null : this.getByPublicId(result.rows[0].public_id);
  }

  async getSummary(filters: DeliveryFilters = {}): Promise<DeliverySummary> {
    const deliveries = await this.list(filters);
    const tenantSlug = filters.tenantSlug ?? deliveries[0]?.tenantSlug ?? this.bootstrapTenantSlug;
    const templates = await this.listTemplatesForTenant(tenantSlug);
    const touchpoints = await this.listTouchpointsForTenant(tenantSlug);
    return buildDeliverySummary(tenantSlug, templates, touchpoints, deliveries);
  }

  async listByTouchpointPublicId(touchpointPublicId: string, tenantSlug?: string): Promise<TouchpointDelivery[]> {
    return this.list({ tenantSlug, touchpointPublicId });
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
      throw new Error("delivery_tenant_not_found");
    }

    const tenantId = Number(result.rows[0].id);
    this.cachedTenantIds.set(slug, tenantId);
    return tenantId;
  }

  private mapRow(row: DeliveryRow): TouchpointDelivery {
    return {
      id: Number(row.id),
      publicId: row.public_id,
      tenantSlug: row.tenant_slug,
      touchpointPublicId: row.touchpoint_public_id,
      templatePublicId: row.template_public_id,
      templateKey: row.template_key,
      channel: row.channel,
      provider: row.provider,
      providerMessageId: row.provider_message_id,
      status: row.status,
      sentBy: row.sent_by,
      errorCode: row.error_code,
      notes: row.notes,
      attemptedAt: row.attempted_at.toISOString(),
      createdAt: row.created_at.toISOString(),
      updatedAt: row.updated_at.toISOString()
    };
  }
}
