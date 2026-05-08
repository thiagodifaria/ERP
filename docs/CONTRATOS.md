# CONTRATOS

## Objetivo

Centralizar governanca de OpenAPI, schemas de eventos, registry, portal e regras de compatibilidade. Contrato aqui e artefato de engenharia: ele guia consumidor, teste, smoke, analytics e revisao de mudanca publica.

## Localizacao

- OpenAPI: `docs/contracts/http/`
- Eventos: `docs/contracts/events/`
- Registry HTTP/eventos/docs: `docs/contracts/registry.json`
- Schema registry: `docs/contracts/schema-registry.json`
- Portal navegavel: `docs/contracts/portal/index.html`

## Compatibilidade

### Mudanca compativel

- adicionar rota;
- adicionar campo opcional;
- adicionar enum com fallback seguro;
- ampliar descricao, exemplo ou metadata;
- adicionar filtro sem alterar comportamento default.

### Mudanca potencialmente breaking

- remover rota;
- renomear campo publico;
- alterar semantica de status;
- trocar tipo de campo;
- tornar campo opcional obrigatorio;
- alterar shape de evento consumido por outro servico;
- mudar codigo de erro publico.

## Regras obrigatorias

- Endpoint publico relevante nasce com OpenAPI.
- Evento compartilhado nasce com JSON Schema.
- Registry muda no mesmo commit do artefato novo.
- Breaking change exige changelog correto quando houver versao real para registrar.
- Historico antigo do changelog nao deve ser reescrito para refletir reorganizacao posterior.
- `docs/contracts/` e a fonte versionada atual; referencias historicas antigas continuam historicas.

## Catalogo HTTP

## `analytics`

- Titulo: ERP Analytics API
- Versao: `0.1.0`
- Arquivo: `docs/contracts/http/analytics.openapi.yaml`
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.
- Endpoints: `9`

### Rotas versionadas

#### `GET /api/analytics/reports/adapter-catalog`

- Summary: Read external adapter capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/integration-readiness`

- Summary: Read external integration readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/saas-control`

- Summary: Read SaaS control posture by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/contract-governance`

- Summary: Read contract governance posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/hardening-review`

- Summary: Read hardening review.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/core-operations`

- Summary: Read core product operations.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/relationship-intelligence`

- Summary: Read relationship intelligence.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/compliance-control`

- Summary: Read fiscal and privacy compliance control.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/analytics/reports/go-live-control`

- Summary: Read go-live rollout control.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `billing`

- Titulo: ERP Billing API
- Versao: `0.9.7`
- Arquivo: `docs/contracts/http/billing.openapi.yaml`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.
- Endpoints: `9`

### Rotas versionadas

#### `GET /health/details`

- Summary: Return readiness details and gateway posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/billing/gateways`

- Summary: List gateway capabilities and Pix posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/billing/gateways/{provider}`

- Summary: Read one gateway capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/billing/plans`

- Summary: List billing plans including flat, hybrid and usage-based pricing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/billing/plans`

- Summary: Create billing plan.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/billing/subscriptions`

- Summary: List subscriptions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/billing/subscriptions`

- Summary: Create subscription.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/billing/subscriptions/{publicId}/usage-pricing`

- Summary: Project usage-based charge for one subscription.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/billing/invoices/{publicId}/attempts`

- Summary: Create payment attempt with idempotency support.
- Parametros: `Idempotency-Key`, `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `catalog`

- Titulo: ERP Catalog API
- Versao: `0.2.0`
- Arquivo: `docs/contracts/http/catalog.openapi.yaml`
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.
- Endpoints: `12`

### Rotas versionadas

#### `GET /api/catalog/capabilities`

- Summary: Read catalog capability posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/consumers`

- Summary: Read catalog consumer contracts across core domains.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/categories`

- Summary: List categories by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/catalog/categories`

- Summary: Create one category.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/categories/page`

- Summary: Cursor-based category listing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/items`

- Summary: List catalog items.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/catalog/items`

- Summary: Create one catalog item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/items/page`

- Summary: Cursor-based item listing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/catalog/items/bulk`

- Summary: Bulk create catalog items with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/items/{publicId}`

- Summary: Read one catalog item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/catalog/items/{publicId}`

- Summary: Update active state, price and attributes.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/catalog/items/{publicId}/versions`

