# OPERACOES

## Objetivo

Documentar comandos, runtime local, validacao, observabilidade, banco, backup, smoke e runbooks operacionais. O projeto usa dois comandos oficiais: `./scripts/build.sh` para operacao local e `./scripts/test.sh` para validacao.

## Comandos oficiais

### Runtime

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh down
./scripts/build.sh restart
./scripts/build.sh ps
./scripts/build.sh logs edge
```

### Banco

```bash
./scripts/build.sh migrate all
./scripts/build.sh migrate identity
./scripts/build.sh seed all
./scripts/build.sh summary crm
./scripts/build.sh backup /tmp/erp-local-backup.sql
./scripts/build.sh restore /tmp/erp-local-backup.sql
./scripts/build.sh psql
```

### Testes

```bash
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh contract
./scripts/test.sh platform
./scripts/test.sh smoke
./scripts/test.sh performance
./scripts/test.sh backup-restore
./scripts/test.sh hardening
```

## Enderecos locais

- Edge: `http://localhost:${EDGE_HTTP_PORT}`
- Identity: `http://localhost:${IDENTITY_HTTP_PORT}`
- CRM: `http://localhost:${CRM_HTTP_PORT}`
- Sales: `http://localhost:${SALES_HTTP_PORT}`
- Workflow Control: `http://localhost:${WORKFLOW_CONTROL_HTTP_PORT}`
- Workflow Runtime: `http://localhost:${WORKFLOW_RUNTIME_HTTP_PORT}`
- Analytics: `http://localhost:${ANALYTICS_HTTP_PORT}`
- Webhook Hub: `http://localhost:${WEBHOOK_HUB_HTTP_PORT}`
- PostgreSQL: `localhost:${POSTGRES_PORT}`
- Redis: `localhost:${REDIS_PORT}`
- Kafka: `localhost:${KAFKA_PORT}`
- Keycloak: `http://localhost:${KEYCLOAK_PORT}`
- OpenFGA: `http://localhost:${OPENFGA_HTTP_PORT}`
- Prometheus: `http://localhost:${PROMETHEUS_PORT}`
- Grafana: `http://localhost:${GRAFANA_PORT}`

## Runbook por servico

## `analytics`

- Stack: Python
- Codigo: `service-api/service-python/analytics`
- Banco: `analytics/simulation`
- Contrato: `docs/contracts/http/analytics.openapi.yaml`
- Responsabilidade: relatorios operacionais, governanca, confiabilidade, hardening, custos e leituras executivas.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs analytics
./scripts/build.sh ps
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/analytics/reports/adapter-catalog` - Read external adapter capability catalog
- `GET /api/analytics/reports/integration-readiness` - Read external integration readiness
- `GET /api/analytics/reports/saas-control` - Read SaaS control posture by tenant
- `GET /api/analytics/reports/contract-governance` - Read contract governance posture
- `GET /api/analytics/reports/hardening-review` - Read hardening review
- `GET /api/analytics/reports/core-operations` - Read core product operations
- `GET /api/analytics/reports/relationship-intelligence` - Read relationship intelligence
- `GET /api/analytics/reports/compliance-control` - Read fiscal and privacy compliance control
- `GET /api/analytics/reports/go-live-control` - Read go-live rollout control

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `billing`

- Stack: .NET
- Codigo: `service-api/service-csharp/billing`
- Banco: `billing`
- Contrato: `docs/contracts/http/billing.openapi.yaml`
- Responsabilidade: planos, assinaturas, invoices recorrentes, tentativas de cobranca e recovery.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs billing
./scripts/build.sh summary billing
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /health/details` - Return readiness details and gateway posture
- `GET /api/billing/gateways` - List gateway capabilities and Pix posture
- `GET /api/billing/gateways/{provider}` - Read one gateway capability
- `GET /api/billing/plans` - List billing plans including flat, hybrid and usage-based pricing
- `POST /api/billing/plans` - Create billing plan
- `GET /api/billing/subscriptions` - List subscriptions
- `POST /api/billing/subscriptions` - Create subscription
- `GET /api/billing/subscriptions/{publicId}/usage-pricing` - Project usage-based charge for one subscription
- `POST /api/billing/invoices/{publicId}/attempts` - Create payment attempt with idempotency support

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `catalog`

