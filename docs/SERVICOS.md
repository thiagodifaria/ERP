# serviços

Este documento descreve os serviços do projeto: ownership, stack, caminho de implementação e responsabilidade. Ele não lista todos os endpoints; para isso, use `docs/API.md` e os OpenAPI em `docs/contracts/http/`.

## inventário

| serviço | Stack | Caminho | Contrato | Endpoints |
|---------|-------|---------|----------|-----------|
| `accounting` | Python | `service-api/service-python/accounting` | `docs/contracts/http/accounting.openapi.yaml` | 25 |
| `ai-governance` | Python | `service-api/service-python/ai-governance` | `docs/contracts/http/ai-governance.openapi.yaml` | 6 |
| `analytics` | Python | `service-api/service-python/analytics` | `docs/contracts/http/analytics.openapi.yaml` | 56 |
| `banking` | Python | `service-api/service-python/banking` | `docs/contracts/http/banking.openapi.yaml` | 33 |
| `billing` | .NET | `service-api/service-csharp/billing` | `docs/contracts/http/billing.openapi.yaml` | 31 |
| `catalog` | Python | `service-api/service-python/catalog` | `docs/contracts/http/catalog.openapi.yaml` | 12 |
| `crm` | Go | `service-api/service-golang/crm` | `docs/contracts/http/crm.openapi.yaml` | 26 |
| `documents` | Go | `service-api/service-golang/documents` | `docs/contracts/http/documents.openapi.yaml` | 19 |
| `edge` | Go | `service-api/service-golang/edge` | `docs/contracts/http/edge.openapi.yaml` | 19 |
| `engagement` | TypeScript | `service-api/service-typescript/engagement` | `docs/contracts/http/engagement.openapi.yaml` | 9 |
| `finance` | .NET | `service-api/service-csharp/finance` | `docs/contracts/http/finance.openapi.yaml` | 26 |
| `fiscal` | Python | `service-api/service-python/fiscal` | `docs/contracts/http/fiscal.openapi.yaml` | 37 |
| `identity` | .NET | `service-api/service-csharp/identity` | `docs/contracts/http/identity.openapi.yaml` | 46 |
| `inventory` | Python | `service-api/service-python/inventory` | `docs/contracts/http/inventory.openapi.yaml` | 23 |
| `notification` | Python | `service-api/service-python/notification` | `docs/contracts/http/notification.openapi.yaml` | 8 |
| `platform-control` | Python | `service-api/service-python/platform-control` | `docs/contracts/http/platform-control.openapi.yaml` | 89 |
| `procurement` | Python | `service-api/service-python/procurement` | `docs/contracts/http/procurement.openapi.yaml` | 25 |
| `rentals` | Go | `service-api/service-golang/rentals` | `docs/contracts/http/rentals.openapi.yaml` | 12 |
| `sales` | Go | `service-api/service-golang/sales` | `docs/contracts/http/sales.openapi.yaml` | 37 |
| `search` | Python | `service-api/service-python/search` | `docs/contracts/http/search.openapi.yaml` | 9 |
| `simulation` | Python | `service-api/service-python/simulation` | `docs/contracts/http/simulation.openapi.yaml` | 6 |
| `supplier` | Python | `service-api/service-python/supplier` | `docs/contracts/http/supplier.openapi.yaml` | 10 |
| `support` | Python | `service-api/service-python/support` | `docs/contracts/http/support.openapi.yaml` | 11 |
| `webhook-hub` | Rust | `service-api/service-rust/webhook-hub` | `docs/contracts/http/webhook-hub.openapi.yaml` | 22 |
| `workflow-control` | TypeScript | `service-api/service-typescript/workflow-control` | `docs/contracts/http/workflow-control.openapi.yaml` | 25 |
| `workflow-runtime` | Elixir | `service-api/service-elixir/workflow-runtime` | `docs/contracts/http/workflow-runtime.openapi.yaml` | 15 |

## Ownership Por serviço

### `accounting`

Responsável por plano de contas, contas gerenciais, lançamentos de diário imutáveis, regras de posting, fechamentos de período, demonstrativos e reconciliação contábil com movimentos financeiros e fiscais.

### `ai-governance`

Responsável por governança de uso de IA/LLM: catálogo de ferramentas aprovadas, políticas de uso, runs auditáveis, redação de prompt/resposta e bloqueio de ferramentas mutantes ou não registradas.

### `analytics`

Responsável por reports executivos, governança contratual, readiness de integrações, hardening, compliance, custos, production-readiness, go-live, BI semântico, risk/compliance scoring, reconciliação, fechamento financeiro, qualidade de dados mestres e manifesto lakehouse. Deve agregar leituras; não deve virar dono transacional dos domínios que observa.

