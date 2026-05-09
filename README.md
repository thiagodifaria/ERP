# ERP

Backend-first ERP platform for tenancy, CRM, sales, finance, billing, documents, workflows, analytics, integrations and operational governance.

This repository is a polyglot, container-first backend platform. The root README is intentionally short; detailed reference lives in the language-specific READMEs and in `docs/`.

## Documentation

- [README_EN.md](README_EN.md): detailed English overview.
- [README_PT.md](README_PT.md): visao detalhada em Portugues.
- [docs/ARQUITETURA.md](docs/ARQUITETURA.md): architecture, boundaries and runtime topology.
- [docs/API.md](docs/API.md): HTTP API conventions and endpoint index.
- [docs/SERVICOS.md](docs/SERVICOS.md): service ownership and implementation map.
- [docs/CONTRATOS.md](docs/CONTRATOS.md): OpenAPI, event schemas and compatibility rules.
- [docs/INTEGRACOES.md](docs/INTEGRACOES.md): providers, webhooks, events and cross-context integration.
- [docs/OPERACOES.md](docs/OPERACOES.md): local runtime, validation, database and runbooks.
- [docs/PADROES.md](docs/PADROES.md): engineering standards.
- [docs/CHANGELOG.md](docs/CHANGELOG.md): chronological project history.

## Current Baseline

| Item | Value |
|------|-------|
| HTTP services with OpenAPI contracts | 20 |
| Versioned HTTP endpoints | 201 |
| Event schemas | 12 |
| Contract catalog | `docs/contracts/` |
| Runtime command | `./scripts/build.sh` |
| Validation command | `./scripts/test.sh` |
| API console | `client-web/client-api` |

## Quick Start

```bash
./scripts/build.sh
./scripts/test.sh contract
```

Common local operations:

```bash
./scripts/build.sh up
./scripts/build.sh ps
./scripts/build.sh logs edge
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh down
```

API console:

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

## Repository Layout

```text
client-web/client-api/     technical web console for API exploration
docs/                      project documentation
docs/contracts/            OpenAPI, event schemas, registry and portal
infra/                     Docker Compose and runtime infrastructure
scripts/build.sh           build, runtime, database, backup and restore
scripts/test.sh            validation suites
service-api/               backend services and PostgreSQL contexts
```

## Service Map

| Service | Stack | Main responsibility |
|---------|-------|---------------------|
| `analytics` | Python | executive reports, governance and operational reads |
| `billing` | .NET | plans, subscriptions, invoices and payment attempts |
| `catalog` | Python | categories, items, versions and consumer contracts |
| `crm` | Go | leads, customers, pipeline and enrichment |
| `documents` | Go | attachments, storage posture, versions and signing |
| `edge` | Go | public entrypoint and cross-service cockpits |
| `engagement` | TypeScript | campaigns, touchpoints, conversations and callbacks |
| `finance` | .NET | receivables, payables, treasury and commissions |
| `fiscal` | Python | fiscal documents, retention, privacy and audit |
| `identity` | .NET | tenants, users, roles, sessions, invites and MFA |
| `notification` | Python | preferences and internal notification center |
| `platform-control` | Python | capabilities, quotas, lifecycle and go-live |
| `rentals` | Go | recurring contracts and rental charges |
| `sales` | Go | opportunities, proposals, sales and invoices |
| `simulation` | Python | scenarios and load benchmarks |
| `supplier` | Python | supplier categories and directory |
| `support` | Python | queues, cases, SLA and comments |
| `webhook-hub` | Rust | inbound/outbound webhooks, idempotency and DLQ |
| `workflow-control` | TypeScript | workflow definitions and catalogs |
| `workflow-runtime` | Elixir | durable executions, actions and transitions |

## Validation

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
./scripts/test.sh hardening
```

Use the detailed READMEs for context and the files in `docs/` for focused reference.

---

**Thiago Di Faria** - thiagodifaria@gmail.com

[GitHub](https://github.com/thiagodifaria) · [LinkedIn](https://linkedin.com/in/thiagodifaria)
