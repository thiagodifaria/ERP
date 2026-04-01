# edge

The edge service is the public backend entry point.

Initial scope:

- health endpoints
- dynamic downstream readiness for `identity`, `crm`, `workflow-control`, `workflow-runtime`, `analytics` and `webhook-hub`
- bootstrap structure
- request correlation
- tenant-aware growth path for future middleware