- Summary: Read catalog item version history.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `crm`

- Titulo: ERP CRM API
- Versao: `0.2.0`
- Arquivo: `docs/contracts/http/crm.openapi.yaml`
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento.
- Endpoints: `5`

### Rotas versionadas

#### `GET /api/crm/enrichment/cnpj/capabilities`

- Summary: Read CNPJ enrichment provider capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/crm/enrichment/cnpj/lookup`

- Summary: Lookup and enrich one CNPJ through provider contract.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/crm/pipeline/config`

- Summary: Read tenant pipeline configuration.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/crm/pipeline/config`

- Summary: Upsert tenant pipeline configuration.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/crm/leads/intelligence/summary`

- Summary: Read lead scoring and pipeline intelligence summary.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `documents`

- Titulo: ERP Documents API
- Versao: `0.9.7`
- Arquivo: `docs/contracts/http/documents.openapi.yaml`
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.
- Endpoints: `10`

### Rotas versionadas

#### `GET /health/details`

- Summary: Return runtime readiness and storage posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/documents/signing/capabilities`

- Summary: List digital signature capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/documents/signing/capabilities/{provider}`

- Summary: Read one signing capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/documents/signing/requests`

- Summary: Queue one digital signature request.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/documents/storage/capabilities`

- Summary: List storage capability registry.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/documents/storage/capabilities/{provider}`

- Summary: Read one storage capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/documents/attachments`

- Summary: List attachments.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/documents/attachments`

- Summary: Create attachment metadata.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/documents/attachments/{publicId}/versions`

- Summary: List attachment versions.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/documents/attachments/{publicId}/versions`

- Summary: Append attachment version.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `edge`

- Titulo: ERP Edge API
- Versao: `0.1.0`
- Arquivo: `docs/contracts/http/edge.openapi.yaml`
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.
- Endpoints: `8`

### Rotas versionadas

#### `GET /api/edge/ops/core-operations`

- Summary: Read executive core product cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/relationship-overview`

- Summary: Read executive relationship cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/compliance-overview`

- Summary: Read executive compliance cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/go-live-overview`

- Summary: Read executive go-live cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/integrations-overview`

- Summary: Read executive integrations cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/saas-overview`

- Summary: Read executive SaaS cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/contracts-overview`

- Summary: Read executive contracts cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/edge/ops/hardening-overview`

- Summary: Read executive hardening cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `engagement`

- Titulo: ERP Engagement API
- Versao: `0.9.7`
- Arquivo: `docs/contracts/http/engagement.openapi.yaml`
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.
- Endpoints: `9`

### Rotas versionadas

#### `GET /health/details`

- Summary: Return readiness details for engagement runtime.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/engagement/providers`

- Summary: List provider capabilities and fallback posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/engagement/providers/{provider}`

- Summary: Read one provider capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/engagement/providers/meta-ads/leads`

- Summary: Ingest inbound lead from Meta Ads.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/engagement/providers/resend/events`

- Summary: Register Resend callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/engagement/providers/whatsapp-cloud/events`

- Summary: Register WhatsApp callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/engagement/providers/telegram-bot/events`

- Summary: Register Telegram callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/engagement/provider-events`

- Summary: List provider events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/engagement/provider-events/{publicId}`

- Summary: Read one provider event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `finance`

- Titulo: ERP Finance API
- Versao: `0.4.0`
- Arquivo: `docs/contracts/http/finance.openapi.yaml`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.
- Endpoints: `5`

### Rotas versionadas

#### `GET /api/finance/receivable-projections`

- Summary: List receivable projections.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/finance/receivable-projections/sync`

- Summary: Sync projections from sales and rentals.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/finance/commission-holds`

- Summary: List commission holds.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/finance/commission-holds/{publicId}/release`

- Summary: Release one commission hold.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/finance/activity`

- Summary: List finance operational activity.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `fiscal`

- Titulo: ERP Fiscal API
- Versao: `0.1.0`
- Arquivo: `docs/contracts/http/fiscal.openapi.yaml`
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.
- Endpoints: `25`

### Rotas versionadas

#### `GET /api/fiscal/capabilities`

