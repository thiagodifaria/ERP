DROP TRIGGER IF EXISTS trg_billing_plans_updated_at ON billing.plans;
CREATE TRIGGER trg_billing_plans_updated_at
BEFORE UPDATE ON billing.plans
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_billing_subscriptions_updated_at ON billing.subscriptions;
CREATE TRIGGER trg_billing_subscriptions_updated_at
BEFORE UPDATE ON billing.subscriptions
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_billing_subscription_invoices_updated_at ON billing.subscription_invoices;
CREATE TRIGGER trg_billing_subscription_invoices_updated_at
BEFORE UPDATE ON billing.subscription_invoices
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
