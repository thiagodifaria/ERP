# ARQUITETURA

Este documento descreve a arquitetura do ERP. Ele nao lista endpoints, nao substitui runbook operacional e nao detalha padroes de codigo. Para esses assuntos, use `docs/API.md`, `docs/OPERACOES.md` e `docs/PADROES.md`.

## Escopo

A arquitetura atual e uma plataforma backend-first, multi-tenant e poliglota, composta por servicos com ownership funcional claro. O monorepo agrupa os servicos, mas cada contexto mantem sua fronteira de API, persistencia e responsabilidade.

## Numeros de Referencia

- 20 servicos HTTP com contrato OpenAPI.
- 201 endpoints HTTP versionados.
- 12 schemas de evento versionados.
- Contratos em `docs/contracts/`.
- Runtime local em `infra/docker-compose.yml`.
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
| Recurrence/financial | `billing`, `finance`, `rentals` | assinaturas, cobrancas, recebiveis, comissoes e contratos recorrentes |
| Documents/compliance | `documents`, `fiscal` | anexos, assinatura, documentos fiscais, privacidade, consentimento e auditoria |
| Automation | `workflow-control`, `workflow-runtime` | definicao, publicacao e execucao duravel de workflows |
| Interaction/integration | `engagement`, `notification`, `webhook-hub` | comunicacao, callbacks, webhooks e notificacoes |
| Platform governance | `platform-control`, `analytics`, `simulation` | capabilities, providers, quotas, lifecycle, reports e simulacoes |
| Administrative | `support`, `supplier`, `catalog` | suporte, fornecedores, catalogo e contratos de consumo |

## Topologia de Runtime

```text
client / client-api
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

## Ownership de Dados

| Contexto | Dono principal |
|----------|----------------|
| `identity` | tenants, companies, users, roles, sessions, invitations |
| `crm` | leads, customers, pipeline, enrichment |
| `sales` | opportunities, proposals, sales, invoices |
| `rentals` | rental contracts and charges |
| `billing` | plans, subscriptions, invoices, payment attempts |
| `finance` | receivables, projections, commission holds, activity |
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

## Comunicacao Entre Contextos

- HTTP versionado e usado para superficie publica e consumo sincronico.
- Eventos versionados sao usados quando a informacao e compartilhada como fato entre contextos.
- Webhooks externos entram por endpoints de provider ou pelo `webhook-hub`.
- Operacoes longas usam recursos de job, rollout ou execution.
- Agregacoes devem preferir `analytics`/`edge` a leitura direta de tabela alheia.

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

Essa escolha evita o custo de varios repositorios nesta fase e ainda preserva disciplina de arquitetura distribuida.

### Poliglotismo orientado por dominio

O projeto usa stacks diferentes porque os contextos evoluiram com necessidades diferentes. A regra arquitetural nao e "cada servico precisa de uma linguagem nova"; a regra e "cada servico deve ser previsivel na stack que ja usa".

O custo do poliglotismo e compensado por:

- contratos HTTP comuns;
- Docker Compose como runtime uniforme;
- PostgreSQL como persistencia local comum;
- `scripts/build.sh` e `scripts/test.sh` como entradas padronizadas;
- documentacao central em `docs/`.

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
  -> analytics/edge
```

`crm` qualifica relacionamento, `sales` registra operacao comercial, `billing` controla recorrencia e cobranca, `finance` consolida consequencias financeiras. `analytics` e `edge` observam.

### Documento e compliance

```text
documents
  -> fiscal
  -> analytics compliance-control
  -> edge compliance-overview
```

`documents` guarda metadata documental e assinatura. `fiscal` governa documento fiscal, retencao, consentimento, privacidade e auditoria.

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
