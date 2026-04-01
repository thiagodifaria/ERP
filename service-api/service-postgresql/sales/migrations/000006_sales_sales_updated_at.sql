-- Mantem o campo updated_at coerente nas vendas do contexto sales.

DROP TRIGGER IF EXISTS trg_sales_sales_updated_at ON sales.sales;

CREATE TRIGGER trg_sales_sales_updated_at
BEFORE UPDATE ON sales.sales
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
