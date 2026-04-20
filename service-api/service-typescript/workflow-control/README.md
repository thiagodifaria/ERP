# workflow-control

The workflow-control service owns workflow definitions, activation rules and control-plane orchestration concerns.

Initial scope:

- TypeScript service bootstrap
- layered folder structure aligned with the monorepo standard
- runtime config isolated in `src/config`
- health and readiness routes
- bootstrap workflow definition list
- repository abstraction ready for memory and PostgreSQL
- domain primitives ready for versioned publication history
- domain primitives ready for workflow execution tracking
- version persistence ready for memory and PostgreSQL
- workflow run repository abstraction ready for memory and PostgreSQL
- workflow run PostgreSQL adapter ready for runtime read and write
- runtime container wired for workflow run persistence and readiness
- workflow run event repositories ready for memory and PostgreSQL
- workflow trigger catalog bootstrap with public read-side
- workflow action catalog bootstrap with public read-side
- workflow definition trigger validation against the published catalog
- workflow action-plan persistence on current definitions and published versions
- room for workflow definitions and future control-plane APIs

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/workflow-control/catalog/triggers`
- `GET /api/workflow-control/catalog/actions`
- `GET /api/workflow-control/editor`
- `GET /api/workflow-control/definitions`
- `POST /api/workflow-control/runs`
- `GET /api/workflow-control/runs`
- `GET /api/workflow-control/runs/summary`
- `GET /api/workflow-control/runs/{publicId}`
- `GET /api/workflow-control/runs/{publicId}/events`
- `GET /api/workflow-control/runs/{publicId}/events/summary`
- `POST /api/workflow-control/runs/{publicId}/events`

Event query params:

- `category=status|note`
- `createdBy=<actor>`
- `POST /api/workflow-control/runs/{publicId}/start`
- `POST /api/workflow-control/runs/{publicId}/complete`
- `POST /api/workflow-control/runs/{publicId}/fail`
- `POST /api/workflow-control/runs/{publicId}/cancel`
- `GET /api/workflow-control/definitions/{key}`
- `GET /api/workflow-control/definitions/{key}/versions`
- `GET /api/workflow-control/definitions/{key}/versions/current`
- `GET /api/workflow-control/definitions/{key}/versions/{versionNumber}`
- `GET /api/workflow-control/definitions/{key}/versions/summary`
- `POST /api/workflow-control/definitions/{key}/versions`
- `POST /api/workflow-control/definitions/{key}/versions/{versionNumber}/restore`
- `POST /api/workflow-control/definitions`
- `PATCH /api/workflow-control/definitions/{key}`
- `PATCH /api/workflow-control/definitions/{key}/status`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-typescript/workflow-control node:22-alpine sh -lc "npm install && npm run test:unit"`
- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-typescript/workflow-control node:22-alpine sh -lc "npm install && npm run test:contract"`
- `docker build -t erp-workflow-control ./service-api/service-typescript/workflow-control`
- `WORKFLOW_CONTROL_REPOSITORY_DRIVER=postgres` habilita o catalogo relacional em cima do contexto `workflow_control`

Current unit scope:

- health and readiness routes
- workflow definition list, create, detail and status transitions
- workflow run list read
- workflow run detail read
- workflow run create linked to current published version
- workflow run summary read
- workflow run event history read
- workflow run note create
- workflow run automatic status event ledger
- workflow run event summary
- workflow run event filters by category and creator
- workflow run start transition
- workflow run complete transition
- workflow run fail transition
- workflow run cancel transition
- workflow run filtered list read
- workflow definition metadata update
- workflow trigger and action catalog read
- workflow definition trigger validation against catalog dependencies
- workflow definition action-plan validation and normalization
- workflow definition version history read
- workflow definition manual publication
- workflow definition current-version read
- workflow definition version detail read
- workflow definition version summary read
- workflow definition restore from published snapshot
- workflow run aggregate validation helpers
- workflow run repository bootstrap in memory
- workflow run PostgreSQL persistence build
- workflow run readiness dependency on runtime boot
- workflow run event repositories for memory and PostgreSQL

Current contract scope:

- public payload shape for workflow definition list
- public payload shape for workflow run list
- public payload shape for workflow run detail
- public payload shape for workflow run event history
- public payload shape for workflow run note create
- public payload shape for workflow run status event ledger
- public payload shape for workflow run event summary
- public payload shape for workflow run event filters
- workflow run create payload and linkage to current version
- workflow run operational summary payload
- workflow run pending-to-running transition
- workflow run running-to-completed transition
- workflow run running-to-failed transition
- workflow run pending-or-running-to-cancelled transition
- workflow run filters by status, definition key, subject type and initiator
- contract coverage for workflow run list, detail and create
- contract coverage for workflow run summary and filters
- contract coverage for workflow run lifecycle transitions
- contract coverage for workflow run event history
- create/update/detail/status lifecycle
- trigger and action catalog payloads
- version list, publish and current-version lifecycle
- version detail and restore lifecycle
- action-plan payload persistence across create, update, publish and restore
- not found error contract
