# billing

The billing service owns ERP subscription charging, invoice collection lifecycle and gateway webhook reconciliation.

Current scope:

- plans
- subscriptions
- subscription invoices
- payment attempts with idempotency keys
- operational batch processing of pending webhook events
- retry counting and grace period transitions
- recovery cases with severity, next action and payment promise governance
- recovery action ledger for contact, promise and resolution lifecycle
- controlled suspension and reactivation
- webhook processing through `webhook-hub`
- billing events and operational reporting
- gateway capability registry with explicit `configured`, `manual` and `unconfigured` posture for Pix-ready adapters

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/billing/gateways`
- `GET /api/billing/gateways/{provider}`
- `GET /api/billing/plans`
- `POST /api/billing/plans`
- `GET /api/billing/subscriptions`
- `POST /api/billing/subscriptions`
- `GET /api/billing/subscriptions/{publicId}`
- `GET /api/billing/subscriptions/{publicId}/events`
- `POST /api/billing/subscriptions/{publicId}/suspend`
- `POST /api/billing/subscriptions/{publicId}/reactivate`
- `GET /api/billing/invoices`
- `GET /api/billing/invoices/{publicId}`
- `POST /api/billing/subscriptions/{publicId}/invoices`
- `GET /api/billing/invoices/{publicId}/attempts`
- `POST /api/billing/invoices/{publicId}/attempts`
- `POST /api/billing/webhook-events/process`
- `GET /api/billing/webhook-events/pending`
- `POST /api/billing/webhook-events/process-batch`
- `GET /api/billing/events`
- `GET /api/billing/reports/operations`
- `GET /api/billing/recovery/cases`
- `GET /api/billing/recovery/cases/{publicId}`
- `GET /api/billing/recovery/cases/{publicId}/actions`
- `POST /api/billing/invoices/{publicId}/recovery/open`
- `POST /api/billing/recovery/cases/{publicId}/touchpoints`
- `POST /api/billing/recovery/cases/{publicId}/promise`
- `POST /api/billing/recovery/cases/{publicId}/close`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/billing mcr.microsoft.com/dotnet/sdk:8.0 dotnet build src/Billing.Api/Billing.Api.csproj -c Release`
- `docker build -t erp-billing ./service-api/service-csharp/billing`
- `bash scripts/test.sh smoke`
