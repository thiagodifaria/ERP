# edge

The edge service is the public backend entry point.

Initial scope:

- health endpoints
- dynamic downstream readiness for `identity`, `crm`, `sales`, `workflow-control`, `workflow-runtime`, `analytics` and `webhook-hub`
- operational health aggregation via `GET /api/edge/ops/health`
- tenant cockpit aggregation via `GET /api/edge/ops/tenant-overview`
- automation cockpit aggregation via `GET /api/edge/ops/automation-overview`
- sales cockpit aggregation via `GET /api/edge/ops/sales-overview`
- revenue cockpit aggregation via `GET /api/edge/ops/revenue-overview`
- bootstrap structure
- request correlation
- tenant session enforcement backed by `identity`
- bearer token resolution against tenant-scoped access snapshots
- tenant-aware growth path for future middleware

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/edge/ops/health`
- `GET /api/edge/ops/tenant-overview`
- `GET /api/edge/ops/automation-overview`
- `GET /api/edge/ops/sales-overview`
- `GET /api/edge/ops/revenue-overview`

Protected routes:

- `GET /api/edge/ops/tenant-overview`
- `GET /api/edge/ops/automation-overview`
- `GET /api/edge/ops/sales-overview`
- `GET /api/edge/ops/revenue-overview`

Protected operational routes require:

- `Authorization: Bearer <sessionToken>`
- `tenantSlug` query parameter
- successful access resolution via `GET /api/identity/tenants/{slug}/access`
