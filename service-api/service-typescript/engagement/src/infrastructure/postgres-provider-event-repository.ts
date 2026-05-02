import pg from "pg";
import { ProviderEventRepository } from "../domain/provider-event-repository.js";
import {
  buildProviderEventSummary,
  CreateProviderEventInput,
  ensureProviderEventInput,
  ProviderEvent,
  ProviderEventFilters,
  ProviderEventSummary
} from "../domain/provider-event.js";

const { Pool } = pg;

type ProviderEventRow = {
  id: number;
  public_id: string;
  tenant_slug: string;
  provider: ProviderEvent["provider"];
  event_type: string;
  direction: ProviderEvent["direction"];
  external_event_id: string | null;
  lead_public_id: string | null;
  business_entity_type: string | null;
  business_entity_public_id: string | null;
  touchpoint_public_id: string | null;
  delivery_public_id: string | null;
  workflow_run_public_id: string | null;
  status: ProviderEvent["status"];
  payload_summary: string;
  response_summary: string;
  created_at: Date;
  processed_at: Date | null;
};

export class PostgresProviderEventRepository implements ProviderEventRepository {
  private readonly cachedTenantIds = new Map<string, number>();

  constructor(
    private readonly pool: InstanceType<typeof Pool>,
    private readonly bootstrapTenantSlug: string
  ) {}

  async list(filters: ProviderEventFilters = {}): Promise<ProviderEvent[]> {
    const tenantId = await this.resolveTenantId(filters.tenantSlug);
    const conditions = ["event.tenant_id = $1"];
    const params: Array<number | string> = [tenantId];

    if (filters.provider) {
      params.push(filters.provider);
      conditions.push(`event.provider = $${params.length}`);
    }
    if (filters.eventType) {
      params.push(filters.eventType);
      conditions.push(`event.event_type = $${params.length}`);
    }
    if (filters.direction) {
      params.push(filters.direction);
      conditions.push(`event.direction = $${params.length}`);
    }
    if (filters.status) {
      params.push(filters.status);
      conditions.push(`event.status = $${params.length}`);
    }
    if (filters.businessEntityType) {
      params.push(filters.businessEntityType);
      conditions.push(`event.business_entity_type = $${params.length}`);
    }
    if (filters.businessEntityPublicId) {
      params.push(filters.businessEntityPublicId);
      conditions.push(`event.business_entity_public_id = $${params.length}::uuid`);
    }

    const result = await this.pool.query<ProviderEventRow>(
      `
        SELECT
          event.id,
          event.public_id::text,
          tenant.slug AS tenant_slug,
          event.provider,
          event.event_type,
          event.direction,
          COALESCE(event.external_event_id, '') AS external_event_id,
          COALESCE(event.lead_public_id::text, '') AS lead_public_id,
          COALESCE(event.business_entity_type, '') AS business_entity_type,
          COALESCE(event.business_entity_public_id::text, '') AS business_entity_public_id,
          COALESCE(event.touchpoint_public_id::text, '') AS touchpoint_public_id,
          COALESCE(event.delivery_public_id::text, '') AS delivery_public_id,
          COALESCE(event.workflow_run_public_id::text, '') AS workflow_run_public_id,
          event.status,
          event.payload_summary,
          event.response_summary,
          event.created_at,
          event.processed_at
        FROM engagement.provider_events AS event
        INNER JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
        WHERE ${conditions.join(" AND ")}
        ORDER BY event.id
      `,
      params
    );

    return result.rows.map((row) => this.mapRow(row));
  }

