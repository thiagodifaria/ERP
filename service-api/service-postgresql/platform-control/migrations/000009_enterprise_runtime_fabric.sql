CREATE SCHEMA IF NOT EXISTS platform_control;

CREATE TABLE IF NOT EXISTS platform_control.event_mesh_events (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    stream_key TEXT NOT NULL,
    event_type TEXT NOT NULL,
    schema_version TEXT NOT NULL,
    producer TEXT NOT NULL,
    consumer TEXT,
    correlation_id TEXT NOT NULL,
    causation_id TEXT,
    status TEXT NOT NULL,
    payload_hash TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS platform_control.event_mesh_dead_letters (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    event_public_id UUID NOT NULL,
    stream_key TEXT NOT NULL,
    event_type TEXT NOT NULL,
    producer TEXT NOT NULL,
    reason TEXT NOT NULL,
    status TEXT NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 1,
    payload_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    replayed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS platform_control.tenant_runtime_profiles (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL UNIQUE,
    plan_key TEXT NOT NULL,
    status TEXT NOT NULL,
    modules JSONB NOT NULL DEFAULT '[]'::jsonb,
    feature_flags JSONB NOT NULL DEFAULT '{}'::jsonb,
    slo_profile JSONB NOT NULL DEFAULT '{}'::jsonb,
    risk_status TEXT NOT NULL,
    policy_set JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS platform_control.tenant_runtime_quotas (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    metric_key TEXT NOT NULL,
    limit_value BIGINT NOT NULL,
    usage_pct INTEGER NOT NULL DEFAULT 0,
    enforcement_mode TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_slug, metric_key)
);

CREATE TABLE IF NOT EXISTS platform_control.tenant_maintenance_windows (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    title TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    impact TEXT NOT NULL,
    status TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS platform_control.contract_snapshots (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    contract_key TEXT NOT NULL,
    version TEXT NOT NULL,
    kind TEXT NOT NULL,
    service TEXT NOT NULL,
    payload_hash TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (contract_key, version)
);

CREATE TABLE IF NOT EXISTS platform_control.contract_diffs (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    contract_key TEXT NOT NULL,
    from_version TEXT NOT NULL,
    to_version TEXT NOT NULL,
    added_operations JSONB NOT NULL DEFAULT '[]'::jsonb,
    removed_operations JSONB NOT NULL DEFAULT '[]'::jsonb,
    changed_schemas JSONB NOT NULL DEFAULT '[]'::jsonb,
    breaking BOOLEAN NOT NULL DEFAULT false,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS platform_control.contract_breaking_changes (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    contract_key TEXT NOT NULL,
    diff_public_id UUID NOT NULL,
    severity TEXT NOT NULL,
    status TEXT NOT NULL,
    summary TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_by TEXT,
    approved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_event_mesh_events_tenant_stream ON platform_control.event_mesh_events (tenant_slug, stream_key, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_event_mesh_events_correlation ON platform_control.event_mesh_events (tenant_slug, correlation_id);
CREATE INDEX IF NOT EXISTS idx_event_mesh_dead_letters_status ON platform_control.event_mesh_dead_letters (tenant_slug, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_runtime_quotas_tenant ON platform_control.tenant_runtime_quotas (tenant_slug, metric_key);
CREATE INDEX IF NOT EXISTS idx_contract_diffs_key ON platform_control.contract_diffs (contract_key, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_contract_breaking_changes_status ON platform_control.contract_breaking_changes (status, created_at DESC);
