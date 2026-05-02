# INTEGRACOES

## Principios

Toda integracao externa deve ter:

- adapter especifico
- contrato interno proprio
- capability registry visivel em runtime
- timeout
- retry com backoff
- circuit breaker quando fizer sentido
- logs contextualizados
- metricas
- idempotencia
- documentacao curta

## Fluxos principais

### Gateways de pagamento

- entrada obrigatoria por `webhook-hub`
- validacao de assinatura
- deduplicacao
- normalizacao de payload
- publicacao de evento interno
- consumo e decisao no `billing`

### WhatsApp e Telegram

- ownership no `engagement`
- providers isolados por adapter
- tracking de entrega e resposta
- ligacao com lead, cliente, venda ou workflow

### Meta Ads

- ownership no `engagement`
- ingestao de leads
- sincronizacao de campanhas
- leitura de metricas
- rastreabilidade de origem

### E-mail

- ownership no `engagement`
- foco em envio transacional
- templates
- tracking basico
- fallback operacional

## Capability registry

Integracoes e adapters criticos devem expor:

- `configured`
- `fallback`
- `manual`
- `unconfigured`

Esses estados aparecem no `health/details` ou em endpoints dedicados de capability para evitar que o sistema esconda ausencia de configuracao real.

## Contratos

- specs OpenAPI ficam em `contracts/http/`
- schemas compartilhados ficam em `contracts/events/`
- mudanca publica de integracao deve atualizar os artefatos versionados junto com o codigo
