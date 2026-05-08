# ERP

![ERP](https://img.shields.io/badge/ERP-Enterprise%20Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Backend-first ERP platform for enterprise operations, automation, contracts, finance, billing, analytics and integrations.**

[![Go](https://img.shields.io/badge/Go-edge%20crm%20sales%20rentals%20documents-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![.NET](https://img.shields.io/badge/.NET-identity%20finance%20billing-512BD4?style=flat&logo=dotnet&logoColor=white)](https://dotnet.microsoft.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-workflow%20engagement-3178C6?style=flat&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Elixir](https://img.shields.io/badge/Elixir-workflow%20runtime-4B275F?style=flat&logo=elixir&logoColor=white)](https://elixir-lang.org/)
[![Python](https://img.shields.io/badge/Python-analytics%20platform%20admin-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![Rust](https://img.shields.io/badge/Rust-webhook%20hub-000000?style=flat&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-domain%20schemas-316192?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-container%20first-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)

---

## Documentation

**Short overview:** [README.md](README.md)  
**English detailed README:** [README_EN.md](README_EN.md)  
**README detalhado em Portugues:** [README_PT.md](README_PT.md)  
**API Reference:** [docs/API.md](docs/API.md)  
**Architecture:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Contracts:** [docs/CONTRATOS.md](docs/CONTRATOS.md)  
**Integrations:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Services:** [docs/SERVICOS.md](docs/SERVICOS.md)

---

## What is ERP?

ERP is a multi-tenant, polyglot and container-first enterprise platform. The repository treats commercial operations, recurring contracts, documents, finance, billing, workflow automation, workflow runtime, engagement, analytics, simulation, SaaS governance, support, suppliers, notifications, fiscal/compliance and webhooks as connected parts of one platform.

The goal is to work as a technical reference and portfolio-grade product for realistic enterprise architecture. The project focuses on domain boundaries, database ownership, versioned contracts, integrated smoke validation, real health checks, external adapters and operational automation.

## Current Scale

| Metric | Value |
|--------|-------|
| Services with OpenAPI contracts | 20 |
| Versioned HTTP endpoints | 201 |
| Contract catalog | `docs/contracts/` |
| Runtime command | `./scripts/build.sh` |
| Validation command | `./scripts/test.sh` |

## Architecture Planes

- administrative plane: `support`
- administrative/notification plane: `notification`
- administrative/procurement plane: `supplier`
- analytics plane: `analytics`
- compliance plane: `fiscal`
- control plane: `workflow-control`
- integration plane: `webhook-hub`
- interaction/control plane: `engagement`
- platform control plane: `platform-control`
- public operations plane: `edge`
- runtime plane: `workflow-runtime`
- simulation plane: `simulation`
- transaction plane: `crm`, `documents`, `rentals`, `sales`
- transaction/billing plane: `billing`
- transaction/catalog plane: `catalog`
- transaction/finance plane: `finance`
- transaction/security plane: `identity`

## Stack Rationale

- .NET: used by `billing`, `finance`, `identity` so each workload stays close to the ecosystem that fits its operational shape instead of forcing every problem into one language.
- Elixir: used by `workflow-runtime` so each workload stays close to the ecosystem that fits its operational shape instead of forcing every problem into one language.
- Go: used by `crm`, `documents`, `edge`, `rentals`, `sales` so each workload stays close to the ecosystem that fits its operational shape instead of forcing every problem into one language.
- Python: used by `analytics`, `catalog`, `fiscal`, `notification`, `platform-control`, `simulation`, `supplier`, `support` so each workload stays close to the ecosystem that fits its operational shape instead of forcing every problem into one language.
- Rust: used by `webhook-hub` so each workload stays close to the ecosystem that fits its operational shape instead of forcing every problem into one language.
- TypeScript: used by `engagement`, `workflow-control` so each workload stays close to the ecosystem that fits its operational shape instead of forcing every problem into one language.

## Service Inventory

| Service | Stack | Plane | Responsibility | Endpoints |
|---------|-------|-------|----------------|-----------|
| `analytics` | Python | analytics plane | operational reports, governance, reliability, hardening, cost and executive reads | 9 |
| `billing` | .NET | transaction/billing plane | plans, subscriptions, recurring invoices, payment attempts and recovery | 9 |
| `catalog` | Python | transaction/catalog plane | categories, items, item versions, bulk creation and consumer contracts | 12 |
| `crm` | Go | transaction plane | leads, customers, ownership, pipeline, notes, history, attachments and enrichment | 5 |
| `documents` | Go | transaction plane | attachments, upload, storage posture, signing, versions, archive and access links | 10 |
| `edge` | Go | public operations plane | public entrypoint, cross-service aggregation and operational cockpits | 8 |
| `engagement` | TypeScript | interaction/control plane | campaigns, templates, touchpoints, conversations, delivery, providers and callbacks | 9 |
| `finance` | .NET | transaction/finance plane | receivables, payables, treasury, costs, commissions and closing | 5 |
| `fiscal` | Python | compliance plane | fiscal profiles, retention, tax documents, privacy, consent and audit | 25 |
| `identity` | .NET | transaction/security plane | tenants, companies, users, teams, roles, sessions, invites, MFA and audit | 6 |
| `notification` | Python | administrative/notification plane | preferences, internal alert center, severity and notification lifecycle | 7 |
| `platform-control` | Python | platform control plane | capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle and go-live | 40 |
| `rentals` | Go | transaction plane | recurring contracts, adjustments, terminations, charges and contractual attachments | 4 |
| `sales` | Go | transaction plane | opportunities, proposals, sales, invoices, commissions, renegotiations and pending items | 6 |
| `simulation` | Python | simulation plane | operational scenarios, load benchmarks and capacity modeling | 3 |
| `supplier` | Python | administrative/procurement plane | supplier categories, supplier directory and procurement ownership | 8 |
| `support` | Python | administrative plane | queues, cases, SLA, comments, bulk operations and support summaries | 9 |
| `webhook-hub` | Rust | integration plane | webhook intake, idempotency, transitions, DLQ, outbound endpoints and deliveries | 13 |
| `workflow-control` | TypeScript | control plane | definitions, published versions, trigger/action catalogs, runs and events | 7 |
| `workflow-runtime` | Elixir | runtime plane | durable executions, timeline, actions, transitions, retries, waits and compensations | 6 |

## Main Vertical Slices

### Commercial journey

lead capture in CRM, opportunity and proposal in sales, conversion, invoicing, finance visibility and edge cockpit

### Recurring contract journey

customer context, rental contract, scheduled charges, adjustments, termination, documents and financial projections

### Automation journey

workflow definition, published version, workflow run, runtime execution, timeline, retries and analytics visibility

### Integration journey

provider event, webhook-hub intake, validation, queueing, processing, forwarding, dead letter and analytics posture

### SaaS governance journey

capability catalog, provider defaults, entitlements, quotas, metering, onboarding/offboarding and go-live controls

## Local Development

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh logs edge
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh down
```

## Validation

```bash
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh contract
./scripts/test.sh platform
./scripts/test.sh smoke
./scripts/test.sh performance
./scripts/test.sh backup-restore
./scripts/test.sh hardening
```

## Contract Catalog

- `docs/contracts/http/`
- `docs/contracts/events/`
- `docs/contracts/registry.json`
- `docs/contracts/schema-registry.json`
- `docs/contracts/portal/index.html`

## Detailed Services

### `analytics`

- Stack: Python
- Plane: analytics plane
- Code: `service-api/service-python/analytics`
- Database context: `analytics/simulation read models`
- Contract: `docs/contracts/http/analytics.openapi.yaml`
- OpenAPI version: `0.1.0`
- Responsibility: operational reports, governance, reliability, hardening, cost and executive reads.

Routes:

- `GET /api/analytics/reports/adapter-catalog` - Read external adapter capability catalog
- `GET /api/analytics/reports/integration-readiness` - Read external integration readiness
- `GET /api/analytics/reports/saas-control` - Read SaaS control posture by tenant
- `GET /api/analytics/reports/contract-governance` - Read contract governance posture
- `GET /api/analytics/reports/hardening-review` - Read hardening review
- `GET /api/analytics/reports/core-operations` - Read core product operations
- `GET /api/analytics/reports/relationship-intelligence` - Read relationship intelligence
- `GET /api/analytics/reports/compliance-control` - Read fiscal and privacy compliance control
- `GET /api/analytics/reports/go-live-control` - Read go-live rollout control

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `billing`

- Stack: .NET
- Plane: transaction/billing plane
- Code: `service-api/service-csharp/billing`
- Database context: `billing`
- Contract: `docs/contracts/http/billing.openapi.yaml`
- OpenAPI version: `0.9.7`
- Responsibility: plans, subscriptions, recurring invoices, payment attempts and recovery.

Routes:

- `GET /health/details` - Return readiness details and gateway posture
- `GET /api/billing/gateways` - List gateway capabilities and Pix posture
- `GET /api/billing/gateways/{provider}` - Read one gateway capability
- `GET /api/billing/plans` - List billing plans including flat, hybrid and usage-based pricing
- `POST /api/billing/plans` - Create billing plan
- `GET /api/billing/subscriptions` - List subscriptions
- `POST /api/billing/subscriptions` - Create subscription
- `GET /api/billing/subscriptions/{publicId}/usage-pricing` - Project usage-based charge for one subscription
- `POST /api/billing/invoices/{publicId}/attempts` - Create payment attempt with idempotency support

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `catalog`

- Stack: Python
- Plane: transaction/catalog plane
- Code: `service-api/service-python/catalog`
- Database context: `catalog`
- Contract: `docs/contracts/http/catalog.openapi.yaml`
- OpenAPI version: `0.2.0`
- Responsibility: categories, items, item versions, bulk creation and consumer contracts.

Routes:

- `GET /api/catalog/capabilities` - Read catalog capability posture
- `GET /api/catalog/consumers` - Read catalog consumer contracts across core domains
- `GET /api/catalog/categories` - List categories by tenant
- `POST /api/catalog/categories` - Create one category
- `GET /api/catalog/categories/page` - Cursor-based category listing
- `GET /api/catalog/items` - List catalog items
- `POST /api/catalog/items` - Create one catalog item
- `GET /api/catalog/items/page` - Cursor-based item listing
- `POST /api/catalog/items/bulk` - Bulk create catalog items with partial success
- `GET /api/catalog/items/{publicId}` - Read one catalog item
- `PATCH /api/catalog/items/{publicId}` - Update active state, price and attributes
- `GET /api/catalog/items/{publicId}/versions` - Read catalog item version history

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `crm`

- Stack: Go
- Plane: transaction plane
- Code: `service-api/service-golang/crm`
- Database context: `crm`
- Contract: `docs/contracts/http/crm.openapi.yaml`
- OpenAPI version: `0.2.0`
- Responsibility: leads, customers, ownership, pipeline, notes, history, attachments and enrichment.

Routes:

- `GET /api/crm/enrichment/cnpj/capabilities` - Read CNPJ enrichment provider capabilities
- `POST /api/crm/enrichment/cnpj/lookup` - Lookup and enrich one CNPJ through provider contract
- `GET /api/crm/pipeline/config` - Read tenant pipeline configuration
- `PUT /api/crm/pipeline/config` - Upsert tenant pipeline configuration
- `GET /api/crm/leads/intelligence/summary` - Read lead scoring and pipeline intelligence summary

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `documents`

- Stack: Go
- Plane: transaction plane
- Code: `service-api/service-golang/documents`
- Database context: `documents`
- Contract: `docs/contracts/http/documents.openapi.yaml`
- OpenAPI version: `0.9.7`
- Responsibility: attachments, upload, storage posture, signing, versions, archive and access links.

Routes:

- `GET /health/details` - Return runtime readiness and storage posture
- `GET /api/documents/signing/capabilities` - List digital signature capabilities
- `GET /api/documents/signing/capabilities/{provider}` - Read one signing capability
- `POST /api/documents/signing/requests` - Queue one digital signature request
- `GET /api/documents/storage/capabilities` - List storage capability registry
- `GET /api/documents/storage/capabilities/{provider}` - Read one storage capability
- `GET /api/documents/attachments` - List attachments
- `POST /api/documents/attachments` - Create attachment metadata
- `GET /api/documents/attachments/{publicId}/versions` - List attachment versions
- `POST /api/documents/attachments/{publicId}/versions` - Append attachment version

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `edge`

- Stack: Go
- Plane: public operations plane
- Code: `service-api/service-golang/edge`
- Database context: `none`
- Contract: `docs/contracts/http/edge.openapi.yaml`
- OpenAPI version: `0.1.0`
- Responsibility: public entrypoint, cross-service aggregation and operational cockpits.

Routes:

- `GET /api/edge/ops/core-operations` - Read executive core product cockpit
- `GET /api/edge/ops/relationship-overview` - Read executive relationship cockpit
- `GET /api/edge/ops/compliance-overview` - Read executive compliance cockpit
- `GET /api/edge/ops/go-live-overview` - Read executive go-live cockpit
- `GET /api/edge/ops/integrations-overview` - Read executive integrations cockpit
- `GET /api/edge/ops/saas-overview` - Read executive SaaS cockpit
- `GET /api/edge/ops/contracts-overview` - Read executive contracts cockpit
- `GET /api/edge/ops/hardening-overview` - Read executive hardening cockpit

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `engagement`

- Stack: TypeScript
- Plane: interaction/control plane
- Code: `service-api/service-typescript/engagement`
- Database context: `engagement`
- Contract: `docs/contracts/http/engagement.openapi.yaml`
- OpenAPI version: `0.9.7`
- Responsibility: campaigns, templates, touchpoints, conversations, delivery, providers and callbacks.

Routes:

- `GET /health/details` - Return readiness details for engagement runtime
- `GET /api/engagement/providers` - List provider capabilities and fallback posture
- `GET /api/engagement/providers/{provider}` - Read one provider capability
- `POST /api/engagement/providers/meta-ads/leads` - Ingest inbound lead from Meta Ads
- `POST /api/engagement/providers/resend/events` - Register Resend callback event
- `POST /api/engagement/providers/whatsapp-cloud/events` - Register WhatsApp callback event
- `POST /api/engagement/providers/telegram-bot/events` - Register Telegram callback event
- `GET /api/engagement/provider-events` - List provider events
- `GET /api/engagement/provider-events/{publicId}` - Read one provider event

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `finance`

- Stack: .NET
- Plane: transaction/finance plane
- Code: `service-api/service-csharp/finance`
- Database context: `finance`
- Contract: `docs/contracts/http/finance.openapi.yaml`
- OpenAPI version: `0.4.0`
- Responsibility: receivables, payables, treasury, costs, commissions and closing.

Routes:

- `GET /api/finance/receivable-projections` - List receivable projections
- `POST /api/finance/receivable-projections/sync` - Sync projections from sales and rentals
- `GET /api/finance/commission-holds` - List commission holds
- `POST /api/finance/commission-holds/{publicId}/release` - Release one commission hold
- `GET /api/finance/activity` - List finance operational activity

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `fiscal`

- Stack: Python
- Plane: compliance plane
- Code: `service-api/service-python/fiscal`
- Database context: `fiscal`
- Contract: `docs/contracts/http/fiscal.openapi.yaml`
- OpenAPI version: `0.1.0`
- Responsibility: fiscal profiles, retention, tax documents, privacy, consent and audit.

Routes:

- `GET /api/fiscal/capabilities` - Read fiscal capability registry
- `GET /api/fiscal/companies/{companyPublicId}/profile` - Read fiscal company profile
- `PUT /api/fiscal/companies/{companyPublicId}/profile` - Upsert fiscal company profile
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies` - List retention policies by company
- `GET /api/fiscal/companies/{companyPublicId}/retention-execution` - Read retention execution plan for one company
- `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute` - Execute retention and anonymization plan
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}` - Upsert retention policy for one data domain
- `GET /api/fiscal/documents` - List fiscal documents
- `POST /api/fiscal/documents` - Issue one fiscal document
- `GET /api/fiscal/documents/{publicId}` - Read one fiscal document
- `POST /api/fiscal/documents/{publicId}/cancel` - Cancel one fiscal document
- `POST /api/fiscal/documents/{publicId}/correction-letter` - Register correction letter for one fiscal document
- `POST /api/fiscal/documents/{publicId}/invalidate` - Register invalidation for one fiscal document
- `GET /api/fiscal/documents/{publicId}/events` - List fiscal document audit events
- `GET /api/fiscal/privacy-requests` - List privacy requests
- `POST /api/fiscal/privacy-requests` - Create privacy request
- `GET /api/fiscal/privacy-requests/{publicId}` - Read one privacy request
- `GET /api/fiscal/privacy-requests/{publicId}/export-package` - Build export package for one privacy request
- `POST /api/fiscal/privacy-requests/{publicId}/execute` - Execute one privacy request with audit trail
- `PATCH /api/fiscal/privacy-requests/{publicId}/status` - Transition privacy request lifecycle status
- `GET /api/fiscal/consents` - List consent ledger
- `POST /api/fiscal/consents` - Create consent record
- `PATCH /api/fiscal/consents/{publicId}` - Transition consent status
- `GET /api/fiscal/audit-events` - List fiscal audit events
- `GET /api/fiscal/compliance/summary` - Read fiscal compliance summary

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `identity`

- Stack: .NET
- Plane: transaction/security plane
- Code: `service-api/service-csharp/identity`
- Database context: `identity`
- Contract: `docs/contracts/http/identity.openapi.yaml`
- OpenAPI version: `0.5.0`
- Responsibility: tenants, companies, users, teams, roles, sessions, invites, MFA and audit.

Routes:

- `GET /api/identity/tenants` - List tenants
- `POST /api/identity/tenants` - Create tenant
- `GET /api/identity/tenants/{slug}/snapshot` - Read one tenant snapshot
- `POST /api/identity/sessions/login` - Authenticate identity session
- `POST /api/identity/sessions/refresh` - Refresh identity session
- `POST /api/identity/invitations` - Create invitation

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `notification`

- Stack: Python
- Plane: administrative/notification plane
- Code: `service-api/service-python/notification`
- Database context: `notification`
- Contract: `docs/contracts/http/notification.openapi.yaml`
- OpenAPI version: `0.1.0`
- Responsibility: preferences, internal alert center, severity and notification lifecycle.

Routes:

- `GET /api/notification/capabilities` - Read notification capability catalog
- `GET /api/notification/preferences/{userPublicId}` - Read one user notification preference
- `PUT /api/notification/preferences/{userPublicId}` - Upsert one user notification preference
- `GET /api/notification/center` - List notification center items with cursor filters
- `POST /api/notification/center` - Create one notification center item
- `PATCH /api/notification/center/{publicId}/status` - Transition notification status
- `GET /api/notification/summary` - Read notification summary

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `platform-control`

- Stack: Python
- Plane: platform control plane
- Code: `service-api/service-python/platform-control`
- Database context: `platform-control`
- Contract: `docs/contracts/http/platform-control.openapi.yaml`
- OpenAPI version: `0.2.0`
- Responsibility: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle and go-live.

Routes:

- `GET /api/platform-control/capabilities/catalog` - List platform capability catalog
- `GET /api/platform-control/providers/catalog` - List provider capability catalog and environment posture
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements` - List tenant entitlements with cursor pagination
- `GET /api/platform-control/tenants/{tenantSlug}/feature-flags` - List tenant feature flags with capability metadata
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}` - Upsert one entitlement
- `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}` - Upsert one feature flag using entitlement governance
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk` - Bulk upsert entitlements with partial success
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults` - List provider defaults selected for one tenant
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}` - Upsert provider default for one tenant capability
- `GET /api/platform-control/tenants/{tenantSlug}/quotas` - List quotas by tenant
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}` - Upsert one quota
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk` - Bulk upsert quotas with partial success
- `GET /api/platform-control/tenants/{tenantSlug}/blocks` - List tenant blocks
- `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}` - Upsert tenant block
- `GET /api/platform-control/tenants/{tenantSlug}/metering` - Read metering snapshots and summary with cursor pagination
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots` - Create one usage snapshot
- `GET /api/platform-control/tenants/{tenantSlug}/usage-summary` - Read quota and metering utilization summary
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness` - Read tenant lifecycle readiness and provider posture
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs` - List onboarding and offboarding jobs with cursor pagination
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}` - Read one lifecycle job with audit events
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview` - Preview onboarding plan, provider defaults and lifecycle readiness
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding` - Queue onboarding job with Idempotency-Key and 202 Accepted
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview` - Preview offboarding plan, retention posture and lifecycle readiness
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding` - Queue offboarding job with Idempotency-Key and 202 Accepted
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start` - Transition lifecycle job to running
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete` - Transition lifecycle job to completed
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail` - Transition lifecycle job to failed
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel` - Transition lifecycle job to cancelled
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` - Read go-live rollout readiness by tenant
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` - Read tenant go-live adoption baseline and gap
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` - List go-live bottlenecks and operational blockers
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` - Read rollout and rollback playbook for one tenant
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` - List recommended go-live adjustments
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply` - Apply one go-live operational adjustment
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - List go-live rollouts
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - Create one go-live rollout
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}` - Read one go-live rollout with events
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` - Transition go-live rollout to running
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` - Transition go-live rollout to completed
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` - Roll back one go-live rollout

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `rentals`