- Summary: Read fiscal capability registry.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Read fiscal company profile.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Upsert fiscal company profile.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/companies/{companyPublicId}/retention-policies`

- Summary: List retention policies by company.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/companies/{companyPublicId}/retention-execution`

- Summary: Read retention execution plan for one company.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute`

- Summary: Execute retention and anonymization plan.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`

- Summary: Upsert retention policy for one data domain.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/documents`

- Summary: List fiscal documents.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/documents`

- Summary: Issue one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/documents/{publicId}`

- Summary: Read one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/documents/{publicId}/cancel`

- Summary: Cancel one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/documents/{publicId}/correction-letter`

- Summary: Register correction letter for one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/documents/{publicId}/invalidate`

- Summary: Register invalidation for one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/documents/{publicId}/events`

- Summary: List fiscal document audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/privacy-requests`

- Summary: List privacy requests.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/privacy-requests`

- Summary: Create privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/privacy-requests/{publicId}`

- Summary: Read one privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/privacy-requests/{publicId}/export-package`

- Summary: Build export package for one privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/privacy-requests/{publicId}/execute`

- Summary: Execute one privacy request with audit trail.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/fiscal/privacy-requests/{publicId}/status`

- Summary: Transition privacy request lifecycle status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/consents`

- Summary: List consent ledger.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/fiscal/consents`

- Summary: Create consent record.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/fiscal/consents/{publicId}`

- Summary: Transition consent status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/audit-events`

- Summary: List fiscal audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/fiscal/compliance/summary`

- Summary: Read fiscal compliance summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `identity`

- Titulo: ERP Identity API
- Versao: `0.5.0`
- Arquivo: `docs/contracts/http/identity.openapi.yaml`
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.
- Endpoints: `6`

### Rotas versionadas

#### `GET /api/identity/tenants`

- Summary: List tenants.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/identity/tenants`

- Summary: Create tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/identity/tenants/{slug}/snapshot`

- Summary: Read one tenant snapshot.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/identity/sessions/login`

- Summary: Authenticate identity session.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/identity/sessions/refresh`

- Summary: Refresh identity session.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/identity/invitations`

- Summary: Create invitation.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `notification`

- Titulo: ERP Notification API
- Versao: `0.1.0`
- Arquivo: `docs/contracts/http/notification.openapi.yaml`
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.
- Endpoints: `7`

### Rotas versionadas

#### `GET /api/notification/capabilities`

- Summary: Read notification capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/notification/preferences/{userPublicId}`

- Summary: Read one user notification preference.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/notification/preferences/{userPublicId}`

- Summary: Upsert one user notification preference.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/notification/center`

- Summary: List notification center items with cursor filters.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/notification/center`

- Summary: Create one notification center item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/notification/center/{publicId}/status`

- Summary: Transition notification status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/notification/summary`

- Summary: Read notification summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `platform-control`

- Titulo: ERP Platform Control API
- Versao: `0.2.0`
- Arquivo: `docs/contracts/http/platform-control.openapi.yaml`
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.
- Endpoints: `40`

### Rotas versionadas

#### `GET /api/platform-control/capabilities/catalog`

- Summary: List platform capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/providers/catalog`

- Summary: List provider capability catalog and environment posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/entitlements`

- Summary: List tenant entitlements with cursor pagination.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/feature-flags`

- Summary: List tenant feature flags with capability metadata.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`

- Summary: Upsert one entitlement.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}`

- Summary: Upsert one feature flag using entitlement governance.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`

- Summary: Bulk upsert entitlements with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`

- Summary: List provider defaults selected for one tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`

- Summary: Upsert provider default for one tenant capability.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/quotas`

- Summary: List quotas by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`

- Summary: Upsert one quota.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`

- Summary: Bulk upsert quotas with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/blocks`

- Summary: List tenant blocks.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`

- Summary: Upsert tenant block.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/metering`

- Summary: Read metering snapshots and summary with cursor pagination.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`

- Summary: Create one usage snapshot.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`

- Summary: Read quota and metering utilization summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`

- Summary: Read tenant lifecycle readiness and provider posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`

- Summary: List onboarding and offboarding jobs with cursor pagination.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`

- Summary: Read one lifecycle job with audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`

- Summary: Preview onboarding plan, provider defaults and lifecycle readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`

