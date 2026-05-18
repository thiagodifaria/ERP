# segurança

Este documento consolida os controles de segurança verificaveis do projeto.

## API e Workload

- Todos os serviços HTTP possuem middleware de JWT HS256 ou token interno de service account.
- `ERP_AUTH_ENFORCEMENT=enforced` ativa rejeicao de requests sem credencial valida fora de rotas públicas.
- `ERP_SECURITY_LEVEL=strict` define o baseline produtivo esperado para borda, workloads, logs e integrações.
- `ERP_OPENFGA_ENFORCEMENT=true` ativa check relacional por OpenFGA.
- `ERP_REQUIRE_REQUEST_SIGNATURE=true` exige que chamadas internas sensíveis e webhooks sigam politica de assinatura/correlação definida pelo domínio.
- Mutações exigem `X-Correlation-Id` quando o enforcement está ativo.
- Gateway aplica rate limit global, por tenant, por actor e por rotas sensíveis.
- Gateway aplica headers defensivos (`Strict-Transport-Security`, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `Permissions-Policy`, CSP), limite de body e timeouts curtos.
- Kubernetes usa Pod Security `restricted`, NetworkPolicy deny-by-default e workloads sem privilege escalation.

## Contrato único De Auth e Tracing

Todo serviço novo ou existente deve tratar autenticação, autorização e rastreabilidade como contrato de plataforma, não como detalhe local do framework.

| Campo | obrigatório | Regra |
| --- | --- | --- |
| `Authorization: Bearer` | sim para rotas protegidas | JWT de usuário ou service-account token; token invalido falha em `401` |
| `ERP_JWT_AUDIENCE` | sim | audience precisa bater com a API ou gateway interno esperado |
| tenant | sim para dado multi-tenant | vem de claim, header validado ou path resolvido pelo domínio |
| actor | sim para mutação iniciada por usuário | registrado em auditoria, timeline ou evento |
| `X-Correlation-Id` | sim para mutação | preservado pelo gateway quando recebido e gerado quando ausente |
| `traceparent` | recomendado/obrigatório em runtime produtivo | propagado entre chamadas HTTP, eventos e logs estruturados |
| OpenFGA | obrigatório quando `ERP_OPENFGA_ENFORCEMENT=true` | falha fechada em erro de policy/check para operação sensível |
| idempotência | obrigatoria para mutação sensível | chave única por tentativa logica, nunca derivada apenas de timestamp |

Falhas devem ser públicas e previsíveis: `401` para credencial ausente/invalida, `403` para policy negada, `400` para correlation ausente em mutação, `409` para conflito de estado/idempotência e `429` para rate limit.

## Dados sensíveis

- Password recovery persiste hash do token e não devolve token fora de local/test.
- Access links de documents são assinados, expiraveis, revogaveis e auditados com persistência relacional de revogações e eventos.
- `DOCUMENTS_MALWARE_SCAN_MODE=required` torna varredura de upload requisito produtivo antes de disponibilizar download.
- Uploads com assinatura maliciosa conhecida são bloqueados antes de virarem attachment.
- Logs, contratos e portal não devem expor segredo, token, senha ou provider credential.
- `ERP_AUDIT_LOG_REDACTION=strict` exige redaction de token, senha, segredo, documento sensível e credential de provider em logs/auditoria.
- Providers externos seguem BYOK: chaves como `OPENAI_API_KEY`, `BILLING_STRIPE_SECRET_KEY`, `ENGAGEMENT_RESEND_API_KEY` e tokens de assinatura/documentos entram por ambiente/secret manager, nunca por código ou contrato.
- Testes de provider retornam `unavailable` quando a credencial não existe e registram timeline/evidência sem expor o valor do segredo.
- OCR, fiscal, consulta cadastral, mercado e feeds de risco podem processar dados sensíveis; cada provider deve declarar finalidade, credencial, postura de fallback e se a API é pública ou BYOK.
- Exportação, portabilidade e anonimização passam por politica de domínio, preservando evidências fiscais, financeiras e auditorias com obrigação legal de retenção.
- Novas features com dado pessoal precisam declarar finalidade, retenção e exposição em contrato/evento/log antes de entrar na superfície pública.

## Modelo De Ameacas 1.0.0