- Stack: Go
- Plane: transaction plane
- Code: `service-api/service-golang/rentals`
- Database context: `rentals`
- Contract: `docs/contracts/http/rentals.openapi.yaml`
- OpenAPI version: `0.8.0`
- Responsibility: recurring contracts, adjustments, terminations, charges and contractual attachments.

Routes:

- `GET /api/rentals/contracts` - List rental contracts
- `POST /api/rentals/contracts` - Create rental contract
- `GET /api/rentals/contracts/{publicId}/charges` - List contract charges
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` - Update charge status

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `sales`

- Stack: Go
- Plane: transaction plane
- Code: `service-api/service-golang/sales`
- Database context: `sales`
- Contract: `docs/contracts/http/sales.openapi.yaml`
- OpenAPI version: `0.7.0`
- Responsibility: opportunities, proposals, sales, invoices, commissions, renegotiations and pending items.

Routes:

- `GET /api/sales/opportunities` - List opportunities
- `POST /api/sales/opportunities` - Create opportunity
- `GET /api/sales/proposals` - List proposals
- `POST /api/sales/proposals` - Create proposal
- `GET /api/sales/sales` - List sales
- `GET /api/sales/invoices` - List commercial invoices

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `simulation`

- Stack: Python
- Plane: simulation plane
- Code: `service-api/service-python/simulation`
- Database context: `simulation`
- Contract: `docs/contracts/http/simulation.openapi.yaml`
- OpenAPI version: `0.7.0`
- Responsibility: operational scenarios, load benchmarks and capacity modeling.

Routes:

- `GET /api/simulation/scenarios` - List scenarios
- `POST /api/simulation/scenarios` - Create scenario run
- `POST /api/simulation/benchmarks/load` - Execute one load benchmark run

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `supplier`

- Stack: Python
- Plane: administrative/procurement plane
- Code: `service-api/service-python/supplier`
- Database context: `supplier`
- Contract: `docs/contracts/http/supplier.openapi.yaml`
- OpenAPI version: `0.1.0`
- Responsibility: supplier categories, supplier directory and procurement ownership.

Routes:

- `GET /api/supplier/capabilities` - Read supplier capability catalog
- `GET /api/supplier/categories` - List supplier categories
- `PUT /api/supplier/categories/{categoryKey}` - Upsert one supplier category
- `GET /api/supplier/suppliers` - List suppliers by tenant and status
- `POST /api/supplier/suppliers` - Create one supplier
- `GET /api/supplier/suppliers/summary` - Read supplier summary
- `GET /api/supplier/suppliers/{publicId}` - Read one supplier
- `PATCH /api/supplier/suppliers/{publicId}` - Update one supplier

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `support`

- Stack: Python
- Plane: administrative plane
- Code: `service-api/service-python/support`
- Database context: `support`
- Contract: `docs/contracts/http/support.openapi.yaml`
- OpenAPI version: `0.1.0`
- Responsibility: queues, cases, SLA, comments, bulk operations and support summaries.

Routes:

- `GET /api/support/capabilities` - Read support capability catalog
- `GET /api/support/queues` - List support queues by tenant
- `PUT /api/support/queues/{queueKey}` - Upsert one support queue
- `GET /api/support/cases` - List support cases with cursor filters
- `POST /api/support/cases` - Create one support case
- `GET /api/support/cases/summary` - Read support case summary
- `GET /api/support/cases/{publicId}` - Read one support case
- `PATCH /api/support/cases/{publicId}/status` - Transition support case status
- `POST /api/support/cases/{publicId}/comments` - Append comment to support case

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `webhook-hub`

- Stack: Rust
- Plane: integration plane
- Code: `service-api/service-rust/webhook-hub`
- Database context: `webhook-hub`
- Contract: `docs/contracts/http/webhook-hub.openapi.yaml`
- OpenAPI version: `0.9.7`
- Responsibility: webhook intake, idempotency, transitions, DLQ, outbound endpoints and deliveries.

Routes:

- `GET /health/details` - Return readiness details for webhook runtime
- `GET /api/webhook-hub/capabilities` - Read outbound webhook capability posture
- `GET /api/webhook-hub/outbound-endpoints` - List tenant outbound endpoints
- `POST /api/webhook-hub/outbound-endpoints` - Register one tenant outbound endpoint
- `GET /api/webhook-hub/outbound-endpoints/{publicId}` - Read one tenant outbound endpoint
- `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - List outbound delivery log for one endpoint
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - Register one outbound delivery attempt
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter` - Move one outbound delivery to dead letter
- `GET /api/webhook-hub/events` - List inbound webhook events
- `POST /api/webhook-hub/events` - Register inbound webhook event
- `GET /api/webhook-hub/events/summary` - Aggregate inbound webhook state
- `POST /api/webhook-hub/events/{publicId}/dead-letter` - Move event to dead letter queue
- `POST /api/webhook-hub/events/{publicId}/requeue` - Requeue dead-letter event

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `workflow-control`

- Stack: TypeScript
- Plane: control plane
- Code: `service-api/service-typescript/workflow-control`
- Database context: `workflow-control`
- Contract: `docs/contracts/http/workflow-control.openapi.yaml`
- OpenAPI version: `0.6.0`
- Responsibility: definitions, published versions, trigger/action catalogs, runs and events.

Routes:

- `GET /api/workflow-control/definitions` - List workflow definitions
- `POST /api/workflow-control/definitions` - Create workflow definition
- `GET /api/workflow-control/definitions/{key}` - Read one workflow definition
- `PATCH /api/workflow-control/definitions/{key}` - Update one workflow definition
- `PATCH /api/workflow-control/definitions/{key}/status` - Update workflow definition status
- `GET /api/workflow-control/capabilities/triggers` - List workflow trigger catalog
- `GET /api/workflow-control/capabilities/actions` - List workflow action catalog

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

### `workflow-runtime`

- Stack: Elixir
- Plane: runtime plane
- Code: `service-api/service-elixir/workflow-runtime`
- Database context: `workflow-runtime`
- Contract: `docs/contracts/http/workflow-runtime.openapi.yaml`
- OpenAPI version: `0.6.0`
- Responsibility: durable executions, timeline, actions, transitions, retries, waits and compensations.

Routes:

- `GET /api/workflow-runtime/executions` - List workflow executions
- `POST /api/workflow-runtime/executions` - Create workflow execution
- `GET /api/workflow-runtime/executions/{publicId}` - Read one workflow execution
- `GET /api/workflow-runtime/executions/{publicId}/actions` - List execution action snapshots
- `POST /api/workflow-runtime/executions/{publicId}/advance` - Advance one workflow execution
- `GET /api/workflow-runtime/capabilities` - List runtime capabilities

Operational notes:

- Keeps clear domain ownership and should not write to another context schema without an explicit decision.
- Must keep health checks real, contracts updated and central validation aligned when the public surface changes.
- Relevant mutations must preserve tenant, actor, correlation id and operational history when applicable.

---

**Thiago Di Faria** - thiagodifaria@gmail.com

## HTTP Contract Detail

This section mirrors the current OpenAPI catalog at a high level for quick repository reading.

### `analytics`

- Contract: `docs/contracts/http/analytics.openapi.yaml`
- Endpoints: `9`

#### `GET /api/analytics/reports/adapter-catalog`

- Summary: Read external adapter capability catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/integration-readiness`

