# OPERACOES

## Principios operacionais

- desenvolvimento e testes devem priorizar execucao em containers
- evitar dependencias instaladas diretamente no host sempre que possivel
- comandos pesados de runtime so entram quando agregarem valor real ao miniupdate
- toda validacao importante deve acontecer perto do ambiente de execucao real

## Observabilidade minima esperada

- logs estruturados
- correlation id obrigatorio
- tenant id quando aplicavel
- readiness com dependencias refletindo o runtime real do servico
- metricas de latencia, erro, throughput e filas
- traces ponta a ponta em fluxos criticos

## Runbook esperado por servico

Cada servico critico deve documentar:

- como subir localmente
- como validar health
- como consultar backlog, filas ou retries
- como diagnosticar falhas
- como validar integracoes externas

## Comandos atuais do monorepo

- `./scripts/up.sh` sobe o ecossistema local definido em compose
- `./scripts/down.sh` derruba os containers e remove orfaos
- `./scripts/logs.sh <servico>` consulta logs do servico informado
- o compose local agora tambem sobe `service-kafka`, `service-keycloak`, `service-openfga`, `service-prometheus`, `service-grafana` e `service-blackbox-exporter` como fundacao de plataforma
- quando o host ja estiver usando portas como `5432`, `6379` ou `808x`, os scripts container-first remapeiam automaticamente as exposicoes locais para manter o fluxo de runtime e smoke previsivel
- `./scripts/db.sh migrate all` aplica a base `common`, `identity` e `crm`
- `./scripts/db.sh migrate crm` aplica apenas o contexto relacional de CRM
- `./scripts/db.sh migrate workflow-control` aplica apenas o contexto relacional de workflow-control
- `./scripts/db.sh seed identity` aplica o bootstrap relacional do contexto `identity`
- `./scripts/db.sh seed crm` aplica o bootstrap relacional do contexto `crm`
- `./scripts/db.sh seed workflow-control` aplica o bootstrap relacional do contexto `workflow-control`
- `./scripts/db.sh backup /tmp/erp-backup.sql` gera um dump completo do PostgreSQL local para restauracao controlada
- `./scripts/db.sh restore /tmp/erp-backup.sql` reaplica um dump completo no banco local ativo
- `./scripts/db.sh summary identity <tenant-slug>` resume companies, users, teams, roles e bindings do tenant
- `./scripts/db.sh summary crm <tenant-slug>` resume total de leads, distribuicao por status e ownership por tenant
- `./scripts/db.sh summary workflow-control <tenant-slug>` resume total de definicoes, versoes publicadas, runs, eventos, distribuicao por status e a ultima versao publicada por tenant
- `workflow-control` no compose local sobe em `postgres`, usando `bootstrap-ops` como tenant bootstrap do catalogo inicial
- `./scripts/test.sh unit` executa Go, TypeScript, .NET e Rust em modo container-first
- `./scripts/test.sh integration` executa a suite HTTP do `identity`
- `./scripts/test.sh contract` executa as suites publicas de contratos de `workflow-control`, `crm` e `identity`
- `./scripts/test.sh platform` valida a plataforma local da Fase 1, checando Keycloak, OpenFGA, Kafka, Prometheus e Grafana em container
- `./scripts/test.sh smoke` agora valida primeiro a plataforma local e depois reseta volume, aplica bootstrap relacional e exercita `workflow-control`, `crm`, `sales`, `engagement`, `analytics`, `simulation`, `identity`, `webhook-hub`, `workflow-runtime` e `edge` ao vivo por HTTP
- `./scripts/test.sh backup-restore` valida um ciclo destrutivo controlado de dump e restauracao do PostgreSQL, preservando registros de workload em `simulation`

## Enderecos locais da plataforma

- Keycloak: `http://localhost:${KEYCLOAK_PORT}` com realm bootstrap em `http://localhost:${KEYCLOAK_PORT}/realms/erp-local`
- OpenFGA HTTP API: `http://localhost:${OPENFGA_HTTP_PORT}`
- OpenFGA Playground: `http://localhost:${OPENFGA_PLAYGROUND_PORT}/playground`
- Prometheus: `http://localhost:${PROMETHEUS_PORT}`
- Grafana: `http://localhost:${GRAFANA_PORT}`
- Kafka: `localhost:${KAFKA_PORT}`

## Observabilidade local bootstrap

- o Prometheus usa `blackbox_exporter` para sondar health HTTP e conectividade TCP do stack, incluindo o plano local de autorizacao em OpenFGA
- o Grafana sobe com datasource provisionado para o Prometheus e com o dashboard `ERP Platform Health`
- o stack inicial ja permite enxergar saude de `edge`, `identity`, `webhook-hub`, Keycloak, OpenFGA, Grafana, PostgreSQL, Redis e Kafka sem depender de instrumentacao pesada por servico

## Entrega incremental

- miniupdates pequenos e frequentes
- changelog atualizado a cada entrega relevante
- testes executados antes de subir uma entrega quando houver runtime suficiente para isso
- commits curtos e versionados de forma progressiva

## Automacao no GitHub

- `.github/workflows/quality.yml` roda `unit`, `integration` e `contract` em `push` e `pull_request`, e executa `smoke` em `main` ou por `workflow_dispatch`
- `.github/workflows/containers.yml` publica as imagens de `edge`, `crm`, `sales`, `identity`, `workflow-control`, `workflow-runtime`, `analytics`, `simulation` e `webhook-hub` no `ghcr.io`
- o registro publica `latest` na branch principal, `sha-*` por commit e tags por release quando o push vier de `v*`
- as imagens passam a aparecer na aba `Packages` do repositório depois da primeira execucao bem-sucedida do workflow de containers

## Enderecos esperados no GHCR

- `ghcr.io/thiagodifaria/erp-edge`
- `ghcr.io/thiagodifaria/erp-crm`
- `ghcr.io/thiagodifaria/erp-sales`
- `ghcr.io/thiagodifaria/erp-identity`
- `ghcr.io/thiagodifaria/erp-workflow-control`
- `ghcr.io/thiagodifaria/erp-workflow-runtime`
- `ghcr.io/thiagodifaria/erp-analytics`
- `ghcr.io/thiagodifaria/erp-simulation`
- `ghcr.io/thiagodifaria/erp-webhook-hub`
