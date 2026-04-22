-- Mantem updated_at coerente nas tabelas operacionais do financeiro.
DROP TRIGGER IF EXISTS trg_finance_receivable_entries_updated_at ON finance.receivable_entries;
CREATE TRIGGER trg_finance_receivable_entries_updated_at
BEFORE UPDATE ON finance.receivable_entries
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_finance_commission_entries_updated_at ON finance.commission_entries;
CREATE TRIGGER trg_finance_commission_entries_updated_at
BEFORE UPDATE ON finance.commission_entries
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_finance_payables_updated_at ON finance.payables;
CREATE TRIGGER trg_finance_payables_updated_at
BEFORE UPDATE ON finance.payables
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_finance_cost_entries_updated_at ON finance.cost_entries;
CREATE TRIGGER trg_finance_cost_entries_updated_at
BEFORE UPDATE ON finance.cost_entries
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
