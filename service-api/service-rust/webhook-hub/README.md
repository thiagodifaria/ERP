# webhook-hub

The webhook-hub service protects the platform from external webhook bursts.

Initial scope:

- runtime bootstrap
- health endpoint
- router boundary
- webhook intake list and ingest routes
- webhook intake summary route
- room for signature validation, normalization and deduplication

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/webhook-hub/events`
- `GET /api/webhook-hub/events/summary`
- `GET /api/webhook-hub/events/{publicId}`
- `POST /api/webhook-hub/events`
