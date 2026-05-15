# PADROES

Este documento define padroes de engenharia do ERP. Ele nao descreve arquitetura completa, endpoints ou runbooks.

## Organizacao

- `README.md`: entrada curta do projeto.
- `README_EN.md` e `README_PT.md`: visao detalhada em cada idioma.
- `docs/ARQUITETURA.md`: arquitetura e fronteiras.
- `docs/API.md`: convencoes HTTP e indice de API.
- `docs/SERVICOS.md`: ownership de servicos.
- `docs/CONTRATOS.md`: governanca contratual.
- `docs/INTEGRACOES.md`: integracoes e providers.
- `docs/OPERACOES.md`: execucao, validacao e troubleshooting.
- `docs/contracts/`: contratos versionados.
- `service-api/`: implementacao backend.
- `infra/`: runtime e observabilidade.
- `scripts/`: comandos oficiais.

## Padroes Globais

- Tenant explicito em operacoes tenant-aware.
- `publicId`, `slug` ou `key` para referencia publica de recurso.
- Erro publico com `code`, `message` e `details` opcional.
- Health com live, ready e details quando o servico participa do runtime HTTP.
- Idempotencia para mutacao sensivel.
- `202 Accepted` para operacao longa.
- Cursor pagination para listagem de volume.
- Bulk com partial success.
- Adapter/capability registry para provider externo.
- Migration por contexto persistente.
- Smoke para fluxo cross-service relevante.

## Padroes HTTP

- `GET` nao deve produzir efeito colateral.
- `POST` cria recurso, agenda trabalho ou executa comando.
- `PUT` substitui/upserta recurso identificado.
- `PATCH` altera parte do recurso ou transiciona status.
- `DELETE` so deve existir quando houver semantica clara de remocao.
- Mutacoes devem registrar ator, tenant e correlation id quando aplicavel.
- Endpoints publicos precisam estar no OpenAPI correspondente.

## Padroes de Erro

Formato recomendado:

```json
{
  "code": "RESOURCE_NOT_FOUND",
  "message": "Resource was not found.",
  "details": {
    "publicId": "..."
  }
}
```

Regras:

- `code` estavel e legivel por maquina;
- `message` curta e operacional;
- `details` sem segredo;
- status HTTP coerente com o problema.

## Padroes de Persistencia

- Cada contexto tem migrations em `service-api/service-postgresql/<contexto>/migrations`.
- Seeds sao para bootstrap, dados de referencia ou smoke.
- Servico nao deve escrever diretamente tabela de outro contexto.
- Historico/auditoria deve existir para transicoes relevantes.
- Dados externos devem guardar provider, external id e erro normalizado quando aplicavel.

## Padroes de Contrato

- OpenAPI muda junto da implementacao.
- Event schema muda junto do produtor/consumidor.
- Registry muda quando artefato contratual novo aparece.
- Breaking change exige decisao explicita.
- Compatibilidade vale mais que conveniencia local.

## Padroes de Teste

Escolha cobertura pelo risco:

- unidade para regra local;
- contrato para superficie publica;
- integracao para dependencia real;
- smoke para fluxo cross-service;
- hardening para readiness, provider e governanca;
- performance para carga/capacidade;
- backup-restore para confianca operacional de banco.

Comandos:

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

## Padroes Por Stack

### .NET

Servicos: `identity`, `billing`, `finance`.

- Separar API, application/domain e infrastructure quando o modulo ja seguir esse desenho.
- Manter DTO publico separado de entidade interna.
- Cobrir regra financeira/identidade com teste unitario.

### Go

Servicos: `edge`, `crm`, `sales`, `documents`, `rentals`.

- Handlers finos.
- Regras em pacotes internos apropriados.
- DTOs publicos claros.
- Testes para handler/regra quando alterar comportamento.

### Python

Servicos: `analytics`, `catalog`, `simulation`, `platform-control`, `support`, `supplier`, `notification`, `fiscal`.

- FastAPI com rotas claras.
- Modelos publicos separados de estrutura interna quando necessario.
- Repository driver explicito (`memory`, PostgreSQL ou equivalente).
- Pytest para regra e rota critica.

### TypeScript

Servicos: `workflow-control`, `engagement`.

- Separar runtime HTTP, dominio e adapters.
- Manter tipos publicos estaveis.
- Testes unitarios e de contrato via npm scripts do servico.

### Elixir

Servico: `workflow-runtime`.

- Modelar execucao como estado duravel.
- Registrar timeline e actions.
- Tratar retry, wait e compensacao como comportamento de dominio.

### Rust

Servico: `webhook-hub`.

- Tratar idempotencia e DLQ como regra central.
- Normalizar erros de provider.
- Manter tipos de payload explicitos.

## Padroes de Documentacao

- Cada arquivo deve ficar no proprio assunto.
- Evite repetir catalogo de endpoint fora de `docs/API.md` ou OpenAPI.
- Evite repetir inventario de servico fora de README e `docs/SERVICOS.md`.
- Atualize docs junto da mudanca que altera contrato, runtime ou arquitetura.
- Changelog registra evolucao cronologica; nao deve virar documentacao de uso.