- Summary: Read external integration readiness. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/saas-control`

- Summary: Read SaaS control posture by tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/contract-governance`

- Summary: Read contract governance posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/hardening-review`

- Summary: Read hardening review. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/core-operations`

- Summary: Read core product operations. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/relationship-intelligence`

- Summary: Read relationship intelligence. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/compliance-control`

- Summary: Read fiscal and privacy compliance control. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/analytics/reports/go-live-control`

- Summary: Read go-live rollout control. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `billing`

- Contract: `docs/contracts/http/billing.openapi.yaml`
- Endpoints: `9`

#### `GET /health/details`

- Summary: Return readiness details and gateway posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/billing/gateways`

- Summary: List gateway capabilities and Pix posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/billing/gateways/{provider}`

- Summary: Read one gateway capability. 
- Parameters: `provider`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/billing/plans`

- Summary: List billing plans including flat, hybrid and usage-based pricing. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/billing/plans`

- Summary: Create billing plan. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/billing/subscriptions`

- Summary: List subscriptions. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/billing/subscriptions`

- Summary: Create subscription. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/billing/subscriptions/{publicId}/usage-pricing`

- Summary: Project usage-based charge for one subscription. 
- Parameters: `publicId`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/billing/invoices/{publicId}/attempts`

- Summary: Create payment attempt with idempotency support. 
- Parameters: `Idempotency-Key`, `publicId`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

### `catalog`

- Contract: `docs/contracts/http/catalog.openapi.yaml`
- Endpoints: `12`

#### `GET /api/catalog/capabilities`

- Summary: Read catalog capability posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/catalog/consumers`

