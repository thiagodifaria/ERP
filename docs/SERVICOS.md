# SERVICOS

## Objetivo

Concentrar a documentacao dos servicos em um unico lugar. Os diretorios de implementacao ficam livres de README isolado; a referencia funcional, operacional e contratual vive em `docs/`.

Este documento descreve ownership, stack, plano arquitetural, contratos HTTP, responsabilidades de banco, validacao e cuidados de evolucao por servico.

## Numeros atuais

- Servicos com contrato HTTP versionado: `20`
- Endpoints versionados no catalogo HTTP: `201`
- Contratos HTTP: `docs/contracts/http/`
- Contratos de eventos: `docs/contracts/events/`
- Banco por contexto: `service-api/service-postgresql/<contexto>/`
- Comando de runtime: `./scripts/build.sh`
- Comando de validacao: `./scripts/test.sh`

## Inventario executivo

| Servico | Stack | Plano | Caminho | Endpoints | Responsabilidade |
|---------|-------|-------|---------|-----------|------------------|
| `analytics` | Python | analytics plane | `service-api/service-python/analytics` | 9 | relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas |
| `billing` | .NET | transaction/billing plane | `service-api/service-csharp/billing` | 9 | planos, assinaturas, invoices recorrentes, tentativas de pagamento e recovery |
| `catalog` | Python | transaction/catalog plane | `service-api/service-python/catalog` | 12 | categorias, itens, versoes de item, bulk e contratos de consumo |
| `crm` | Go | transaction plane | `service-api/service-golang/crm` | 5 | leads, customers, ownership, pipeline, notas, historico, anexos e enrichment |
| `documents` | Go | transaction plane | `service-api/service-golang/documents` | 10 | anexos, upload, storage posture, assinatura, versoes, archive e access links |
| `edge` | Go | public operations plane | `service-api/service-golang/edge` | 8 | entrada publica, agregacao cross-service e cockpits operacionais |
| `engagement` | TypeScript | interaction/control plane | `service-api/service-typescript/engagement` | 9 | campanhas, templates, touchpoints, conversas, delivery, providers e callbacks |
| `finance` | .NET | transaction/finance plane | `service-api/service-csharp/finance` | 5 | recebiveis, payables, caixa, custos, comissoes e fechamento |
| `fiscal` | Python | compliance plane | `service-api/service-python/fiscal` | 25 | perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria |
| `identity` | .NET | transaction/security plane | `service-api/service-csharp/identity` | 6 | tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria |
| `notification` | Python | administrative/notification plane | `service-api/service-python/notification` | 7 | preferencias, centro interno de alertas, severidade e lifecycle de notificacoes |
| `platform-control` | Python | platform control plane | `service-api/service-python/platform-control` | 40 | capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live |
| `rentals` | Go | transaction plane | `service-api/service-golang/rentals` | 4 | contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais |
| `sales` | Go | transaction plane | `service-api/service-golang/sales` | 6 | oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias |
| `simulation` | Python | simulation plane | `service-api/service-python/simulation` | 3 | cenarios operacionais, benchmark de carga e modelagem de capacidade |
| `supplier` | Python | administrative/procurement plane | `service-api/service-python/supplier` | 8 | categorias de fornecedor, diretorio de fornecedores e procurement ownership |
| `support` | Python | administrative plane | `service-api/service-python/support` | 9 | filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento |
| `webhook-hub` | Rust | integration plane | `service-api/service-rust/webhook-hub` | 13 | intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries |
| `workflow-control` | TypeScript | control plane | `service-api/service-typescript/workflow-control` | 7 | definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos |
| `workflow-runtime` | Elixir | runtime plane | `service-api/service-elixir/workflow-runtime` | 6 | execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes |

## Regras globais de servico

- Cada servico deve ter health live, ready e details quando participa do runtime HTTP.
- Cada servico que expoe superficie publica deve ter OpenAPI em `docs/contracts/http/`.
- Cada contexto persistente deve ter migrations em `service-api/service-postgresql/<contexto>/migrations`.
- Seeds devem existir apenas para bootstrap, smoke e dados de referencia controlados.
- Testes locais ficam junto ao servico; validacoes cross-service ficam em `scripts/test.sh`.
- A documentacao funcional fica aqui e em `docs/API.md`, nao em README dentro do servico.
- Provider externo deve declarar modo configurado, fallback, manual, disabled ou unconfigured.
- Mutacao relevante deve produzir historico, outbox ou evento quando houver consumidor operacional.

