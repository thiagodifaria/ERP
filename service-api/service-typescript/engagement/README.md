# engagement

The engagement service owns omnichannel campaigns, touchpoints and operational follow-up across inbound and outbound journeys.

Initial scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- campaign catalog with channel, budget and workflow linkage
- template catalog with provider, subject and body by channel
- touchpoint stream linked to campaigns and CRM lead public ids
- touchpoint stream with business entity linkage for CRM-driven journeys
- delivery stream linked to touchpoints and reusable templates
- touchpoint status transitions for delivery, response and conversion
- touchpoint operational summary by channel and lifecycle
- delivery summary by provider, status and template linkage
- provider capability catalog with fallback visibility
- inbound lead ingestion linked to CRM
- workflow dispatch orchestration linked to workflow runtime ids
- provider callback ledger for inbound and outbound events
- provider events linked to touchpoints, workflow runs and business entities
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
- `GET /api/engagement/providers`
- `GET /api/engagement/providers/{provider}`
- `GET /api/engagement/provider-events`
- `GET /api/engagement/provider-events/{publicId}`
- `GET /api/engagement/provider-events/summary`
- `POST /api/engagement/providers/inbound-leads`
- `POST /api/engagement/providers/meta-ads/leads`
- `POST /api/engagement/workflow-dispatches`
- `POST /api/engagement/providers/events`
- `POST /api/engagement/providers/resend/events`
- `POST /api/engagement/providers/whatsapp-cloud/events`
- `POST /api/engagement/providers/telegram-bot/events`
- `POST /api/engagement/providers/meta-ads/events`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-typescript/engagement node:22-alpine sh -lc "npm install && npm run test:unit && npm run test:contract && npm run build"`
- `docker build -t erp-engagement ./service-api/service-typescript/engagement`

Runtime switch:

- `ENGAGEMENT_REPOSITORY_DRIVER=postgres`

Notes:

- touchpoints persist `businessEntityType` and `businessEntityPublicId` when the interaction is directly attached to a business aggregate such as `crm.lead`
- provider events inherit business linkage from their touchpoint when available, which improves callback traceability across CRM, workflow and analytics
- provider capability detail now exposes credential key, fallback posture and readiness mode, which is reused in `/health/details` and in cross-service integration readiness
