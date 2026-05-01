# finance

The finance service now owns the first operational finance cycle of the ERP, bridging commercial activity and downstream financial control.

Current scope:

- health and readiness endpoints
- receivable projection ingestion from `sales.outbox_events`
- idempotent projection persistence in PostgreSQL
- renegotiation and invoice status reflected into projected receivable amounts
- operational receivables synchronized from `sales.invoices`
- operational receivables synchronized from `rentals.charges`
- settlement flow with idempotent `settlementReference`
- operational commissions synchronized from `sales.commissions`
- commission block and release governance
- payables with creation, payment and cancellation lifecycle
- operational cost entries
- treasury cash accounts with opening balance and provider metadata
- idempotent treasury sync from receivable settlements, paid payables and operational costs
- cash movement summary with inflow, outflow and live balance
- treasury liquidity report with current balance and projected net position
- monthly period closures with persisted financial snapshots
- operational activity ledger for receivable sync, settlements, payables, treasury and closures
- operational report and database summaries for receivables, commissions, payables, costs, treasury and closures

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
- `POST /api/finance/commissions/{publicId}/block`
- `POST /api/finance/commissions/{publicId}/release`
- `GET /api/finance/payables`
- `POST /api/finance/payables`
- `PATCH /api/finance/payables/{publicId}/status`
- `GET /api/finance/costs`
- `POST /api/finance/costs`
- `GET /api/finance/cash-accounts`
- `POST /api/finance/cash-accounts`
- `POST /api/finance/treasury/sync`
- `GET /api/finance/cash-movements`
- `GET /api/finance/cash-movements/summary`
- `GET /api/finance/reports/treasury`
- `POST /api/finance/period-closures`
- `GET /api/finance/period-closures`
- `GET /api/finance/period-closures/{periodKey}`
- `GET /api/finance/reports/operations`
- `GET /api/finance/activity`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/finance mcr.microsoft.com/dotnet/sdk:8.0 dotnet build src/Finance.Api/Finance.Api.csproj -c Release`
- `docker build -t erp-finance ./service-api/service-csharp/finance`
- `bash scripts/test.sh smoke`
