CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.reconciliation_runs (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    status TEXT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    snapshot_hash TEXT NOT NULL,
    summary JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS analytics.reconciliation_findings (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    run_public_id UUID,
    finding_key TEXT NOT NULL,
    domain TEXT NOT NULL,
    severity TEXT NOT NULL,
    count_value INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analytics.financial_close_periods (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    period_key TEXT NOT NULL,
    status TEXT NOT NULL,
    readiness TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at TIMESTAMPTZ,
    UNIQUE (tenant_slug, period_key)
);

CREATE TABLE IF NOT EXISTS analytics.financial_close_snapshots (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    period_public_id UUID NOT NULL,
    period_key TEXT NOT NULL,
    snapshot_hash TEXT NOT NULL,
    totals JSONB NOT NULL DEFAULT '{}'::jsonb,
    controls JSONB NOT NULL DEFAULT '[]'::jsonb,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analytics.master_data_quality_scores (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    owner_service TEXT NOT NULL,
    score NUMERIC(5,2) NOT NULL,
    dimension_scores JSONB NOT NULL DEFAULT '{}'::jsonb,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analytics.master_data_findings (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    rule_key TEXT NOT NULL,
    severity TEXT NOT NULL,
    status TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analytics.master_data_merge_proposals (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    source_records JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT NOT NULL,
    proposal_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analytics.lakehouse_datasets (
    id BIGSERIAL PRIMARY KEY,
    dataset_key TEXT NOT NULL UNIQUE,
    owner_service TEXT NOT NULL,
    classification JSONB NOT NULL DEFAULT '[]'::jsonb,
    freshness_target TEXT NOT NULL,
    retention TEXT NOT NULL,
    source_refs JSONB NOT NULL DEFAULT '[]'::jsonb,
    export_policy TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reconciliation_findings_tenant_status ON analytics.reconciliation_findings (tenant_slug, status, severity);
CREATE INDEX IF NOT EXISTS idx_financial_close_tenant_period ON analytics.financial_close_periods (tenant_slug, period_key);
CREATE INDEX IF NOT EXISTS idx_master_data_scores_tenant_entity ON analytics.master_data_quality_scores (tenant_slug, entity_type);
CREATE INDEX IF NOT EXISTS idx_master_data_findings_status ON analytics.master_data_findings (tenant_slug, status, severity);
CREATE INDEX IF NOT EXISTS idx_lakehouse_datasets_owner ON analytics.lakehouse_datasets (owner_service);
