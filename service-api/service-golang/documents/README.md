# documents

The documents service owns attachment metadata, owner associations and future storage lifecycle policies.

Initial scope:

- health and readiness endpoints
- attachment metadata registry for tenant-owned aggregates
- generic owner references such as `crm.lead` and `crm.customer`
- repository driver selection between memory and PostgreSQL
- first public routes for listing and creating attachment references
- tenant-aware filtering by `tenantSlug`, `ownerType` and `ownerPublicId`
- PostgreSQL-backed attachment metadata ready for CRM attachment ownership flows
- container-first runtime validation against `bootstrap-ops` and `northwind-group`

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/documents/attachments`
- `POST /api/documents/attachments`

Query and payload conventions:

- `GET /api/documents/attachments` accepts `tenantSlug`, `ownerType` and `ownerPublicId`
- `POST /api/documents/attachments` expects tenant-aware attachment metadata and stores only attachment references for now
