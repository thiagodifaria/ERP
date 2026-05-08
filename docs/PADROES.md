# PADROES

## Objetivo

Definir padroes de engenharia para manter o ERP consistente apesar de ser grande, poliglota e distribuido em varios dominios.

## Organizacao

- README raiz e curto e serve como entrada.
- `README_EN.md` e `README_PT.md` sao detalhados.
- Documentacao central fica em `docs/`.
- Contratos ficam em `docs/contracts/`.
- Codigo fica em `service-api/`.
- Runtime fica em `infra/`.
- Scripts oficiais sao `scripts/build.sh` e `scripts/test.sh`.
- Testes ficam nos modulos que exercitam.

## Padroes globais

- tenant explicito;
- public id em rotas publicas;
- erro com `code`, `message` e `details`;
- health live, ready e details;
- idempotencia em mutacao sensivel;
- async com `202 Accepted` quando a operacao nao termina na requisicao;
- cursor pagination para volume transacional;
- bulk com partial success;
- adapter para provider externo;
- migration por contexto;
- smoke para fluxo cross-service.

## Padroes por stack

## .NET

Servicos:

- `billing`
- `finance`
- `identity`

Regras:

- respeitar idiomatica da stack;
- manter bootstrap fino;
- separar API, application, domain, infrastructure e config quando a stack suportar;
- usar testes de unidade para regra local;
- usar contrato e smoke para superficie publica;

## Elixir

Servicos:

- `workflow-runtime`

Regras:

- respeitar idiomatica da stack;
- manter bootstrap fino;
- separar API, application, domain, infrastructure e config quando a stack suportar;
- usar testes de unidade para regra local;
- usar contrato e smoke para superficie publica;

## Go

Servicos:

- `crm`
- `documents`
- `edge`
- `rentals`
- `sales`

Regras:

- respeitar idiomatica da stack;
- manter bootstrap fino;
- separar API, application, domain, infrastructure e config quando a stack suportar;
- usar testes de unidade para regra local;
- usar contrato e smoke para superficie publica;

## Python

Servicos:

- `analytics`
- `catalog`
- `fiscal`
- `notification`
- `platform-control`
- `simulation`
- `supplier`
- `support`

Regras:

- respeitar idiomatica da stack;
- manter bootstrap fino;
- separar API, application, domain, infrastructure e config quando a stack suportar;
- usar testes de unidade para regra local;
- usar contrato e smoke para superficie publica;

## Rust

Servicos:

- `webhook-hub`

Regras:

- respeitar idiomatica da stack;
- manter bootstrap fino;
- separar API, application, domain, infrastructure e config quando a stack suportar;
- usar testes de unidade para regra local;
- usar contrato e smoke para superficie publica;

## TypeScript

Servicos:

- `engagement`
- `workflow-control`

Regras:

- respeitar idiomatica da stack;
- manter bootstrap fino;
- separar API, application, domain, infrastructure e config quando a stack suportar;
- usar testes de unidade para regra local;
- usar contrato e smoke para superficie publica;

## Padroes por servico

## `analytics`

- Stack: Python
- Codigo: `service-api/service-python/analytics`
- Contrato: `docs/contracts/http/analytics.openapi.yaml`
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/analytics/reports/adapter-catalog` - Read external adapter capability catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/integration-readiness` - Read external integration readiness
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/saas-control` - Read SaaS control posture by tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/contract-governance` - Read contract governance posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/hardening-review` - Read hardening review
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/core-operations` - Read core product operations
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/relationship-intelligence` - Read relationship intelligence
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/compliance-control` - Read fiscal and privacy compliance control
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/analytics/reports/go-live-control` - Read go-live rollout control
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `billing`

- Stack: .NET
- Codigo: `service-api/service-csharp/billing`
- Contrato: `docs/contracts/http/billing.openapi.yaml`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /health/details` - Return readiness details and gateway posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/billing/gateways` - List gateway capabilities and Pix posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/billing/gateways/{provider}` - Read one gateway capability
  - Metodo: `GET`.
  - Parametros: `provider`.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/billing/plans` - List billing plans including flat, hybrid and usage-based pricing
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/billing/plans` - Create billing plan
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/billing/subscriptions` - List subscriptions
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/billing/subscriptions` - Create subscription
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/billing/subscriptions/{publicId}/usage-pricing` - Project usage-based charge for one subscription
  - Metodo: `GET`.
  - Parametros: `publicId`.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/billing/invoices/{publicId}/attempts` - Create payment attempt with idempotency support
  - Metodo: `POST`.
  - Parametros: `Idempotency-Key`, `publicId`.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `catalog`

