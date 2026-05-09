# SERVICOS

Este documento descreve os servicos do ERP: ownership, stack, caminho de implementacao e responsabilidade. Ele nao lista todos os endpoints; para isso, use `docs/API.md` e os OpenAPI em `docs/contracts/http/`.

## Inventario

| Servico | Stack | Caminho | Contrato | Endpoints |
|---------|-------|---------|----------|-----------|
| `analytics` | Python | `service-api/service-python/analytics` | `docs/contracts/http/analytics.openapi.yaml` | 9 |
| `billing` | .NET | `service-api/service-csharp/billing` | `docs/contracts/http/billing.openapi.yaml` | 9 |
| `catalog` | Python | `service-api/service-python/catalog` | `docs/contracts/http/catalog.openapi.yaml` | 12 |
| `crm` | Go | `service-api/service-golang/crm` | `docs/contracts/http/crm.openapi.yaml` | 5 |
| `documents` | Go | `service-api/service-golang/documents` | `docs/contracts/http/documents.openapi.yaml` | 10 |
| `edge` | Go | `service-api/service-golang/edge` | `docs/contracts/http/edge.openapi.yaml` | 8 |
| `engagement` | TypeScript | `service-api/service-typescript/engagement` | `docs/contracts/http/engagement.openapi.yaml` | 9 |
| `finance` | .NET | `service-api/service-csharp/finance` | `docs/contracts/http/finance.openapi.yaml` | 5 |
| `fiscal` | Python | `service-api/service-python/fiscal` | `docs/contracts/http/fiscal.openapi.yaml` | 25 |
| `identity` | .NET | `service-api/service-csharp/identity` | `docs/contracts/http/identity.openapi.yaml` | 6 |
| `notification` | Python | `service-api/service-python/notification` | `docs/contracts/http/notification.openapi.yaml` | 8 |
| `platform-control` | Python | `service-api/service-python/platform-control` | `docs/contracts/http/platform-control.openapi.yaml` | 40 |
| `rentals` | Go | `service-api/service-golang/rentals` | `docs/contracts/http/rentals.openapi.yaml` | 4 |
| `sales` | Go | `service-api/service-golang/sales` | `docs/contracts/http/sales.openapi.yaml` | 6 |
| `simulation` | Python | `service-api/service-python/simulation` | `docs/contracts/http/simulation.openapi.yaml` | 3 |
| `supplier` | Python | `service-api/service-python/supplier` | `docs/contracts/http/supplier.openapi.yaml` | 10 |
| `support` | Python | `service-api/service-python/support` | `docs/contracts/http/support.openapi.yaml` | 11 |
| `webhook-hub` | Rust | `service-api/service-rust/webhook-hub` | `docs/contracts/http/webhook-hub.openapi.yaml` | 13 |
| `workflow-control` | TypeScript | `service-api/service-typescript/workflow-control` | `docs/contracts/http/workflow-control.openapi.yaml` | 7 |
| `workflow-runtime` | Elixir | `service-api/service-elixir/workflow-runtime` | `docs/contracts/http/workflow-runtime.openapi.yaml` | 6 |

## Ownership Por Servico

### `analytics`

Responsavel por reports executivos, governanca contratual, readiness de integracoes, hardening, compliance, custos e go-live. Deve agregar leituras; nao deve virar dono transacional dos dominios que observa.

### `billing`

Responsavel por planos, assinaturas, invoices recorrentes, pricing por uso, tentativas de cobranca e recovery. Mutacoes financeiras sensiveis devem ser idempotentes.

### `catalog`

Responsavel por categorias, itens, historico de versoes, criacao em lote e contratos de consumo entre dominios. Deve preservar versionamento de item e leitura por consumidores.

### `crm`

Responsavel por leads, customers, pipeline, ownership, historico e enriquecimento de CNPJ. E a origem operacional da relacao comercial antes da venda.

### `documents`

Responsavel por anexos, storage posture, assinatura digital, versoes e metadata documental. Deve separar metadata de documento de armazenamento fisico real.