## Detalhe por servico

## `analytics`

- Stack: Python
- Plano arquitetural: analytics plane
- Codigo: `service-api/service-python/analytics`
- Contrato HTTP: `docs/contracts/http/analytics.openapi.yaml`
- Titulo OpenAPI: ERP Analytics API
- Versao OpenAPI: `0.1.0`
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.
- Descricao do contrato: Executive reports, adapter catalog, SaaS control and contract governance..
- Contexto PostgreSQL: `service-api/service-postgresql/analytics`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `9`.

### Ownership

- O servico e dono operacional de relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/analytics/reports/adapter-catalog`

- Summary: Read external adapter capability catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/analytics/reports/integration-readiness`

- Summary: Read external integration readiness.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/analytics/reports/saas-control`

- Summary: Read SaaS control posture by tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/analytics/reports/contract-governance`

- Summary: Read contract governance posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/analytics/reports/hardening-review`

- Summary: Read hardening review.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/analytics/reports/core-operations`

- Summary: Read core product operations.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/analytics/reports/relationship-intelligence`

- Summary: Read relationship intelligence.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/analytics/reports/compliance-control`

- Summary: Read fiscal and privacy compliance control.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `GET /api/analytics/reports/go-live-control`

- Summary: Read go-live rollout control.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `analytics`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `billing`

- Stack: .NET
- Plano arquitetural: transaction/billing plane
- Codigo: `service-api/service-csharp/billing`
- Contrato HTTP: `docs/contracts/http/billing.openapi.yaml`
- Titulo OpenAPI: ERP Billing API
- Versao OpenAPI: `0.9.7`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de pagamento e recovery.
- Descricao do contrato: Subscription billing, payment recovery and gateway capabilities..
- Contexto PostgreSQL: `service-api/service-postgresql/billing`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `9`.

### Ownership

- O servico e dono operacional de planos, assinaturas, invoices recorrentes, tentativas de pagamento e recovery.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /health/details`

- Summary: Return readiness details and gateway posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/billing/gateways`

- Summary: List gateway capabilities and Pix posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/billing/gateways/{provider}`

- Summary: Read one gateway capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/billing/plans`

- Summary: List billing plans including flat, hybrid and usage-based pricing.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/billing/plans`

- Summary: Create billing plan.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/billing/subscriptions`

- Summary: List subscriptions.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `POST /api/billing/subscriptions`

- Summary: Create subscription.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/billing/subscriptions/{publicId}/usage-pricing`

- Summary: Project usage-based charge for one subscription.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `POST /api/billing/invoices/{publicId}/attempts`

- Summary: Create payment attempt with idempotency support.
- Parametros: `Idempotency-Key`, `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `billing`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `catalog`

- Stack: Python
- Plano arquitetural: transaction/catalog plane
- Codigo: `service-api/service-python/catalog`
- Contrato HTTP: `docs/contracts/http/catalog.openapi.yaml`
- Titulo OpenAPI: ERP Catalog API
- Versao OpenAPI: `0.2.0`
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.
- Descricao do contrato: Product and service catalog with categories, activation, versioned items, cursor pagination and bulk creation..
- Contexto PostgreSQL: `service-api/service-postgresql/catalog`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `12`.

### Ownership

- O servico e dono operacional de categorias, itens, versoes de item, bulk e contratos de consumo.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/catalog/capabilities`

- Summary: Read catalog capability posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/catalog/consumers`

- Summary: Read catalog consumer contracts across core domains.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/catalog/categories`

- Summary: List categories by tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/catalog/categories`

- Summary: Create one category.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/catalog/categories/page`

- Summary: Cursor-based category listing.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/catalog/items`

- Summary: List catalog items.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `POST /api/catalog/items`

