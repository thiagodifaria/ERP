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
./scripts/build.sh restore /tmp/erp-local-backup.sql
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
./scripts/test.sh hardening-secrets
```

## Runtime Local

`./scripts/build.sh` usa `infra/docker-compose.yml` e carrega `.env` quando existir. Se `.env` nao existir, usa `.env.example`.

O script tambem tenta remapear portas locais ocupadas para evitar falha simples de subida do stack. Para diagnostico de porta, confira a saida do proprio comando e use:

```bash
./scripts/build.sh ps
```

Em `ERP_ENV=production`, fallbacks inseguros devem falhar: `ERP_ALLOW_BOOTSTRAP_TENANT_FALLBACK=false`, `DOCUMENTS_ACCESS_TOKEN_SECRET` obrigatorio e secrets criticos diferentes dos defaults locais.

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

Para o hardening enterprise, o conjunto minimo de evidencia operacional e formado por `contract`, `smoke`, `performance`, `backup-restore` e `hardening`. O relatorio `GET /api/analytics/reports/hardening-review` consolida esse fechamento como `operationalRunbooks`, permitindo validar rapidamente se seguranca operacional, observabilidade, DLQ/retry, backup/restore, SLOs, multi-tenant, failover, performance e permissoes possuem cobertura operacional rastreavel.

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

## Backup e Restore

Backup local:

```bash
./scripts/build.sh backup /tmp/erp-local-backup.sql
```

Restore local:

```bash
./scripts/build.sh restore /tmp/erp-local-backup.sql
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