| Ameaca | Controle obrigatório | evidência |
| --- | --- | --- |
| Bypass do gateway | perfil corporate-like, serviços internos sem porta públicada e NetworkPolicy deny-by-default | `infra/docker-compose.corporate-like.yml`, `infra/kubernetes/base/network-policy.yaml` |
| BOLA/BFLA | JWT/service account, tenant explícito, scopes/capabilities e OpenFGA quando habilitado | middlewares por stack, OpenAPI security, `scripts/test.sh security` |
| Replay de webhook/callback | assinatura, janela curta de replay, idempotência e DLQ | `WEBHOOK_HUB_REQUIRE_SIGNATURE`, `WEBHOOK_HUB_REPLAY_WINDOW_SECONDS`, event envelope |
| Exfiltração por log | redaction estrita e erro público sem segredo | `ERP_AUDIT_LOG_REDACTION=strict`, padrão de erro público |
| Upload malicioso | scan obrigatório, tamanho máximo, retention e revogação | `DOCUMENTS_MALWARE_SCAN_MODE=required`, gateway body limit, documents audit |
| exposição por busca/export | redação padrão, auditoria de query, legal hold e export controlado | `search`, `/api/search/audit-events`, `/api/search/legal-holds` |
| Uso indevido de IA | allowlist de ferramentas, politica read-only, redação e trilha de run | `ai-governance`, `/api/ai-governance/policies`, `/api/ai-governance/audit-events` |
| Chamada externa sem chave | provider activation BYOK, status indisponível e auditoria sem secret | `platform-control` provider activation |
| Sinal externo usado como verdade | reports tratam news/mercado/cadastro como sinal, não como mutação de domínio | `analytics` external intelligence |
| Falha operacional sem resposta | timeline imutavel, action items, resolução e postmortem | `platform-control` incident command |
| Comando crítico sem autorização | policy decision center, approval workflow e evidence vault | `platform-control` policies, approvals e evidence |
| auditoria sem prova verificável | hash lógico de evidência, retention e classificação | `platform-control` audit evidence vault |
| Evento operacional sem rastreio | event mesh com payload hash, dead letter, replay e lineage | `platform-control` event mesh |
| Fechamento financeiro sem prova | reconciliação, snapshot hash e readiness de fechamento | `analytics` financial close |
| mudança breaking sem controle | contract evolution com diff, breaking change e aprovação | `platform-control` contract evolution |
| Abuso de auth/recovery | rate limit dedicado, token com hash, single-use e expiração | gateway, identity recovery, security audit |
| Movimento lateral | service account por workload, Pod Security restricted e egress controlado | Kubernetes manifests e secret template |
| configuração insegura | defaults locais bloqueados em produção | `.env.production.example`, `hardening-secrets` |

## Gate De segurança 1.4.6

Para considerar a postura produtiva respeitavel, rode:

```bash
./scripts/test.sh security
./scripts/test.sh hardening-secrets
./scripts/test.sh supply-chain
./scripts/test.sh production-readiness
```

O aceite exige que:

- headers defensivos estejam ativos no gateway;
- `ERP_SECURITY_LEVEL=strict` esteja definido no perfil produtivo;
- request signature e redaction estrita estejam documentadas e presentes no ambiente;
- webhooks tenham assinatura obrigatoria e janela de replay;
- Pod Security restricted e NetworkPolicy deny-by-default existam;
- secrets default não passem em ambiente não-local;
- contratos, portal e eventos não exponham segredo.
- `MAX(id)+1` não exista em código transacional;
- cURL gerado pelo console redija Bearer Token por padrão;
- imagens operacionais não dependam de tag `latest`;
- auth, tenant, actor, correlation id e `traceparent` estejam documentados como contrato transversal.

## inventário LGPD

| domínio | Dados pessoais | Sensibilidade | Finalidade |
| --- | --- | --- | --- |
| `identity` | nome, é-mail, papeis, sessões, MFA, auditoria | alta | autenticação, autorização e trilha de acesso |
| `crm` | leads, clientes, contatos, notas, anexos vinculados | alta | relacionamento comercial e funil |
| `engagement` | touchpoints, entregas, callbacks e respostas | alta | comunicação transacional/comercial |
| `support` | casos, comentários, filas e SLA | media/alta | atendimento e evidências de suporte |
| `billing`/`finance` | cobrança, recebíveis, caixa, movimentos e reconciliação | alta | controle financeiro e auditoria |
| `documents`/`fiscal` | arquivos, documentos fiscais, consentimentos e privacy requests | critica | guarda documental, compliance e retenção legal |
| `search` | índice operacional, evidências de discovery, legal holds e exports | alta | descoberta controlada e auditoria operacional |
| `ai-governance` | prompts redigidos, decisoes de ferramenta, runs e auditoria | media/alta | governança de assistentes e rastreabilidade |

## Suite De segurança

`./scripts/test.sh security` executa validações estaticas de:

- middleware JWT/OpenFGA por stack;
- rate limit do gateway;
- headers defensivos e body limit no gateway;
- hardening de secrets;
- Pod Security restricted e NetworkPolicy deny-by-default;
- recovery token com hash;
- documents access links, revogação persistente, auditoria e scan;
- assinatura/replay de webhooks;
- envelope/event registry;
- portal de integração;
- inventário LGPD;
- busca operacional, redação, e-discovery e governança de IA;
- readiness de incident command e trilha de postmortem;
- policy decisions, command approvals, runbook automation, evidence vault, event mesh, financial close e contract evolution;
- provider activation BYOK e AI Governance com fallback deterministico local;
- OCR/document intelligence, fiscal Brazil, registry enrichment, market macro risk e external risk feed;
- static policy para `MAX(id)+1`, bearer em cURL e tag `latest`;
- contrato único de auth/tracing com `traceparent`, tenant, actor e correlation id;
- risk/compliance scoring como sinal de aceite operacional;
- SLOs/runbooks em `docs/OPERações.md`;
- guardrails de poliglotismo.
