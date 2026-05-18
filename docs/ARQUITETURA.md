# ARQUITETURA

Este documento descreve a arquitetura do projeto. Ele não lista endpoints, não substitui runbook operacional e não detalha padrões de código. Para esses assuntos, use `docs/API.md`, `docs/OPERações.md` e `docs/PADROES.md`.

## Escopo

A arquitetura atual e uma plataforma backend-first, multi-tenant e poliglota, composta por serviços com ownership funcional claro. O monorepo agrupa os serviços, mas cada contexto mantem sua fronteira de API, persistência e responsabilidade.

## números de referência

- 26 serviços HTTP com contrato OpenAPI.
- 646 endpoints HTTP versionados.
- 17 schemas de evento versionados.
- Contratos em `docs/contracts/`.
- Runtime local em `infra/docker-compose.yml`.
- Deploy corporativo em `infra/kubernetes/`.
- Banco PostgreSQL dividido por contextos.

## princípios Arquiteturais

- Tenant explícito: operações tenant-aware não devem inferir tenant de forma invisível.
- Ownership por contexto: um serviço não escreve diretamente a verdade transacional de outro.
- Contrato primeiro: superfície pública relevante nasce ou evolui junto do OpenAPI/event schema.
- agregação separada de transação: leituras executivas ficam em `analytics` ou `edge`.
- Provider isolado: dependência externa passa por adapter, capability registry ou `webhook-hub`.
- operação verificável: health, readiness, smoke e hardening fazem parte da arquitetura, não são detalhe posterior.
- persistência por domínio: cada contexto PostgreSQL tem migrations e seeds próprios quando necessário.

## Planos

| Plano | serviços | Papel |
|-------|----------|-------|
| Public operations | `edge` | entrada pública e cockpits cross-service |
| Identity/security | `identity` | tenancy, usuário, sessão, role, convite e MFA |
| Commercial | `crm`, `sales` | relacionamento, pipeline, oportunidades, propostas e vendas |
| Supply chain | `inventory`, `procurement`, `supplier`, `catalog` | estoque, compras, fornecedores, recebimento e catálogo de itens |
| Recurrence/financial | `billing`, `finance`, `banking`, `accounting`, `rentals` | assinaturas, cobranças, bancos, ledger, recebíveis, comissões e contratos recorrentes |
| Documents/compliance | `documents`, `fiscal` | anexos, assinatura, documentos fiscais, emissão, certificados, privacidade, consentimento e auditoria |
| Automation | `workflow-control`, `workflow-runtime` | definicao, públicação e execução durável de workflows |
| Interaction/integration | `engagement`, `notification`, `webhook-hub` | comunicação, callbacks, webhooks e notificações |
| Platform governance | `platform-control`, `analytics`, `simulation`, `search`, `ai-governance` | capabilities, providers, quotas, lifecycle, incident command, reports, busca, IA governada e simulações |
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

`edge` oferece leituras operacionais consolidadas. Ele não deve virar dono das regras transacionais dos domínios.

`analytics` consolida relatórios e governança usando contratos, fixtures e leituras operacionais. Ele não substitui os serviços donos de escrita.

`client-web/client-api` e uma ferramenta técnica para explorar a API e a documentação; ele não e o frontend empresarial do produto.

O gateway local em `infra/gateway/nginx.conf` concentra roteamento `/api/<serviço>/`, health público, cache de leituras, rate limit, timeouts, correlação de request e failover passivo por dependência. Ele não substitui um API management corporativo completo, mas torna o stack local mais próximo de um ponto único de entrada verificável.

## Versão 1.4.6

A arquitetura da versão 1.4.6 prova que o projeto consegue consultar, validar e classificar sinais externos em modo BYOK/public API sem embutir chaves, sem simular provider produtivo e sem esconder quando uma capacidade está indisponível por falta de credencial. A mesma versão também acrescenta hardening de raiz/env, CI, IDs transacionais, console técnico, gateway, Kubernetes e conformance de auth/observabilidade.

O pacote arquitetural de 1.4.6 inclui:

