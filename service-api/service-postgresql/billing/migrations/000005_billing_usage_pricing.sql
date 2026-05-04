ALTER TABLE billing.plans
  ADD COLUMN IF NOT EXISTS pricing_model VARCHAR(24) NOT NULL DEFAULT 'flat',
  ADD COLUMN IF NOT EXISTS meter_key VARCHAR(120) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS included_quantity BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS overage_unit_cents BIGINT NOT NULL DEFAULT 0;

ALTER TABLE billing.plans
  DROP CONSTRAINT IF EXISTS ck_billing_plans_pricing_model;

ALTER TABLE billing.plans
  ADD CONSTRAINT ck_billing_plans_pricing_model
  CHECK (pricing_model IN ('flat', 'hybrid', 'usage'));

ALTER TABLE billing.plans
  DROP CONSTRAINT IF EXISTS ck_billing_plans_included_quantity;

ALTER TABLE billing.plans
  ADD CONSTRAINT ck_billing_plans_included_quantity
  CHECK (included_quantity >= 0);

ALTER TABLE billing.plans
  DROP CONSTRAINT IF EXISTS ck_billing_plans_overage_unit_cents;

ALTER TABLE billing.plans
  ADD CONSTRAINT ck_billing_plans_overage_unit_cents
  CHECK (overage_unit_cents >= 0);

CREATE INDEX IF NOT EXISTS idx_billing_plans_pricing_model
  ON billing.plans (pricing_model);
