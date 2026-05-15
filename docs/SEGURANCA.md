# SEGURANCA

Este documento consolida os controles de seguranca verificaveis do ERP.

## API E Workload

- Todos os servicos HTTP possuem middleware de JWT HS256 ou token interno de service account.
- `ERP_AUTH_ENFORCEMENT=enforced` ativa rejeicao de requests sem credencial valida fora de rotas publicas.
- `ERP_SECURITY_LEVEL=strict` define o baseline produtivo esperado para borda, workloads, logs e integracoes.
- `ERP_OPENFGA_ENFORCEMENT=true` ativa check relacional por OpenFGA.
- `ERP_REQUIRE_REQUEST_SIGNATURE=true` exige que chamadas internas sensiveis e webhooks sigam politica de assinatura/correlacao definida pelo dominio.
- Mutacoes exigem `X-Correlation-Id` quando o enforcement esta ativo.
- Gateway aplica rate limit global, por tenant, por actor e por rotas sensiveis.
- Gateway aplica headers defensivos (`Strict-Transport-Security`, `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `Permissions-Policy`, CSP), limite de body e timeouts curtos.
- Kubernetes usa Pod Security `restricted`, NetworkPolicy deny-by-default e workloads sem privilege escalation.

## Dados Sensiveis

- Password recovery persiste hash do token e nao devolve token fora de local/test.
- Access links de documents sao assinados, expiraveis, revogaveis e auditados com persistencia relacional de revogacoes e eventos.
- `DOCUMENTS_MALWARE_SCAN_MODE=required` torna varredura de upload requisito produtivo antes de disponibilizar download.
- Uploads com assinatura maliciosa conhecida sao bloqueados antes de virarem attachment.
- Logs, contratos e portal nao devem expor segredo, token, senha ou provider credential.
- `ERP_AUDIT_LOG_REDACTION=strict` exige redaction de token, senha, segredo, documento sensivel e credential de provider em logs/auditoria.
- Exportacao, portabilidade e anonimizacao passam por politica de dominio, preservando evidencias fiscais, financeiras e auditorias com obrigacao legal de retencao.
- Novas features com dado pessoal precisam declarar finalidade, retencao e exposicao em contrato/evento/log antes de entrar na superficie publica.

## Modelo De Ameacas 1.0.0

| Ameaca | Controle obrigatorio | Evidencia |
| --- | --- | --- |
| Bypass do gateway | perfil corporate-like, servicos internos sem porta publicada e NetworkPolicy deny-by-default | `infra/docker-compose.corporate-like.yml`, `infra/kubernetes/base/network-policy.yaml` |
| BOLA/BFLA | JWT/service account, tenant explicito, scopes/capabilities e OpenFGA quando habilitado | middlewares por stack, OpenAPI security, `scripts/test.sh security` |
| Replay de webhook/callback | assinatura, janela curta de replay, idempotencia e DLQ | `WEBHOOK_HUB_REQUIRE_SIGNATURE`, `WEBHOOK_HUB_REPLAY_WINDOW_SECONDS`, event envelope |
| Exfiltracao por log | redaction estrita e erro publico sem segredo | `ERP_AUDIT_LOG_REDACTION=strict`, padrao de erro publico |
| Upload malicioso | scan obrigatorio, tamanho maximo, retention e revogacao | `DOCUMENTS_MALWARE_SCAN_MODE=required`, gateway body limit, documents audit |
| Abuso de auth/recovery | rate limit dedicado, token com hash, single-use e expiracao | gateway, identity recovery, security audit |
| Movimento lateral | service account por workload, Pod Security restricted e egress controlado | Kubernetes manifests e secret template |
| Configuracao insegura | defaults locais bloqueados em producao | `.env.production.example`, `hardening-secrets` |

## Gate De Seguranca 1.0.0

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
- secrets default nao passem em ambiente nao-local;
- contratos, portal e eventos nao exponham segredo.

## Inventario LGPD

| Dominio | Dados pessoais | Sensibilidade | Finalidade |
| --- | --- | --- | --- |
| `identity` | nome, e-mail, papeis, sessoes, MFA, auditoria | alta | autenticacao, autorizacao e trilha de acesso |
| `crm` | leads, clientes, contatos, notas, anexos vinculados | alta | relacionamento comercial e funil |
| `engagement` | touchpoints, entregas, callbacks e respostas | alta | comunicacao transacional/comercial |
| `support` | casos, comentarios, filas e SLA | media/alta | atendimento e evidencias de suporte |
| `billing`/`finance` | cobranca, recebiveis, caixa, movimentos e reconciliacao | alta | controle financeiro e auditoria |
| `documents`/`fiscal` | arquivos, documentos fiscais, consentimentos e privacy requests | critica | guarda documental, compliance e retencao legal |

## Suite De Seguranca

`./scripts/test.sh security` executa validacoes estaticas de:

- middleware JWT/OpenFGA por stack;
- rate limit do gateway;
- headers defensivos e body limit no gateway;
- hardening de secrets;
- Pod Security restricted e NetworkPolicy deny-by-default;
- recovery token com hash;
- documents access links, revogacao persistente, auditoria e scan;
- assinatura/replay de webhooks;
- envelope/event registry;
- portal de integracao;
- inventario LGPD;
- SLOs/runbooks em `docs/OPERACOES.md`;
- guardrails de poliglotismo.
