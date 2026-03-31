# PADROES

## Codigo

- codigo em ingles
- comentarios em portugues do Brasil
- nomes de pastas curtos, secos e previsiveis
- arquivos de bootstrap chamados `server.*` ou equivalente idiomatico
- regra de negocio em `domain` e `application`
- handlers, controllers e APIs finos

## Comentarios

Comentarios devem explicar:

- intencao
- limite de responsabilidade
- regra nao trivial
- motivo de uma decisao
- risco de alteracao

Comentarios nao devem explicar o obvio.

## Dados e identidade

- `id BIGINT` como identificador interno
- `public_id UUIDv7` como identificador externo
- timestamps operacionais em UTC
- `tenant_id` em todos os agregados relevantes
- eventos carregam os identificadores necessarios para rastreabilidade

## System design

- database per service
- transactional outbox para eventos relevantes
- inbox e consumo idempotente
- anti-corruption layer para APIs externas
- CQRS pragmatico quando a leitura exigir outro shape
- Temporal para fluxos longos e orquestrados

## O que evitar

- handler falando direto com SQL
- regra de negocio em bootstrap
- SDK externo dentro de dominio
- banco compartilhado como atalho entre servicos
- comentario em excesso
- criacao de pasta sem necessidade real
