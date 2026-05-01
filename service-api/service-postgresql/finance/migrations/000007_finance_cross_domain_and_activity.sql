ALTER TABLE finance.receivable_entries
  ADD COLUMN IF NOT EXISTS source_kind VARCHAR(40) NOT NULL DEFAULT 'sales_invoice',
  ADD COLUMN IF NOT EXISTS contract_public_id UUID NULL;

ALTER TABLE finance.receivable_entries
  ALTER COLUMN sale_public_id DROP NOT NULL;

ALTER TABLE finance.receivable_entries
  DROP CONSTRAINT IF EXISTS ck_finance_receivable_entries_source_kind;

ALTER TABLE finance.receivable_entries
  ADD CONSTRAINT ck_finance_receivable_entries_source_kind
  CHECK (source_kind IN ('sales_invoice', 'rental_charge'));

CREATE INDEX IF NOT EXISTS idx_finance_receivable_entries_contract_public_id
  ON finance.receivable_entries (contract_public_id);

CREATE TABLE IF NOT EXISTS finance.activity_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  activity_type VARCHAR(80) NOT NULL,
  entity_type VARCHAR(80) NOT NULL,
  entity_public_id UUID NULL,
  summary TEXT NOT NULL,
  actor VARCHAR(120) NOT NULL DEFAULT 'system',
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_activity_events_public_id UNIQUE (public_id),
  CONSTRAINT ck_finance_activity_events_type CHECK (
    activity_type IN (
      'receivable_synced',
      'receivable_settled',
      'commission_blocked',
      'commission_released',
      'payable_created',
      'payable_status_changed',
      'cost_created',
      'treasury_synced',
      'period_closed'
    )
  )
);

CREATE INDEX IF NOT EXISTS idx_finance_activity_events_tenant_created_at
  ON finance.activity_events (tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_finance_activity_events_entity
  ON finance.activity_events (tenant_id, entity_type, entity_public_id, created_at DESC);
