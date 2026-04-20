-- Aplica trigger padrao de updated_at para clientes do contexto CRM.
DROP TRIGGER IF EXISTS set_crm_customers_updated_at ON crm.customers;

CREATE TRIGGER set_crm_customers_updated_at
BEFORE UPDATE ON crm.customers
FOR EACH ROW
EXECUTE FUNCTION common.common_set_updated_at();
