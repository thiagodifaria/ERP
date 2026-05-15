# ARQUITETURA

Este documento descreve a arquitetura do ERP. Ele nao lista endpoints, nao substitui runbook operacional e nao detalha padroes de codigo. Para esses assuntos, use `docs/API.md`, `docs/OPERACOES.md` e `docs/PADROES.md`.

## Escopo

A arquitetura atual e uma plataforma backend-first, multi-tenant e poliglota, composta por servicos com ownership funcional claro. O monorepo agrupa os servicos, mas cada contexto mantem sua fronteira de API, persistencia e responsabilidade.

## Numeros de Referencia

- 24 servicos HTTP com contrato OpenAPI.
- 543 endpoints HTTP versionados.
- 15 schemas de evento versionados.
- Contratos em `docs/contracts/`.
- Runtime local em `infra/docker-compose.yml`.
- Deploy corporativo em `infra/kubernetes/`.
- Banco PostgreSQL dividido por contextos.

## Principios Arquiteturais

- Tenant explicito: operacoes tenant-aware nao devem inferir tenant de forma invisivel.
- Ownership por contexto: um servico nao escreve diretamente a verdade transacional de outro.
- Contrato primeiro: superficie publica relevante nasce ou evolui junto do OpenAPI/event schema.
- Agregacao separada de transacao: leituras executivas ficam em `analytics` ou `edge`.
- Provider isolado: dependencia externa passa por adapter, capability registry ou `webhook-hub`.
- Operacao verificavel: health, readiness, smoke e hardening fazem parte da arquitetura, nao sao detalhe posterior.
- Persistencia por dominio: cada contexto PostgreSQL tem migrations e seeds proprios quando necessario.

## Planos

| Plano | Servicos | Papel |
|-------|----------|-------|
| Public operations | `edge` | entrada publica e cockpits cross-service |
| Identity/security | `identity` | tenancy, usuario, sessao, role, convite e MFA |
| Commercial | `crm`, `sales` | relacionamento, pipeline, oportunidades, propostas e vendas |
| Supply chain | `inventory`, `procurement`, `supplier`, `catalog` | estoque, compras, fornecedores, recebimento e catalogo de itens |
| Recurrence/financial | `billing`, `finance`, `banking`, `accounting`, `rentals` | assinaturas, cobrancas, bancos, ledger, recebiveis, comissoes e contratos recorrentes |
| Documents/compliance | `documents`, `fiscal` | anexos, assinatura, documentos fiscais, emissao, certificados, privacidade, consentimento e auditoria |
| Automation | `workflow-control`, `workflow-runtime` | definicao, publicacao e execucao duravel de workflows |
| Interaction/integration | `engagement`, `notification`, `webhook-hub` | comunicacao, callbacks, webhooks e notificacoes |
| Platform governance | `platform-control`, `analytics`, `simulation` | capabilities, providers, quotas, lifecycle, reports e simulacoes |
| Administrative | `support` | suporte e atendimento operacional |

## Topologia de Runtime

```text
client / client-api
        |
        v
   local gateway
        |
        v
      edge  ----------------------+
        |                         |
        v                         v
 domain services              analytics
        |                         |
        v                         v
 PostgreSQL contexts       operational reads
```

`edge` oferece leituras operacionais consolidadas. Ele nao deve virar dono das regras transacionais dos dominios.

`analytics` consolida relatorios e governanca usando contratos, fixtures e leituras operacionais. Ele nao substitui os servicos donos de escrita.

`client-web/client-api` e uma ferramenta tecnica para explorar a API e a documentacao; ele nao e o frontend empresarial do produto.

O gateway local em `infra/gateway/nginx.conf` concentra roteamento `/api/<servico>/`, health publico, cache de leituras, rate limit, timeouts, correlacao de request e failover passivo por dependencia. Ele nao substitui um API management corporativo completo, mas torna o stack local mais proximo de um ponto unico de entrada verificavel.

## Release 1.0.0

A arquitetura da versao 1.0.0 e orientada por production readiness. O objetivo do release e provar que o ERP nao e apenas amplo em dominios, mas tambem operavel em ambiente corporativo.

O pacote arquitetural de 1.0.0 inclui:

- `infra/docker-compose.corporate-like.yml` para validar topologia com somente gateway/edge publicados;
- `infra/kubernetes/base` com namespace, config, secret template, service account, migration job, deployments, services, ingress TLS, NetworkPolicy e HPA;
- `infra/kubernetes/overlays/production` como entrada declarativa de deploy produtivo;
- `GET /api/analytics/reports/production-readiness` como evidencia runtime do gate;
- `./scripts/test.sh production-readiness` como validacao estatica de release;
- documentacao operacional do aceite em `docs/OPERACOES.md`.

