# OPERACOES

Este documento descreve como executar, validar e diagnosticar o ERP localmente. Ele nao explica arquitetura de dominio nem lista endpoints.

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

Esse overlay remove a publicacao direta dos servicos internos e deixa gateway/edge como pontos de entrada.

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

Validacao:

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

`./scripts/build.sh` usa `infra/docker-compose.yml` e carrega `.env` quando existir. Se `.env` nao existir, usa `.env.example`.

O script tambem tenta remapear portas locais ocupadas para evitar falha simples de subida do stack. Para diagnostico de porta, confira a saida do proprio comando e use:

```bash
./scripts/build.sh ps
```

Em `ERP_ENV=production`, fallbacks inseguros devem falhar: `ERP_ALLOW_BOOTSTRAP_TENANT_FALLBACK=false`, `DOCUMENTS_ACCESS_TOKEN_SECRET` obrigatorio e secrets criticos diferentes dos defaults locais.

Em perfil produtivo, a postura esperada tambem inclui `ERP_SECURITY_LEVEL=strict`, `ERP_REQUIRE_REQUEST_SIGNATURE=true`, `ERP_SECURITY_HEADERS=enabled`, `ERP_AUDIT_LOG_REDACTION=strict`, `WEBHOOK_HUB_REQUIRE_SIGNATURE=true`, janela curta de replay em webhooks e `DOCUMENTS_MALWARE_SCAN_MODE=required`.

## Console Tecnico da API

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

Use o console para navegar documentacao, explorar contratos e testar endpoints contra o backend local.

## Gateway Local

O compose sobe `gateway` como ponto unico de entrada HTTP para `/api/<servico>/`. Ele aplica cache de leitura, rate limit, timeouts, failover passivo por dependencia e encaminha `X-Request-ID`.

Validacao rapida:

```bash
curl -i "http://localhost:${GATEWAY_HTTP_PORT}/gateway/health"
curl -i "http://localhost:${GATEWAY_HTTP_PORT}/api/edge/ops/service-pulse"
```

## Probes Esperadas

Quando um servico HTTP estiver ativo:

- `/health/live` indica processo vivo;
- `/health/ready` indica readiness operacional;
- `/health/details` ajuda diagnostico sem expor segredo.

Nem todo contrato lista os tres probes, mas servicos de runtime devem seguir esse padrao quando participam do stack.

## SLOs E Alertas

| Jornada | Indicador | Objetivo local/corporativo |
| --- | --- | --- |
| Auth/login | taxa de sucesso e p95 | 99.5% de sucesso, p95 abaixo de 800ms |
| CRM lead intake | criacao de lead e outbox | 99.0% de sucesso, sem backlog critico |
| Sales creation | criacao de oportunidade/sale/invoice | 99.0% de sucesso, idempotencia preservada |
| Billing payment attempt | tentativa e callback | 99.0% processado ou DLQ explicita |
| Webhook ingestion | validacao, fila, forward/DLQ | 99.0% sem perda, replay auditavel |
| Document download | access link valido para redirect | 99.5% de sucesso, revogacao imediata |
| Fiscal issue/cancel | documento/evento fiscal | 99.0% com trilha de auditoria |

Alertas basicos: erro 5xx acima do limiar do dominio, p95/p99 fora do SLO por 10 minutos, DLQ crescendo sem drenagem, provider fora, backup/restore falhando, divergencia financeira e repeticao anormal de login, recovery, upload ou webhook.

## Banco de Dados

Migrations vivem em:

```text
service-api/service-postgresql/<contexto>/migrations
```

Seeds vivem em:

```text
service-api/service-postgresql/<contexto>/seeds
```

Dominios migrados pelo comando central incluem:

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

## Suites de Validacao

| Suite | Uso |
|-------|-----|
| `unit` | regras locais por stack/servico |
| `integration` | integracoes de runtime quando aplicavel |
| `contract` | OpenAPI, event schemas e registry |
| `platform` | baseline de plataforma e infraestrutura |
| `smoke` | fluxo real de ponta a ponta |
| `performance` | carga e capacidade local |
| `backup-restore` | backup e restauracao do PostgreSQL |
| `hardening` | readiness, contratos, providers e postura operacional |
| `security` | guardrails de auth, secrets, eventos, documents, dados pessoais e padroes |
| `supply-chain` | secret scan de alta confianca, inventario SBOM e pinning de imagens |
| `production-readiness` | gate 1.0.0, manifests Kubernetes, docs oficiais e postura de ownership |

