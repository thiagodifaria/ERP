# operações

Este documento descreve como executar, validar e diagnosticar O projeto localmente. Ele não explica arquitetura de domínio nem lista endpoints.

## Comandos Oficiais

Runtime e infraestrutura:

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh down
./scripts/build.sh restart
./scripts/build.sh ps
./scripts/build.sh logs edge
```

Perfil corporate-like:

```bash
docker compose \
  --env-file .env.production.example \
  -f infra/docker-compose.yml \
  -f infra/docker-compose.corporate-like.yml \
  config
```

Esse overlay remove a públicação direta dos serviços internos e deixa gateway/edge como pontos de entrada.

Banco:

```bash
./scripts/build.sh migrate all
./scripts/build.sh migrate identity
./scripts/build.sh seed all
./scripts/build.sh seed crm
./scripts/build.sh summary crm
./scripts/build.sh psql
./scripts/build.sh backup /tmp/erp-local-backup.sql
ERP_BACKUP_ENCRYPTION_KEY="$SECRET_VALUE" ./scripts/build.sh backup-encrypted /tmp/erp-local-backup.sql.enc
./scripts/build.sh restore /tmp/erp-local-backup.sql
ERP_BACKUP_ENCRYPTION_KEY="$SECRET_VALUE" ./scripts/build.sh restore-encrypted /tmp/erp-local-backup.sql.enc
```

validação:

```bash
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh contract
./scripts/test.sh platform
./scripts/test.sh smoke
./scripts/test.sh performance
./scripts/test.sh backup-restore
./scripts/test.sh hardening
./scripts/test.sh security
./scripts/test.sh supply-chain
./scripts/test.sh production-readiness
./scripts/test.sh hardening-secrets
```

## Runtime Local

`./scripts/build.sh` usa `infra/docker-compose.yml` e carrega `.env` quando existir. Se `.env` não existir, usa `.env.example`.

O script também tenta remapear portas locais ocupadas para evitar falha simples de subida do stack. Para diagnóstico de porta, confira a saída do próprio comando e use:

```bash
./scripts/build.sh ps
```

Em `ERP_ENV=production`, fallbacks inseguros devem falhar: `ERP_ALLOW_BOOTSTRAP_TENANT_FALLBACK=false`, `DOCUMENTS_ACCESS_TOKEN_SECRET` obrigatório e secrets críticos diferentes dos defaults locais.

Em perfil produtivo, a postura esperada também inclui `ERP_SECURITY_LEVEL=strict`, `ERP_REQUIRE_REQUEST_SIGNATURE=true`, `ERP_SECURITY_HEADERS=enabled`, `ERP_AUDIT_LOG_REDACTION=strict`, `WEBHOOK_HUB_REQUIRE_SIGNATURE=true`, janela curta de replay em webhooks e `DOCUMENTS_MALWARE_SCAN_MODE=required`.

## Console técnico da API

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

Use o console para navegar documentação, explorar contratos e testar endpoints contra o backend local.

## Capacidades Operacionais 1.5.0

Busca operacional e e-discovery:

```bash
curl "http://localhost:${SEARCH_HTTP_PORT}/api/search/query?tenant_slug=bootstrap-ops&q=contrato"
curl "http://localhost:${SEARCH_HTTP_PORT}/api/search/facets?tenant_slug=bootstrap-ops"
```

BI semântico e catálogo de métricas:

```bash
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/metrics"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/data-quality"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/lineage"
```

governança de IA:

```bash
curl "http://localhost:${AI_GOVERNANCE_HTTP_PORT}/api/ai-governance/tools"
curl "http://localhost:${AI_GOVERNANCE_HTTP_PORT}/api/ai-governance/policies"
```

Comando de incidentes:

```bash
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/tenants/bootstrap-ops/incidents"
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/tenants/bootstrap-ops/incident-command/readiness"
```

Esses fluxos entram no gate `production-readiness` como capacidades produtivas: busca auditada com redação, catálogo semântico com qualidade/lineage, IA limitada a ferramentas aprovadas e incidente com timeline, action items, resolução e postmortem.

Hardening poliglota da versão:

```bash
./scripts/test.sh security
./scripts/test.sh production-readiness
```

A suite bloqueia regressões de geração insegura de IDs, panic por UUID em persistência Go, alocação insegura de clients HTTP em middleware, prefixo duplicado de migration e ausência do manifesto de contrato de tenant. Se esse gate falhar, corrija a causa no domínio dono antes de tentar subir o runtime.

governança operacional autonoma:

```bash
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/policies/catalog"
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/tenants/bootstrap-ops/timeline"
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/runbooks/catalog"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/risk/tenant-score?tenant_slug=bootstrap-ops"
```

Runtime empresarial:

```bash
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/event-mesh/catalog"
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/tenants/bootstrap-ops/event-mesh/lineage"
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/tenants/bootstrap-ops/runtime/profile"
curl "http://localhost:${PLATFORM_CONTROL_HTTP_PORT}/api/platform-control/contracts/evolution"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/reconciliation/findings?tenant_slug=bootstrap-ops"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/financial-close/readiness?tenant_slug=bootstrap-ops"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/master-data/quality-score?tenant_slug=bootstrap-ops"
curl "http://localhost:${ANALYTICS_HTTP_PORT}/api/analytics/lakehouse/datasets"
```

Na `v1.5.0`, policy decision center, command approvals, runbook automation, audit evidence vault, event mesh, tenant runtime, contract evolution, reconciliação, fechamento financeiro, lakehouse, ativação BYOK, inteligência externa e hardening transversal entram no aceite operacional. Comandos, eventos, chamadas externas e sinais de risco devem registrar decisão, timeline, hash ou evidência antes de serem tratados como executados.

## Gateway Local

O Compose sobe `gateway` como ponto único de entrada HTTP para `/api/<serviço>/`. Ele aplica cache de leitura, rate limit, timeouts, failover passivo por dependência e encaminha `X-Request-ID`.

validação rápida:

```bash
curl -i "http://localhost:${GATEWAY_HTTP_PORT}/gateway/health"
curl -i "http://localhost:${GATEWAY_HTTP_PORT}/api/edge/ops/service-pulse"
```

## Probes Esperadas

Quando um serviço HTTP estiver ativo:

- `/health/live` indica processo vivo;
- `/health/ready` indica readiness operacional;
- `/health/details` ajuda diagnóstico sem expor segredo.

Nem todo contrato lista os tres probes, mas serviços de runtime devem seguir esse padrão quando participam do stack.

## SLOs e Alertas

| Jornada | Indicador | Objetivo local/corporativo |
| --- | --- | --- |
| Auth/login | taxa de sucesso e p95 | 99.5% de sucesso, p95 abaixo de 800ms |
| CRM lead intake | criação de lead e outbox | 99.0% de sucesso, sem backlog crítico |
| Sales creation | criação de oportunidade/sale/invoice | 99.0% de sucesso, idempotência preservada |
| Billing payment attempt | tentativa e callback | 99.0% processado ou DLQ explicita |
| Webhook ingestion | validação, fila, forward/DLQ | 99.0% sem perda, replay auditável |
| Document download | access link valido para redirect | 99.5% de sucesso, revogação imediata |
| Fiscal issue/cancel | documento/evento fiscal | 99.0% com trilha de auditoria |

Alertas basicos: erro 5xx acima do limiar do domínio, p95/p99 fora do SLO por 10 minutos, DLQ crescendo sem drenagem, provider fora, backup/restore falhando, divergencia financeira e repetição anormal de login, recovery, upload ou webhook.

## Banco de Dados

Migrations vivem em:

```text
service-api/service-postgresql/<contexto>/migrations
```

Seeds vivem em:

```text
service-api/service-postgresql/<contexto>/seeds
```

domínios migrados pelo comando central incluem:

```text
common
identity
crm
sales
rentals
finance
billing
documents
analytics
simulation
catalog
platform-control
support
supplier
notification
fiscal
engagement
webhook-hub
workflow-control
workflow-runtime
tenant-security
```

## Suites de validação

| Suite | Uso |
|-------|-----|
| `unit` | regras locais por stack/serviço |
| `integration` | integrações de runtime quando aplicável |
| `contract` | OpenAPI, event schemas e registry |
| `platform` | baseline de plataforma e infraestrutura |
| `smoke` | fluxo real de ponta a ponta |
| `performance` | carga e capacidade local |
| `backup-restore` | backup e restauração do PostgreSQL |
| `hardening` | readiness, contratos, providers e postura operacional |
| `security` | guardrails de auth, secrets, eventos, documents, dados pessoais e padrões |
| `supply-chain` | secret scan de alta confiança, inventário SBOM e pinning de imagens |
| `production-readiness` | aceite 1.5.0, manifests Kubernetes, docs oficiais e postura de ownership |

Para o hardening enterprise, o conjunto minimo de evidência operacional e formado por `contract`, `smoke`, `performance`, `backup-restore` e `hardening`. O relatório `GET /api/analytics/reports/hardening-review` consolida esse fechamento como `operationalRunbooks`, permitindo validar rápidamente se segurança operacional, observabilidade, DLQ/retry, backup/restore, SLOs, multi-tenant, failover, performance e permissões possuem cobertura operacional rastreável.

Para o hardening de segurança 1.0.0, `security` valida também headers defensivos, limite de body, correlação na borda, Pod Security restricted, NetworkPolicy deny-by-default, request signature, redaction estrita, scan obrigatório de documentos e assinatura/replay de webhooks.

## Production Readiness 1.5.0

A versão 1.5.0 preserva governança operacional, event mesh, reconciliação, fechamento financeiro, dados mestres, lakehouse, tenant runtime, ativação BYOK, OCR/document intelligence, fiscal Brasil, enriquecimento cadastral brasileiro, mercado/macro e feeds externos de risco, e acrescenta a linha de hardening `1.4.x`: raiz/env, static policy, geração transacional de IDs, console seguro, infra runtime e conformance de auth/observabilidade.

O gate oficial é exposto em:

```bash
GET /api/analytics/reports/production-readiness?tenant_slug=bootstrap-ops
```

O relatório retorna `release.version=1.5.0`, `release.releaseReady`, gates de tráfego, auth/tenant, secrets, root hardening, static policy, geração transacional de IDs, contratos, observabilidade, backup/DR, deploy, providers, provider activation, LLM BYOK, document intelligence, fiscal Brazil, registry enrichment, market macro risk, external risk feed, go-live, policies, timeline, approvals, runbooks, evidence vault, risk scoring, event mesh, financial close, master data, lakehouse, tenant runtime, contract evolution, console técnico e conformance de plataforma, além das evidências que precisam existir para aceitar a entrega.

validação completa da versão:

```bash
docker compose --env-file .env.production.example -f infra/docker-compose.yml -f infra/docker-compose.corporate-like.yml config >/dev/null
./scripts/test.sh contract
./scripts/test.sh security
./scripts/test.sh auth-conformance
./scripts/test.sh observability
./scripts/test.sh supply-chain
./scripts/test.sh hardening-secrets
./scripts/test.sh backup-restore
./scripts/test.sh hardening
./scripts/test.sh smoke
./scripts/test.sh production-readiness
```

O caminho Kubernetes oficial fica em `infra/kubernetes/`:

```bash
kubectl apply --dry-run=server -k infra/kubernetes/overlays/production
kubectl -n erp rollout status deployment/erp-edge
kubectl -n erp get networkpolicy
```

O deploy corporativo considera obrigatório:

- namespace dedicado;
- secrets injetados fora do repositório;
- config sem credenciais;
- job de migration antes do rollout;
- `edge` exposto por ingress TLS;
- serviços internos como `ClusterIP`;
- NetworkPolicy com deny-by-default;
- probes de live/readiness;
- `readOnlyRootFilesystem`, `allowPrivilegeEscalation=false` e drop de capabilities;
- HPA inicial para o ponto de entrada;
- matriz de deployabilidade em `infra/kubernetes/base/deployability-matrix.yaml`;
- rollback por revisão de deployment ou estrategia controlada.

Se provider externo estiver em fallback/manual/unconfigured, a versão pode continuar técnicamente pronta, mas a capability deve aparecer como não produtiva no readiness de providers. O sistema não deve mascarar fallback local como integração real.

## Runbook rápido

### Stack não sobe

```bash
./scripts/build.sh ps
./scripts/build.sh logs service-postgresql
./scripts/build.sh logs edge
```

Verifique Docker ativo, portas ocupadas e `.env`.

### serviço responde live mas não ready

```bash
./scripts/build.sh logs <serviço>
./scripts/build.sh psql
./scripts/test.sh contract
```

Normalmente o problema está em dependência, migration ou variavel de ambiente.

### API Explorer não chama backend

Confirme que o backend está ativo e que o console foi iniciado pelo Vite:

```bash
./scripts/build.sh ps
cd client-web/client-api
npm run dev
```

O console usa proxy local; se o serviço de destino não estiver no ar, a chamada falha corretamente.

### Contrato fora de sincronia

```bash
./scripts/test.sh contract
cd client-web/client-api
npm run generate
```

Atualize OpenAPI, registry e implementação no mesmo fluxo de mudança.

### DLQ ou retry crescendo

1. Conferir `/api/webhook-hub/events/dead-letter` e summary.
2. Validar `correlationId`, provider, payload e schemaRef.
3. Reprocessar apenas eventos idempotentes.
4. Se o erro for contrato, bloquear replay ate corrigir schema/adapter.

### Divergencia financeira

1. Conferir cash movements, settlement/payment reference e period closure.
2. Validar se a mutação usou `Idempotency-Key`.
3. Registrar ajuste ou estorno em vez de update destrutivo.
4. Bloquear fechamento de período se a divergencia continuar.

## Backup e Restore

Backup local:

```bash
./scripts/build.sh backup /tmp/erp-local-backup.sql
ERP_BACKUP_ENCRYPTION_KEY="$SECRET_VALUE" ./scripts/build.sh backup-encrypted /tmp/erp-local-backup.sql.enc
```

Restore local:

```bash
./scripts/build.sh restore /tmp/erp-local-backup.sql
ERP_BACKUP_ENCRYPTION_KEY="$SECRET_VALUE" ./scripts/build.sh restore-encrypted /tmp/erp-local-backup.sql.enc
```

Depois de restore, rode uma validação proporcional:

```bash
./scripts/test.sh smoke
```

## Go-live Operacional

O go-live progressivo é controlado por tenant. O objetivo operacional é liberar ondas pequenas, observar métricas reais, manter rollback claro e registrar ajustes finos sem depender de planilhas externas.

O `platform-control` concentra o controle transacional:

- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` mostra readiness de rollout, rollback e métricas;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` cria uma onda com `targetEnv`, `waveKey`, `rollbackPlaybook` e `adoptionTargetPct`;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` inicia a onda;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` conclui uma onda iniciada;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` registra rollback controlado;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` acompanha adoção e gap contra alvo;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` lista gargalos de providers, quotas e bloqueios;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` devolve o checklist operacional;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` sugere ajustes aplicaveis;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/{adjustmentKey}/apply` aplica ajustes suportados.

O `analytics` consolida o fechamento em `GET /api/analytics/reports/go-live-control`, incluindo `rollouts`, `adoption`, `bottlenecks`, `adjustments`, `readiness` e `releaseControls`. O bloco `releaseControls` e a evidência rápida de que rollout por tenant, monitoramento de adoção, rollback, observação de gargalos e ajuste por uso estão cobertos por runbook e suites oficiais.

O `edge` pública a visão executiva em `GET /api/edge/ops/go-live-overview?tenantSlug=bootstrap-ops`, agregando `service-pulse`, `saas-control`, `go-live-control` e um `executiveSummary` com `rolloutReady`, `rollbackReady`, `metricsObserved` e `acceptanceReady`.

validação recomendada:

```bash
./scripts/test.sh smoke
./scripts/test.sh hardening
```

Considere a entrega pronta quando existir pelo menos uma onda auditável, as métricas de uso estiverem observadas, o rollback estiver disponível, os gargalos estiverem visíveis e os ajustes suportados puderem ser aplicados por API.

## observabilidade Local

O stack possui assets de observabilidade em `infra/`. Use logs por serviço para diagnóstico imediato e Prometheus/Grafana quando o stack completo estiver ativo.

Mutações importantes devem carregar correlation id e registrar status final. Providers externos devem registrar provider, tentativa e erro normalizado.

## Encerramento Limpo

```bash
./scripts/build.sh down
```

Use `down` antes de trocar de branch ou alterar migrations grandes, para evitar estado local enganoso.

## Fluxo diário Recomendado

Para desenvolvimento comum:

```bash
./scripts/build.sh up
./scripts/test.sh contract
./scripts/test.sh unit
```

Antes de validar fluxo completo:

```bash
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/test.sh smoke
```

Antes de mexer em contrato:

```bash
./scripts/test.sh contract
cd client-web/client-api
npm run generate
```

## Quando Usar Cada Suite

Use `unit` quando alterar regra local de serviço.

Use `integration` quando alterar dependência real, repository driver, banco ou runtime entre componentes.

Use `contract` quando alterar OpenAPI, event schema, registry ou documentação contratual.

Use `platform` quando alterar compose, infraestrutura local, portas, health, observabilidade ou bootstrap.

Use `smoke` quando alterar fluxo que cruza domínios.

Use `hardening` quando alterar readiness, provider posture, governança, segurança operacional ou go-live.

Use `backup-restore` quando alterar migrations, scripts de banco ou formato de persistência.

## diagnóstico Por Sintoma

| Sintoma | Primeira checagem | Segunda checagem |
|---------|-------------------|------------------|
| container reiniciando | `./scripts/build.sh logs <serviço>` | variaveis e migration |
| ready falhando | `/health/details` | dependência indisponível |
| contrato falhando | `./scripts/test.sh contract` | OpenAPI e registry |
| smoke falhando | logs de `edge` e serviço alvo | seed/migration |
| API console sem resposta | `./scripts/build.sh ps` | proxy/Vite e porta do serviço |
| erro de porta | saída do `build.sh` | `.env` e processos locais |
| banco inconsistente | `migrate all` e `summary` | backup/restore |

## Ordem Segura Para mudanças Grandes

1. Rode `./scripts/test.sh contract` para saber baseline.
2. Altere contrato ou schema.
3. Altere implementação.
4. Altere migration/seed quando houver persistência.
5. Rode unidade do serviço.
6. Rode contrato.
7. Rode smoke se cruzar domínio.
8. Atualize docs no arquivo certo.
9. Atualize changelog quando a mudança estiver fechada.

## operação do Client-api

O console técnico tem dois modos de uso:

- leitura local de documentação e contratos;
- execução de chamadas reais contra backend local.

Quando o backend está desligado, a UI pode abrir e navegar docs, mas chamadas reais falham. Isso e esperado e ajuda a diferenciar problema de frontend de problema de runtime.

Comandos uteis:

```bash
cd client-web/client-api
npm run generate
npm run typecheck
npm run build
npm run dev
```

## Cuidados Com Estado Local

Estado local antigo pode mascarar bug ou criar bug falso. Antes de investigar falha estranha:

- confirme branch e diff atual;
- confirme se migrations foram aplicadas;
- confirme se seeds foram reaplicados;
- confirme se containers antigos foram derrubados;
- confira se portas foram remapeadas;
- rode smoke apenas depois do stack estar coerente.

## Politica de Logs

Logs locais devem ajudar diagnóstico sem expor segredo.

Procure por:

- correlation id;
- tenant;
- rota;
- status;
- provider externo;
- erro normalizado.

Evite registrar:

- token;
- senha;
- segredo de provider;
- payload fiscal sensível completo;
- documento ou dado pessoal desnecessário.
