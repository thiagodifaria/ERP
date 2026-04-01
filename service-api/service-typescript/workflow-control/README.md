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
- room for workflow definitions and future control-plane APIs

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/workflow-control/definitions`
- `POST /api/workflow-control/runs`
- `GET /api/workflow-control/runs`
- `GET /api/workflow-control/runs/summary`
- `GET /api/workflow-control/runs/{publicId}`
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
- workflow definition metadata update
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

Current contract scope:

- public payload shape for workflow definition list
- public payload shape for workflow run list
- public payload shape for workflow run detail
- workflow run create payload and linkage to current version
- workflow run operational summary payload
- create/update/detail/status lifecycle
- version list, publish and current-version lifecycle
- version detail and restore lifecycle
- not found error contract
