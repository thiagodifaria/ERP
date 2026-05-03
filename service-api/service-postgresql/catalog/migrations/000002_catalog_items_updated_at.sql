DROP TRIGGER IF EXISTS trg_catalog_items_updated_at ON catalog.items;
CREATE TRIGGER trg_catalog_items_updated_at
BEFORE UPDATE ON catalog.items
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
