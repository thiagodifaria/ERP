# workflow-control

The workflow-control service owns workflow definitions, activation rules and control-plane orchestration concerns.

Initial scope:

- TypeScript service bootstrap
- layered folder structure aligned with the monorepo standard
- health and readiness routes
- room for workflow definitions and future control-plane APIs

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-typescript/workflow-control node:22-alpine sh -lc "npm install && npm run test:unit"`
