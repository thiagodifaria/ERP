# billing

This context owns ERP subscription billing persistence, covering plans, subscriptions,
collection invoices, payment attempts, billing event audit and webhook reconciliation state.

Current relational scope:

- `billing.plans`
- `billing.subscriptions`
- `billing.subscription_invoices`
- `billing.payment_attempts`
- `billing.subscription_events`

Operational notes:

- `plans` keep the catalog of billable offerings and grace-period defaults.
- `subscriptions` persist tenant-facing billing contracts, lifecycle status and gateway references.
- `subscription_invoices` track invoice amount, due date and collection state.
- `payment_attempts` enforce idempotent retry registration with provider metadata.
- `subscription_events` preserve the billing audit trail for creation, grace period, suspension, payment and reactivation.

Validation:

- `bash scripts/db.sh migrate billing`
- `bash scripts/db.sh summary billing`
- `bash scripts/test.sh smoke`
