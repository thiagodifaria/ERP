# crm

The crm service owns leads, customers, relationship history and ownership flows.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- first domain entity for `Lead`
- bootstrap list and creation flow for leads
- filterable lead list by `status`, `source`, `ownerUserId`, `assigned` and `q`
- pipeline summary endpoint for operational funnel reading
- owner assignment endpoint for taking or clearing lead responsibility
- lead public ids and owner references aligned with UUID public ids from identity
- selectable repository driver with PostgreSQL-backed lead persistence for a bootstrap tenant
- tenant-aware repository resolution by request, allowing `tenantSlug` to switch live CRM context without restarting the service
- runtime smoke now exercises the live HTTP API against PostgreSQL
- runtime health details reflect the active repository dependency state
- lead lookup by public id and controlled status transitions
- partial lead profile update endpoint for name, email and source
- lead note domain started for next-step relationship history flows
- bootstrap endpoint `GET /api/crm/leads/{publicId}/notes` for read-side relationship context
- bootstrap endpoint `POST /api/crm/leads/{publicId}/notes` for operational follow-up capture
- selectable repository driver now covers lead notes as well when `crm` runs on PostgreSQL
- contract and smoke coverage now exercise live lead note history over HTTP
- contract coverage for public HTTP routes and public error shape
- unit validation for bootstrap and domain basics
- customer conversion flow from qualified lead into active customer
- customer list and detail endpoints backed by memory and PostgreSQL
- bootstrap customer seed aligned with the first converted commercial account
- multi-tenant runtime validation now covers isolated list/create flows for `bootstrap-ops` and `northwind-group`
- relationship history and outbox event ledgers persisted for leads and customers
- every relevant mutation now appends traceable history and prepares pending integration events
- initial attachments now flow through the `documents` service for leads and customers
- smoke coverage now validates history, outbox and attachment contracts in live PostgreSQL runtime

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/crm/leads`
- `GET /api/crm/leads/summary`
- `POST /api/crm/leads`
- `GET /api/crm/leads/{publicId}`
- `POST /api/crm/leads/{publicId}/convert`
- `GET /api/crm/leads/{publicId}/history`
- `PATCH /api/crm/leads/{publicId}`
- `PATCH /api/crm/leads/{publicId}/owner`
- `PATCH /api/crm/leads/{publicId}/status`
- `GET /api/crm/leads/{publicId}/notes`
- `POST /api/crm/leads/{publicId}/notes`
- `GET /api/crm/leads/{publicId}/attachments`
- `POST /api/crm/leads/{publicId}/attachments`
- `GET /api/crm/customers`
- `GET /api/crm/customers/{publicId}`
- `GET /api/crm/customers/{publicId}/history`
- `GET /api/crm/customers/{publicId}/attachments`
- `POST /api/crm/customers/{publicId}/attachments`
- `GET /api/crm/outbox/pending`

Query conventions:

- `tenantSlug` can be provided in list/detail routes through the query string
- `POST /api/crm/leads` also accepts `tenantSlug` in the payload for explicit tenant creation
- when `tenantSlug` is omitted, the service falls back to the configured bootstrap tenant
- history and outbox routes stay tenant-scoped the same way as the read/write CRM routes
- attachments are proxied to `documents`, preserving tenant and aggregate ownership contracts