- Stack: Python
- Codigo: `service-api/service-python/catalog`
- Banco: `catalog`
- Contrato: `docs/contracts/http/catalog.openapi.yaml`
- Responsabilidade: categorias, itens, versoes de item, bulk e contratos de consumo.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs catalog
./scripts/build.sh summary catalog
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/catalog/capabilities` - Read catalog capability posture
- `GET /api/catalog/consumers` - Read catalog consumer contracts across core domains
- `GET /api/catalog/categories` - List categories by tenant
- `POST /api/catalog/categories` - Create one category
- `GET /api/catalog/categories/page` - Cursor-based category listing
- `GET /api/catalog/items` - List catalog items
- `POST /api/catalog/items` - Create one catalog item
- `GET /api/catalog/items/page` - Cursor-based item listing
- `POST /api/catalog/items/bulk` - Bulk create catalog items with partial success
- `GET /api/catalog/items/{publicId}` - Read one catalog item
- `PATCH /api/catalog/items/{publicId}` - Update active state, price and attributes
- `GET /api/catalog/items/{publicId}/versions` - Read catalog item version history

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `crm`

- Stack: Go
- Codigo: `service-api/service-golang/crm`
- Banco: `crm`
- Contrato: `docs/contracts/http/crm.openapi.yaml`
- Responsabilidade: leads, customers, ownership, pipeline, notas, historico, anexos e enriquecimento.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs crm
./scripts/build.sh summary crm
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/crm/enrichment/cnpj/capabilities` - Read CNPJ enrichment provider capabilities
- `POST /api/crm/enrichment/cnpj/lookup` - Lookup and enrich one CNPJ through provider contract
- `GET /api/crm/pipeline/config` - Read tenant pipeline configuration
- `PUT /api/crm/pipeline/config` - Upsert tenant pipeline configuration
- `GET /api/crm/leads/intelligence/summary` - Read lead scoring and pipeline intelligence summary

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `documents`

- Stack: Go
- Codigo: `service-api/service-golang/documents`
- Banco: `documents`
- Contrato: `docs/contracts/http/documents.openapi.yaml`
- Responsabilidade: anexos, upload, storage posture, assinatura, versoes, archive e access links.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs documents
./scripts/build.sh summary documents
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /health/details` - Return runtime readiness and storage posture
- `GET /api/documents/signing/capabilities` - List digital signature capabilities
- `GET /api/documents/signing/capabilities/{provider}` - Read one signing capability
- `POST /api/documents/signing/requests` - Queue one digital signature request
- `GET /api/documents/storage/capabilities` - List storage capability registry
- `GET /api/documents/storage/capabilities/{provider}` - Read one storage capability
- `GET /api/documents/attachments` - List attachments
- `POST /api/documents/attachments` - Create attachment metadata
- `GET /api/documents/attachments/{publicId}/versions` - List attachment versions
- `POST /api/documents/attachments/{publicId}/versions` - Append attachment version

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `edge`

- Stack: Go
- Codigo: `service-api/service-golang/edge`
- Banco: `none`
- Contrato: `docs/contracts/http/edge.openapi.yaml`
- Responsabilidade: entrada publica, agregacao cross-service e cockpits operacionais.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs edge
./scripts/build.sh ps
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/edge/ops/core-operations` - Read executive core product cockpit
- `GET /api/edge/ops/relationship-overview` - Read executive relationship cockpit
- `GET /api/edge/ops/compliance-overview` - Read executive compliance cockpit
- `GET /api/edge/ops/go-live-overview` - Read executive go-live cockpit
- `GET /api/edge/ops/integrations-overview` - Read executive integrations cockpit
- `GET /api/edge/ops/saas-overview` - Read executive SaaS cockpit
- `GET /api/edge/ops/contracts-overview` - Read executive contracts cockpit
- `GET /api/edge/ops/hardening-overview` - Read executive hardening cockpit

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `engagement`

- Stack: TypeScript
- Codigo: `service-api/service-typescript/engagement`
- Banco: `engagement`
- Contrato: `docs/contracts/http/engagement.openapi.yaml`
- Responsabilidade: campanhas, templates, touchpoints, conversas, delivery, providers e callbacks.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs engagement
./scripts/build.sh summary engagement
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /health/details` - Return readiness details for engagement runtime
- `GET /api/engagement/providers` - List provider capabilities and fallback posture
- `GET /api/engagement/providers/{provider}` - Read one provider capability
- `POST /api/engagement/providers/meta-ads/leads` - Ingest inbound lead from Meta Ads
- `POST /api/engagement/providers/resend/events` - Register Resend callback event
- `POST /api/engagement/providers/whatsapp-cloud/events` - Register WhatsApp callback event
- `POST /api/engagement/providers/telegram-bot/events` - Register Telegram callback event
- `GET /api/engagement/provider-events` - List provider events
- `GET /api/engagement/provider-events/{publicId}` - Read one provider event

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `finance`

- Stack: .NET
- Codigo: `service-api/service-csharp/finance`
- Banco: `finance`
- Contrato: `docs/contracts/http/finance.openapi.yaml`
- Responsabilidade: recebiveis, payables, caixa, custos, comissoes e fechamento.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs finance
./scripts/build.sh summary finance
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/finance/receivable-projections` - List receivable projections
- `POST /api/finance/receivable-projections/sync` - Sync projections from sales and rentals
- `GET /api/finance/commission-holds` - List commission holds
- `POST /api/finance/commission-holds/{publicId}/release` - Release one commission hold
- `GET /api/finance/activity` - List finance operational activity

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `fiscal`

