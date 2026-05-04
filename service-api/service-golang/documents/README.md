# documents

The documents service owns attachment metadata, upload orchestration, owner associations, archive governance and future storage lifecycle policies.

Initial scope:

- health and readiness endpoints
- attachment metadata registry for tenant-owned aggregates
- generic owner references such as `crm.lead` and `crm.customer`
- repository driver selection between memory and PostgreSQL
- digital signature capability registry with local fallback and provider posture
- first public routes for listing and creating attachment references
- upload sessions for storage adapters and deferred attachment completion
- tenant-aware filtering by `tenantSlug`, `ownerType` and `ownerPublicId`
- attachment detail lookup, archive control, signed access-link issuance and download handoff
- attachment version history with immutable snapshots of file metadata and storage pointers
- governance metadata such as file size, checksum, visibility and retention
- PostgreSQL-backed attachment metadata ready for CRM attachment ownership flows
- container-first runtime validation against `bootstrap-ops` and `northwind-group`

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/documents/storage/capabilities`
- `GET /api/documents/storage/capabilities/{provider}`
- `GET /api/documents/signing/capabilities`
- `GET /api/documents/signing/capabilities/{provider}`
- `POST /api/documents/signing/requests`
- `GET /api/documents/attachments`
- `POST /api/documents/attachments`
- `GET /api/documents/attachments/{publicId}`
- `GET /api/documents/attachments/{publicId}/versions`
- `GET /api/documents/attachments/{publicId}/download`
- `POST /api/documents/attachments/{publicId}/archive`
- `POST /api/documents/attachments/{publicId}/access-links`
- `POST /api/documents/attachments/{publicId}/versions`
- `POST /api/documents/upload-sessions`
- `GET /api/documents/upload-sessions/{publicId}`
- `POST /api/documents/upload-sessions/{publicId}/complete`

Query and payload conventions:

- `GET /api/documents/attachments` accepts `tenantSlug`, `ownerType`, `ownerPublicId`, `source`, `visibility` and `archived`
- `POST /api/documents/attachments` expects tenant-aware attachment metadata and stores only attachment references for now
- `GET /api/documents/attachments/{publicId}/versions` lists immutable attachment revisions ordered by version number
- `POST /api/documents/attachments/{publicId}/versions` appends a new revision and updates the attachment's current version pointer
- `POST /api/documents/attachments/{publicId}/archive` archives the attachment reference without deleting metadata history
- `POST /api/documents/attachments/{publicId}/access-links` returns a signed operational access URL for downstream storage adapters
- `GET /api/documents/attachments/{publicId}/download` validates the signed token and hands off the request to the configured storage path
- `POST /api/documents/upload-sessions` opens a tenant-aware upload session with storage key, expiration and completion URL
- `POST /api/documents/upload-sessions/{publicId}/complete` finalizes the session and materializes the attachment metadata with file size and checksum
- `GET /api/documents/storage/capabilities` exposes local, S3-compatible and R2 capability posture, including fallback state for environments without a paid provider
- `GET /api/documents/signing/capabilities` exposes local, Clicksign and DocuSign posture with safe fallback when no paid provider is configured
- `POST /api/documents/signing/requests` queues a provider-backed or local-simulated signing flow for attachments linked to sales, rentals or any other aggregate
