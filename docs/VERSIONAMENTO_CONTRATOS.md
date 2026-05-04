# VERSIONAMENTO DE CONTRATOS

## Objetivo

Definir um baseline simples para evolucao de contratos HTTP e eventos sem depender de memoria informal.

## Regras

- mudanca publica em rota relevante deve atualizar a spec correspondente em `contracts/http/`
- evento compartilhado novo ou alterado deve atualizar o schema correspondente em `contracts/events/`
- breaking change exige registro explicito no changelog e revisao de compatibilidade
- operacoes assincronas devem preferir `202 Accepted`, polling e artefato contratual correspondente
- operacoes sensiveis com risco de duplicidade devem explicitar `Idempotency-Key`
- listagens novas de volume transacional devem preferir cursor pagination
- operacoes em lote devem documentar partial success com `results`, `errors` e `summary`

## Compatibilidade

- adicionar campo opcional: compativel
- adicionar rota nova: compativel
- mudar nome, semantica ou remover campo publico: potencialmente breaking
- mudar shape de evento consumido por outro servico: potencialmente breaking

## Baseline atual

- `contracts/registry.json` funciona como indice central minimo
- `docs/API_PORTAL.md` descreve a base para UI agregada
- `docs/ADR-001-http-interno-vs-grpc.md` registra a escolha atual do backbone interno
