# sales

The sales service owns opportunities, proposals, closing flows, invoicing and early revenue operations.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- opportunity pipeline with controlled stage transitions
- explicit customer linkage from CRM and typed sales such as `new`, `upsell`, `renewal` and `expansion`
- proposal catalog linked to opportunities
- sale conversion flow linked to accepted proposals
- renegotiation flow with commercial amount history
- installment scheduling for recurring or split commercial agreements
- operational commissions with status transitions such as `pending`, `blocked` and `released`
- operational pending items with open, resolved and cancelled states
- invoice creation linked to closed sales
- collection summary with open, paid and overdue footprints
- commercial history ledger by aggregate
- transactional outbox for downstream consumers
- selectable repository driver with PostgreSQL-backed persistence for a bootstrap tenant
- contract and smoke coverage for opportunities, proposals, sales, invoices and the new operational surfaces

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/sales/opportunities`
- `GET /api/sales/opportunities/summary`
- `POST /api/sales/opportunities`
- `GET /api/sales/opportunities/{publicId}`
- `GET /api/sales/opportunities/{publicId}/history`
- `PATCH /api/sales/opportunities/{publicId}`
- `PATCH /api/sales/opportunities/{publicId}/stage`
- `GET /api/sales/opportunities/{publicId}/proposals`
- `POST /api/sales/opportunities/{publicId}/proposals`
- `GET /api/sales/proposals/{publicId}`
- `GET /api/sales/proposals/{publicId}/history`
- `PATCH /api/sales/proposals/{publicId}/status`
- `POST /api/sales/proposals/{publicId}/convert`
- `GET /api/sales/sales`
- `GET /api/sales/sales/summary`
- `GET /api/sales/sales/{publicId}`
- `GET /api/sales/sales/{publicId}/history`
- `GET /api/sales/sales/{publicId}/installments`
- `POST /api/sales/sales/{publicId}/installments`
- `GET /api/sales/sales/{publicId}/commissions`
- `POST /api/sales/sales/{publicId}/commissions`
- `PATCH /api/sales/commissions/{publicId}/status`
- `GET /api/sales/sales/{publicId}/pending-items`
- `POST /api/sales/sales/{publicId}/pending-items`
- `PATCH /api/sales/pending-items/{publicId}/status`
- `GET /api/sales/sales/{publicId}/renegotiations`
- `POST /api/sales/sales/{publicId}/renegotiations`
- `PATCH /api/sales/sales/{publicId}/status`
- `POST /api/sales/sales/{publicId}/invoice`
- `GET /api/sales/invoices`
- `GET /api/sales/invoices/summary`
- `GET /api/sales/invoices/{publicId}`
- `GET /api/sales/invoices/{publicId}/history`
- `PATCH /api/sales/invoices/{publicId}/status`
- `GET /api/sales/outbox/pending`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-golang/sales golang:1.24-alpine go test ./...`
- `docker build -t erp-sales ./service-api/service-golang/sales`
- `bash scripts/test.sh smoke`

Runtime switch:

- `SALES_REPOSITORY_DRIVER=postgres`
