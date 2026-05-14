# SEGURANCA

Este documento consolida os controles de seguranca verificaveis do ERP.

## API E Workload

- Todos os servicos HTTP possuem middleware de JWT HS256 ou token interno de service account.
- `ERP_AUTH_ENFORCEMENT=enforced` ativa rejeicao de requests sem credencial valida fora de rotas publicas.
- `ERP_OPENFGA_ENFORCEMENT=true` ativa check relacional por OpenFGA.
- Mutacoes exigem `X-Correlation-Id` quando o enforcement esta ativo.
- Gateway aplica rate limit global, por tenant, por actor e por rotas sensiveis.

## Dados Sensiveis

- Password recovery persiste hash do token e nao devolve token fora de local/test.
- Access links de documents sao assinados, expiraveis, revogaveis e auditados com persistencia relacional de revogacoes e eventos.
- Uploads com assinatura maliciosa conhecida sao bloqueados antes de virarem attachment.
- Logs, contratos e portal nao devem expor segredo, token, senha ou provider credential.

## Suite De Seguranca

`./scripts/test.sh security` executa validacoes estaticas de:

- middleware JWT/OpenFGA por stack;
- rate limit do gateway;
- hardening de secrets;
- recovery token com hash;
- documents access links, revogacao persistente, auditoria e scan;
- envelope/event registry;
- portal de integracao;
- inventario LGPD;
- SRE/runbooks;
- guardrails de poliglotismo.
