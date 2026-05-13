# SERVICOS

Este documento descreve os servicos do ERP: ownership, stack, caminho de implementacao e responsabilidade. Ele nao lista todos os endpoints; para isso, use `docs/API.md` e os OpenAPI em `docs/contracts/http/`.

## Inventario

| Servico | Stack | Caminho | Contrato | Endpoints |
|---------|-------|---------|----------|-----------|
| `analytics` | Python | `service-api/service-python/analytics` | `docs/contracts/http/analytics.openapi.yaml` | 9 |
| `billing` | .NET | `service-api/service-csharp/billing` | `docs/contracts/http/billing.openapi.yaml` | 28 |
| `catalog` | Python | `service-api/service-python/catalog` | `docs/contracts/http/catalog.openapi.yaml` | 9 |
| `crm` | Go | `service-api/service-golang/crm` | `docs/contracts/http/crm.openapi.yaml` | 20 |
| `documents` | Go | `service-api/service-golang/documents` | `docs/contracts/http/documents.openapi.yaml` | 8 |
| `edge` | Go | `service-api/service-golang/edge` | `docs/contracts/http/edge.openapi.yaml` | 8 |
| `engagement` | TypeScript | `service-api/service-typescript/engagement` | `docs/contracts/http/engagement.openapi.yaml` | 9 |
| `finance` | .NET | `service-api/service-csharp/finance` | `docs/contracts/http/finance.openapi.yaml` | 22 |
| `fiscal` | Python | `service-api/service-python/fiscal` | `docs/contracts/http/fiscal.openapi.yaml` | 21 |
| `identity` | .NET | `service-api/service-csharp/identity` | `docs/contracts/http/identity.openapi.yaml` | 35 |
| `notification` | Python | `service-api/service-python/notification` | `docs/contracts/http/notification.openapi.yaml` | 6 |
| `platform-control` | Python | `service-api/service-python/platform-control` | `docs/contracts/http/platform-control.openapi.yaml` | 39 |
| `rentals` | Go | `service-api/service-golang/rentals` | `docs/contracts/http/rentals.openapi.yaml` | 9 |
| `sales` | Go | `service-api/service-golang/sales` | `docs/contracts/http/sales.openapi.yaml` | 30 |
| `simulation` | Python | `service-api/service-python/simulation` | `docs/contracts/http/simulation.openapi.yaml` | 6 |
| `supplier` | Python | `service-api/service-python/supplier` | `docs/contracts/http/supplier.openapi.yaml` | 8 |
| `support` | Python | `service-api/service-python/support` | `docs/contracts/http/support.openapi.yaml` | 10 |
| `webhook-hub` | Rust | `service-api/service-rust/webhook-hub` | `docs/contracts/http/webhook-hub.openapi.yaml` | 10 |
| `workflow-control` | TypeScript | `service-api/service-typescript/workflow-control` | `docs/contracts/http/workflow-control.openapi.yaml` | 20 |
| `workflow-runtime` | Elixir | `service-api/service-elixir/workflow-runtime` | `docs/contracts/http/workflow-runtime.openapi.yaml` | 14 |

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
| `identity` | tenancy, empresas, usuarios, times, roles, convites, MFA, sessoes, recuperacao de senha, resolucao de acesso, Keycloak/OpenFGA e auditoria operacionais | modelar novas politicas de negocio por dominio conforme surgirem regras mais especificas |
| `crm` | relacao comercial, pipeline, deduplicacao por email, bulk/import, export filtrado, historico e outbox operacional | evoluir scoring preditivo e enriquecimento externo real |
| `sales` | oportunidade, proposta, conversao em venda, parcelas, comissao, pendencias, renegociacao, cancelamento, historico e outbox operacionais | evoluir politicas comerciais por segmento e aprovacao avancada |
| `billing` | recorrencia, cobranca, pricing flat/hybrid/usage, idempotencia de tentativas, recovery e suspend/reactivate operacionais | conectar gateways reais e conciliacao externa controlada |
| `finance` | recebiveis, liquidacao idempotente, comissoes, custos, contas a pagar, tesouraria, ledger de caixa, sync financeiro e fechamento de periodo operacionais | ampliar conciliacao bancaria/provider real e demonstrativos gerenciais |
| `documents` | metadata, upload sessions, assinatura, versoes, retencao, arquivamento e links seguros operacionais | conectar storage real e varredura automatica de retencao |
| `fiscal` | perfis fiscais, documentos, eventos, consentimentos, privacidade, retencao e auditoria operacional | conectar provider fiscal real, certificado digital e regras por regime |
| `platform-control` | capabilities, entitlements, quotas, metering, provider defaults, blocks, lifecycle e go-live operacionais | evoluir enforcement distribuido de quotas e offboarding produtivo |
| `analytics` | reports executivos, tenant 360, service pulse, hardening, simulacao, benchmark, estimativa de custo e read models operacionais quase em tempo real | evoluir streaming dedicado quando houver requisito comprovado de evento em baixa latencia |
| `edge` | cockpits consolidados com auth, health, go-live, SaaS, contratos e visoes cross-service, apoiados por gateway local com cache, rate limit, timeouts e failover passivo | evoluir replicas produtivas e politicas de trafego por ambiente |
| `workflow-control` | definicoes, catalogos, versionamento, publish/restore, runs, eventos, ledger e diagnostico por workflow operacionais | melhorar autoria visual e validacao pre-publicacao |
| `workflow-runtime` | execucao duravel, timeline, transicoes, delays, retries, capacidades e compensacoes basicas operacionais | ampliar observabilidade de execucoes longas e cargas concorrentes |
| `engagement` | campanhas, templates, callbacks idempotentes, touchpoints, deliveries, provider events e conversas estruturadas operacionais | conectar adapters externos reais por canal |
| `webhook-hub` | inbound/outbound, assinatura preparada, retries, DLQ, requeue e ledger de transicoes operacionais | endurecer seguranca de endpoint e politicas por tenant |
| `catalog` | itens, versoes imutaveis, consumers, contratos de consumo e governanca de produto operacionais | evoluir disponibilidade, pricing e politicas comerciais |
| `support` | casos, filas, SLA, comentarios, exportacao/bulk e resumo operacional | evoluir automacoes de atendimento e escalonamento |
| `supplier` | diretorio, categorias, exportacao/bulk e contratos administrativos operacionais | evoluir procurement workflow e avaliacao continua |
| `notification` | preferencias, central, severidade, ciclo de vida e bulk operacional | conectar canais externos e templates avancados |
| `simulation` | catalogo de cenarios, execucao what-if, listagem de runs, benchmark de carga e insumos de sizing operacionais | alimentar planejamento com series historicas reais |
| `rentals` | contratos, charges, reajustes, encerramento, historico, outbox, anexos e contrato HTTP completo para operacao recorrente | evoluir regras avancadas de reajuste e integracao contabil/fiscal |

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