- Stack: Python
- Codigo: `service-api/service-python/catalog`
- Contrato: `docs/contracts/http/catalog.openapi.yaml`
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/catalog/capabilities` - Read catalog capability posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/catalog/consumers` - Read catalog consumer contracts across core domains
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/catalog/categories` - List categories by tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/catalog/categories` - Create one category
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/catalog/categories/page` - Cursor-based category listing
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/catalog/items` - List catalog items
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/catalog/items` - Create one catalog item
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/catalog/items/page` - Cursor-based item listing
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/catalog/items/bulk` - Bulk create catalog items with partial success
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/catalog/items/{publicId}` - Read one catalog item
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PATCH /api/catalog/items/{publicId}` - Update active state, price and attributes
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/catalog/items/{publicId}/versions` - Read catalog item version history
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `crm`

- Stack: Go
- Codigo: `service-api/service-golang/crm`
- Contrato: `docs/contracts/http/crm.openapi.yaml`
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/crm/enrichment/cnpj/capabilities` - Read CNPJ enrichment provider capabilities
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/crm/enrichment/cnpj/lookup` - Lookup and enrich one CNPJ through provider contract
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/crm/pipeline/config` - Read tenant pipeline configuration
  - Metodo: `GET`.
  - Parametros: `tenantSlug`.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/crm/pipeline/config` - Upsert tenant pipeline configuration
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/crm/leads/intelligence/summary` - Read lead scoring and pipeline intelligence summary
  - Metodo: `GET`.
  - Parametros: `tenantSlug`.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `documents`

- Stack: Go
- Codigo: `service-api/service-golang/documents`
- Contrato: `docs/contracts/http/documents.openapi.yaml`
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /health/details` - Return runtime readiness and storage posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/documents/signing/capabilities` - List digital signature capabilities
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/documents/signing/capabilities/{provider}` - Read one signing capability
  - Metodo: `GET`.
  - Parametros: `provider`.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/documents/signing/requests` - Queue one digital signature request
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/documents/storage/capabilities` - List storage capability registry
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/documents/storage/capabilities/{provider}` - Read one storage capability
  - Metodo: `GET`.
  - Parametros: `provider`.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/documents/attachments` - List attachments
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/documents/attachments` - Create attachment metadata
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/documents/attachments/{publicId}/versions` - List attachment versions
  - Metodo: `GET`.
  - Parametros: `publicId`.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/documents/attachments/{publicId}/versions` - Append attachment version
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `edge`

- Stack: Go
- Codigo: `service-api/service-golang/edge`
- Contrato: `docs/contracts/http/edge.openapi.yaml`
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/edge/ops/core-operations` - Read executive core product cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/relationship-overview` - Read executive relationship cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/compliance-overview` - Read executive compliance cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/go-live-overview` - Read executive go-live cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/integrations-overview` - Read executive integrations cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/saas-overview` - Read executive SaaS cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/contracts-overview` - Read executive contracts cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/edge/ops/hardening-overview` - Read executive hardening cockpit
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `engagement`

- Stack: TypeScript
- Codigo: `service-api/service-typescript/engagement`
- Contrato: `docs/contracts/http/engagement.openapi.yaml`
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /health/details` - Return readiness details for engagement runtime
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/engagement/providers` - List provider capabilities and fallback posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/engagement/providers/{provider}` - Read one provider capability
  - Metodo: `GET`.
  - Parametros: `provider`.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/engagement/providers/meta-ads/leads` - Ingest inbound lead from Meta Ads
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/engagement/providers/resend/events` - Register Resend callback event
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/engagement/providers/whatsapp-cloud/events` - Register WhatsApp callback event
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/engagement/providers/telegram-bot/events` - Register Telegram callback event
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/engagement/provider-events` - List provider events
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/engagement/provider-events/{publicId}` - Read one provider event
  - Metodo: `GET`.
  - Parametros: `publicId`.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `finance`

