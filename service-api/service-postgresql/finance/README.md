# finance

This domain owns the first financial projections derived from the commercial stack and will later absorb receivables, payables, costs, commissions and closures.

Current structure:

- migrations
- indexes

Current scope:

- `finance.receivable_projections`
- idempotent linkage to `sales.outbox_events`
- operational summary support for the initial finance API
