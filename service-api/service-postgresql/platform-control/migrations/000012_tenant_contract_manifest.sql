CREATE TABLE IF NOT EXISTS platform_control.tenant_contract_manifest (
  id BIGSERIAL PRIMARY KEY,
  schema_name TEXT NOT NULL,
  table_name TEXT NOT NULL,
  tenant_strategy TEXT NOT NULL,
  tenant_column TEXT,
  justification TEXT NOT NULL DEFAULT '',
  reviewed_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_platform_control_tenant_contract_manifest UNIQUE (schema_name, table_name),
  CONSTRAINT ck_platform_control_tenant_contract_manifest_strategy
    CHECK (tenant_strategy IN ('tenant_id', 'global_reference', 'system_operational', 'derived_by_join'))
);

INSERT INTO platform_control.tenant_contract_manifest (schema_name, table_name, tenant_strategy, tenant_column, justification)
VALUES
  ('identity', 'tenants', 'global_reference', NULL, 'Tenant registry is the source of tenant identity.'),
  ('platform_control', 'provider_activation_runs', 'tenant_id', 'tenant_id', 'Provider activation evidence is scoped by tenant.'),
  ('platform_control', 'event_mesh_events', 'tenant_id', 'tenant_id', 'Event mesh records are tenant operational data.'),
  ('analytics', 'semantic_metrics', 'tenant_id', 'tenant_id', 'Semantic metrics can be tenant scoped or global by design.'),
  ('workflow_control', 'workflow_runs', 'tenant_id', 'tenant_id', 'Workflow runs are tenant scoped runtime state.'),
  ('workflow_runtime', 'executions', 'tenant_id', 'tenant_id', 'Runtime executions are tenant scoped operational state.'),
  ('crm', 'outbox_events', 'tenant_id', 'tenant_id', 'CRM outbox records are tenant scoped integration state.'),
  ('sales', 'outbox_events', 'tenant_id', 'tenant_id', 'Sales outbox records are tenant scoped integration state.'),
  ('rentals', 'outbox_events', 'tenant_id', 'tenant_id', 'Rentals outbox records are tenant scoped integration state.')
ON CONFLICT (schema_name, table_name) DO UPDATE SET
  tenant_strategy = EXCLUDED.tenant_strategy,
  tenant_column = EXCLUDED.tenant_column,
  justification = EXCLUDED.justification,
  reviewed_at = timezone('utc', now());
