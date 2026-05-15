DROP TRIGGER IF EXISTS trg_banking_records_updated_at ON banking.records;
CREATE TRIGGER trg_banking_records_updated_at
BEFORE UPDATE ON banking.records
FOR EACH ROW EXECUTE FUNCTION common.set_updated_at();
