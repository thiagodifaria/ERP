# platform-control

Contexto de plataforma SaaS para capabilities, entitlements, metering e lifecycle de tenant.

## Escopo atual

- catalogo publico de capabilities
- entitlements por tenant
- snapshots simples de metering
- jobs de onboarding e offboarding
- persistencia PostgreSQL ou bootstrap em memoria

## Rotas publicas

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/platform-control/capabilities/catalog`
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements`
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`
- `GET /api/platform-control/tenants/{tenantSlug}/metering`
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`
