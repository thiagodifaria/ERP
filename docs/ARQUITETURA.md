# ARQUITETURA

## Objetivo

Este documento registra a arquitetura do ERP com profundidade suficiente para orientar evolucao real do monorepo. A arquitetura nao e descrita como desenho decorativo: cada decisao abaixo existe para proteger ownership, reduzir acoplamento, manter contratos versionados e permitir operacao local confiavel.

## Escala atual

- Servicos HTTP com contrato versionado: `20`
- Endpoints versionados: `201`
- Schemas de evento versionados: `12`
- Contratos: `docs/contracts/`
- Runtime local: `infra/docker-compose.yml`
- Comando operacional: `./scripts/build.sh`
- Validacao central: `./scripts/test.sh`

## Principios arquiteturais

- O tenant e contexto operacional explicito.
- Cada servico tem ownership funcional claro.
- Cada schema PostgreSQL pertence a um contexto.
- Contratos sao artefatos de governanca versionada.
- Health e readiness precisam refletir dependencias reais.
- Smoke integrado valida comportamento de plataforma, nao apenas compilacao.
- Adapters externos devem declarar postura configurada, fallback, manual, disabled ou unconfigured.
- A documentacao acompanha a arquitetura e nao fica espalhada por README de servico.

## Decisao: localizacao dos contratos

Contratos ficam em `docs/contracts/`. Eles nao ficam em `infra/` porque nao sao runtime infrastructure como compose, Prometheus, Grafana, Keycloak ou manifests operacionais. Tambem nao ficam em `service-api/` porque nao sao codigo de implementacao de servico. Aqui os contratos sao tratados como fonte versionada de governanca tecnica, descoberta de API, schema de eventos e baseline de compatibilidade.

A decisao pratica e simples: codigo implementa comportamento em `service-api/`; infraestrutura executa ambiente em `infra/`; documentacao e contratos que explicam e governam a superficie publica ficam em `docs/`.

## Planos arquiteturais

### administrative plane

`support`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### administrative/notification plane

`notification`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### administrative/procurement plane

`supplier`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### analytics plane

`analytics`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### compliance plane

`fiscal`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### control plane

`workflow-control`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### integration plane

`webhook-hub`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### interaction/control plane

`engagement`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### platform control plane

`platform-control`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### public operations plane

`edge`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### runtime plane

`workflow-runtime`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### simulation plane

`simulation`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### transaction plane

`crm`, `documents`, `rentals`, `sales`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### transaction/billing plane

`billing`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### transaction/catalog plane

`catalog`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### transaction/finance plane

`finance`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

### transaction/security plane

`identity`

Responsabilidades:

- manter fronteira operacional clara;
- reduzir acoplamento entre dominio e consumidor;
- permitir evolucao independente por contrato;
- facilitar smoke e diagnostico local;

## Fluxos macro

### Requisicao HTTP externa

1. A requisicao chega no gateway ou diretamente no `edge` em ambiente local.
2. O contexto de tenant e autenticacao e resolvido quando a rota exige protecao.
3. O servico dono do dominio executa a regra.
4. O schema PostgreSQL do contexto persiste a verdade transacional quando existe persistencia.
5. Eventos, historico ou outbox sao registrados quando ha consumidor operacional.
6. Analytics e edge agregam leitura sem virar donos do dado transacional.

### Evento externo

1. O provider externo chama `webhook-hub`.
2. A entrada e normalizada por provider e external id.
3. O evento passa por transicoes auditaveis.
4. Falhas persistentes seguem para dead letter.
5. Consumidores internos usam contrato, nao acesso informal ao payload cru.

### Automacao

1. `workflow-control` define e publica versoes.
2. Runs registram intencao operacional.
3. `workflow-runtime` executa actions, waits, retries e compensacoes.
4. Analytics e edge consolidam visibilidade.

## Servicos em detalhe

## `analytics` - Analytics

