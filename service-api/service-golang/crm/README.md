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
- runtime smoke now exercises the live HTTP API against PostgreSQL
- runtime health details reflect the active repository dependency state
- lead lookup by public id and controlled status transitions
- partial lead profile update endpoint for name, email and source
- lead note domain started for next-step relationship history flows
- bootstrap endpoint `GET /api/crm/leads/{publicId}/notes` for read-side relationship context
- bootstrap endpoint `POST /api/crm/leads/{publicId}/notes` for operational follow-up capture
- contract coverage for public HTTP routes and public error shape
- unit validation for bootstrap and domain basics
