# INTEGRACOES

Este documento descreve como o ERP integra contextos internos, providers externos, eventos e webhooks. Ele nao e catalogo de endpoints; para rotas, use `docs/API.md`.

## Objetivo

Impedir acoplamento invisivel. Integracao deve ser feita por contrato, adapter, evento ou webhook documentado, nao por acesso livre a tabela ou dependencia implicita entre servicos.

## Principios

- Consumidor conhece contrato, nao implementacao.
- Provider externo passa por adapter, capability registry ou endpoint de fronteira.
- Webhook critico entra por `webhook-hub` ou por endpoint de provider explicitamente documentado.
- Evento compartilhado usa JSON Schema versionado.
- Leitura operacional agregada pertence a `analytics` ou `edge`.
- Mutacao cross-context deve carregar tenant, ator e correlation id quando aplicavel.

## Postura de Provider

| Postura | Significado |
|---------|-------------|
| `configured` | provider configurado e pronto para uso real |
| `fallback` | caminho local/simulado permitido para desenvolvimento ou contingencia |
| `manual` | operacao depende de intervencao humana |
| `disabled` | capacidade desligada por decisao operacional |
| `unconfigured` | dependencia ausente e sem fallback suficiente |

## Tipos de Integracao

### HTTP interno/publico

Usado quando o consumidor precisa consultar ou comandar um recurso de outro servico. Deve respeitar OpenAPI e compatibilidade.

### Evento

Usado para publicar fato ocorrido em um dominio. O evento deve ser pequeno, versionado e independente de tabela interna.

### Webhook inbound

Usado quando provider externo chama a plataforma. Deve ter idempotencia, normalizacao de payload, status de processamento e DLQ quando aplicavel.

### Webhook outbound

Usado quando o ERP notifica sistemas externos. Deve registrar endpoint, tentativa, resultado, erro normalizado e dead-letter.

### Aggregation/read model

Usado por `analytics` e `edge` para consolidar leitura operacional. Nao deve virar caminho de escrita transacional.

## Fronteiras Principais

| Fronteira | Dono | Observacao |
|-----------|------|------------|
| entrada publica e cockpits | `edge` | agrega respostas de outros dominios |
| reports e governanca | `analytics` | leitura derivada e operacional |
| capabilities/providers | `platform-control` | postura de providers e tenants |
| webhooks inbound/outbound | `webhook-hub` | idempotencia, delivery e DLQ |
| provider de comunicacao | `engagement` | Meta, WhatsApp, Telegram, email e callbacks |
| assinatura/storage | `documents` | capabilities de storage e assinatura |
| pagamento/cobranca | `billing` | gateways, tentativas e recorrencia |
| fiscal/privacidade | `fiscal` | compliance, consentimento e auditoria |

## Fluxos de Integracao

### Provider callback

```text
provider externo
  -> endpoint de provider ou webhook-hub
  -> normalizacao e idempotencia
  -> persistencia do evento recebido
  -> leitura operacional em analytics/edge quando aplicavel
```

### Webhook outbound

```text
servico de dominio
  -> webhook-hub outbound endpoint
  -> delivery attempt
  -> success / retry / dead-letter
```

### Readiness de provider

