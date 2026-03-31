# ARQUITETURA

## Objetivo

Registrar a arquitetura-alvo do ERP de forma curta, objetiva e versionavel.

## Visao geral

O ERP esta sendo construido como uma plataforma empresarial modular, multi-tenant e poliglota.

Fluxo macro:

1. a requisicao externa entra pelo gateway
2. o gateway encaminha para `edge`
3. `edge` autentica, autoriza, correlaciona e roteia
4. o servico dono do dominio executa a regra
5. eventos relevantes sao publicados de forma assincrona
6. servicos consumidores reagem com idempotencia
7. logs, traces e metricas acompanham todos os saltos

## Planos arquiteturais

- transaction plane para verdade operacional
- control plane para definicao de automacoes
- runtime plane para execucao duravel
- analytics plane para leitura pesada e projecoes

## Servicos por stack

- Go: `edge`, `crm`, `sales`, `rentals`, `documents`
- C#: `identity`, `finance`, `billing`
- TypeScript: `workflow-control`, `engagement`
- Elixir: `workflow-runtime`
- Python: `analytics`, `simulation`
- Rust: `webhook-hub`
- PostgreSQL: ownership de dados por dominio

## Fundacoes de plataforma

- Keycloak para autenticacao
- OpenFGA para autorizacao fina
- Kafka para backbone de eventos
- Redis para cache, rate limiting e locks leves
- PostgreSQL por contexto de dominio
- MinIO ou S3 para documentos
- OpenSearch para busca operacional
- Loki, Tempo, Prometheus e Grafana para observabilidade
- Temporal para fluxos duraveis e orquestracoes longas

## Regras arquiteturais

- bootstrap fino, sem regra de negocio
- dominio e aplicacao concentram as regras importantes
- adapters isolam integracoes externas
- cada contexto e dono dos seus dados e migrations
- nenhum servico de negocio recebe webhook critico diretamente
