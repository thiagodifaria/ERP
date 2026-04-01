# ERP

![ERP](https://img.shields.io/badge/ERP-Enterprise%20Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Polyglot, multi-tenant and domain-driven ERP platform for serious enterprise operations**

[![Go](https://img.shields.io/badge/Go-edge%20services-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![.NET](https://img.shields.io/badge/.NET-identity%20finance%20billing-512BD4?style=flat&logo=dotnet&logoColor=white)](https://dotnet.microsoft.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-control%20plane-3178C6?style=flat&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Elixir](https://img.shields.io/badge/Elixir-workflow%20runtime-4B275F?style=flat&logo=elixir&logoColor=white)](https://elixir-lang.org/)
[![Python](https://img.shields.io/badge/Python-analytics%20simulation-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![Rust](https://img.shields.io/badge/Rust-webhook%20ingress-000000?style=flat&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-domain%20data-316192?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)

---

## Documentation / Documentacao

**Architecture Reference:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Engineering Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Integration Map:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations Reference:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

---

## What Is ERP?

ERP is being built as a modular enterprise platform, not as a generic backoffice CRUD application. The target is a cloud-native product with strong domain boundaries, multi-tenant ownership, high-concurrency services, auditability, observability and room for long-term evolution.

The platform is intentionally polyglot:

- Go for edge and high-concurrency transactional services
- .NET for identity, finance and billing
- TypeScript for workflow control and engagement integrations
- Elixir for workflow runtime orchestration
- Python for analytics and simulation
- Rust for critical webhook ingress
- PostgreSQL for domain-owned transactional storage

---

## Key Highlights

- multi-tenant architecture from day one
- polyglot backend split by service responsibility
- strong domain ownership across transactional services and PostgreSQL structures
- workflow separation between control plane and runtime plane
- external integrations isolated through adapters
- auditability, observability and idempotent event processing as core platform concerns

### What Makes It Special?

```text
OK Go, .NET, TypeScript, Elixir, Python and Rust mapped to clear operational roles
OK service-api organized first by language instead of mixed deployment folders
OK domain-owned PostgreSQL layout ready for migrations, seeds, views, functions and indexes
OK edge, identity and webhook-hub already started as concrete service templates
OK event-driven architecture planned with Kafka, retries and transactional outbox
OK enterprise observability stack designed from the foundation instead of late-stage retrofit
```

---

## Current Status

The repository is no longer just in the foundation stage. The platform already has an operational backend MVP slice running in containers with real PostgreSQL-backed flows.

Current active services:

- `identity` for tenants, companies, users, teams and roles
- `crm` for leads, notes, ownership and qualification flow
- `sales` for opportunities, proposals, conversions and revenue transitions
- `workflow-control` for workflow catalog, versions, runs and operational events
- `workflow-runtime` for execution lifecycle, retries and runtime transitions
- `analytics` for pipeline, sales, tenant and automation reports
- `webhook-hub` for controlled webhook intake and transition tracking
- `edge` for operational cockpits and cross-service aggregation

Current MVP vertical slice:

- lead enters through `crm`
- opportunity and proposal move through `sales`
- conversion creates a sale and booked revenue
- workflow run and runtime execution reflect the business subject
- analytics exposes `pipeline-summary`, `sales-journey`, `tenant-360` and automation reports
- edge exposes `tenant-overview`, `automation-overview` and `sales-overview`

---

## Contact

**Thiago Di Faria** - thiagodifaria@gmail.com

[![GitHub](https://img.shields.io/badge/GitHub-@thiagodifaria-black?style=flat&logo=github)](https://github.com/thiagodifaria)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Thiago_Di_Faria-blue?style=flat&logo=linkedin)](https://linkedin.com/in/thiagodifaria)

---

### Star this project if you find it useful

**Made by [Thiago Di Faria](https://github.com/thiagodifaria)**
