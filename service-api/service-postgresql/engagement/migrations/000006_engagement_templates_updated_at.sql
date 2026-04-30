DROP TRIGGER IF EXISTS trg_engagement_templates_updated_at
  ON engagement.templates;

CREATE TRIGGER trg_engagement_templates_updated_at
BEFORE UPDATE ON engagement.templates
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
