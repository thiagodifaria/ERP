# Business Operating System

![Business Operating System](https://img.shields.io/badge/Business%20Operating%20System-Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**A complete operating ecosystem for companies that need sales, customers, billing, finance, fiscal routines, documents, reports, automated processes and external services working as one business platform.**

[![Version](https://img.shields.io/badge/Version-1.4.6-2563EB?style=flat)](docs/CHANGELOG.md)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-contracts-6BA539?style=flat&logo=openapiinitiative&logoColor=white)](docs/contracts/http)
[![Services](https://img.shields.io/badge/Services-26%20HTTP%20APIs-111827?style=flat)](docs/SERVICOS.md)
[![Console](https://img.shields.io/badge/API%20Console-client--api-2563EB?style=flat)](client-web/client-api)
[![Runtime](https://img.shields.io/badge/Runtime-Docker%20Compose-2496ED?style=flat&logo=docker&logoColor=white)](infra)

## Documentation

**Leia em português:** [README_PT.md](README_PT.md)  
**Read the detailed English README:** [README_EN.md](README_EN.md)  
**Architecture:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**API:** [docs/API.md](docs/API.md)  
**Services:** [docs/SERVICOS.md](docs/SERVICOS.md)  
**Contracts:** [docs/CONTRATOS.md](docs/CONTRATOS.md)  
**Integrations:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

## What Is The Project?

The project is more than a traditional ERP module set. It is a business operating system for running company operations from the first customer contact to billing, finance, documents, reports and audit trails.

The project brings together capabilities that normally live across several tools. It helps manage customers, sales, contracts, subscriptions, invoices, payments, documents, fiscal routines, reports and controlled communication with external services.

The goal is to support complete business journeys. A commercial opportunity can become a proposal, generate a contract, create billing obligations, affect finance, attach documents, trigger workflows, emit webhooks and appear in operational reports.

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

## Validation

| Scope | Command |
|-------|---------|
| Unit tests | `./scripts/test.sh unit` |
| Integration tests | `./scripts/test.sh integration` |
| HTTP and event contracts | `./scripts/test.sh contract` |
| Platform checks | `./scripts/test.sh platform` |
| Smoke tests | `./scripts/test.sh smoke` |
| Performance checks | `./scripts/test.sh performance` |
| Backup and restore | `./scripts/test.sh backup-restore` |
| Hardening checks | `./scripts/test.sh hardening` |
| Production readiness | `./scripts/test.sh production-readiness` |

## Repository Layout

```text
client-web/client-api/     technical API console
docs/                      project documentation
docs/contracts/            OpenAPI, event schemas, registry and portal
infra/                     Docker Compose and Kubernetes runtime assets
scripts/                   build, runtime and validation entrypoints
service-api/               backend services and PostgreSQL contexts
```

## Private Ownership

This repository is privately maintained. Code changes are controlled directly by the maintainer.

## Contact

**Thiago Di Faria**  
thiagodifaria@gmail.com

[GitHub](https://github.com/thiagodifaria)  
[LinkedIn](https://linkedin.com/in/thiagodifaria)
