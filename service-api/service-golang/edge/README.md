# edge

The edge service is the public backend entry point.

Initial scope:

- health endpoints
- dynamic downstream readiness for `identity`, `crm`, `workflow-control`, `workflow-runtime`, `analytics` and `webhook-hub`
- operational health aggregation via `GET /api/edge/ops/health`
- tenant cockpit aggregation via `GET /api/edge/ops/tenant-overview`
- automation cockpit aggregation via `GET /api/edge/ops/automation-overview`
- bootstrap structure
- request correlation
- tenant-aware growth path for future middleware
