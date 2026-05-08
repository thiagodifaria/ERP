# API

## Objetivo

Este documento e a referencia central da superficie HTTP do ERP. Ele e deliberadamente extenso porque o projeto ja possui uma malha grande de servicos, contratos, eventos, governanca operacional e smoke cross-service.

A fonte executavel dos contratos fica em `docs/contracts/http/`. Este arquivo consolida esses contratos para leitura humana, onboarding, revisao de PR, governanca de breaking changes e operacao do monorepo.

## Numeros atuais

- Servicos com OpenAPI versionado: `20`
- Endpoints HTTP versionados: `201`
- Schema registry de eventos: `docs/contracts/schema-registry.json`
- Portal navegavel baseline: `docs/contracts/portal/index.html`
- Validacao oficial: `./scripts/test.sh contract`

## Indice de servicos

| Servico | Contrato | Endpoints | Responsabilidade |
|---------|----------|-----------|------------------|
| Analytics | `docs/contracts/http/analytics.openapi.yaml` | 9 | leituras operacionais pesadas, relatorios executivos, governanca, custos e hardening |
| Billing | `docs/contracts/http/billing.openapi.yaml` | 9 | planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery |
| Catalog | `docs/contracts/http/catalog.openapi.yaml` | 12 | categorias, itens, versoes de item e contratos de consumo cross-context |
| CRM | `docs/contracts/http/crm.openapi.yaml` | 5 | leads, customers, ownership, historico, notas, anexos e enriquecimento |
| Documents | `docs/contracts/http/documents.openapi.yaml` | 10 | anexos, upload, versoes, archive, access links e assinatura |
| Edge | `docs/contracts/http/edge.openapi.yaml` | 8 | entrada publica, agregacao operacional e overviews cross-service |
| Engagement | `docs/contracts/http/engagement.openapi.yaml` | 9 | campanhas, templates, touchpoints, conversas, providers e callbacks |
| Finance | `docs/contracts/http/finance.openapi.yaml` | 5 | recebiveis, payables, caixa, custos, comissoes e fechamento |
| Fiscal | `docs/contracts/http/fiscal.openapi.yaml` | 25 | perfil fiscal, documentos fiscais, retencao, consentimentos, LGPD e auditoria |
| Identity | `docs/contracts/http/identity.openapi.yaml` | 6 | tenancy, usuarios, empresas, times, roles, sessoes, convites, MFA e auditoria |
| Notification | `docs/contracts/http/notification.openapi.yaml` | 7 | preferencias, centro interno de alertas e lifecycle de notificacoes |
| Platform Control | `docs/contracts/http/platform-control.openapi.yaml` | 40 | capabilities, providers, entitlements, quotas, metering, lifecycle e go-live |
| Rentals | `docs/contracts/http/rentals.openapi.yaml` | 4 | contratos recorrentes, reajustes, terminacoes, cobrancas e anexos |
| Sales | `docs/contracts/http/sales.openapi.yaml` | 6 | oportunidades, propostas, vendas, invoices, comissoes e pendencias |
| Simulation | `docs/contracts/http/simulation.openapi.yaml` | 3 | cenarios what-if, benchmark de carga e modelagem operacional |
| Supplier | `docs/contracts/http/supplier.openapi.yaml` | 8 | categorias de fornecedores, diretorio e procurement ownership |
| Support | `docs/contracts/http/support.openapi.yaml` | 9 | filas, casos, SLA, comentarios e summary de atendimento |
| Webhook Hub | `docs/contracts/http/webhook-hub.openapi.yaml` | 13 | intake de webhooks, idempotencia, transicoes, DLQ e outbound delivery |
| Workflow Control | `docs/contracts/http/workflow-control.openapi.yaml` | 7 | definicoes, versoes publicadas, runs e eventos de controle |
| Workflow Runtime | `docs/contracts/http/workflow-runtime.openapi.yaml` | 6 | execucao duravel, transicoes, timeline, actions, retries e compensacoes |

## Padroes transversais

### Health

- `GET /health/live` deve indicar que o processo HTTP esta vivo.
- `GET /health/ready` deve refletir dependencias necessarias para atender trafego real.
- `GET /health/details` deve expor postura operacional sem vazar segredo.
- Probes devem ser baratos, idempotentes e seguros para chamada frequente.

### Tenant

- Rotas tenant-aware devem receber `tenantSlug` por path, query ou payload de forma explicita.
- Quando existir fallback de bootstrap, ele deve ser documentado e nunca mascarar isolamento em producao.
- Leitura agregada deve preservar a fronteira do tenant mesmo quando o dado vem de multiplos servicos.
- Mutacao cross-context deve carregar tenant, ator e correlation id.

### Idempotencia

