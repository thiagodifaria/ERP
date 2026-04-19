-- Mantem o campo updated_at coerente nas invoices do contexto sales.

DROP TRIGGER IF EXISTS trg_sales_invoices_updated_at ON sales.invoices;

CREATE TRIGGER trg_sales_invoices_updated_at
BEFORE UPDATE ON sales.invoices
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
