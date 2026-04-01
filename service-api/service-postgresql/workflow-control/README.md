# workflow-control

This domain owns workflow definitions, triggers, actions and versioned published plans.

Current structure:

- migrations
- seeds
- workflow definition catalog
- workflow definition version catalog
- workflow run ledger
- workflow run event log
- updated_at trigger for mutable records
- subject and trigger indexes for workflow run lookups
- bootstrap seed for lead follow-up orchestration
- bootstrap seed for version 1 of each base workflow
- bootstrap seed for first workflow run per tenant
- views
- functions
- indexes
