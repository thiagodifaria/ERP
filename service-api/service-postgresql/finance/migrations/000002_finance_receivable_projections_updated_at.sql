-- Mantem updated_at coerente nas projecoes financeiras.
DROP TRIGGER IF EXISTS trg_finance_receivable_projections_updated_at ON finance.receivable_projections;

CREATE TRIGGER trg_finance_receivable_projections_updated_at
BEFORE UPDATE ON finance.receivable_projections
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
