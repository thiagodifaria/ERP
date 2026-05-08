# INTEGRACOES

## Objetivo

Mapear integracoes internas, externas, eventos, providers, callbacks, webhooks e leitura cross-context. O objetivo e impedir que integracao vire acoplamento invisivel.

## Principios

- Provider externo sempre passa por adapter ou servico de fronteira.
- Webhook critico entra por `webhook-hub`.
- Evento compartilhado usa schema em `docs/contracts/events/`.
- Superficie HTTP compartilhada usa OpenAPI em `docs/contracts/http/`.
- Consumidor conhece contrato, nao tabela.
- Analytics agrega leitura, nao assume ownership transacional.

## Postura de provider

- `configured`: credencial presente e adapter pronto.
- `fallback`: caminho local ou simulado aceitavel.
- `manual`: operacao exige intervencao humana.
- `disabled`: capacidade desligada por decisao.
- `unconfigured`: dependencia ausente e sem fallback suficiente.

## Matriz de integracao por servico

## `analytics`

- Plano: analytics plane
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.
- Contrato HTTP: `docs/contracts/http/analytics.openapi.yaml`

### Pontos de integracao

- leitura operacional de multiplos contextos;
- consolidacao para edge;
- nao escreve verdade transacional de outros dominios.

### Rotas com potencial de integracao

- `GET /api/analytics/reports/adapter-catalog` - Read external adapter capability catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/integration-readiness` - Read external integration readiness
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/saas-control` - Read SaaS control posture by tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/contract-governance` - Read contract governance posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/hardening-review` - Read hardening review
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/core-operations` - Read core product operations
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/relationship-intelligence` - Read relationship intelligence
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/compliance-control` - Read fiscal and privacy compliance control
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/analytics/reports/go-live-control` - Read go-live rollout control
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `billing`

- Plano: transaction/billing plane
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.
- Contrato HTTP: `docs/contracts/http/billing.openapi.yaml`

### Pontos de integracao

- capabilities de provider ou adapters externos;
- fallback explicito para desenvolvimento;
- relatorio de postura via analytics ou health details;
- eventos e callbacks quando ha provider externo.

### Rotas com potencial de integracao

- `GET /health/details` - Return readiness details and gateway posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/billing/gateways` - List gateway capabilities and Pix posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/billing/gateways/{provider}` - Read one gateway capability
  - Parametros: `provider`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/billing/plans` - List billing plans including flat, hybrid and usage-based pricing
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/billing/plans` - Create billing plan
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/billing/subscriptions` - List subscriptions
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/billing/subscriptions` - Create subscription
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/billing/subscriptions/{publicId}/usage-pricing` - Project usage-based charge for one subscription
  - Parametros: `publicId`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/billing/invoices/{publicId}/attempts` - Create payment attempt with idempotency support
  - Parametros: `Idempotency-Key`, `publicId`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `catalog`

- Plano: transaction/catalog plane
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.
- Contrato HTTP: `docs/contracts/http/catalog.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/catalog/capabilities` - Read catalog capability posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/consumers` - Read catalog consumer contracts across core domains
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/categories` - List categories by tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/catalog/categories` - Create one category
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/categories/page` - Cursor-based category listing
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/items` - List catalog items
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/catalog/items` - Create one catalog item
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/items/page` - Cursor-based item listing
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/catalog/items/bulk` - Bulk create catalog items with partial success
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/items/{publicId}` - Read one catalog item
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/catalog/items/{publicId}` - Update active state, price and attributes
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/catalog/items/{publicId}/versions` - Read catalog item version history
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `crm`

- Plano: transaction plane
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento.
- Contrato HTTP: `docs/contracts/http/crm.openapi.yaml`

### Pontos de integracao

- capabilities de provider ou adapters externos;
- fallback explicito para desenvolvimento;
- relatorio de postura via analytics ou health details;
- eventos e callbacks quando ha provider externo.

### Rotas com potencial de integracao

- `GET /api/crm/enrichment/cnpj/capabilities` - Read CNPJ enrichment provider capabilities
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/crm/enrichment/cnpj/lookup` - Lookup and enrich one CNPJ through provider contract
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/crm/pipeline/config` - Read tenant pipeline configuration
  - Parametros: `tenantSlug`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/crm/pipeline/config` - Upsert tenant pipeline configuration
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/crm/leads/intelligence/summary` - Read lead scoring and pipeline intelligence summary
  - Parametros: `tenantSlug`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `documents`

- Plano: transaction plane
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.
- Contrato HTTP: `docs/contracts/http/documents.openapi.yaml`

### Pontos de integracao

- capabilities de provider ou adapters externos;
- fallback explicito para desenvolvimento;
- relatorio de postura via analytics ou health details;
- eventos e callbacks quando ha provider externo.

### Rotas com potencial de integracao

- `GET /health/details` - Return runtime readiness and storage posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/documents/signing/capabilities` - List digital signature capabilities
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/documents/signing/capabilities/{provider}` - Read one signing capability
  - Parametros: `provider`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/documents/signing/requests` - Queue one digital signature request
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/documents/storage/capabilities` - List storage capability registry
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/documents/storage/capabilities/{provider}` - Read one storage capability
  - Parametros: `provider`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/documents/attachments` - List attachments
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/documents/attachments` - Create attachment metadata
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/documents/attachments/{publicId}/versions` - List attachment versions
  - Parametros: `publicId`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/documents/attachments/{publicId}/versions` - Append attachment version
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `edge`

