# SRE

Este documento define SLOs, sinais de alerta e runbooks curtos para operacao do ERP. Ele nao substitui monitoramento real, mas torna explicitos os criterios operacionais esperados.

## SLOs Iniciais

| Jornada | Indicador | Objetivo local/corporativo |
| --- | --- | --- |
| Auth/login | taxa de sucesso e p95 | 99.5% de sucesso, p95 abaixo de 800ms |
| CRM lead intake | criacao de lead e outbox | 99.0% de sucesso, sem backlog critico |
| Sales creation | criacao de oportunidade/sale/invoice | 99.0% de sucesso, idempotencia preservada |
| Billing payment attempt | tentativa e callback | 99.0% processado ou DLQ explicita |
| Webhook ingestion | validacao, fila, forward/DLQ | 99.0% sem perda, replay auditavel |
| Document download | access link valido para redirect | 99.5% de sucesso, revogacao imediata |
| Fiscal issue/cancel | documento/evento fiscal | 99.0% com trilha de auditoria |

## Sinais De Alerta

- Erro 5xx por rota acima do limiar do dominio.
- Latencia p95/p99 fora do SLO por 10 minutos.
- Crescimento de DLQ ou retry sem drenagem.
- Queda de provider ou callback com assinatura invalida.
- Falha de backup ou restore drill.
- Divergencia de reconciliacao financeira.
- Tentativas repetidas de login, recovery, upload ou webhook acima do rate limit.

## Runbooks Curtos

### DLQ Ou Retry Crescendo

1. Conferir `/api/webhook-hub/events/dead-letter` e summary.
2. Validar `correlationId`, provider, payload e schemaRef.
3. Reprocessar apenas eventos idempotentes.
4. Se o erro for contrato, bloquear replay ate corrigir schema/adapter.

### Download De Documento Bloqueado

1. Verificar se o access link expirou, foi revogado ou falhou tenant.
2. Consultar `/api/documents/audit-events`.
3. Gerar novo link somente com ator autorizado.
4. Se a retencao expirou, seguir politica de fiscal/LGPD.

### Falha De Backup

1. Rodar `./scripts/test.sh backup-restore`.
2. Conferir espaco, permissao e conexao PostgreSQL.
3. Repetir restore em ambiente limpo.
4. Registrar RPO/RTO afetado e janela de dados impactada.

### Divergencia Financeira

1. Conferir cash movements, settlement/payment reference e period closure.
2. Validar se a mutacao usou `Idempotency-Key`.
3. Registrar ajuste/estorno em vez de update destrutivo.
4. Bloquear fechamento de periodo se a divergencia continuar.
