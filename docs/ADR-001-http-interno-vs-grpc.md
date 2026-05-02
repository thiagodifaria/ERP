# ADR-001 - HTTP interno vs gRPC

## Status

Aceita.

## Contexto

O monorepo tem servicos em Go, .NET, TypeScript, Elixir, Python e Rust. O principal objetivo nesta etapa e manter integracao simples, debuggable e facil de operar em ambiente local e em containers.

## Decisao

O backend vai manter `HTTP/JSON` como padrao interno por enquanto.

`gRPC` permanece opcao futura para trechos especificos de alta frequencia ou necessidades de streaming, mas nao e o padrao atual do projeto.

## Motivos

- integra melhor com a superficie publica que ja existe
- reduz atrito entre linguagens diferentes
- simplifica smoke, curl, troubleshooting e onboarding
- facilita a futura geracao de OpenAPI e UI agregada de contratos
- encaixa melhor no estagio atual do produto, que ainda esta consolidando dominio e contratos

## Consequencias

- contratos HTTP precisam ser tratados como artefatos versionados
- padroes de idempotencia, async e paginacao precisam ser explicitados para nao virar caos por servico
- se algum fluxo justificar gRPC depois, ele deve nascer com ADR propria e sem quebrar a malha HTTP ja existente
