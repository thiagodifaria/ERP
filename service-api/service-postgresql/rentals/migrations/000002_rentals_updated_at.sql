DROP TRIGGER IF EXISTS trg_rentals_contracts_updated_at ON rentals.contracts;

CREATE TRIGGER trg_rentals_contracts_updated_at
BEFORE UPDATE ON rentals.contracts
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_rentals_contract_charges_updated_at ON rentals.contract_charges;

CREATE TRIGGER trg_rentals_contract_charges_updated_at
BEFORE UPDATE ON rentals.contract_charges
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_rentals_outbox_events_updated_at ON rentals.outbox_events;

CREATE TRIGGER trg_rentals_outbox_events_updated_at
BEFORE UPDATE ON rentals.outbox_events
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
