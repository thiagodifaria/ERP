# crm

The crm service owns leads, customers, relationship history and ownership flows.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- first domain entity for `Lead`
- bootstrap list and creation flow for leads
- filterable lead list by `status`, `source`, `ownerUserId`, `assigned` and `q`
- lead lookup by public id and controlled status transitions
- unit validation for bootstrap and domain basics
