# workflow-runtime

This domain owns durable execution state, retries, timers and workflow runtime traces.

Current relational scope:

- `workflow_runtime.executions`
- `workflow_runtime.execution_transitions`
- lifecycle timestamps and retry count
- tenant-scoped indexes for execution reads
- audit trail for status changes
