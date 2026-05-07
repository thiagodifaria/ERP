DROP TRIGGER IF EXISTS trg_supplier_categories_updated_at ON supplier.categories;
CREATE TRIGGER trg_supplier_categories_updated_at
BEFORE UPDATE ON supplier.categories
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_supplier_suppliers_updated_at ON supplier.suppliers;
CREATE TRIGGER trg_supplier_suppliers_updated_at
BEFORE UPDATE ON supplier.suppliers
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
