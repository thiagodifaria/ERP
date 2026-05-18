# padrões

Este documento define padrões de engenharia do projeto. Ele não descreve arquitetura completa, endpoints ou runbooks.

## Organização

- `README.md`: entrada curta do projeto.
- `README_EN.md` e `README_PT.md`: visão detalhada em cada idioma.
- `docs/ARQUITETURA.md`: arquitetura e fronteiras.
- `docs/API.md`: convenções HTTP e índice de API.
- `docs/SERVICOS.md`: ownership de serviços.
- `docs/CONTRATOS.md`: governança contratual.
- `docs/INTEGRações.md`: integrações e providers.
- `docs/OPERações.md`: execução, validação e troubleshooting.
- `docs/contracts/`: contratos versionados.
- `service-api/`: implementação backend.
- `infra/`: runtime e observabilidade.
- `scripts/`: comandos oficiais.

## padrões Globais

- Tenant explícito em operações tenant-aware.
- `publicId`, `slug` ou `key` para referência pública de recurso.
- Erro público com `code`, `message` e `details` opcional.
- Health com live, ready e details quando o serviço participa do runtime HTTP.
- idempotência para mutação sensível.
- `202 Accepted` para operação longa.
- Cursor pagination para listagem de volume.
- Bulk com partial success.
- Adapter/capability registry para provider externo.
- Migration por contexto persistente.
- Smoke para fluxo cross-service relevante.

## padrões HTTP

- `GET` não deve produzir efeito colateral.
- `POST` cria recurso, agenda trabalho ou executa comando.
- `PUT` substitui/upserta recurso identificado.
- `PATCH` altera parte do recurso ou transiciona status.
- `DELETE` só deve existir quando houver semantica clara de remocao.
- Mutações devem registrar ator, tenant e correlation id quando aplicável.
- Endpoints públicos precisam estar no OpenAPI correspondente.

## padrões de Erro

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

## padrões de persistência

- Cada contexto tem migrations em `service-api/service-postgresql/<contexto>/migrations`.
- Seeds são para bootstrap, dados de referência ou smoke.
- serviço não deve escrever diretamente tabela de outro contexto.
- histórico/auditoria deve existir para transições relevantes.
- Dados externos devem guardar provider, external id e erro normalizado quando aplicável.

## padrões de Contrato

- OpenAPI muda junto da implementação.
- Event schema muda junto do produtor/consumidor.
- Registry muda quando artefato contratual novo aparece.
- Breaking change exige decisão explicita.
- Compatibilidade vale mais que conveniencia local.

## padrões de Teste

Escolha cobertura pelo risco:

- unidade para regra local;
- contrato para superfície pública;
- integração para dependência real;
- smoke para fluxo cross-service;
- hardening para readiness, provider e governança;
- performance para carga/capacidade;
- backup-restore para confiança operacional de banco.

Comandos:

```bash
./scripts/test.sh unit
./scripts/test.sh contract
./scripts/test.sh smoke
```

## padrões Por Stack

### .NET

serviços: `identity`, `billing`, `finance`.

- Separar API, application/domain e infrastructure quando o módulo já seguir esse desenho.
- Manter DTO público separado de entidade interna.
- Cobrir regra financeira/identidade com teste unitario.

### Go

serviços: `edge`, `crm`, `sales`, `documents`, `rentals`.

- Handlers finos.
- Regras em pacotes internos apropriados.
- DTOs públicos claros.
- Testes para handler/regra quando alterar comportamento.

### Python

serviços: `ai-governance`, `analytics`, `catalog`, `search`, `simulation`, `platform-control`, `support`, `supplier`, `notification`, `fiscal`.

Na `v1.5.0`, comandos críticos, eventos operacionais, fechamento financeiro, mudanças breaking, chamadas reais a provider externo, sinais externos de risco/verificação e mudanças de hardening devem passar por policy decision, approval quando aplicável, timeline, evidence vault, event mesh, snapshot hash, contract evolution, provider activation ou static policy antes de serem tratados como finalizados.

- FastAPI com rotas claras.
- Modelos públicos separados de estrutura interna quando necessário.
- Repository driver explícito (`memory`, PostgreSQL ou equivalente).
- Pytest para regra e rota critica.

### TypeScript

serviços: `workflow-control`, `engagement`.

- Separar runtime HTTP, domínio e adapters.
- Manter tipos públicos estaveis.
- Testes unitários e de contrato via npm scripts do serviço.

### Elixir

serviço: `workflow-runtime`.

- Modelar execução como estado durável.
- Registrar timeline e actions.
- Tratar retry, wait e compensação como comportamento de domínio.

### Rust

serviço: `webhook-hub`.

- Tratar idempotência e DLQ como regra central.
- Normalizar erros de provider.
- Manter tipos de payload explícitos.

## padrões de documentação

- Cada arquivo deve ficar no próprio assunto.
- Evite repetir catálogo de endpoint fora de `docs/API.md` ou OpenAPI.
- Evite repetir inventário de serviço fora de README e `docs/SERVICOS.md`.
- Atualize docs junto da mudança que altera contrato, runtime ou arquitetura.
- Changelog registra evolucao cronologica; não deve virar documentação de uso.

## Controle Interno De mudança

- O serviço dono está correto?
- O contrato foi atualizado?
- A mudança preserva compatibilidade?
- A operação sensível e idempotente?
- Existe observabilidade minima?
- Existe teste proporcional ao risco?
- A documentação alterada está no arquivo certo?

## Checklist Para Novo serviço

Este conteudo substitui checklists soltos que antes ficavam fora do fluxo principal de documentação; agora o criterio oficial fica aqui.