- Summary: Read catalog consumer contracts across core domains. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/catalog/categories`

- Summary: List categories by tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/catalog/categories`

- Summary: Create one category. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/catalog/categories/page`

- Summary: Cursor-based category listing. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/catalog/items`

- Summary: List catalog items. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/catalog/items`

- Summary: Create one catalog item. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/catalog/items/page`

- Summary: Cursor-based item listing. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/catalog/items/bulk`

- Summary: Bulk create catalog items with partial success. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/catalog/items/{publicId}`

- Summary: Read one catalog item. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `PATCH /api/catalog/items/{publicId}`

- Summary: Update active state, price and attributes. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `GET /api/catalog/items/{publicId}/versions`

- Summary: Read catalog item version history. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `crm`

- Contract: `docs/contracts/http/crm.openapi.yaml`
- Endpoints: `5`

#### `GET /api/crm/enrichment/cnpj/capabilities`

- Summary: Read CNPJ enrichment provider capabilities. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/crm/enrichment/cnpj/lookup`

- Summary: Lookup and enrich one CNPJ through provider contract. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/crm/pipeline/config`

- Summary: Read tenant pipeline configuration. 
- Parameters: `tenantSlug`.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/crm/pipeline/config`

- Summary: Upsert tenant pipeline configuration. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/crm/leads/intelligence/summary`