- Plano: public operations plane
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.
- Contrato HTTP: `docs/contracts/http/edge.openapi.yaml`

### Pontos de integracao

- agregacao de health e overviews;
- composicao com identity, analytics e servicos transacionais;
- enforcement de tenant/session para rotas protegidas.

### Rotas com potencial de integracao

- `GET /api/edge/ops/core-operations` - Read executive core product cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/relationship-overview` - Read executive relationship cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/compliance-overview` - Read executive compliance cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/go-live-overview` - Read executive go-live cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/integrations-overview` - Read executive integrations cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/saas-overview` - Read executive SaaS cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/contracts-overview` - Read executive contracts cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/edge/ops/hardening-overview` - Read executive hardening cockpit
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `engagement`

- Plano: interaction/control plane
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.
- Contrato HTTP: `docs/contracts/http/engagement.openapi.yaml`

### Pontos de integracao

- capabilities de provider ou adapters externos;
- fallback explicito para desenvolvimento;
- relatorio de postura via analytics ou health details;
- eventos e callbacks quando ha provider externo.

### Rotas com potencial de integracao

- `GET /health/details` - Return readiness details for engagement runtime
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/engagement/providers` - List provider capabilities and fallback posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/engagement/providers/{provider}` - Read one provider capability
  - Parametros: `provider`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/engagement/providers/meta-ads/leads` - Ingest inbound lead from Meta Ads
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/engagement/providers/resend/events` - Register Resend callback event
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/engagement/providers/whatsapp-cloud/events` - Register WhatsApp callback event
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/engagement/providers/telegram-bot/events` - Register Telegram callback event
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/engagement/provider-events` - List provider events
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/engagement/provider-events/{publicId}` - Read one provider event
  - Parametros: `publicId`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `finance`

- Plano: transaction/finance plane
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.
- Contrato HTTP: `docs/contracts/http/finance.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/finance/receivable-projections` - List receivable projections
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/finance/receivable-projections/sync` - Sync projections from sales and rentals
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/finance/commission-holds` - List commission holds
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/finance/commission-holds/{publicId}/release` - Release one commission hold
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/finance/activity` - List finance operational activity
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `fiscal`

- Plano: compliance plane
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.
- Contrato HTTP: `docs/contracts/http/fiscal.openapi.yaml`

### Pontos de integracao

- capabilities de provider ou adapters externos;
- fallback explicito para desenvolvimento;
- relatorio de postura via analytics ou health details;
- eventos e callbacks quando ha provider externo.

### Rotas com potencial de integracao

- `GET /api/fiscal/capabilities` - Read fiscal capability registry
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/companies/{companyPublicId}/profile` - Read fiscal company profile
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/fiscal/companies/{companyPublicId}/profile` - Upsert fiscal company profile
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies` - List retention policies by company
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/companies/{companyPublicId}/retention-execution` - Read retention execution plan for one company
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute` - Execute retention and anonymization plan
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}` - Upsert retention policy for one data domain
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/documents` - List fiscal documents
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/documents` - Issue one fiscal document
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/documents/{publicId}` - Read one fiscal document
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/documents/{publicId}/cancel` - Cancel one fiscal document
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/documents/{publicId}/correction-letter` - Register correction letter for one fiscal document
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/documents/{publicId}/invalidate` - Register invalidation for one fiscal document
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/documents/{publicId}/events` - List fiscal document audit events
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/privacy-requests` - List privacy requests
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/privacy-requests` - Create privacy request
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/privacy-requests/{publicId}` - Read one privacy request
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/privacy-requests/{publicId}/export-package` - Build export package for one privacy request
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/privacy-requests/{publicId}/execute` - Execute one privacy request with audit trail
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/fiscal/privacy-requests/{publicId}/status` - Transition privacy request lifecycle status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/consents` - List consent ledger
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/fiscal/consents` - Create consent record
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/fiscal/consents/{publicId}` - Transition consent status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/audit-events` - List fiscal audit events
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/fiscal/compliance/summary` - Read fiscal compliance summary
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `identity`

- Plano: transaction/security plane
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.
- Contrato HTTP: `docs/contracts/http/identity.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/identity/tenants` - List tenants
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/identity/tenants` - Create tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/identity/tenants/{slug}/snapshot` - Read one tenant snapshot
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/identity/sessions/login` - Authenticate identity session
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/identity/sessions/refresh` - Refresh identity session
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/identity/invitations` - Create invitation
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `notification`

