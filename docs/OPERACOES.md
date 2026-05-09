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
```

## Runtime Local

`./scripts/build.sh` usa `infra/docker-compose.yml` e carrega `.env` quando existir. Se `.env` nao existir, usa `.env.example`.

O script tambem tenta remapear portas locais ocupadas para evitar falha simples de subida do stack. Para diagnostico de porta, confira a saida do proprio comando e use:

```bash
./scripts/build.sh ps
```

## Console Tecnico da API

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

Use o console para navegar documentacao, explorar contratos e testar endpoints contra o backend local.

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
