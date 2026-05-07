DROP TRIGGER IF EXISTS trg_fiscal_company_profiles_updated_at ON fiscal.company_profiles;
CREATE TRIGGER trg_fiscal_company_profiles_updated_at
BEFORE UPDATE ON fiscal.company_profiles
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_fiscal_retention_policies_updated_at ON fiscal.retention_policies;
CREATE TRIGGER trg_fiscal_retention_policies_updated_at
BEFORE UPDATE ON fiscal.retention_policies
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_fiscal_documents_updated_at ON fiscal.documents;
CREATE TRIGGER trg_fiscal_documents_updated_at
BEFORE UPDATE ON fiscal.documents
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_fiscal_privacy_requests_updated_at ON fiscal.privacy_requests;
CREATE TRIGGER trg_fiscal_privacy_requests_updated_at
BEFORE UPDATE ON fiscal.privacy_requests
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
