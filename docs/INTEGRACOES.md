# integrações

Este documento descreve como O projeto integra contextos internos, providers externos, eventos e webhooks. Ele não e catálogo de endpoints; para rotas, use `docs/API.md`.

## Objetivo

Impedir acoplamento invisível. integração deve ser feita por contrato, adapter, evento ou webhook documentado, não por acesso livre a tabela ou dependência implicita entre serviços.

## princípios

- Consumidor conhece contrato, não implementação.
- Provider externo passa por adapter, capability registry ou endpoint de fronteira.
- Webhook crítico entra por `webhook-hub` ou por endpoint de provider explicitamente documentado.
- Evento compartilhado usa JSON Schema versionado.
- Leitura operacional agregada pertence a `analytics` ou `edge`.
- mutação cross-context deve carregar tenant, ator e correlation id quando aplicável.

## Postura de Provider

| Postura | Significado |
|---------|-------------|
| `configured` | provider configurado e pronto para uso real |
| `fallback` | caminho local/simulado permitido para desenvolvimento ou contingência |
| `manual` | operação depende de intervencao humana |
| `disabled` | capacidade desligada por decisão operacional |
| `unconfigured` | dependência ausente e sem fallback suficiente |

## Ativação BYOK De Provider

Provider externo real e sempre BYOK: a plataforma não distribui chave de terceiro, não mascara ausência de credencial e não chama API externa quando a chave não existe. O operador configura a chave em ambiente/secret manager e a capacidade passa a poder ser testada no `platform-control`.

superfície oficial:

```http
GET /api/platform-control/providers/activation/catalog
GET /api/platform-control/tenants/bootstrap-ops/providers/activation/runs
POST /api/platform-control/tenants/bootstrap-ops/providers/activation/stripe/test
```

Providers cobertos pelo gate:

| Provider | domínio | Credencial | ações |
|----------|---------|------------|-------|
| Stripe | billing | `BILLING_STRIPE_SECRET_KEY` | `connection_test`, `payment_intent.create` |
| Asaas | billing | `BILLING_ASAAS_API_KEY` | `connection_test` |
| Mercado Pago | billing | `BILLING_MERCADO_PAGO_ACCESS_TOKEN` | `connection_test` |
| Resend | engagement | `ENGAGEMENT_RESEND_API_KEY` | `connection_test`, `email.send` |
| OpenAI | ai-governance | `OPENAI_API_KEY` | `connection_test`, `response.create` |
| DocuSign | documents | `DOCUMENTS_DOCUSIGN_ACCESS_TOKEN` | `connection_test` |
| Clicksign | documents | `DOCUMENTS_CLICKSIGN_API_KEY` | `connection_test` |
| WhatsApp Cloud | engagement | `ENGAGEMENT_WHATSAPP_ACCESS_TOKEN` | `connection_test` |
| AWS Textract | document_intelligence | `AWS_TEXTRACT_ACCESS_KEY_ID` + secret/region | `credential_check` |
| Google Document AI | document_intelligence | `GOOGLE_DOCUMENT_AI_CREDENTIALS_JSON` | `credential_check` |
| Focus NFe | fiscal | `FISCAL_FOCUS_NFE_API_KEY` | `connection_test` |
| eNotas | fiscal | `FISCAL_ENOTAS_API_KEY` | `connection_test` |
| Serpro CNPJ | registry_enrichment | `CRM_SERPRO_CLIENT_ID` + secret | `connection_test` |
| BrasilAPI | registry_enrichment | API pública | `connection_test`, `cnpj.lookup` |
| ViaCEP | registry_enrichment | API pública | `connection_test`, `cep.lookup` |
| Alpha Vantage | market_macro_risk | `MARKET_ALPHA_VANTAGE_API_KEY` | `connection_test`, `fx.lookup` |
| Fixer | market_macro_risk | `MARKET_FIXER_API_KEY` | `connection_test`, `fx.latest` |
| Banco Central SGS/PTAX | market_macro_risk | API pública | `connection_test`, `series.latest` |
| NewsAPI | external_risk_feed | `NEWSAPI_KEY` | `connection_test`, `news.search` |
| GDELT | external_risk_feed | API pública | `connection_test`, `news.search` |
| Alpha Vantage News | external_risk_feed | `MARKET_ALPHA_VANTAGE_API_KEY` | `connection_test`, `news.sentiment` |

