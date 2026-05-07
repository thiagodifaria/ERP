# fiscal

Servico de contexto fiscal, privacidade e governanca de dados.

## Escopo atual

- perfil fiscal por empresa
- politicas de retencao e classificacao por dominio
- emissao/cancelamento basicos de documentos fiscais com trilha auditavel
- requests de LGPD com rastreabilidade
- eventos de auditoria para operacoes fiscais e administrativas sensiveis

## Endpoints

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/fiscal/capabilities`
- `GET /api/fiscal/companies/{companyPublicId}/profile`
- `PUT /api/fiscal/companies/{companyPublicId}/profile`
- `GET /api/fiscal/companies/{companyPublicId}/retention-policies`
- `PUT /api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}`
- `GET /api/fiscal/documents`
- `POST /api/fiscal/documents`
- `POST /api/fiscal/documents/{publicId}/cancel`
- `POST /api/fiscal/documents/{publicId}/correction-letter`
- `POST /api/fiscal/documents/{publicId}/invalidate`
- `GET /api/fiscal/documents/{publicId}/events`
- `POST /api/fiscal/privacy-requests`
- `GET /api/fiscal/privacy-requests`
- `PATCH /api/fiscal/privacy-requests/{publicId}/status`
- `POST /api/fiscal/consents`
- `GET /api/fiscal/consents`
- `PATCH /api/fiscal/consents/{publicId}`
- `GET /api/fiscal/audit-events`
- `GET /api/fiscal/compliance/summary`