- A responsabilidade funcional não pertence claramente a serviço existente.
- O runtime define `ERP_ENV`, `ERP_AUTH_ENFORCEMENT`, `ERP_JWT_HS256_SECRET`, `ERP_INTERNAL_SERVICE_TOKEN`, `ERP_OPENFGA_ENFORCEMENT` e `OPENFGA_STORE_ID`.
- Auth middleware valida JWT/service account e chama OpenFGA quando habilitado.
- Tenant resolver e erro `tenant_slug_required` quando aplicável.
- Mutações exigem `X-Correlation-Id` e idempotência quando sensíveis.
- Erros públicos usam `code`, `message` e `details` sem segredo.
- `/health/live`, `/health/ready` e `/health/details` existem quando o serviço participa do runtime.
- OpenAPI registra security, erros, idempotência e rotas internas quando houver.
- Eventos usam envelope padrão e schema versionado.
- Logs não imprimem senha, token, access link, documento fiscal ou provider secret.
- Testes cobrem unauthorized, forbidden, tenant mismatch, idempotência e contrato público.
- Dockerfile usa imagem versionada e compose recebe env comum de segurança.
- Hooks de logs, métricas e traces propagam `X-Correlation-Id` e `traceparent`.
- repositórios relacionais nunca calculam identificador por `MAX(id)+1`; use `BIGSERIAL`, identity column, sequence ou `RETURNING id`.
- migrations de um mesmo domínio nunca repetem prefixo numérico; a suite `security` bloqueia duplicidade antes de runtime.
- tabela operacional nova precisa declarar `tenant_id` ou justificativa no manifesto de contrato de tenant.
- UI técnica nunca copia `Authorization: Bearer` real por padrão; cURL e histórico devem redigir segredo salvo ação explicita.
- Imagem de runtime e dependência de infra devem ter versão explicita; `latest` não e artefato de deploy.

## Conformance De Plataforma

A linha `1.4.x` exige que todo serviço passe por conformance minima antes de ser tratado como produtivo:

- Auth: JWT/service account, tenant explícito, actor quando houver usuário e OpenFGA quando habilitado.
- observabilidade: `X-Correlation-Id`, `traceparent`, logs sem segredo, health endpoints e erro público consistente.
- Dados: identificadores por sequence/identity, constraints de unicidade e FK quando houver relação transacional.
- Contratos: OpenAPI atualizado, security schemes declarados e catálogo do `client-api` regenerado.
- Infra: imagem versionada, env comum, probes, resources e documentação de cobertura Kubernetes/Compose.
- segurança: secrets locais bloqueados fora de local/test e provider BYOK nunca mascarado como real sem credencial.
- modularidade: health, autenticação, validação de payload e clients externos devem ficar fora do roteador principal sempre que o serviço passar de uma superfície pequena.

## padrões de Nome

Use nomes estaveis e explícitos:

- `publicId` para identificador público gerado;
- `tenantSlug` para tenant legivel;
- `provider` para identificador de fornecedor externo;
- `externalId` para id vindo de sistema externo;
- `status` para estado de lifecycle;
- `occurredAt` para momento do fato;
- `createdAt` e `updatedAt` para auditoria basica.

Evite:

- abreviação obscura;
- nome que mistura domínio e apresentação;
- campo público com nome de coluna interna;
- enum sem fallback operacional.

## padrões de Lifecycle

Recursos com transicao devem ter estados previsíveis. Exemplos:

```text
queued -> running -> completed
queued -> running -> failed
queued -> cancelled
running -> failed
running -> completed
```

Transicao invalida deve retornar erro claro, preferêncialmente `409 Conflict` quando o problema for estado atual.

## padrões de Bulk

operação bulk deve evitar "tudo ou nada" quando o domínio aceitar sucesso parcial.

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

Cada erro deve apontar o item afetado e um código estavel.

## padrões de Provider

Provider externo deve ser tratado como dependência falivel.

obrigatório quando aplicável:

- capability catalog;
- modo `configured`, `fallback`, `manual`, `disabled` ou `unconfigured`;
- timeout;
- retry policy quando seguro;
- erro normalizado;
- log sem segredo;
- fallback local apenas quando documentado.

## padrões de Read Model

Read model deve deixar claro que e leitura derivada. Ele pode otimizar operação, mas não deve virar fonte transacional.

Bom uso:

- dashboard executivo;
- readiness;
- governança contratual;
- hardening;
- compliance overview;
- go-live overview.

Mau uso:

- atualizar assinatura a partir de report;
- corrigir invoice dentro de analytics;
- aplicar entitlement a partir de edge.

## padrões de Comentario

Comente quando o código expressa decisão que não e obvia:

- workaround temporario;
- compatibilidade;
- regra fiscal ou financeira;
- retry/idempotência;
- tradeoff de performance;
- formato exigido por provider externo.

não comente o obvio. Comentario ruim envelhece mais rápido que código.

## Checklist Por Tipo de mudança

### Nova rota HTTP

- pertence ao serviço certo;
- OpenAPI atualizado;
- request/response com shape claro;
- erro público definido;
- teste unitario ou de handler;
- contract validation;
- doc especifica atualizada se necessário.

### Nova tabela ou migration

- contexto correto;
- migration ordenada;
- seed apenas se necessário;
- rollback manual entendido;
- backup/restore considerado se mudar dado sensível.

### Novo provider

- capability declarada;
- configuração sem segredo no repositório;
- fallback ou modo unconfigured;
- timeout e erro normalizado;
- teste sem chamar provider real por padrão.

### Novo fluxo cross-service

- dono de cada etapa definido;
- contrato entre etapas;
- correlation id;
- smoke ou teste integrado;
- observabilidade minima.

### mudança de documentação

- arquivo certo;
- sem duplicar OpenAPI inteiro;
- sem prometer feature inexistente;
- links internos validos;
- changelog apenas quando registrar evolucao concluida.
