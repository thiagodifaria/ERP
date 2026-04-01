# analytics

The analytics service owns heavy operational reads, summarized reports and future forecasting workloads.

Initial scope:

- FastAPI service bootstrap
- health and readiness routes
- first operational pipeline summary report
- room for warehouse ingestion, ETL and forecasting models

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/analytics/reports/pipeline-summary`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-python/analytics python:3.12-slim sh -lc "pip install -e .[dev] && pytest"`
- `docker build -t erp-analytics ./service-api/service-python/analytics`
