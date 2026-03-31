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
- `./scripts/db.sh migrate all` aplica a base `common`, `identity` e `crm`
- `./scripts/db.sh migrate crm` aplica apenas o contexto relacional de CRM
- `./scripts/db.sh seed identity` aplica o bootstrap relacional do contexto `identity`
- `./scripts/db.sh seed crm` aplica o bootstrap relacional do contexto `crm`
- `./scripts/db.sh summary identity <tenant-slug>` resume companies, users, teams, roles e bindings do tenant
- `./scripts/db.sh summary crm <tenant-slug>` resume total de leads, distribuicao por status e ownership por tenant
- `./scripts/test.sh unit` executa Go, .NET e Rust em modo container-first
- `./scripts/test.sh integration` executa a suite HTTP do `identity`
- `./scripts/test.sh contract` executa a suite publica de contratos do `identity`
- `./scripts/test.sh smoke` valida compose e o bootstrap relacional de `identity` e `crm`

## Entrega incremental

- miniupdates pequenos e frequentes
- changelog atualizado a cada entrega relevante
- testes executados antes de subir uma entrega quando houver runtime suficiente para isso
- commits curtos e versionados de forma progressiva