### `banking`

Responsável por arquivos CNAB, boletos, lifecycle de cobrança bancária, extratos, conciliação, Pix cobrança/devolução/webhook e trilha opcional de Open Finance.

### `billing`

Responsável por planos, assinaturas, invoices recorrentes, pricing por uso, tentativas de cobrança e recovery. Mutações financeiras sensíveis devem ser idempotentes.

### `catalog`

Responsável por categorias, itens, histórico de versões, criação em lote e contratos de consumo entre domínios. Deve preservar versionamento de item e leitura por consumidores.

### `crm`

Responsável por leads, customers, pipeline, ownership, histórico e enriquecimento de CNPJ. É a origem operacional da relação comercial antes da venda.

### `documents`

Responsável por anexos, storage posture, assinatura digital, versões e metadata documental. Deve separar metadata de documento de armazenamento físico real.

### `edge`

Responsável por entrada pública e cockpits cross-service. Deve orquestrar leitura e exposição, não concentrar regra transacional de domínio.

### `engagement`

Responsável por providers de comunicação, inbound leads, callbacks, touchpoints, conversations e provider events. Deve normalizar eventos externos sem vazar detalhe de provider para outros domínios.

### `finance`

Responsável por projeções de recebíveis, atividade financeira, bloqueios/liberação de comissao, contas a pagar, tesouraria e consolidação financeira de vendas, recorrência e contratos.

### `fiscal`

Responsável por perfil fiscal, políticas de retenção, documentos fiscais, fila de emissão, certificados, contingência, eventos fiscais, SPED, consentimentos, privacidade, auditoria e resumo de compliance.

### `identity`

Responsável por tenants, companies, usuários, roles, times, sessões, convites, MFA e auditoria de acesso. Outros serviços não devem duplicar modelo de identidade.

### `inventory`

Responsável por locais, depósitos, saldos por item, movimentos, reservas, custo, contagem cíclica e integração operacional com compras, vendas, documentos e fiscal.

### `notification`

Responsável por preferências de notificação, central interna de alertas, severidade e lifecycle de notificações.

### `platform-control`

Responsável por capabilities, providers, entitlements, feature flags, quotas, metering, tenant blocks, lifecycle, go-live, incident command, policy decision center, command approvals, runbook automation, operational timeline e audit evidence vault. É o plano de governança SaaS e operacional da plataforma.

### `procurement`

Responsável por requisições de compra, RFQ/cotações, pedidos de compra, aprovações, recebimento, 3-way matching e relacionamento operacional com fornecedores.

### `rentals`

Responsável por contratos recorrentes, cobranças, reajustes, terminações e anexos contratuais ligados a locação/recorrência.

### `sales`

Responsável por oportunidades, propostas, vendas, invoices comerciais, comissões, renegociações e pendências comerciais.

### `search`

Responsável por busca operacional unificada, facets, auditoria de consulta, saved queries, e-discovery, legal hold e exports controlados com redação de dados sensíveis.

### `simulation`

Responsável por cenários what-if e benchmarks de carga. Deve apoiar planejamento e capacidade, não substituir métricas reais de produção.

### `supplier`

Responsável por categorias de fornecedores, diretório de fornecedores e dados administrativos de supplier.

### `support`

Responsável por filas, casos, comentários, SLA, operações em massa e resumo de atendimento.

### `webhook-hub`

Responsável por intake de webhooks, idempotência, estado de processamento, outbound endpoints, delivery log, dead-letter e requeue.

### `workflow-control`

Responsável por definições, versões públicadas, catálogos de trigger/action, runs e eventos do plano de controle de workflows.

### `workflow-runtime`

Responsável por execuções duraveis, timeline, actions, transições, retries, waits e compensações.

## Regras de serviço

- Cada serviço deve manter sua regra no próprio domínio.
- Cada API pública deve ter OpenAPI em `docs/contracts/http/`.
- Cada evento compartilhado deve ter JSON Schema em `docs/contracts/events/`.
- Cada contexto persistente deve ter migration em `service-api/service-postgresql/<contexto>/migrations`.
- Seeds devem ser dados de bootstrap, referência ou smoke, não substituto de configuração real.
- Health deve refletir processo, readiness e diagnóstico sem expor segredo.
- Providers externos devem declarar postura: `configured`, `fallback`, `manual`, `disabled` ou `unconfigured`.
- mudanças cross-service devem atualizar contrato, teste e documentação correspondente.

## Banco Por Contexto