- Stack: Python
- Plano: analytics plane
- Codigo: `service-api/service-python/analytics`
- Contexto de banco: `analytics/simulation`
- Contrato: `docs/contracts/http/analytics.openapi.yaml`
- Endpoints versionados: `9`
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/analytics/reports/adapter-catalog` - Read external adapter capability catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/integration-readiness` - Read external integration readiness
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/saas-control` - Read SaaS control posture by tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/contract-governance` - Read contract governance posture
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/hardening-review` - Read hardening review
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/core-operations` - Read core product operations
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/relationship-intelligence` - Read relationship intelligence
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/compliance-control` - Read fiscal and privacy compliance control
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/analytics/reports/go-live-control` - Read go-live rollout control
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `analytics` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `billing` - Billing

- Stack: .NET
- Plano: transaction/billing plane
- Codigo: `service-api/service-csharp/billing`
- Contexto de banco: `billing`
- Contrato: `docs/contracts/http/billing.openapi.yaml`
- Endpoints versionados: `9`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /health/details` - Return readiness details and gateway posture
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/billing/gateways` - List gateway capabilities and Pix posture
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/billing/gateways/{provider}` - Read one gateway capability
  - Parametros: `provider`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/billing/plans` - List billing plans including flat, hybrid and usage-based pricing
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/billing/plans` - Create billing plan
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/billing/subscriptions` - List subscriptions
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/billing/subscriptions` - Create subscription
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/billing/subscriptions/{publicId}/usage-pricing` - Project usage-based charge for one subscription
  - Parametros: `publicId`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/billing/invoices/{publicId}/attempts` - Create payment attempt with idempotency support
  - Parametros: `Idempotency-Key`, `publicId`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `billing` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `catalog` - Catalog

- Stack: Python
- Plano: transaction/catalog plane
- Codigo: `service-api/service-python/catalog`
- Contexto de banco: `catalog`
- Contrato: `docs/contracts/http/catalog.openapi.yaml`
- Endpoints versionados: `12`
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/catalog/capabilities` - Read catalog capability posture
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/consumers` - Read catalog consumer contracts across core domains
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/categories` - List categories by tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/catalog/categories` - Create one category
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/categories/page` - Cursor-based category listing
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/items` - List catalog items
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/catalog/items` - Create one catalog item
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/items/page` - Cursor-based item listing
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/catalog/items/bulk` - Bulk create catalog items with partial success
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/items/{publicId}` - Read one catalog item
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/catalog/items/{publicId}` - Update active state, price and attributes
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/catalog/items/{publicId}/versions` - Read catalog item version history
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `catalog` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `crm` - CRM

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/crm`
- Contexto de banco: `crm`
- Contrato: `docs/contracts/http/crm.openapi.yaml`
- Endpoints versionados: `5`
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/crm/enrichment/cnpj/capabilities` - Read CNPJ enrichment provider capabilities
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `crm` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/crm/enrichment/cnpj/lookup` - Lookup and enrich one CNPJ through provider contract
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `crm` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/crm/pipeline/config` - Read tenant pipeline configuration
  - Parametros: `tenantSlug`.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `crm` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/crm/pipeline/config` - Upsert tenant pipeline configuration
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `crm` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/crm/leads/intelligence/summary` - Read lead scoring and pipeline intelligence summary
  - Parametros: `tenantSlug`.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `crm` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `documents` - Documents

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/documents`
- Contexto de banco: `documents`
- Contrato: `docs/contracts/http/documents.openapi.yaml`
- Endpoints versionados: `10`
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /health/details` - Return runtime readiness and storage posture
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/documents/signing/capabilities` - List digital signature capabilities
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/documents/signing/capabilities/{provider}` - Read one signing capability
  - Parametros: `provider`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/documents/signing/requests` - Queue one digital signature request
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/documents/storage/capabilities` - List storage capability registry
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/documents/storage/capabilities/{provider}` - Read one storage capability
  - Parametros: `provider`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/documents/attachments` - List attachments
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/documents/attachments` - Create attachment metadata
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/documents/attachments/{publicId}/versions` - List attachment versions
  - Parametros: `publicId`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/documents/attachments/{publicId}/versions` - Append attachment version
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `documents` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `edge` - Edge

- Stack: Go
- Plano: public operations plane
- Codigo: `service-api/service-golang/edge`
- Contexto de banco: `none`
- Contrato: `docs/contracts/http/edge.openapi.yaml`
- Endpoints versionados: `8`
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/edge/ops/core-operations` - Read executive core product cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/relationship-overview` - Read executive relationship cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/compliance-overview` - Read executive compliance cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/go-live-overview` - Read executive go-live cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/integrations-overview` - Read executive integrations cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/saas-overview` - Read executive SaaS cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/contracts-overview` - Read executive contracts cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/edge/ops/hardening-overview` - Read executive hardening cockpit
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `edge` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `engagement` - Engagement

- Stack: TypeScript
- Plano: interaction/control plane
- Codigo: `service-api/service-typescript/engagement`
- Contexto de banco: `engagement`
- Contrato: `docs/contracts/http/engagement.openapi.yaml`
- Endpoints versionados: `9`
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /health/details` - Return readiness details for engagement runtime
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/engagement/providers` - List provider capabilities and fallback posture
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/engagement/providers/{provider}` - Read one provider capability
  - Parametros: `provider`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/engagement/providers/meta-ads/leads` - Ingest inbound lead from Meta Ads
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/engagement/providers/resend/events` - Register Resend callback event
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/engagement/providers/whatsapp-cloud/events` - Register WhatsApp callback event
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/engagement/providers/telegram-bot/events` - Register Telegram callback event
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/engagement/provider-events` - List provider events
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/engagement/provider-events/{publicId}` - Read one provider event
  - Parametros: `publicId`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `engagement` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `finance` - Finance

