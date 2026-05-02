# CONTRATOS

## Objetivo

Centralizar os artefatos de integracao do backend em um lugar versionado e previsivel.

## Estrutura

- `contracts/http/`: specs OpenAPI por servico
- `contracts/events/`: schemas JSON para eventos e payloads compartilhados

## Regras operacionais

- endpoint publico novo deve nascer com artefato correspondente quando ele fizer parte da superficie de integracao
- payload de evento entre servicos ou com providers externos deve ter schema versionado
- breaking change em contrato exige update explicito de changelog e revisao de compatibilidade
- operacoes sensiveis devem aceitar `Idempotency-Key` por header ou contrato equivalente quando isso fizer sentido
- operacoes longas devem convergir para padrao `202 Accepted`, polling e callback conforme a maturidade do servico
- listagens novas devem preferir cursor pagination quando o volume esperado justificar

## Proximo passo natural

Os artefatos daqui servem de base para uma UI central de navegacao e teste de API, sem depender apenas de README ou memoria operacional.