O diretório `service-api/service-postgresql/` guarda os contextos persistentes. O nome do contexto deve acompanhar o dono lógico sempre que possivel:

```text
service-api/service-postgresql/<contexto>/migrations
service-api/service-postgresql/<contexto>/seeds
```

Quando um serviço le dados derivados de outro contexto, essa leitura deve ser justificada pela arquitetura de report/operação, não por conveniencia de escrita.

## relação Entre serviços

| Fluxo | serviços envolvidos | Observação |
|-------|---------------------|------------|
| onboarding de tenant | `identity`, `platform-control`, `analytics`, `edge` | tenancy nasce em identity e postura SaaS fica em platform-control |
| pipeline comercial | `crm`, `sales`, `billing`, `finance` | cada etapa tem owner próprio |
| recorrência | `rentals`, `billing`, `finance` | contrato recorrente e consequência financeira não devem se confundir |
| compras e estoque | `procurement`, `supplier`, `inventory`, `documents`, `fiscal` | compra, fornecedor, recebimento, anexo e documento fiscal tem ciclos separados |
| contábilidade e bancos | `accounting`, `finance`, `billing`, `banking`, `fiscal` | ledger, tesouraria, cobrança bancária e fiscal se reconciliam sem misturar ownership |
| documentos e fiscal | `documents`, `fiscal` | anexo/assinatura e documento fiscal tem ownership separado |
| comunicação | `engagement`, `notification`, `webhook-hub` | provider event, alerta interno e webhook tem ciclos diferentes |
| workflow | `workflow-control`, `workflow-runtime`, domínios | runtime executa, domínios validam regras |
| go-live, incidentes e enterprise runtime | `platform-control`, `analytics`, `edge` | controle operacional, incident command, policy decisions, approvals, runbooks, event mesh, tenant runtime, contract evolution, financial close, risk scoring, evidence vault, report e cockpit |
| descoberta operacional | `search`, `ai-governance`, `analytics` | busca auditada, e-discovery, governança de IA e catálogo semântico |
| administrativo | `support`, `supplier`, `catalog` | domínios administrativos com contratos próprios |

## Nivel de Maturidade Funcional

está tabela não promete produção final; ela ajuda a entender o papel atual de cada módulo.