- Stack: .NET
- Codigo: `service-api/service-csharp/finance`
- Contrato: `docs/contracts/http/finance.openapi.yaml`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/finance/receivable-projections` - List receivable projections
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/finance/receivable-projections/sync` - Sync projections from sales and rentals
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/finance/commission-holds` - List commission holds
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/finance/commission-holds/{publicId}/release` - Release one commission hold
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/finance/activity` - List finance operational activity
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `fiscal`

- Stack: Python
- Codigo: `service-api/service-python/fiscal`
- Contrato: `docs/contracts/http/fiscal.openapi.yaml`
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/fiscal/capabilities` - Read fiscal capability registry
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/fiscal/companies/{companyPublicId}/profile` - Read fiscal company profile
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/fiscal/companies/{companyPublicId}/profile` - Upsert fiscal company profile
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies` - List retention policies by company
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/fiscal/companies/{companyPublicId}/retention-execution` - Read retention execution plan for one company
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute` - Execute retention and anonymization plan
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}` - Upsert retention policy for one data domain
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/documents` - List fiscal documents
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/fiscal/documents` - Issue one fiscal document
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/documents/{publicId}` - Read one fiscal document
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/fiscal/documents/{publicId}/cancel` - Cancel one fiscal document
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/fiscal/documents/{publicId}/correction-letter` - Register correction letter for one fiscal document
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/fiscal/documents/{publicId}/invalidate` - Register invalidation for one fiscal document
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/documents/{publicId}/events` - List fiscal document audit events
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/fiscal/privacy-requests` - List privacy requests
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/fiscal/privacy-requests` - Create privacy request
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/privacy-requests/{publicId}` - Read one privacy request
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/fiscal/privacy-requests/{publicId}/export-package` - Build export package for one privacy request
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/fiscal/privacy-requests/{publicId}/execute` - Execute one privacy request with audit trail
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `PATCH /api/fiscal/privacy-requests/{publicId}/status` - Transition privacy request lifecycle status
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/consents` - List consent ledger
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/fiscal/consents` - Create consent record
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `PATCH /api/fiscal/consents/{publicId}` - Transition consent status
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/fiscal/audit-events` - List fiscal audit events
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/fiscal/compliance/summary` - Read fiscal compliance summary
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `identity`

- Stack: .NET
- Codigo: `service-api/service-csharp/identity`
- Contrato: `docs/contracts/http/identity.openapi.yaml`
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/identity/tenants` - List tenants
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/identity/tenants` - Create tenant
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/identity/tenants/{slug}/snapshot` - Read one tenant snapshot
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/identity/sessions/login` - Authenticate identity session
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/identity/sessions/refresh` - Refresh identity session
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/identity/invitations` - Create invitation
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `notification`

- Stack: Python
- Codigo: `service-api/service-python/notification`
- Contrato: `docs/contracts/http/notification.openapi.yaml`
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/notification/capabilities` - Read notification capability catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/notification/preferences/{userPublicId}` - Read one user notification preference
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/notification/preferences/{userPublicId}` - Upsert one user notification preference
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/notification/center` - List notification center items with cursor filters
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/notification/center` - Create one notification center item
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `PATCH /api/notification/center/{publicId}/status` - Transition notification status
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/notification/summary` - Read notification summary
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `platform-control`

