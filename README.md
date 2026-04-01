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
[![Tests](https://img.shields.io/badge/Tests-Unit%20Contract%20Smoke-success?style=flat)]()

---

## Documentation / Documentacao

**Detailed English README:** [README_EN.md](README_EN.md)  
**README detalhado em Portugues:** [README_PT.md](README_PT.md)  
**Architecture Reference:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Engineering Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Integration Map:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations Reference:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

---

## What is ERP?

ERP is a portfolio-grade enterprise platform built to feel like a serious internal product instead of a generic CRUD sample. The repository is organized as a multi-tenant, container-first and polyglot backend where each service owns a clear operational responsibility.

The project is also the public evolution of ERP structures already used in real business environments. Here, that experience is being turned into a reusable template, a study case, a portfolio centerpiece and a public enterprise ERP reference.

The current milestone is already a real MVP. The platform now has an end-to-end commercial and automation slice connecting `identity`, `crm`, `sales`, `workflow-control`, `workflow-runtime`, `analytics`, `webhook-hub` and `edge` with PostgreSQL-backed flows and container-first validation.

### Key Highlights

- multi-tenant architecture from day one
- evolution of real ERP implementation experience into a reusable platform reference
- identity, CRM, sales and automation already connected in the same operational flow
- bounded PostgreSQL schemas with migrations and seeds per context
- `edge` exposing platform, automation and sales cockpits
- `analytics` exposing pipeline, sales and tenant-wide operational reports
- unit, integration, contract and smoke validation running in containers
- GitHub Actions and GitHub Container Registry configured for CI/CD

### What Makes It Special?

```text
Go, .NET, TypeScript, Elixir, Python and Rust mapped to clear operational roles
Monorepo organized by language and bounded service ownership
Domain-owned PostgreSQL schemas, migrations and bootstrap seeds
Container-first validation across unit, integration, contract and smoke layers
Commercial flow already connected from lead capture to sale conversion
Operational dashboards exposed through analytics and edge aggregation
CI/CD workflows prepared for GitHub Actions and GHCR publishing
```

### Project Positioning

- template for future enterprise systems and internal platform accelerators
- study case for multi-tenant, polyglot and container-first backend architecture
- portfolio project with real operational scope instead of isolated code samples
- public enterprise ERP reference that can evolve in an open and inspectable way

---

## Quick Start

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

### Validation

```bash
./scripts/test.sh unit
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
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

---

## MVP Scope

The current MVP is backend-first and already covers a complete vertical slice:

- lead capture, owner assignment, status updates and note history in `crm`
- opportunity pipeline, proposal lifecycle, sale conversion and revenue transitions in `sales`
- workflow catalog, publication, runs and operational events in `workflow-control`
- execution lifecycle, retries and transitions in `workflow-runtime`
- sales, tenant and automation reports in `analytics`
- aggregated operational cockpits in `edge`
- webhook intake and delivery transition tracking in `webhook-hub`

For the complete service-by-service presentation, use the detailed guides:

- [README_EN.md](README_EN.md)
- [README_PT.md](README_PT.md)

---

## API Overview

| Area | Endpoint | Description |
|------|----------|-------------|
| Identity | `GET /api/identity/tenants/{slug}/snapshot` | Read the tenant operational snapshot |
| CRM | `POST /api/crm/leads` | Create a lead in the commercial funnel |
| CRM | `POST /api/crm/leads/{publicId}/notes` | Register relationship notes for a lead |
| Sales | `POST /api/sales/opportunities` | Create an opportunity linked to a lead |
| Sales | `POST /api/sales/opportunities/{publicId}/proposals` | Create a proposal for an opportunity |
| Sales | `POST /api/sales/proposals/{publicId}/convert` | Convert an accepted proposal into a sale |
| Workflow Control | `POST /api/workflow-control/runs` | Create a workflow run for a business subject |
| Workflow Runtime | `POST /api/workflow-runtime/executions` | Create a runtime execution from a workflow definition |
| Analytics | `GET /api/analytics/reports/sales-journey` | Read the commercial funnel from lead to sale |
| Analytics | `GET /api/analytics/reports/tenant-360` | Read a consolidated tenant snapshot |
| Edge | `GET /api/edge/ops/automation-overview` | Read the operational automation cockpit |
| Edge | `GET /api/edge/ops/sales-overview` | Read the commercial cockpit |
| Webhook Hub | `POST /api/webhook-hub/events` | Ingest external webhook events |
| Ops | `GET /health/live` | Liveness endpoint |
| Ops | `GET /health/ready` | Readiness endpoint |

More details are available in the service references:

- [service-api/service-csharp/identity/README.md](service-api/service-csharp/identity/README.md)
- [service-api/service-golang/crm/README.md](service-api/service-golang/crm/README.md)
- [service-api/service-golang/sales/README.md](service-api/service-golang/sales/README.md)
- [service-api/service-typescript/workflow-control/README.md](service-api/service-typescript/workflow-control/README.md)
- [service-api/service-elixir/workflow-runtime/README.md](service-api/service-elixir/workflow-runtime/README.md)
- [service-api/service-python/analytics/README.md](service-api/service-python/analytics/README.md)
- [service-api/service-rust/webhook-hub/README.md](service-api/service-rust/webhook-hub/README.md)
- [service-api/service-golang/edge/README.md](service-api/service-golang/edge/README.md)

---

## Automation Scripts

### `scripts/build.sh`

- builds the local Docker images defined in Compose

### `scripts/up.sh`

- starts the local stack in detached mode with build

### `scripts/down.sh`

- stops the stack and removes orphan containers

### `scripts/logs.sh <service>`

- tails service logs from the Compose stack

### `scripts/db.sh`

- applies migrations, seeds and relational summaries by bounded context

### `scripts/test.sh`

- executes container-first `unit`, `integration`, `contract` and `smoke`

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