- Stack: Python
- Codigo: `service-api/service-python/fiscal`
- Banco: `fiscal`
- Contrato: `docs/contracts/http/fiscal.openapi.yaml`
- Responsabilidade: perfil fiscal, retencao, documentos fiscais, privacidade, consentimentos e auditoria.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs fiscal
./scripts/build.sh summary fiscal
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/fiscal/capabilities` - Read fiscal capability registry
- `GET /api/fiscal/companies/{companyPublicId}/profile` - Read fiscal company profile
- `PUT /api/fiscal/companies/{companyPublicId}/profile` - Upsert fiscal company profile
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies` - List retention policies by company
- `GET /api/fiscal/companies/{companyPublicId}/retention-execution` - Read retention execution plan for one company
- `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute` - Execute retention and anonymization plan
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}` - Upsert retention policy for one data domain
- `GET /api/fiscal/documents` - List fiscal documents
- `POST /api/fiscal/documents` - Issue one fiscal document
- `GET /api/fiscal/documents/{publicId}` - Read one fiscal document
- `POST /api/fiscal/documents/{publicId}/cancel` - Cancel one fiscal document
- `POST /api/fiscal/documents/{publicId}/correction-letter` - Register correction letter for one fiscal document

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `identity`

- Stack: .NET
- Codigo: `service-api/service-csharp/identity`
- Banco: `identity`
- Contrato: `docs/contracts/http/identity.openapi.yaml`
- Responsabilidade: tenants, empresas, usuarios, times, roles, sessoes, convites, MFA e auditoria.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs identity
./scripts/build.sh summary identity
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/identity/tenants` - List tenants
- `POST /api/identity/tenants` - Create tenant
- `GET /api/identity/tenants/{slug}/snapshot` - Read one tenant snapshot
- `POST /api/identity/sessions/login` - Authenticate identity session
- `POST /api/identity/sessions/refresh` - Refresh identity session
- `POST /api/identity/invitations` - Create invitation

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `notification`

- Stack: Python
- Codigo: `service-api/service-python/notification`
- Banco: `notification`
- Contrato: `docs/contracts/http/notification.openapi.yaml`
- Responsabilidade: preferencias, centro interno de alertas, severidade e lifecycle de notificacoes.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs notification
./scripts/build.sh summary notification
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/notification/capabilities` - Read notification capability catalog
- `GET /api/notification/preferences/{userPublicId}` - Read one user notification preference
- `PUT /api/notification/preferences/{userPublicId}` - Upsert one user notification preference
- `GET /api/notification/center` - List notification center items with cursor filters
- `POST /api/notification/center` - Create one notification center item
- `PATCH /api/notification/center/{publicId}/status` - Transition notification status
- `GET /api/notification/summary` - Read notification summary

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `platform-control`

- Stack: Python
- Codigo: `service-api/service-python/platform-control`
- Banco: `platform-control`
- Contrato: `docs/contracts/http/platform-control.openapi.yaml`
- Responsabilidade: capabilities, providers, entitlements, feature flags, quotas, metering, lifecycle e go-live.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs platform-control
./scripts/build.sh summary platform-control
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/platform-control/capabilities/catalog` - List platform capability catalog
- `GET /api/platform-control/providers/catalog` - List provider capability catalog and environment posture
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements` - List tenant entitlements with cursor pagination
- `GET /api/platform-control/tenants/{tenantSlug}/feature-flags` - List tenant feature flags with capability metadata
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}` - Upsert one entitlement
- `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}` - Upsert one feature flag using entitlement governance
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk` - Bulk upsert entitlements with partial success
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults` - List provider defaults selected for one tenant
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}` - Upsert provider default for one tenant capability
- `GET /api/platform-control/tenants/{tenantSlug}/quotas` - List quotas by tenant
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}` - Upsert one quota
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk` - Bulk upsert quotas with partial success

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `rentals`

- Stack: Go
- Codigo: `service-api/service-golang/rentals`
- Banco: `rentals`
- Contrato: `docs/contracts/http/rentals.openapi.yaml`
- Responsabilidade: contratos recorrentes, reajustes, terminacoes, cobrancas e anexos contratuais.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs rentals
./scripts/build.sh summary rentals
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/rentals/contracts` - List rental contracts
- `POST /api/rentals/contracts` - Create rental contract
- `GET /api/rentals/contracts/{publicId}/charges` - List contract charges
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` - Update charge status

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `sales`

- Stack: Go
- Codigo: `service-api/service-golang/sales`
- Banco: `sales`
- Contrato: `docs/contracts/http/sales.openapi.yaml`
- Responsabilidade: oportunidades, propostas, vendas, invoices, comissoes, renegociacoes e pendencias.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs sales
./scripts/build.sh summary sales
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/sales/opportunities` - List opportunities
- `POST /api/sales/opportunities` - Create opportunity
- `GET /api/sales/proposals` - List proposals
- `POST /api/sales/proposals` - Create proposal
- `GET /api/sales/sales` - List sales
- `GET /api/sales/invoices` - List commercial invoices

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `simulation`

