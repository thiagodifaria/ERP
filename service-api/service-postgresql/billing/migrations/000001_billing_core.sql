CREATE SCHEMA IF NOT EXISTS billing;

CREATE TABLE IF NOT EXISTS billing.plans (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL,
  code VARCHAR(80) NOT NULL,
  name VARCHAR(160) NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  amount_cents BIGINT NOT NULL,
  currency_code VARCHAR(8) NOT NULL DEFAULT 'BRL',
  interval_unit VARCHAR(24) NOT NULL,
  interval_count INT NOT NULL DEFAULT 1,
  grace_period_days INT NOT NULL DEFAULT 0,
  max_retries INT NOT NULL DEFAULT 0,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_plans_public_id UNIQUE (public_id),
  CONSTRAINT uq_billing_plans_code UNIQUE (code),
  CONSTRAINT ck_billing_plans_amount_cents CHECK (amount_cents > 0),
  CONSTRAINT ck_billing_plans_interval_unit CHECK (interval_unit IN ('monthly', 'yearly')),
  CONSTRAINT ck_billing_plans_interval_count CHECK (interval_count > 0),
  CONSTRAINT ck_billing_plans_grace_period_days CHECK (grace_period_days >= 0),
  CONSTRAINT ck_billing_plans_max_retries CHECK (max_retries >= 0)
);

CREATE INDEX IF NOT EXISTS idx_billing_plans_active
  ON billing.plans (active);

CREATE TABLE IF NOT EXISTS billing.subscriptions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  public_id UUID NOT NULL,
  plan_id BIGINT NOT NULL REFERENCES billing.plans (id),
  external_reference VARCHAR(160) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  current_period_start DATE NOT NULL,
  current_period_end DATE NOT NULL,
  grace_ends_at TIMESTAMPTZ,
  suspended_at TIMESTAMPTZ,
  cancelled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_subscriptions_public_id UNIQUE (public_id),
  CONSTRAINT ck_billing_subscriptions_status CHECK (status IN ('active', 'grace_period', 'suspended', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_billing_subscriptions_tenant_status
  ON billing.subscriptions (tenant_id, status);

CREATE TABLE IF NOT EXISTS billing.subscription_invoices (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  subscription_id BIGINT NOT NULL REFERENCES billing.subscriptions (id),
  public_id UUID NOT NULL,
  number VARCHAR(80) NOT NULL,
  status VARCHAR(24) NOT NULL,
  amount_cents BIGINT NOT NULL,
  currency_code VARCHAR(8) NOT NULL DEFAULT 'BRL',
  retry_count INT NOT NULL DEFAULT 0,
  due_date DATE NOT NULL,
  issued_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  paid_at TIMESTAMPTZ,
  last_attempt_at TIMESTAMPTZ,
  gateway_reference VARCHAR(160) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_subscription_invoices_public_id UNIQUE (public_id),
  CONSTRAINT uq_billing_subscription_invoices_tenant_number UNIQUE (tenant_id, number),
  CONSTRAINT ck_billing_subscription_invoices_status CHECK (status IN ('draft', 'open', 'paid', 'failed', 'void')),
  CONSTRAINT ck_billing_subscription_invoices_amount_cents CHECK (amount_cents > 0),
  CONSTRAINT ck_billing_subscription_invoices_retry_count CHECK (retry_count >= 0)
);

CREATE INDEX IF NOT EXISTS idx_billing_subscription_invoices_tenant_status
  ON billing.subscription_invoices (tenant_id, status);

CREATE TABLE IF NOT EXISTS billing.payment_attempts (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  invoice_id BIGINT NOT NULL REFERENCES billing.subscription_invoices (id),
  public_id UUID NOT NULL,
  attempt_number INT NOT NULL,
  provider VARCHAR(80) NOT NULL,
  status VARCHAR(24) NOT NULL,
  idempotency_key VARCHAR(160) NOT NULL,
  external_reference VARCHAR(160) NOT NULL DEFAULT '',
  failure_reason TEXT NOT NULL DEFAULT '',
  attempted_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_payment_attempts_public_id UNIQUE (public_id),
  CONSTRAINT uq_billing_payment_attempts_idempotency_key UNIQUE (idempotency_key),
  CONSTRAINT ck_billing_payment_attempts_attempt_number CHECK (attempt_number > 0),
  CONSTRAINT ck_billing_payment_attempts_status CHECK (status IN ('succeeded', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_billing_payment_attempts_invoice
  ON billing.payment_attempts (invoice_id, attempt_number);

CREATE TABLE IF NOT EXISTS billing.subscription_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  public_id UUID NOT NULL,
  subscription_public_id UUID NOT NULL,
  invoice_public_id UUID,
  event_code VARCHAR(80) NOT NULL,
  actor VARCHAR(80) NOT NULL,
  summary TEXT NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_subscription_events_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_billing_subscription_events_subscription
  ON billing.subscription_events (subscription_public_id, created_at DESC);
