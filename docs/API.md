# API

Este documento descreve a superfície HTTP do projeto: convenções, índice de contratos e regras de uso. A fonte executável dos contratos é `docs/contracts/http/*.openapi.yaml`.

## Escopo

Use este arquivo para entender como chamar e evoluir a API. Para ownership dos serviços, veja `docs/SERVICOS.md`. Para compatibilidade contratual, veja `docs/CONTRATOS.md`.

## Fonte de Verdade

- OpenAPI: `docs/contracts/http/`
- Event schemas: `docs/contracts/events/`
- Contract registry: `docs/contracts/registry.json`
- Schema registry: `docs/contracts/schema-registry.json`
- Console técnico: `client-web/client-api`

## índice HTTP

| serviço | Arquivo | Endpoints | Base funcional |
|---------|---------|-----------|----------------|
| `accounting` | `docs/contracts/http/accounting.openapi.yaml` | 25 | contas contábeis, centros de custo, journal entries imutáveis, regras de posting, razão, DRE, balanço e fechamento |
| `ai-governance` | `docs/contracts/http/ai-governance.openapi.yaml` | 6 | ferramentas aprovadas, políticas, runs auditáveis, redação e governança de IA |
| `analytics` | `docs/contracts/http/analytics.openapi.yaml` | 62 | relatórios executivos, governança, readiness, production-readiness, BI semântico, risk scoring, conciliação, fechamento financeiro, dados mestres, lakehouse e inteligência externa |
| `banking` | `docs/contracts/http/banking.openapi.yaml` | 33 | CNAB, boletos, extratos, conciliação bancária, Pix cobrança/devolução/webhooks e Open Finance |
| `billing` | `docs/contracts/http/billing.openapi.yaml` | 31 | planos, assinaturas, invoices, pricing por uso e tentativas de pagamento |
| `catalog` | `docs/contracts/http/catalog.openapi.yaml` | 12 | categorias, itens, histórico de versões, bulk e contratos de consumo |
| `crm` | `docs/contracts/http/crm.openapi.yaml` | 26 | leads, customers, pipeline e enriquecimento de CNPJ |
| `documents` | `docs/contracts/http/documents.openapi.yaml` | 19 | anexos, storage, assinatura e histórico de versões |
| `edge` | `docs/contracts/http/edge.openapi.yaml` | 19 | entrada pública e cockpits cross-service |
| `engagement` | `docs/contracts/http/engagement.openapi.yaml` | 9 | providers, eventos inbound, touchpoints e conversas |
| `finance` | `docs/contracts/http/finance.openapi.yaml` | 26 | recebíveis, projeções, contas a pagar, tesouraria e comissões |
| `fiscal` | `docs/contracts/http/fiscal.openapi.yaml` | 37 | perfil fiscal, documentos fiscais, emissão, certificados, contingência, SPED, retenção, privacidade e auditoria |
| `identity` | `docs/contracts/http/identity.openapi.yaml` | 46 | tenants, usuários, sessões, convites, roles e MFA |
| `inventory` | `docs/contracts/http/inventory.openapi.yaml` | 23 | saldos por local/deposito, movimentos, reservas, custo medio/FIFO e contagem cíclica |
| `notification` | `docs/contracts/http/notification.openapi.yaml` | 8 | preferências, central de notificações e estado de entrega |
| `platform-control` | `docs/contracts/http/platform-control.openapi.yaml` | 92 | capabilities, providers, ativação externa, entitlements, quotas, lifecycle, go-live, incident command, policies, approvals, runbooks, timeline, evidências, event mesh, tenant runtime e contract evolution |
| `procurement` | `docs/contracts/http/procurement.openapi.yaml` | 25 | requisições, cotações, pedidos de compra, aprovações, recebimento e 3-way matching real |
| `rentals` | `docs/contracts/http/rentals.openapi.yaml` | 12 | contratos recorrentes e ciclo de cobranças |
| `sales` | `docs/contracts/http/sales.openapi.yaml` | 37 | oportunidades, propostas, vendas, invoices e comissões |
| `search` | `docs/contracts/http/search.openapi.yaml` | 9 | busca operacional, facets, auditoria, legal hold, discovery cases e exports |
| `simulation` | `docs/contracts/http/simulation.openapi.yaml` | 6 | cenários e benchmarks de carga |
| `supplier` | `docs/contracts/http/supplier.openapi.yaml` | 10 | categorias e diretório de fornecedores |
| `support` | `docs/contracts/http/support.openapi.yaml` | 11 | filas, casos, SLA, comentários e resumo de atendimento |
| `webhook-hub` | `docs/contracts/http/webhook-hub.openapi.yaml` | 22 | webhooks inbound/outbound, delivery log e DLQ |
| `workflow-control` | `docs/contracts/http/workflow-control.openapi.yaml` | 25 | definições de workflow, catálogos e estado de controle |
| `workflow-runtime` | `docs/contracts/http/workflow-runtime.openapi.yaml` | 15 | execuções, actions, transições, retries e compensações |

