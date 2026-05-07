# support

Servico de atendimento e `case management` do ERP.

## Escopo atual

- filas de atendimento por tenant
- casos com prioridade, responsavel, SLA e vinculo com entidade de negocio
- historico auditavel por eventos
- resumo operacional por status e prioridade

## Endpoints

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/support/capabilities`
- `GET /api/support/queues`
- `PUT /api/support/queues/{queueKey}`
- `GET /api/support/cases`
- `GET /api/support/cases/export`
- `POST /api/support/cases`
- `POST /api/support/cases/bulk`
- `GET /api/support/cases/summary`
- `GET /api/support/cases/{publicId}`
- `PATCH /api/support/cases/{publicId}/status`
- `POST /api/support/cases/{publicId}/comments`