```text
platform-control
  -> capability/provider catalog
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

Eventos devem ser usados para fatos de dominio, nao para comandos ambiguos. Quando a intencao for comandar outro servico, prefira HTTP ou workflow explicito.

## Cuidados Por Area

### Comercial

`crm`, `sales`, `billing` e `finance` participam de fluxos comerciais, mas cada um mantem sua verdade. Sincronizacao financeira deve ser auditable e nao depender de leitura informal entre tabelas.

### Automacao

`workflow-control` define o que pode ser executado. `workflow-runtime` executa e registra transicoes. Dominios chamados por workflow ainda preservam suas proprias regras.

### Compliance

`fiscal` deve manter trilha de auditoria para documentos, consentimentos, privacidade e retencao. Integracoes fiscais externas precisam declarar provider, status e erro normalizado.

### Comunicacao

`engagement` cuida de canais e provider events. `notification` cuida da central interna e preferencias. Os dois nao devem compartilhar estado implicito.

## Checklist de Integracao

- Existe contrato HTTP ou schema de evento?
- Existe dono claro do dado?
- Existe idempotencia quando a chamada pode repetir?
- Existe correlation id?
- Existe tratamento de retry/DLQ quando aplicavel?
- Provider externo tem postura declarada?
- O fluxo aparece em smoke, contract test ou teste unitario proporcional ao risco?

## Matriz de Providers e Fronteiras

| Area | Servico dono | Tipo de integracao | Observacao |
|------|--------------|--------------------|------------|
| pagamento/cobranca | `billing` | gateway capability e payment attempt | exige idempotencia |
| assinatura digital | `documents` | signing capability e request | provider externo deve ficar encapsulado |
| storage documental | `documents` | storage capability | metadata nao deve vazar segredo de storage |
| comunicacao/callbacks | `engagement` | provider event | normalizar payload de canal |
| webhook generico | `webhook-hub` | inbound/outbound delivery | DLQ e requeue sao centrais |
| CNPJ/enrichment | `crm` | provider lookup | fallback precisa ser explicito |
| fiscal/privacy | `fiscal` | provider fiscal ou execucao interna | auditoria e compliance primeiro |
| tenant/provider posture | `platform-control` | provider catalog | controla defaults e readiness |

## Idempotencia em Integracoes

Integracao externa quase sempre pode repetir: provider reenvia callback, usuario clica duas vezes, job falha no meio, rede cai depois de gravar.

Padrao esperado:

- receber chave de idempotencia quando o consumidor controla a chamada;
- derivar chave estavel quando provider nao envia chave;
- persistir tentativa;
- retornar o mesmo recurso quando a chamada ja foi processada;
- registrar erro normalizado quando falhar;
- permitir requeue apenas quando semanticamente seguro.

No baseline atual, callbacks de provider do `engagement` usam `externalEventId` como chave de replay por tenant e provider. Um reenvio do mesmo evento retorna o `provider_event` ja processado, sem duplicar touchpoint, delivery ou efeito de callback.

## Retry e Dead-letter

Retry deve ser usado para falha temporaria. Dead-letter deve ser usado quando o evento nao pode ser processado com seguranca naquele momento.

Falhas temporarias:

- timeout;
- provider indisponivel;
- rate limit;
- conexao recusada.

Falhas nao temporarias:

- payload invalido;
- assinatura invalida;
- tenant desconhecido;
- recurso de destino inexistente;
- violacao de regra de dominio.

`webhook-hub` e o lugar natural para delivery log e DLQ generico. Dominios especificos podem ter trilha propria quando a regra de negocio exigir.

## Seguranca de Integracao

Toda integracao externa deve considerar:

- autenticidade do emissor;
- replay;
- payload malformado;
- segredo em log;
- rate limit;
- tenant correto;
- autorizacao do destino;
- erro sem vazamento de informacao sensivel.

O fechamento atual cobre os pontos centrais de integracao do monorepo: capability/readiness por provider, ledger de callbacks em `engagement`, DLQ/requeue em `webhook-hub`, gateways de cobranca em `billing`, storage/assinatura em `documents`, enrichment em `crm`, governanca fiscal em `fiscal` e postura por tenant em `platform-control`.

## Integracao com Workflows

Workflows nao devem furar fronteira de dominio. O runtime pode orquestrar passos, mas cada passo deve chamar o contrato publico do servico dono.

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

## Observabilidade Cross-context

Para investigar fluxo distribuido, os seguintes sinais precisam aparecer quando aplicavel:

- correlation id;
- tenant;
- public id do recurso;
- provider;
- external id;
- tentativa;
- status final;
- erro normalizado;
- timestamp de entrada e saida.

Sem esses sinais, integracao funciona em demo mas fica dificil de operar.
