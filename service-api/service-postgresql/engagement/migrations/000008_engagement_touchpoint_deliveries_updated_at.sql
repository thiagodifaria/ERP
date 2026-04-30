DROP TRIGGER IF EXISTS trg_engagement_touchpoint_deliveries_updated_at
  ON engagement.touchpoint_deliveries;

CREATE TRIGGER trg_engagement_touchpoint_deliveries_updated_at
BEFORE UPDATE ON engagement.touchpoint_deliveries
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