- Stack: Python
- Codigo: `service-api/service-python/simulation`
- Banco: `simulation`
- Contrato: `docs/contracts/http/simulation.openapi.yaml`
- Responsabilidade: cenarios operacionais, benchmark de carga e modelagem de capacidade.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs simulation
./scripts/build.sh summary simulation
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/simulation/scenarios` - List scenarios
- `POST /api/simulation/scenarios` - Create scenario run
- `POST /api/simulation/benchmarks/load` - Execute one load benchmark run

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `supplier`

- Stack: Python
- Codigo: `service-api/service-python/supplier`
- Banco: `supplier`
- Contrato: `docs/contracts/http/supplier.openapi.yaml`
- Responsabilidade: categorias de fornecedor, diretorio de fornecedores e procurement ownership.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs supplier
./scripts/build.sh summary supplier
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/supplier/capabilities` - Read supplier capability catalog
- `GET /api/supplier/categories` - List supplier categories
- `PUT /api/supplier/categories/{categoryKey}` - Upsert one supplier category
- `GET /api/supplier/suppliers` - List suppliers by tenant and status
- `POST /api/supplier/suppliers` - Create one supplier
- `GET /api/supplier/suppliers/summary` - Read supplier summary
- `GET /api/supplier/suppliers/{publicId}` - Read one supplier
- `PATCH /api/supplier/suppliers/{publicId}` - Update one supplier

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `support`

- Stack: Python
- Codigo: `service-api/service-python/support`
- Banco: `support`
- Contrato: `docs/contracts/http/support.openapi.yaml`
- Responsabilidade: filas, casos, SLA, comentarios, bulk e resumo operacional de atendimento.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs support
./scripts/build.sh summary support
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/support/capabilities` - Read support capability catalog
- `GET /api/support/queues` - List support queues by tenant
- `PUT /api/support/queues/{queueKey}` - Upsert one support queue
- `GET /api/support/cases` - List support cases with cursor filters
- `POST /api/support/cases` - Create one support case
- `GET /api/support/cases/summary` - Read support case summary
- `GET /api/support/cases/{publicId}` - Read one support case
- `PATCH /api/support/cases/{publicId}/status` - Transition support case status
- `POST /api/support/cases/{publicId}/comments` - Append comment to support case

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `webhook-hub`

- Stack: Rust
- Codigo: `service-api/service-rust/webhook-hub`
- Banco: `webhook-hub`
- Contrato: `docs/contracts/http/webhook-hub.openapi.yaml`
- Responsabilidade: intake de webhooks, idempotencia, transicoes, DLQ, outbound endpoints e deliveries.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs webhook-hub
./scripts/build.sh summary webhook-hub
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /health/details` - Return readiness details for webhook runtime
- `GET /api/webhook-hub/capabilities` - Read outbound webhook capability posture
- `GET /api/webhook-hub/outbound-endpoints` - List tenant outbound endpoints
- `POST /api/webhook-hub/outbound-endpoints` - Register one tenant outbound endpoint
- `GET /api/webhook-hub/outbound-endpoints/{publicId}` - Read one tenant outbound endpoint
- `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - List outbound delivery log for one endpoint
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries` - Register one outbound delivery attempt
- `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter` - Move one outbound delivery to dead letter
- `GET /api/webhook-hub/events` - List inbound webhook events
- `POST /api/webhook-hub/events` - Register inbound webhook event
- `GET /api/webhook-hub/events/summary` - Aggregate inbound webhook state
- `POST /api/webhook-hub/events/{publicId}/dead-letter` - Move event to dead letter queue

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `workflow-control`

