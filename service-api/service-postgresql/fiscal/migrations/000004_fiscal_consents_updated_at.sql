DROP TRIGGER IF EXISTS trg_fiscal_consents_updated_at ON fiscal.consents;
CREATE TRIGGER trg_fiscal_consents_updated_at
BEFORE UPDATE ON fiscal.consents
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
