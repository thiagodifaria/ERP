# ERP

![ERP](https://img.shields.io/badge/ERP-Enterprise%20Platform-111827?style=for-the-badge&logo=github&logoColor=white)

**Plataforma ERP backend-first para operacao empresarial, automacao, contratos, financeiro, billing, analytics e integracoes.**

[![Go](https://img.shields.io/badge/Go-edge%20crm%20sales%20rentals%20documents-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![.NET](https://img.shields.io/badge/.NET-identity%20finance%20billing-512BD4?style=flat&logo=dotnet&logoColor=white)](https://dotnet.microsoft.com/)
[![TypeScript](https://img.shields.io/badge/TypeScript-workflow%20engagement-3178C6?style=flat&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Elixir](https://img.shields.io/badge/Elixir-workflow%20runtime-4B275F?style=flat&logo=elixir&logoColor=white)](https://elixir-lang.org/)
[![Python](https://img.shields.io/badge/Python-analytics%20platform%20admin-3776AB?style=flat&logo=python&logoColor=white)](https://www.python.org/)
[![Rust](https://img.shields.io/badge/Rust-webhook%20hub-000000?style=flat&logo=rust&logoColor=white)](https://www.rust-lang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-domain%20schemas-316192?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-container%20first-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)

---

## Documentacao

**Short overview:** [README.md](README.md)  
**English detailed README:** [README_EN.md](README_EN.md)  
**README detalhado em Portugues:** [README_PT.md](README_PT.md)  
**API Reference:** [docs/API.md](docs/API.md)  
**Architecture:** [docs/ARQUITETURA.md](docs/ARQUITETURA.md)  
**Contracts:** [docs/CONTRATOS.md](docs/CONTRATOS.md)  
**Integrations:** [docs/INTEGRACOES.md](docs/INTEGRACOES.md)  
**Operations:** [docs/OPERACOES.md](docs/OPERACOES.md)  
**Standards:** [docs/PADROES.md](docs/PADROES.md)  
**Services:** [docs/SERVICOS.md](docs/SERVICOS.md)

---

## O que e o ERP?

ERP e uma plataforma empresarial multi-tenant, poliglota e container-first. O repositorio trata operacao comercial, contratos recorrentes, documentos, financeiro, billing, automacao, runtime de workflows, campanhas, analytics, simulacao, governanca SaaS, suporte, fornecedores, notificacoes, fiscal/compliance e webhooks como partes conectadas de uma mesma plataforma.

A proposta e servir como referencia tecnica e produto de portfolio para uma arquitetura enterprise realista. O foco esta em fronteiras de dominio, ownership de banco, contratos versionados, smoke integrado, health real, adapters externos e automacao operacional.

## Escala Atual

| Metrica | Valor |
|--------|-------|
| Servicos com OpenAPI versionado | 20 |
| Endpoints HTTP versionados | 201 |
| Contract catalog | `docs/contracts/` |
| Runtime command | `./scripts/build.sh` |
| Validation command | `./scripts/test.sh` |

## Planos Arquiteturais

- administrative plane: `support`
- administrative/notification plane: `notification`
- administrative/procurement plane: `supplier`
- analytics plane: `analytics`
- compliance plane: `fiscal`
- control plane: `workflow-control`
- integration plane: `webhook-hub`
- interaction/control plane: `engagement`
- platform control plane: `platform-control`
- public operations plane: `edge`
- runtime plane: `workflow-runtime`
- simulation plane: `simulation`
- transaction plane: `crm`, `documents`, `rentals`, `sales`
- transaction/billing plane: `billing`
- transaction/catalog plane: `catalog`
- transaction/finance plane: `finance`
- transaction/security plane: `identity`

## Racional da Stack

- .NET: usado em `billing`, `finance`, `identity` para manter cada tipo de carga no ecossistema mais coerente com sua operacao, sem forcar uma unica linguagem para problemas diferentes.
- Elixir: usado em `workflow-runtime` para manter cada tipo de carga no ecossistema mais coerente com sua operacao, sem forcar uma unica linguagem para problemas diferentes.
- Go: usado em `crm`, `documents`, `edge`, `rentals`, `sales` para manter cada tipo de carga no ecossistema mais coerente com sua operacao, sem forcar uma unica linguagem para problemas diferentes.
- Python: usado em `analytics`, `catalog`, `fiscal`, `notification`, `platform-control`, `simulation`, `supplier`, `support` para manter cada tipo de carga no ecossistema mais coerente com sua operacao, sem forcar uma unica linguagem para problemas diferentes.
- Rust: usado em `webhook-hub` para manter cada tipo de carga no ecossistema mais coerente com sua operacao, sem forcar uma unica linguagem para problemas diferentes.
- TypeScript: usado em `engagement`, `workflow-control` para manter cada tipo de carga no ecossistema mais coerente com sua operacao, sem forcar uma unica linguagem para problemas diferentes.

## Inventario de Servicos

| Servico | Stack | Plano | Responsabilidade | Endpoints |
|---------|-------|-------|----------------|-----------|
| `analytics` | Python | analytics plane | relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas | 9 |
| `billing` | .NET | transaction/billing plane | planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery | 9 |
| `catalog` | Python | transaction/catalog plane | categorias, itens, versoes de item, bulk e contratos de consumo | 12 |
| `crm` | Go | transaction plane | leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento | 5 |
| `documents` | Go | transaction plane | anexos, upload, storage posture, assinatura, versoes, archive e access links | 10 |
| `edge` | Go | public operations plane | entrada publica, agregacao cross-service e cockpits operacionais | 8 |
| `engagement` | TypeScript | interaction/control plane | campanhas, templates, touchpoints, conversas, delivery, providers e callbacks | 9 |
| `finance` | .NET | transaction/finance plane | recebiveis, payables, caixa, custos, comissoes e fechamento | 5 |
| `fiscal` | Python | compliance plane | perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria | 25 |
| `identity` | .NET | transaction/security plane | tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria | 6 |
| `notification` | Python | administrative/notification plane | preferencias, centro interno de alertas, severidade e lifecycle de notificacoes | 7 |
| `platform-control` | Python | platform control plane | capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live | 40 |
| `rentals` | Go | transaction plane | contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais | 4 |
| `sales` | Go | transaction plane | oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias | 6 |
| `simulation` | Python | simulation plane | cenarios operacionais, benchmark de carga e modelagem de capacidade | 3 |
| `supplier` | Python | administrative/procurement plane | categorias de fornecedor, diretorio de fornecedores e procurement ownership | 8 |
| `support` | Python | administrative plane | filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento | 9 |
| `webhook-hub` | Rust | integration plane | intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries | 13 |
| `workflow-control` | TypeScript | control plane | definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos | 7 |
| `workflow-runtime` | Elixir | runtime plane | execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes | 6 |

## Cortes Verticais Principais

### Jornada comercial

lead no CRM, oportunidade e proposta em sales, conversao, invoice, visibilidade financeira e cockpit no edge

### Jornada contratual recorrente

cliente, contrato de locacao/recorrencia, cobrancas futuras, reajustes, terminacao, documentos e projecoes financeiras

### Jornada de automacao

definicao de workflow, versao publicada, run de controle, execucao runtime, timeline, retries e visibilidade analitica

### Jornada de integracao

evento externo, intake no webhook-hub, validacao, fila, processamento, forwarding, dead letter e postura analitica

### Jornada de governanca SaaS

catalogo de capabilities, provider defaults, entitlements, quotas, metering, onboarding/offboarding e go-live

## Desenvolvimento Local

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh logs edge
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh down
```

## Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh contract
./scripts/test.sh platform
./scripts/test.sh smoke
./scripts/test.sh performance
./scripts/test.sh backup-restore
./scripts/test.sh hardening
```

## Catalogo de Contratos

- `docs/contracts/http/`
- `docs/contracts/events/`
- `docs/contracts/registry.json`
- `docs/contracts/schema-registry.json`
- `docs/contracts/portal/index.html`

## Servicos em Detalhe

### `analytics`

- Stack: Python
- Plano: analytics plane
- Codigo: `service-api/service-python/analytics`
- Contexto de banco: `analytics/simulation read models`
- Contrato: `docs/contracts/http/analytics.openapi.yaml`
- Versao OpenAPI: `0.1.0`
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.

Rotas:

- `GET /api/analytics/reports/adapter-catalog` - Read external adapter capability catalog
- `GET /api/analytics/reports/integration-readiness` - Read external integration readiness
- `GET /api/analytics/reports/saas-control` - Read SaaS control posture by tenant
- `GET /api/analytics/reports/contract-governance` - Read contract governance posture
- `GET /api/analytics/reports/hardening-review` - Read hardening review
- `GET /api/analytics/reports/core-operations` - Read core product operations
- `GET /api/analytics/reports/relationship-intelligence` - Read relationship intelligence
- `GET /api/analytics/reports/compliance-control` - Read fiscal and privacy compliance control
- `GET /api/analytics/reports/go-live-control` - Read go-live rollout control

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `billing`

- Stack: .NET
- Plano: transaction/billing plane
- Codigo: `service-api/service-csharp/billing`
- Contexto de banco: `billing`
- Contrato: `docs/contracts/http/billing.openapi.yaml`
- Versao OpenAPI: `0.9.7`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.

Rotas:

- `GET /health/details` - Return readiness details and gateway posture
- `GET /api/billing/gateways` - List gateway capabilities and Pix posture
- `GET /api/billing/gateways/{provider}` - Read one gateway capability
- `GET /api/billing/plans` - List billing plans including flat, hybrid and usage-based pricing
- `POST /api/billing/plans` - Create billing plan
- `GET /api/billing/subscriptions` - List subscriptions
- `POST /api/billing/subscriptions` - Create subscription
- `GET /api/billing/subscriptions/{publicId}/usage-pricing` - Project usage-based charge for one subscription
- `POST /api/billing/invoices/{publicId}/attempts` - Create payment attempt with idempotency support

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `catalog`

- Stack: Python
- Plano: transaction/catalog plane
- Codigo: `service-api/service-python/catalog`
- Contexto de banco: `catalog`
- Contrato: `docs/contracts/http/catalog.openapi.yaml`
- Versao OpenAPI: `0.2.0`
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.

Rotas:

- `GET /api/catalog/capabilities` - Read catalog capability posture
- `GET /api/catalog/consumers` - Read catalog consumer contracts across core domains
- `GET /api/catalog/categories` - List categories by tenant
- `POST /api/catalog/categories` - Create one category
- `GET /api/catalog/categories/page` - Cursor-based category listing
- `GET /api/catalog/items` - List catalog items
- `POST /api/catalog/items` - Create one catalog item
- `GET /api/catalog/items/page` - Cursor-based item listing
- `POST /api/catalog/items/bulk` - Bulk create catalog items with partial success
- `GET /api/catalog/items/{publicId}` - Read one catalog item
- `PATCH /api/catalog/items/{publicId}` - Update active state, price and attributes
- `GET /api/catalog/items/{publicId}/versions` - Read catalog item version history

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `crm`

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/crm`
- Contexto de banco: `crm`
- Contrato: `docs/contracts/http/crm.openapi.yaml`
- Versao OpenAPI: `0.2.0`
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento.

Rotas:

- `GET /api/crm/enrichment/cnpj/capabilities` - Read CNPJ enrichment provider capabilities
- `POST /api/crm/enrichment/cnpj/lookup` - Lookup and enrich one CNPJ through provider contract
- `GET /api/crm/pipeline/config` - Read tenant pipeline configuration
- `PUT /api/crm/pipeline/config` - Upsert tenant pipeline configuration
- `GET /api/crm/leads/intelligence/summary` - Read lead scoring and pipeline intelligence summary

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `documents`

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/documents`
- Contexto de banco: `documents`
- Contrato: `docs/contracts/http/documents.openapi.yaml`
- Versao OpenAPI: `0.9.7`
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.

Rotas:

- `GET /health/details` - Return runtime readiness and storage posture
- `GET /api/documents/signing/capabilities` - List digital signature capabilities
- `GET /api/documents/signing/capabilities/{provider}` - Read one signing capability
- `POST /api/documents/signing/requests` - Queue one digital signature request
- `GET /api/documents/storage/capabilities` - List storage capability registry
- `GET /api/documents/storage/capabilities/{provider}` - Read one storage capability
- `GET /api/documents/attachments` - List attachments
- `POST /api/documents/attachments` - Create attachment metadata
- `GET /api/documents/attachments/{publicId}/versions` - List attachment versions
- `POST /api/documents/attachments/{publicId}/versions` - Append attachment version

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `edge`

- Stack: Go
- Plano: public operations plane
- Codigo: `service-api/service-golang/edge`
- Contexto de banco: `none`
- Contrato: `docs/contracts/http/edge.openapi.yaml`
- Versao OpenAPI: `0.1.0`
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.

Rotas:

- `GET /api/edge/ops/core-operations` - Read executive core product cockpit
- `GET /api/edge/ops/relationship-overview` - Read executive relationship cockpit
- `GET /api/edge/ops/compliance-overview` - Read executive compliance cockpit
- `GET /api/edge/ops/go-live-overview` - Read executive go-live cockpit
- `GET /api/edge/ops/integrations-overview` - Read executive integrations cockpit
- `GET /api/edge/ops/saas-overview` - Read executive SaaS cockpit
- `GET /api/edge/ops/contracts-overview` - Read executive contracts cockpit
- `GET /api/edge/ops/hardening-overview` - Read executive hardening cockpit

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `engagement`

- Stack: TypeScript
- Plano: interaction/control plane
- Codigo: `service-api/service-typescript/engagement`
- Contexto de banco: `engagement`
- Contrato: `docs/contracts/http/engagement.openapi.yaml`
- Versao OpenAPI: `0.9.7`
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.

Rotas:

- `GET /health/details` - Return readiness details for engagement runtime
- `GET /api/engagement/providers` - List provider capabilities and fallback posture
- `GET /api/engagement/providers/{provider}` - Read one provider capability
- `POST /api/engagement/providers/meta-ads/leads` - Ingest inbound lead from Meta Ads
- `POST /api/engagement/providers/resend/events` - Register Resend callback event
- `POST /api/engagement/providers/whatsapp-cloud/events` - Register WhatsApp callback event
- `POST /api/engagement/providers/telegram-bot/events` - Register Telegram callback event
- `GET /api/engagement/provider-events` - List provider events
- `GET /api/engagement/provider-events/{publicId}` - Read one provider event

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `finance`

- Stack: .NET
- Plano: transaction/finance plane
- Codigo: `service-api/service-csharp/finance`
- Contexto de banco: `finance`
- Contrato: `docs/contracts/http/finance.openapi.yaml`
- Versao OpenAPI: `0.4.0`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.

Rotas:

- `GET /api/finance/receivable-projections` - List receivable projections
- `POST /api/finance/receivable-projections/sync` - Sync projections from sales and rentals
- `GET /api/finance/commission-holds` - List commission holds
- `POST /api/finance/commission-holds/{publicId}/release` - Release one commission hold
- `GET /api/finance/activity` - List finance operational activity

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `fiscal`

- Stack: Python
- Plano: compliance plane
- Codigo: `service-api/service-python/fiscal`
- Contexto de banco: `fiscal`
- Contrato: `docs/contracts/http/fiscal.openapi.yaml`
- Versao OpenAPI: `0.1.0`
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.

Rotas:

- `GET /api/fiscal/capabilities` - Read fiscal capability registry
- `GET /api/fiscal/companies/{companyPublicId}/profile` - Read fiscal company profile
- `PUT /api/fiscal/companies/{companyPublicId}/profile` - Upsert fiscal company profile
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies` - List retention policies by company
- `GET /api/fiscal/companies/{companyPublicId}/retention-execution` - Read retention execution plan for one company
- `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute` - Execute retention and anonymization plan
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}` - Upsert retention policy for one data domain
- `GET /api/fiscal/documents` - List fiscal documents
- `POST /api/fiscal/documents` - Issue one fiscal document
- `GET /api/fiscal/documents/{publicId}` - Read one fiscal document
- `POST /api/fiscal/documents/{publicId}/cancel` - Cancel one fiscal document
- `POST /api/fiscal/documents/{publicId}/correction-letter` - Register correction letter for one fiscal document
- `POST /api/fiscal/documents/{publicId}/invalidate` - Register invalidation for one fiscal document
- `GET /api/fiscal/documents/{publicId}/events` - List fiscal document audit events
- `GET /api/fiscal/privacy-requests` - List privacy requests
- `POST /api/fiscal/privacy-requests` - Create privacy request
- `GET /api/fiscal/privacy-requests/{publicId}` - Read one privacy request
- `GET /api/fiscal/privacy-requests/{publicId}/export-package` - Build export package for one privacy request
- `POST /api/fiscal/privacy-requests/{publicId}/execute` - Execute one privacy request with audit trail
- `PATCH /api/fiscal/privacy-requests/{publicId}/status` - Transition privacy request lifecycle status
- `GET /api/fiscal/consents` - List consent ledger
- `POST /api/fiscal/consents` - Create consent record
- `PATCH /api/fiscal/consents/{publicId}` - Transition consent status
- `GET /api/fiscal/audit-events` - List fiscal audit events
- `GET /api/fiscal/compliance/summary` - Read fiscal compliance summary

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `identity`

- Stack: .NET
- Plano: transaction/security plane
- Codigo: `service-api/service-csharp/identity`
- Contexto de banco: `identity`
- Contrato: `docs/contracts/http/identity.openapi.yaml`
- Versao OpenAPI: `0.5.0`
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.

Rotas:

- `GET /api/identity/tenants` - List tenants
- `POST /api/identity/tenants` - Create tenant
- `GET /api/identity/tenants/{slug}/snapshot` - Read one tenant snapshot
- `POST /api/identity/sessions/login` - Authenticate identity session
- `POST /api/identity/sessions/refresh` - Refresh identity session
- `POST /api/identity/invitations` - Create invitation

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `notification`

- Stack: Python
- Plano: administrative/notification plane
- Codigo: `service-api/service-python/notification`
- Contexto de banco: `notification`
- Contrato: `docs/contracts/http/notification.openapi.yaml`
- Versao OpenAPI: `0.1.0`
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.

Rotas:

- `GET /api/notification/capabilities` - Read notification capability catalog
- `GET /api/notification/preferences/{userPublicId}` - Read one user notification preference
- `PUT /api/notification/preferences/{userPublicId}` - Upsert one user notification preference
- `GET /api/notification/center` - List notification center items with cursor filters
- `POST /api/notification/center` - Create one notification center item
- `PATCH /api/notification/center/{publicId}/status` - Transition notification status
- `GET /api/notification/summary` - Read notification summary

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `platform-control`

- Stack: Python
- Plano: platform control plane
- Codigo: `service-api/service-python/platform-control`
- Contexto de banco: `platform-control`
- Contrato: `docs/contracts/http/platform-control.openapi.yaml`
- Versao OpenAPI: `0.2.0`
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.

Rotas:

- `GET /api/platform-control/capabilities/catalog` - List platform capability catalog
- `GET /api/platform-control/providers/catalog` - List provider capability catalog and environment posture
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements` - List tenant entitlements with cursor pagination
- `GET /api/platform-control/tenants/{tenantSlug}/feature-flags` - List tenant feature flags with capability metadata
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}` - Upsert one entitlement
- `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}` - Upsert one feature flag using entitlement governance
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk` - Bulk upsert entitlements with partial success
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults` - List provider defaults selected for one tenant
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}` - Upsert provider default for one tenant capability
- `GET /api/platform-control/tenants/{tenantSlug}/quotas` - List quotas by tenant
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}` - Upsert one quota
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk` - Bulk upsert quotas with partial success
- `GET /api/platform-control/tenants/{tenantSlug}/blocks` - List tenant blocks
- `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}` - Upsert tenant block
- `GET /api/platform-control/tenants/{tenantSlug}/metering` - Read metering snapshots and summary with cursor pagination
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots` - Create one usage snapshot
- `GET /api/platform-control/tenants/{tenantSlug}/usage-summary` - Read quota and metering utilization summary
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness` - Read tenant lifecycle readiness and provider posture
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs` - List onboarding and offboarding jobs with cursor pagination
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}` - Read one lifecycle job with audit events
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview` - Preview onboarding plan, provider defaults and lifecycle readiness
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding` - Queue onboarding job with Idempotency-Key and 202 Accepted
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview` - Preview offboarding plan, retention posture and lifecycle readiness
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding` - Queue offboarding job with Idempotency-Key and 202 Accepted
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start` - Transition lifecycle job to running
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete` - Transition lifecycle job to completed
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail` - Transition lifecycle job to failed
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel` - Transition lifecycle job to cancelled
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` - Read go-live rollout readiness by tenant
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` - Read tenant go-live adoption baseline and gap
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` - List go-live bottlenecks and operational blockers
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` - Read rollout and rollback playbook for one tenant
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` - List recommended go-live adjustments
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply` - Apply one go-live operational adjustment
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - List go-live rollouts
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - Create one go-live rollout
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}` - Read one go-live rollout with events
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` - Transition go-live rollout to running
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` - Transition go-live rollout to completed
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` - Roll back one go-live rollout

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `rentals`

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/rentals`
- Contexto de banco: `rentals`
- Contrato: `docs/contracts/http/rentals.openapi.yaml`
- Versao OpenAPI: `0.8.0`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.

Rotas:

- `GET /api/rentals/contracts` - List rental contracts
- `POST /api/rentals/contracts` - Create rental contract
- `GET /api/rentals/contracts/{publicId}/charges` - List contract charges
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` - Update charge status

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `sales`

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/sales`
- Contexto de banco: `sales`
- Contrato: `docs/contracts/http/sales.openapi.yaml`
- Versao OpenAPI: `0.7.0`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.

Rotas:

- `GET /api/sales/opportunities` - List opportunities
- `POST /api/sales/opportunities` - Create opportunity
- `GET /api/sales/proposals` - List proposals
- `POST /api/sales/proposals` - Create proposal
- `GET /api/sales/sales` - List sales
- `GET /api/sales/invoices` - List commercial invoices

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `simulation`

- Stack: Python
- Plano: simulation plane
- Codigo: `service-api/service-python/simulation`
- Contexto de banco: `simulation`
- Contrato: `docs/contracts/http/simulation.openapi.yaml`
- Versao OpenAPI: `0.7.0`
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.

Rotas:

- `GET /api/simulation/scenarios` - List scenarios
- `POST /api/simulation/scenarios` - Create scenario run
- `POST /api/simulation/benchmarks/load` - Execute one load benchmark run

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `supplier`

- Stack: Python
- Plano: administrative/procurement plane
- Codigo: `service-api/service-python/supplier`
- Contexto de banco: `supplier`
- Contrato: `docs/contracts/http/supplier.openapi.yaml`
- Versao OpenAPI: `0.1.0`
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.

Rotas:

- `GET /api/supplier/capabilities` - Read supplier capability catalog
- `GET /api/supplier/categories` - List supplier categories
- `PUT /api/supplier/categories/{categoryKey}` - Upsert one supplier category
- `GET /api/supplier/suppliers` - List suppliers by tenant and status
- `POST /api/supplier/suppliers` - Create one supplier
- `GET /api/supplier/suppliers/summary` - Read supplier summary
- `GET /api/supplier/suppliers/{publicId}` - Read one supplier
- `PATCH /api/supplier/suppliers/{publicId}` - Update one supplier

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `support`

- Stack: Python
- Plano: administrative plane
- Codigo: `service-api/service-python/support`
- Contexto de banco: `support`
- Contrato: `docs/contracts/http/support.openapi.yaml`
- Versao OpenAPI: `0.1.0`
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.

Rotas:

- `GET /api/support/capabilities` - Read support capability catalog
- `GET /api/support/queues` - List support queues by tenant
- `PUT /api/support/queues/{queueKey}` - Upsert one support queue
- `GET /api/support/cases` - List support cases with cursor filters
- `POST /api/support/cases` - Create one support case
- `GET /api/support/cases/summary` - Read support case summary
- `GET /api/support/cases/{publicId}` - Read one support case
- `PATCH /api/support/cases/{publicId}/status` - Transition support case status
- `POST /api/support/cases/{publicId}/comments` - Append comment to support case

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `webhook-hub`

- Stack: Rust
- Plano: integration plane
- Codigo: `service-api/service-rust/webhook-hub`
- Contexto de banco: `webhook-hub`
- Contrato: `docs/contracts/http/webhook-hub.openapi.yaml`
- Versao OpenAPI: `0.9.7`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.

Rotas:

- `GET /health/details` - Return readiness details for webhook runtime
- `GET /api/webhook-hub/capabilities` - Read outbound webhook capability posture
- `GET /api/webhook-hub/outbound-endpoints` - List tenant outbound endpoints
- `POST /api/webhook-hub/outbound-endpoints` - Register one tenant outbound endpoint
- `GET /api/webhook-hub/outbound-endpoints/{publicId}` - Read one tenant outbound endpoint
- `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - List outbound delivery log for one endpoint
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - Register one outbound delivery attempt
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter` - Move one outbound delivery to dead letter
- `GET /api/webhook-hub/events` - List inbound webhook events
- `POST /api/webhook-hub/events` - Register inbound webhook event
- `GET /api/webhook-hub/events/summary` - Aggregate inbound webhook state
- `POST /api/webhook-hub/events/{publicId}/dead-letter` - Move event to dead letter queue
- `POST /api/webhook-hub/events/{publicId}/requeue` - Requeue dead-letter event

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `workflow-control`

- Stack: TypeScript
- Plano: control plane
- Codigo: `service-api/service-typescript/workflow-control`
- Contexto de banco: `workflow-control`
- Contrato: `docs/contracts/http/workflow-control.openapi.yaml`
- Versao OpenAPI: `0.6.0`
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.

Rotas:

- `GET /api/workflow-control/definitions` - List workflow definitions
- `POST /api/workflow-control/definitions` - Create workflow definition
- `GET /api/workflow-control/definitions/{key}` - Read one workflow definition
- `PATCH /api/workflow-control/definitions/{key}` - Update one workflow definition
- `PATCH /api/workflow-control/definitions/{key}/status` - Update workflow definition status
- `GET /api/workflow-control/capabilities/triggers` - List workflow trigger catalog
- `GET /api/workflow-control/capabilities/actions` - List workflow action catalog

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

### `workflow-runtime`

- Stack: Elixir
- Plano: runtime plane
- Codigo: `service-api/service-elixir/workflow-runtime`
- Contexto de banco: `workflow-runtime`
- Contrato: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Versao OpenAPI: `0.6.0`
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.

Rotas:

- `GET /api/workflow-runtime/executions` - List workflow executions
- `POST /api/workflow-runtime/executions` - Create workflow execution
- `GET /api/workflow-runtime/executions/{publicId}` - Read one workflow execution
- `GET /api/workflow-runtime/executions/{publicId}/actions` - List execution action snapshots
- `POST /api/workflow-runtime/executions/{publicId}/advance` - Advance one workflow execution
- `GET /api/workflow-runtime/capabilities` - List runtime capabilities

Notas operacionais:

- Mantem ownership claro do dominio e nao deve gravar em schema de outro contexto sem decisao explicita.
- Deve manter health real, contrato atualizado e validacao no fluxo central quando a superficie publica mudar.
- Mutacoes relevantes precisam preservar tenant, ator, correlation id e historico operacional quando aplicavel.

---

**Thiago Di Faria** - thiagodifaria@gmail.com

## Detalhe dos Contratos HTTP

Esta secao espelha o catalogo OpenAPI atual em alto nivel para leitura rapida do repositorio.

### `analytics`

- Contrato: `docs/contracts/http/analytics.openapi.yaml`
- Endpoints: `9`

#### `GET /api/analytics/reports/adapter-catalog`

- Resumo: Read external adapter capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/integration-readiness`

- Resumo: Read external integration readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/saas-control`

- Resumo: Read SaaS control posture by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/contract-governance`

- Resumo: Read contract governance posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/hardening-review`

- Resumo: Read hardening review.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/core-operations`

- Resumo: Read core product operations.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/relationship-intelligence`

- Resumo: Read relationship intelligence.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/compliance-control`

- Resumo: Read fiscal and privacy compliance control.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/analytics/reports/go-live-control`

- Resumo: Read go-live rollout control.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `billing`

- Contrato: `docs/contracts/http/billing.openapi.yaml`
- Endpoints: `9`

#### `GET /health/details`

- Resumo: Return readiness details and gateway posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/billing/gateways`

- Resumo: List gateway capabilities and Pix posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/billing/gateways/{provider}`

- Resumo: Read one gateway capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/billing/plans`

- Resumo: List billing plans including flat, hybrid and usage-based pricing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/billing/plans`

- Resumo: Create billing plan.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/billing/subscriptions`

- Resumo: List subscriptions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/billing/subscriptions`

- Resumo: Create subscription.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/billing/subscriptions/{publicId}/usage-pricing`

- Resumo: Project usage-based charge for one subscription.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/billing/invoices/{publicId}/attempts`

- Resumo: Create payment attempt with idempotency support.
- Parametros: `Idempotency-Key`, `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

### `catalog`

- Contrato: `docs/contracts/http/catalog.openapi.yaml`
- Endpoints: `12`

#### `GET /api/catalog/capabilities`

- Resumo: Read catalog capability posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/catalog/consumers`

- Resumo: Read catalog consumer contracts across core domains.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/catalog/categories`

- Resumo: List categories by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/catalog/categories`

- Resumo: Create one category.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/catalog/categories/page`

- Resumo: Cursor-based category listing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/catalog/items`

- Resumo: List catalog items.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/catalog/items`

- Resumo: Create one catalog item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/catalog/items/page`

- Resumo: Cursor-based item listing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/catalog/items/bulk`

- Resumo: Bulk create catalog items with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/catalog/items/{publicId}`

- Resumo: Read one catalog item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `PATCH /api/catalog/items/{publicId}`

- Resumo: Update active state, price and attributes.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `GET /api/catalog/items/{publicId}/versions`

- Resumo: Read catalog item version history.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `crm`

- Contrato: `docs/contracts/http/crm.openapi.yaml`
- Endpoints: `5`

#### `GET /api/crm/enrichment/cnpj/capabilities`

- Resumo: Read CNPJ enrichment provider capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/crm/enrichment/cnpj/lookup`

- Resumo: Lookup and enrich one CNPJ through provider contract.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/crm/pipeline/config`

- Resumo: Read tenant pipeline configuration.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/crm/pipeline/config`

- Resumo: Upsert tenant pipeline configuration.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/crm/leads/intelligence/summary`

- Resumo: Read lead scoring and pipeline intelligence summary.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.

### `documents`

- Contrato: `docs/contracts/http/documents.openapi.yaml`
- Endpoints: `10`

#### `GET /health/details`

- Resumo: Return runtime readiness and storage posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/documents/signing/capabilities`

- Resumo: List digital signature capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/documents/signing/capabilities/{provider}`

- Resumo: Read one signing capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/documents/signing/requests`

- Resumo: Queue one digital signature request.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/documents/storage/capabilities`

- Resumo: List storage capability registry.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/documents/storage/capabilities/{provider}`

- Resumo: Read one storage capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/documents/attachments`

- Resumo: List attachments.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/documents/attachments`

- Resumo: Create attachment metadata.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/documents/attachments/{publicId}/versions`

- Resumo: List attachment versions.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/documents/attachments/{publicId}/versions`

- Resumo: Append attachment version.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

### `edge`

- Contrato: `docs/contracts/http/edge.openapi.yaml`
- Endpoints: `8`

#### `GET /api/edge/ops/core-operations`

- Resumo: Read executive core product cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/relationship-overview`

- Resumo: Read executive relationship cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/compliance-overview`

- Resumo: Read executive compliance cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/go-live-overview`

- Resumo: Read executive go-live cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/integrations-overview`

- Resumo: Read executive integrations cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/saas-overview`

- Resumo: Read executive SaaS cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/contracts-overview`

- Resumo: Read executive contracts cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/edge/ops/hardening-overview`

- Resumo: Read executive hardening cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `engagement`

- Contrato: `docs/contracts/http/engagement.openapi.yaml`
- Endpoints: `9`

#### `GET /health/details`

- Resumo: Return readiness details for engagement runtime.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/engagement/providers`

- Resumo: List provider capabilities and fallback posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/engagement/providers/{provider}`

- Resumo: Read one provider capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/engagement/providers/meta-ads/leads`

- Resumo: Ingest inbound lead from Meta Ads.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `POST /api/engagement/providers/resend/events`

- Resumo: Register Resend callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `POST /api/engagement/providers/whatsapp-cloud/events`

- Resumo: Register WhatsApp callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `POST /api/engagement/providers/telegram-bot/events`

- Resumo: Register Telegram callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/engagement/provider-events`

- Resumo: List provider events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/engagement/provider-events/{publicId}`

- Resumo: Read one provider event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

### `finance`

- Contrato: `docs/contracts/http/finance.openapi.yaml`
- Endpoints: `5`

#### `GET /api/finance/receivable-projections`

- Resumo: List receivable projections.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/finance/receivable-projections/sync`

- Resumo: Sync projections from sales and rentals.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/finance/commission-holds`

- Resumo: List commission holds.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/finance/commission-holds/{publicId}/release`

- Resumo: Release one commission hold.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/finance/activity`

- Resumo: List finance operational activity.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `fiscal`

- Contrato: `docs/contracts/http/fiscal.openapi.yaml`
- Endpoints: `25`

#### `GET /api/fiscal/capabilities`

- Resumo: Read fiscal capability registry.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/companies/{companyPublicId}/profile`

- Resumo: Read fiscal company profile.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/fiscal/companies/{companyPublicId}/profile`

- Resumo: Upsert fiscal company profile.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/companies/{companyPublicId}/retention-policies`

- Resumo: List retention policies by company.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/companies/{companyPublicId}/retention-execution`

- Resumo: Read retention execution plan for one company.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute`

- Resumo: Execute retention and anonymization plan.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`

- Resumo: Upsert retention policy for one data domain.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/documents`

- Resumo: List fiscal documents.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/fiscal/documents`

- Resumo: Issue one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.

#### `GET /api/fiscal/documents/{publicId}`

- Resumo: Read one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `POST /api/fiscal/documents/{publicId}/cancel`

- Resumo: Cancel one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/fiscal/documents/{publicId}/correction-letter`

- Resumo: Register correction letter for one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/fiscal/documents/{publicId}/invalidate`

- Resumo: Register invalidation for one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/documents/{publicId}/events`

- Resumo: List fiscal document audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/privacy-requests`

- Resumo: List privacy requests.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/fiscal/privacy-requests`

- Resumo: Create privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.

#### `GET /api/fiscal/privacy-requests/{publicId}`

- Resumo: Read one privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `GET /api/fiscal/privacy-requests/{publicId}/export-package`

- Resumo: Build export package for one privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `POST /api/fiscal/privacy-requests/{publicId}/execute`

- Resumo: Execute one privacy request with audit trail.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `PATCH /api/fiscal/privacy-requests/{publicId}/status`

- Resumo: Transition privacy request lifecycle status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `GET /api/fiscal/consents`

- Resumo: List consent ledger.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/fiscal/consents`

- Resumo: Create consent record.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.

#### `PATCH /api/fiscal/consents/{publicId}`

- Resumo: Transition consent status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.

#### `GET /api/fiscal/audit-events`

- Resumo: List fiscal audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/fiscal/compliance/summary`

- Resumo: Read fiscal compliance summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `identity`

- Contrato: `docs/contracts/http/identity.openapi.yaml`
- Endpoints: `6`

#### `GET /api/identity/tenants`

- Resumo: List tenants.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/identity/tenants`

- Resumo: Create tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/identity/tenants/{slug}/snapshot`

- Resumo: Read one tenant snapshot.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/identity/sessions/login`

- Resumo: Authenticate identity session.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/identity/sessions/refresh`

- Resumo: Refresh identity session.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/identity/invitations`

- Resumo: Create invitation.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `notification`

- Contrato: `docs/contracts/http/notification.openapi.yaml`
- Endpoints: `7`

#### `GET /api/notification/capabilities`

- Resumo: Read notification capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/notification/preferences/{userPublicId}`

- Resumo: Read one user notification preference.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/notification/preferences/{userPublicId}`

- Resumo: Upsert one user notification preference.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/notification/center`

- Resumo: List notification center items with cursor filters.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/notification/center`

- Resumo: Create one notification center item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.

#### `PATCH /api/notification/center/{publicId}/status`

- Resumo: Transition notification status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/notification/summary`

- Resumo: Read notification summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `platform-control`

- Contrato: `docs/contracts/http/platform-control.openapi.yaml`
- Endpoints: `40`

#### `GET /api/platform-control/capabilities/catalog`

- Resumo: List platform capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/providers/catalog`

- Resumo: List provider capability catalog and environment posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/entitlements`

- Resumo: List tenant entitlements with cursor pagination.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/feature-flags`

- Resumo: List tenant feature flags with capability metadata.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`

- Resumo: Upsert one entitlement.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}`

- Resumo: Upsert one feature flag using entitlement governance.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`

- Resumo: Bulk upsert entitlements with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`

- Resumo: List provider defaults selected for one tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`

- Resumo: Upsert provider default for one tenant capability.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/quotas`

- Resumo: List quotas by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`

- Resumo: Upsert one quota.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`

- Resumo: Bulk upsert quotas with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/blocks`

- Resumo: List tenant blocks.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`

- Resumo: Upsert tenant block.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/metering`

- Resumo: Read metering snapshots and summary with cursor pagination.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`

- Resumo: Create one usage snapshot.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`

- Resumo: Read quota and metering utilization summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`

- Resumo: Read tenant lifecycle readiness and provider posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`

- Resumo: List onboarding and offboarding jobs with cursor pagination.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`

- Resumo: Read one lifecycle job with audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`

- Resumo: Preview onboarding plan, provider defaults and lifecycle readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`

- Resumo: Queue onboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `202`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`

- Resumo: Preview offboarding plan, retention posture and lifecycle readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`

- Resumo: Queue offboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `202`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`

- Resumo: Transition lifecycle job to running.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`

- Resumo: Transition lifecycle job to completed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`

- Resumo: Transition lifecycle job to failed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`

- Resumo: Transition lifecycle job to cancelled.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness`

- Resumo: Read go-live rollout readiness by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption`

- Resumo: Read tenant go-live adoption baseline and gap.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks`

- Resumo: List go-live bottlenecks and operational blockers.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook`

- Resumo: Read rollout and rollback playbook for one tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments`

- Resumo: List recommended go-live adjustments.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply`

- Resumo: Apply one go-live operational adjustment.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Resumo: List go-live rollouts.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Resumo: Create one go-live rollout.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}`

- Resumo: Read one go-live rollout with events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start`

- Resumo: Transition go-live rollout to running.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete`

- Resumo: Transition go-live rollout to completed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback`

- Resumo: Roll back one go-live rollout.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `rentals`

- Contrato: `docs/contracts/http/rentals.openapi.yaml`
- Endpoints: `4`

#### `GET /api/rentals/contracts`

- Resumo: List rental contracts.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/rentals/contracts`

- Resumo: Create rental contract.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/rentals/contracts/{publicId}/charges`

- Resumo: List contract charges.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`

- Resumo: Update charge status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `sales`

- Contrato: `docs/contracts/http/sales.openapi.yaml`
- Endpoints: `6`

#### `GET /api/sales/opportunities`

- Resumo: List opportunities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/sales/opportunities`

- Resumo: Create opportunity.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/sales/proposals`

- Resumo: List proposals.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/sales/proposals`

- Resumo: Create proposal.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/sales/sales`

- Resumo: List sales.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/sales/invoices`

- Resumo: List commercial invoices.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `simulation`

- Contrato: `docs/contracts/http/simulation.openapi.yaml`
- Endpoints: `3`

#### `GET /api/simulation/scenarios`

- Resumo: List scenarios.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/simulation/scenarios`

- Resumo: Create scenario run.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/simulation/benchmarks/load`

- Resumo: Execute one load benchmark run.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `supplier`

- Contrato: `docs/contracts/http/supplier.openapi.yaml`
- Endpoints: `8`

#### `GET /api/supplier/capabilities`

- Resumo: Read supplier capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/supplier/categories`

- Resumo: List supplier categories.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/supplier/categories/{categoryKey}`

- Resumo: Upsert one supplier category.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/supplier/suppliers`

- Resumo: List suppliers by tenant and status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/supplier/suppliers`

- Resumo: Create one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.

#### `GET /api/supplier/suppliers/summary`

- Resumo: Read supplier summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/supplier/suppliers/{publicId}`

- Resumo: Read one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PATCH /api/supplier/suppliers/{publicId}`

- Resumo: Update one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `support`

- Contrato: `docs/contracts/http/support.openapi.yaml`
- Endpoints: `9`

#### `GET /api/support/capabilities`

- Resumo: Read support capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/support/queues`

- Resumo: List support queues by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PUT /api/support/queues/{queueKey}`

- Resumo: Upsert one support queue.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/support/cases`

- Resumo: List support cases with cursor filters.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/support/cases`

- Resumo: Create one support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.

#### `GET /api/support/cases/summary`

- Resumo: Read support case summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/support/cases/{publicId}`

- Resumo: Read one support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PATCH /api/support/cases/{publicId}/status`

- Resumo: Transition support case status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/support/cases/{publicId}/comments`

- Resumo: Append comment to support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `webhook-hub`

- Contrato: `docs/contracts/http/webhook-hub.openapi.yaml`
- Endpoints: `13`

#### `GET /health/details`

- Resumo: Return readiness details for webhook runtime.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/webhook-hub/capabilities`

- Resumo: Read outbound webhook capability posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/webhook-hub/outbound-endpoints`

- Resumo: List tenant outbound endpoints.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/webhook-hub/outbound-endpoints`

- Resumo: Register one tenant outbound endpoint.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/webhook-hub/outbound-endpoints/{publicId}`

- Resumo: Read one tenant outbound endpoint.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Resumo: List outbound delivery log for one endpoint.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Resumo: Register one outbound delivery attempt.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter`

- Resumo: Move one outbound delivery to dead letter.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `GET /api/webhook-hub/events`

- Resumo: List inbound webhook events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/webhook-hub/events`

- Resumo: Register inbound webhook event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.

#### `GET /api/webhook-hub/events/summary`

- Resumo: Aggregate inbound webhook state.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/webhook-hub/events/{publicId}/dead-letter`

- Resumo: Move event to dead letter queue.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

#### `POST /api/webhook-hub/events/{publicId}/requeue`

- Resumo: Requeue dead-letter event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.

### `workflow-control`

- Contrato: `docs/contracts/http/workflow-control.openapi.yaml`
- Endpoints: `7`

#### `GET /api/workflow-control/definitions`

- Resumo: List workflow definitions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/workflow-control/definitions`

- Resumo: Create workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/workflow-control/definitions/{key}`

- Resumo: Read one workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PATCH /api/workflow-control/definitions/{key}`

- Resumo: Update one workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `PATCH /api/workflow-control/definitions/{key}/status`

- Resumo: Update workflow definition status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/workflow-control/capabilities/triggers`

- Resumo: List workflow trigger catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/workflow-control/capabilities/actions`

- Resumo: List workflow action catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

### `workflow-runtime`

- Contrato: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Endpoints: `6`

#### `GET /api/workflow-runtime/executions`

- Resumo: List workflow executions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/workflow-runtime/executions`

- Resumo: Create workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/workflow-runtime/executions/{publicId}`

- Resumo: Read one workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/workflow-runtime/executions/{publicId}/actions`

- Resumo: List execution action snapshots.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `POST /api/workflow-runtime/executions/{publicId}/advance`

- Resumo: Advance one workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

#### `GET /api/workflow-runtime/capabilities`

- Resumo: List runtime capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.

