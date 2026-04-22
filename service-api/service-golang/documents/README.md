# documents

The documents service owns attachment metadata, owner associations, archive governance and future storage lifecycle policies.

Initial scope:

- health and readiness endpoints
- attachment metadata registry for tenant-owned aggregates
- generic owner references such as `crm.lead` and `crm.customer`
- repository driver selection between memory and PostgreSQL
- first public routes for listing and creating attachment references
- tenant-aware filtering by `tenantSlug`, `ownerType` and `ownerPublicId`
- attachment detail lookup, archive control and temporary access-link issuance
- governance metadata such as file size, checksum, visibility and retention
- PostgreSQL-backed attachment metadata ready for CRM attachment ownership flows
- container-first runtime validation against `bootstrap-ops` and `northwind-group`

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/documents/attachments`
- `POST /api/documents/attachments`
- `GET /api/documents/attachments/{publicId}`
- `POST /api/documents/attachments/{publicId}/archive`
- `POST /api/documents/attachments/{publicId}/access-links`

Query and payload conventions:

- `GET /api/documents/attachments` accepts `tenantSlug`, `ownerType`, `ownerPublicId`, `source`, `visibility` and `archived`
- `POST /api/documents/attachments` expects tenant-aware attachment metadata and stores only attachment references for now
- `POST /api/documents/attachments/{publicId}/archive` archives the attachment reference without deleting metadata history
- `POST /api/documents/attachments/{publicId}/access-links` returns a temporary operational access URL for downstream storage adapters
