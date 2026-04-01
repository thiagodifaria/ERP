# edge

The edge service is the public backend entry point.

Initial scope:

- health endpoints
- dynamic downstream readiness for `identity`, `crm`, `workflow-control`, `workflow-runtime`, `analytics` and `webhook-hub`
- operational health aggregation via `GET /api/edge/ops/health`
- bootstrap structure
- request correlation
- tenant-aware growth path for future middleware
