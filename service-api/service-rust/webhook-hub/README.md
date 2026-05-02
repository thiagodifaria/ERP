# webhook-hub

The webhook-hub service protects the platform from external webhook bursts.

Initial scope:

- runtime bootstrap
- repository driver selection between memory and PostgreSQL
- health endpoint
- router boundary
- webhook intake list and ingest routes
- webhook intake summary route
- idempotent duplicate guard on intake by provider and external id
- operational lifecycle summary by provider and status
- dead-letter queue governance with retry metadata
- room for signature validation, normalization and deduplication

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/webhook-hub/events`
- `GET /api/webhook-hub/events/summary`
- `GET /api/webhook-hub/events/{publicId}`
- `GET /api/webhook-hub/events/{publicId}/transitions`
- `GET /api/webhook-hub/events/dead-letter`
- `POST /api/webhook-hub/events`
- `POST /api/webhook-hub/events/{publicId}/validate`
- `POST /api/webhook-hub/events/{publicId}/queue`
- `POST /api/webhook-hub/events/{publicId}/process`
- `POST /api/webhook-hub/events/{publicId}/forward`
- `POST /api/webhook-hub/events/{publicId}/fail`
- `POST /api/webhook-hub/events/{publicId}/reject`
- `POST /api/webhook-hub/events/{publicId}/dead-letter`
- `POST /api/webhook-hub/events/{publicId}/requeue`

Event query params:

- `provider=<provider>`
- `event_type=<eventType>`
- `status=<status>`

Runtime switch:

- `WEBHOOK_HUB_REPOSITORY_DRIVER=postgres`
