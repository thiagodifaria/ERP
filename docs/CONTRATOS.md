# CONTRATOS

## Objetivo

Centralizar os artefatos de integracao do backend em um lugar versionado e previsivel.

## Estrutura

- `contracts/http/`: specs OpenAPI por servico
- `contracts/events/`: schemas JSON para eventos e payloads compartilhados
- `contracts/registry.json`: indice simples para agregacao automatica
- `contracts/schema-registry.json`: registry materializado dos schemas de evento
- `contracts/portal/index.html`: baseline navegavel da superficie HTTP e dos schemas
- `docs/VERSIONAMENTO_CONTRATOS.md`: baseline de compatibilidade e versionamento

## Baseline atual

- specs HTTP publicos ja versionados para `identity`, `crm`, `sales`, `finance`, `billing`, `documents`, `rentals`, `engagement`, `catalog`, `platform-control`, `analytics`, `simulation`, `edge`, `webhook-hub`, `workflow-control` e `workflow-runtime`
- schemas de evento para trilhas de lifecycle de tenant, quotas, assinatura documental, enriquecimento de CNPJ, publicacao de item de catalogo e entrega outbound de webhook
- artefatos suficientes para alimentar governanca de integracao, smoke mais forte e uma futura UI agregada de navegacao e teste

## Regras operacionais

- endpoint publico novo deve nascer com artefato correspondente quando ele fizer parte da superficie de integracao
- payload de evento entre servicos ou com providers externos deve ter schema versionado
- breaking change em contrato exige update explicito de changelog e revisao de compatibilidade
- operacoes sensiveis devem aceitar `Idempotency-Key` por header ou contrato equivalente quando isso fizer sentido
- operacoes longas devem convergir para padrao `202 Accepted`, polling e callback conforme a maturidade do servico
- listagens novas devem preferir cursor pagination quando o volume esperado justificar
- operacoes em lote devem explicitar `results`, `errors` e `summary` quando houver `partial success`
- `docs/API_PORTAL.md` registra a diretriz e `contracts/portal/index.html` materializa o baseline versionado da UI central de navegacao e teste

## Proximo passo natural

Os artefatos daqui servem de base para uma UI central de navegacao e teste de API, sem depender apenas de README ou memoria operacional.
