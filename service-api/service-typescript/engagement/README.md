# engagement

The engagement service owns omnichannel campaigns, touchpoints and operational follow-up across inbound and outbound journeys.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- campaign catalog with channel, budget and workflow linkage
- template catalog with provider, subject and body by channel
- touchpoint stream linked to campaigns and CRM lead public ids
- delivery stream linked to touchpoints and reusable templates
- touchpoint status transitions for delivery, response and conversion
- touchpoint operational summary by channel and lifecycle
- delivery summary by provider, status and template linkage
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
- `GET /api/engagement/templates`
- `POST /api/engagement/templates`
- `GET /api/engagement/templates/{publicId}`
- `PATCH /api/engagement/templates/{publicId}/status`
- `GET /api/engagement/touchpoints`
- `GET /api/engagement/touchpoints/summary`
- `POST /api/engagement/touchpoints`
- `GET /api/engagement/touchpoints/{publicId}`
- `PATCH /api/engagement/touchpoints/{publicId}/status`
- `GET /api/engagement/touchpoints/{publicId}/deliveries`
- `POST /api/engagement/touchpoints/{publicId}/deliveries`
- `PATCH /api/engagement/touchpoints/{publicId}/deliveries/{deliveryPublicId}/status`
- `GET /api/engagement/deliveries/summary`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-typescript/engagement node:22-alpine sh -lc "npm install && npm run test:unit && npm run test:contract && npm run build"`
- `docker build -t erp-engagement ./service-api/service-typescript/engagement`

Runtime switch:

- `ENGAGEMENT_REPOSITORY_DRIVER=postgres`