O principio e deny-by-default: servicos internos nao sao pontos publicos de entrada, secrets nao vivem no repositorio, mutacoes sensiveis exigem identidade/correlacao, e providers sem credencial real nao sao apresentados como produtivos.

## Ownership de Dados

| Contexto | Dono principal |
|----------|----------------|
| `identity` | tenants, companies, users, roles, sessions, invitations |
| `crm` | leads, customers, pipeline, enrichment |
| `sales` | opportunities, proposals, sales, invoices |
| `inventory` | locations, balances, movements, reservations, FIFO/average costing, cycle counts |
| `procurement` | requisitions, quotations, purchase orders, approvals, receiving, 3-way matching |
| `rentals` | rental contracts and charges |
| `billing` | plans, subscriptions, invoices, payment attempts |
| `finance` | receivables, projections, commission holds, activity |
| `accounting` | chart of accounts, cost centers, immutable journal entries, posting rules, ledger, close, statements |
| `banking` | CNAB, boletos, statements, reconciliation, Pix refunds/webhooks, Open Finance |
| `documents` | attachments, storage metadata, signing requests, versions |
| `engagement` | providers, touchpoints, conversations, provider events |
| `workflow-control` | definitions, catalogs, control-plane runs |
| `workflow-runtime` | executions, timeline, actions, transitions |
| `webhook-hub` | inbound events, outbound endpoints, deliveries, DLQ |
| `catalog` | categories, items, versions, consumer contracts |
| `platform-control` | capabilities, entitlements, quotas, lifecycle, go-live |
| `support` | queues, cases, comments, SLA |
| `supplier` | supplier categories and suppliers |
| `notification` | preferences and notification center |
| `fiscal` | fiscal profiles, documents, retention, consent, privacy, audit |
| `analytics` | derived reports and governance views |
| `simulation` | scenarios and benchmark results |

## Contratos Como Fronteira

Contratos ficam em `docs/contracts/` porque sao artefatos de governanca e interoperabilidade. Eles nao pertencem a `infra/`, pois nao sao infraestrutura de runtime, e nao pertencem a um unico diretorio de servico, pois tambem orientam consumidores, smoke, console tecnico e revisao de compatibilidade.

### Decisao: HTTP Interno Antes de gRPC

O padrao oficial do projeto, no baseline atual, e usar HTTP versionado com OpenAPI para a superficie interna/publica entre servicos e consumidores tecnicos.

Essa decisao acompanha a realidade do ERP: o monorepo e poliglota, com servicos em Go, .NET, Elixir, TypeScript, Python e Rust. A plataforma precisa de contratos navegaveis, testes simples em ambiente local, onboarding rapido e uma forma objetiva de expor a API completa no `client-web/client-api`.

gRPC continua sendo uma opcao tecnica valida para comunicacao interna de alto volume, streaming ou contratos binarios fortemente tipados entre servicos controlados. Mesmo assim, o custo operacional de introduzir gRPC em todos os contextos agora seria maior que o ganho imediato.

Consequencias praticas:

- cada servico com superficie publica mantem OpenAPI em `docs/contracts/http/`;
- endpoint novo relevante nasce com contrato ou atualiza contrato existente;
- mudanca de contrato deve ser compativel ou explicitamente versionada;
- operacoes longas usam `202 Accepted`, recurso de job/rollout/execution e endpoint de leitura;
- mutacoes sensiveis usam `Idempotency-Key` quando houver risco de duplicidade;
- eventos compartilhados continuam em JSON Schema dentro de `docs/contracts/events/`;
- `client-web/client-api` pode gerar catalogo navegavel a partir dos contratos HTTP;
- `./scripts/test.sh contract` valida registry, schema registry e portal;
- consumidores humanos e assistidos por LLM conseguem inspecionar contratos sem tooling especifico de gRPC.

A decisao deve ser reavaliada quando houver streaming real entre servicos, gargalo comprovado de chamadas internas de altissimo volume, necessidade concreta de contrato binario entre dois servicos controlados ou requisito de multiplexacao/performance que o HTTP atual nao resolva bem.

## Comunicacao Entre Contextos

- HTTP versionado e usado para superficie publica e consumo sincronico.
- Eventos versionados sao usados quando a informacao e compartilhada como fato entre contextos.
- Webhooks externos entram por endpoints de provider ou pelo `webhook-hub`.
- Operacoes longas usam recursos de job, rollout ou execution.
- Agregacoes devem preferir `analytics`/`edge` a leitura direta de tabela alheia.