- Stack: TypeScript
- Codigo: `service-api/service-typescript/workflow-control`
- Banco: `workflow-control`
- Contrato: `docs/contracts/http/workflow-control.openapi.yaml`
- Responsabilidade: definicoes, versoes publicadas, catalogos de trigger/action, runs e eventos.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs workflow-control
./scripts/build.sh summary workflow-control
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/workflow-control/definitions` - List workflow definitions
- `POST /api/workflow-control/definitions` - Create workflow definition
- `GET /api/workflow-control/definitions/{key}` - Read one workflow definition
- `PATCH /api/workflow-control/definitions/{key}` - Update one workflow definition
- `PATCH /api/workflow-control/definitions/{key}/status` - Update workflow definition status
- `GET /api/workflow-control/capabilities/triggers` - List workflow trigger catalog
- `GET /api/workflow-control/capabilities/actions` - List workflow action catalog

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## `workflow-runtime`

- Stack: Elixir
- Codigo: `service-api/service-elixir/workflow-runtime`
- Banco: `workflow-runtime`
- Contrato: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Responsabilidade: execucoes duraveis, timeline, actions, transicoes, retries, waits e compensacoes.

### Diagnostico rapido

```bash
./scripts/build.sh ps
./scripts/build.sh logs workflow-runtime
./scripts/build.sh summary workflow-runtime
./scripts/test.sh contract
```

### Health esperado

- `/health/live` precisa responder quando o processo HTTP esta vivo.
- `/health/ready` precisa refletir dependencias reais.
- `/health/details` precisa ajudar diagnostico sem expor segredo.

### Rotas criticas para smoke ou diagnostico

- `GET /api/workflow-runtime/executions` - List workflow executions
- `POST /api/workflow-runtime/executions` - Create workflow execution
- `GET /api/workflow-runtime/executions/{publicId}` - Read one workflow execution
- `GET /api/workflow-runtime/executions/{publicId}/actions` - List execution action snapshots
- `POST /api/workflow-runtime/executions/{publicId}/advance` - Advance one workflow execution
- `GET /api/workflow-runtime/capabilities` - List runtime capabilities

### Falhas comuns

- porta ocupada no host: o build remapeia quando a variavel permite;
- banco sem migration: rodar `./scripts/build.sh migrate all`;
- dado bootstrap ausente: rodar `./scripts/build.sh seed all`;
- contrato divergente: rodar `./scripts/test.sh contract`;
- fluxo cross-service quebrado: rodar `./scripts/test.sh smoke`;

## Observabilidade minima

- correlation id em requisicoes;
- tenant em rotas tenant-aware;
- status HTTP final;
- provider e external id em webhooks;
- job/run/execution id em lifecycle e automacao;
- erro normalizado em falha de provider;
- metricas de latencia, erro, throughput e filas.

## Rotas Operacionais Por Servico

Esta secao serve como checklist rapido durante incidente, revisao de smoke ou diagnostico de contrato. Cada rota abaixo deve ser considerada parte da superficie operacional versionada enquanto existir no OpenAPI atual.

## Rotas de `analytics`

- Contrato: `docs/contracts/http/analytics.openapi.yaml`
- Total: `9`

### `GET /api/analytics/reports/adapter-catalog`

- Summary: Read external adapter capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/integration-readiness`

- Summary: Read external integration readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/saas-control`

- Summary: Read SaaS control posture by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/contract-governance`

- Summary: Read contract governance posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/hardening-review`

- Summary: Read hardening review.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/core-operations`

- Summary: Read core product operations.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/relationship-intelligence`

- Summary: Read relationship intelligence.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/compliance-control`

- Summary: Read fiscal and privacy compliance control.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/analytics/reports/go-live-control`

- Summary: Read go-live rollout control.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `billing`

- Contrato: `docs/contracts/http/billing.openapi.yaml`
- Total: `9`

### `GET /health/details`

- Summary: Return readiness details and gateway posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/billing/gateways`

- Summary: List gateway capabilities and Pix posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/billing/gateways/{provider}`

- Summary: Read one gateway capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/billing/plans`

- Summary: List billing plans including flat, hybrid and usage-based pricing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/billing/plans`

- Summary: Create billing plan.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/billing/subscriptions`

- Summary: List subscriptions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/billing/subscriptions`

- Summary: Create subscription.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/billing/subscriptions/{publicId}/usage-pricing`

- Summary: Project usage-based charge for one subscription.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/billing/invoices/{publicId}/attempts`

- Summary: Create payment attempt with idempotency support.
- Parametros: `Idempotency-Key`, `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `catalog`

- Contrato: `docs/contracts/http/catalog.openapi.yaml`
- Total: `12`

### `GET /api/catalog/capabilities`

- Summary: Read catalog capability posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/consumers`

- Summary: Read catalog consumer contracts across core domains.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/categories`

- Summary: List categories by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/catalog/categories`

- Summary: Create one category.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/categories/page`

- Summary: Cursor-based category listing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/items`

- Summary: List catalog items.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/catalog/items`

- Summary: Create one catalog item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/items/page`

- Summary: Cursor-based item listing.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/catalog/items/bulk`

- Summary: Bulk create catalog items with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/items/{publicId}`

- Summary: Read one catalog item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/catalog/items/{publicId}`

