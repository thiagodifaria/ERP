DROP TRIGGER IF EXISTS trg_documents_attachments_updated_at ON documents.attachments;

CREATE TRIGGER trg_documents_attachments_updated_at
BEFORE UPDATE ON documents.attachments
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
