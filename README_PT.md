# ERP

ERP e uma plataforma backend-first, multi-tenant e poliglota para dominios de ERP. O projeto e organizado por servicos com ownership claro, contratos versionados, execucao local em containers e validacao automatizada.

Este README traz uma visao completa, mas compacta. Os detalhes especificos ficam em `docs/`.

## Mapa de Documentacao

| Arquivo | Finalidade |
|---------|------------|
| `docs/ARQUITETURA.md` | arquitetura, fronteiras, ownership de dados e topologia de runtime |
| `docs/API.md` | convencoes HTTP, indice de endpoints e regras de uso da API |
| `docs/SERVICOS.md` | inventario de servicos, ownership e caminhos de implementacao |
| `docs/CONTRATOS.md` | OpenAPI, schemas de eventos, registry e politica de compatibilidade |
| `docs/INTEGRACOES.md` | providers, webhooks, eventos e integracao cross-context |
| `docs/OPERACOES.md` | runtime local, scripts, banco, validacao e troubleshooting |
| `docs/PADROES.md` | padroes de engenharia para backend, testes e documentacao |
| `docs/CHANGELOG.md` | historico cronologico de mudancas |

## Estado Atual

- 24 servicos HTTP com contratos OpenAPI.
- 542 endpoints HTTP versionados.
- 15 schemas de evento versionados.
- Contratos em `docs/contracts/`.
- Runtime em `infra/docker-compose.yml`.
- Entrada operacional em `./scripts/build.sh`.
- Entrada de validacao em `./scripts/test.sh`.
- Console tecnico da API em `client-web/client-api`.

## Resumo Arquitetural

O sistema e dividido por fronteiras de servico, nao por uma unica aplicacao monolitica. Cada servico possui seu contrato publico, sua implementacao e seu contexto de persistencia quando aplicavel.

Regras principais:

- contexto de tenant deve ser explicito em operacoes tenant-aware;
- contratos publicos vivem em `docs/contracts/`;
- implementacao vive em `service-api/`;
- infraestrutura de runtime vive em `infra/`;
- ownership de banco e separado por schema/contexto PostgreSQL;
- agregacao cross-service pertence a `analytics` ou `edge`;
- callbacks externos e entrega de webhooks passam por `webhook-hub` ou adapter de provider;
- operacoes longas devem expor job, rollout ou execution em vez de bloquear a requisicao.

## Runtime Local

```bash
./scripts/build.sh
./scripts/build.sh up
./scripts/build.sh ps
./scripts/build.sh logs edge
./scripts/build.sh down
```

Banco:

```bash
./scripts/build.sh migrate all
./scripts/build.sh seed all
./scripts/build.sh backup /tmp/erp-local-backup.sql
./scripts/build.sh restore /tmp/erp-local-backup.sql
./scripts/build.sh psql
```

Validacao:

```bash
./scripts/test.sh unit
./scripts/test.sh integration
./scripts/test.sh contract
./scripts/test.sh platform
./scripts/test.sh smoke
./scripts/test.sh performance
./scripts/test.sh backup-restore
./scripts/test.sh hardening
```

## Console da API

`client-web/client-api` e o console tecnico para testar e navegar a API backend. Ele fica separado de qualquer futuro frontend empresarial.

Ele oferece:

- overview da plataforma;
- catalogo de endpoints gerado pelos arquivos OpenAPI;
- construtor de requisicoes com headers, params e JSON body;
- leitor de documentacao baseado nos markdowns reais do repositorio;
- visoes de contratos, ambientes, jornadas e operacoes.

```bash
cd client-web/client-api
npm install
npm run generate
npm run dev
```

O console usa proxy local do Vite para chamar os servicos backend quando o stack esta ativo.

## Inventario de Servicos

