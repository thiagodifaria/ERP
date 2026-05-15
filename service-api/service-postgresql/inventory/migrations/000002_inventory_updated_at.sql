DROP TRIGGER IF EXISTS trg_inventory_records_updated_at ON inventory.records;
CREATE TRIGGER trg_inventory_records_updated_at
BEFORE UPDATE ON inventory.records
FOR EACH ROW EXECUTE FUNCTION common.set_updated_at();