## convenções HTTP

### segurança

Todos os contratos OpenAPI declaram `bearerAuth` e `internalServiceToken`. A regra operacional é:

- chamadas de usuário usam token OIDC/JWT;
- chamadas internas usam service account com audience do destino;
- health live e capability pública podem declarar excecao explicita;
- rotas tenant-aware devem receber tenant explícito e validar permissão no serviço dono;
- mutação sensível não deve confiar apenas na rede interna.

### Health

serviços HTTP devem expor probes baratas e seguras:

- `GET /health/live`: processo vivo;
- `GET /health/ready`: dependências essenciais prontas;
- `GET /health/details`: diagnóstico operacional sem segredo.

### Production Readiness

O gate HTTP da versão 1.5.0 fica no `analytics`:

```http
GET /api/analytics/reports/production-readiness?tenant_slug=bootstrap-ops
```

A camada de runtime empresarial usa também:

```http
POST /api/analytics/reconciliation/run
GET /api/analytics/financial-close/readiness?tenant_slug=bootstrap-ops
GET /api/analytics/master-data/quality-score?tenant_slug=bootstrap-ops
GET /api/analytics/lakehouse/datasets
GET /api/platform-control/event-mesh/catalog
POST /api/platform-control/tenants/bootstrap-ops/event-mesh/events
GET /api/platform-control/tenants/bootstrap-ops/runtime/profile
GET /api/platform-control/contracts/evolution
GET /api/platform-control/tenants/bootstrap-ops/enterprise-runtime/readiness
GET /api/platform-control/policies/catalog
POST /api/platform-control/tenants/bootstrap-ops/policies/evaluate
GET /api/platform-control/tenants/bootstrap-ops/timeline
POST /api/platform-control/tenants/bootstrap-ops/approvals
GET /api/platform-control/runbooks/catalog
GET /api/platform-control/tenants/bootstrap-ops/evidence
GET /api/platform-control/providers/activation/catalog
POST /api/platform-control/tenants/bootstrap-ops/providers/activation/stripe/test
GET /api/analytics/external-intelligence/readiness
GET /api/analytics/document-intelligence/readiness
GET /api/analytics/fiscal-brazil/readiness
GET /api/analytics/registry-enrichment/brazil
GET /api/analytics/market-macro-risk
GET /api/analytics/external-risk-feed
```

Esse endpoint consolida a decisão de aceite com `release.version`, `release.releaseReady`, `blockingGates`, evidências de teste, artefatos Kubernetes, postura de providers, ativação BYOK, inteligência externa, verificação cadastral/fiscal e fechamento de go-live. Ele deve ser lido junto de `hardening-review`, `adapter-catalog`, `external-intelligence/readiness` e `go-live-control`.

### Tenant

operações tenant-aware devem receber tenant de forma explicita por path, query, header ou payload conforme o contrato do serviço. O consumidor não deve depender de fallback implicito.

### Identificadores

- Recursos públicos usam `publicId`, `slug`, `key` ou identificador equivalente declarado no contrato.
- IDs internos de banco não devem aparecer como requisito de consumidor externo.

### idempotência

Mutações sensíveis devem aceitar `Idempotency-Key` quando repetição puder gerar duplicidade:

- pagamento;
- onboarding/offboarding;
- rollout;
- webhook;
- assinatura;
- processamento fiscal;
- lancamento contábil;
- conciliação bancária;
- reserva de estoque;
- requisição, pedido de compra e recebimento;
- operação bulk.

### operações Longas

Quando a operação não termina dentro da requisição, o contrato deve preferir:

- `202 Accepted`;
- identificador público de job, execution ou rollout;
- endpoint de leitura/polling;
- status e motivo de falha auditáveis.

### Paginação e Bulk

- Listagens de volume devem usar cursor quando aplicável.
- Bulk deve retornar sucesso parcial com `results`, `errors` e `summary`.
- Exportação deve ser endpoint próprio quando o custo operacional for alto.

### Erros