- `POST`, `PUT`, `PATCH` e operacoes de lifecycle devem avaliar `Idempotency-Key` quando houver risco de duplicidade.
- Webhooks, onboarding, offboarding, rollouts, cobrancas, assinatura e provider callbacks exigem cuidado maior.
- Resposta idempotente deve manter o mesmo recurso publico quando a mesma chave for repetida.

### Async

- Operacoes longas devem preferir `202 Accepted`.
- O payload deve incluir identificador publico do job, execution, rollout ou request.
- Deve existir rota de polling ou leitura de detalhe com status operacional.
- Falhas devem ser persistidas com codigo e motivo auditavel.

### Paginacao e bulk

- Listagens transacionais de volume devem preferir cursor pagination.
- Bulk deve retornar `results`, `errors` e `summary` para partial success.
- Exportacoes devem ser separadas de listagens normais quando o custo for alto.

### Erros

- Erros publicos devem expor `code`, `message` e `details` opcional.
- Codigo de erro deve ser estavel e legivel por maquina.
- Mensagem deve ser curta e operacional.
- Status HTTP deve refletir o problema real: validacao, conflito, ausente, nao autorizado, dependencia indisponivel ou falha interna.

### Observabilidade

- Toda requisicao relevante deve carregar correlation id.
- Rotas tenant-aware devem registrar tenant.
- Mutacoes devem registrar ator, recurso publico, transicao e resultado.
- Providers externos devem registrar provider, external id, tentativa, status e erro normalizado.

## Mapa detalhado por servico

## Analytics

- Servico: `analytics`
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`
- Titulo OpenAPI: ERP Analytics API
- Versao do contrato: `0.1.0`
- Responsabilidade: leituras operacionais pesadas, relatorios executivos, governanca, custos e hardening.
- Descricao do contrato: Executive reports, adapter catalog, SaaS control and contract governance.
- Total de endpoints versionados: `9`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /api/analytics/reports/adapter-catalog`

- Resumo OpenAPI: Read external adapter capability catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/analytics/reports/integration-readiness`

- Resumo OpenAPI: Read external integration readiness.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/analytics/reports/saas-control`

- Resumo OpenAPI: Read SaaS control posture by tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/analytics/reports/contract-governance`

- Resumo OpenAPI: Read contract governance posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/analytics/reports/hardening-review`

- Resumo OpenAPI: Read hardening review.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/analytics/reports/core-operations`

- Resumo OpenAPI: Read core product operations.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/analytics/reports/relationship-intelligence`

- Resumo OpenAPI: Read relationship intelligence.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/analytics/reports/compliance-control`

- Resumo OpenAPI: Read fiscal and privacy compliance control.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `GET /api/analytics/reports/go-live-control`

- Resumo OpenAPI: Read go-live rollout control.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/analytics.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/analytics.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Billing

- Servico: `billing`
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`
- Titulo OpenAPI: ERP Billing API
- Versao do contrato: `0.9.7`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.
- Descricao do contrato: Subscription billing, payment recovery and gateway capabilities.
- Total de endpoints versionados: `9`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /health/details`

- Resumo OpenAPI: Return readiness details and gateway posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: health deve refletir dependencia real quando for readiness/detail e deve permanecer barato para probes frequentes.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/billing/gateways`

- Resumo OpenAPI: List gateway capabilities and Pix posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/billing/gateways/{provider}`

- Resumo OpenAPI: Read one gateway capability.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `provider`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/billing/plans`

- Resumo OpenAPI: List billing plans including flat, hybrid and usage-based pricing.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/billing/plans`

- Resumo OpenAPI: Create billing plan.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/billing/subscriptions`

- Resumo OpenAPI: List subscriptions.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `POST /api/billing/subscriptions`

- Resumo OpenAPI: Create subscription.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/billing/subscriptions/{publicId}/usage-pricing`

- Resumo OpenAPI: Project usage-based charge for one subscription.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `publicId`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `POST /api/billing/invoices/{publicId}/attempts`

- Resumo OpenAPI: Create payment attempt with idempotency support.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/billing.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `Idempotency-Key`, `publicId`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/billing.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Catalog

- Servico: `catalog`
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`
- Titulo OpenAPI: ERP Catalog API
- Versao do contrato: `0.2.0`
- Responsabilidade: categorias, itens, versoes de item e contratos de consumo cross-context.
- Descricao do contrato: Product and service catalog with categories, activation, versioned items, cursor pagination and bulk creation.
- Total de endpoints versionados: `12`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/catalog/capabilities`

- Resumo OpenAPI: Read catalog capability posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/catalog/consumers`

- Resumo OpenAPI: Read catalog consumer contracts across core domains.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/catalog/categories`

- Resumo OpenAPI: List categories by tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/catalog/categories`

- Resumo OpenAPI: Create one category.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/catalog/categories/page`

- Resumo OpenAPI: Cursor-based category listing.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/catalog/items`

