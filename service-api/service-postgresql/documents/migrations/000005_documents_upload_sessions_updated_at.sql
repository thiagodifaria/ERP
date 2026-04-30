DROP TRIGGER IF EXISTS trg_documents_upload_sessions_updated_at ON documents.upload_sessions;
CREATE TRIGGER trg_documents_upload_sessions_updated_at
BEFORE UPDATE ON documents.upload_sessions
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
