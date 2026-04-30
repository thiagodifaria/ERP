ALTER TABLE rentals.contract_charges
  ADD COLUMN IF NOT EXISTS paid_at TIMESTAMPTZ NULL,
  ADD COLUMN IF NOT EXISTS payment_reference VARCHAR(160) NOT NULL DEFAULT '';

ALTER TABLE rentals.contract_charges
  DROP CONSTRAINT IF EXISTS ck_rentals_contract_charges_status;

ALTER TABLE rentals.contract_charges
  ADD CONSTRAINT ck_rentals_contract_charges_status
  CHECK (status IN ('scheduled', 'paid', 'cancelled'));

ALTER TABLE rentals.contract_charges
  DROP CONSTRAINT IF EXISTS ck_rentals_contract_charges_paid_state;

ALTER TABLE rentals.contract_charges
  ADD CONSTRAINT ck_rentals_contract_charges_paid_state
  CHECK (
    (status = 'paid' AND paid_at IS NOT NULL AND length(trim(payment_reference)) > 0)
    OR (status <> 'paid' AND paid_at IS NULL)
  );