- Summary: Read lead scoring and pipeline intelligence summary. 
- Parameters: `tenantSlug`.
- Request body: not declared.
- Responses: `200`.

### `documents`

- Contract: `docs/contracts/http/documents.openapi.yaml`
- Endpoints: `10`

#### `GET /health/details`

- Summary: Return runtime readiness and storage posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/documents/signing/capabilities`

- Summary: List digital signature capabilities. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/documents/signing/capabilities/{provider}`

- Summary: Read one signing capability. 
- Parameters: `provider`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/documents/signing/requests`

- Summary: Queue one digital signature request. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/documents/storage/capabilities`

- Summary: List storage capability registry. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/documents/storage/capabilities/{provider}`

- Summary: Read one storage capability. 
- Parameters: `provider`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/documents/attachments`

- Summary: List attachments. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/documents/attachments`

- Summary: Create attachment metadata. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/documents/attachments/{publicId}/versions`

- Summary: List attachment versions. 
- Parameters: `publicId`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/documents/attachments/{publicId}/versions`

- Summary: Append attachment version. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

### `edge`

- Contract: `docs/contracts/http/edge.openapi.yaml`
- Endpoints: `8`

#### `GET /api/edge/ops/core-operations`

- Summary: Read executive core product cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/relationship-overview`