Erros públicos devem ser legiveis por humanos e maquinas:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid request payload.",
  "details": {
    "field": "tenantSlug"
  }
}
```

Use status HTTP coerente: `400`, `401`, `403`, `404`, `409`, `422`, `429`, `503` ou `500`.

### observabilidade

Mutações e chamadas cross-context devem registrar:

- correlation id;
- tenant quando aplicável;
- ator quando aplicável;
- recurso público;
- provider externo quando existir;
- status final e erro normalizado.

## Grupos de Endpoints

### operação executiva

- `edge`: core, relacionamento, compliance, go-live, integrações, SaaS, contratos e hardening.
- `analytics`: reports de readiness, governança, hardening, compliance, go-live, catálogo semântico de métricas e scoring de risco.
- `ai-governance`: execução de assistente governada por políticas, ferramentas aprovadas, redação e auditoria.

### Plataforma e tenancy

- `identity`: tenants, snapshots, login, refresh e convites.
- `platform-control`: capabilities, providers, entitlements, flags, quotas, blocks, metering, lifecycle, go-live, incident command, policy decision center, approvals, runbooks, timeline e evidence vault.
- `search`: busca operacional auditada, e-discovery, legal holds e exports controlados.

### Comercial, financeiro e contábilidade

- `crm`: enrichment e pipeline intelligence.
- `sales`: opportunities, proposals, sales e invoices.
- `billing`: gateways, plans, subscriptions, usage pricing e payment attempts.
- `finance`: receivable projections, commission holds, treasury e activity.
- `accounting`: chart of accounts, cost centers, immutable journal entries, posting rules, ledger, DRE/balance sheet e period close.
- `banking`: CNAB, boletos, statements, reconciliation, Pix charges/refunds/webhooks e Open Finance.
- `rentals`: contracts e charges.

### Supply chain, compliance e documentos

- `inventory`: location balances, stock movements, reservations, FIFO/average costing e cycle counts.
- `procurement`: purchase requisitions, quotations, purchase orders, approvals, receiving e 3-way matching.
- `documents`: storage, signing, attachments e versions.
- `fiscal`: company profile, tax document lifecycle, issuance queue, certificates, contingency, SPED, privacy requests, consents, audit events e compliance summary.

### integração e automação

- `engagement`: provider capabilities, inbound leads, callbacks e provider events.
- `webhook-hub`: inbound events, outbound endpoints, deliveries, dead-letter e requeue.
- `workflow-control`: workflow definitions, status e catalogs.
- `workflow-runtime`: executions, actions, advance e capabilities.

### Administrativo

- `catalog`: capabilities, consumers, categories, items, bulk e item versions.
- `support`: queues, cases, comments, status e summary.
- `supplier`: capabilities, categories, suppliers e summary.
- `notification`: capabilities, preferences, center e summary.
- `simulation`: scenários e load benchmarks.

## validação

```bash
./scripts/test.sh contract
./scripts/test.sh smoke
```

Para testar visualmente e executar chamadas:

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

## Checklist de mudança de API

- O OpenAPI correspondente foi atualizado?
- O endpoint pertence ao serviço dono correto?
- Tenant, idempotência e erro público estão claros?
- A mudança e compatível com consumidores existentes?
- Existe teste de contrato, unidade ou smoke proporcional ao risco?
- A documentação especifica foi atualizada sem duplicar conteudo de outros docs?

## Como Ler os Contratos

Cada OpenAPI em `docs/contracts/http/` deve ser lido como a interface pública daquele serviço. Quando houver duvida entre implementação e contrato, a diferenca deve ser tratada como bug de sincronização.

Leitura recomendada:

1. Abra o arquivo do serviço em `docs/contracts/http/<serviço>.openapi.yaml`.
2. Verifique `paths` para superfície exposta.
3. Verifique `parameters` para tenant, filtros, cursor e identificadores.
4. Verifique `requestBody` para Mutações.
5. Verifique `responses` para status esperado.
6. Rode `./scripts/test.sh contract` antes de concluir a mudança.

## Base URLs Locais

O runtime local usa portas vindas do `.env` ou `.env.example`. A forma exata pode variar porque `scripts/build.sh` pode remapear portas ocupadas.

padrão conceitual:

```text
edge                 http://localhost:${EDGE_HTTP_PORT}
gateway              http://localhost:${GATEWAY_HTTP_PORT}
identity             http://localhost:${IDENTITY_HTTP_PORT}
accounting           http://localhost:${ACCOUNTING_HTTP_PORT}
banking              http://localhost:${BANKING_HTTP_PORT}
crm                  http://localhost:${CRM_HTTP_PORT}
sales                http://localhost:${SALES_HTTP_PORT}
billing              http://localhost:${BILLING_HTTP_PORT}
finance              http://localhost:${FINANCE_HTTP_PORT}
fiscal               http://localhost:${FISCAL_HTTP_PORT}
inventory            http://localhost:${INVENTORY_HTTP_PORT}
procurement          http://localhost:${PROCUREMENT_HTTP_PORT}
platform-control     http://localhost:${PLATFORM_CONTROL_HTTP_PORT}
search               http://localhost:${SEARCH_HTTP_PORT}
ai-governance        http://localhost:${AI_GOVERNANCE_HTTP_PORT}
```

Para confirmar o estado real:

```bash
./scripts/build.sh ps
```

O gateway local em `infra/gateway/nginx.conf` também pública `/gateway/health` e roteia `/api/<serviço>/` com cache para leituras, rate limit, timeouts de downstream, correlação de request e failover passivo por dependência.

No `client-web/client-api`, chamadas podem passar pelo proxy local do Vite. Isso evita CORS no desenvolvimento e permite testar endpoints de serviços diferentes em uma única interface.

## Headers Recomendados

| Header | Quando usar | Observação |
|--------|-------------|------------|
| `Content-Type: application/json` | requests com body JSON | obrigatório na maioria das Mutações |
| `Accept: application/json` | chamadas de API | deixa expectativa de resposta explicita |
| `X-Correlation-Id` | chamadas rastreaveis | recomendado para diagnóstico cross-service |
| `Idempotency-Key` | Mutações sensíveis | evita duplicidade em retry |
| `Authorization` | endpoints protegidos | depende do endurecimento de auth por serviço |

## Exemplos de Uso

### Health

```bash
curl -i http://localhost:${EDGE_HTTP_PORT}/health/live
curl -i http://localhost:${EDGE_HTTP_PORT}/health/ready
curl -i http://localhost:${EDGE_HTTP_PORT}/health/details
```

### mutação idempotente

```bash
curl -i \
  -X POST "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/tenants/bootstrap-ops/lifecycle/onboarding" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: onboarding-bootstrap-ops-001" \
  -d '{"requestedBy":"local-admin"}'
