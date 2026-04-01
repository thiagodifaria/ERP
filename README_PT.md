# ERP

![ERP](https://img.shields.io/badge/ERP-Enterprise%20Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Plataforma ERP poliglota para identidade, CRM, vendas, workflows, analytics e automacao operacional**

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

## Documentacao

**Read in English:** [README_EN.md](README_EN.md)  
**Leia em Portugues:** [README_PT.md](README_PT.md)  
**Architecture Reference:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Engineering Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Integration Map:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations Reference:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

---

## O que e o ERP?

ERP e uma plataforma enterprise pensada para parecer um produto interno serio, e nao um CRUD generico de estudo. O repositorio esta organizado como um backend multi-tenant, container-first e poliglota, no qual cada servico tem uma responsabilidade operacional clara.

O marco atual ja e um MVP real. A plataforma agora possui um corte vertical comercial e de automacao ligando `identity`, `crm`, `sales`, `workflow-control`, `workflow-runtime`, `analytics`, `webhook-hub` e `edge` com fluxos em PostgreSQL e validacao container-first.

### Destaques

- arquitetura multi-tenant desde o inicio
- identidade, CRM, vendas e automacao ja conectados no mesmo fluxo operacional
- schemas PostgreSQL, migrations e seeds separados por contexto
- `edge` expondo cockpits de plataforma, automacao e vendas
- `analytics` expondo relatorios de pipeline, vendas e visao 360 por tenant
- validacao `unit`, `integration`, `contract` e `smoke` rodando em containers
- GitHub Actions e GitHub Container Registry configurados para CI/CD

### O que faz este projeto se destacar?

```text
Go, .NET, TypeScript, Elixir, Python e Rust mapeados para papeis operacionais claros
Monorepo organizado por linguagem e ownership de servico
Schemas PostgreSQL, migrations e bootstrap seeds por contexto de dominio
Validacao container-first nas camadas unit, integration, contract e smoke
Fluxo comercial ja conectado de captura de lead ate conversao em venda
Dashboards operacionais expostos por analytics e agregacao no edge
Pipelines de CI/CD preparadas para GitHub Actions e GHCR
```

---

## Inicio Rapido

### Opcao 1: Scripts Shell

```bash
./scripts/build.sh
./scripts/up.sh
./scripts/down.sh
```

### Opcao 2: Docker Compose

```bash
docker compose --env-file .env.example -f infra/docker-compose.yml up --build -d
```

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Servicos Disponiveis

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

## Escopo do MVP

O MVP atual e focado em backend e ja cobre um corte vertical completo:

- captura de lead, atribuicao de owner, status e historico de notas em `crm`
- pipeline de oportunidade, ciclo de proposta, conversao em venda e transicoes de receita em `sales`
- catalogo de workflows, publicacao, runs e eventos operacionais em `workflow-control`
- ciclo de vida de execucao, retries e transicoes em `workflow-runtime`
- relatorios de vendas, tenant e automacao em `analytics`
- cockpits operacionais agregados em `edge`
- intake de webhooks e trilha de transicoes de entrega em `webhook-hub`

---

## Visao Geral da API

| Area | Endpoint | Descricao |
|------|----------|-----------|
| Identity | `GET /api/identity/tenants/{slug}/snapshot` | Le o snapshot operacional do tenant |
| CRM | `POST /api/crm/leads` | Cria um lead no funil comercial |
| CRM | `POST /api/crm/leads/{publicId}/notes` | Registra notas de relacionamento para um lead |
| Sales | `POST /api/sales/opportunities` | Cria uma oportunidade vinculada a um lead |
| Sales | `POST /api/sales/opportunities/{publicId}/proposals` | Cria uma proposta para uma oportunidade |
| Sales | `POST /api/sales/proposals/{publicId}/convert` | Converte uma proposta aceita em venda |
| Workflow Control | `POST /api/workflow-control/runs` | Cria uma workflow run para um sujeito de negocio |
| Workflow Runtime | `POST /api/workflow-runtime/executions` | Cria uma execucao runtime a partir de uma definicao |
| Analytics | `GET /api/analytics/reports/sales-journey` | Le o funil comercial do lead ate a venda |
| Analytics | `GET /api/analytics/reports/tenant-360` | Le uma visao consolidada do tenant |
| Edge | `GET /api/edge/ops/automation-overview` | Le o cockpit operacional de automacao |
| Edge | `GET /api/edge/ops/sales-overview` | Le o cockpit comercial |
| Webhook Hub | `POST /api/webhook-hub/events` | Recebe eventos externos por webhook |
| Ops | `GET /health/live` | Liveness |
| Ops | `GET /health/ready` | Readiness |

Mais detalhes estao nas referencias por servico:

- [service-api/service-csharp/identity/README.md](service-api/service-csharp/identity/README.md)
- [service-api/service-golang/crm/README.md](service-api/service-golang/crm/README.md)
- [service-api/service-golang/sales/README.md](service-api/service-golang/sales/README.md)
- [service-api/service-typescript/workflow-control/README.md](service-api/service-typescript/workflow-control/README.md)
- [service-api/service-elixir/workflow-runtime/README.md](service-api/service-elixir/workflow-runtime/README.md)
- [service-api/service-python/analytics/README.md](service-api/service-python/analytics/README.md)
- [service-api/service-rust/webhook-hub/README.md](service-api/service-rust/webhook-hub/README.md)
- [service-api/service-golang/edge/README.md](service-api/service-golang/edge/README.md)

---

## Scripts de Automacao

### `scripts/build.sh`

- builda as imagens Docker locais definidas no Compose

### `scripts/up.sh`

- sobe a stack local em modo detached com build

### `scripts/down.sh`

- derruba a stack e remove containers orfaos

### `scripts/logs.sh <service>`

- acompanha os logs de um servico do Compose

### `scripts/db.sh`

- aplica migrations, seeds e resumos relacionais por contexto de dominio

### `scripts/test.sh`

- executa `unit`, `integration`, `contract` e `smoke` em modo container-first

---

## Contato

**Thiago Di Faria** - thiagodifaria@gmail.com

[![GitHub](https://img.shields.io/badge/GitHub-@thiagodifaria-black?style=flat&logo=github)](https://github.com/thiagodifaria)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Thiago_Di_Faria-blue?style=flat&logo=linkedin)](https://linkedin.com/in/thiagodifaria)

---

## Agradecimentos

Agradecimentos especiais para:

- os ecossistemas Go, .NET, TypeScript, Elixir, Python e Rust
- mantenedores de PostgreSQL, Redis e Docker
- as bibliotecas open source usadas ao longo da plataforma

---

### Se este projeto te ajudou, deixa uma estrela

**Made by [Thiago Di Faria](https://github.com/thiagodifaria)**
