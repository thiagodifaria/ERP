-- Mantem updated_at coerente nas tabelas de tesouraria.
DROP TRIGGER IF EXISTS trg_finance_cash_accounts_updated_at ON finance.cash_accounts;
CREATE TRIGGER trg_finance_cash_accounts_updated_at
BEFORE UPDATE ON finance.cash_accounts
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_finance_cash_movements_updated_at ON finance.cash_movements;
CREATE TRIGGER trg_finance_cash_movements_updated_at
BEFORE UPDATE ON finance.cash_movements
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
