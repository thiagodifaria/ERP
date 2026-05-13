# API

Este documento descreve a superficie HTTP do ERP: convencoes, indice de contratos e regras de uso. A fonte executavel dos contratos e `docs/contracts/http/*.openapi.yaml`.

## Escopo

Use este arquivo para entender como chamar e evoluir a API. Para ownership dos servicos, veja `docs/SERVICOS.md`. Para compatibilidade contratual, veja `docs/CONTRATOS.md`.

## Fonte de Verdade

- OpenAPI: `docs/contracts/http/`
- Event schemas: `docs/contracts/events/`
- Contract registry: `docs/contracts/registry.json`
- Schema registry: `docs/contracts/schema-registry.json`
- Console tecnico: `client-web/client-api`

## Indice HTTP

| Servico | Arquivo | Endpoints | Base funcional |
|---------|---------|-----------|----------------|
| `analytics` | `docs/contracts/http/analytics.openapi.yaml` | 9 | reports operacionais e governanca |
| `billing` | `docs/contracts/http/billing.openapi.yaml` | 28 | planos, assinaturas e cobranca |
| `catalog` | `docs/contracts/http/catalog.openapi.yaml` | 9 | categorias, itens e consumidores |
| `crm` | `docs/contracts/http/crm.openapi.yaml` | 20 | pipeline, leads, customers, historico, anexos e enrichment |
| `documents` | `docs/contracts/http/documents.openapi.yaml` | 8 | anexos, storage e assinatura |
| `edge` | `docs/contracts/http/edge.openapi.yaml` | 8 | cockpits cross-service |
| `engagement` | `docs/contracts/http/engagement.openapi.yaml` | 9 | providers, callbacks e conversas |
| `finance` | `docs/contracts/http/finance.openapi.yaml` | 22 | recebiveis, comissoes, tesouraria, fechamentos e atividade financeira |
| `fiscal` | `docs/contracts/http/fiscal.openapi.yaml` | 21 | fiscal, retencao, privacidade e auditoria |
| `identity` | `docs/contracts/http/identity.openapi.yaml` | 35 | tenants, sessoes e convites |
| `notification` | `docs/contracts/http/notification.openapi.yaml` | 6 | preferencias e central de notificacoes |
| `platform-control` | `docs/contracts/http/platform-control.openapi.yaml` | 39 | capabilities, quotas, lifecycle e go-live |
| `rentals` | `docs/contracts/http/rentals.openapi.yaml` | 9 | contratos recorrentes e cobrancas |
| `sales` | `docs/contracts/http/sales.openapi.yaml` | 30 | oportunidades, propostas, vendas, invoices e operacoes comerciais |
| `simulation` | `docs/contracts/http/simulation.openapi.yaml` | 6 | cenarios e benchmarks |
| `supplier` | `docs/contracts/http/supplier.openapi.yaml` | 8 | categorias e fornecedores |
| `support` | `docs/contracts/http/support.openapi.yaml` | 10 | filas, casos e comentarios |
| `webhook-hub` | `docs/contracts/http/webhook-hub.openapi.yaml` | 10 | inbound/outbound webhooks e DLQ |
| `workflow-control` | `docs/contracts/http/workflow-control.openapi.yaml` | 20 | definicoes, versionamento, runs e eventos de workflow |
| `workflow-runtime` | `docs/contracts/http/workflow-runtime.openapi.yaml` | 14 | execucoes, timeline, transicoes e actions |

## Convencoes HTTP

### Seguranca

Todos os contratos OpenAPI declaram `bearerAuth` e `internalServiceToken`. A regra operacional e:

- chamadas de usuario usam token OIDC/JWT;
- chamadas internas usam service account com audience do destino;
- health live e capability publica podem declarar excecao explicita;
- rotas tenant-aware devem receber tenant explicito e validar permissao no servico dono;
- mutacao sensivel nao deve confiar apenas na rede interna.

### Health

Servicos HTTP devem expor probes baratas e seguras:

- `GET /health/live`: processo vivo;
- `GET /health/ready`: dependencias essenciais prontas;
- `GET /health/details`: diagnostico operacional sem segredo.

### Tenant

Operacoes tenant-aware devem receber tenant de forma explicita por path, query, header ou payload conforme o contrato do servico. O consumidor nao deve depender de fallback implicito.

