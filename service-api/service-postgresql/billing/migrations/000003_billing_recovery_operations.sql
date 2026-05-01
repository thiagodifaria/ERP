CREATE TABLE IF NOT EXISTS billing.recovery_cases (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  public_id UUID NOT NULL,
  subscription_public_id UUID NOT NULL,
  invoice_public_id UUID NOT NULL,
  source_attempt_public_id UUID,
  workflow_definition_key VARCHAR(120) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL,
  severity VARCHAR(24) NOT NULL,
  contact_channel VARCHAR(40) NOT NULL DEFAULT '',
  provider_hint VARCHAR(60) NOT NULL DEFAULT '',
  last_failed_attempt_number INT NOT NULL DEFAULT 0,
  next_action_at TIMESTAMPTZ,
  promised_payment_date DATE,
  resolved_at TIMESTAMPTZ,
  resolution_code VARCHAR(60) NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_recovery_cases_public_id UNIQUE (public_id),
  CONSTRAINT uq_billing_recovery_cases_invoice UNIQUE (tenant_id, invoice_public_id),
  CONSTRAINT ck_billing_recovery_cases_status CHECK (status IN ('open', 'contacted', 'promised', 'recovered', 'defaulted', 'closed')),
  CONSTRAINT ck_billing_recovery_cases_severity CHECK (severity IN ('attention', 'critical')),
  CONSTRAINT ck_billing_recovery_cases_contact_channel CHECK (contact_channel IN ('', 'email', 'whatsapp', 'phone', 'manual')),
  CONSTRAINT ck_billing_recovery_cases_last_failed_attempt_number CHECK (last_failed_attempt_number >= 0)
);

CREATE INDEX IF NOT EXISTS idx_billing_recovery_cases_tenant_status
  ON billing.recovery_cases (tenant_id, status, severity);

CREATE INDEX IF NOT EXISTS idx_billing_recovery_cases_next_action
  ON billing.recovery_cases (tenant_id, next_action_at);

CREATE TABLE IF NOT EXISTS billing.recovery_actions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  recovery_case_id BIGINT NOT NULL REFERENCES billing.recovery_cases (id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  action_code VARCHAR(60) NOT NULL,
  actor VARCHAR(80) NOT NULL,
  channel VARCHAR(40) NOT NULL DEFAULT '',
  provider VARCHAR(60) NOT NULL DEFAULT '',
  touchpoint_public_id UUID,
  delivery_public_id UUID,
  promised_payment_date DATE,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_billing_recovery_actions_public_id UNIQUE (public_id),
  CONSTRAINT ck_billing_recovery_actions_action_code CHECK (action_code IN ('case_opened', 'case_refreshed', 'touchpoint_registered', 'promise_registered', 'case_recovered', 'case_closed', 'note_added')),
  CONSTRAINT ck_billing_recovery_actions_channel CHECK (channel IN ('', 'email', 'whatsapp', 'phone', 'manual'))
);

CREATE INDEX IF NOT EXISTS idx_billing_recovery_actions_case
  ON billing.recovery_actions (recovery_case_id, created_at DESC);