Para o hardening enterprise, o conjunto minimo de evidencia operacional e formado por `contract`, `smoke`, `performance`, `backup-restore` e `hardening`. O relatorio `GET /api/analytics/reports/hardening-review` consolida esse fechamento como `operationalRunbooks`, permitindo validar rapidamente se seguranca operacional, observabilidade, DLQ/retry, backup/restore, SLOs, multi-tenant, failover, performance e permissoes possuem cobertura operacional rastreavel.

Para o hardening de seguranca 1.0.0, `security` valida tambem headers defensivos, limite de body, correlacao na borda, Pod Security restricted, NetworkPolicy deny-by-default, request signature, redaction estrita, scan obrigatorio de documentos e assinatura/replay de webhooks.

## Production Readiness 1.0.0

A versao 1.0.0 e o gate de producao do projeto. Ela nao adiciona um novo dominio funcional; ela consolida o ERP como plataforma implantavel, auditavel, observavel, recuperavel e segura por padrao.

O gate oficial e exposto em:

```bash
GET /api/analytics/reports/production-readiness?tenant_slug=bootstrap-ops
```

O relatorio retorna `release.version=1.0.0`, `release.releaseReady`, gates de trafego, auth/tenant, secrets, contratos, observabilidade, backup/DR, deploy, providers e go-live, alem das evidencias que precisam existir para aceitar a entrega.

Validacao completa do release:

```bash
docker compose --env-file .env.production.example -f infra/docker-compose.yml -f infra/docker-compose.corporate-like.yml config >/dev/null
./scripts/test.sh contract
./scripts/test.sh security
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

O deploy corporativo considera obrigatorio:

- namespace dedicado;
- secrets injetados fora do repositorio;
- config sem credenciais;
- job de migration antes do rollout;
- `edge` exposto por ingress TLS;
- servicos internos como `ClusterIP`;
- NetworkPolicy com deny-by-default;
- probes de live/readiness;
- `readOnlyRootFilesystem`, `allowPrivilegeEscalation=false` e drop de capabilities;
- HPA inicial para o ponto de entrada;
- rollback por revisao de deployment ou estrategia controlada.

Se provider externo estiver em fallback/manual/unconfigured, o release pode continuar tecnicamente pronto, mas a capability deve aparecer como nao produtiva no readiness de providers. O sistema nao deve mascarar fallback local como integracao real.

## Runbook Rapido

### Stack nao sobe

```bash
./scripts/build.sh ps
./scripts/build.sh logs service-postgresql
./scripts/build.sh logs edge
```

Verifique Docker ativo, portas ocupadas e `.env`.

### Servico responde live mas nao ready

```bash
./scripts/build.sh logs <servico>
./scripts/build.sh psql
./scripts/test.sh contract
```

Normalmente o problema esta em dependencia, migration ou variavel de ambiente.

### API Explorer nao chama backend

Confirme que o backend esta ativo e que o console foi iniciado pelo Vite:

```bash
./scripts/build.sh ps
cd client-web/client-api
npm run dev
```

O console usa proxy local; se o servico de destino nao estiver no ar, a chamada falha corretamente.

### Contrato fora de sincronia

```bash
./scripts/test.sh contract
cd client-web/client-api
npm run generate
```

Atualize OpenAPI, registry e implementacao no mesmo fluxo de mudanca.

### DLQ ou retry crescendo

1. Conferir `/api/webhook-hub/events/dead-letter` e summary.
2. Validar `correlationId`, provider, payload e schemaRef.
3. Reprocessar apenas eventos idempotentes.
4. Se o erro for contrato, bloquear replay ate corrigir schema/adapter.

### Divergencia financeira

1. Conferir cash movements, settlement/payment reference e period closure.
2. Validar se a mutacao usou `Idempotency-Key`.
3. Registrar ajuste ou estorno em vez de update destrutivo.
4. Bloquear fechamento de periodo se a divergencia continuar.

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

Depois de restore, rode uma validacao proporcional:

```bash
./scripts/test.sh smoke
```

## Go-live Operacional

O go-live progressivo e controlado por tenant. O objetivo operacional e liberar ondas pequenas, observar metricas reais, manter rollback claro e registrar ajustes finos sem depender de planilhas externas.

O `platform-control` concentra o controle transacional:

- `GET /api/platform-control/tenants/{tenantSlug}/go-live/readiness` mostra readiness de rollout, rollback e metricas;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts` cria uma onda com `targetEnv`, `waveKey`, `rollbackPlaybook` e `adoptionTargetPct`;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start` inicia a onda;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete` conclui uma onda iniciada;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback` registra rollback controlado;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adoption` acompanha adocao e gap contra alvo;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks` lista gargalos de providers, quotas e bloqueios;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/playbook` devolve o checklist operacional;
- `GET /api/platform-control/tenants/{tenantSlug}/go-live/adjustments` sugere ajustes aplicaveis;
- `POST /api/platform-control/tenants/{tenantSlug}/go-live/adjustments/{adjustmentKey}/apply` aplica ajustes suportados.

