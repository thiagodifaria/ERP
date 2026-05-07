DROP TRIGGER IF EXISTS trg_notification_preferences_updated_at ON notification.preferences;
CREATE TRIGGER trg_notification_preferences_updated_at
BEFORE UPDATE ON notification.preferences
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_notification_notifications_updated_at ON notification.notifications;
CREATE TRIGGER trg_notification_notifications_updated_at
BEFORE UPDATE ON notification.notifications
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