- Summary: Read executive relationship cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/compliance-overview`

- Summary: Read executive compliance cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/go-live-overview`

- Summary: Read executive go-live cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/integrations-overview`

- Summary: Read executive integrations cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/saas-overview`

- Summary: Read executive SaaS cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/contracts-overview`

- Summary: Read executive contracts cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/edge/ops/hardening-overview`

- Summary: Read executive hardening cockpit. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `engagement`

- Contract: `docs/contracts/http/engagement.openapi.yaml`
- Endpoints: `9`

#### `GET /health/details`

- Summary: Return readiness details for engagement runtime. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/engagement/providers`

- Summary: List provider capabilities and fallback posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/engagement/providers/{provider}`

- Summary: Read one provider capability. 
- Parameters: `provider`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/engagement/providers/meta-ads/leads`

- Summary: Ingest inbound lead from Meta Ads. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `POST /api/engagement/providers/resend/events`

- Summary: Register Resend callback event. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `POST /api/engagement/providers/whatsapp-cloud/events`

- Summary: Register WhatsApp callback event. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `POST /api/engagement/providers/telegram-bot/events`

- Summary: Register Telegram callback event. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/engagement/provider-events`

- Summary: List provider events. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/engagement/provider-events/{publicId}`

- Summary: Read one provider event. 
- Parameters: `publicId`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

### `finance`

- Contract: `docs/contracts/http/finance.openapi.yaml`
- Endpoints: `5`

#### `GET /api/finance/receivable-projections`

- Summary: List receivable projections. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/finance/receivable-projections/sync`