| Servico | Stack | Caminho | Responsabilidade |
|---------|-------|---------|------------------|
| `accounting` | Python | `service-api/service-python/accounting` | contabilidade gerencial, centros de custo, regras de posting, razao, DRE, balanco e fechamento |
| `analytics` | Python | `service-api/service-python/analytics` | relatorios executivos, governanca, readiness e leituras operacionais |
| `banking` | Python | `service-api/service-python/banking` | CNAB, boletos, extratos, conciliacao, Pix cobranca/devolucao/webhooks e Open Finance |
| `billing` | .NET | `service-api/service-csharp/billing` | planos, assinaturas, invoices, pricing por uso e tentativas de pagamento |
| `catalog` | Python | `service-api/service-python/catalog` | categorias, itens, historico de versoes, bulk e contratos de consumo |
| `crm` | Go | `service-api/service-golang/crm` | leads, customers, pipeline e enriquecimento de CNPJ |
| `documents` | Go | `service-api/service-golang/documents` | anexos, storage, assinatura e historico de versoes |
| `edge` | Go | `service-api/service-golang/edge` | entrada publica e cockpits cross-service |
| `engagement` | TypeScript | `service-api/service-typescript/engagement` | providers, eventos inbound, touchpoints e conversas |
| `finance` | .NET | `service-api/service-csharp/finance` | recebiveis, projecoes, bloqueios de comissao e atividade financeira |
| `fiscal` | Python | `service-api/service-python/fiscal` | perfil fiscal, emissao, certificados, contingencia, SPED, documentos fiscais, privacidade e auditoria |
| `identity` | .NET | `service-api/service-csharp/identity` | tenants, usuarios, sessoes, convites, roles e MFA |
| `inventory` | Python | `service-api/service-python/inventory` | saldos por local, movimentos, reservas, custo medio/FIFO e contagem ciclica |
| `notification` | Python | `service-api/service-python/notification` | preferencias, central de notificacoes e estado de entrega |
| `platform-control` | Python | `service-api/service-python/platform-control` | capabilities, providers, entitlements, quotas, lifecycle e go-live |
| `procurement` | Python | `service-api/service-python/procurement` | requisicoes, cotacoes, pedidos de compra, aprovacoes, recebimento e 3-way matching real |
| `rentals` | Go | `service-api/service-golang/rentals` | contratos recorrentes e ciclo de cobrancas |
| `sales` | Go | `service-api/service-golang/sales` | oportunidades, propostas, vendas, invoices e comissoes |
| `simulation` | Python | `service-api/service-python/simulation` | cenarios e benchmarks de carga |
| `supplier` | Python | `service-api/service-python/supplier` | categorias e diretorio de fornecedores |
| `support` | Python | `service-api/service-python/support` | filas, casos, SLA, comentarios e resumo de atendimento |
| `webhook-hub` | Rust | `service-api/service-rust/webhook-hub` | webhooks inbound/outbound, delivery log e DLQ |
| `workflow-control` | TypeScript | `service-api/service-typescript/workflow-control` | definicoes de workflow, catalogos e estado de controle |
| `workflow-runtime` | Elixir | `service-api/service-elixir/workflow-runtime` | execucoes, actions, transicoes, retries e compensacoes |

## Catalogo de Contratos

Contratos sao artefatos de engenharia versionados:

```text
docs/contracts/http/              arquivos OpenAPI
docs/contracts/events/            schemas JSON de eventos
docs/contracts/registry.json      registry de contratos
docs/contracts/schema-registry.json
docs/contracts/portal/index.html
```

Execute `./scripts/test.sh contract` antes de alterar a superficie publica da API.

## Regras de Desenvolvimento

- Mantenha comportamento de servico dentro do dono do dominio.
- Atualize OpenAPI ao mudar shape HTTP publico.
- Atualize schemas de eventos ao publicar ou consumir evento compartilhado.
- Mantenha docs com escopo claro: arquitetura em arquitetura, API em API, operacao em operacao.
- Prefira fluxos pequenos e explicitos a acoplamento invisivel.
- Adicione testes proporcionais ao impacto da mudanca.

## Fluxos Principais

Fluxo comercial:

`crm` cria demanda qualificada, `sales` conduz propostas e vendas, `billing` cria obrigacoes recorrentes e `finance` consolida projecoes, bloqueios e atividade.

Fluxo de contrato recorrente:

`rentals` e `billing` representam obrigacoes recorrentes enquanto `finance` le as consequencias financeiras.

Fluxo de automacao:

`workflow-control` define automacoes, `workflow-runtime` executa, e os servicos de dominio expoem os recursos operados pelos fluxos.