- Resumo OpenAPI: List catalog items.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `POST /api/catalog/items`

- Resumo OpenAPI: Create one catalog item.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/catalog/items/page`

- Resumo OpenAPI: Cursor-based item listing.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `POST /api/catalog/items/bulk`

- Resumo OpenAPI: Bulk create catalog items with partial success.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 10. `GET /api/catalog/items/{publicId}`

- Resumo OpenAPI: Read one catalog item.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 11. `PATCH /api/catalog/items/{publicId}`

- Resumo OpenAPI: Update active state, price and attributes.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 12. `GET /api/catalog/items/{publicId}/versions`

- Resumo OpenAPI: Read catalog item version history.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/catalog.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/catalog.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## CRM

- Servico: `crm`
- Contrato fonte: `docs/contracts/http/crm.openapi.yaml`
- Titulo OpenAPI: ERP CRM API
- Versao do contrato: `0.2.0`
- Responsabilidade: leads, customers, ownership, historico, notas, anexos e enriquecimento.
- Descricao do contrato: CRM leads, customers, activity, attachments and commercial intelligence.
- Total de endpoints versionados: `5`
- Teste recomendado: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.

### Rotas

#### 1. `GET /api/crm/enrichment/cnpj/capabilities`

- Resumo OpenAPI: Read CNPJ enrichment provider capabilities.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/crm.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/crm.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/crm/enrichment/cnpj/lookup`

- Resumo OpenAPI: Lookup and enrich one CNPJ through provider contract.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/crm.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/crm.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/crm/pipeline/config`

- Resumo OpenAPI: Read tenant pipeline configuration.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/crm.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `tenantSlug`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/crm.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `PUT /api/crm/pipeline/config`

- Resumo OpenAPI: Upsert tenant pipeline configuration.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/crm.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/crm.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/crm/leads/intelligence/summary`

- Resumo OpenAPI: Read lead scoring and pipeline intelligence summary.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/crm.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `tenantSlug`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/crm.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Documents

- Servico: `documents`
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`
- Titulo OpenAPI: ERP Documents API
- Versao do contrato: `0.9.7`
- Responsabilidade: anexos, upload, versoes, archive, access links e assinatura.
- Descricao do contrato: Attachment governance, upload orchestration and storage capabilities.
- Total de endpoints versionados: `10`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /health/details`

- Resumo OpenAPI: Return runtime readiness and storage posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: health deve refletir dependencia real quando for readiness/detail e deve permanecer barato para probes frequentes.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/documents/signing/capabilities`

- Resumo OpenAPI: List digital signature capabilities.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/documents/signing/capabilities/{provider}`

- Resumo OpenAPI: Read one signing capability.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `provider`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/documents/signing/requests`

- Resumo OpenAPI: Queue one digital signature request.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/documents/storage/capabilities`

- Resumo OpenAPI: List storage capability registry.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/documents/storage/capabilities/{provider}`

- Resumo OpenAPI: Read one storage capability.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `provider`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/documents/attachments`

- Resumo OpenAPI: List attachments.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `POST /api/documents/attachments`

- Resumo OpenAPI: Create attachment metadata.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `GET /api/documents/attachments/{publicId}/versions`

- Resumo OpenAPI: List attachment versions.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `publicId`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 10. `POST /api/documents/attachments/{publicId}/versions`

- Resumo OpenAPI: Append attachment version.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/documents.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/documents.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Edge

- Servico: `edge`
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`
- Titulo OpenAPI: ERP Edge API
- Versao do contrato: `0.1.0`
- Responsabilidade: entrada publica, agregacao operacional e overviews cross-service.
- Descricao do contrato: Aggregated operational cockpits for tenants, contracts and SaaS control.
- Total de endpoints versionados: `8`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/edge/ops/core-operations`

- Resumo OpenAPI: Read executive core product cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/edge/ops/relationship-overview`

- Resumo OpenAPI: Read executive relationship cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/edge/ops/compliance-overview`

- Resumo OpenAPI: Read executive compliance cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/edge/ops/go-live-overview`

- Resumo OpenAPI: Read executive go-live cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/edge/ops/integrations-overview`

- Resumo OpenAPI: Read executive integrations cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/edge/ops/saas-overview`

- Resumo OpenAPI: Read executive SaaS cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/edge/ops/contracts-overview`

- Resumo OpenAPI: Read executive contracts cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/edge/ops/hardening-overview`

- Resumo OpenAPI: Read executive hardening cockpit.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/edge.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/edge.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Engagement

- Servico: `engagement`
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`
- Titulo OpenAPI: ERP Engagement API
- Versao do contrato: `0.9.7`
- Responsabilidade: campanhas, templates, touchpoints, conversas, providers e callbacks.
- Descricao do contrato: Omnichannel engagement, provider callbacks and campaign operations.
- Total de endpoints versionados: `9`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /health/details`

- Resumo OpenAPI: Return readiness details for engagement runtime.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: health deve refletir dependencia real quando for readiness/detail e deve permanecer barato para probes frequentes.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/engagement/providers`