## Confianca Entre Servicos

O padrao de chamadas internas e service-account first: cada servico chamador deve usar token curto com audience do destino, correlation id obrigatorio, tenant explicito e actor propagado quando a acao nasceu de um usuario. A rede interna reduz exposicao, mas nao e autorizacao suficiente para mutacoes sensiveis.

Regras atuais:

- `bearerAuth` representa sessao/OIDC de usuario nos OpenAPI.
- `internalServiceToken` representa token de service account entre workloads.
- `X-Correlation-Id` deve acompanhar chamadas e eventos cross-service.
- `tenantSlug` nao deve cair em bootstrap fora de `local`/`test`.
- rotas financeiras, fiscais, documentais, identity, platform-control e webhook-hub devem ser tratadas como sensiveis mesmo quando chamadas internamente.

Para um ambiente corporativo, `infra/docker-compose.corporate-like.yml` remove a publicacao direta dos servicos de dominio e mantem o trafego no gateway/edge. mTLS ou service mesh seguem como evolucao natural quando houver runtime produtivo com certificados e identidade de workload gerenciados.

## Persistencia

O PostgreSQL local e compartilhado como infraestrutura, mas a propriedade logica e separada por contexto. Migrations vivem em `service-api/service-postgresql/<contexto>/migrations`.

Essa escolha facilita:

- smoke local de ponta a ponta;
- backup/restore uniforme;
- evolucao independente por dominio;
- leitura operacional controlada;
- reducao de acoplamento entre tabelas de servicos diferentes.

## Fronteiras Que Nao Devem Ser Misturadas

- `docs/API.md` descreve API, nao arquitetura geral.
- `docs/SERVICOS.md` descreve servicos e ownership, nao cada endpoint em detalhe.
- `docs/OPERACOES.md` descreve comandos e runbooks, nao regras de dominio.
- `docs/CONTRATOS.md` descreve governanca contratual, nao implementacao.
- `docs/PADROES.md` descreve padroes de engenharia, nao estado de produto.

## Decisoes Arquiteturais Atuais

### Monorepo com fronteiras fortes

O repositorio e unico para facilitar evolucao coordenada, smoke local e governanca de contratos. Isso nao significa que os servicos possam depender livremente uns dos outros. A fronteira real e formada por:

- contrato HTTP;
- schema de evento;
- schema PostgreSQL de ownership;
- runtime do servico;
- suite de testes que valida aquela superficie.

Essa escolha evita o custo de varios repositorios no momento atual e ainda preserva disciplina de arquitetura distribuida.

### Poliglotismo orientado por dominio

O projeto usa stacks diferentes porque os contextos evoluiram com necessidades diferentes. A regra arquitetural nao e "cada servico precisa de uma linguagem nova"; a regra e "cada servico deve ser previsivel na stack que ja usa".

O custo do poliglotismo e compensado por:

- contratos HTTP comuns;
- Docker Compose como runtime uniforme;
- PostgreSQL como persistencia local comum;
- `scripts/build.sh` e `scripts/test.sh` como entradas padronizadas;
- documentacao central em `docs/`.

| Stack | Servicos atuais | Uso preferencial | Evitar quando |
| --- | --- | --- | --- |
| Go | `edge`, `crm`, `sales`, `documents`, `rentals` | APIs leves, IO direto, gateways e dominios com handlers simples | regra ja depende fortemente de .NET ou Python |
| C#/.NET | `identity`, `billing`, `finance` | dominios transacionais, contratos ricos e forte tipagem | scripts pequenos ou adapters descartaveis |
| Python/FastAPI | `accounting`, `analytics`, `banking`, `catalog`, `fiscal`, `inventory`, `notification`, `platform-control`, `procurement`, `simulation`, `supplier`, `support` | reports, catalogos operacionais, simulacao, integracoes simples e dominios administrativos parametrizados | transacao financeira critica sem camada forte de teste |
| TypeScript | `workflow-control`, `engagement` | orquestracao HTTP, catalogo de workflows e integracoes web/eventos | processamento numerico/financeiro critico |
| Rust | `webhook-hub` | ingestao/eventos com foco em robustez e controle de runtime | CRUD simples com alta mudanca de regra |
| Elixir | `workflow-runtime` | runtime concorrente, timers, retries e processos de longa duracao | rotas CRUD convencionais |

Nova linguagem so entra se houver ganho claro e se auth, tracing, contrato, teste, build e dependency scan ficarem equivalentes aos stacks existentes.

