DROP TRIGGER IF EXISTS trg_procurement_records_updated_at ON procurement.records;
CREATE TRIGGER trg_procurement_records_updated_at
BEFORE UPDATE ON procurement.records
FOR EACH ROW EXECUTE FUNCTION common.set_updated_at();