- Summary: Sync projections from sales and rentals. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/finance/commission-holds`

- Summary: List commission holds. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/finance/commission-holds/{publicId}/release`

- Summary: Release one commission hold. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/finance/activity`

- Summary: List finance operational activity. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `fiscal`

- Contract: `docs/contracts/http/fiscal.openapi.yaml`
- Endpoints: `25`

#### `GET /api/fiscal/capabilities`

- Summary: Read fiscal capability registry. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Read fiscal company profile. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Upsert fiscal company profile. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/companies/{companyPublicId}/retention-policies`

- Summary: List retention policies by company. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/companies/{companyPublicId}/retention-execution`

- Summary: Read retention execution plan for one company. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute`

- Summary: Execute retention and anonymization plan. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`

- Summary: Upsert retention policy for one data domain. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/documents`

- Summary: List fiscal documents. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/fiscal/documents`

- Summary: Issue one fiscal document. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `201`.

#### `GET /api/fiscal/documents/{publicId}`

- Summary: Read one fiscal document. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `POST /api/fiscal/documents/{publicId}/cancel`

- Summary: Cancel one fiscal document. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/fiscal/documents/{publicId}/correction-letter`

- Summary: Register correction letter for one fiscal document. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/fiscal/documents/{publicId}/invalidate`

- Summary: Register invalidation for one fiscal document. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/documents/{publicId}/events`

- Summary: List fiscal document audit events. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/privacy-requests`

- Summary: List privacy requests. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/fiscal/privacy-requests`

- Summary: Create privacy request. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `201`.

#### `GET /api/fiscal/privacy-requests/{publicId}`

- Summary: Read one privacy request. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `GET /api/fiscal/privacy-requests/{publicId}/export-package`

- Summary: Build export package for one privacy request. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `POST /api/fiscal/privacy-requests/{publicId}/execute`

- Summary: Execute one privacy request with audit trail. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `PATCH /api/fiscal/privacy-requests/{publicId}/status`

- Summary: Transition privacy request lifecycle status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `GET /api/fiscal/consents`

- Summary: List consent ledger. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/fiscal/consents`

- Summary: Create consent record. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `201`.

#### `PATCH /api/fiscal/consents/{publicId}`

- Summary: Transition consent status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`, `404`.

#### `GET /api/fiscal/audit-events`

- Summary: List fiscal audit events. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/fiscal/compliance/summary`

- Summary: Read fiscal compliance summary. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `identity`

- Contract: `docs/contracts/http/identity.openapi.yaml`
- Endpoints: `6`

#### `GET /api/identity/tenants`

- Summary: List tenants. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/identity/tenants`

- Summary: Create tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/identity/tenants/{slug}/snapshot`

- Summary: Read one tenant snapshot. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/identity/sessions/login`

- Summary: Authenticate identity session. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/identity/sessions/refresh`

- Summary: Refresh identity session. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/identity/invitations`

- Summary: Create invitation. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `notification`

- Contract: `docs/contracts/http/notification.openapi.yaml`
- Endpoints: `7`

#### `GET /api/notification/capabilities`

- Summary: Read notification capability catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/notification/preferences/{userPublicId}`

- Summary: Read one user notification preference. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/notification/preferences/{userPublicId}`

- Summary: Upsert one user notification preference. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/notification/center`

- Summary: List notification center items with cursor filters. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/notification/center`

- Summary: Create one notification center item. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `201`.

#### `PATCH /api/notification/center/{publicId}/status`

- Summary: Transition notification status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/notification/summary`

- Summary: Read notification summary. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `platform-control`

- Contract: `docs/contracts/http/platform-control.openapi.yaml`
- Endpoints: `40`

#### `GET /api/platform-control/capabilities/catalog`

- Summary: List platform capability catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/providers/catalog`

- Summary: List provider capability catalog and environment posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/entitlements`

- Summary: List tenant entitlements with cursor pagination. 
- Parameters: `tenantSlug`.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/feature-flags`

- Summary: List tenant feature flags with capability metadata. 
- Parameters: `tenantSlug`.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`

- Summary: Upsert one entitlement. 
- Parameters: `capabilityKey`, `tenantSlug`.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}`

- Summary: Upsert one feature flag using entitlement governance. 
- Parameters: `capabilityKey`, `tenantSlug`.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`

- Summary: Bulk upsert entitlements with partial success. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`

- Summary: List provider defaults selected for one tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`

- Summary: Upsert provider default for one tenant capability. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/quotas`

- Summary: List quotas by tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`

- Summary: Upsert one quota. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`

- Summary: Bulk upsert quotas with partial success. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/blocks`

- Summary: List tenant blocks. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`

- Summary: Upsert tenant block. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/metering`

- Summary: Read metering snapshots and summary with cursor pagination. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`

- Summary: Create one usage snapshot. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`

- Summary: Read quota and metering utilization summary. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`

- Summary: Read tenant lifecycle readiness and provider posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`

- Summary: List onboarding and offboarding jobs with cursor pagination. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`

- Summary: Read one lifecycle job with audit events. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`

- Summary: Preview onboarding plan, provider defaults and lifecycle readiness. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`

- Summary: Queue onboarding job with Idempotency-Key and 202 Accepted. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `202`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`

- Summary: Preview offboarding plan, retention posture and lifecycle readiness. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`

- Summary: Queue offboarding job with Idempotency-Key and 202 Accepted. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `202`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`

- Summary: Transition lifecycle job to running. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`

- Summary: Transition lifecycle job to completed. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`

- Summary: Transition lifecycle job to failed. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`

- Summary: Transition lifecycle job to cancelled. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness`

- Summary: Read go-live rollout readiness by tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption`

- Summary: Read tenant go-live adoption baseline and gap. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks`

- Summary: List go-live bottlenecks and operational blockers. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook`