- Stack: Python
- Codigo: `service-api/service-python/platform-control`
- Contrato: `docs/contracts/http/platform-control.openapi.yaml`
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/platform-control/capabilities/catalog` - List platform capability catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/providers/catalog` - List provider capability catalog and environment posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements` - List tenant entitlements with cursor pagination
  - Metodo: `GET`.
  - Parametros: `tenantSlug`.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/feature-flags` - List tenant feature flags with capability metadata
  - Metodo: `GET`.
  - Parametros: `tenantSlug`.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}` - Upsert one entitlement
  - Metodo: `PUT`.
  - Parametros: `capabilityKey`, `tenantSlug`.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}` - Upsert one feature flag using entitlement governance
  - Metodo: `PUT`.
  - Parametros: `capabilityKey`, `tenantSlug`.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk` - Bulk upsert entitlements with partial success
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults` - List provider defaults selected for one tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}` - Upsert provider default for one tenant capability
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/quotas` - List quotas by tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}` - Upsert one quota
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk` - Bulk upsert quotas with partial success
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/blocks` - List tenant blocks
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}` - Upsert tenant block
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/metering` - Read metering snapshots and summary with cursor pagination
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots` - Create one usage snapshot
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/usage-summary` - Read quota and metering utilization summary
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness` - Read tenant lifecycle readiness and provider posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs` - List onboarding and offboarding jobs with cursor pagination
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}` - Read one lifecycle job with audit events
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview` - Preview onboarding plan, provider defaults and lifecycle readiness
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding` - Queue onboarding job with Idempotency-Key and 202 Accepted
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview` - Preview offboarding plan, retention posture and lifecycle readiness
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding` - Queue offboarding job with Idempotency-Key and 202 Accepted
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start` - Transition lifecycle job to running
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete` - Transition lifecycle job to completed
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail` - Transition lifecycle job to failed
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel` - Transition lifecycle job to cancelled
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` - Read go-live rollout readiness by tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` - Read tenant go-live adoption baseline and gap
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` - List go-live bottlenecks and operational blockers
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` - Read rollout and rollback playbook for one tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` - List recommended go-live adjustments
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply` - Apply one go-live operational adjustment
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - List go-live rollouts
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - Create one go-live rollout
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}` - Read one go-live rollout with events
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` - Transition go-live rollout to running
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` - Transition go-live rollout to completed
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` - Roll back one go-live rollout
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `rentals`

- Stack: Go
- Codigo: `service-api/service-golang/rentals`
- Contrato: `docs/contracts/http/rentals.openapi.yaml`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/rentals/contracts` - List rental contracts
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/rentals/contracts` - Create rental contract
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/rentals/contracts/{publicId}/charges` - List contract charges
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` - Update charge status
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `sales`

- Stack: Go
- Codigo: `service-api/service-golang/sales`
- Contrato: `docs/contracts/http/sales.openapi.yaml`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/sales/opportunities` - List opportunities
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/sales/opportunities` - Create opportunity
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/sales/proposals` - List proposals
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/sales/proposals` - Create proposal
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/sales/sales` - List sales
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/sales/invoices` - List commercial invoices
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `simulation`

- Stack: Python
- Codigo: `service-api/service-python/simulation`
- Contrato: `docs/contracts/http/simulation.openapi.yaml`
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/simulation/scenarios` - List scenarios
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/simulation/scenarios` - Create scenario run
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/simulation/benchmarks/load` - Execute one load benchmark run
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `supplier`

- Stack: Python
- Codigo: `service-api/service-python/supplier`
- Contrato: `docs/contracts/http/supplier.openapi.yaml`
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/supplier/capabilities` - Read supplier capability catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/supplier/categories` - List supplier categories
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/supplier/categories/{categoryKey}` - Upsert one supplier category
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/supplier/suppliers` - List suppliers by tenant and status
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/supplier/suppliers` - Create one supplier
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/supplier/suppliers/summary` - Read supplier summary
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/supplier/suppliers/{publicId}` - Read one supplier
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PATCH /api/supplier/suppliers/{publicId}` - Update one supplier
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `support`

- Stack: Python
- Codigo: `service-api/service-python/support`
- Contrato: `docs/contracts/http/support.openapi.yaml`
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/support/capabilities` - Read support capability catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/support/queues` - List support queues by tenant
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PUT /api/support/queues/{queueKey}` - Upsert one support queue
  - Metodo: `PUT`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/support/cases` - List support cases with cursor filters
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/support/cases` - Create one support case
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/support/cases/summary` - Read support case summary
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/support/cases/{publicId}` - Read one support case
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PATCH /api/support/cases/{publicId}/status` - Transition support case status
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/support/cases/{publicId}/comments` - Append comment to support case
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `webhook-hub`