- Plano: administrative/notification plane
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.
- Contrato HTTP: `docs/contracts/http/notification.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/notification/capabilities` - Read notification capability catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/notification/preferences/{userPublicId}` - Read one user notification preference
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/notification/preferences/{userPublicId}` - Upsert one user notification preference
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/notification/center` - List notification center items with cursor filters
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/notification/center` - Create one notification center item
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/notification/center/{publicId}/status` - Transition notification status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/notification/summary` - Read notification summary
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `platform-control`

- Plano: platform control plane
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.
- Contrato HTTP: `docs/contracts/http/platform-control.openapi.yaml`

### Pontos de integracao

- capabilities de provider ou adapters externos;
- fallback explicito para desenvolvimento;
- relatorio de postura via analytics ou health details;
- eventos e callbacks quando ha provider externo.

### Rotas com potencial de integracao

- `GET /api/platform-control/capabilities/catalog` - List platform capability catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/providers/catalog` - List provider capability catalog and environment posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements` - List tenant entitlements with cursor pagination
  - Parametros: `tenantSlug`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/feature-flags` - List tenant feature flags with capability metadata
  - Parametros: `tenantSlug`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}` - Upsert one entitlement
  - Parametros: `capabilityKey`, `tenantSlug`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}` - Upsert one feature flag using entitlement governance
  - Parametros: `capabilityKey`, `tenantSlug`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk` - Bulk upsert entitlements with partial success
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults` - List provider defaults selected for one tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}` - Upsert provider default for one tenant capability
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/quotas` - List quotas by tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}` - Upsert one quota
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk` - Bulk upsert quotas with partial success
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/blocks` - List tenant blocks
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}` - Upsert tenant block
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/metering` - Read metering snapshots and summary with cursor pagination
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots` - Create one usage snapshot
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/usage-summary` - Read quota and metering utilization summary
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness` - Read tenant lifecycle readiness and provider posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs` - List onboarding and offboarding jobs with cursor pagination
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}` - Read one lifecycle job with audit events
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview` - Preview onboarding plan, provider defaults and lifecycle readiness
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding` - Queue onboarding job with Idempotency-Key and 202 Accepted
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview` - Preview offboarding plan, retention posture and lifecycle readiness
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding` - Queue offboarding job with Idempotency-Key and 202 Accepted
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start` - Transition lifecycle job to running
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete` - Transition lifecycle job to completed
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail` - Transition lifecycle job to failed
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel` - Transition lifecycle job to cancelled
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` - Read go-live rollout readiness by tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` - Read tenant go-live adoption baseline and gap
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` - List go-live bottlenecks and operational blockers
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` - Read rollout and rollback playbook for one tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` - List recommended go-live adjustments
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply` - Apply one go-live operational adjustment
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - List go-live rollouts
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - Create one go-live rollout
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}` - Read one go-live rollout with events
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` - Transition go-live rollout to running
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` - Transition go-live rollout to completed
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` - Roll back one go-live rollout
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `rentals`

- Plano: transaction plane
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.
- Contrato HTTP: `docs/contracts/http/rentals.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/rentals/contracts` - List rental contracts
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/rentals/contracts` - Create rental contract
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/rentals/contracts/{publicId}/charges` - List contract charges
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` - Update charge status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `sales`

- Plano: transaction plane
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.
- Contrato HTTP: `docs/contracts/http/sales.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/sales/opportunities` - List opportunities
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/sales/opportunities` - Create opportunity
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/sales/proposals` - List proposals
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/sales/proposals` - Create proposal
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/sales/sales` - List sales
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/sales/invoices` - List commercial invoices
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `simulation`

- Plano: simulation plane
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.
- Contrato HTTP: `docs/contracts/http/simulation.openapi.yaml`

### Pontos de integracao

- leitura operacional de multiplos contextos;
- consolidacao para edge;
- nao escreve verdade transacional de outros dominios.

### Rotas com potencial de integracao

- `GET /api/simulation/scenarios` - List scenarios
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/simulation/scenarios` - Create scenario run
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/simulation/benchmarks/load` - Execute one load benchmark run
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `supplier`

- Plano: administrative/procurement plane
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.
- Contrato HTTP: `docs/contracts/http/supplier.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/supplier/capabilities` - Read supplier capability catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/supplier/categories` - List supplier categories
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/supplier/categories/{categoryKey}` - Upsert one supplier category
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/supplier/suppliers` - List suppliers by tenant and status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/supplier/suppliers` - Create one supplier
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/supplier/suppliers/summary` - Read supplier summary
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/supplier/suppliers/{publicId}` - Read one supplier
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/supplier/suppliers/{publicId}` - Update one supplier
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `support`