- Summary: Queue onboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `202`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`

- Summary: Preview offboarding plan, retention posture and lifecycle readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`

- Summary: Queue offboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `202`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`

- Summary: Transition lifecycle job to running.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`

- Summary: Transition lifecycle job to completed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`

- Summary: Transition lifecycle job to failed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`

- Summary: Transition lifecycle job to cancelled.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness`

- Summary: Read go-live rollout readiness by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption`

- Summary: Read tenant go-live adoption baseline and gap.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks`

- Summary: List go-live bottlenecks and operational blockers.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook`

- Summary: Read rollout and rollback playbook for one tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments`

- Summary: List recommended go-live adjustments.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply`

- Summary: Apply one go-live operational adjustment.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: List go-live rollouts.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: Create one go-live rollout.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}`

- Summary: Read one go-live rollout with events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start`

- Summary: Transition go-live rollout to running.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete`

- Summary: Transition go-live rollout to completed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback`

- Summary: Roll back one go-live rollout.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `rentals`

- Titulo: ERP Rentals API
- Versao: `0.8.0`
- Arquivo: `docs/contracts/http/rentals.openapi.yaml`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.
- Endpoints: `4`

### Rotas versionadas

#### `GET /api/rentals/contracts`

- Summary: List rental contracts.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/rentals/contracts`

- Summary: Create rental contract.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/rentals/contracts/{publicId}/charges`

- Summary: List contract charges.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`

- Summary: Update charge status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `sales`

- Titulo: ERP Sales API
- Versao: `0.7.0`
- Arquivo: `docs/contracts/http/sales.openapi.yaml`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.
- Endpoints: `6`

### Rotas versionadas

#### `GET /api/sales/opportunities`

- Summary: List opportunities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/sales/opportunities`

- Summary: Create opportunity.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/sales/proposals`

- Summary: List proposals.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/sales/proposals`

- Summary: Create proposal.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/sales/sales`

- Summary: List sales.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/sales/invoices`

- Summary: List commercial invoices.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `simulation`

- Titulo: ERP Simulation API
- Versao: `0.7.0`
- Arquivo: `docs/contracts/http/simulation.openapi.yaml`
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.
- Endpoints: `3`

### Rotas versionadas

#### `GET /api/simulation/scenarios`

- Summary: List scenarios.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/simulation/scenarios`

- Summary: Create scenario run.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/simulation/benchmarks/load`

- Summary: Execute one load benchmark run.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `supplier`

- Titulo: ERP Supplier API
- Versao: `0.1.0`
- Arquivo: `docs/contracts/http/supplier.openapi.yaml`
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.
- Endpoints: `8`

### Rotas versionadas

#### `GET /api/supplier/capabilities`

- Summary: Read supplier capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/supplier/categories`

- Summary: List supplier categories.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/supplier/categories/{categoryKey}`

- Summary: Upsert one supplier category.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/supplier/suppliers`

- Summary: List suppliers by tenant and status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/supplier/suppliers`

- Summary: Create one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/supplier/suppliers/summary`

- Summary: Read supplier summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/supplier/suppliers/{publicId}`

- Summary: Read one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/supplier/suppliers/{publicId}`

- Summary: Update one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `support`

- Titulo: ERP Support API
- Versao: `0.1.0`
- Arquivo: `docs/contracts/http/support.openapi.yaml`
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.
- Endpoints: `9`

### Rotas versionadas

#### `GET /api/support/capabilities`

- Summary: Read support capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/support/queues`

- Summary: List support queues by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PUT /api/support/queues/{queueKey}`

- Summary: Upsert one support queue.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/support/cases`

- Summary: List support cases with cursor filters.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/support/cases`

- Summary: Create one support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/support/cases/summary`

- Summary: Read support case summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/support/cases/{publicId}`

- Summary: Read one support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/support/cases/{publicId}/status`

- Summary: Transition support case status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/support/cases/{publicId}/comments`

- Summary: Append comment to support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `webhook-hub`

- Titulo: ERP Webhook Hub API
- Versao: `0.9.7`
- Arquivo: `docs/contracts/http/webhook-hub.openapi.yaml`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.
- Endpoints: `13`

### Rotas versionadas

#### `GET /health/details`

- Summary: Return readiness details for webhook runtime.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/webhook-hub/capabilities`

