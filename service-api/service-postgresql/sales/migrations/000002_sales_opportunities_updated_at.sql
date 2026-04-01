-- Mantem o campo updated_at coerente nas oportunidades do contexto sales.

DROP TRIGGER IF EXISTS trg_sales_opportunities_updated_at ON sales.opportunities;

CREATE TRIGGER trg_sales_opportunities_updated_at
BEFORE UPDATE ON sales.opportunities
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