- Plano: administrative plane
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.
- Contrato HTTP: `docs/contracts/http/support.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/support/capabilities` - Read support capability catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/support/queues` - List support queues by tenant
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PUT /api/support/queues/{queueKey}` - Upsert one support queue
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/support/cases` - List support cases with cursor filters
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/support/cases` - Create one support case
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/support/cases/summary` - Read support case summary
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/support/cases/{publicId}` - Read one support case
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/support/cases/{publicId}/status` - Transition support case status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/support/cases/{publicId}/comments` - Append comment to support case
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `webhook-hub`

- Plano: integration plane
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.
- Contrato HTTP: `docs/contracts/http/webhook-hub.openapi.yaml`

### Pontos de integracao

- intake de eventos externos;
- deduplicacao por provider e external id;
- transicoes auditaveis;
- dead letter;
- outbound endpoints por tenant;
- delivery log e reprocessamento.

### Rotas com potencial de integracao

- `GET /health/details` - Return readiness details for webhook runtime
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/webhook-hub/capabilities` - Read outbound webhook capability posture
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/webhook-hub/outbound-endpoints` - List tenant outbound endpoints
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/webhook-hub/outbound-endpoints` - Register one tenant outbound endpoint
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/webhook-hub/outbound-endpoints/{publicId}` - Read one tenant outbound endpoint
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - List outbound delivery log for one endpoint
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - Register one outbound delivery attempt
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter` - Move one outbound delivery to dead letter
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/webhook-hub/events` - List inbound webhook events
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/webhook-hub/events` - Register inbound webhook event
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/webhook-hub/events/summary` - Aggregate inbound webhook state
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/webhook-hub/events/{publicId}/dead-letter` - Move event to dead letter queue
  - Parametros: `publicId`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/webhook-hub/events/{publicId}/requeue` - Requeue dead-letter event
  - Parametros: `publicId`.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `workflow-control`

- Plano: control plane
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.
- Contrato HTTP: `docs/contracts/http/workflow-control.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/workflow-control/definitions` - List workflow definitions
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/workflow-control/definitions` - Create workflow definition
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/workflow-control/definitions/{key}` - Read one workflow definition
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/workflow-control/definitions/{key}` - Update one workflow definition
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `PATCH /api/workflow-control/definitions/{key}/status` - Update workflow definition status
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/workflow-control/capabilities/triggers` - List workflow trigger catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/workflow-control/capabilities/actions` - List workflow action catalog
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## `workflow-runtime`

- Plano: runtime plane
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.
- Contrato HTTP: `docs/contracts/http/workflow-runtime.openapi.yaml`

### Pontos de integracao

- consumo por contrato HTTP;
- possivel producao de historico, evento ou outbox;
- isolamento por tenant e correlation id.

### Rotas com potencial de integracao

- `GET /api/workflow-runtime/executions` - List workflow executions
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/workflow-runtime/executions` - Create workflow execution
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/workflow-runtime/executions/{publicId}` - Read one workflow execution
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/workflow-runtime/executions/{publicId}/actions` - List execution action snapshots
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `POST /api/workflow-runtime/executions/{publicId}/advance` - Advance one workflow execution
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.
- `GET /api/workflow-runtime/capabilities` - List runtime capabilities
  - Parametros: nenhum parametro declarado.
  - Cuidados: idempotencia em mutacao, tenant explicito e rastreabilidade quando houver consumidor externo.

### Checklist

- Existe contrato versionado?
- Existe fallback ou modo unconfigured claro?
- Existe correlation id?
- Existe estrategia de retry ou dead letter quando aplicavel?
- Existe teste de contrato ou smoke?

## Eventos versionados

- `docs/contracts/events/catalog.item.schema.json`
- `docs/contracts/events/crm.cnpj-enrichment.schema.json`
- `docs/contracts/events/documents.signing-request.schema.json`
- `docs/contracts/events/engagement.provider-event.schema.json`
- `docs/contracts/events/fiscal.consent.schema.json`
- `docs/contracts/events/fiscal.document-event.schema.json`
- `docs/contracts/events/platform-control.go-live-rollout.schema.json`
- `docs/contracts/events/platform-control.lifecycle-job.schema.json`
- `docs/contracts/events/platform-control.quota.schema.json`
- `docs/contracts/events/support.case.schema.json`
- `docs/contracts/events/webhook-hub.inbound-event.schema.json`
- `docs/contracts/events/webhook-hub.outbound-delivery.schema.json`