- `infra/docker-compose.corporate-like.yml` para validar topologia com somente gateway/edge públicados;
- `infra/kubernetes/base` com namespace, config, secret template, service account, migration job, deployments, services, ingress TLS, NetworkPolicy e HPA;
- `infra/kubernetes/overlays/production` como entrada declarativa de deploy produtivo;
- `GET /api/analytics/reports/production-readiness` como evidência runtime do gate;
- `GET /api/search/query`, `GET /api/analytics/metrics`, `GET /api/ai-governance/tools` e `GET /api/platform-control/tenants/{tenantSlug}/incident-command/readiness` como evidências runtime das novas capacidades operacionais;
- `POST /api/platform-control/tenants/{tenantSlug}/policies/evaluate`, approvals, runbooks, timeline e evidence vault como plano de governança autonoma;
- `GET /api/analytics/risk/tenant-score` como evidência analitica de risco e compliance;
- `GET /api/platform-control/event-mesh/catalog`, dead letters, replay e lineage como malha de eventos enterprise;
- `GET /api/analytics/financial-close/readiness` e snapshots de fechamento como evidência financeira;
- `GET /api/analytics/master-data/quality-score` e manifesto lakehouse como governança de dados;
- `GET /api/platform-control/tenants/{tenantSlug}/runtime/profile` e contract evolution como controle runtime e evolucao segura;
- `GET /api/platform-control/providers/activation/catalog` para listar Stripe, Asaas, Mercado Pago, Resend, OpenAI, DocuSign, Clicksign e WhatsApp Cloud com postura de credencial;
- `POST /api/platform-control/tenants/{tenantSlug}/providers/activation/{providerKey}/test` para executar connection test ou ação suportada apenas quando a chave do operador existir;
- `POST /api/ai-governance/runs` usando OpenAI Responses API somente quando `OPENAI_API_KEY` estiver configurada, com fallback deterministico local em modo read-only;
- `GET /api/analytics/external-intelligence/readiness` como consolidado de OCR, fiscal Brasil, enriquecimento cadastral, mercado/macro e risco externo;
- `GET /api/analytics/document-intelligence/readiness` para AWS Textract e Google Document AI;
- `GET /api/analytics/fiscal-brazil/readiness` para Focus NFe, eNotas e certificado digital;
- `GET /api/analytics/registry-enrichment/brazil`, `market-macro-risk` e `external-risk-feed` como sinais operacionais governados;
- `./scripts/test.sh production-readiness` como validação estatica da versão;
- documentação operacional do aceite em `docs/OPERações.md`.

O princípio é deny-by-default: serviços internos não são pontos públicos de entrada, secrets não vivem no repositório, mutações sensíveis exigem identidade/correlação, e providers sem credencial real permanecem indisponíveis ou em fallback local explicitamente declarado.

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
| `platform-control` | capabilities, entitlements, quotas, lifecycle, go-live, incident command, policies, approvals, runbooks, timeline, evidence, event mesh, tenant runtime, contract evolution |
| `search` | operational search, e-discovery, legal hold, exports |
| `ai-governance` | approved tools, policies, assistant audit, redaction |
| `support` | queues, cases, comments, SLA |
| `supplier` | supplier categories and suppliers |
| `notification` | preferences and notification center |
| `fiscal` | fiscal profiles, documents, retention, consent, privacy, audit |
| `analytics` | derived reports, semantic BI, risk scoring, reconciliation, financial close, master data quality and lakehouse views |
| `simulation` | scenários and benchmark results |

## Contratos Como Fronteira

Contratos ficam em `docs/contracts/` porque são artefatos de governança e interoperabilidade. Eles não pertencem a `infra/`, pois não são infraestrutura de runtime, e não pertencem a um único diretório de serviço, pois também orientam consumidores, smoke, console técnico e revisão de compatibilidade.

### decisão: HTTP Interno Antes de gRPC

O padrão oficial do projeto, no baseline atual, é usar HTTP versionado com OpenAPI para a superfície interna/pública entre serviços e consumidores técnicos.

Essa decisão acompanha a realidade do projeto: o monorepo é poliglota, com serviços em Go, .NET, Elixir, TypeScript, Python e Rust. A plataforma precisa de contratos navegáveis, testes simples em ambiente local, onboarding rápido e uma forma objetiva de expor a API completa no `client-web/client-api`.

gRPC continua sendo uma opcao técnica valida para comunicação interna de alto volume, streaming ou contratos binarios fortemente tipados entre serviços controlados. Mesmo assim, o custo operacional de introduzir gRPC em todos os contextos agora seria maior que o ganho imediato.

Consequências praticas:

- cada serviço com superfície pública mantem OpenAPI em `docs/contracts/http/`;
- endpoint novo relevante nasce com contrato ou atualiza contrato existente;
- mudança de contrato deve ser compatível ou explicitamente versionada;
- operações longas usam `202 Accepted`, recurso de job/rollout/execution e endpoint de leitura;
- Mutações sensíveis usam `Idempotency-Key` quando houver risco de duplicidade;
- eventos compartilhados continuam em JSON Schema dentro de `docs/contracts/events/`;
- `client-web/client-api` pode gerar catálogo navegavel a partir dos contratos HTTP;
- `./scripts/test.sh contract` valida registry, schema registry e portal;
- consumidores humanos e assistidos por LLM conseguem inspecionar contratos sem tooling especifico de gRPC.

