# edge

The edge service is the public backend entry point.

Initial scope:

- health endpoints
- dynamic downstream readiness for `identity`, `crm`, `sales`, `workflow-control`, `workflow-runtime`, `analytics` and `webhook-hub`
- operational health aggregation via `GET /api/edge/ops/health`
- tenant cockpit aggregation via `GET /api/edge/ops/tenant-overview`
- automation cockpit aggregation via `GET /api/edge/ops/automation-overview`
- engagement cockpit aggregation via `GET /api/edge/ops/engagement-overview`
- integration readiness cockpit aggregation via `GET /api/edge/ops/integrations-overview`
- document governance cockpit aggregation via `GET /api/edge/ops/documents-overview`
- collections and recovery cockpit aggregation via `GET /api/edge/ops/collections-overview`
- platform reliability cockpit aggregation via `GET /api/edge/ops/platform-reliability`
- hardening review cockpit aggregation via `GET /api/edge/ops/hardening-overview`
- executive summaries now surface capability gaps and contract artifact counts for external adapters
- sales cockpit aggregation via `GET /api/edge/ops/sales-overview`
- revenue cockpit aggregation via `GET /api/edge/ops/revenue-overview`
- finance cockpit aggregation via `GET /api/edge/ops/finance-overview`
- rentals cockpit aggregation via `GET /api/edge/ops/rentals-overview`
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
- `GET /api/edge/ops/engagement-overview`
- `GET /api/edge/ops/integrations-overview`
- `GET /api/edge/ops/documents-overview`
- `GET /api/edge/ops/collections-overview`
- `GET /api/edge/ops/platform-reliability`
- `GET /api/edge/ops/hardening-overview`
- `GET /api/edge/ops/sales-overview`
- `GET /api/edge/ops/revenue-overview`
- `GET /api/edge/ops/finance-overview`
- `GET /api/edge/ops/rentals-overview`

Protected routes:

- `GET /api/edge/ops/tenant-overview`
- `GET /api/edge/ops/automation-overview`
- `GET /api/edge/ops/engagement-overview`
- `GET /api/edge/ops/integrations-overview`
- `GET /api/edge/ops/documents-overview`
- `GET /api/edge/ops/collections-overview`
- `GET /api/edge/ops/platform-reliability`
- `GET /api/edge/ops/hardening-overview`
- `GET /api/edge/ops/sales-overview`
- `GET /api/edge/ops/revenue-overview`
- `GET /api/edge/ops/finance-overview`
- `GET /api/edge/ops/rentals-overview`

Protected operational routes require:

- `Authorization: Bearer <sessionToken>`
- `tenantSlug` query parameter
- successful access resolution via `GET /api/identity/tenants/{slug}/access`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-golang/edge golang:1.24-alpine go test ./...`
- `bash scripts/test.sh smoke`
