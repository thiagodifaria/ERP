# ERP

![ERP](https://img.shields.io/badge/ERP-Enterprise%20Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Backend-first ERP platform for tenancy, CRM, sales, contracts, finance, billing, workflows, analytics, integrations and operational automation.**

[![Go](https://img.shields.io/badge/Go-edge%20crm%20sales%20rentals%20documents-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![.NET](https://img.shields.io/badge/.NET-identity%20finance%20billing-512BD4?style=flat&logo=dotnet&logoColor=white)](https://dotnet.microsoft.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-workflow%20engagement-3178C6?style=flat&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Elixir](https://img.shields.io/badge/Elixir-workflow%20runtime-4B275F?style=flat&logo=elixir&logoColor=white)](https://elixir-lang.org/)
[![Python](https://img.shields.io/badge/Python-analytics%20platform%20admin-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![Rust](https://img.shields.io/badge/Rust-webhook%20hub-000000?style=flat&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-domain%20schemas-316192?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-container%20first-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)

---

## Documentation / Documentacao

**Detailed English README:** [README_EN.md](README_EN.md)  
**README detalhado em Portugues:** [README_PT.md](README_PT.md)  
**API Reference:** [docs/API.md](docs/API.md)  
**Architecture:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Contracts:** [docs/CONTRATOS.md](docs/CONTRATOS.md)  
**Integrations:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Services:** [docs/SERVICOS.md](docs/SERVICOS.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

---

## What is ERP?

ERP is a multi-tenant, polyglot and container-first backend platform. It is organized as a serious enterprise system instead of a set of isolated CRUD samples. The current backend covers commercial operation, recurring contracts, documents, finance, billing, workflow control, workflow runtime, engagement, analytics, simulation, platform governance, support, suppliers, notifications, fiscal/compliance and webhook operations.

## Current Shape

| Area | Current baseline |
|------|------------------|
| Services with OpenAPI contracts | 20 |
| Versioned HTTP endpoints | 201 |
| Contract catalog | `docs/contracts/` |
| Runtime command | `./scripts/build.sh` |
| Validation command | `./scripts/test.sh` |

## Quick Start

```bash
./scripts/build.sh
./scripts/test.sh contract
```

Common operations:

```bash
./scripts/build.sh up
./scripts/build.sh logs edge
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh down
```

## Repository Layout

```text
docs/                    central documentation and contracts
docs/contracts/          OpenAPI, event schemas, registry and portal
infra/                   compose, observability and runtime assets
scripts/build.sh         build, runtime, logs, database, backup and restore
scripts/test.sh          unit, integration, contract, platform, smoke and hardening suites
service-api/             backend services and PostgreSQL ownership
```

## Service Map

| Service | Stack | Responsibility |
|---------|-------|----------------|
| `analytics` | Python | operational reports, governance, reliability, hardening, cost and executive reads |
| `billing` | .NET | plans, subscriptions, recurring invoices, payment attempts and recovery |
| `catalog` | Python | categories, items, item versions, bulk creation and consumer contracts |
| `crm` | Go | leads, customers, ownership, pipeline, notes, history, attachments and enrichment |
| `documents` | Go | attachments, upload, storage posture, signing, versions, archive and access links |
| `edge` | Go | public entrypoint, cross-service aggregation and operational cockpits |
| `engagement` | TypeScript | campaigns, templates, touchpoints, conversations, delivery, providers and callbacks |
| `finance` | .NET | receivables, payables, treasury, costs, commissions and closing |
| `fiscal` | Python | fiscal profiles, retention, tax documents, privacy, consent and audit |
| `identity` | .NET | tenants, companies, users, teams, roles, sessions, invites, MFA and audit |
| `notification` | Python | preferences, internal alert center, severity and notification lifecycle |
| `platform-control` | Python | capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle and go-live |
| `rentals` | Go | recurring contracts, adjustments, terminations, charges and contractual attachments |
| `sales` | Go | opportunities, proposals, sales, invoices, commissions, renegotiations and pending items |
| `simulation` | Python | operational scenarios, load benchmarks and capacity modeling |
| `supplier` | Python | supplier categories, supplier directory and procurement ownership |
| `support` | Python | queues, cases, SLA, comments, bulk operations and support summaries |
| `webhook-hub` | Rust | webhook intake, idempotency, transitions, DLQ, outbound endpoints and deliveries |
| `workflow-control` | TypeScript | definitions, published versions, trigger/action catalogs, runs and events |
| `workflow-runtime` | Elixir | durable executions, timeline, actions, transitions, retries, waits and compensations |

---

**Thiago Di Faria** - thiagodifaria@gmail.com

[![GitHub](https://img.shields.io/badge/GitHub-@thiagodifaria-black?style=flat&logo=github)](https://github.com/thiagodifaria)  
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Thiago_Di_Faria-blue?style=flat&logo=linkedin)](https://linkedin.com/in/thiagodifaria)
