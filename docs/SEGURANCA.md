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
- Exportacao, portabilidade e anonimizacao passam por politica de dominio, preservando evidencias fiscais, financeiras e auditorias com obrigacao legal de retencao.
- Novas features com dado pessoal precisam declarar finalidade, retencao e exposicao em contrato/evento/log antes de entrar na superficie publica.

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
- hardening de secrets;
- recovery token com hash;
- documents access links, revogacao persistente, auditoria e scan;
- envelope/event registry;
- portal de integracao;
- inventario LGPD;
- SLOs/runbooks em `docs/OPERACOES.md`;
- guardrails de poliglotismo.
