# workflow-runtime

The workflow-runtime service owns durable execution, timers and high-resilience automation flow.

Initial scope:

- Elixir OTP service bootstrap
- Plug HTTP server
- in-memory execution store
- health and readiness routes
- execution list route
- execution summary route
- room for timers, retries and durable orchestration

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/workflow-runtime/executions`
- `GET /api/workflow-runtime/executions/{publicId}`
- `GET /api/workflow-runtime/executions/summary`
- `POST /api/workflow-runtime/executions`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-elixir/workflow-runtime elixir:1.17-alpine sh -lc "apk add --no-cache build-base git && mix local.hex --force && mix local.rebar --force && mix deps.get && mix test"`
- `docker build -t erp-workflow-runtime ./service-api/service-elixir/workflow-runtime`
