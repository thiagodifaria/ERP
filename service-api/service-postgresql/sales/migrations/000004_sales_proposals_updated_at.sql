-- Mantem o campo updated_at coerente nas propostas do contexto sales.

DROP TRIGGER IF EXISTS trg_sales_proposals_updated_at ON sales.proposals;

CREATE TRIGGER trg_sales_proposals_updated_at
BEFORE UPDATE ON sales.proposals
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
