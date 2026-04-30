-- Abre a trilha de tesouraria e caixa operacional do financeiro.
CREATE TABLE IF NOT EXISTS finance.cash_accounts (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  code VARCHAR(60) NOT NULL,
  display_name VARCHAR(120) NOT NULL,
  currency_code CHAR(3) NOT NULL DEFAULT 'BRL',
  provider VARCHAR(60) NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'active',
  opening_balance_cents BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_cash_accounts_public_id UNIQUE (public_id),
  CONSTRAINT uq_finance_cash_accounts_tenant_code UNIQUE (tenant_id, code),
  CONSTRAINT ck_finance_cash_accounts_currency_code CHECK (currency_code ~ '^[A-Z]{3}$'),
  CONSTRAINT ck_finance_cash_accounts_status CHECK (status IN ('active', 'inactive'))
);

CREATE INDEX IF NOT EXISTS idx_finance_cash_accounts_tenant_status
  ON finance.cash_accounts (tenant_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS finance.cash_movements (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  cash_account_id BIGINT NOT NULL REFERENCES finance.cash_accounts(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  movement_type VARCHAR(60) NOT NULL,
  direction VARCHAR(20) NOT NULL,
  source_public_id UUID NULL,
  reference_code VARCHAR(120) NOT NULL,
  amount_cents BIGINT NOT NULL,
  counterparty_name VARCHAR(160) NOT NULL,
  description TEXT NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL,
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_cash_movements_public_id UNIQUE (public_id),
  CONSTRAINT ck_finance_cash_movements_direction CHECK (direction IN ('inflow', 'outflow')),
  CONSTRAINT ck_finance_cash_movements_amount CHECK (amount_cents > 0),
  CONSTRAINT ck_finance_cash_movements_type CHECK (movement_type IN ('receivable_settlement', 'payable_payment', 'cost_entry', 'manual_adjustment'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_finance_cash_movements_source
  ON finance.cash_movements (tenant_id, movement_type, source_public_id)
  WHERE source_public_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_finance_cash_movements_account_effective_at
  ON finance.cash_movements (cash_account_id, effective_at DESC, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_finance_cash_movements_tenant_direction
  ON finance.cash_movements (tenant_id, direction, effective_at DESC);
