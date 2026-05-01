DROP TRIGGER IF EXISTS trg_billing_recovery_cases_updated_at ON billing.recovery_cases;
CREATE TRIGGER trg_billing_recovery_cases_updated_at
BEFORE UPDATE ON billing.recovery_cases
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