- Summary: Update active state, price and attributes.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/catalog/items/{publicId}/versions`

- Summary: Read catalog item version history.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `crm`

- Contrato: `docs/contracts/http/crm.openapi.yaml`
- Total: `5`

### `GET /api/crm/enrichment/cnpj/capabilities`

- Summary: Read CNPJ enrichment provider capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/crm/enrichment/cnpj/lookup`

- Summary: Lookup and enrich one CNPJ through provider contract.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/crm/pipeline/config`

- Summary: Read tenant pipeline configuration.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/crm/pipeline/config`

- Summary: Upsert tenant pipeline configuration.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/crm/leads/intelligence/summary`

- Summary: Read lead scoring and pipeline intelligence summary.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `documents`

- Contrato: `docs/contracts/http/documents.openapi.yaml`
- Total: `10`

### `GET /health/details`

- Summary: Return runtime readiness and storage posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/documents/signing/capabilities`

- Summary: List digital signature capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/documents/signing/capabilities/{provider}`

- Summary: Read one signing capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/documents/signing/requests`

- Summary: Queue one digital signature request.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/documents/storage/capabilities`

- Summary: List storage capability registry.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/documents/storage/capabilities/{provider}`

- Summary: Read one storage capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/documents/attachments`

- Summary: List attachments.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/documents/attachments`

- Summary: Create attachment metadata.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/documents/attachments/{publicId}/versions`

- Summary: List attachment versions.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/documents/attachments/{publicId}/versions`

- Summary: Append attachment version.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `edge`

- Contrato: `docs/contracts/http/edge.openapi.yaml`
- Total: `8`

### `GET /api/edge/ops/core-operations`

- Summary: Read executive core product cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/relationship-overview`

- Summary: Read executive relationship cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/compliance-overview`

- Summary: Read executive compliance cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/go-live-overview`

- Summary: Read executive go-live cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/integrations-overview`

- Summary: Read executive integrations cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/saas-overview`

- Summary: Read executive SaaS cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/contracts-overview`

- Summary: Read executive contracts cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/edge/ops/hardening-overview`

- Summary: Read executive hardening cockpit.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `engagement`

- Contrato: `docs/contracts/http/engagement.openapi.yaml`
- Total: `9`

### `GET /health/details`

- Summary: Return readiness details for engagement runtime.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/engagement/providers`

- Summary: List provider capabilities and fallback posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/engagement/providers/{provider}`

- Summary: Read one provider capability.
- Parametros: `provider`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/engagement/providers/meta-ads/leads`

- Summary: Ingest inbound lead from Meta Ads.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/engagement/providers/resend/events`

- Summary: Register Resend callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/engagement/providers/whatsapp-cloud/events`

- Summary: Register WhatsApp callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/engagement/providers/telegram-bot/events`

- Summary: Register Telegram callback event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/engagement/provider-events`

- Summary: List provider events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/engagement/provider-events/{publicId}`

- Summary: Read one provider event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `finance`

- Contrato: `docs/contracts/http/finance.openapi.yaml`
- Total: `5`

### `GET /api/finance/receivable-projections`

- Summary: List receivable projections.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/finance/receivable-projections/sync`

- Summary: Sync projections from sales and rentals.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/finance/commission-holds`

- Summary: List commission holds.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/finance/commission-holds/{publicId}/release`

- Summary: Release one commission hold.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/finance/activity`

- Summary: List finance operational activity.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `fiscal`

- Contrato: `docs/contracts/http/fiscal.openapi.yaml`
- Total: `25`

### `GET /api/fiscal/capabilities`

- Summary: Read fiscal capability registry.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Read fiscal company profile.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/fiscal/companies/{companyPublicId}/profile`

- Summary: Upsert fiscal company profile.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/companies/{companyPublicId}/retention-policies`

- Summary: List retention policies by company.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/companies/{companyPublicId}/retention-execution`

- Summary: Read retention execution plan for one company.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/companies/{companyPublicId}/retention-execution/execute`

- Summary: Execute retention and anonymization plan.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`

- Summary: Upsert retention policy for one data domain.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/documents`

- Summary: List fiscal documents.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/documents`

- Summary: Issue one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/documents/{publicId}`

- Summary: Read one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/documents/{publicId}/cancel`

- Summary: Cancel one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/documents/{publicId}/correction-letter`

- Summary: Register correction letter for one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/documents/{publicId}/invalidate`

- Summary: Register invalidation for one fiscal document.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/documents/{publicId}/events`

- Summary: List fiscal document audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/privacy-requests`

- Summary: List privacy requests.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/privacy-requests`

- Summary: Create privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/privacy-requests/{publicId}`