  async getByPublicId(publicId: string): Promise<ProviderEvent | null> {
    const result = await this.pool.query<ProviderEventRow>(
      `
        SELECT
          event.id,
          event.public_id::text,
          tenant.slug AS tenant_slug,
          event.provider,
          event.event_type,
          event.direction,
          COALESCE(event.external_event_id, '') AS external_event_id,
          COALESCE(event.lead_public_id::text, '') AS lead_public_id,
          COALESCE(event.business_entity_type, '') AS business_entity_type,
          COALESCE(event.business_entity_public_id::text, '') AS business_entity_public_id,
          COALESCE(event.touchpoint_public_id::text, '') AS touchpoint_public_id,
          COALESCE(event.delivery_public_id::text, '') AS delivery_public_id,
          COALESCE(event.workflow_run_public_id::text, '') AS workflow_run_public_id,
          event.status,
          event.payload_summary,
          event.response_summary,
          event.created_at,
          event.processed_at
        FROM engagement.provider_events AS event
        INNER JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
        WHERE event.public_id = $1::uuid
        LIMIT 1
      `,
      [publicId]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  async findByProviderAndExternalEventId(tenantSlug: string, provider: string, externalEventId: string): Promise<ProviderEvent | null> {
    const normalizedExternalEventId = externalEventId.trim();
    if (normalizedExternalEventId.length === 0) {
      return null;
    }

    const tenantId = await this.resolveTenantId(tenantSlug);
    const result = await this.pool.query<ProviderEventRow>(
      `
        SELECT
          event.id,
          event.public_id::text,
          tenant.slug AS tenant_slug,
          event.provider,
          event.event_type,
          event.direction,
          COALESCE(event.external_event_id, '') AS external_event_id,
          COALESCE(event.lead_public_id::text, '') AS lead_public_id,
          COALESCE(event.business_entity_type, '') AS business_entity_type,
          COALESCE(event.business_entity_public_id::text, '') AS business_entity_public_id,
          COALESCE(event.touchpoint_public_id::text, '') AS touchpoint_public_id,
          COALESCE(event.delivery_public_id::text, '') AS delivery_public_id,
          COALESCE(event.workflow_run_public_id::text, '') AS workflow_run_public_id,
          event.status,
          event.payload_summary,
          event.response_summary,
          event.created_at,
          event.processed_at
        FROM engagement.provider_events AS event
        INNER JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
        WHERE event.tenant_id = $1
          AND event.provider = $2
          AND event.external_event_id = $3
        LIMIT 1
      `,
      [tenantId, provider, normalizedExternalEventId]
    );

    return result.rowCount === 0 ? null : this.mapRow(result.rows[0]);
  }

  async create(input: CreateProviderEventInput): Promise<ProviderEvent> {
    const normalized = ensureProviderEventInput(input);
    const tenantId = await this.resolveTenantId(normalized.tenantSlug);

    if (normalized.externalEventId && (await this.findByProviderAndExternalEventId(normalized.tenantSlug, normalized.provider, normalized.externalEventId)) !== null) {
      throw new Error("provider_event_conflict");
    }

    const result = await this.pool.query<ProviderEventRow>(
      `
        INSERT INTO engagement.provider_events (
          tenant_id,
          public_id,
          provider,
          event_type,
          direction,
          external_event_id,
          lead_public_id,
          business_entity_type,
          business_entity_public_id,
          touchpoint_public_id,
          delivery_public_id,
          workflow_run_public_id,
          status,
          payload_summary,
          response_summary,
          processed_at
        )
        VALUES (
          $1,
          gen_random_uuid(),
          $2,
          $3,
          $4,
          NULLIF($5, ''),
          CASE WHEN $6 = '' THEN NULL ELSE $6::uuid END,
          NULLIF($7, ''),
          CASE WHEN $8 = '' THEN NULL ELSE $8::uuid END,
          CASE WHEN $9 = '' THEN NULL ELSE $9::uuid END,
          CASE WHEN $10 = '' THEN NULL ELSE $10::uuid END,
          CASE WHEN $11 = '' THEN NULL ELSE $11::uuid END,
          $12,
          $13,
          $14,
          CASE WHEN $15 = '' THEN NULL ELSE $15::timestamptz END
        )
        RETURNING
          id,
          public_id::text,
          $16::text AS tenant_slug,
          provider,
          event_type,
          direction,
          COALESCE(external_event_id, '') AS external_event_id,
          COALESCE(lead_public_id::text, '') AS lead_public_id,
          COALESCE(business_entity_type, '') AS business_entity_type,
          COALESCE(business_entity_public_id::text, '') AS business_entity_public_id,
          COALESCE(touchpoint_public_id::text, '') AS touchpoint_public_id,
          COALESCE(delivery_public_id::text, '') AS delivery_public_id,
          COALESCE(workflow_run_public_id::text, '') AS workflow_run_public_id,
          status,
          payload_summary,
          response_summary,
          created_at,
          processed_at
      `,
      [
        tenantId,
        normalized.provider,
        normalized.eventType,
        normalized.direction,
        normalized.externalEventId ?? "",
        normalized.leadPublicId ?? "",
        normalized.businessEntityType ?? "",
        normalized.businessEntityPublicId ?? "",
        normalized.touchpointPublicId ?? "",
        normalized.deliveryPublicId ?? "",
        normalized.workflowRunPublicId ?? "",
        normalized.status,
        normalized.payloadSummary,
        normalized.responseSummary,
        normalized.processedAt ?? "",
        normalized.tenantSlug
      ]
    );

    return this.mapRow(result.rows[0]);
  }

  async getSummary(filters: ProviderEventFilters = {}): Promise<ProviderEventSummary> {
    return buildProviderEventSummary(await this.list(filters), filters.tenantSlug ?? this.bootstrapTenantSlug);
  }

  private async resolveTenantId(tenantSlug?: string): Promise<number> {
    const slug = tenantSlug ?? this.bootstrapTenantSlug;
    const cached = this.cachedTenantIds.get(slug);
    if (cached) {
      return cached;
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
      throw new Error("provider_event_tenant_not_found");
    }

    const tenantId = Number(result.rows[0].id);
    this.cachedTenantIds.set(slug, tenantId);
    return tenantId;
  }

  private mapRow(row: ProviderEventRow): ProviderEvent {
    return {
      id: row.id,
      publicId: row.public_id,
      tenantSlug: row.tenant_slug,
      provider: row.provider,
      eventType: row.event_type,
      direction: row.direction,
      externalEventId: row.external_event_id || null,
      leadPublicId: row.lead_public_id || null,
      businessEntityType: row.business_entity_type || null,
      businessEntityPublicId: row.business_entity_public_id || null,
      touchpointPublicId: row.touchpoint_public_id || null,
      deliveryPublicId: row.delivery_public_id || null,
      workflowRunPublicId: row.workflow_run_public_id || null,
      status: row.status,
      payloadSummary: row.payload_summary,
      responseSummary: row.response_summary,
      createdAt: row.created_at.toISOString(),
      processedAt: row.processed_at ? row.processed_at.toISOString() : null
    };
  }
}