- Summary: Create one catalog item.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/catalog/items/page`

- Summary: Cursor-based item listing.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `POST /api/catalog/items/bulk`

- Summary: Bulk create catalog items with partial success.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 10. `GET /api/catalog/items/{publicId}`

- Summary: Read one catalog item.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 11. `PATCH /api/catalog/items/{publicId}`

- Summary: Update active state, price and attributes.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 12. `GET /api/catalog/items/{publicId}/versions`

- Summary: Read catalog item version history.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `catalog`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `crm`

- Stack: Go
- Plano arquitetural: transaction plane
- Codigo: `service-api/service-golang/crm`
- Contrato HTTP: `docs/contracts/http/crm.openapi.yaml`
- Titulo OpenAPI: ERP CRM API
- Versao OpenAPI: `0.2.0`
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enrichment.
- Descricao do contrato: CRM leads, customers, activity, attachments and commercial intelligence..
- Contexto PostgreSQL: `service-api/service-postgresql/crm`
- Migrations: sim.
- Seeds: sim.
- Endpoints versionados: `5`.

### Ownership

- O servico e dono operacional de leads, customers, ownership, pipeline, notas, historico, anexos e enrichment.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/crm/enrichment/cnpj/capabilities`

- Summary: Read CNPJ enrichment provider capabilities.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `crm`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/crm/enrichment/cnpj/lookup`

- Summary: Lookup and enrich one CNPJ through provider contract.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `crm`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/crm/pipeline/config`

- Summary: Read tenant pipeline configuration.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `crm`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `PUT /api/crm/pipeline/config`

- Summary: Upsert tenant pipeline configuration.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `crm`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/crm/leads/intelligence/summary`

- Summary: Read lead scoring and pipeline intelligence summary.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `crm`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `documents`

- Stack: Go
- Plano arquitetural: transaction plane
- Codigo: `service-api/service-golang/documents`
- Contrato HTTP: `docs/contracts/http/documents.openapi.yaml`
- Titulo OpenAPI: ERP Documents API
- Versao OpenAPI: `0.9.7`
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.
- Descricao do contrato: Attachment governance, upload orchestration and storage capabilities..
- Contexto PostgreSQL: `service-api/service-postgresql/documents`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `10`.

### Ownership

- O servico e dono operacional de anexos, upload, storage posture, assinatura, versoes, archive e access links.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /health/details`

- Summary: Return runtime readiness and storage posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/documents/signing/capabilities`

- Summary: List digital signature capabilities.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/documents/signing/capabilities/{provider}`

- Summary: Read one signing capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/documents/signing/requests`

- Summary: Queue one digital signature request.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/documents/storage/capabilities`

- Summary: List storage capability registry.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/documents/storage/capabilities/{provider}`

- Summary: Read one storage capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/documents/attachments`

- Summary: List attachments.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `POST /api/documents/attachments`

- Summary: Create attachment metadata.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `GET /api/documents/attachments/{publicId}/versions`

- Summary: List attachment versions.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 10. `POST /api/documents/attachments/{publicId}/versions`

- Summary: Append attachment version.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `documents`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `edge`

- Stack: Go
- Plano arquitetural: public operations plane
- Codigo: `service-api/service-golang/edge`
- Contrato HTTP: `docs/contracts/http/edge.openapi.yaml`
- Titulo OpenAPI: ERP Edge API
- Versao OpenAPI: `0.1.0`
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.
- Descricao do contrato: Aggregated operational cockpits for tenants, contracts and SaaS control..
- Contexto PostgreSQL: nao ha pasta dedicada detectada para este servico.
- Migrations: nao detectadas.
- Seeds: nao detectados.
- Endpoints versionados: `8`.

### Ownership

- O servico e dono operacional de entrada publica, agregacao cross-service e cockpits operacionais.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/edge/ops/core-operations`

- Summary: Read executive core product cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/edge/ops/relationship-overview`

- Summary: Read executive relationship cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/edge/ops/compliance-overview`

- Summary: Read executive compliance cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/edge/ops/go-live-overview`

- Summary: Read executive go-live cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/edge/ops/integrations-overview`

- Summary: Read executive integrations cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/edge/ops/saas-overview`

- Summary: Read executive SaaS cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/edge/ops/contracts-overview`

- Summary: Read executive contracts cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/edge/ops/hardening-overview`

- Summary: Read executive hardening cockpit.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `edge`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `engagement`

