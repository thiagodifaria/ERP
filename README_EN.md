# Business Operating System

![Business Operating System](https://img.shields.io/badge/Business%20Operating%20System-Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**A complete operating ecosystem for companies that need sales, customers, billing, finance, fiscal routines, documents, reports, automated processes and external services working as one business platform.**

[![Version](https://img.shields.io/badge/Version-1.5.0-2563EB?style=flat)](docs/CHANGELOG.md)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-contracts-6BA539?style=flat&logo=openapiinitiative&logoColor=white)](docs/contracts/http)
[![Services](https://img.shields.io/badge/Services-26%20HTTP%20APIs-111827?style=flat)](docs/SERVICOS.md)
[![Console](https://img.shields.io/badge/API%20Console-client--api-2563EB?style=flat)](client-web/client-api)
[![Runtime](https://img.shields.io/badge/Runtime-Docker%20Compose-2496ED?style=flat&logo=docker&logoColor=white)](infra)

## Documentation

| File | Purpose |
|------|---------|
| [README.md](README.md) | concise repository overview |
| [README_PT.md](README_PT.md) | detailed Portuguese overview |
| [docs/ARQUITETURA.md](docs/ARQUITETURA.md) | architecture, boundaries, data ownership and runtime topology |
| [docs/API.md](docs/API.md) | HTTP conventions, endpoint index and API usage rules |
| [docs/SERVICOS.md](docs/SERVICOS.md) | service inventory, ownership and implementation paths |
| [docs/CONTRATOS.md](docs/CONTRATOS.md) | OpenAPI, event schemas, registry and compatibility policy |
| [docs/INTEGRACOES.md](docs/INTEGRACOES.md) | providers, webhooks, events and cross context integration |
| [docs/OPERACOES.md](docs/OPERACOES.md) | local runtime, scripts, database, validation and troubleshooting |
| [docs/PADROES.md](docs/PADROES.md) | engineering standards for backend, tests, contracts and docs |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | chronological history by version |

## What Is The Project?

The project is more than a traditional ERP module set. It is a business operating system for running company operations from the first customer contact to billing, finance, documents, reports and audit trails.

The project brings together capabilities that normally live across several tools. It helps manage customers, sales, contracts, subscriptions, invoices, payments, documents, fiscal routines, reports and controlled communication with external services.

The goal is to support complete business journeys. A commercial opportunity can become a proposal, generate a contract, create billing obligations, affect finance, attach documents, trigger workflows, emit webhooks and appear in operational reports.

## Product Scope

This repository models the operating core of a company. It can represent acquisition, sales, recurring contracts, invoicing, receivables, payables, fiscal obligations, procurement, inventory, documents, workflows, support, search, analytics and provider governance.

The system also prepares optional external integrations. Payment gateways, banking providers, AI services, document reading tools, company lookup services, news feeds, signing tools and communication providers can be connected when the user supplies the required credentials. When a credential is missing, the platform should show that clearly instead of pretending the integration is active.

## Core Capabilities

| Capability | Description |
|------------|-------------|
| Commercial operation | Leads, customers, opportunities, proposals, recurring contracts and commercial catalog. |
| Billing and finance | Subscriptions, invoices, payment attempts, money to receive, money to pay, treasury, commissions and reconciliation. |
| Fiscal and banking | Fiscal documents, certificates, SPED posture, Pix, boletos, Open Finance, statements and banking reconciliation. |
| Documents | Attachments, versions, storage posture, signing, document intelligence and audit trails. |
| Workflows | Step by step business processes that the system can execute and track. |
| Reports and analysis | Operational reports, quality indicators, risk views, financial close and platform health. |
| External integrations | Payment gateways, AI tools, document reading, company lookup, market data, news, communication and digital signing. |
| Governance | Users, permissions, limits, activation status, lifecycle and operational evidence. |
| Search and evidence | Operational search, information discovery, legal holds, controlled exports and audit evidence. |

## Main Modules

| Area | Services | Responsibility |
|------|----------|----------------|
| Identity and tenancy | `identity`, `platform-control` | tenants, users, sessions, roles, MFA, capabilities, quotas, lifecycle and go live |
| Commercial operation | `crm`, `sales`, `rentals`, `catalog` | leads, customers, opportunities, proposals, recurring contracts and commercial catalog |
| Billing and finance | `billing`, `finance`, `accounting`, `banking` | subscriptions, invoices, receivables, payables, treasury, commissions, ledger views and reconciliation |
| Fiscal and procurement | `fiscal`, `procurement`, `supplier`, `inventory` | fiscal documents, certificates, SPED posture, requisitions, purchase orders, supplier records and inventory |
| Documents and workflows | `documents`, `workflow-control`, `workflow-runtime` | attachments, storage posture, signing, workflow definitions, executions, retries and compensations |
| Integrations | `webhook-hub`, `engagement`, `notification`, `ai-governance` | webhooks, provider callbacks, touchpoints, notifications, approved AI tools and redaction |
| Intelligence and operations | `analytics`, `search`, `simulation`, `edge`, `support` | reports, operational search, e discovery, scenarios, technical gateway and support operations |

## Business Flows

1. Commercial flow. CRM qualifies demand, Sales manages opportunities and proposals, Billing converts contracts into recurring obligations, and Finance follows projections, receivables and commissions.

2. Financial flow. Billing, Finance, Accounting and Banking connect invoices, payments, statements, reconciliation, ledger views and financial close.

3. Document flow. Documents manages attachments, versions, signing and document intelligence, while Fiscal and Procurement use those records for formal obligations and evidence.

4. Automation flow. Workflow Control defines processes, Workflow Runtime executes steps, and domain services expose the resources used in each journey.

5. Integration flow. External providers enter through adapters, webhooks or specific APIs. The platform keeps BYOK posture and shows whether each capability is configured, manual, in fallback or unavailable.

6. Governance flow. Platform Control, Analytics and Edge consolidate tenants, quotas, providers, readiness, go live, incidents, risks and operational evidence.

## Technical API Console

`client-web/client-api` is a technical control console for the API. It is closer to a modern, project specific Swagger UI than to the future business application. Its role is to document, inspect and exercise the platform from one place.

The console provides a platform overview, endpoint catalog generated from OpenAPI contracts, request builder, local environments, documentation reader, contract views, test journeys and operational screens.

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

## Local Runtime

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh ps
./scripts/build.sh logs edge
./scripts/build.sh down
```

Database:

```bash
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh psql
./scripts/build.sh backup /tmp/erp-local-backup.sql
./scripts/build.sh restore /tmp/erp-local-backup.sql
```

## Contracts And Governance

Contracts are versioned engineering artifacts. They define the public API surface, shared events, compatibility rules and the source used by the technical console.

```text
docs/contracts/http/              OpenAPI files
docs/contracts/events/            JSON Schema event contracts
docs/contracts/registry.json      contract registry
docs/contracts/schema-registry.json
docs/contracts/portal/index.html
```

Before changing public API shape, run:

```bash
./scripts/test.sh contract
```

## Validation

| Scope | Command |
|-------|---------|
| unit tests | `./scripts/test.sh unit` |
| integration tests | `./scripts/test.sh integration` |
| HTTP and event contracts | `./scripts/test.sh contract` |
| platform checks | `./scripts/test.sh platform` |
| smoke tests | `./scripts/test.sh smoke` |
| performance | `./scripts/test.sh performance` |
| backup and restore | `./scripts/test.sh backup-restore` |
| hardening | `./scripts/test.sh hardening` |
| production readiness | `./scripts/test.sh production-readiness` |

API console:

```bash
cd client-web/client-api
npm run generate
npm run typecheck
npm run build
```

## Repository Layout

```text
client-web/client-api/     technical API console
docs/                      project documentation
docs/contracts/            OpenAPI, events, registry and portal
infra/                     Docker Compose and Kubernetes
scripts/                   runtime, build and validation entrypoints
service-api/               backend services and PostgreSQL contexts
```

## Private Ownership

This repository is privately maintained. Code changes are controlled directly by the maintainer.

## Contact

**Thiago Di Faria**  
thiagodifaria@gmail.com

[GitHub](https://github.com/thiagodifaria)  
[LinkedIn](https://linkedin.com/in/thiagodifaria)