- Stack: Rust
- Codigo: `service-api/service-rust/webhook-hub`
- Contrato: `docs/contracts/http/webhook-hub.openapi.yaml`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /health/details` - Return readiness details for webhook runtime
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/webhook-hub/capabilities` - Read outbound webhook capability posture
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/webhook-hub/outbound-endpoints` - List tenant outbound endpoints
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/webhook-hub/outbound-endpoints` - Register one tenant outbound endpoint
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/webhook-hub/outbound-endpoints/{publicId}` - Read one tenant outbound endpoint
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - List outbound delivery log for one endpoint
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - Register one outbound delivery attempt
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter` - Move one outbound delivery to dead letter
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/webhook-hub/events` - List inbound webhook events
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/webhook-hub/events` - Register inbound webhook event
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/webhook-hub/events/summary` - Aggregate inbound webhook state
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/webhook-hub/events/{publicId}/dead-letter` - Move event to dead letter queue
  - Metodo: `POST`.
  - Parametros: `publicId`.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `POST /api/webhook-hub/events/{publicId}/requeue` - Requeue dead-letter event
  - Metodo: `POST`.
  - Parametros: `publicId`.
  - Padrao: mutacao idempotente ou com transicao auditavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `workflow-control`

- Stack: TypeScript
- Codigo: `service-api/service-typescript/workflow-control`
- Contrato: `docs/contracts/http/workflow-control.openapi.yaml`
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/workflow-control/definitions` - List workflow definitions
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/workflow-control/definitions` - Create workflow definition
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/workflow-control/definitions/{key}` - Read one workflow definition
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `PATCH /api/workflow-control/definitions/{key}` - Update one workflow definition
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `PATCH /api/workflow-control/definitions/{key}/status` - Update workflow definition status
  - Metodo: `PATCH`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/workflow-control/capabilities/triggers` - List workflow trigger catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/workflow-control/capabilities/actions` - List workflow action catalog
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## `workflow-runtime`

- Stack: Elixir
- Codigo: `service-api/service-elixir/workflow-runtime`
- Contrato: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.

### Padroes obrigatorios

- regra de negocio fora do bootstrap;
- DTO sem regra pesada;
- repository ou adapter isolado;
- health refletindo dependencia real;
- contrato atualizado para rota publica;
- teste proporcional ao risco;
- logs com correlation id;
- tenant em operacoes tenant-aware;

### Rotas e padroes esperados

- `GET /api/workflow-runtime/executions` - List workflow executions
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/workflow-runtime/executions` - Create workflow execution
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/workflow-runtime/executions/{publicId}` - Read one workflow execution
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `GET /api/workflow-runtime/executions/{publicId}/actions` - List execution action snapshots
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.
- `POST /api/workflow-runtime/executions/{publicId}/advance` - Advance one workflow execution
  - Metodo: `POST`.
  - Parametros: nenhum parametro declarado.
  - Padrao: mutacao idempotente ou com transicao auditavel.
- `GET /api/workflow-runtime/capabilities` - List runtime capabilities
  - Metodo: `GET`.
  - Parametros: nenhum parametro declarado.
  - Padrao: leitura barata, filtravel e observavel.

### Checklist de PR

- O contrato mudou junto com a implementacao?
- O teste cobre o comportamento relevante?
- O erro publico tem codigo estavel?
- O tenant esta preservado?
- A documentacao central continua coerente?

## Padrao de comentarios em codigo

Comentarios devem parecer escritos durante desenvolvimento real: claros, diretos e em primeira pessoa operacional quando fizer sentido. Evitar comentario cerimonial, terceira pessoa artificial e explicacao obvia. Comentario bom explica decisao, restricao ou risco; comentario ruim narra uma linha que o codigo ja mostra.
