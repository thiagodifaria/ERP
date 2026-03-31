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

## Documentation

- Architecture: [docs/ARQUITETURA.md](docs/ARQUITETURA.md)
- Standards: [docs/PADROES.md](docs/PADROES.md)
- Integrations: [docs/INTEGRACOES.md](docs/INTEGRACOES.md)
- Operations: [docs/OPERACOES.md](docs/OPERACOES.md)
- Delivery phases: [docs/FASES.md](docs/FASES.md)
- Changelog: [docs/CHANGELOG.md](docs/CHANGELOG.md)

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

## Current Status

The repository is in **Phase 0 - Architectural Foundation**.

This update establishes:

- the root monorepo layout
- the backend split by language under `service-api`
- the committed documentation baseline under `docs`
- the initial versioning policy through `docs/CHANGELOG.md`
- the local-only workspace rules for private planning and progress artifacts

---

## Repository Layout

```text
ERP/
  README.md
  .env.example
  .gitignore
  .dockerignore
  .editorconfig

  service-api/
    service-golang/
    service-csharp/
    service-elixir/
    service-typescript/
    service-python/
    service-rust/
    service-postgresql/

  infra/
  docs/
  scripts/
  tests/
```

---

## Delivery Principles

- keep the repository root clean
- keep code in English
- keep comments short and in Brazilian Portuguese when code starts landing
- keep business rules inside domain and application, never in bootstrap
- keep each domain responsible for its own data and migrations
- prefer container-first execution for local development and tests
- document every meaningful update in `docs/CHANGELOG.md`