- Resumo OpenAPI: List provider capabilities and fallback posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/engagement/providers/{provider}`

- Resumo OpenAPI: Read one provider capability.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `provider`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/engagement/providers/meta-ads/leads`

- Resumo OpenAPI: Ingest inbound lead from Meta Ads.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/engagement/providers/resend/events`

- Resumo OpenAPI: Register Resend callback event.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `POST /api/engagement/providers/whatsapp-cloud/events`

- Resumo OpenAPI: Register WhatsApp callback event.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `POST /api/engagement/providers/telegram-bot/events`

- Resumo OpenAPI: Register Telegram callback event.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/engagement/provider-events`

- Resumo OpenAPI: List provider events.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `GET /api/engagement/provider-events/{publicId}`

- Resumo OpenAPI: Read one provider event.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/engagement.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `publicId`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/engagement.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Finance

- Servico: `finance`
- Contrato fonte: `docs/contracts/http/finance.openapi.yaml`
- Titulo OpenAPI: ERP Finance API
- Versao do contrato: `0.4.0`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.
- Descricao do contrato: Receivables, commission holds, cash control and cross-domain financial activity.
- Total de endpoints versionados: `5`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/finance/receivable-projections`

- Resumo OpenAPI: List receivable projections.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/finance.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/finance.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/finance/receivable-projections/sync`

- Resumo OpenAPI: Sync projections from sales and rentals.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/finance.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/finance.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/finance/commission-holds`

- Resumo OpenAPI: List commission holds.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/finance.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/finance.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/finance/commission-holds/{publicId}/release`

- Resumo OpenAPI: Release one commission hold.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/finance.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/finance.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/finance/activity`

- Resumo OpenAPI: List finance operational activity.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/finance.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/finance.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Fiscal

- Servico: `fiscal`
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`
- Titulo OpenAPI: ERP Fiscal API
- Versao do contrato: `0.1.0`
- Responsabilidade: perfil fiscal, documentos fiscais, retencao, consentimentos, LGPD e auditoria.
- Descricao do contrato: Fiscal profile, document operations, privacy rights and compliance governance.
- Total de endpoints versionados: `25`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/fiscal/capabilities`

- Resumo OpenAPI: Read fiscal capability registry.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/fiscal/companies/{companyPublicId}/profile`

- Resumo OpenAPI: Read fiscal company profile.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `PUT /api/fiscal/companies/{companyPublicId}/profile`

- Resumo OpenAPI: Upsert fiscal company profile.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/fiscal/companies/{companyPublicId}/retention-policies`

- Resumo OpenAPI: List retention policies by company.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/fiscal/companies/{companyPublicId}/retention-execution`

- Resumo OpenAPI: Read retention execution plan for one company.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute`

- Resumo OpenAPI: Execute retention and anonymization plan.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`

- Resumo OpenAPI: Upsert retention policy for one data domain.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/fiscal/documents`

- Resumo OpenAPI: List fiscal documents.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `POST /api/fiscal/documents`

- Resumo OpenAPI: Issue one fiscal document.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `201`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 10. `GET /api/fiscal/documents/{publicId}`

- Resumo OpenAPI: Read one fiscal document.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 11. `POST /api/fiscal/documents/{publicId}/cancel`

- Resumo OpenAPI: Cancel one fiscal document.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 12. `POST /api/fiscal/documents/{publicId}/correction-letter`

- Resumo OpenAPI: Register correction letter for one fiscal document.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 13. `POST /api/fiscal/documents/{publicId}/invalidate`

- Resumo OpenAPI: Register invalidation for one fiscal document.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 14. `GET /api/fiscal/documents/{publicId}/events`

- Resumo OpenAPI: List fiscal document audit events.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 15. `GET /api/fiscal/privacy-requests`

- Resumo OpenAPI: List privacy requests.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 16. `POST /api/fiscal/privacy-requests`

- Resumo OpenAPI: Create privacy request.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `201`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 17. `GET /api/fiscal/privacy-requests/{publicId}`

- Resumo OpenAPI: Read one privacy request.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 18. `GET /api/fiscal/privacy-requests/{publicId}/export-package`

- Resumo OpenAPI: Build export package for one privacy request.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 19. `POST /api/fiscal/privacy-requests/{publicId}/execute`

- Resumo OpenAPI: Execute one privacy request with audit trail.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 20. `PATCH /api/fiscal/privacy-requests/{publicId}/status`

- Resumo OpenAPI: Transition privacy request lifecycle status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 21. `GET /api/fiscal/consents`

- Resumo OpenAPI: List consent ledger.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 22. `POST /api/fiscal/consents`

- Resumo OpenAPI: Create consent record.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `201`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 23. `PATCH /api/fiscal/consents/{publicId}`

- Resumo OpenAPI: Transition consent status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`, `404`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 24. `GET /api/fiscal/audit-events`