### Identificadores

- Recursos publicos usam `publicId`, `slug`, `key` ou identificador equivalente declarado no contrato.
- IDs internos de banco nao devem aparecer como requisito de consumidor externo.

### Idempotencia

Mutacoes sensiveis devem aceitar `Idempotency-Key` quando repeticao puder gerar duplicidade:

- pagamento;
- onboarding/offboarding;
- rollout;
- webhook;
- assinatura;
- processamento fiscal;
- operacao bulk.

### Operacoes Longas

Quando a operacao nao termina dentro da requisicao, o contrato deve preferir:

- `202 Accepted`;
- identificador publico de job, execution ou rollout;
- endpoint de leitura/polling;
- status e motivo de falha auditaveis.

### Paginacao e Bulk

- Listagens de volume devem usar cursor quando aplicavel.
- Bulk deve retornar sucesso parcial com `results`, `errors` e `summary`.
- Exportacao deve ser endpoint proprio quando o custo operacional for alto.

### Erros

Erros publicos devem ser legiveis por humanos e maquinas:

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

### Observabilidade

Mutacoes e chamadas cross-context devem registrar:

- correlation id;
- tenant quando aplicavel;
- ator quando aplicavel;
- recurso publico;
- provider externo quando existir;
- status final e erro normalizado.

## Grupos de Endpoints

### Operacao executiva

- `edge`: core, relacionamento, compliance, go-live, integracoes, SaaS, contratos e hardening.
- `analytics`: reports de readiness, governanca, hardening, compliance e go-live.

### Plataforma e tenancy

- `identity`: tenants, snapshots, login, refresh e convites.
- `platform-control`: capabilities, providers, entitlements, flags, quotas, blocks, metering, lifecycle e go-live.

### Comercial e financeiro

- `crm`: enrichment e pipeline intelligence.
- `sales`: opportunities, proposals, sales e invoices.
- `billing`: gateways, plans, subscriptions, usage pricing e payment attempts.
- `finance`: receivable projections, commission holds e activity.
- `rentals`: contracts e charges.

### Compliance e documentos

- `documents`: storage, signing, attachments e versions.
- `fiscal`: company profile, retention, fiscal documents, privacy requests, consents, audit events e compliance summary.

### Integracao e automacao

- `engagement`: provider capabilities, inbound leads, callbacks e provider events.
- `webhook-hub`: inbound events, outbound endpoints, deliveries, dead-letter e requeue.
- `workflow-control`: workflow definitions, status e catalogs.
- `workflow-runtime`: executions, actions, advance e capabilities.

### Administrativo

- `catalog`: capabilities, consumers, categories, items, bulk e item versions.
- `support`: queues, cases, comments, status e summary.
- `supplier`: capabilities, categories, suppliers e summary.
- `notification`: capabilities, preferences, center e summary.
- `simulation`: scenarios e load benchmarks.

## Validacao

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

## Checklist de Mudanca de API

- O OpenAPI correspondente foi atualizado?
- O endpoint pertence ao servico dono correto?
- Tenant, idempotencia e erro publico estao claros?
- A mudanca e compativel com consumidores existentes?
- Existe teste de contrato, unidade ou smoke proporcional ao risco?
- A documentacao especifica foi atualizada sem duplicar conteudo de outros docs?

## Como Ler os Contratos

Cada OpenAPI em `docs/contracts/http/` deve ser lido como a interface publica daquele servico. Quando houver duvida entre implementacao e contrato, a diferenca deve ser tratada como bug de sincronizacao.

Leitura recomendada:

1. Abra o arquivo do servico em `docs/contracts/http/<servico>.openapi.yaml`.
2. Verifique `paths` para superficie exposta.
3. Verifique `parameters` para tenant, filtros, cursor e identificadores.
4. Verifique `requestBody` para mutacoes.
5. Verifique `responses` para status esperado.
6. Rode `./scripts/test.sh contract` antes de concluir a mudanca.

## Base URLs Locais

O runtime local usa portas vindas do `.env` ou `.env.example`. A forma exata pode variar porque `scripts/build.sh` pode remapear portas ocupadas.

Padrao conceitual:

```text
edge                 http://localhost:${EDGE_HTTP_PORT}
gateway              http://localhost:${GATEWAY_HTTP_PORT}
identity             http://localhost:${IDENTITY_HTTP_PORT}
crm                  http://localhost:${CRM_HTTP_PORT}
sales                http://localhost:${SALES_HTTP_PORT}
billing              http://localhost:${BILLING_HTTP_PORT}
finance              http://localhost:${FINANCE_HTTP_PORT}
platform-control     http://localhost:${PLATFORM_CONTROL_HTTP_PORT}
```

Para confirmar o estado real:

```bash
./scripts/build.sh ps
```

O gateway local em `infra/gateway/nginx.conf` tambem publica `/gateway/health` e roteia `/api/<servico>/` com cache para leituras, rate limit, timeouts de downstream, correlacao de request e failover passivo por dependencia.

No `client-web/client-api`, chamadas podem passar pelo proxy local do Vite. Isso evita CORS no desenvolvimento e permite testar endpoints de servicos diferentes em uma unica interface.

## Headers Recomendados

| Header | Quando usar | Observacao |
|--------|-------------|------------|
| `Content-Type: application/json` | requests com body JSON | obrigatorio na maioria das mutacoes |
| `Accept: application/json` | chamadas de API | deixa expectativa de resposta explicita |
| `X-Correlation-Id` | chamadas rastreaveis | recomendado para diagnostico cross-service |
| `Idempotency-Key` | mutacoes sensiveis | evita duplicidade em retry |
| `Authorization` | endpoints protegidos | depende do endurecimento de auth por servico |

## Exemplos de Uso

### Health

```bash
curl -i http://localhost:${EDGE_HTTP_PORT}/health/live
curl -i http://localhost:${EDGE_HTTP_PORT}/health/ready
curl -i http://localhost:${EDGE_HTTP_PORT}/health/details
```

### Mutacao idempotente

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

### Catalogos e capabilities

Endpoints de catalogo e capability devem ser baratos, cacheaveis quando possivel e seguros para leitura frequente. Eles orientam UI, integracao e readiness operacional.

Exemplos:

- provider catalog;
- trigger/action catalog;
- storage/signing capabilities;
- notification capabilities;
- fiscal capabilities.

### Lifecycle

Endpoints de lifecycle devem registrar transicoes, preservar idempotencia e expor leitura de job.

Estados comuns:

- `queued`;
- `running`;
- `completed`;
- `failed`;
- `cancelled`;
- `rolled_back` quando fizer sentido.

### Status transition

Endpoints de transicao (`/status`, `/start`, `/complete`, `/fail`, `/cancel`, `/rollback`) devem validar estado atual. Repetir uma transicao ja aplicada deve ser previsivel.

### Reports

Reports podem consolidar dados de multiplos dominios, mas devem deixar claro que sao leitura derivada. Eles nao devem ser usados como fonte para gravar verdade transacional.

### Webhooks

Webhooks precisam de:

- idempotencia;
- registro de tentativa;
- normalizacao de erro;
- status de processamento;
- dead-letter ou requeue quando aplicavel.

## Familias de API

| Familia | Servicos | Caracteristica |
|---------|----------|----------------|
| Core business | `crm`, `sales`, `catalog`, `rentals` | recursos operacionais de negocio |
| Money movement | `billing`, `finance` | cobranca, recorrencia, recebiveis e comissoes |
| Compliance | `documents`, `fiscal` | documentos, fiscal, privacidade e auditoria |
| Platform | `identity`, `platform-control` | tenant, acesso, entitlement, quota e lifecycle |
| Automation | `workflow-control`, `workflow-runtime` | definicao e execucao de workflows |
| Integration | `engagement`, `webhook-hub`, `notification` | callbacks, webhooks, comunicacao e alertas |
| Operations | `analytics`, `edge`, `simulation`, `support`, `supplier` | leitura executiva, suporte, fornecedor e capacidade |

## Quando Criar Endpoint Novo

Crie endpoint novo quando:

- o recurso tem consumidor claro;
- a responsabilidade pertence ao servico dono;
- o comportamento nao cabe em endpoint existente sem ambiguidade;
- existe contrato para request/response;
- existe plano de teste proporcional.

Evite endpoint novo quando:

- ele apenas contorna falta de modelagem interna;
- ele escreve estado de outro servico;
- ele duplica report ja existente;
- ele depende de tabela interna de outro contexto;
- ele existe apenas para uma tela temporaria sem contrato estavel.