- Stack: TypeScript
- Plano arquitetural: interaction/control plane
- Codigo: `service-api/service-typescript/engagement`
- Contrato HTTP: `docs/contracts/http/engagement.openapi.yaml`
- Titulo OpenAPI: ERP Engagement API
- Versao OpenAPI: `0.9.7`
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.
- Descricao do contrato: Omnichannel engagement, provider callbacks and campaign operations..
- Contexto PostgreSQL: `service-api/service-postgresql/engagement`
- Migrations: sim.
- Seeds: sim.
- Endpoints versionados: `9`.

### Ownership

- O servico e dono operacional de campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /health/details`

- Summary: Return readiness details for engagement runtime.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/engagement/providers`

- Summary: List provider capabilities and fallback posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/engagement/providers/{provider}`

- Summary: Read one provider capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/engagement/providers/meta-ads/leads`

- Summary: Ingest inbound lead from Meta Ads.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/engagement/providers/resend/events`

- Summary: Register Resend callback event.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `POST /api/engagement/providers/whatsapp-cloud/events`

- Summary: Register WhatsApp callback event.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `POST /api/engagement/providers/telegram-bot/events`

- Summary: Register Telegram callback event.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/engagement/provider-events`

- Summary: List provider events.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `GET /api/engagement/provider-events/{publicId}`

- Summary: Read one provider event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `engagement`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `finance`

- Stack: .NET
- Plano arquitetural: transaction/finance plane
- Codigo: `service-api/service-csharp/finance`
- Contrato HTTP: `docs/contracts/http/finance.openapi.yaml`
- Titulo OpenAPI: ERP Finance API
- Versao OpenAPI: `0.4.0`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.
- Descricao do contrato: Receivables, commission holds, cash control and cross-domain financial activity..
- Contexto PostgreSQL: `service-api/service-postgresql/finance`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `5`.

### Ownership

- O servico e dono operacional de recebiveis, payables, caixa, custos, comissoes e fechamento.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/finance/receivable-projections`

- Summary: List receivable projections.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `finance`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/finance/receivable-projections/sync`

- Summary: Sync projections from sales and rentals.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `finance`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/finance/commission-holds`

- Summary: List commission holds.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `finance`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/finance/commission-holds/{publicId}/release`

- Summary: Release one commission hold.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `finance`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/finance/activity`

- Summary: List finance operational activity.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `finance`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `fiscal`

- Stack: Python
- Plano arquitetural: compliance plane
- Codigo: `service-api/service-python/fiscal`
- Contrato HTTP: `docs/contracts/http/fiscal.openapi.yaml`
- Titulo OpenAPI: ERP Fiscal API
- Versao OpenAPI: `0.1.0`
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.
- Descricao do contrato: Fiscal profile, document operations, privacy rights and compliance governance..
- Contexto PostgreSQL: `service-api/service-postgresql/fiscal`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `25`.

### Ownership

- O servico e dono operacional de perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/fiscal/capabilities`

- Summary: Read fiscal capability registry.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Read fiscal company profile.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `PUT /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Upsert fiscal company profile.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/fiscal/companies/{companyPublicId}/retention-policies`

- Summary: List retention policies by company.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/fiscal/companies/{companyPublicId}/retention-execution`

- Summary: Read retention execution plan for one company.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute`

- Summary: Execute retention and anonymization plan.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`

- Summary: Upsert retention policy for one data domain.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/fiscal/documents`

- Summary: List fiscal documents.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `POST /api/fiscal/documents`

- Summary: Issue one fiscal document.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `201`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 10. `GET /api/fiscal/documents/{publicId}`

- Summary: Read one fiscal document.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 11. `POST /api/fiscal/documents/{publicId}/cancel`

- Summary: Cancel one fiscal document.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 12. `POST /api/fiscal/documents/{publicId}/correction-letter`

- Summary: Register correction letter for one fiscal document.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 13. `POST /api/fiscal/documents/{publicId}/invalidate`

- Summary: Register invalidation for one fiscal document.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 14. `GET /api/fiscal/documents/{publicId}/events`

- Summary: List fiscal document audit events.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 15. `GET /api/fiscal/privacy-requests`

- Summary: List privacy requests.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 16. `POST /api/fiscal/privacy-requests`

- Summary: Create privacy request.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `201`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 17. `GET /api/fiscal/privacy-requests/{publicId}`

