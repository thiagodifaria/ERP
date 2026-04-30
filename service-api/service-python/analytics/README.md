# analytics

The analytics service owns heavy operational reads, summarized reports and future forecasting workloads.

Initial scope:

- FastAPI service bootstrap
- health and readiness routes
- first operational pipeline summary report
- automation board with workflow-level operational breakdown
- workflow definition health report with catalog and runtime alignment
- revenue operations report crossing booked revenue, invoice coverage and collection risk
- rental operations report crossing contracts, recurring charges and attachment governance
- load benchmark report backed by simulation runs
- scenario-driven cost estimator
- selectable repository driver between memory and PostgreSQL
- room for warehouse ingestion, ETL and forecasting models

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/analytics/reports/pipeline-summary`
- `GET /api/analytics/reports/service-pulse`
- `GET /api/analytics/reports/sales-journey`
- `GET /api/analytics/reports/tenant-360`
- `GET /api/analytics/reports/automation-board`
- `GET /api/analytics/reports/workflow-definition-health`
- `GET /api/analytics/reports/delivery-reliability`
- `GET /api/analytics/reports/revenue-operations`
- `GET /api/analytics/reports/rental-operations`
- `GET /api/analytics/reports/load-benchmark`
- `GET /api/analytics/reports/cost-estimator`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-python/analytics python:3.12-slim sh -lc "pip install -e .[dev] && pytest"`
- `docker build -t erp-analytics ./service-api/service-python/analytics`

Runtime switch:

- `ANALYTICS_REPOSITORY_DRIVER=postgres`
