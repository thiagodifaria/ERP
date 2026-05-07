DROP TRIGGER IF EXISTS trg_support_queues_updated_at ON support.queues;
CREATE TRIGGER trg_support_queues_updated_at
BEFORE UPDATE ON support.queues
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_support_cases_updated_at ON support.cases;
CREATE TRIGGER trg_support_cases_updated_at
BEFORE UPDATE ON support.cases
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