- Summary: Read one privacy request.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 18. `GET /api/fiscal/privacy-requests/{publicId}/export-package`

- Summary: Build export package for one privacy request.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 19. `POST /api/fiscal/privacy-requests/{publicId}/execute`

- Summary: Execute one privacy request with audit trail.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 20. `PATCH /api/fiscal/privacy-requests/{publicId}/status`

- Summary: Transition privacy request lifecycle status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 21. `GET /api/fiscal/consents`

- Summary: List consent ledger.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 22. `POST /api/fiscal/consents`

- Summary: Create consent record.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `201`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 23. `PATCH /api/fiscal/consents/{publicId}`

- Summary: Transition consent status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 24. `GET /api/fiscal/audit-events`

- Summary: List fiscal audit events.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 25. `GET /api/fiscal/compliance/summary`

- Summary: Read fiscal compliance summary.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `fiscal`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `identity`

- Stack: .NET
- Plano arquitetural: transaction/security plane
- Codigo: `service-api/service-csharp/identity`
- Contrato HTTP: `docs/contracts/http/identity.openapi.yaml`
- Titulo OpenAPI: ERP Identity API
- Versao OpenAPI: `0.5.0`
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.
- Descricao do contrato: Tenancy, access, sessions, invitations and tenant-scoped identity governance..
- Contexto PostgreSQL: `service-api/service-postgresql/identity`
- Migrations: sim.
- Seeds: sim.
- Endpoints versionados: `6`.

### Ownership

- O servico e dono operacional de tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/identity/tenants`

- Summary: List tenants.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `identity`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/identity/tenants`

- Summary: Create tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `identity`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/identity/tenants/{slug}/snapshot`

- Summary: Read one tenant snapshot.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `identity`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/identity/sessions/login`

- Summary: Authenticate identity session.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `identity`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/identity/sessions/refresh`

- Summary: Refresh identity session.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `identity`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `POST /api/identity/invitations`

- Summary: Create invitation.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `identity`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `notification`

- Stack: Python
- Plano arquitetural: administrative/notification plane
- Codigo: `service-api/service-python/notification`
- Contrato HTTP: `docs/contracts/http/notification.openapi.yaml`
- Titulo OpenAPI: ERP Notification API
- Versao OpenAPI: `0.1.0`
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.
- Descricao do contrato: Internal notification center, preferences and reusable dispatch shape..
- Contexto PostgreSQL: `service-api/service-postgresql/notification`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `7`.

### Ownership

- O servico e dono operacional de preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/notification/capabilities`

- Summary: Read notification capability catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/notification/preferences/{userPublicId}`

- Summary: Read one user notification preference.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `PUT /api/notification/preferences/{userPublicId}`

- Summary: Upsert one user notification preference.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/notification/center`

- Summary: List notification center items with cursor filters.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/notification/center`

- Summary: Create one notification center item.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `201`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `PATCH /api/notification/center/{publicId}/status`

- Summary: Transition notification status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/notification/summary`

- Summary: Read notification summary.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `notification`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `platform-control`

- Stack: Python
- Plano arquitetural: platform control plane
- Codigo: `service-api/service-python/platform-control`
- Contrato HTTP: `docs/contracts/http/platform-control.openapi.yaml`
- Titulo OpenAPI: ERP Platform Control API
- Versao OpenAPI: `0.2.0`
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.
- Descricao do contrato: Tenant capabilities, entitlements, quotas, metering, lifecycle jobs and SaaS governance..
- Contexto PostgreSQL: `service-api/service-postgresql/platform-control`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `40`.

### Ownership

- O servico e dono operacional de capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/platform-control/capabilities/catalog`

- Summary: List platform capability catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/platform-control/providers/catalog`

- Summary: List provider capability catalog and environment posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/platform-control/tenants/{tenantSlug}/entitlements`

- Summary: List tenant entitlements with cursor pagination.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/platform-control/tenants/{tenantSlug}/feature-flags`

- Summary: List tenant feature flags with capability metadata.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`

- Summary: Upsert one entitlement.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}`

- Summary: Upsert one feature flag using entitlement governance.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`

- Summary: Bulk upsert entitlements with partial success.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`

