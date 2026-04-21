# simulation

The simulation service owns operational what-if scenarios, sizing inputs and future workload modeling.

Initial scope:

- FastAPI service bootstrap
- health and readiness routes
- scenario catalog for workload modeling
- selectable repository driver between memory and PostgreSQL
- persistence foundation for future scenario runs and load benchmarks

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/simulation/scenarios/catalog`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-python/simulation python:3.12-slim sh -lc "pip install -e .[dev] && pytest"`
- `docker build -t erp-simulation ./service-api/service-python/simulation`

Runtime switch:

- `SIMULATION_REPOSITORY_DRIVER=postgres`
