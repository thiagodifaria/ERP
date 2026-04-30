CREATE SCHEMA IF NOT EXISTS rentals;

CREATE TABLE IF NOT EXISTS rentals.contracts (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  customer_public_id UUID NOT NULL,
  title VARCHAR(160) NOT NULL,
  property_code VARCHAR(80) NOT NULL DEFAULT '',
  currency_code CHAR(3) NOT NULL DEFAULT 'BRL',
  amount_cents BIGINT NOT NULL,
  billing_day SMALLINT NOT NULL,
  starts_at DATE NOT NULL,
  ends_at DATE NOT NULL,
  status VARCHAR(30) NOT NULL,
  terminated_at TIMESTAMPTZ NULL,
  termination_reason TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_rentals_contracts_public_id UNIQUE (public_id),
  CONSTRAINT ck_rentals_contracts_amount CHECK (amount_cents > 0),
  CONSTRAINT ck_rentals_contracts_billing_day CHECK (billing_day BETWEEN 1 AND 31),
  CONSTRAINT ck_rentals_contracts_status CHECK (status IN ('active', 'terminated')),
  CONSTRAINT ck_rentals_contracts_dates CHECK (ends_at >= starts_at)
);

CREATE INDEX IF NOT EXISTS idx_rentals_contracts_tenant_status
  ON rentals.contracts (tenant_id, status, created_at);

CREATE INDEX IF NOT EXISTS idx_rentals_contracts_customer_public_id
  ON rentals.contracts (customer_public_id);

CREATE TABLE IF NOT EXISTS rentals.contract_charges (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  contract_id BIGINT NOT NULL REFERENCES rentals.contracts(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  due_date DATE NOT NULL,
  amount_cents BIGINT NOT NULL,
  status VARCHAR(30) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_rentals_contract_charges_public_id UNIQUE (public_id),
  CONSTRAINT ck_rentals_contract_charges_amount CHECK (amount_cents > 0),
  CONSTRAINT ck_rentals_contract_charges_status CHECK (status IN ('scheduled', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_rentals_contract_charges_tenant_status
  ON rentals.contract_charges (tenant_id, status, due_date);

CREATE INDEX IF NOT EXISTS idx_rentals_contract_charges_contract
  ON rentals.contract_charges (contract_id, due_date);

CREATE TABLE IF NOT EXISTS rentals.contract_adjustments (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  contract_id BIGINT NOT NULL REFERENCES rentals.contracts(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  effective_at DATE NOT NULL,
  previous_amount_cents BIGINT NOT NULL,
  new_amount_cents BIGINT NOT NULL,
  reason TEXT NOT NULL,
  recorded_by VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_rentals_contract_adjustments_public_id UNIQUE (public_id),
  CONSTRAINT ck_rentals_contract_adjustments_previous_amount CHECK (previous_amount_cents > 0),
  CONSTRAINT ck_rentals_contract_adjustments_new_amount CHECK (new_amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_rentals_contract_adjustments_contract
  ON rentals.contract_adjustments (contract_id, effective_at, created_at);

CREATE TABLE IF NOT EXISTS rentals.contract_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  contract_id BIGINT NOT NULL REFERENCES rentals.contracts(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  event_code VARCHAR(60) NOT NULL,
  summary TEXT NOT NULL,
  recorded_by VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_rentals_contract_events_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_rentals_contract_events_contract
  ON rentals.contract_events (contract_id, created_at);

CREATE TABLE IF NOT EXISTS rentals.outbox_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  aggregate_type VARCHAR(60) NOT NULL,
  aggregate_public_id UUID NOT NULL,
  event_type VARCHAR(80) NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  processed_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_rentals_outbox_events_public_id UNIQUE (public_id),
  CONSTRAINT ck_rentals_outbox_events_status CHECK (status IN ('pending', 'processed', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_rentals_outbox_events_tenant_status
  ON rentals.outbox_events (tenant_id, status, created_at);
