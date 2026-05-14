# POLIGLOTISMO

O ERP e poliglota por decisao arquitetural. Este documento define a licenca de uso de cada stack e o checklist minimo para criar ou manter servicos.

## Licenca Arquitetural

| Stack | Servicos atuais | Quando usar | Quando evitar |
| --- | --- | --- | --- |
| Go | `edge`, `crm`, `sales`, `documents`, `rentals` | APIs leves, IO direto, gateways e dominios com handlers simples | regra que ja depende fortemente de ecossistema .NET ou Python |
| C#/.NET | `identity`, `billing`, `finance` | dominios com regras transacionais, contratos ricos e forte tipagem | scripts pequenos ou adapters descartaveis |
| Python/FastAPI | `analytics`, `catalog`, `fiscal`, `notification`, `platform-control`, `simulation`, `supplier`, `support` | reports, catalogos operacionais, simulacao, integrações simples | transacao financeira critica sem camada forte de teste |
| TypeScript | `workflow-control`, `engagement` | orquestracao HTTP, catalogo de workflows, integrações web/eventos | processamento numerico/financeiro critico |
| Rust | `webhook-hub` | ingestao/eventos com foco em robustez e controle de runtime | CRUD simples com alta mudanca de regra |
| Elixir | `workflow-runtime` | runtime concorrente, timers, retries e processos de longa duracao | rotas CRUD convencionais |

## Checklist De Novo Servico

- Auth middleware JWT/service token e OpenFGA.
- Tenant resolver e erro `tenant_slug_required` quando aplicavel.
- Error model publico com `code` e `message`.
- `/health/live`, `/health/ready` e `/health/details`.
- Propagacao de `X-Correlation-Id` e `traceparent`.
- OpenAPI e, se publicar eventos, schema em `docs/contracts/events/`.
- Teste unitario e contrato minimo.
- Config validation para secrets e providers.
- Dockerfile e compose com env comum de seguranca.

## Regra De Contencao

Nao adicionar nova linguagem enquanto as stacks atuais nao tiverem guardrails equivalentes de auth, tracing, contrato, teste, build e dependency scan.