- Summary: Read rollout and rollback playbook for one tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments`

- Summary: List recommended go-live adjustments. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply`

- Summary: Apply one go-live operational adjustment. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: List go-live rollouts. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: Create one go-live rollout. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}`

- Summary: Read one go-live rollout with events. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start`

- Summary: Transition go-live rollout to running. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete`

- Summary: Transition go-live rollout to completed. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback`

- Summary: Roll back one go-live rollout. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `rentals`

- Contract: `docs/contracts/http/rentals.openapi.yaml`
- Endpoints: `4`

#### `GET /api/rentals/contracts`

- Summary: List rental contracts. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/rentals/contracts`

- Summary: Create rental contract. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/rentals/contracts/{publicId}/charges`

- Summary: List contract charges. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`

- Summary: Update charge status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `sales`

- Contract: `docs/contracts/http/sales.openapi.yaml`
- Endpoints: `6`

#### `GET /api/sales/opportunities`

- Summary: List opportunities. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/sales/opportunities`

- Summary: Create opportunity. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/sales/proposals`

- Summary: List proposals. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/sales/proposals`

- Summary: Create proposal. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/sales/sales`

- Summary: List sales. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/sales/invoices`

- Summary: List commercial invoices. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `simulation`

- Contract: `docs/contracts/http/simulation.openapi.yaml`
- Endpoints: `3`

#### `GET /api/simulation/scenarios`

- Summary: List scenarios. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/simulation/scenarios`

- Summary: Create scenario run. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/simulation/benchmarks/load`

- Summary: Execute one load benchmark run. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `supplier`

- Contract: `docs/contracts/http/supplier.openapi.yaml`
- Endpoints: `8`

#### `GET /api/supplier/capabilities`

- Summary: Read supplier capability catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/supplier/categories`

- Summary: List supplier categories. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/supplier/categories/{categoryKey}`

- Summary: Upsert one supplier category. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/supplier/suppliers`

- Summary: List suppliers by tenant and status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/supplier/suppliers`

- Summary: Create one supplier. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `201`.

#### `GET /api/supplier/suppliers/summary`

- Summary: Read supplier summary. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/supplier/suppliers/{publicId}`

- Summary: Read one supplier. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PATCH /api/supplier/suppliers/{publicId}`

- Summary: Update one supplier. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `support`

- Contract: `docs/contracts/http/support.openapi.yaml`
- Endpoints: `9`

#### `GET /api/support/capabilities`

- Summary: Read support capability catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/support/queues`

- Summary: List support queues by tenant. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PUT /api/support/queues/{queueKey}`

- Summary: Upsert one support queue. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/support/cases`

- Summary: List support cases with cursor filters. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/support/cases`

- Summary: Create one support case. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `201`.

#### `GET /api/support/cases/summary`

- Summary: Read support case summary. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/support/cases/{publicId}`

- Summary: Read one support case. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PATCH /api/support/cases/{publicId}/status`

- Summary: Transition support case status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/support/cases/{publicId}/comments`

- Summary: Append comment to support case. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `webhook-hub`

- Contract: `docs/contracts/http/webhook-hub.openapi.yaml`
- Endpoints: `13`

#### `GET /health/details`

- Summary: Return readiness details for webhook runtime. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/webhook-hub/capabilities`

- Summary: Read outbound webhook capability posture. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/webhook-hub/outbound-endpoints`

- Summary: List tenant outbound endpoints. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/webhook-hub/outbound-endpoints`

- Summary: Register one tenant outbound endpoint. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/webhook-hub/outbound-endpoints/{publicId}`

- Summary: Read one tenant outbound endpoint. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: List outbound delivery log for one endpoint. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: Register one outbound delivery attempt. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter`

- Summary: Move one outbound delivery to dead letter. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `GET /api/webhook-hub/events`

- Summary: List inbound webhook events. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/webhook-hub/events`

- Summary: Register inbound webhook event. 
- Parameters: nenhum parametro declarado.
- Request body: yes.
- Responses: nenhuma resposta declarada.

#### `GET /api/webhook-hub/events/summary`

- Summary: Aggregate inbound webhook state. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/webhook-hub/events/{publicId}/dead-letter`

- Summary: Move event to dead letter queue. 
- Parameters: `publicId`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

#### `POST /api/webhook-hub/events/{publicId}/requeue`

- Summary: Requeue dead-letter event. 
- Parameters: `publicId`.
- Request body: not declared.
- Responses: nenhuma resposta declarada.

### `workflow-control`

- Contract: `docs/contracts/http/workflow-control.openapi.yaml`
- Endpoints: `7`

#### `GET /api/workflow-control/definitions`

- Summary: List workflow definitions. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/workflow-control/definitions`

- Summary: Create workflow definition. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/workflow-control/definitions/{key}`

- Summary: Read one workflow definition. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PATCH /api/workflow-control/definitions/{key}`

- Summary: Update one workflow definition. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `PATCH /api/workflow-control/definitions/{key}/status`

- Summary: Update workflow definition status. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/workflow-control/capabilities/triggers`

- Summary: List workflow trigger catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/workflow-control/capabilities/actions`

- Summary: List workflow action catalog. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

### `workflow-runtime`

- Contract: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Endpoints: `6`

#### `GET /api/workflow-runtime/executions`

- Summary: List workflow executions. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/workflow-runtime/executions`

- Summary: Create workflow execution. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/workflow-runtime/executions/{publicId}`

- Summary: Read one workflow execution. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/workflow-runtime/executions/{publicId}/actions`

- Summary: List execution action snapshots. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `POST /api/workflow-runtime/executions/{publicId}/advance`

- Summary: Advance one workflow execution. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

#### `GET /api/workflow-runtime/capabilities`

- Summary: List runtime capabilities. 
- Parameters: nenhum parametro declarado.
- Request body: not declared.
- Responses: `200`.

