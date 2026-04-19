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
[![Tests](https://img.shields.io/badge/Tests-Unit%20Integration%20Contract%20Smoke-success?style=flat)]()

---

## Documentacao

**Visao Geral do Projeto:** [README.md](README.md)  
**Read in English:** [README_EN.md](README_EN.md)  
**Leia em Portugues:** [README_PT.md](README_PT.md)  
**Architecture Reference:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Engineering Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Integration Map:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations Reference:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Changelog:** [docs/CHANGELOG.md](docs/CHANGELOG.md)

---

## O que e o ERP?

ERP e uma plataforma enterprise pensada para parecer um produto interno real, e nao um CRUD generico de estudo. O repositorio esta estruturado como um backend multi-tenant, container-first e poliglota, em que cada servico possui uma responsabilidade operacional delimitada e um contexto proprio de persistencia quando necessario.

Este repositorio tambem e a evolucao publica de estruturas de ERP ja aplicadas em empresas reais. A ideia nao e apenas construir um sistema novo, mas transformar essa experiencia acumulada em:

- um template reutilizavel para sistemas enterprise
- um study case de arquitetura e design de sistemas
- um projeto de portfolio com escopo realista
- uma referencia de ERP empresarial em formato publico e extensivel

O marco atual ja e um MVP backend-first real. A plataforma agora expoe um corte vertical comercial e de automacao conectando `identity`, `crm`, `sales`, `workflow-control`, `workflow-runtime`, `analytics`, `webhook-hub` e `edge`.

---

## Estado Atual do MVP

Hoje o MVP ja inclui:

- gestao de tenants, empresas, usuarios, times e roles em `identity`
- captura de leads, ownership, progressao de status e historico de notas em `crm`
- ciclo de oportunidade, proposta e venda em `sales`
- catalogo de definicoes, versionamento, runs e eventos operacionais em `workflow-control`
- execucoes runtime, retries e ledger de transicoes em `workflow-runtime`
- relatorios operacionais agregados em `analytics`
- intake de webhooks e rastreio de transicoes em `webhook-hub`
- cockpits operacionais em `edge`

O projeto ainda e backend-first. Isso significa que o ponto mais forte do MVP hoje esta em arquitetura, comportamento de servico, contratos, ownership de banco e visibilidade operacional, e nao ainda em uma camada frontend completa.

---

## Por que este projeto existe

ERP foi desenhado para cumprir quatro papeis ao mesmo tempo:

- template enterprise: base para futuros sistemas e aceleradores internos
- study case: exemplo concreto de arquitetura multi-tenant, poliglota e container-first
- portfolio: projeto que demonstra profundidade tecnica, consistencia e realismo operacional
- referencia publica: uma forma de expor como um ERP empresarial moderno pode ser estruturado

Por isso o repositorio se preocupa tanto com fronteiras de dominio, health endpoints, ownership de banco, estabilidade de contratos, smoke real e composicao entre servicos.

---

## Arquitetura em alto nivel

A plataforma esta separada intencionalmente em planos operacionais:

- plano transacional: `identity`, `crm`, `sales`
- plano de controle: `workflow-control`
- plano de execucao: `workflow-runtime`
- plano analitico: `analytics`
- plano de integracao: `webhook-hub`
- plano de agregacao e operacao publica: `edge`

Essa separacao evita misturar no mesmo servico responsabilidades de catalogo, execucao duravel, escrita transacional e leitura analitica.

---

## Stack e racional tecnico

A stack e poliglota por decisao arquitetural, nao por estetica.

- `Go` foi usado em `edge`, `crm` e `sales` porque esses servicos se beneficiam de alta performance HTTP, deploy simples e baixo atrito operacional.
- `.NET` foi usado em `identity` porque controle de acesso, tenancy e futuros dominios financeiros combinam bem com ergonomia enterprise e uma camada de aplicacao madura.
- `TypeScript` foi usado em `workflow-control` porque APIs de controle e catalogo com muita modelagem de metadado ganham produtividade nessa stack.
- `Elixir` foi usado em `workflow-runtime` porque execucao duravel, retries, timers e automacao concorrente combinam naturalmente com OTP.
- `Python` foi usado em `analytics` porque read models pesados, relatorios e futuras cargas de forecasting ficam muito confortaveis nesse ecossistema.
- `Rust` foi usado em `webhook-hub` porque intake de eventos externos, idempotencia e controle rigoroso de transicao ganham previsibilidade extra.
- `PostgreSQL` e o armazenamento transacional por dominio, com schemas, migrations e seeds separados por contexto.

---

## Inventario de Servicos

| Servico | Linguagem | Responsabilidade | Superficie principal |
|---------|-----------|------------------|----------------------|
| `identity` | .NET | tenancy, empresas, usuarios, times, roles e base de acesso | bootstrap de tenant e snapshot de acesso |
| `crm` | Go | leads, ownership, status, resumo e notas | funil de leads e historico de relacionamento |
| `sales` | Go | oportunidades, propostas, conversao em venda e receita | ciclo comercial de oportunidade ate venda |
| `workflow-control` | TypeScript | definicoes, versoes, runs e eventos | APIs de controle e ledger operacional |
| `workflow-runtime` | Elixir | execucao duravel, transicoes e retries | runtime de automacao e resumo de execucoes |
| `analytics` | Python | leituras operacionais pesadas e relatorios | relatorios de vendas, tenant e automacao |
| `webhook-hub` | Rust | intake de webhooks e transicoes de entrega | ingestao de eventos externos |
| `edge` | Go | agregacao entre servicos e cockpit publico | health, tenant, automacao e overview comercial |

---

## Corte Vertical Comercial e de Automacao

A maneira mais clara de entender o MVP atual e seguir o fluxo principal:

1. Um lead entra em `crm`.
2. O lead pode receber owner, notas e progressao de status.
3. Uma oportunidade vinculada e criada em `sales`.
4. Uma proposta e criada para essa oportunidade.
5. A proposta pode passar por transicoes de status e ser convertida em venda.
6. `workflow-control` pode criar uma workflow run ligada ao mesmo sujeito de negocio.
7. `workflow-runtime` pode executar a definicao correspondente com transicoes e retries.
8. `analytics` le esse footprint e expoe relatorios como `pipeline-summary`, `sales-journey`, `tenant-360`, `automation-board` e `workflow-definition-health`.
9. `edge` agrega essas leituras em `tenant-overview`, `automation-overview` e `sales-overview`.

Com isso, o repositorio ganha uma narrativa de negocio ponta a ponta, em vez de apenas demonstracoes isoladas por servico.

---

## Snapshot da API Publica

### Identity

Principais rotas:

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

Principais rotas:

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

Principais rotas:

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

Principais rotas:

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

Principais rotas:

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

Principais rotas:

- `GET /api/analytics/reports/pipeline-summary`
- `GET /api/analytics/reports/service-pulse`
- `GET /api/analytics/reports/sales-journey`
- `GET /api/analytics/reports/tenant-360`
- `GET /api/analytics/reports/automation-board`
- `GET /api/analytics/reports/workflow-definition-health`
- `GET /api/analytics/reports/delivery-reliability`

### Edge

Principais rotas:

- `GET /api/edge/ops/health`
- `GET /api/edge/ops/tenant-overview`
- `GET /api/edge/ops/automation-overview`
- `GET /api/edge/ops/sales-overview`

### Webhook Hub

Principais rotas:

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

Todos os servicos tambem expoem:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`

---

## Desenvolvimento Local

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

### Comandos uteis

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

### Servicos disponiveis

- Edge: `http://localhost:8080`
- Identity: `http://localhost:8081`
- Webhook Hub: `http://localhost:8082`
- CRM: `http://localhost:8083`
- Workflow Control: `http://localhost:8084`
- Workflow Runtime: `http://localhost:8085`
- Analytics: `http://localhost:8086`
- Sales: `http://localhost:8087`
- PostgreSQL: `localhost:5432` por padrao
- Redis: `localhost:6379` por padrao

Se alguma dessas portas do host ja estiver ocupada, os scripts container-first remapeiam automaticamente durante a subida local e a validacao smoke.

---

## Estrategia de Validacao

O repositorio e intencionalmente container-first. Hoje a validacao cobre:

- `unit`: testes de servico em Go, TypeScript, Elixir, Python, .NET e Rust
- `integration`: suite HTTP dedicada de `identity`
- `contract`: suites de contrato publico para `workflow-control`, `crm`, `sales` e `identity`
- `smoke`: validacao end-to-end em Docker Compose com PostgreSQL e Redis

Isso significa que o MVP nao esta apenas modelado em codigo, mas tambem verificado perto da topologia real de runtime.

---

## CI/CD e Publicacao de Containers

O repositorio ja inclui:

- workflow `Quality` no GitHub Actions para `unit`, `integration`, `contract` e `smoke`
- workflow `Containers` para publicar imagens no `ghcr.io`

Imagens esperadas:

- `ghcr.io/thiagodifaria/erp-edge`
- `ghcr.io/thiagodifaria/erp-crm`
- `ghcr.io/thiagodifaria/erp-sales`
- `ghcr.io/thiagodifaria/erp-identity`
- `ghcr.io/thiagodifaria/erp-workflow-control`
- `ghcr.io/thiagodifaria/erp-workflow-runtime`
- `ghcr.io/thiagodifaria/erp-analytics`
- `ghcr.io/thiagodifaria/erp-webhook-hub`

---

## O que pode ser estudado aqui

Este repositorio e util se voce quiser estudar:

- fronteiras de servico em ambiente multi-tenant
- arquitetura poliglota com atribuicao explicita de papeis por linguagem
- ownership de PostgreSQL por contexto de dominio
- separacao entre control plane e runtime plane
- intake de webhooks e trilha operacional de transicoes
- agregacao de servicos e read models executivos
- validacao container-first dentro de um monorepo
- apresentacao de backend enterprise em formato portfolio

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