- Stack: .NET
- Plano: transaction/finance plane
- Codigo: `service-api/service-csharp/finance`
- Contexto de banco: `finance`
- Contrato: `docs/contracts/http/finance.openapi.yaml`
- Endpoints versionados: `5`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/finance/receivable-projections` - List receivable projections
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `finance` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/finance/receivable-projections/sync` - Sync projections from sales and rentals
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `finance` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/finance/commission-holds` - List commission holds
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `finance` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/finance/commission-holds/{publicId}/release` - Release one commission hold
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `finance` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/finance/activity` - List finance operational activity
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `finance` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `fiscal` - Fiscal

- Stack: Python
- Plano: compliance plane
- Codigo: `service-api/service-python/fiscal`
- Contexto de banco: `fiscal`
- Contrato: `docs/contracts/http/fiscal.openapi.yaml`
- Endpoints versionados: `25`
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/fiscal/capabilities` - Read fiscal capability registry
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/companies/{companyPublicId}/profile` - Read fiscal company profile
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/fiscal/companies/{companyPublicId}/profile` - Upsert fiscal company profile
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies` - List retention policies by company
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/companies/{companyPublicId}/retention-execution` - Read retention execution plan for one company
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute` - Execute retention and anonymization plan
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}` - Upsert retention policy for one data domain
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/documents` - List fiscal documents
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/documents` - Issue one fiscal document
  - Parametros: nenhum parametro declarado.
  - Respostas: `201`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/documents/{publicId}` - Read one fiscal document
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/documents/{publicId}/cancel` - Cancel one fiscal document
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/documents/{publicId}/correction-letter` - Register correction letter for one fiscal document
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/documents/{publicId}/invalidate` - Register invalidation for one fiscal document
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/documents/{publicId}/events` - List fiscal document audit events
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/privacy-requests` - List privacy requests
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/privacy-requests` - Create privacy request
  - Parametros: nenhum parametro declarado.
  - Respostas: `201`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/privacy-requests/{publicId}` - Read one privacy request
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/privacy-requests/{publicId}/export-package` - Build export package for one privacy request
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/privacy-requests/{publicId}/execute` - Execute one privacy request with audit trail
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/fiscal/privacy-requests/{publicId}/status` - Transition privacy request lifecycle status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/consents` - List consent ledger
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/fiscal/consents` - Create consent record
  - Parametros: nenhum parametro declarado.
  - Respostas: `201`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/fiscal/consents/{publicId}` - Transition consent status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`, `404`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/audit-events` - List fiscal audit events
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/fiscal/compliance/summary` - Read fiscal compliance summary
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `fiscal` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `identity` - Identity

- Stack: .NET
- Plano: transaction/security plane
- Codigo: `service-api/service-csharp/identity`
- Contexto de banco: `identity`
- Contrato: `docs/contracts/http/identity.openapi.yaml`
- Endpoints versionados: `6`
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/identity/tenants` - List tenants
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `identity` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/identity/tenants` - Create tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `identity` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/identity/tenants/{slug}/snapshot` - Read one tenant snapshot
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `identity` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/identity/sessions/login` - Authenticate identity session
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `identity` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/identity/sessions/refresh` - Refresh identity session
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `identity` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/identity/invitations` - Create invitation
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `identity` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `notification` - Notification

- Stack: Python
- Plano: administrative/notification plane
- Codigo: `service-api/service-python/notification`
- Contexto de banco: `notification`
- Contrato: `docs/contracts/http/notification.openapi.yaml`
- Endpoints versionados: `7`
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/notification/capabilities` - Read notification capability catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/notification/preferences/{userPublicId}` - Read one user notification preference
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/notification/preferences/{userPublicId}` - Upsert one user notification preference
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/notification/center` - List notification center items with cursor filters
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/notification/center` - Create one notification center item
  - Parametros: nenhum parametro declarado.
  - Respostas: `201`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/notification/center/{publicId}/status` - Transition notification status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/notification/summary` - Read notification summary
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `notification` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `platform-control` - Platform Control

