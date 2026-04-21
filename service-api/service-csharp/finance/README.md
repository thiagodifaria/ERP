# finance

The finance service starts as the first downstream consumer of commercial events.

Initial scope:

- health and readiness endpoints
- receivable projection ingestion from `sales.outbox_events`
- idempotent projection persistence in PostgreSQL
- commercial renegotiation reflected into projected receivable amounts
- sale-booking and invoice projections with `forecast`, `open`, `paid` and `cancelled` states
- projection list and summary routes for operational follow-up
- room for future receivables, commissions, closures and financial snapshots

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `POST /api/finance/projections/ingest`
- `GET /api/finance/projections`
- `GET /api/finance/projections/summary`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/finance mcr.microsoft.com/dotnet/sdk:8.0 dotnet build src/Finance.Api/Finance.Api.csproj -c Release`
- `docker build -t erp-finance ./service-api/service-csharp/finance`