- Summary: List provider defaults selected for one tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`

- Summary: Upsert provider default for one tenant capability.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 10. `GET /api/platform-control/tenants/{tenantSlug}/quotas`

- Summary: List quotas by tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 11. `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`

- Summary: Upsert one quota.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 12. `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`

- Summary: Bulk upsert quotas with partial success.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 13. `GET /api/platform-control/tenants/{tenantSlug}/blocks`

- Summary: List tenant blocks.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 14. `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`

- Summary: Upsert tenant block.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 15. `GET /api/platform-control/tenants/{tenantSlug}/metering`

- Summary: Read metering snapshots and summary with cursor pagination.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 16. `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`

- Summary: Create one usage snapshot.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 17. `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`

- Summary: Read quota and metering utilization summary.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 18. `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`

- Summary: Read tenant lifecycle readiness and provider posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 19. `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`

- Summary: List onboarding and offboarding jobs with cursor pagination.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 20. `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`

- Summary: Read one lifecycle job with audit events.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 21. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`

- Summary: Preview onboarding plan, provider defaults and lifecycle readiness.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 22. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`

- Summary: Queue onboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `202`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 23. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`

- Summary: Preview offboarding plan, retention posture and lifecycle readiness.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 24. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`

- Summary: Queue offboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `202`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 25. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`

- Summary: Transition lifecycle job to running.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 26. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`

- Summary: Transition lifecycle job to completed.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 27. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`

- Summary: Transition lifecycle job to failed.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 28. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`

- Summary: Transition lifecycle job to cancelled.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 29. `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness`

- Summary: Read go-live rollout readiness by tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 30. `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption`

- Summary: Read tenant go-live adoption baseline and gap.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 31. `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks`

- Summary: List go-live bottlenecks and operational blockers.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 32. `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook`

- Summary: Read rollout and rollback playbook for one tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 33. `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments`

- Summary: List recommended go-live adjustments.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 34. `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply`

- Summary: Apply one go-live operational adjustment.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 35. `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: List go-live rollouts.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 36. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: Create one go-live rollout.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 37. `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}`

- Summary: Read one go-live rollout with events.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 38. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start`

- Summary: Transition go-live rollout to running.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 39. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete`

- Summary: Transition go-live rollout to completed.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 40. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback`

- Summary: Roll back one go-live rollout.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `platform-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `rentals`

- Stack: Go
- Plano arquitetural: transaction plane
- Codigo: `service-api/service-golang/rentals`
- Contrato HTTP: `docs/contracts/http/rentals.openapi.yaml`
- Titulo OpenAPI: ERP Rentals API
- Versao OpenAPI: `0.8.0`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.
- Descricao do contrato: Rental contracts, recurring charges, adjustments, rescission and attachment linkage..
- Contexto PostgreSQL: `service-api/service-postgresql/rentals`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `4`.

### Ownership

- O servico e dono operacional de contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/rentals/contracts`

- Summary: List rental contracts.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `rentals`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/rentals/contracts`

- Summary: Create rental contract.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `rentals`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/rentals/contracts/{publicId}/charges`

- Summary: List contract charges.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `rentals`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`

- Summary: Update charge status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `rentals`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `sales`

- Stack: Go
- Plano arquitetural: transaction plane
- Codigo: `service-api/service-golang/sales`
- Contrato HTTP: `docs/contracts/http/sales.openapi.yaml`
- Titulo OpenAPI: ERP Sales API
- Versao OpenAPI: `0.7.0`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.
- Descricao do contrato: Opportunities, proposals, sales, invoices and commercial lifecycle control..
- Contexto PostgreSQL: `service-api/service-postgresql/sales`
- Migrations: sim.
- Seeds: sim.
- Endpoints versionados: `6`.

### Ownership

- O servico e dono operacional de oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/sales/opportunities`

- Summary: List opportunities.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `sales`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/sales/opportunities`

- Summary: Create opportunity.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `sales`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/sales/proposals`

- Summary: List proposals.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `sales`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/sales/proposals`

- Summary: Create proposal.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `sales`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/sales/sales`

- Summary: List sales.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `sales`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/sales/invoices`

- Summary: List commercial invoices.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `sales`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `simulation`

