-- Aplica trigger padrao de updated_at para leads do contexto CRM.
-- A funcao compartilhada fica no contexto common.

DROP TRIGGER IF EXISTS set_crm_leads_updated_at ON crm.leads;

CREATE TRIGGER set_crm_leads_updated_at
BEFORE UPDATE ON crm.leads
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