- Resumo OpenAPI: List fiscal audit events.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 25. `GET /api/fiscal/compliance/summary`

- Resumo OpenAPI: Read fiscal compliance summary.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/fiscal.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/fiscal.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Identity

- Servico: `identity`
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`
- Titulo OpenAPI: ERP Identity API
- Versao do contrato: `0.5.0`
- Responsabilidade: tenancy, usuarios, empresas, times, roles, sessoes, convites, MFA e auditoria.
- Descricao do contrato: Tenancy, access, sessions, invitations and tenant-scoped identity governance.
- Total de endpoints versionados: `6`
- Teste recomendado: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.

### Rotas

#### 1. `GET /api/identity/tenants`

- Resumo OpenAPI: List tenants.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/identity.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/identity/tenants`

- Resumo OpenAPI: Create tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/identity.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/identity/tenants/{slug}/snapshot`

- Resumo OpenAPI: Read one tenant snapshot.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant aparece como recurso de primeira classe ou contexto de agregacao operacional.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/identity.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/identity/sessions/login`

- Resumo OpenAPI: Authenticate identity session.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/identity.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/identity/sessions/refresh`

- Resumo OpenAPI: Refresh identity session.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/identity.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `POST /api/identity/invitations`

- Resumo OpenAPI: Create invitation.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/identity.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/identity.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Notification

- Servico: `notification`
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`
- Titulo OpenAPI: ERP Notification API
- Versao do contrato: `0.1.0`
- Responsabilidade: preferencias, centro interno de alertas e lifecycle de notificacoes.
- Descricao do contrato: Internal notification center, preferences and reusable dispatch shape.
- Total de endpoints versionados: `7`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/notification/capabilities`

- Resumo OpenAPI: Read notification capability catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/notification/preferences/{userPublicId}`

- Resumo OpenAPI: Read one user notification preference.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `PUT /api/notification/preferences/{userPublicId}`

- Resumo OpenAPI: Upsert one user notification preference.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/notification/center`

- Resumo OpenAPI: List notification center items with cursor filters.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/notification/center`

- Resumo OpenAPI: Create one notification center item.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `201`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `PATCH /api/notification/center/{publicId}/status`

- Resumo OpenAPI: Transition notification status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/notification/summary`

- Resumo OpenAPI: Read notification summary.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/notification.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/notification.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Platform Control

- Servico: `platform-control`
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`
- Titulo OpenAPI: ERP Platform Control API
- Versao do contrato: `0.2.0`
- Responsabilidade: capabilities, providers, entitlements, quotas, metering, lifecycle e go-live.
- Descricao do contrato: Tenant capabilities, entitlements, quotas, metering, lifecycle jobs and SaaS governance.
- Total de endpoints versionados: `40`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /api/platform-control/capabilities/catalog`

- Resumo OpenAPI: List platform capability catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/platform-control/providers/catalog`

- Resumo OpenAPI: List provider capability catalog and environment posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/platform-control/tenants/{tenantSlug}/entitlements`

- Resumo OpenAPI: List tenant entitlements with cursor pagination.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `tenantSlug`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/platform-control/tenants/{tenantSlug}/feature-flags`

- Resumo OpenAPI: List tenant feature flags with capability metadata.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `tenantSlug`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`

- Resumo OpenAPI: Upsert one entitlement.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}`

- Resumo OpenAPI: Upsert one feature flag using entitlement governance.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`

- Resumo OpenAPI: Bulk upsert entitlements with partial success.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`

- Resumo OpenAPI: List provider defaults selected for one tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`

- Resumo OpenAPI: Upsert provider default for one tenant capability.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 10. `GET /api/platform-control/tenants/{tenantSlug}/quotas`

- Resumo OpenAPI: List quotas by tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 11. `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`

- Resumo OpenAPI: Upsert one quota.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 12. `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`

- Resumo OpenAPI: Bulk upsert quotas with partial success.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 13. `GET /api/platform-control/tenants/{tenantSlug}/blocks`

- Resumo OpenAPI: List tenant blocks.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 14. `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`

- Resumo OpenAPI: Upsert tenant block.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 15. `GET /api/platform-control/tenants/{tenantSlug}/metering`

- Resumo OpenAPI: Read metering snapshots and summary with cursor pagination.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 16. `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`

- Resumo OpenAPI: Create one usage snapshot.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 17. `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`

- Resumo OpenAPI: Read quota and metering utilization summary.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 18. `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`

- Resumo OpenAPI: Read tenant lifecycle readiness and provider posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 19. `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`

- Resumo OpenAPI: List onboarding and offboarding jobs with cursor pagination.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 20. `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`

- Resumo OpenAPI: Read one lifecycle job with audit events.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 21. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`

