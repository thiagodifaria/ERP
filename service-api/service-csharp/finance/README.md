# finance

The finance service now owns the first operational finance cycle of the ERP, bridging commercial activity and downstream financial control.

Current scope:

- health and readiness endpoints
- receivable projection ingestion from `sales.outbox_events`
- idempotent projection persistence in PostgreSQL
- renegotiation and invoice status reflected into projected receivable amounts
- operational receivables synchronized from `sales.invoices`
- settlement flow with idempotent `settlementReference`
- operational commissions synchronized from `sales.commissions`
- payables with creation, payment and cancellation lifecycle
- operational cost entries
- monthly period closures with persisted financial snapshots
- operational report and database summaries for receivables, commissions, payables, costs and closures

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `POST /api/finance/projections/ingest`
- `GET /api/finance/projections`
- `GET /api/finance/projections/summary`
- `POST /api/finance/operations/sync`
- `GET /api/finance/receivables`
- `POST /api/finance/receivables/{publicId}/settlements`
- `GET /api/finance/commissions`
- `GET /api/finance/commissions/summary`
- `GET /api/finance/payables`
- `POST /api/finance/payables`
- `PATCH /api/finance/payables/{publicId}/status`
- `GET /api/finance/costs`
- `POST /api/finance/costs`
- `POST /api/finance/period-closures`
- `GET /api/finance/period-closures`
- `GET /api/finance/period-closures/{periodKey}`
- `GET /api/finance/reports/operations`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/finance mcr.microsoft.com/dotnet/sdk:8.0 dotnet build src/Finance.Api/Finance.Api.csproj -c Release`
- `docker build -t erp-finance ./service-api/service-csharp/finance`
- `bash scripts/test.sh smoke`
