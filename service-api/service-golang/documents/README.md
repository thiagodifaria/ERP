# documents

The documents service owns attachment metadata, upload orchestration, owner associations, archive governance and future storage lifecycle policies.

Initial scope:

- health and readiness endpoints
- attachment metadata registry for tenant-owned aggregates
- generic owner references such as `crm.lead` and `crm.customer`
- repository driver selection between memory and PostgreSQL
- first public routes for listing and creating attachment references
- upload sessions for storage adapters and deferred attachment completion
- tenant-aware filtering by `tenantSlug`, `ownerType` and `ownerPublicId`
- attachment detail lookup, archive control, signed access-link issuance and download handoff
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
- `GET /api/documents/attachments/{publicId}/download`
- `POST /api/documents/attachments/{publicId}/archive`
- `POST /api/documents/attachments/{publicId}/access-links`
- `POST /api/documents/upload-sessions`
- `GET /api/documents/upload-sessions/{publicId}`
- `POST /api/documents/upload-sessions/{publicId}/complete`

Query and payload conventions:

- `GET /api/documents/attachments` accepts `tenantSlug`, `ownerType`, `ownerPublicId`, `source`, `visibility` and `archived`
- `POST /api/documents/attachments` expects tenant-aware attachment metadata and stores only attachment references for now
- `POST /api/documents/attachments/{publicId}/archive` archives the attachment reference without deleting metadata history
- `POST /api/documents/attachments/{publicId}/access-links` returns a signed operational access URL for downstream storage adapters
- `GET /api/documents/attachments/{publicId}/download` validates the signed token and hands off the request to the configured storage path
- `POST /api/documents/upload-sessions` opens a tenant-aware upload session with storage key, expiration and completion URL
- `POST /api/documents/upload-sessions/{publicId}/complete` finalizes the session and materializes the attachment metadata with file size and checksum
