# notification

Servico de centro interno de alertas e preferencias.

## Escopo atual

- preferencias por usuario e tenant
- notificacoes internas com severidade, canal e vinculo com entidade de negocio
- transicao `unread -> read -> archived`
- resumo operacional para consumo por `analytics` e `edge`

## Endpoints

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/notification/capabilities`
- `GET /api/notification/preferences/{userPublicId}`
- `PUT /api/notification/preferences/{userPublicId}`
- `GET /api/notification/center`
- `POST /api/notification/center`
- `POST /api/notification/center/bulk`
- `PATCH /api/notification/center/{publicId}/status`
- `GET /api/notification/summary`