```

### Consulta operacional

```bash
curl -i "http://localhost:${EDGE_HTTP_PORT}/api/edge/ops/go-live-overview"
```

## Regras Por Tipo de Endpoint

### catálogos e capabilities

Endpoints de catálogo e capability devem ser baratos, cacheaveis quando possivel e seguros para leitura frequente. Eles orientam UI, integração e readiness operacional.

Exemplos:

- provider catalog;
- trigger/action catalog;
- storage/signing capabilities;
- notification capabilities;
- fiscal capabilities;
- accounting close capabilities;
- banking provider capabilities;
- inventory availability capabilities;
- procurement approval capabilities.

### Lifecycle

Endpoints de lifecycle devem registrar transições, preservar idempotência e expor leitura de job.

Estados comuns:

- `queued`;
- `running`;
- `completed`;
- `failed`;
- `cancelled`;
- `rolled_back` quando fizer sentido.

### Status transition

Endpoints de transicao (`/status`, `/start`, `/complete`, `/fail`, `/cancel`, `/rollback`) devem validar estado atual. Repetir uma transicao já aplicada deve ser previsivel.

### Reports

Reports podem consolidar dados de multiplos domínios, mas devem deixar claro que são leitura derivada. Eles não devem ser usados como fonte para gravar verdade transacional.

### Webhooks

Webhooks precisam de:

- idempotência;
- registro de tentativa;
- normalização de erro;
- status de processamento;
- dead-letter ou requeue quando aplicável.

## Familias de API

| Familia | serviços | Caracteristica |
|---------|----------|----------------|
| Core business | `crm`, `sales`, `catalog`, `rentals` | recursos operacionais de negocio |
| Supply chain | `inventory`, `procurement`, `supplier` | estoque, compras, fornecedores e recebimento |
| Money movement | `billing`, `finance`, `banking`, `accounting` | cobrança, tesouraria, bancos, ledger e comissões |
| Compliance | `documents`, `fiscal` | documentos, fiscal, privacidade e auditoria |
| Platform | `identity`, `platform-control` | tenant, acesso, entitlement, quota e lifecycle |
| Automation | `workflow-control`, `workflow-runtime` | definicao e execução de workflows |
| Integration | `engagement`, `webhook-hub`, `notification` | callbacks, webhooks, comunicação e alertas |
| Operations | `analytics`, `edge`, `search`, `ai-governance`, `simulation`, `support`, `supplier` | leitura executiva, busca operacional, governança de IA, suporte, fornecedor e capacidade |

## Quando Criar Endpoint Novo

Crie endpoint novo quando:

- o recurso tem consumidor claro;
- a responsabilidade pertence ao serviço dono;
- o comportamento não cabe em endpoint existente sem ambiguidade;
- existe contrato para request/response;
- existe plano de teste proporcional.

Evite endpoint novo quando:

- ele apenas contorna falta de modelagem interna;
- ele escreve estado de outro serviço;
- ele duplica report já existente;
- ele depende de tabela interna de outro contexto;
- ele existe apenas para uma tela temporaria sem contrato estavel.