### `edge`

Responsavel por entrada publica e cockpits cross-service. Deve orquestrar leitura e exposicao, nao concentrar regra transacional de dominio.

### `engagement`

Responsavel por providers de comunicacao, inbound leads, callbacks, touchpoints, conversations e provider events. Deve normalizar eventos externos sem vazar detalhe de provider para outros dominios.

### `finance`

Responsavel por projecoes de recebiveis, atividade financeira, bloqueios/liberacao de comissao e consolidacao financeira de vendas, recorrencia e contratos.

### `fiscal`

Responsavel por perfil fiscal, politicas de retencao, documentos fiscais, eventos fiscais, consentimentos, privacidade, auditoria e resumo de compliance.

### `identity`

Responsavel por tenants, companies, usuarios, roles, times, sessoes, convites, MFA e auditoria de acesso. Outros servicos nao devem duplicar modelo de identidade.

### `notification`

Responsavel por preferencias de notificacao, central interna de alertas, severidade e lifecycle de notificacoes.

### `platform-control`

Responsavel por capabilities, providers, entitlements, feature flags, quotas, metering, tenant blocks, lifecycle e go-live. E o plano de governanca SaaS da plataforma.

### `rentals`

Responsavel por contratos recorrentes, cobrancas, reajustes, terminacoes e anexos contratuais ligados a locacao/recorrencia.

### `sales`

Responsavel por oportunidades, propostas, vendas, invoices comerciais, comissoes, renegociacoes e pendencias comerciais.

### `simulation`

Responsavel por cenarios what-if e benchmarks de carga. Deve apoiar planejamento e capacidade, nao substituir metricas reais de producao.

### `supplier`

Responsavel por categorias de fornecedores, diretorio de fornecedores e ownership de procurement.

### `support`

Responsavel por filas, casos, comentarios, SLA, operacoes em massa e resumo de atendimento.

### `webhook-hub`

Responsavel por intake de webhooks, idempotencia, estado de processamento, outbound endpoints, delivery log, dead-letter e requeue.

### `workflow-control`

Responsavel por definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos do plano de controle de workflows.

### `workflow-runtime`

Responsavel por execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.

## Regras de Servico

- Cada servico deve manter sua regra no proprio dominio.
- Cada API publica deve ter OpenAPI em `docs/contracts/http/`.
- Cada evento compartilhado deve ter JSON Schema em `docs/contracts/events/`.
- Cada contexto persistente deve ter migration em `service-api/service-postgresql/<contexto>/migrations`.
- Seeds devem ser dados de bootstrap, referencia ou smoke, nao substituto de configuracao real.
- Health deve refletir processo, readiness e diagnostico sem expor segredo.
- Providers externos devem declarar postura: `configured`, `fallback`, `manual`, `disabled` ou `unconfigured`.
- Mudancas cross-service devem atualizar contrato, teste e documentacao correspondente.

## Banco Por Contexto

O diretorio `service-api/service-postgresql/` guarda os contextos persistentes. O nome do contexto deve acompanhar o dono logico sempre que possivel:

```text
service-api/service-postgresql/<contexto>/migrations
service-api/service-postgresql/<contexto>/seeds
```

Quando um servico le dados derivados de outro contexto, essa leitura deve ser justificada pela arquitetura de report/operacao, nao por conveniencia de escrita.

## Relacao Entre Servicos

| Fluxo | Servicos envolvidos | Observacao |
|-------|---------------------|------------|
| onboarding de tenant | `identity`, `platform-control`, `analytics`, `edge` | tenancy nasce em identity e postura SaaS fica em platform-control |
| pipeline comercial | `crm`, `sales`, `billing`, `finance` | cada etapa tem owner proprio |
| recorrencia | `rentals`, `billing`, `finance` | contrato recorrente e consequencia financeira nao devem se confundir |
| documentos e fiscal | `documents`, `fiscal` | anexo/assinatura e documento fiscal tem ownership separado |
| comunicacao | `engagement`, `notification`, `webhook-hub` | provider event, alerta interno e webhook tem ciclos diferentes |
| workflow | `workflow-control`, `workflow-runtime`, dominios | runtime executa, dominios validam regras |
| go-live | `platform-control`, `analytics`, `edge` | controle operacional, report e cockpit |
| suporte/procurement | `support`, `supplier`, `catalog` | dominios administrativos com contratos proprios |

