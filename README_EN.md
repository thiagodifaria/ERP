# ERP

ERP is a backend-first, multi-tenant and polyglot ERP platform. It is organized around bounded services, explicit contracts, local container execution and automated validation.

This README gives a complete but compact view of the project. Topic-specific details live in `docs/`.

## Documentation Map

| File | Purpose |
|------|---------|
| `docs/ARQUITETURA.md` | Architecture, boundaries, data ownership and runtime topology |
| `docs/API.md` | HTTP conventions, endpoint index and API usage rules |
| `docs/SERVICOS.md` | Service inventory, ownership and implementation paths |
| `docs/CONTRATOS.md` | OpenAPI, event schemas, registry and compatibility policy |
| `docs/INTEGRACOES.md` | Providers, webhooks, events and cross-context integration |
| `docs/OPERACOES.md` | Local runtime, scripts, database, validation and troubleshooting |
| `docs/PADROES.md` | Engineering standards for backend, tests and documentation |
| `docs/CHANGELOG.md` | Chronological change history |

## Project Shape

- 24 HTTP services with OpenAPI contracts.
- 542 versioned HTTP endpoints.
- 15 versioned event schemas.
- Contracts under `docs/contracts/`.
- Container runtime through `infra/docker-compose.yml`.
- Operational entrypoint through `./scripts/build.sh`.
- Validation entrypoint through `./scripts/test.sh`.
- Technical API console under `client-web/client-api`.

## Architecture Summary

The system is split into service ownership boundaries instead of being a single modular monolith. Each service owns its public API contract, runtime implementation and persistence context where applicable.

Main architectural rules:

- tenant context must be explicit in tenant-aware operations;
- public contracts live in `docs/contracts/`;
- implementation lives under `service-api/`;
- runtime infrastructure lives under `infra/`;
- database ownership is split by PostgreSQL schema/context;
- cross-service aggregation belongs to `analytics` or `edge`;
- external callbacks and webhook delivery belong to `webhook-hub` or a provider-facing adapter;
- long-running operations should expose a job, rollout or execution resource instead of blocking the request.