- Summary: Read outbound webhook capability posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/webhook-hub/outbound-endpoints`

- Summary: List tenant outbound endpoints.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/webhook-hub/outbound-endpoints`

- Summary: Register one tenant outbound endpoint.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/webhook-hub/outbound-endpoints/{publicId}`

- Summary: Read one tenant outbound endpoint.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: List outbound delivery log for one endpoint.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: Register one outbound delivery attempt.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter`

- Summary: Move one outbound delivery to dead letter.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/webhook-hub/events`

- Summary: List inbound webhook events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/webhook-hub/events`

- Summary: Register inbound webhook event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/webhook-hub/events/summary`

- Summary: Aggregate inbound webhook state.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/webhook-hub/events/{publicId}/dead-letter`

- Summary: Move event to dead letter queue.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/webhook-hub/events/{publicId}/requeue`

- Summary: Requeue dead-letter event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `workflow-control`

- Titulo: ERP Workflow Control API
- Versao: `0.6.0`
- Arquivo: `docs/contracts/http/workflow-control.openapi.yaml`
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.
- Endpoints: `7`

### Rotas versionadas

#### `GET /api/workflow-control/definitions`

- Summary: List workflow definitions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/workflow-control/definitions`

- Summary: Create workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/workflow-control/definitions/{key}`

- Summary: Read one workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/workflow-control/definitions/{key}`

- Summary: Update one workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `PATCH /api/workflow-control/definitions/{key}/status`

- Summary: Update workflow definition status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/workflow-control/capabilities/triggers`

- Summary: List workflow trigger catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/workflow-control/capabilities/actions`

- Summary: List workflow action catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## `workflow-runtime`

- Titulo: ERP Workflow Runtime API
- Versao: `0.6.0`
- Arquivo: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.
- Endpoints: `6`

### Rotas versionadas

#### `GET /api/workflow-runtime/executions`

- Summary: List workflow executions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/workflow-runtime/executions`

- Summary: Create workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/workflow-runtime/executions/{publicId}`

- Summary: Read one workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/workflow-runtime/executions/{publicId}/actions`

- Summary: List execution action snapshots.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `POST /api/workflow-runtime/executions/{publicId}/advance`

- Summary: Advance one workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

#### `GET /api/workflow-runtime/capabilities`

- Summary: List runtime capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Regra de compatibilidade: preservar shape publico e codigo de erro.
- Regra de revisao: qualquer mudanca de semantica precisa passar por contrato, teste e avaliacao de consumidor.

## Catalogo de eventos

### `catalog.item.schema.json`

- Arquivo: `docs/contracts/events/catalog.item.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `crm.cnpj-enrichment.schema.json`

- Arquivo: `docs/contracts/events/crm.cnpj-enrichment.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `documents.signing-request.schema.json`

- Arquivo: `docs/contracts/events/documents.signing-request.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `engagement.provider-event.schema.json`

- Arquivo: `docs/contracts/events/engagement.provider-event.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `fiscal.consent.schema.json`

- Arquivo: `docs/contracts/events/fiscal.consent.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `fiscal.document-event.schema.json`

- Arquivo: `docs/contracts/events/fiscal.document-event.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `platform-control.go-live-rollout.schema.json`

- Arquivo: `docs/contracts/events/platform-control.go-live-rollout.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `platform-control.lifecycle-job.schema.json`

- Arquivo: `docs/contracts/events/platform-control.lifecycle-job.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `platform-control.quota.schema.json`

- Arquivo: `docs/contracts/events/platform-control.quota.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `support.case.schema.json`

- Arquivo: `docs/contracts/events/support.case.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `webhook-hub.inbound-event.schema.json`

- Arquivo: `docs/contracts/events/webhook-hub.inbound-event.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

### `webhook-hub.outbound-delivery.schema.json`

- Arquivo: `docs/contracts/events/webhook-hub.outbound-delivery.schema.json`
- Formato: JSON Schema.
- Uso esperado: payload compartilhado entre servicos, provider externo ou governanca de integracao.
- Regra: mudanca em campo consumido e potencialmente breaking e exige revisao de consumidores.

## Validacao

```bash
./scripts/test.sh contract
```
