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
- version persistence ready for memory and PostgreSQL
- room for workflow definitions and future control-plane APIs

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/workflow-control/definitions`
- `GET /api/workflow-control/definitions/{key}`
- `GET /api/workflow-control/definitions/{key}/versions`
- `GET /api/workflow-control/definitions/{key}/versions/current`
- `POST /api/workflow-control/definitions/{key}/versions`
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
- workflow definition metadata update
- workflow definition version history read
- workflow definition manual publication
- workflow definition current-version read

Current contract scope:

- public payload shape for workflow definition list
- create/update/detail/status lifecycle
- version list, publish and current-version lifecycle
- not found error contract
