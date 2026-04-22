-- Expande o contexto financeiro para o ciclo operacional.
CREATE TABLE IF NOT EXISTS finance.receivable_entries (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  source_invoice_public_id UUID NOT NULL,
  sale_public_id UUID NOT NULL,
  customer_public_id UUID NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'open',
  amount_cents BIGINT NOT NULL,
  due_date DATE NOT NULL,
  paid_at TIMESTAMPTZ NULL,
  last_synced_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_receivable_entries_public_id UNIQUE (public_id),
  CONSTRAINT uq_finance_receivable_entries_source_invoice UNIQUE (source_invoice_public_id),
  CONSTRAINT ck_finance_receivable_entries_status CHECK (status IN ('open', 'paid', 'cancelled')),
  CONSTRAINT ck_finance_receivable_entries_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_finance_receivable_entries_tenant_status
  ON finance.receivable_entries (tenant_id, status, due_date);

CREATE INDEX IF NOT EXISTS idx_finance_receivable_entries_sale_public_id
  ON finance.receivable_entries (sale_public_id);

CREATE TABLE IF NOT EXISTS finance.receivable_settlements (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  receivable_entry_id BIGINT NOT NULL REFERENCES finance.receivable_entries(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  settlement_reference VARCHAR(120) NOT NULL,
  amount_cents BIGINT NOT NULL,
  settled_at TIMESTAMPTZ NOT NULL,
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_receivable_settlements_public_id UNIQUE (public_id),
  CONSTRAINT uq_finance_receivable_settlements_reference UNIQUE (tenant_id, settlement_reference),
  CONSTRAINT uq_finance_receivable_settlements_receivable UNIQUE (receivable_entry_id),
  CONSTRAINT ck_finance_receivable_settlements_amount CHECK (amount_cents > 0)
);

CREATE TABLE IF NOT EXISTS finance.commission_entries (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  source_commission_public_id UUID NOT NULL,
  sale_public_id UUID NOT NULL,
  recipient_user_public_id UUID NOT NULL,
  role_code VARCHAR(60) NOT NULL,
  rate_bps INTEGER NOT NULL,
  amount_cents BIGINT NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_commission_entries_public_id UNIQUE (public_id),
  CONSTRAINT uq_finance_commission_entries_source UNIQUE (source_commission_public_id),
  CONSTRAINT ck_finance_commission_entries_rate_bps CHECK (rate_bps > 0 AND rate_bps <= 10000),
  CONSTRAINT ck_finance_commission_entries_amount CHECK (amount_cents > 0),
  CONSTRAINT ck_finance_commission_entries_status CHECK (status IN ('pending', 'blocked', 'released'))
);

CREATE INDEX IF NOT EXISTS idx_finance_commission_entries_tenant_status
  ON finance.commission_entries (tenant_id, status, created_at);

CREATE TABLE IF NOT EXISTS finance.payables (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  category VARCHAR(60) NOT NULL,
  vendor_name VARCHAR(120) NOT NULL,
  description TEXT NOT NULL,
  amount_cents BIGINT NOT NULL,
  due_date DATE NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'open',
  payment_reference VARCHAR(120) NULL,
  paid_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_payables_public_id UNIQUE (public_id),
  CONSTRAINT ck_finance_payables_amount CHECK (amount_cents > 0),
  CONSTRAINT ck_finance_payables_status CHECK (status IN ('open', 'paid', 'cancelled'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_finance_payables_payment_reference
  ON finance.payables (tenant_id, payment_reference)
  WHERE payment_reference IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_finance_payables_tenant_status
  ON finance.payables (tenant_id, status, due_date);

CREATE TABLE IF NOT EXISTS finance.cost_entries (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  category VARCHAR(60) NOT NULL,
  summary TEXT NOT NULL,
  amount_cents BIGINT NOT NULL,
  incurred_on DATE NOT NULL,
  sale_public_id UUID NULL,
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_cost_entries_public_id UNIQUE (public_id),
  CONSTRAINT ck_finance_cost_entries_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_finance_cost_entries_tenant_incurred_on
  ON finance.cost_entries (tenant_id, incurred_on DESC, created_at DESC);

CREATE TABLE IF NOT EXISTS finance.period_closures (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  period_key VARCHAR(7) NOT NULL,
  closed_at TIMESTAMPTZ NOT NULL,
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_period_closures_public_id UNIQUE (public_id),
  CONSTRAINT uq_finance_period_closures_tenant_period UNIQUE (tenant_id, period_key)
);

CREATE INDEX IF NOT EXISTS idx_finance_period_closures_tenant_created_at
  ON finance.period_closures (tenant_id, created_at DESC);