O `analytics` consolida o fechamento em `GET /api/analytics/reports/go-live-control`, incluindo `rollouts`, `adoption`, `bottlenecks`, `adjustments`, `readiness` e `releaseControls`. O bloco `releaseControls` e a evidencia rapida de que rollout por tenant, monitoramento de adocao, rollback, observacao de gargalos e ajuste por uso estao cobertos por runbook e suites oficiais.

O `edge` publica a visao executiva em `GET /api/edge/ops/go-live-overview?tenantSlug=bootstrap-ops`, agregando `service-pulse`, `saas-control`, `go-live-control` e um `executiveSummary` com `rolloutReady`, `rollbackReady`, `metricsObserved` e `acceptanceReady`.

Validacao recomendada:

```bash
./scripts/test.sh smoke
./scripts/test.sh hardening
```

Considere a entrega pronta quando existir pelo menos uma onda auditavel, as metricas de uso estiverem observadas, o rollback estiver disponivel, os gargalos estiverem visiveis e os ajustes suportados puderem ser aplicados por API.

## Observabilidade Local

O stack possui assets de observabilidade em `infra/`. Use logs por servico para diagnostico imediato e Prometheus/Grafana quando o stack completo estiver ativo.

Mutacoes importantes devem carregar correlation id e registrar status final. Providers externos devem registrar provider, tentativa e erro normalizado.

## Encerramento Limpo

```bash
./scripts/build.sh down
```

Use `down` antes de trocar de branch ou alterar migrations grandes, para evitar estado local enganoso.

## Fluxo Diario Recomendado

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

Use `unit` quando alterar regra local de servico.

Use `integration` quando alterar dependencia real, repository driver, banco ou runtime entre componentes.

Use `contract` quando alterar OpenAPI, event schema, registry ou documentacao contratual.

Use `platform` quando alterar compose, infraestrutura local, portas, health, observabilidade ou bootstrap.

Use `smoke` quando alterar fluxo que cruza dominios.

Use `hardening` quando alterar readiness, provider posture, governanca, seguranca operacional ou go-live.

Use `backup-restore` quando alterar migrations, scripts de banco ou formato de persistencia.

## Diagnostico Por Sintoma

| Sintoma | Primeira checagem | Segunda checagem |
|---------|-------------------|------------------|
| container reiniciando | `./scripts/build.sh logs <servico>` | variaveis e migration |
| ready falhando | `/health/details` | dependencia indisponivel |
| contrato falhando | `./scripts/test.sh contract` | OpenAPI e registry |
| smoke falhando | logs de `edge` e servico alvo | seed/migration |
| API console sem resposta | `./scripts/build.sh ps` | proxy/Vite e porta do servico |
| erro de porta | saida do `build.sh` | `.env` e processos locais |
| banco inconsistente | `migrate all` e `summary` | backup/restore |

## Ordem Segura Para Mudancas Grandes

1. Rode `./scripts/test.sh contract` para saber baseline.
2. Altere contrato ou schema.
3. Altere implementacao.
4. Altere migration/seed quando houver persistencia.
5. Rode unidade do servico.
6. Rode contrato.
7. Rode smoke se cruzar dominio.
8. Atualize docs no arquivo certo.
9. Atualize changelog quando a mudanca estiver fechada.

## Operacao do Client-api

O console tecnico tem dois modos de uso:

- leitura local de documentacao e contratos;
- execucao de chamadas reais contra backend local.

Quando o backend esta desligado, a UI pode abrir e navegar docs, mas chamadas reais falham. Isso e esperado e ajuda a diferenciar problema de frontend de problema de runtime.

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

Logs locais devem ajudar diagnostico sem expor segredo.

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
- payload fiscal sensivel completo;
- documento ou dado pessoal desnecessario.