A decisão deve ser reavaliada quando houver streaming real entre serviços, gargalo comprovado de chamadas internas de altissimo volume, necessidade concreta de contrato binario entre dois serviços controlados ou requisito de multiplexação/performance que o HTTP atual não resolva bem.

## comunicação Entre Contextos

- HTTP versionado e usado para superfície pública e consumo sincronico.
- Eventos versionados são usados quando a informação e compartilhada como fato entre contextos.
- Webhooks externos entram por endpoints de provider ou pelo `webhook-hub`.
- operações longas usam recursos de job, rollout ou execution.
- Agregações devem preferir `analytics`/`edge` a leitura direta de tabela alheia.

## confiança Entre serviços

O padrão de chamadas internas é service-account first: cada serviço chamador deve usar token curto com audience do destino, correlation id obrigatório, tenant explícito e actor propagado quando a ação nasceu de um usuário. A rede interna reduz exposição, mas não é autorização suficiente para Mutações sensíveis.

Regras atuais:

- `bearerAuth` representa sessão/OIDC de usuário nos OpenAPI.
- `internalServiceToken` representa token de service account entre workloads.
- `X-Correlation-Id` deve acompanhar chamadas e eventos cross-service.
- `tenantSlug` não deve cair em bootstrap fora de `local`/`test`.
- rotas financeiras, fiscais, documentais, identity, platform-control e webhook-hub devem ser tratadas como sensíveis mesmo quando chamadas internamente.

Para um ambiente corporativo, `infra/docker-compose.corporate-like.yml` remove a públicação direta dos serviços de domínio e mantem o tráfego no gateway/edge. mTLS ou service mesh seguem como evolucao natural quando houver runtime produtivo com certificados e identidade de workload gerenciados.

## persistência

O PostgreSQL local e compartilhado como infraestrutura, mas a propriedade logica e separada por contexto. Migrations vivem em `service-api/service-postgresql/<contexto>/migrations`.

Essa escolha facilita:

- smoke local de ponta a ponta;
- backup/restore uniforme;
- evolucao independente por domínio;
- leitura operacional controlada;
- reducao de acoplamento entre tabelas de serviços diferentes.

## Fronteiras Que não Devem Ser Misturadas

- `docs/API.md` descreve API, não arquitetura geral.
- `docs/SERVICOS.md` descreve serviços e ownership, não cada endpoint em detalhe.
- `docs/OPERações.md` descreve comandos e runbooks, não regras de domínio.
- `docs/CONTRATOS.md` descreve governança contratual, não implementação.
- `docs/PADROES.md` descreve padrões de engenharia, não estado de produto.

## Decisoes Arquiteturais Atuais

### Monorepo com fronteiras fortes

O repositório e único para facilitar evolucao coordenada, smoke local e governança de contratos. Isso não significa que os serviços possam depender livremente uns dos outros. A fronteira real e formada por:

- contrato HTTP;
- schema de evento;
- schema PostgreSQL de ownership;
- runtime do serviço;
- suite de testes que valida aquela superfície.

Essa escolha evita o custo de varios repositórios no momento atual e ainda preserva disciplina de arquitetura distribuida.

### Poliglotismo orientado por domínio

O projeto usa stacks diferentes porque os contextos evoluíram com necessidades diferentes. A regra arquitetural não é "cada serviço precisa de uma linguagem nova"; a regra é "cada serviço deve ser previsível na stack que já usa".

O custo do poliglotismo e compensado por:

- contratos HTTP comuns;
- Docker Compose como runtime uniforme;
- PostgreSQL como persistência local comum;
- `scripts/build.sh` e `scripts/test.sh` como entradas padronizadas;
- documentação central em `docs/`.

| Stack | serviços atuais | Uso preferêncial | Evitar quando |
| --- | --- | --- | --- |
| Go | `edge`, `crm`, `sales`, `documents`, `rentals` | APIs leves, IO direto, gateways e domínios com handlers simples | regra já depende fortemente de .NET ou Python |
| C#/.NET | `identity`, `billing`, `finance` | domínios transacionais, contratos ricos e forte tipagem | scripts pequenos ou adapters descartaveis |
| Python/FastAPI | `accounting`, `ai-governance`, `analytics`, `banking`, `catalog`, `fiscal`, `inventory`, `notification`, `platform-control`, `procurement`, `search`, `simulation`, `supplier`, `support` | reports, catálogos operacionais, busca, IA governada, simulação, integrações simples e domínios administrativos parametrizados | transação financeira critica sem camada forte de teste |
| TypeScript | `workflow-control`, `engagement` | orquestração HTTP, catálogo de workflows e integrações web/eventos | processamento numerico/financeiro crítico |
| Rust | `webhook-hub` | ingestao/eventos com foco em robustez e controle de runtime | CRUD simples com alta mudança de regra |
| Elixir | `workflow-runtime` | runtime concorrente, timers, retries e processos de longa duração | rotas CRUD convencionais |