- Resumo OpenAPI: Preview onboarding plan, provider defaults and lifecycle readiness.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 22. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`

- Resumo OpenAPI: Queue onboarding job with Idempotency-Key and 202 Accepted.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `202`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 23. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`

- Resumo OpenAPI: Preview offboarding plan, retention posture and lifecycle readiness.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 24. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`

- Resumo OpenAPI: Queue offboarding job with Idempotency-Key and 202 Accepted.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `202`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 25. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`

- Resumo OpenAPI: Transition lifecycle job to running.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 26. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`

- Resumo OpenAPI: Transition lifecycle job to completed.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 27. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`

- Resumo OpenAPI: Transition lifecycle job to failed.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 28. `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`

- Resumo OpenAPI: Transition lifecycle job to cancelled.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 29. `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness`

- Resumo OpenAPI: Read go-live rollout readiness by tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 30. `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption`

- Resumo OpenAPI: Read tenant go-live adoption baseline and gap.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 31. `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks`

- Resumo OpenAPI: List go-live bottlenecks and operational blockers.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 32. `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook`

- Resumo OpenAPI: Read rollout and rollback playbook for one tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 33. `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments`

- Resumo OpenAPI: List recommended go-live adjustments.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 34. `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply`

- Resumo OpenAPI: Apply one go-live operational adjustment.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 35. `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Resumo OpenAPI: List go-live rollouts.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 36. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Resumo OpenAPI: Create one go-live rollout.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 37. `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}`

- Resumo OpenAPI: Read one go-live rollout with events.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 38. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start`

- Resumo OpenAPI: Transition go-live rollout to running.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 39. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete`

- Resumo OpenAPI: Transition go-live rollout to completed.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 40. `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback`

- Resumo OpenAPI: Roll back one go-live rollout.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/platform-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: tenant explicito por path ou parametro; a chamada deve preservar isolamento e auditoria do tenant.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/platform-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Rentals

- Servico: `rentals`
- Contrato fonte: `docs/contracts/http/rentals.openapi.yaml`
- Titulo OpenAPI: ERP Rentals API
- Versao do contrato: `0.8.0`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos.
- Descricao do contrato: Rental contracts, recurring charges, adjustments, rescission and attachment linkage.
- Total de endpoints versionados: `4`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /api/rentals/contracts`

- Resumo OpenAPI: List rental contracts.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/rentals.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/rentals.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/rentals/contracts`

- Resumo OpenAPI: Create rental contract.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/rentals.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/rentals.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/rentals/contracts/{publicId}/charges`

- Resumo OpenAPI: List contract charges.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/rentals.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/rentals.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`

- Resumo OpenAPI: Update charge status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/rentals.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/rentals.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Sales

- Servico: `sales`
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`
- Titulo OpenAPI: ERP Sales API
- Versao do contrato: `0.7.0`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes e pendencias.
- Descricao do contrato: Opportunities, proposals, sales, invoices and commercial lifecycle control.
- Total de endpoints versionados: `6`
- Teste recomendado: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.

### Rotas

#### 1. `GET /api/sales/opportunities`

- Resumo OpenAPI: List opportunities.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/sales.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/sales/opportunities`

- Resumo OpenAPI: Create opportunity.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/sales.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/sales/proposals`

- Resumo OpenAPI: List proposals.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/sales.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/sales/proposals`

- Resumo OpenAPI: Create proposal.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/sales.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/sales/sales`

- Resumo OpenAPI: List sales.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/sales.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/sales/invoices`

- Resumo OpenAPI: List commercial invoices.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/sales.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/sales.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Simulation

- Servico: `simulation`
- Contrato fonte: `docs/contracts/http/simulation.openapi.yaml`
- Titulo OpenAPI: ERP Simulation API
- Versao do contrato: `0.7.0`
- Responsabilidade: cenarios what-if, benchmark de carga e modelagem operacional.
- Descricao do contrato: Scenario simulation, load benchmark and cost estimation runtime.
- Total de endpoints versionados: `3`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/simulation/scenarios`

- Resumo OpenAPI: List scenarios.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/simulation.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/simulation.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/simulation/scenarios`

- Resumo OpenAPI: Create scenario run.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/simulation.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/simulation.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `POST /api/simulation/benchmarks/load`

- Resumo OpenAPI: Execute one load benchmark run.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/simulation.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/simulation.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Supplier

- Servico: `supplier`
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`
- Titulo OpenAPI: ERP Supplier API
- Versao do contrato: `0.1.0`
- Responsabilidade: categorias de fornecedores, diretorio e procurement ownership.
- Descricao do contrato: Supplier directory, categories and payables-oriented vendor governance.
- Total de endpoints versionados: `8`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/supplier/capabilities`

- Resumo OpenAPI: Read supplier capability catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/supplier/categories`