## Nivel de Maturidade Funcional

Esta tabela nao promete producao final; ela ajuda a entender o papel atual de cada modulo.

| Servico | Estado atual | Proximo cuidado natural |
|---------|--------------|-------------------------|
| `identity` | base de tenancy/acesso estabelecida | endurecer autorizacao por rota e claims |
| `crm` | relacao comercial e pipeline operacional | ampliar fluxos de import/export e deduplicacao |
| `sales` | operacao comercial basica integrada | aprofundar estados comerciais e aprovacao |
| `billing` | recorrencia e cobranca estruturadas | integrar gateways reais com fallback controlado |
| `finance` | consolidacao financeira operacional | ampliar reconciliacao e fechamento |
| `documents` | metadata, assinatura e versoes | ligar storage real e politicas de retencao |
| `fiscal` | compliance fiscal/privacidade amplo | validar provider fiscal real e regras por regime |
| `platform-control` | governanca SaaS avancada | endurecer lifecycle real e auditoria por ator |
| `analytics` | reports executivos e hardening | persistir read models quando necessario |
| `edge` | cockpits consolidados | garantir auth, cache e resiliencia por dependencia |
| `workflow-control` | definicao e catalogos | versionamento e validacao visual de fluxos |
| `workflow-runtime` | execucao duravel | robustez de retry/wait em cenarios longos |
| `engagement` | callbacks e provider events | adapters externos reais por canal |
| `webhook-hub` | inbound/outbound e DLQ | assinatura, retry policy e seguranca de endpoint |
| `catalog` | itens, versoes e consumers | relacao com pricing e disponibilidade |
| `support` | casos e filas | SLA real e automacoes |
| `supplier` | diretorio e categorias | procurement workflow |
| `notification` | preferencias e central | canais externos e templates |
| `simulation` | cenarios e benchmark | alimentar planejamento com dados reais |
| `rentals` | contratos e charges | reajuste e encerramento mais completos |

## Checklist Para Novo Servico

Um servico novo so deve entrar no monorepo se houver:

- responsabilidade funcional que nao pertence claramente a servico existente;
- contrato HTTP ou evento quando tiver consumidor externo ao modulo;
- caminho de codigo em `service-api/<stack>/<servico>`;
- contexto PostgreSQL quando houver persistencia;
- health/readiness quando exposto no runtime;
- teste unitario minimo;
- entrada no compose quando participar do stack local;
- entrada em `docs/SERVICOS.md`, `docs/API.md` e `docs/contracts/registry.json` quando aplicavel.

## Checklist Para Evoluir Servico Existente

- A mudanca respeita ownership?
- O OpenAPI foi atualizado?
- O banco mudou? Se sim, ha migration?
- Existe seed ou fixture afetada?
- O smoke precisa conhecer o novo comportamento?
- O `edge` ou `analytics` precisam refletir novo sinal?
- O `client-api` precisa regenerar catalogo?
- A documentacao alterada esta no arquivo certo?

## Sinais de Que a Fronteira Esta Errada

- O endpoint precisa consultar muitas tabelas de dominios diferentes para escrever um recurso.
- O servico passa a conhecer detalhes internos de provider que outro servico ja encapsula.
- A alteracao exige mudar varios servicos sem mudar contrato algum.
- Um report vira fonte de verdade.
- Uma tela nova dita o modelo de dominio sem contrato.

Quando isso acontecer, reavalie se a responsabilidade deveria estar em outro servico, em `analytics`, em `edge`, ou em um contrato/evento novo.
