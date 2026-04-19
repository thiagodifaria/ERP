# sales

The sales service owns opportunities, proposals, closing flows, invoicing and early revenue operations.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- opportunity pipeline with controlled stage transitions
- proposal catalog linked to opportunities
- sale conversion flow linked to accepted proposals
- invoice creation linked to closed sales
- collection summary with open, paid and overdue footprints
- selectable repository driver with PostgreSQL-backed persistence for a bootstrap tenant
- contract and smoke coverage for opportunities, proposals, sales and invoices
- room for commissions and deeper revenue operations

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/sales/opportunities`
- `GET /api/sales/opportunities/summary`
- `POST /api/sales/opportunities`
- `GET /api/sales/opportunities/{publicId}`
- `PATCH /api/sales/opportunities/{publicId}`
- `PATCH /api/sales/opportunities/{publicId}/stage`
- `GET /api/sales/opportunities/{publicId}/proposals`
- `POST /api/sales/opportunities/{publicId}/proposals`
- `GET /api/sales/proposals/{publicId}`
- `PATCH /api/sales/proposals/{publicId}/status`
- `POST /api/sales/proposals/{publicId}/convert`
- `GET /api/sales/sales`
- `GET /api/sales/sales/summary`
- `GET /api/sales/sales/{publicId}`
- `PATCH /api/sales/sales/{publicId}/status`
- `POST /api/sales/sales/{publicId}/invoice`
- `GET /api/sales/invoices`
- `GET /api/sales/invoices/summary`
- `GET /api/sales/invoices/{publicId}`
- `PATCH /api/sales/invoices/{publicId}/status`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-golang/sales golang:1.24-alpine go test ./...`
- `docker build -t erp-sales ./service-api/service-golang/sales`

Runtime switch:

- `SALES_REPOSITORY_DRIVER=postgres`