- Resumo OpenAPI: List supplier categories.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `PUT /api/supplier/categories/{categoryKey}`

- Resumo OpenAPI: Upsert one supplier category.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/supplier/suppliers`

- Resumo OpenAPI: List suppliers by tenant and status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/supplier/suppliers`

- Resumo OpenAPI: Create one supplier.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `201`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/supplier/suppliers/summary`

- Resumo OpenAPI: Read supplier summary.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/supplier/suppliers/{publicId}`

- Resumo OpenAPI: Read one supplier.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `PATCH /api/supplier/suppliers/{publicId}`

- Resumo OpenAPI: Update one supplier.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/supplier.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/supplier.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Support

- Servico: `support`
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`
- Titulo OpenAPI: ERP Support API
- Versao do contrato: `0.1.0`
- Responsabilidade: filas, casos, SLA, comentarios e summary de atendimento.
- Descricao do contrato: Queue-based support cases with SLA, comments and lifecycle history.
- Total de endpoints versionados: `9`
- Teste recomendado: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.

### Rotas

#### 1. `GET /api/support/capabilities`

- Resumo OpenAPI: Read support capability catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/support/queues`

- Resumo OpenAPI: List support queues by tenant.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `PUT /api/support/queues/{queueKey}`

- Resumo OpenAPI: Upsert one support queue.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/support/cases`

- Resumo OpenAPI: List support cases with cursor filters.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/support/cases`

- Resumo OpenAPI: Create one support case.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `201`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/support/cases/summary`

- Resumo OpenAPI: Read support case summary.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/support/cases/{publicId}`

- Resumo OpenAPI: Read one support case.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `PATCH /api/support/cases/{publicId}/status`

- Resumo OpenAPI: Transition support case status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `POST /api/support/cases/{publicId}/comments`

- Resumo OpenAPI: Append comment to support case.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/support.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar correlation id, tenant quando aplicavel, ator e status HTTP final.
- Validacao recomendada: `./scripts/test.sh unit` valida a unidade do servico e `./scripts/test.sh contract` garante registro da superficie no catalogo central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/support.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Webhook Hub

- Servico: `webhook-hub`
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`
- Titulo OpenAPI: ERP Webhook Hub API
- Versao do contrato: `0.9.7`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ e outbound delivery.
- Descricao do contrato: Inbound webhook intake, DLQ and operator recovery surface.
- Total de endpoints versionados: `13`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /health/details`

- Resumo OpenAPI: Return readiness details for webhook runtime.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: health deve refletir dependencia real quando for readiness/detail e deve permanecer barato para probes frequentes.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `GET /api/webhook-hub/capabilities`

- Resumo OpenAPI: Read outbound webhook capability posture.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/webhook-hub/outbound-endpoints`

- Resumo OpenAPI: List tenant outbound endpoints.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `POST /api/webhook-hub/outbound-endpoints`

- Resumo OpenAPI: Register one tenant outbound endpoint.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `GET /api/webhook-hub/outbound-endpoints/{publicId}`

- Resumo OpenAPI: Read one tenant outbound endpoint.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Resumo OpenAPI: List outbound delivery log for one endpoint.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Resumo OpenAPI: Register one outbound delivery attempt.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 8. `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter`

- Resumo OpenAPI: Move one outbound delivery to dead letter.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 9. `GET /api/webhook-hub/events`

- Resumo OpenAPI: List inbound webhook events.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 10. `POST /api/webhook-hub/events`

- Resumo OpenAPI: Register inbound webhook event.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: declarado.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 11. `GET /api/webhook-hub/events/summary`

- Resumo OpenAPI: Aggregate inbound webhook state.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 12. `POST /api/webhook-hub/events/{publicId}/dead-letter`

- Resumo OpenAPI: Move event to dead letter queue.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `publicId`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 13. `POST /api/webhook-hub/events/{publicId}/requeue`

- Resumo OpenAPI: Requeue dead-letter event.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/webhook-hub.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: `publicId`.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: nao declarado no OpenAPI atual.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: operacao com natureza assincrona ou de lifecycle; preferir `202 Accepted`, polling e rastreabilidade por public id quando houver enfileiramento.
- Observabilidade minima: registrar provider, external id, status de transicao, correlation id e motivo de falha quando existir.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/webhook-hub.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Workflow Control

- Servico: `workflow-control`
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`
- Titulo OpenAPI: ERP Workflow Control API
- Versao do contrato: `0.6.0`
- Responsabilidade: definicoes, versoes publicadas, runs e eventos de controle.
- Descricao do contrato: Workflow definition catalog, status, publication and action snapshots.
- Total de endpoints versionados: `7`
- Teste recomendado: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.

