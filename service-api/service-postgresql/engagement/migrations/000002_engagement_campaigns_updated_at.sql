-- Mantem updated_at coerente nas campanhas do contexto engagement.

DROP TRIGGER IF EXISTS trg_engagement_campaigns_updated_at ON engagement.campaigns;

CREATE TRIGGER trg_engagement_campaigns_updated_at
BEFORE UPDATE ON engagement.campaigns
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