## Runtime

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh ps
./scripts/build.sh logs edge
./scripts/build.sh down
```

Database operations:

```bash
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh backup /tmp/erp-local-backup.sql
./scripts/build.sh restore /tmp/erp-local-backup.sql
./scripts/build.sh psql
```

Validation:

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

## API Console

`client-web/client-api` is the technical console for the backend API. It is separate from any future business-facing frontend.

It provides:

- API overview;
- generated endpoint catalog from OpenAPI files;
- request builder with headers, params and JSON body;
- documentation reader backed by the markdown files in this repository;
- contract, environment, journey and operations views.

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

The console uses the local Vite proxy to call backend services when the stack is running.

## Service Inventory

| Service | Stack | Path | Responsibility |
|---------|-------|------|----------------|
| `accounting` | Python | `service-api/service-python/accounting` | management accounting, cost centers, posting rules, ledger, statements and close |
| `analytics` | Python | `service-api/service-python/analytics` | executive reports, contract governance, readiness and operational views |
| `banking` | Python | `service-api/service-python/banking` | CNAB, boletos, bank statements, reconciliation, Pix charges/refunds/webhooks and Open Finance |
| `billing` | .NET | `service-api/service-csharp/billing` | plans, subscriptions, invoices, usage pricing and payment attempts |
| `catalog` | Python | `service-api/service-python/catalog` | categories, items, version history, bulk creation and consumer contracts |
| `crm` | Go | `service-api/service-golang/crm` | leads, customers, pipeline configuration and CNPJ enrichment |
| `documents` | Go | `service-api/service-golang/documents` | attachments, storage capabilities, signing and version history |
| `edge` | Go | `service-api/service-golang/edge` | public operational cockpit and cross-service reads |
| `engagement` | TypeScript | `service-api/service-typescript/engagement` | provider capabilities, inbound events, touchpoints and conversations |
| `finance` | .NET | `service-api/service-csharp/finance` | receivables, projections, commission holds and financial activity |
| `fiscal` | Python | `service-api/service-python/fiscal` | fiscal profiles, issuance, certificates, contingency, SPED, tax documents, privacy and audit |
| `identity` | .NET | `service-api/service-csharp/identity` | tenants, users, sessions, invitations, roles and MFA posture |
| `inventory` | Python | `service-api/service-python/inventory` | location balances, stock movements, reservations, FIFO/average costing and cycle counts |
| `notification` | Python | `service-api/service-python/notification` | preferences, notification center and delivery state |
| `platform-control` | Python | `service-api/service-python/platform-control` | capabilities, providers, entitlements, quotas, lifecycle and go-live |
| `procurement` | Python | `service-api/service-python/procurement` | requisitions, quotations, purchase orders, approvals, receiving and 3-way matching |
| `rentals` | Go | `service-api/service-golang/rentals` | recurring contracts and charge lifecycle |
| `sales` | Go | `service-api/service-golang/sales` | opportunities, proposals, sales, invoices and commissions |
| `simulation` | Python | `service-api/service-python/simulation` | scenario runs and load benchmarks |
| `supplier` | Python | `service-api/service-python/supplier` | supplier categories, supplier directory and summary |
| `support` | Python | `service-api/service-python/support` | queues, cases, SLA, comments and support summary |
| `webhook-hub` | Rust | `service-api/service-rust/webhook-hub` | inbound webhooks, outbound endpoints, delivery log and DLQ |
| `workflow-control` | TypeScript | `service-api/service-typescript/workflow-control` | workflow definitions, catalogs and control-plane state |
| `workflow-runtime` | Elixir | `service-api/service-elixir/workflow-runtime` | executions, actions, transitions, retries and compensations |

## Contract Catalog

Contracts are source-controlled engineering artifacts:

```text
docs/contracts/http/              OpenAPI files
docs/contracts/events/            JSON Schema event contracts
docs/contracts/registry.json      contract registry
docs/contracts/schema-registry.json
docs/contracts/portal/index.html
```

Use `./scripts/test.sh contract` before changing public API shape.

## Development Rules

- Keep service behavior inside the owning service.
- Update OpenAPI when changing public HTTP shape.
- Update event schemas when publishing or consuming shared events.
- Keep docs scoped: architecture in architecture docs, API rules in API docs, operations in operations docs.
- Prefer small, explicit runtime flows over hidden coupling.
- Add tests proportional to the blast radius.

## Main Flows

Commercial flow:

`crm` creates qualified demand, `sales` manages proposals and sales, `billing` creates recurring commercial obligations, and `finance` consolidates projections, holds and activity.

Recurring contract flow:

`rentals` and `billing` represent recurring obligations while `finance` reads the financial consequences.

Automation flow:

`workflow-control` defines automation, `workflow-runtime` executes it, and domain services expose the resources that workflows operate on.

Integration flow:

Provider callbacks enter through provider-specific endpoints or `webhook-hub`, are normalized, and can be observed through `analytics` and `edge`.

Go-live flow:

`platform-control` owns lifecycle, quotas, blocks and rollout posture; `analytics` and `edge` expose executive views.

## Repository Layout

```text
client-web/client-api/     technical API console
docs/                      documentation
docs/contracts/            contracts and schemas
infra/                     Docker Compose and runtime assets
scripts/                   build/runtime and validation entrypoints
service-api/               backend services and PostgreSQL contexts
```

## What This Repository Is

This repository is the backend and technical platform layer of the ERP. It is not a marketing website, not a single CRUD API and not a frontend-first application. The main product surface today is the API, its contracts, its service boundaries and the local runtime used to validate the platform.

The project is intentionally backend-heavy because the hard part being modeled is not the screen layout. The hard part is keeping many business contexts consistent enough to evolve:

- identity and tenancy;
- commercial operation;
- recurring contracts and billing;
- finance and commissions;
- document metadata and signing;
- fiscal, privacy and audit;
- workflow definition and execution;
- provider callbacks and webhooks;
- platform entitlements and go-live;
- analytics and operational control.

## What This Repository Is Not Yet

The existence of a service or endpoint does not mean every production concern is fully closed. Some areas are intentionally structured before being connected to real external providers.

Examples of remaining production concerns:

- stronger authentication and authorization on every public route;
- real external payment, fiscal, communication and signing providers;
- distributed tracing across all services;
- stricter secret handling for every integration path;
- production deployment manifests beyond local Docker Compose;
- business-facing frontend separated from the technical API console.

This distinction matters: the platform is advanced structurally, but documentation should not pretend that local fallback behavior is the same as production integration.

## Quality Gates

The project has several validation levels. Use the smallest suite that proves the change, then broaden when the blast radius grows.

| Change | Recommended validation |
|--------|------------------------|
| local domain rule | `./scripts/test.sh unit` |
| public HTTP shape | `./scripts/test.sh contract` |
| cross-service behavior | `./scripts/test.sh smoke` |
| provider/readiness/go-live posture | `./scripts/test.sh hardening` |
| infrastructure/runtime change | `./scripts/test.sh platform` |
| database migration risk | `./scripts/test.sh backup-restore` |

For the API console:

```bash
cd client-web/client-api
npm run generate
npm run typecheck
npm run build
```

## Documentation Principles

The documentation is split by responsibility:

- architecture explains boundaries and decisions;
- API explains HTTP usage and conventions;
- services explain ownership;
- contracts explain compatibility;
- integrations explain cross-context and provider flows;
- operations explain how to run and diagnose;
- standards explain engineering rules.

Avoid adding the same endpoint catalog to every file. If a change affects an endpoint, update the OpenAPI and the focused documentation that explains why the endpoint exists.

## Typical Contributor Flow

1. Identify the owning service.
2. Check its OpenAPI contract.
3. Change implementation and tests together.
4. Update migrations/seeds if persistence changed.
5. Regenerate the API console catalog if HTTP contracts changed.
6. Run the validation suite that matches the risk.
7. Update the correct documentation file.
8. Add changelog entries only when the change is ready to be recorded as progress.

## Maintainer

Thiago Di Faria - thiagodifaria@gmail.com
