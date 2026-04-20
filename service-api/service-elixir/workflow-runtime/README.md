# workflow-runtime

The workflow-runtime service owns durable execution, timers and high-resilience automation flow.

Initial scope:

- Elixir OTP service bootstrap
- Plug HTTP server
- in-memory execution store
- health and readiness routes
- execution list route
- execution summary route
- execution transition ledger
- execution action ledger
- filtered operational reads by tenant and status
- runtime catalog validation against published workflow definitions
- runtime capability snapshot for timers, retries and compensations
- published action-plan execution with delay handling
- compensation append on failed or cancelled flows
- room for timers, retries and durable orchestration
- selectable repository driver between memory and PostgreSQL

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/workflow-runtime/capabilities`
- `GET /api/workflow-runtime/executions`
- `GET /api/workflow-runtime/executions/{publicId}`
- `GET /api/workflow-runtime/executions/{publicId}/actions`
- `GET /api/workflow-runtime/executions/{publicId}/timeline`
- `GET /api/workflow-runtime/executions/{publicId}/transitions`
- `GET /api/workflow-runtime/executions/summary`
- `GET /api/workflow-runtime/executions/summary?tenantSlug=...&workflowDefinitionKey=...`
- `GET /api/workflow-runtime/executions/summary/by-workflow`
- `POST /api/workflow-runtime/executions`
- `POST /api/workflow-runtime/executions/{publicId}/start`
- `POST /api/workflow-runtime/executions/{publicId}/advance`
- `POST /api/workflow-runtime/executions/{publicId}/complete`
- `POST /api/workflow-runtime/executions/{publicId}/fail`
- `POST /api/workflow-runtime/executions/{publicId}/cancel`
- `POST /api/workflow-runtime/executions/{publicId}/retry`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-elixir/workflow-runtime elixir:1.17-alpine sh -lc "apk add --no-cache build-base git && mix local.hex --force && mix local.rebar --force && mix deps.get && mix test"`
- `docker build -t erp-workflow-runtime ./service-api/service-elixir/workflow-runtime`

Runtime switch:

- `WORKFLOW_RUNTIME_REPOSITORY_DRIVER=postgres`

Execution model:

- each execution stores the published workflow definition version consumed at creation time
- runtime actions are replayed from the published `snapshotActions` plan persisted by `workflow-control`
- `advance` executes the next planned action, records the action ledger and completes the execution when the plan finishes
- `delay.wait` actions pause the execution with `waitingUntil` until the configured delay elapses
- failed or cancelled executions append compensation entries for completed actions that declare `compensationActionKey`