Sem chave, a resposta esperada e `status=unavailable` com `reason=credential_not_configured`. Com chave, o resultado e auditado por tenant, sem expor valor de secret, com evento de timeline e evidência operacional.

## inteligência Externa e Verificação

A `v1.5.0` organiza sinais externos em cinco blocos:

- OCR/document intelligence para documentos, contratos, invoices e evidências de matching;
- Fiscal Brasil para emissão, homologação, certificado e contingência;
- Consulta cadastral Brasil para CNPJ, CEP e qualidade de master data;
- Cambio, mercado e risco macro para sinais financeiros;
- News/external risk feed para risco reputacional de clientes, fornecedores e mercado.

Esses sinais não substituem os domínios donos. Eles alimentam analytics, risk scoring, search/e-discovery e evidence vault quando um alerta precisa virar prova operacional.

## Tipos de integração

### HTTP interno/público

Usado quando o consumidor precisa consultar ou comandar um recurso de outro serviço. Deve respeitar OpenAPI e compatibilidade.

### Evento

Usado para públicar fato ocorrido em um domínio. O evento deve ser pequeno, versionado e independente de tabela interna.

### Webhook inbound

Usado quando provider externo chama a plataforma. Deve ter idempotência, normalização de payload, status de processamento e DLQ quando aplicável.

### Webhook outbound

Usado quando O projeto notifica sistemas externos. Deve registrar endpoint, tentativa, resultado, erro normalizado e dead-letter.

### Aggregation/read model

Usado por `analytics` e `edge` para consolidar leitura operacional. não deve virar caminho de escrita transacional.

## Fronteiras Principais

| Fronteira | Dono | Observação |
|-----------|------|------------|
| entrada pública e cockpits | `edge` | agrega respostas de outros domínios |
| reports e governança | `analytics` | leitura derivada e operacional |
| capabilities/providers | `platform-control` | postura de providers e tenants |
| webhooks inbound/outbound | `webhook-hub` | idempotência, delivery e DLQ |
| provider de comunicação | `engagement` | Meta, WhatsApp, Telegram, email e callbacks |
| assinatura/storage | `documents` | capabilities de storage e assinatura |
| pagamento/cobrança | `billing` | gateways, tentativas e recorrência |
| fiscal/privacidade | `fiscal` | compliance, consentimento e auditoria |

## Fluxos de integração

### Provider callback

```text
provider externo
  -> endpoint de provider ou webhook-hub
  -> normalização e idempotência
  -> persistência do evento recebido
  -> leitura operacional em analytics/edge quando aplicável
```

### Webhook outbound

```text
serviço de domínio
  -> webhook-hub outbound endpoint
  -> delivery attempt
  -> success / retry / dead-letter
```

### Readiness de provider

```text
platform-control
  -> capability/provider catalog
  -> provider activation catalog/run
  -> analytics integration readiness
  -> edge integrations overview
```

### Go-live

```text
platform-control lifecycle/go-live
  -> analytics go-live-control
  -> edge go-live-overview
```

## Eventos Compartilhados

Schemas ficam em `docs/contracts/events/`.

Eventos devem ser usados para fatos de domínio, não para comandos ambiguos. Quando a intencao for comandar outro serviço, prefira HTTP ou workflow explícito.

## Cuidados Por Area

### Comercial

`crm`, `sales`, `billing` e `finance` participam de fluxos comerciais, mas cada um mantem sua verdade. Sincronização financeira deve ser auditable e não depender de leitura informal entre tabelas.

### Automação

`workflow-control` define o que pode ser executado. `workflow-runtime` executa e registra transições. domínios chamados por workflow ainda preservam suas proprias regras.

### Compliance

`fiscal` deve manter trilha de auditoria para documentos, consentimentos, privacidade e retenção. integrações fiscais externas precisam declarar provider, status e erro normalizado.

### comunicação

`engagement` cuida de canais e provider events. `notification` cuida da central interna e preferências. Os dois não devem compartilhar estado implicito.

