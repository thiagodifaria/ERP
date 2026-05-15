DROP TRIGGER IF EXISTS trg_accounting_records_updated_at ON accounting.records;
CREATE TRIGGER trg_accounting_records_updated_at
BEFORE UPDATE ON accounting.records
FOR EACH ROW EXECUTE FUNCTION common.set_updated_at();
