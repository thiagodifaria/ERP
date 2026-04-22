# finance

This domain now owns both the projection layer and the first operational finance cycle of the ERP.

Current structure:

- migrations
- indexes
- operational snapshots

Current scope:

- `finance.receivable_projections`
- idempotent linkage to `sales.outbox_events`
- `finance.receivable_entries`
- `finance.receivable_settlements`
- `finance.commission_entries`
- `finance.payables`
- `finance.cost_entries`
- `finance.period_closures`
- shared `updated_at` triggers for the operational cycle