- Stack: Python
- Plano arquitetural: simulation plane
- Codigo: `service-api/service-python/simulation`
- Contrato HTTP: `docs/contracts/http/simulation.openapi.yaml`
- Titulo OpenAPI: ERP Simulation API
- Versao OpenAPI: `0.7.0`
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.
- Descricao do contrato: Scenario simulation, load benchmark and cost estimation runtime..
- Contexto PostgreSQL: `service-api/service-postgresql/simulation`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `3`.

### Ownership

- O servico e dono operacional de cenarios operacionais, benchmark de carga e modelagem de capacidade.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/simulation/scenarios`

- Summary: List scenarios.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `simulation`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/simulation/scenarios`

- Summary: Create scenario run.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `simulation`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `POST /api/simulation/benchmarks/load`

- Summary: Execute one load benchmark run.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `simulation`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `supplier`

- Stack: Python
- Plano arquitetural: administrative/procurement plane
- Codigo: `service-api/service-python/supplier`
- Contrato HTTP: `docs/contracts/http/supplier.openapi.yaml`
- Titulo OpenAPI: ERP Supplier API
- Versao OpenAPI: `0.1.0`
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.
- Descricao do contrato: Supplier directory, categories and payables-oriented vendor governance..
- Contexto PostgreSQL: `service-api/service-postgresql/supplier`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `8`.

### Ownership

- O servico e dono operacional de categorias de fornecedor, diretorio de fornecedores e procurement ownership.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/supplier/capabilities`

- Summary: Read supplier capability catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/supplier/categories`

- Summary: List supplier categories.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `PUT /api/supplier/categories/{categoryKey}`

- Summary: Upsert one supplier category.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/supplier/suppliers`

- Summary: List suppliers by tenant and status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/supplier/suppliers`

- Summary: Create one supplier.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `201`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/supplier/suppliers/summary`

- Summary: Read supplier summary.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/supplier/suppliers/{publicId}`

- Summary: Read one supplier.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `PATCH /api/supplier/suppliers/{publicId}`

- Summary: Update one supplier.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `supplier`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `support`

- Stack: Python
- Plano arquitetural: administrative plane
- Codigo: `service-api/service-python/support`
- Contrato HTTP: `docs/contracts/http/support.openapi.yaml`
- Titulo OpenAPI: ERP Support API
- Versao OpenAPI: `0.1.0`
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.
- Descricao do contrato: Queue-based support cases with SLA, comments and lifecycle history..
- Contexto PostgreSQL: `service-api/service-postgresql/support`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `9`.

### Ownership

- O servico e dono operacional de filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/support/capabilities`

- Summary: Read support capability catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/support/queues`

- Summary: List support queues by tenant.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `PUT /api/support/queues/{queueKey}`

- Summary: Upsert one support queue.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/support/cases`

- Summary: List support cases with cursor filters.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/support/cases`

- Summary: Create one support case.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `201`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/support/cases/summary`

- Summary: Read support case summary.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/support/cases/{publicId}`

- Summary: Read one support case.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `PATCH /api/support/cases/{publicId}/status`

- Summary: Transition support case status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `POST /api/support/cases/{publicId}/comments`

- Summary: Append comment to support case.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `support`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `webhook-hub`

- Stack: Rust
- Plano arquitetural: integration plane
- Codigo: `service-api/service-rust/webhook-hub`
- Contrato HTTP: `docs/contracts/http/webhook-hub.openapi.yaml`
- Titulo OpenAPI: ERP Webhook Hub API
- Versao OpenAPI: `0.9.7`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.
- Descricao do contrato: Inbound webhook intake, DLQ and operator recovery surface..
- Contexto PostgreSQL: `service-api/service-postgresql/webhook-hub`
- Migrations: sim.
- Seeds: sim.
- Endpoints versionados: `13`.

### Ownership

- O servico e dono operacional de intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /health/details`

- Summary: Return readiness details for webhook runtime.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `GET /api/webhook-hub/capabilities`

- Summary: Read outbound webhook capability posture.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/webhook-hub/outbound-endpoints`

- Summary: List tenant outbound endpoints.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `POST /api/webhook-hub/outbound-endpoints`

- Summary: Register one tenant outbound endpoint.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `GET /api/webhook-hub/outbound-endpoints/{publicId}`