- Stack: Python
- Plano: platform control plane
- Codigo: `service-api/service-python/platform-control`
- Contexto de banco: `platform-control`
- Contrato: `docs/contracts/http/platform-control.openapi.yaml`
- Endpoints versionados: `40`
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/platform-control/capabilities/catalog` - List platform capability catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/providers/catalog` - List provider capability catalog and environment posture
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements` - List tenant entitlements with cursor pagination
  - Parametros: `tenantSlug`.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/feature-flags` - List tenant feature flags with capability metadata
  - Parametros: `tenantSlug`.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}` - Upsert one entitlement
  - Parametros: `capabilityKey`, `tenantSlug`.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}` - Upsert one feature flag using entitlement governance
  - Parametros: `capabilityKey`, `tenantSlug`.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk` - Bulk upsert entitlements with partial success
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults` - List provider defaults selected for one tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}` - Upsert provider default for one tenant capability
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/quotas` - List quotas by tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}` - Upsert one quota
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk` - Bulk upsert quotas with partial success
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/blocks` - List tenant blocks
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}` - Upsert tenant block
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/metering` - Read metering snapshots and summary with cursor pagination
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots` - Create one usage snapshot
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/usage-summary` - Read quota and metering utilization summary
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness` - Read tenant lifecycle readiness and provider posture
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs` - List onboarding and offboarding jobs with cursor pagination
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}` - Read one lifecycle job with audit events
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview` - Preview onboarding plan, provider defaults and lifecycle readiness
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding` - Queue onboarding job with Idempotency-Key and 202 Accepted
  - Parametros: nenhum parametro declarado.
  - Respostas: `202`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview` - Preview offboarding plan, retention posture and lifecycle readiness
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding` - Queue offboarding job with Idempotency-Key and 202 Accepted
  - Parametros: nenhum parametro declarado.
  - Respostas: `202`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start` - Transition lifecycle job to running
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete` - Transition lifecycle job to completed
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail` - Transition lifecycle job to failed
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel` - Transition lifecycle job to cancelled
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` - Read go-live rollout readiness by tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` - Read tenant go-live adoption baseline and gap
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` - List go-live bottlenecks and operational blockers
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` - Read rollout and rollback playbook for one tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` - List recommended go-live adjustments
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply` - Apply one go-live operational adjustment
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - List go-live rollouts
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` - Create one go-live rollout
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}` - Read one go-live rollout with events
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` - Transition go-live rollout to running
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` - Transition go-live rollout to completed
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` - Roll back one go-live rollout
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `platform-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `rentals` - Rentals

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/rentals`
- Contexto de banco: `rentals`
- Contrato: `docs/contracts/http/rentals.openapi.yaml`
- Endpoints versionados: `4`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/rentals/contracts` - List rental contracts
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `rentals` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/rentals/contracts` - Create rental contract
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `rentals` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/rentals/contracts/{publicId}/charges` - List contract charges
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `rentals` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` - Update charge status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `rentals` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `sales` - Sales

- Stack: Go
- Plano: transaction plane
- Codigo: `service-api/service-golang/sales`
- Contexto de banco: `sales`
- Contrato: `docs/contracts/http/sales.openapi.yaml`
- Endpoints versionados: `6`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/sales/opportunities` - List opportunities
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `sales` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/sales/opportunities` - Create opportunity
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `sales` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/sales/proposals` - List proposals
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `sales` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/sales/proposals` - Create proposal
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `sales` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/sales/sales` - List sales
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `sales` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/sales/invoices` - List commercial invoices
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `sales` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `simulation` - Simulation

- Stack: Python
- Plano: simulation plane
- Codigo: `service-api/service-python/simulation`
- Contexto de banco: `simulation`
- Contrato: `docs/contracts/http/simulation.openapi.yaml`
- Endpoints versionados: `3`
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/simulation/scenarios` - List scenarios
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `simulation` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/simulation/scenarios` - Create scenario run
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `simulation` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/simulation/benchmarks/load` - Execute one load benchmark run
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `simulation` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `supplier` - Supplier