- Summary: Read one privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/privacy-requests/{publicId}/export-package`

- Summary: Build export package for one privacy request.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/privacy-requests/{publicId}/execute`

- Summary: Execute one privacy request with audit trail.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/fiscal/privacy-requests/{publicId}/status`

- Summary: Transition privacy request lifecycle status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/consents`

- Summary: List consent ledger.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/fiscal/consents`

- Summary: Create consent record.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/fiscal/consents/{publicId}`

- Summary: Transition consent status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`, `404`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/audit-events`

- Summary: List fiscal audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/fiscal/compliance/summary`

- Summary: Read fiscal compliance summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `identity`

- Contrato: `docs/contracts/http/identity.openapi.yaml`
- Total: `6`

### `GET /api/identity/tenants`

- Summary: List tenants.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/identity/tenants`

- Summary: Create tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/identity/tenants/{slug}/snapshot`

- Summary: Read one tenant snapshot.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/identity/sessions/login`

- Summary: Authenticate identity session.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/identity/sessions/refresh`

- Summary: Refresh identity session.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/identity/invitations`

- Summary: Create invitation.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `notification`

- Contrato: `docs/contracts/http/notification.openapi.yaml`
- Total: `7`

### `GET /api/notification/capabilities`

- Summary: Read notification capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/notification/preferences/{userPublicId}`

- Summary: Read one user notification preference.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/notification/preferences/{userPublicId}`

- Summary: Upsert one user notification preference.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/notification/center`

- Summary: List notification center items with cursor filters.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/notification/center`

- Summary: Create one notification center item.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/notification/center/{publicId}/status`

- Summary: Transition notification status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/notification/summary`

- Summary: Read notification summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `platform-control`

- Contrato: `docs/contracts/http/platform-control.openapi.yaml`
- Total: `40`

### `GET /api/platform-control/capabilities/catalog`

- Summary: List platform capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/providers/catalog`

- Summary: List provider capability catalog and environment posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/entitlements`

- Summary: List tenant entitlements with cursor pagination.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/feature-flags`

- Summary: List tenant feature flags with capability metadata.
- Parametros: `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`

- Summary: Upsert one entitlement.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}`

- Summary: Upsert one feature flag using entitlement governance.
- Parametros: `capabilityKey`, `tenantSlug`.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`

- Summary: Bulk upsert entitlements with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`

- Summary: List provider defaults selected for one tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`

- Summary: Upsert provider default for one tenant capability.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/quotas`

- Summary: List quotas by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`

- Summary: Upsert one quota.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`

- Summary: Bulk upsert quotas with partial success.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/blocks`

- Summary: List tenant blocks.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`

- Summary: Upsert tenant block.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/metering`

- Summary: Read metering snapshots and summary with cursor pagination.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`

- Summary: Create one usage snapshot.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`

- Summary: Read quota and metering utilization summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`

- Summary: Read tenant lifecycle readiness and provider posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`

- Summary: List onboarding and offboarding jobs with cursor pagination.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`

- Summary: Read one lifecycle job with audit events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`

- Summary: Preview onboarding plan, provider defaults and lifecycle readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`

- Summary: Queue onboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `202`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`

- Summary: Preview offboarding plan, retention posture and lifecycle readiness.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`

- Summary: Queue offboarding job with Idempotency-Key and 202 Accepted.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `202`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`

- Summary: Transition lifecycle job to running.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`

- Summary: Transition lifecycle job to completed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`

- Summary: Transition lifecycle job to failed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`

- Summary: Transition lifecycle job to cancelled.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness`

- Summary: Read go-live rollout readiness by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption`

- Summary: Read tenant go-live adoption baseline and gap.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks`

- Summary: List go-live bottlenecks and operational blockers.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook`

- Summary: Read rollout and rollback playbook for one tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments`

- Summary: List recommended go-live adjustments.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply`

- Summary: Apply one go-live operational adjustment.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: List go-live rollouts.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts`

- Summary: Create one go-live rollout.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}`

- Summary: Read one go-live rollout with events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start`

- Summary: Transition go-live rollout to running.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete`

- Summary: Transition go-live rollout to completed.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback`

- Summary: Roll back one go-live rollout.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `rentals`

- Contrato: `docs/contracts/http/rentals.openapi.yaml`
- Total: `4`

### `GET /api/rentals/contracts`

- Summary: List rental contracts.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/rentals/contracts`

- Summary: Create rental contract.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/rentals/contracts/{publicId}/charges`

- Summary: List contract charges.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`

- Summary: Update charge status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `sales`

- Contrato: `docs/contracts/http/sales.openapi.yaml`
- Total: `6`

### `GET /api/sales/opportunities`

- Summary: List opportunities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/sales/opportunities`

- Summary: Create opportunity.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/sales/proposals`

