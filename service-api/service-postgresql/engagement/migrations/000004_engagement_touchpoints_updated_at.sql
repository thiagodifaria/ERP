-- Mantem updated_at coerente nos touchpoints do contexto engagement.

DROP TRIGGER IF EXISTS trg_engagement_touchpoints_updated_at ON engagement.touchpoints;

CREATE TRIGGER trg_engagement_touchpoints_updated_at
BEFORE UPDATE ON engagement.touchpoints
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