- Stack: Python
- Plano: administrative/procurement plane
- Codigo: `service-api/service-python/supplier`
- Contexto de banco: `supplier`
- Contrato: `docs/contracts/http/supplier.openapi.yaml`
- Endpoints versionados: `8`
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/supplier/capabilities` - Read supplier capability catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/supplier/categories` - List supplier categories
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/supplier/categories/{categoryKey}` - Upsert one supplier category
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/supplier/suppliers` - List suppliers by tenant and status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/supplier/suppliers` - Create one supplier
  - Parametros: nenhum parametro declarado.
  - Respostas: `201`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/supplier/suppliers/summary` - Read supplier summary
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/supplier/suppliers/{publicId}` - Read one supplier
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/supplier/suppliers/{publicId}` - Update one supplier
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `supplier` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `support` - Support

- Stack: Python
- Plano: administrative plane
- Codigo: `service-api/service-python/support`
- Contexto de banco: `support`
- Contrato: `docs/contracts/http/support.openapi.yaml`
- Endpoints versionados: `9`
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/support/capabilities` - Read support capability catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/support/queues` - List support queues by tenant
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PUT /api/support/queues/{queueKey}` - Upsert one support queue
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/support/cases` - List support cases with cursor filters
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/support/cases` - Create one support case
  - Parametros: nenhum parametro declarado.
  - Respostas: `201`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/support/cases/summary` - Read support case summary
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/support/cases/{publicId}` - Read one support case
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/support/cases/{publicId}/status` - Transition support case status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/support/cases/{publicId}/comments` - Append comment to support case
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `support` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `webhook-hub` - Webhook Hub

- Stack: Rust
- Plano: integration plane
- Codigo: `service-api/service-rust/webhook-hub`
- Contexto de banco: `webhook-hub`
- Contrato: `docs/contracts/http/webhook-hub.openapi.yaml`
- Endpoints versionados: `13`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /health/details` - Return readiness details for webhook runtime
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/webhook-hub/capabilities` - Read outbound webhook capability posture
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/webhook-hub/outbound-endpoints` - List tenant outbound endpoints
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/webhook-hub/outbound-endpoints` - Register one tenant outbound endpoint
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/webhook-hub/outbound-endpoints/{publicId}` - Read one tenant outbound endpoint
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - List outbound delivery log for one endpoint
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - Register one outbound delivery attempt
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter` - Move one outbound delivery to dead letter
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/webhook-hub/events` - List inbound webhook events
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/webhook-hub/events` - Register inbound webhook event
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/webhook-hub/events/summary` - Aggregate inbound webhook state
  - Parametros: nenhum parametro declarado.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/webhook-hub/events/{publicId}/dead-letter` - Move event to dead letter queue
  - Parametros: `publicId`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/webhook-hub/events/{publicId}/requeue` - Requeue dead-letter event
  - Parametros: `publicId`.
  - Respostas: nenhuma resposta declarada.
  - Impacto arquitetural: manter ownership em `webhook-hub` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `workflow-control` - Workflow Control