### Edge como leitura operacional

`edge` existe para expor cockpits e leituras consolidadas. Ele pode chamar ou compor dados de outros dominios, mas nao deve se tornar o lugar onde regras de `billing`, `finance`, `crm` ou `platform-control` sao implementadas.

Quando uma regra nasce no `edge`, ela precisa ser questionada:

- e uma regra de apresentacao/agregacao?
- ou e uma regra transacional que pertence a outro servico?

Se for transacional, deve voltar para o dono do dominio.

### Analytics como derivacao

`analytics` produz leituras derivadas: readiness, governanca, hardening, compliance, custos e go-live. Ele pode correlacionar sinais, mas nao deve assumir ownership de escrita de dados operacionais.

O papel dele e responder perguntas como:

- quais contratos estao cobertos?
- quais providers estao configurados?
- quais riscos bloqueiam go-live?
- quais dominios exigem atencao operacional?

### Platform-control como plano SaaS

`platform-control` e o servico que centraliza capabilities, providers, entitlements, feature flags, quotas, metering, blocks, lifecycle e go-live por tenant. Ele e o lugar certo para governanca de plataforma, mas nao para regras comerciais como venda, cobranca ou documento fiscal.

## Fluxos Arquiteturais

### Login e tenancy

```text
client
  -> identity
  -> session/token
  -> chamadas tenant-aware aos dominios
```

`identity` e o dono da identidade. Outros servicos podem receber tenant/ator, mas nao recriam usuario, role ou MFA.

### Venda ate financeiro

```text
crm
  -> sales
  -> billing
  -> finance
  -> accounting
  -> analytics/edge
```

`crm` qualifica relacionamento, `sales` registra operacao comercial, `billing` controla recorrencia e cobranca, `finance` consolida consequencias financeiras, e `accounting` registra impactos gerenciais. `analytics` e `edge` observam.

### Documento e compliance

```text
documents
  -> fiscal
  -> banking/accounting quando houver conciliacao
  -> analytics compliance-control
  -> edge compliance-overview
```

`documents` guarda metadata documental e assinatura. `fiscal` governa documento fiscal, emissao, certificado, retencao, consentimento, privacidade e auditoria. `banking` e `accounting` entram quando ha conciliacao ou reflexo financeiro/contabil.

### Automacao

```text
workflow-control
  -> workflow-runtime
  -> servicos de dominio
  -> timeline/actions
```

O controle define o fluxo. O runtime executa e registra estado. Os dominios chamados continuam donos de suas regras.

### Webhook

```text
provider externo
  -> engagement ou webhook-hub
  -> persistencia normalizada
  -> DLQ/requeue quando aplicavel
  -> analytics/edge para leitura operacional
```

Callback externo sem idempotencia e sem trilha de tentativa e risco arquitetural.

## Dependencias Aceitas

| Origem | Pode depender de | Condicao |
|--------|------------------|----------|
| `edge` | dominios e `analytics` | apenas leitura/agregacao operacional |
| `analytics` | contratos, fixtures e leituras operacionais | sem ownership transacional |
| dominios transacionais | `identity` | para contexto de tenant/ator quando necessario |
| `workflow-runtime` | dominios | por comando/contrato explicito |
| `webhook-hub` | consumers externos | por endpoint registrado e delivery log |
| `client-api` | OpenAPI/docs/backend local | como ferramenta tecnica, nao runtime do produto |

## Dependencias Que Devem Ser Evitadas

- Servico escrevendo tabela de outro contexto.
- Endpoint usando ID interno de banco como contrato publico.
- Analytics corrigindo dado transacional.
- Edge implementando regra de negocio de dominio.
- Workflow alterando estado sem passar pelo contrato do servico dono.
- Provider externo chamado diretamente por varios servicos sem adapter/capability comum.
- Documentacao de arquitetura tentando substituir OpenAPI ou runbook.

## Evolucao Esperada

A arquitetura atual ainda e de plataforma em construcao, nao de produto final fechado. Evolucoes provaveis:

- endurecer autenticacao/autorizacao em todos os endpoints publicos;
- aumentar cobertura real de events/outbox onde houver integracao assíncrona critica;
- transformar mais leituras de `analytics` em read models persistidos quando o custo de consulta crescer;
- separar frontend empresarial de `client-web/client-api`;
- aumentar observabilidade distribuida com tracing e correlation id de ponta a ponta;
- versionar publicamente contratos quando houver consumidores externos reais.

Essas evolucoes devem preservar a regra principal: cada dominio continua dono do proprio estado e a integracao acontece por contrato.
