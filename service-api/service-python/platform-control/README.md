# platform-control

Contexto de plataforma SaaS para capability registry, entitlements, quotas, metering, tenant lifecycle e bloqueios operacionais.

## Escopo atual

- catalogo publico de capabilities por modulo e capability
- catalogo publico de providers externos por capability
- entitlements por tenant com bulk e partial success
- defaults de provider por tenant com governanca de modo `configured`, `fallback`, `manual`, `disabled` e `unconfigured`
- quotas por tenant com enforcement `soft` e `hard`
- snapshots de metering e resumo de utilizacao
- bloqueios administrativos por tenant
- readiness de lifecycle por tenant e previews de onboarding/offboarding
- jobs de onboarding e offboarding com `Idempotency-Key`, `202 Accepted`, polling e trilha de eventos
- persistencia PostgreSQL ou bootstrap em memoria

## Rotas publicas

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/platform-control/capabilities/catalog`
- `GET /api/platform-control/providers/catalog`
- `GET /api/platform-control/tenants/{tenantSlug}/entitlements`
- `PUT /api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}`
- `POST /api/platform-control/tenants/{tenantSlug}/entitlements/bulk`
- `GET /api/platform-control/tenants/{tenantSlug}/provider-defaults`
- `PUT /api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}`
- `GET /api/platform-control/tenants/{tenantSlug}/quotas`
- `PUT /api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}`
- `POST /api/platform-control/tenants/{tenantSlug}/quotas/bulk`
- `GET /api/platform-control/tenants/{tenantSlug}/blocks`
- `PUT /api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}`
- `GET /api/platform-control/tenants/{tenantSlug}/metering`
- `POST /api/platform-control/tenants/{tenantSlug}/metering/snapshots`
- `GET /api/platform-control/tenants/{tenantSlug}/usage-summary`
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness`
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs`
- `GET /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail`
- `POST /api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel`
