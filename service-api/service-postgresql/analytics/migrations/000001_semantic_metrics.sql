CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.metric_definitions (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL UNIQUE,
  code TEXT NOT NULL UNIQUE,
  domain TEXT NOT NULL,
  owner TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  unit TEXT NOT NULL,
  grain TEXT NOT NULL,
  formula TEXT NOT NULL,
  sources_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  dimensions_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  freshness_target_minutes INTEGER NOT NULL DEFAULT 60,
  quality_policy TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics.metric_snapshots (
  id BIGSERIAL PRIMARY KEY,
  metric_id BIGINT NOT NULL REFERENCES analytics.metric_definitions(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  tenant_id BIGINT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  period_key TEXT NOT NULL,
  value_numeric NUMERIC NOT NULL,
  quality TEXT NOT NULL DEFAULT 'passed',
  captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics.dataset_freshness (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL UNIQUE,
  dataset_key TEXT NOT NULL UNIQUE,
  domain TEXT NOT NULL,
  freshness_minutes INTEGER NOT NULL DEFAULT 0,
  target_minutes INTEGER NOT NULL DEFAULT 60,
  status TEXT NOT NULL DEFAULT 'fresh',
  captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics.data_quality_checks (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL UNIQUE,
  check_key TEXT NOT NULL UNIQUE,
  domain TEXT NOT NULL,
  severity TEXT NOT NULL,
  status TEXT NOT NULL,
  failed_rows BIGINT NOT NULL DEFAULT 0,
  captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics.metric_lineage (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL UNIQUE,
  from_node TEXT NOT NULL,
  to_node TEXT NOT NULL,
  edge_type TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
