# API PORTAL

## Objetivo

Servir como base para uma UI central de navegacao, teste e descoberta da API completa do ERP.

## Fonte da documentacao

- `contracts/http/`: specs OpenAPI por servico
- `contracts/events/`: schemas JSON de eventos
- `contracts/registry.json`: indice simples para agregacao automatica

## Padrao esperado

- endpoint publico novo relevante para integracao deve atualizar a spec correspondente
- eventos compartilhados devem possuir schema versionado
- operacoes assicronas devem sinalizar `202 Accepted` e caminho de polling quando aplicavel
- operacoes sensiveis devem explicitar `Idempotency-Key`
- listagens novas devem preferir cursor pagination
- operacoes em lote devem expor contrato de `partial success`

## Proximo passo natural

Consumir o `registry.json` em uma UI agregadora com Swagger UI, Redoc, Scalar ou portal proprio.