| serviço | Estado atual | próximo cuidado natural |
|---------|--------------|-------------------------|
| `identity` | tenancy, empresas, usuários, times, roles, convites, MFA, sessões, recuperação de senha, resolução de acesso, Keycloak/OpenFGA e auditoria operacionais | modelar novas políticas de negocio por domínio conforme surgirem regras mais especificas |
| `crm` | relação comercial, pipeline, deduplicação por email, bulk/import, export filtrado, histórico e outbox operacional | evoluir scoring preditivo e enriquecimento externo real |
| `sales` | oportunidade, proposta, conversão em venda, parcelas, comissao, pendências, renegociação, cancelamento, histórico e outbox operacionais | evoluir políticas comerciais por segmento e aprovação avancada |
| `billing` | recorrência, cobrança, pricing flat/hybrid/usage, idempotência de tentativas, recovery e suspend/reactivate operacionais | conectar gateways reais e conciliação externa controlada |
| `finance` | recebíveis, liquidação idempotente, comissões, custos, contas a pagar, tesouraria, ledger de caixa, sync financeiro e fechamento de período operacionais | ampliar conciliação bancária/provider real e demonstrativos gerenciais |
| `accounting` | plano de contas, centros de custo, diário imutavel, regras de posting, razão, DRE/balanço, fechamento de período e reconciliação contábil operacionais | substituir regras genericas por políticas contábeis por regime, empresa e integração de evento contábil |
| `banking` | CNAB, boleto, extrato, conciliação, Pix cobrança/devolução/webhook e trilha de Open Finance operacionais | conectar bancos reais, assinatura de arquivos e conciliação automatica por provider |
| `inventory` | locais, saldos por SKU/local, ledger de movimentos, reservas, custo medio/FIFO, contagem cíclica e variancias operacionais | ampliar costing por metodo, integração fisica com WMS e políticas de alocação |
| `procurement` | requisição, cotação, pedido de compra, aprovação aplicada, recebimento e 3-way matching operacional | evoluir matriz de aprovação, contratos de fornecimento e avaliação continua |
| `documents` | metadata, upload sessions, assinatura, versões, retenção, arquivamento e links seguros operacionais | conectar storage real e varredura automatica de retenção |
| `fiscal` | perfis fiscais, documentos, fila de emissão, certificados, contingência, SPED, eventos, consentimentos, privacidade, retenção e auditoria operacional | conectar provider fiscal real, certificado digital gerenciado e regras por regime/UF/municipio |
| `platform-control` | capabilities, entitlements, quotas, metering, provider defaults, blocks, lifecycle, go-live, incident command, policies, approvals, runbooks, timeline, evidence vault, event mesh, tenant runtime e contract evolution operacionais | evoluir enforcement distribuido de quotas e automações de incidente por SLO real |
| `search` | busca operacional, facets, saved queries, auditoria, discovery cases, legal hold e exports controlados | conectar indexador incremental a eventos reais de todos os domínios |
| `ai-governance` | ferramentas aprovadas, políticas, runs auditáveis, redação e bloqueio de ferramentas mutantes/não registradas | conectar providers LLM reais mantendo policy enforcement e auditoria |
| `analytics` | reports executivos, tenant 360, service pulse, hardening, production-readiness, BI semântico, risk/compliance scoring, conciliação, fechamento financeiro, dados mestres, lakehouse, simulação, benchmark, estimativa de custo e read models operacionais quase em tempo real | evoluir streaming dedicado quando houver requisito comprovado de evento em baixa latencia |
| `edge` | cockpits consolidados com auth, health, go-live, SaaS, contratos e visoes cross-service, apoiados por gateway local com cache, rate limit, timeouts e failover passivo | evoluir replicas produtivas e políticas de tráfego por ambiente |
| `workflow-control` | definições, catálogos, versionamento, publish/restore, runs, eventos, ledger e diagnóstico por workflow operacionais | melhorar autoria visual e validação pre-públicação |
| `workflow-runtime` | execução durável, timeline, transições, delays, retries, capacidades e compensações basicas operacionais | ampliar observabilidade de execuções longas e cargas concorrentes |
| `engagement` | campanhas, templates, callbacks idempotentes, touchpoints, deliveries, provider events e conversas estruturadas operacionais | conectar adapters externos reais por canal |
| `webhook-hub` | inbound/outbound, assinatura preparada, retries, DLQ, requeue e ledger de transições operacionais | endurecer segurança de endpoint e políticas por tenant |
| `catalog` | itens, versões imutáveis, consumers, contratos de consumo e governança de produto operacionais | evoluir disponibilidade, pricing e políticas comerciais |
| `support` | casos, filas, SLA, comentários, exportação/bulk e resumo operacional | evoluir automações de atendimento e escalonamento |
| `supplier` | diretório, categorias, exportação/bulk e contratos administrativos operacionais | evoluir avaliação continua e relação com procurement sem duplicar ownership |
| `notification` | preferências, central, severidade, ciclo de vida e bulk operacional | conectar canais externos e templates avancados |
| `simulation` | catálogo de cenários, execução what-if, listagem de runs, benchmark de carga e insumos de sizing operacionais | alimentar planejamento com series históricas reais |
| `rentals` | contratos, charges, reajustes, encerramento, histórico, outbox, anexos e contrato HTTP completo para operação recorrente | evoluir regras avancadas de reajuste e integração contábil/fiscal |

## Checklist Para Novo serviço

Um serviço novo só deve entrar no monorepo se houver:

- responsabilidade funcional que não pertence claramente a serviço existente;
- contrato HTTP ou evento quando tiver consumidor externo ao módulo;
- caminho de código em `service-api/<stack>/<serviço>`;
- contexto PostgreSQL quando houver persistência;
- health/readiness quando exposto no runtime;
- teste unitario minimo;
- entrada no compose quando participar do stack local;
- entrada em `docs/SERVICOS.md`, `docs/API.md` e `docs/contracts/registry.json` quando aplicável.

## Checklist Para Evoluir serviço Existente

- A mudança respeita ownership?
- O OpenAPI foi atualizado?
- O banco mudou? Se sim, ha migration?
- Existe seed ou fixture afetada?
- O smoke precisa conhecer o novo comportamento?
- O `edge` ou `analytics` precisam refletir novo sinal?
- O `client-api` precisa regenerar catálogo?
- A documentação alterada está no arquivo certo?

## Sinais de Que a Fronteira está Errada

- O endpoint precisa consultar muitas tabelas de domínios diferentes para escrever um recurso.
- O serviço passa a conhecer detalhes internos de provider que outro serviço já encapsula.
- A alteração exige mudar varios serviços sem mudar contrato algum.
- Um report vira fonte de verdade.
- Uma tela nova dita o modelo de domínio sem contrato.

Quando isso acontecer, reavalie se a responsabilidade deveria estar em outro serviço, em `analytics`, em `edge`, ou em um contrato/evento novo.