## Checklist de integração

- Existe contrato HTTP ou schema de evento?
- Existe dono claro do dado?
- Existe idempotência quando a chamada pode repetir?
- Existe correlation id?
- Existe tratamento de retry/DLQ quando aplicável?
- Provider externo tem postura declarada?
- O fluxo aparece em smoke, contract test ou teste unitario proporcional ao risco?

## Matriz de Providers e Fronteiras

| Area | serviço dono | Tipo de integração | Observação |
|------|--------------|--------------------|------------|
| pagamento/cobrança | `billing` | gateway capability e payment attempt | exige idempotência |
| assinatura digital | `documents` | signing capability e request | provider externo deve ficar encapsulado |
| storage documental | `documents` | storage capability | metadata não deve vazar segredo de storage |
| comunicação/callbacks | `engagement` | provider event | normalizar payload de canal |
| webhook generico | `webhook-hub` | inbound/outbound delivery | DLQ e requeue são centrais |
| CNPJ/enrichment | `crm` | provider lookup | fallback precisa ser explícito |
| fiscal/privacy | `fiscal` | provider fiscal ou execução interna | auditoria e compliance primeiro |
| tenant/provider posture | `platform-control` | provider catalog | controla defaults e readiness |
| ativação externa BYOK | `platform-control` | provider activation | testá chamadas reais apenas com chave do operador |
| AI/LLM | `ai-governance` | OpenAI Responses API opcional | fallback deterministico local quando não ha chave |

## idempotência em integrações

integração externa quase sempre pode repetir: provider reenvia callback, usuário clica duas vezes, job falha no meio, rede cai depois de gravar.

padrão esperado:

- receber chave de idempotência quando o consumidor controla a chamada;
- derivar chave estavel quando provider não envia chave;
- persistir tentativa;
- retornar o mesmo recurso quando a chamada já foi processada;
- registrar erro normalizado quando falhar;
- permitir requeue apenas quando semanticamente seguro.

No baseline atual, callbacks de provider do `engagement` usam `externalEventId` como chave de replay por tenant e provider. Um reenvio do mesmo evento retorna o `provider_event` já processado, sem duplicar touchpoint, delivery ou efeito de callback.

## Retry e Dead-letter

Retry deve ser usado para falha temporaria. Dead-letter deve ser usado quando o evento não pode ser processado com segurança naquele momento.

Falhas temporarias:

- timeout;
- provider indisponível;
- rate limit;
- conexao recusada.

Falhas não temporarias:

- payload invalido;
- assinatura invalida;
- tenant desconhecido;
- recurso de destino inexistente;
- violação de regra de domínio.

`webhook-hub` e o lugar natural para delivery log e DLQ generico. domínios especificos podem ter trilha propria quando a regra de negocio exigir.

## segurança de integração

Toda integração externa deve considerar:

- autenticidade do emissor;
- replay;
- payload malformado;
- segredo em log;
- rate limit;
- tenant correto;
- autorização do destino;
- erro sem vazamento de informação sensível.

O fechamento atual cobre os pontos centrais de integração do monorepo: capability/readiness por provider, ledger de callbacks em `engagement`, DLQ/requeue em `webhook-hub`, gateways de cobrança em `billing`, storage/assinatura em `documents`, enrichment em `crm`, governança fiscal em `fiscal` e postura por tenant em `platform-control`.

## integração com Workflows

Workflows não devem furar fronteira de domínio. O runtime pode orquestrar passos, mas cada passo deve chamar o contrato público do serviço dono.

Bom:

```text
workflow-runtime -> billing POST /api/billing/subscriptions
workflow-runtime -> notification POST /api/notification/center
```

Ruim:

```text
workflow-runtime -> tabela billing.subscriptions
workflow-runtime -> tabela notification.items
```

## observabilidade Cross-context

Para investigar fluxo distribuido, os seguintes sinais precisam aparecer quando aplicável:

- correlation id;
- tenant;
- public id do recurso;
- provider;
- external id;
- tentativa;
- status final;
- erro normalizado;
- timestamp de entrada e saída.

Sem esses sinais, integração funciona em demo mas fica dificil de operar.
