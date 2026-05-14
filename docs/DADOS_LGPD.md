# DADOS E LGPD

Este documento registra a governanca de dados pessoais do ERP. Ele complementa `docs/SERVICOS.md` e `docs/OPERACOES.md`; aqui o foco e inventario, sensibilidade, finalidade, retencao e direitos do titular.

## Principios

- Todo dado pessoal deve ter tenant, finalidade e base operacional clara.
- Exportacao, portabilidade e anonimizacao passam por politica de dominio, sem apagar evidencias fiscais, financeiras ou auditorias que tenham obrigacao legal de retencao.
- Logs e analytics nao devem registrar senha, token, access link, documento fiscal completo, segredo de provider ou conteudo sensivel de documento.
- Novas features precisam passar pelo checklist de impacto de privacidade antes de entrar no contrato publico.

## Inventario Resumido

| Dominio | Dados pessoais | Sensibilidade | Finalidade | Retencao |
| --- | --- | --- | --- | --- |
| `identity` | nome, e-mail, papeis, sessoes, MFA, auditoria | alta | autenticacao, autorizacao e trilha de acesso | enquanto conta ativa e janela legal de auditoria |
| `crm` | leads, clientes, contatos, notas, anexos vinculados | alta | relacionamento comercial e funil | politica comercial por tenant |
| `engagement` | touchpoints, entregas, callbacks e respostas | alta | comunicacao transacional/comercial | consentimento e necessidade operacional |
| `support` | casos, comentarios, filas e SLA | media/alta | atendimento e evidencias de suporte | janela contratual do suporte |
| `supplier` | contatos de fornecedores e dados fiscais de fornecedor | media | operacao de compras/fornecimento | obrigacao contratual/fiscal |
| `billing` | assinaturas, invoices, tentativas de pagamento, recovery | alta | cobranca, inadimplencia e receita | obrigacao financeira |
| `finance` | recebiveis, payables, caixa, movimentos, reconciliacao | alta | controle financeiro e auditoria | obrigacao financeira/fiscal |
| `documents` | arquivos, metadata, links temporarios e assinatura | critica quando restrito | guarda documental e evidencia operacional | politica por classificacao/retencao |
| `fiscal` | documentos fiscais, consentimentos, requests LGPD, auditoria | critica | compliance fiscal, privacidade e retencao legal | conforme obrigacao legal |

## Checklist De Impacto

- O dado novo identifica pessoa direta ou indiretamente?
- O dado precisa aparecer no contrato HTTP ou pode ficar interno?
- Existe base legal/finalidade documentada?
- Existe retencao definida e expurgo/anonimizacao possivel?
- O dado pode aparecer em log, analytics, evento ou webhook?
- Exportacao/portabilidade preserva segregacao por tenant?
- A anonimizacao preserva integridade fiscal/financeira quando a retencao legal exige?

## Controles Verificaveis

- `fiscal` concentra privacy requests, consentimento, retencao e auditoria.
- `documents` bloqueia access link revogado, registra eventos de acesso e impede download depois da janela de retencao.
- `scripts/test.sh security` valida a existencia deste inventario e dos guardrails de seguranca transversais.