- Summary: List proposals.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/sales/proposals`

- Summary: Create proposal.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/sales/sales`

- Summary: List sales.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/sales/invoices`

- Summary: List commercial invoices.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `simulation`

- Contrato: `docs/contracts/http/simulation.openapi.yaml`
- Total: `3`

### `GET /api/simulation/scenarios`

- Summary: List scenarios.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/simulation/scenarios`

- Summary: Create scenario run.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/simulation/benchmarks/load`

- Summary: Execute one load benchmark run.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `supplier`

- Contrato: `docs/contracts/http/supplier.openapi.yaml`
- Total: `8`

### `GET /api/supplier/capabilities`

- Summary: Read supplier capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/supplier/categories`

- Summary: List supplier categories.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/supplier/categories/{categoryKey}`

- Summary: Upsert one supplier category.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/supplier/suppliers`

- Summary: List suppliers by tenant and status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/supplier/suppliers`

- Summary: Create one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/supplier/suppliers/summary`

- Summary: Read supplier summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/supplier/suppliers/{publicId}`

- Summary: Read one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/supplier/suppliers/{publicId}`

- Summary: Update one supplier.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `support`

- Contrato: `docs/contracts/http/support.openapi.yaml`
- Total: `9`

### `GET /api/support/capabilities`

- Summary: Read support capability catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/support/queues`

- Summary: List support queues by tenant.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PUT /api/support/queues/{queueKey}`

- Summary: Upsert one support queue.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/support/cases`

- Summary: List support cases with cursor filters.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/support/cases`

- Summary: Create one support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `201`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/support/cases/summary`

- Summary: Read support case summary.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/support/cases/{publicId}`

- Summary: Read one support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/support/cases/{publicId}/status`

- Summary: Transition support case status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/support/cases/{publicId}/comments`

- Summary: Append comment to support case.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `webhook-hub`

- Contrato: `docs/contracts/http/webhook-hub.openapi.yaml`
- Total: `13`

### `GET /health/details`

- Summary: Return readiness details for webhook runtime.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/webhook-hub/capabilities`

- Summary: Read outbound webhook capability posture.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/webhook-hub/outbound-endpoints`

- Summary: List tenant outbound endpoints.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/webhook-hub/outbound-endpoints`

- Summary: Register one tenant outbound endpoint.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/webhook-hub/outbound-endpoints/{publicId}`

- Summary: Read one tenant outbound endpoint.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: List outbound delivery log for one endpoint.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries`

- Summary: Register one outbound delivery attempt.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter`

- Summary: Move one outbound delivery to dead letter.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/webhook-hub/events`

- Summary: List inbound webhook events.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/webhook-hub/events`

- Summary: Register inbound webhook event.
- Parametros: nenhum parametro declarado.
- Request body: sim.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/webhook-hub/events/summary`

- Summary: Aggregate inbound webhook state.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/webhook-hub/events/{publicId}/dead-letter`

- Summary: Move event to dead letter queue.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/webhook-hub/events/{publicId}/requeue`

- Summary: Requeue dead-letter event.
- Parametros: `publicId`.
- Request body: nao declarado.
- Respostas: nenhuma resposta declarada.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `workflow-control`

- Contrato: `docs/contracts/http/workflow-control.openapi.yaml`
- Total: `7`

### `GET /api/workflow-control/definitions`

- Summary: List workflow definitions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/workflow-control/definitions`

- Summary: Create workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/workflow-control/definitions/{key}`

- Summary: Read one workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/workflow-control/definitions/{key}`

- Summary: Update one workflow definition.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `PATCH /api/workflow-control/definitions/{key}/status`

- Summary: Update workflow definition status.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/workflow-control/capabilities/triggers`

- Summary: List workflow trigger catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/workflow-control/capabilities/actions`

- Summary: List workflow action catalog.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

## Rotas de `workflow-runtime`

- Contrato: `docs/contracts/http/workflow-runtime.openapi.yaml`
- Total: `6`

### `GET /api/workflow-runtime/executions`

- Summary: List workflow executions.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/workflow-runtime/executions`

- Summary: Create workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/workflow-runtime/executions/{publicId}`

- Summary: Read one workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/workflow-runtime/executions/{publicId}/actions`

- Summary: List execution action snapshots.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `POST /api/workflow-runtime/executions/{publicId}/advance`

- Summary: Advance one workflow execution.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

### `GET /api/workflow-runtime/capabilities`

- Summary: List runtime capabilities.
- Parametros: nenhum parametro declarado.
- Request body: nao declarado.
- Respostas: `200`.
- Diagnostico: conferir logs do servico, status HTTP, tenant e correlation id.
- Validacao: usar contrato para divergencia de shape e smoke para fluxo integrado.