Fluxo de integracao:

Callbacks de provider entram por endpoints dedicados ou por `webhook-hub`, sao normalizados e podem ser observados por `analytics` e `edge`.

Fluxo de go-live:

`platform-control` controla lifecycle, quotas, bloqueios e rollout; `analytics` e `edge` expoem visoes executivas.

## Estrutura do Repositorio

```text
client-web/client-api/     console tecnico da API
docs/                      documentacao
docs/contracts/            contratos e schemas
infra/                     Docker Compose e assets de runtime
scripts/                   entradas de runtime/build e validacao
service-api/               servicos backend e contextos PostgreSQL
```

## O Que Este Repositorio E

Este repositorio e a camada backend e tecnica da plataforma ERP. Ele nao e um site institucional, nao e uma API CRUD simples e nao e uma aplicacao frontend-first. A superficie principal hoje e a API, seus contratos, suas fronteiras de servico e o runtime local usado para validar a plataforma.

O projeto e pesado em backend de proposito, porque a parte dificil modelada aqui nao e tela. A parte dificil e manter varios contextos de negocio consistentes o bastante para evoluir:

- identidade e tenancy;
- operacao comercial;
- contratos recorrentes e billing;
- financeiro e comissoes;
- metadata documental e assinatura;
- fiscal, privacidade e auditoria;
- definicao e execucao de workflows;
- callbacks de provider e webhooks;
- entitlements de plataforma e go-live;
- analytics e controle operacional.

## O Que Ainda Nao Deve Ser Vendido Como Final

A existencia de um servico ou endpoint nao significa que toda preocupacao de producao esta fechada. Algumas areas ja estao estruturadas antes de serem conectadas a providers externos reais.

Exemplos de preocupacoes ainda naturais de produto:

- autenticacao e autorizacao mais fortes em todas as rotas publicas;
- providers externos reais de pagamento, fiscal, comunicacao e assinatura;
- tracing distribuido entre todos os servicos;
- tratamento mais rigoroso de segredo em todo caminho de integracao;
- manifests de deploy de producao alem do Docker Compose local;
- frontend empresarial separado do console tecnico da API.

Essa distincao importa: a plataforma esta estruturalmente avancada, mas a documentacao nao deve fingir que fallback local e a mesma coisa que integracao de producao.

## Portoes de Qualidade

O projeto tem varios niveis de validacao. Use a menor suite que prova a mudanca e amplie quando o impacto crescer.

| Mudanca | Validacao recomendada |
|---------|-----------------------|
| regra local de dominio | `./scripts/test.sh unit` |
| shape HTTP publico | `./scripts/test.sh contract` |
| comportamento cross-service | `./scripts/test.sh smoke` |
| provider/readiness/go-live | `./scripts/test.sh hardening` |
| infraestrutura/runtime | `./scripts/test.sh platform` |
| risco de migration/banco | `./scripts/test.sh backup-restore` |

Para o console da API:

```bash
cd client-web/client-api
npm run generate
npm run typecheck
npm run build
```

## Principios de Documentacao

A documentacao e separada por responsabilidade:

- arquitetura explica fronteiras e decisoes;
- API explica uso HTTP e convencoes;
- servicos explicam ownership;
- contratos explicam compatibilidade;
- integracoes explicam fluxos cross-context e providers;
- operacoes explicam como rodar e diagnosticar;
- padroes explicam regras de engenharia.

Evite adicionar o mesmo catalogo de endpoints em todo arquivo. Se uma mudanca afeta endpoint, atualize o OpenAPI e a documentacao focada que explica por que aquele endpoint existe.

## Fluxo Tipico de Contribuicao

1. Identifique o servico dono.
2. Confira o contrato OpenAPI.
3. Altere implementacao e testes juntos.
4. Atualize migrations/seeds se persistencia mudou.
5. Regenere o catalogo do console se contratos HTTP mudaram.
6. Rode a suite de validacao proporcional ao risco.
7. Atualize o arquivo de documentacao correto.
8. Adicione changelog apenas quando a mudanca estiver pronta para ser registrada como progresso.

## Mantenedor

Thiago Di Faria - thiagodifaria@gmail.com