## Checklist de PR

- O servico dono esta correto?
- O contrato foi atualizado?
- A mudanca preserva compatibilidade?
- A operacao sensivel e idempotente?
- Existe observabilidade minima?
- Existe teste proporcional ao risco?
- A documentacao alterada esta no arquivo certo?

## Checklist Para Novo Servico

Este conteudo substitui checklists soltos que antes ficavam fora do fluxo principal de documentacao; agora o criterio oficial fica aqui.

- A responsabilidade funcional nao pertence claramente a servico existente.
- O runtime define `ERP_ENV`, `ERP_AUTH_ENFORCEMENT`, `ERP_JWT_HS256_SECRET`, `ERP_INTERNAL_SERVICE_TOKEN`, `ERP_OPENFGA_ENFORCEMENT` e `OPENFGA_STORE_ID`.
- Auth middleware valida JWT/service account e chama OpenFGA quando habilitado.
- Tenant resolver e erro `tenant_slug_required` quando aplicavel.
- Mutacoes exigem `X-Correlation-Id` e idempotencia quando sensiveis.
- Erros publicos usam `code`, `message` e `details` sem segredo.
- `/health/live`, `/health/ready` e `/health/details` existem quando o servico participa do runtime.
- OpenAPI registra security, erros, idempotencia e rotas internas quando houver.
- Eventos usam envelope padrao e schema versionado.
- Logs nao imprimem senha, token, access link, documento fiscal ou provider secret.
- Testes cobrem unauthorized, forbidden, tenant mismatch, idempotencia e contrato publico.
- Dockerfile usa imagem versionada e compose recebe env comum de seguranca.
- Hooks de logs, metricas e traces propagam `X-Correlation-Id` e `traceparent`.

## Padroes de Nome

Use nomes estaveis e explicitos:

- `publicId` para identificador publico gerado;
- `tenantSlug` para tenant legivel;
- `provider` para identificador de fornecedor externo;
- `externalId` para id vindo de sistema externo;
- `status` para estado de lifecycle;
- `occurredAt` para momento do fato;
- `createdAt` e `updatedAt` para auditoria basica.

Evite:

- abreviacao obscura;
- nome que mistura dominio e apresentacao;
- campo publico com nome de coluna interna;
- enum sem fallback operacional.

## Padroes de Lifecycle

Recursos com transicao devem ter estados previsiveis. Exemplos:

```text
queued -> running -> completed
queued -> running -> failed
queued -> cancelled
running -> failed
running -> completed
```

Transicao invalida deve retornar erro claro, preferencialmente `409 Conflict` quando o problema for estado atual.

## Padroes de Bulk

Operacao bulk deve evitar "tudo ou nada" quando o dominio aceitar sucesso parcial.

Formato recomendado:

```json
{
  "results": [],
  "errors": [],
  "summary": {
    "requested": 10,
    "succeeded": 8,
    "failed": 2
  }
}
```

Cada erro deve apontar o item afetado e um codigo estavel.

## Padroes de Provider

Provider externo deve ser tratado como dependencia falivel.

Obrigatorio quando aplicavel:

- capability catalog;
- modo `configured`, `fallback`, `manual`, `disabled` ou `unconfigured`;
- timeout;
- retry policy quando seguro;
- erro normalizado;
- log sem segredo;
- fallback local apenas quando documentado.

## Padroes de Read Model

Read model deve deixar claro que e leitura derivada. Ele pode otimizar operacao, mas nao deve virar fonte transacional.

Bom uso:

- dashboard executivo;
- readiness;
- governanca contratual;
- hardening;
- compliance overview;
- go-live overview.

Mau uso:

- atualizar assinatura a partir de report;
- corrigir invoice dentro de analytics;
- aplicar entitlement a partir de edge.

## Padroes de Comentario

Comente quando o codigo expressa decisao que nao e obvia:

- workaround temporario;
- compatibilidade;
- regra fiscal ou financeira;
- retry/idempotencia;
- tradeoff de performance;
- formato exigido por provider externo.

Nao comente o obvio. Comentario ruim envelhece mais rapido que codigo.

## Checklist Por Tipo de Mudanca

### Nova rota HTTP

- pertence ao servico certo;
- OpenAPI atualizado;
- request/response com shape claro;
- erro publico definido;
- teste unitario ou de handler;
- contract validation;
- doc especifica atualizada se necessario.

### Nova tabela ou migration

- contexto correto;
- migration ordenada;
- seed apenas se necessario;
- rollback manual entendido;
- backup/restore considerado se mudar dado sensivel.

### Novo provider

- capability declarada;
- configuracao sem segredo no repositorio;
- fallback ou modo unconfigured;
- timeout e erro normalizado;
- teste sem chamar provider real por padrao.

### Novo fluxo cross-service

- dono de cada etapa definido;
- contrato entre etapas;
- correlation id;
- smoke ou teste integrado;
- observabilidade minima.

### Mudanca de documentacao

- arquivo certo;
- sem duplicar OpenAPI inteiro;
- sem prometer feature inexistente;
- links internos validos;
- changelog apenas quando registrar evolucao concluida.