- Summary: Read one tenant outbound endpoint.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: List outbound delivery log for one endpoint.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: Register one outbound delivery attempt.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 8. `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter`

- Summary: Move one outbound delivery to dead letter.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 9. `GET /api/webhook-hub/events`

- Summary: List inbound webhook events.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 10. `POST /api/webhook-hub/events`

- Summary: Register inbound webhook event.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: sim.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 11. `GET /api/webhook-hub/events/summary`

- Summary: Aggregate inbound webhook state.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 12. `POST /api/webhook-hub/events/{publicId}/dead-letter`

- Summary: Move event to dead letter queue.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 13. `POST /api/webhook-hub/events/{publicId}/requeue`

- Summary: Requeue dead-letter event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada no OpenAPI atual.
- Dono: `webhook-hub`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `workflow-control`

- Stack: TypeScript
- Plano arquitetural: control plane
- Codigo: `service-api/service-typescript/workflow-control`
- Contrato HTTP: `docs/contracts/http/workflow-control.openapi.yaml`
- Titulo OpenAPI: ERP Workflow Control API
- Versao OpenAPI: `0.6.0`
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.
- Descricao do contrato: Workflow definition catalog, status, publication and action snapshots..
- Contexto PostgreSQL: `service-api/service-postgresql/workflow-control`
- Migrations: sim.
- Seeds: sim.
- Endpoints versionados: `7`.

### Ownership

- O servico e dono operacional de definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/workflow-control/definitions`

- Summary: List workflow definitions.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/workflow-control/definitions`

- Summary: Create workflow definition.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/workflow-control/definitions/{key}`

- Summary: Read one workflow definition.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `PATCH /api/workflow-control/definitions/{key}`

- Summary: Update one workflow definition.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `PATCH /api/workflow-control/definitions/{key}/status`

- Summary: Update workflow definition status.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/workflow-control/capabilities/triggers`

- Summary: List workflow trigger catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 7. `GET /api/workflow-control/capabilities/actions`

- Summary: List workflow action catalog.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-control`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.

## `workflow-runtime`

- Stack: Elixir
- Plano arquitetural: runtime plane
- Codigo: `service-api/service-elixir/workflow-runtime`
- Contrato HTTP: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Titulo OpenAPI: ERP Workflow Runtime API
- Versao OpenAPI: `0.6.0`
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.
- Descricao do contrato: Workflow execution runtime with action snapshots, retries, delays and compensations..
- Contexto PostgreSQL: `service-api/service-postgresql/workflow-runtime`
- Migrations: sim.
- Seeds: nao detectados.
- Endpoints versionados: `6`.

### Ownership

- O servico e dono operacional de execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.
- Regras de negocio devem ficar em domain/application, nao em bootstrap HTTP.
- Repositorios e providers devem ficar atras de adapters de infrastructure.
- Health details devem mostrar dependencia ativa, modo de repository driver e postura de provider quando aplicavel.

### Contratos HTTP

#### 1. `GET /api/workflow-runtime/executions`

- Summary: List workflow executions.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-runtime`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 2. `POST /api/workflow-runtime/executions`

- Summary: Create workflow execution.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-runtime`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 3. `GET /api/workflow-runtime/executions/{publicId}`

- Summary: Read one workflow execution.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-runtime`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 4. `GET /api/workflow-runtime/executions/{publicId}/actions`

- Summary: List execution action snapshots.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-runtime`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 5. `POST /api/workflow-runtime/executions/{publicId}/advance`

- Summary: Advance one workflow execution.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-runtime`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

#### 6. `GET /api/workflow-runtime/capabilities`

- Summary: List runtime capabilities.
- Parametros: nenhum parametro declarado no OpenAPI atual.
- Request body: nao declarado.
- Respostas: `200`.
- Dono: `workflow-runtime`.
- Verificacao: manter contrato sincronizado com implementacao, testes e smoke quando alterar semantica.
- Observabilidade: registrar correlation id, tenant quando aplicavel e status final.

### Validacao

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

### Evolucao segura

- Criar rota nova junto com contrato.
- Evitar mudanca breaking em payload publico sem changelog.
- Evitar escrever em schema de outro contexto diretamente.
- Atualizar `docs/API.md` indiretamente pelo contrato OpenAPI quando a superficie mudar.