### Rotas

#### 1. `GET /api/workflow-control/definitions`

- Resumo OpenAPI: List workflow definitions.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/workflow-control/definitions`

- Resumo OpenAPI: Create workflow definition.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/workflow-control/definitions/{key}`

- Resumo OpenAPI: Read one workflow definition.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `PATCH /api/workflow-control/definitions/{key}`

- Resumo OpenAPI: Update one workflow definition.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `PATCH /api/workflow-control/definitions/{key}/status`

- Resumo OpenAPI: Update workflow definition status.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/workflow-control/capabilities/triggers`

- Resumo OpenAPI: List workflow trigger catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 7. `GET /api/workflow-control/capabilities/actions`

- Resumo OpenAPI: List workflow action catalog.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-control.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh contract` cobre a superficie publica principal; `./scripts/test.sh smoke` valida comportamento vivo no stack.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-control.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Workflow Runtime

- Servico: `workflow-runtime`
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Titulo OpenAPI: ERP Workflow Runtime API
- Versao do contrato: `0.6.0`
- Responsabilidade: execucao duravel, transicoes, timeline, actions, retries e compensacoes.
- Descricao do contrato: Workflow execution runtime with action snapshots, retries, delays and compensations.
- Total de endpoints versionados: `6`
- Teste recomendado: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.

### Rotas

#### 1. `GET /api/workflow-runtime/executions`

- Resumo OpenAPI: List workflow executions.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-runtime.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 2. `POST /api/workflow-runtime/executions`

- Resumo OpenAPI: Create workflow execution.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-runtime.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 3. `GET /api/workflow-runtime/executions/{publicId}`

- Resumo OpenAPI: Read one workflow execution.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-runtime.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 4. `GET /api/workflow-runtime/executions/{publicId}/actions`

- Resumo OpenAPI: List execution action snapshots.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-runtime.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 5. `POST /api/workflow-runtime/executions/{publicId}/advance`

- Resumo OpenAPI: Advance one workflow execution.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: mutacao; avaliar `Idempotency-Key`, validacao de transicao, historico operacional e outbox quando houver efeito cross-context.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-runtime.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

#### 6. `GET /api/workflow-runtime/capabilities`

- Resumo OpenAPI: List runtime capabilities.
- Descricao: sem description longa declarada no OpenAPI atual.
- Contrato fonte: `docs/contracts/http/workflow-runtime.openapi.yaml`.
- Operation id: `nao declarado`.
- Parametros declarados: nao declarado no OpenAPI atual.
- Request body: nao declarado no contrato atual.
- Respostas declaradas: `200`.
- Contexto de tenant: sem tenant obrigatorio detectado no contrato; verificar fallback do servico antes de promover uso externo.
- Idempotencia e lifecycle: leitura; manter filtros estaveis, paginacao previsivel e custo operacional compativel com uso repetido.
- Observabilidade minima: registrar run/job/execution id, status anterior, status novo, ator e correlation id.
- Validacao recomendada: `./scripts/test.sh smoke` exercita o servico no runtime integrado; contratos sao validados pelo registry central.
- Governanca de contrato: qualquer alteracao nesta rota deve atualizar `docs/contracts/http/workflow-runtime.openapi.yaml` e, se mudar semantica publica, registrar compatibilidade em `docs/CONTRATOS.md` e `docs/CHANGELOG.md`.

## Checklist para rota nova

- Definir dono do dominio e motivo da rota existir.
- Atualizar OpenAPI em `docs/contracts/http/`.
- Atualizar `docs/contracts/registry.json` se o contrato for novo.
- Adicionar schema de evento em `docs/contracts/events/` quando houver payload compartilhado.
- Garantir health/readiness se a rota introduzir dependencia nova.
- Adicionar teste unitario, contrato ou smoke de acordo com risco.
- Registrar breaking change no changelog quando houver impacto em consumidor.

## Checklist para revisao de PR

- A rota respeita tenant e authorization boundary?
- A resposta tem shape estavel?
- O erro tem codigo legivel por maquina?
- Existe idempotencia quando a mutacao pode duplicar efeito?
- A operacao deveria ser assincrona?
- A listagem precisa de cursor?
- O contrato OpenAPI esta atualizado?
- O smoke ou contrato cobre o fluxo critico?
- O servico correto e dono da regra?
- O dado esta sendo escrito no contexto certo?

## Comandos relacionados

```bash
./scripts/build.sh up
./scripts/build.sh logs edge
./scripts/build.sh migrate all
./scripts/test.sh contract
./scripts/test.sh smoke
```
