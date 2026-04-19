# ERP

![ERP](https://img.shields.io/badge/ERP-Enterprise%20Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Polyglot ERP platform for identity, CRM, sales, workflows, analytics and operational automation**

[![Go](https://img.shields.io/badge/Go-edge%20crm%20sales-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![.NET](https://img.shields.io/badge/.NET-identity-512BD4?style=flat&logo=dotnet&logoColor=white)](https://dotnet.microsoft.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-workflow%20control-3178C6?style=flat&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Elixir](https://img.shields.io/badge/Elixir-workflow%20runtime-4B275F?style=flat&logo=elixir&logoColor=white)](https://elixir-lang.org/)
[![Python](https://img.shields.io/badge/Python-analytics-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![Rust](https://img.shields.io/badge/Rust-webhook%20hub-000000?style=flat&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-domain%20storage-316192?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Container%20First-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)
[![Tests](https://img.shields.io/badge/Tests-Unit%20Integration%20Contract%20Smoke-success?style=flat)]()

---

## Documentation / Documentacao

**Project Overview:** [README.md](README.md)  
**Read in English:** [README_EN.md](README_EN.md)  
**Leia em Portugues:** [README_PT.md](README_PT.md)  
**Architecture Reference:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Engineering Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Integration Map:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations Reference:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

---

## What is ERP?

ERP is a portfolio-grade enterprise platform built to feel like a real internal product instead of a generic CRUD sample. The repository is structured as a multi-tenant, container-first and polyglot backend, where each service owns a bounded operational responsibility and a dedicated PostgreSQL context when persistence is required.

This repository is also the public evolution of ERP work already applied in real companies. The goal is not only to build a new system, but to turn that accumulated operational experience into:

- a reusable enterprise template
- a study case for architecture and systems design
- a portfolio project with realistic scope
- a semi-open enterprise ERP reference that can be inspected and extended

The current milestone is already a real backend-first MVP. The platform now exposes an end-to-end commercial and automation slice connecting `identity`, `crm`, `sales`, `workflow-control`, `workflow-runtime`, `analytics`, `webhook-hub` and `edge`.

---

## Current MVP Status

Today the MVP already includes:

- tenant, company, user, team and role management in `identity`
- lead capture, ownership, status progression and note history in `crm`
- opportunity, proposal and sale lifecycle in `sales`
- workflow definition catalog, versioning, runs and operational events in `workflow-control`
- runtime executions, retries and transition ledger in `workflow-runtime`
- aggregated operational reports in `analytics`
- external webhook intake and transition tracking in `webhook-hub`
- operational cockpits in `edge`

The project is still backend-first. That means the strongest part of the MVP today is system design, service behavior, contracts, database ownership and operational visibility rather than frontend presentation.

---

## Why This Project Exists

ERP is designed to serve four roles at the same time:

- enterprise template: a base for future systems and internal accelerators
- study case: a concrete example of multi-tenant, polyglot and container-first architecture
- portfolio: a project that demonstrates depth, consistency and operational realism
- public reference: a way to expose how an enterprise ERP can be structured in a modern modular stack

That is why the repository cares so much about domain boundaries, health endpoints, database ownership, contract stability, smoke validation and service composition.

---

## Architecture at a Glance

The platform is intentionally separated into operational planes:

- transaction plane: `identity`, `crm`, `sales`
- control plane: `workflow-control`
- runtime plane: `workflow-runtime`
- analytics plane: `analytics`
- integration plane: `webhook-hub`
- aggregation and public ops plane: `edge`

This separation keeps catalog concerns, execution concerns, transactional write concerns and reporting concerns from collapsing into the same service.

---

## Stack and Rationale

The stack is polyglot by design, not by aesthetics.

- `Go` powers `edge`, `crm` and `sales` because these services benefit from strong HTTP performance, simple deployment and low operational friction.
- `.NET` powers `identity` because access control, tenancy and future financial domains benefit from strong enterprise ergonomics and mature application structuring.
- `TypeScript` powers `workflow-control` because control-plane APIs and metadata-heavy orchestration benefit from fast iteration and a productive application layer.
- `Elixir` powers `workflow-runtime` because durable execution, retries, timers and high-concurrency orchestration are a natural fit for OTP.
- `Python` powers `analytics` because heavy read models, reporting, ETL-like logic and future forecasting workloads fit well in that ecosystem.
- `Rust` powers `webhook-hub` because external event ingress, idempotency and strict transition control benefit from a highly predictable runtime.
- `PostgreSQL` is used as domain-owned transactional storage, with schemas, migrations and seeds separated by context.

---

## Service Inventory

| Service | Language | Responsibility | Main Public Surface |
|--------|----------|----------------|---------------------|
| `identity` | .NET | tenancy, companies, users, teams, roles and access structure | tenant bootstrap and access snapshot |
| `crm` | Go | leads, ownership, status, summary and notes | lead pipeline and relationship history |
| `sales` | Go | opportunities, proposals, sale conversion and revenue transitions | commercial lifecycle from opportunity to booked sale |
| `workflow-control` | TypeScript | workflow definitions, versions, runs and run events | control-plane APIs and operational event ledger |
| `workflow-runtime` | Elixir | durable execution, lifecycle transitions and retries | runtime orchestration and execution summary |
| `analytics` | Python | heavy operational reads and business reports | sales, tenant and automation reports |
| `webhook-hub` | Rust | inbound webhook intake and transition tracking | event ingestion and delivery traceability |
| `edge` | Go | cross-service aggregation and public operational cockpit | health, tenant, automation and sales overview |

---

## Commercial and Automation Vertical Slice

The clearest way to understand the current MVP is to follow its main flow:

1. A lead enters `crm`.
2. The lead can receive ownership, notes and status progression.
3. A linked opportunity is created in `sales`.
4. A proposal is created for that opportunity.
5. The proposal can move through status transitions and be converted into a sale.
6. `workflow-control` can create a workflow run bound to the same business subject.
7. `workflow-runtime` can execute the corresponding definition with lifecycle transitions and retries.
8. `analytics` reads the resulting footprint and exposes reports such as `pipeline-summary`, `sales-journey`, `tenant-360`, `automation-board` and `workflow-definition-health`.
9. `edge` aggregates these reports into `tenant-overview`, `automation-overview` and `sales-overview`.

This gives the repository a real end-to-end business narrative instead of isolated service demos.

---

## Public API Snapshot

### Identity

Main routes:

- `GET /api/identity/tenants`
- `POST /api/identity/tenants`
- `GET /api/identity/tenants/{slug}`
- `GET /api/identity/tenants/{slug}/snapshot`
- `GET /api/identity/tenants/{slug}/companies`
- `POST /api/identity/tenants/{slug}/companies`
- `PATCH /api/identity/tenants/{slug}/companies/{companyPublicId}`
- `GET /api/identity/tenants/{slug}/users`
- `POST /api/identity/tenants/{slug}/users`
- `PATCH /api/identity/tenants/{slug}/users/{userPublicId}`
- `GET /api/identity/tenants/{slug}/teams`
- `POST /api/identity/tenants/{slug}/teams`
- `PATCH /api/identity/tenants/{slug}/teams/{teamPublicId}`
- `GET /api/identity/tenants/{slug}/roles`

### CRM

Main routes:

- `GET /api/crm/leads`
- `GET /api/crm/leads/summary`
- `POST /api/crm/leads`
- `GET /api/crm/leads/{publicId}`
- `PATCH /api/crm/leads/{publicId}`
- `PATCH /api/crm/leads/{publicId}/owner`
- `PATCH /api/crm/leads/{publicId}/status`
- `GET /api/crm/leads/{publicId}/notes`
- `POST /api/crm/leads/{publicId}/notes`

### Sales

Main routes:

- `GET /api/sales/opportunities`
- `GET /api/sales/opportunities/summary`
- `POST /api/sales/opportunities`
- `GET /api/sales/opportunities/{publicId}`
- `PATCH /api/sales/opportunities/{publicId}`
- `PATCH /api/sales/opportunities/{publicId}/stage`
- `GET /api/sales/opportunities/{publicId}/proposals`
- `POST /api/sales/opportunities/{publicId}/proposals`
- `GET /api/sales/proposals/{publicId}`
- `PATCH /api/sales/proposals/{publicId}/status`
- `POST /api/sales/proposals/{publicId}/convert`
- `GET /api/sales/sales`
- `GET /api/sales/sales/summary`
- `GET /api/sales/sales/{publicId}`
- `PATCH /api/sales/sales/{publicId}/status`

### Workflow Control

Main routes:

- `GET /api/workflow-control/definitions`
- `POST /api/workflow-control/definitions`
- `PATCH /api/workflow-control/definitions/{key}`
- `PATCH /api/workflow-control/definitions/{key}/status`
- `GET /api/workflow-control/definitions/{key}/versions`
- `POST /api/workflow-control/definitions/{key}/versions`
- `POST /api/workflow-control/definitions/{key}/versions/{versionNumber}/restore`
- `GET /api/workflow-control/runs`
- `GET /api/workflow-control/runs/summary`
- `POST /api/workflow-control/runs`
- `GET /api/workflow-control/runs/{publicId}`
- `GET /api/workflow-control/runs/{publicId}/events`
- `GET /api/workflow-control/runs/{publicId}/events/summary`
- `POST /api/workflow-control/runs/{publicId}/events`
- `POST /api/workflow-control/runs/{publicId}/start`
- `POST /api/workflow-control/runs/{publicId}/complete`
- `POST /api/workflow-control/runs/{publicId}/fail`
- `POST /api/workflow-control/runs/{publicId}/cancel`

### Workflow Runtime

Main routes:

- `GET /api/workflow-runtime/executions`
- `GET /api/workflow-runtime/executions/{publicId}`
- `GET /api/workflow-runtime/executions/{publicId}/transitions`
- `GET /api/workflow-runtime/executions/summary`
- `GET /api/workflow-runtime/executions/summary/by-workflow`
- `POST /api/workflow-runtime/executions`
- `POST /api/workflow-runtime/executions/{publicId}/start`
- `POST /api/workflow-runtime/executions/{publicId}/complete`
- `POST /api/workflow-runtime/executions/{publicId}/fail`
- `POST /api/workflow-runtime/executions/{publicId}/cancel`
- `POST /api/workflow-runtime/executions/{publicId}/retry`

### Analytics

Main routes:

- `GET /api/analytics/reports/pipeline-summary`
- `GET /api/analytics/reports/service-pulse`
- `GET /api/analytics/reports/sales-journey`
- `GET /api/analytics/reports/tenant-360`
- `GET /api/analytics/reports/automation-board`
- `GET /api/analytics/reports/workflow-definition-health`
- `GET /api/analytics/reports/delivery-reliability`

### Edge

Main routes:

- `GET /api/edge/ops/health`
- `GET /api/edge/ops/tenant-overview`
- `GET /api/edge/ops/automation-overview`
- `GET /api/edge/ops/sales-overview`

### Webhook Hub

Main routes:

- `GET /api/webhook-hub/events`
- `POST /api/webhook-hub/events`
- `GET /api/webhook-hub/events/{publicId}`
- `GET /api/webhook-hub/events/{publicId}/transitions`
- `GET /api/webhook-hub/events/summary`
- `POST /api/webhook-hub/events/{publicId}/validate`
- `POST /api/webhook-hub/events/{publicId}/queue`
- `POST /api/webhook-hub/events/{publicId}/process`
- `POST /api/webhook-hub/events/{publicId}/forward`
- `POST /api/webhook-hub/events/{publicId}/fail`
- `POST /api/webhook-hub/events/{publicId}/reject`

All services also expose:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`

---

## Local Development

### Option 1: Shell Scripts

```bash
./scripts/build.sh
./scripts/up.sh
./scripts/down.sh
```

### Option 2: Docker Compose

```bash
docker compose --env-file .env.example -f infra/docker-compose.yml up --build -d
```

### Common Utility Commands

```bash
./scripts/logs.sh edge
./scripts/db.sh migrate all
./scripts/db.sh seed all
./scripts/db.sh summary sales bootstrap-ops
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Available Services

- Edge: `http://localhost:8080`
- Identity: `http://localhost:8081`
- Webhook Hub: `http://localhost:8082`
- CRM: `http://localhost:8083`
- Workflow Control: `http://localhost:8084`
- Workflow Runtime: `http://localhost:8085`
- Analytics: `http://localhost:8086`
- Sales: `http://localhost:8087`
- PostgreSQL: `localhost:5432` by default
- Redis: `localhost:6379` by default

If one of these host ports is already in use, the container-first scripts remap it automatically during local runtime and smoke validation.

---

## Validation Strategy

The repository is intentionally container-first. Validation currently covers:

- `unit`: Go, TypeScript, Elixir, Python, .NET and Rust service-level tests
- `integration`: dedicated HTTP integration suite for `identity`
- `contract`: public API contract suites for `workflow-control`, `crm`, `sales` and `identity`
- `smoke`: end-to-end runtime validation in Docker Compose with PostgreSQL and Redis

This means the MVP is not only modeled in code, but also verified close to the real runtime topology.

---

## CI/CD and Container Publishing

The repository already includes:

- `Quality` workflow in GitHub Actions for `unit`, `integration`, `contract` and `smoke`
- `Containers` workflow for publishing service images to `ghcr.io`

Expected published images:

- `ghcr.io/thiagodifaria/erp-edge`
- `ghcr.io/thiagodifaria/erp-crm`
- `ghcr.io/thiagodifaria/erp-sales`
- `ghcr.io/thiagodifaria/erp-identity`
- `ghcr.io/thiagodifaria/erp-workflow-control`
- `ghcr.io/thiagodifaria/erp-workflow-runtime`
- `ghcr.io/thiagodifaria/erp-analytics`
- `ghcr.io/thiagodifaria/erp-webhook-hub`

---

## What Can Be Studied Here

This repository is useful if you want to study:

- multi-tenant service boundaries
- polyglot architecture with explicit role allocation by language
- PostgreSQL ownership by domain context
- control plane versus runtime plane separation
- webhook intake and operational transition tracking
- service aggregation and executive operational read models
- container-first validation in a monorepo
- portfolio-ready enterprise backend presentation

---

## Contact

**Thiago Di Faria** - thiagodifaria@gmail.com

[![GitHub](https://img.shields.io/badge/GitHub-@thiagodifaria-black?style=flat&logo=github)](https://github.com/thiagodifaria)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Thiago_Di_Faria-blue?style=flat&logo=linkedin)](https://linkedin.com/in/thiagodifaria)

---

## Acknowledgments

Special thanks to:

- the Go, .NET, TypeScript, Elixir, Python and Rust ecosystems
- PostgreSQL, Redis and Docker maintainers
- the open-source libraries used across the platform

---

### Star this project if you find it useful

**Made by [Thiago Di Faria](https://github.com/thiagodifaria)**
