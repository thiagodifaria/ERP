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

## Entrega incremental

- miniupdates pequenos e frequentes
- changelog atualizado a cada entrega relevante
- testes executados antes de subir uma entrega quando houver runtime suficiente para isso
- commits curtos e versionados de forma progressiva
