# simulation

The simulation service owns operational what-if scenarios, sizing inputs and future workload modeling.

Initial scope:

- FastAPI service bootstrap
- health and readiness routes
- scenario catalog for workload modeling
- operational-load scenario runs with PostgreSQL persistence
- load benchmark runs with projected latency, throughput and capacity status
- selectable repository driver between memory and PostgreSQL

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/simulation/scenarios/catalog`
- `POST /api/simulation/scenarios/operational-load`
- `GET /api/simulation/scenarios/runs`
- `GET /api/simulation/scenarios/runs/{publicId}`
- `POST /api/simulation/benchmarks/load`
- `GET /api/simulation/benchmarks/runs`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-python/simulation python:3.12-slim sh -lc "pip install -e .[dev] && pytest"`
- `docker build -t erp-simulation ./service-api/service-python/simulation`

Runtime switch:

- `SIMULATION_REPOSITORY_DRIVER=postgres`