Nova linguagem só entra se houver ganho claro e se auth, tracing, contrato, teste, build e dependency scan ficarem equivalentes aos stacks existentes.

### Edge como leitura operacional

`edge` existe para expor cockpits e leituras consolidadas. Ele pode chamar ou compor dados de outros domínios, mas não deve se tornar o lugar onde regras de `billing`, `finance`, `crm` ou `platform-control` são implementadas.

Quando uma regra nasce no `edge`, ela precisa ser questionada:

- e uma regra de apresentação/agregação?
- ou e uma regra transacional que pertence a outro serviço?

Se for transacional, deve voltar para o dono do domínio.

### Analytics como derivação

`analytics` produz leituras derivadas: readiness, governança, hardening, compliance, custos, go-live, catálogo semântico de métricas e scoring de risco. Ele pode correlacionar sinais, mas não deve assumir ownership de escrita de dados operacionais.

O papel dele e responder perguntas como:

- quais contratos estão cobertos?
- quais providers estão configurados?
- quais riscos bloqueiam go-live?
- quais domínios exigem atencao operacional?

### Platform-control como plano SaaS

`platform-control` e o serviço que centraliza capabilities, providers, entitlements, feature flags, quotas, metering, blocks, lifecycle, go-live, incident command, policy decisions, approvals, runbooks, timeline e evidências por tenant. Ele e o lugar certo para governança de plataforma e operação, mas não para regras comerciais como venda, cobrança ou documento fiscal.

## Fluxos Arquiteturais

### Login e tenancy

```text
client
  -> identity
  -> session/token
  -> chamadas tenant-aware aos domínios
```

`identity` e o dono da identidade. Outros serviços podem receber tenant/ator, mas não recriam usuário, role ou MFA.

### Venda ate financeiro

```text
crm
  -> sales
  -> billing
  -> finance
  -> accounting
  -> analytics/edge
```

`crm` qualifica relacionamento, `sales` registra operação comercial, `billing` controla recorrência e cobrança, `finance` consolida consequências financeiras, e `accounting` registra impactos gerenciais. `analytics` e `edge` observam.

### Documento e compliance

```text
documents
  -> fiscal
  -> banking/accounting quando houver conciliação
  -> analytics compliance-control
  -> edge compliance-overview
```

`documents` guarda metadata documental e assinatura. `fiscal` governa documento fiscal, emissão, certificado, retenção, consentimento, privacidade e auditoria. `banking` e `accounting` entram quando ha conciliação ou reflexo financeiro/contábil.

### Automação

```text
workflow-control
  -> workflow-runtime
  -> serviços de domínio
  -> timeline/actions
```

O controle define o fluxo. O runtime executa e registra estado. Os domínios chamados continuam donos de suas regras.

### Webhook

```text
provider externo
  -> engagement ou webhook-hub
  -> persistência normalizada
  -> DLQ/requeue quando aplicável
  -> analytics/edge para leitura operacional
```

Callback externo sem idempotência e sem trilha de tentativa e risco arquitetural.

## dependências Aceitas

| Origem | Pode depender de | Condicao |
|--------|------------------|----------|
| `edge` | domínios e `analytics` | apenas leitura/agregação operacional |
| `analytics` | contratos, fixtures e leituras operacionais | sem ownership transacional |
| domínios transacionais | `identity` | para contexto de tenant/ator quando necessário |
| `workflow-runtime` | domínios | por comando/contrato explícito |
| `webhook-hub` | consumers externos | por endpoint registrado e delivery log |
| `client-api` | OpenAPI/docs/backend local | como ferramenta técnica, não runtime do produto |

## dependências Que Devem Ser Evitadas

- serviço escrevendo tabela de outro contexto.
- Endpoint usando ID interno de banco como contrato público.
- Analytics corrigindo dado transacional.
- Edge implementando regra de negocio de domínio.
- Workflow alterando estado sem passar pelo contrato do serviço dono.
- Provider externo chamado diretamente por varios serviços sem adapter/capability comum.
- documentação de arquitetura tentando substituir OpenAPI ou runbook.

## Evolucao Esperada

A arquitetura atual ainda e de plataforma em construcao, não de produto final fechado. Evolucoes provaveis:

- endurecer autenticação/autorização em todos os endpoints públicos;
- aumentar cobertura real de events/outbox onde houver integração assíncrona critica;
- transformar mais leituras de `analytics` em read models persistidos quando o custo de consulta crescer;
- separar frontend empresarial de `client-web/client-api`;
- aumentar observabilidade distribuida com tracing e correlation id de ponta a ponta;
- versionar públicamente contratos quando houver consumidores externos reais.

Essas evolucoes devem preservar a regra principal: cada domínio continua dono do próprio estado e a integração acontece por contrato.
