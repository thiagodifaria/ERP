# billing

The billing service owns ERP subscription charging, invoice collection lifecycle and gateway webhook reconciliation.

Current scope:

- plans
- subscriptions
- subscription invoices
- payment attempts with idempotency keys
- retry counting and grace period transitions
- controlled suspension and reactivation
- webhook processing through `webhook-hub`
- billing events and operational reporting

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/billing/plans`
- `POST /api/billing/plans`
- `GET /api/billing/subscriptions`
- `POST /api/billing/subscriptions`
- `GET /api/billing/subscriptions/{publicId}`
- `GET /api/billing/subscriptions/{publicId}/events`
- `POST /api/billing/subscriptions/{publicId}/suspend`
- `POST /api/billing/subscriptions/{publicId}/reactivate`
- `GET /api/billing/invoices`
- `POST /api/billing/subscriptions/{publicId}/invoices`
- `GET /api/billing/invoices/{publicId}/attempts`
- `POST /api/billing/invoices/{publicId}/attempts`
- `POST /api/billing/webhook-events/process`
- `GET /api/billing/reports/operations`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/billing mcr.microsoft.com/dotnet/sdk:8.0 dotnet build src/Billing.Api/Billing.Api.csproj -c Release`
- `docker build -t erp-billing ./service-api/service-csharp/billing`
- `bash scripts/test.sh smoke`