- Stack: TypeScript
- Plano: control plane
- Codigo: `service-api/service-typescript/workflow-control`
- Contexto de banco: `workflow-control`
- Contrato: `docs/contracts/http/workflow-control.openapi.yaml`
- Endpoints versionados: `7`
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/workflow-control/definitions` - List workflow definitions
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/workflow-control/definitions` - Create workflow definition
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/workflow-control/definitions/{key}` - Read one workflow definition
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/workflow-control/definitions/{key}` - Update one workflow definition
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `PATCH /api/workflow-control/definitions/{key}/status` - Update workflow definition status
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/workflow-control/capabilities/triggers` - List workflow trigger catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/workflow-control/capabilities/actions` - List workflow action catalog
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-control` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## `workflow-runtime` - Workflow Runtime

- Stack: Elixir
- Plano: runtime plane
- Codigo: `service-api/service-elixir/workflow-runtime`
- Contexto de banco: `workflow-runtime`
- Contrato: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Endpoints versionados: `6`
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.

### Fronteira arquitetural

- Este servico deve manter regra de negocio no dominio ou camada de aplicacao.
- Bootstrap HTTP deve permanecer fino.
- A camada de infraestrutura deve isolar banco, provider externo, cache, fila e clients HTTP.
- O contrato publico precisa evoluir junto com testes e smoke.
- Qualquer dependencia cross-context deve ser explicita por contrato, evento, read model ou adapter.

### Dados e ownership

- O schema transacional pertence ao contexto declarado acima quando existir.
- Leitura analitica pode agregar dados, mas nao assume ownership de escrita.
- Mutacoes relevantes devem preservar tenant, ator, correlation id e trilha operacional.
- Seeds existem para bootstrap e smoke, nao para esconder regra incompleta.

### Rotas e comportamento arquitetural

- `GET /api/workflow-runtime/executions` - List workflow executions
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-runtime` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/workflow-runtime/executions` - Create workflow execution
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-runtime` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/workflow-runtime/executions/{publicId}` - Read one workflow execution
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-runtime` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/workflow-runtime/executions/{publicId}/actions` - List execution action snapshots
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-runtime` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `POST /api/workflow-runtime/executions/{publicId}/advance` - Advance one workflow execution
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-runtime` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.
- `GET /api/workflow-runtime/capabilities` - List runtime capabilities
  - Parametros: nenhum parametro declarado.
  - Respostas: `200`.
  - Impacto arquitetural: manter ownership em `workflow-runtime` e nao vazar regra para consumidor.
  - Observabilidade: correlation id, tenant quando aplicavel, recurso publico e status final.

### Riscos de evolucao

- acoplamento por tabela em vez de contrato;
- fallback de tenant usado fora de bootstrap;
- provider externo sem postura operacional explicita;
- endpoint novo sem contrato e sem smoke quando for fluxo critico;

## Checklist de arquitetura

- O dono do dominio esta claro?
- O schema de banco esta no contexto correto?
- A rota nova esta em `docs/contracts/http/`?
- O evento novo esta em `docs/contracts/events/`?
- O smoke cobre o fluxo quando ha composicao entre servicos?
- O health detail reflete dependencias reais?
- A documentacao foi atualizada em `docs/`?
