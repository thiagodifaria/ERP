# engagement

The engagement service owns omnichannel campaigns, touchpoints and operational follow-up across inbound and outbound journeys.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- campaign catalog with channel, budget and workflow linkage
- touchpoint stream linked to campaigns and CRM lead public ids
- touchpoint status transitions for delivery, response and conversion
- touchpoint operational summary by channel and lifecycle
- in-memory bootstrap data for local exploration
- unit and contract coverage for the public HTTP surface

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/engagement/campaigns`
- `POST /api/engagement/campaigns`
- `GET /api/engagement/campaigns/{publicId}`
- `PATCH /api/engagement/campaigns/{publicId}/status`
- `GET /api/engagement/touchpoints`
- `GET /api/engagement/touchpoints/summary`
- `POST /api/engagement/touchpoints`
- `GET /api/engagement/touchpoints/{publicId}`
- `PATCH /api/engagement/touchpoints/{publicId}/status`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-typescript/engagement node:22-alpine sh -lc "npm install && npm run test:unit && npm run test:contract"`
- `docker build -t erp-engagement ./service-api/service-typescript/engagement`

Runtime switch:

- `ENGAGEMENT_REPOSITORY_DRIVER=memory`
