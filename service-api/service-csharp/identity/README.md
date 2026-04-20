# identity

The identity service owns tenants, companies, users, teams, roles and access foundations.

Current scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- tenant bootstrap with default structure provisioning
- company, user and team management in bootstrap mode
- team memberships and direct user-role assignments
- consolidated tenant access snapshot
- invite lifecycle for tenant onboarding
- session login and refresh flow
- tenant access resolution for downstream gateways
- MFA enrollment, verification and disable flow
- password recovery and reset flow with strong password validation
- session governance with tenant-scoped session list and revocation
- security audit trail for access-sensitive actions
- selectable repository driver with PostgreSQL-backed runtime support
- Keycloak-backed identity provider integration
- OpenFGA-backed tenant authorization graph
- runtime smoke now exercises the live HTTP API against PostgreSQL
- unit, integration and contract coverage in container-first flow

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/identity/tenants`
- `POST /api/identity/tenants`
- `GET /api/identity/tenants/{slug}`
- `GET /api/identity/tenants/{slug}/snapshot`
- `GET /api/identity/tenants/{slug}/companies`
- `GET /api/identity/tenants/{slug}/companies/{companyPublicId}`
- `POST /api/identity/tenants/{slug}/companies`
- `PATCH /api/identity/tenants/{slug}/companies/{companyPublicId}`
- `GET /api/identity/tenants/{slug}/users`
- `POST /api/identity/tenants/{slug}/users`
- `GET /api/identity/tenants/{slug}/users/{userPublicId}`
- `PATCH /api/identity/tenants/{slug}/users/{userPublicId}`
- `GET /api/identity/tenants/{slug}/users/{userPublicId}/roles`
- `POST /api/identity/tenants/{slug}/users/{userPublicId}/roles`
- `DELETE /api/identity/tenants/{slug}/users/{userPublicId}/roles/{roleCode}`
- `PATCH /api/identity/tenants/{slug}/users/{userPublicId}/access`
- `POST /api/identity/tenants/{slug}/users/{userPublicId}/mfa/enroll`
- `POST /api/identity/tenants/{slug}/users/{userPublicId}/mfa/verify`
- `DELETE /api/identity/tenants/{slug}/users/{userPublicId}/mfa`
- `GET /api/identity/tenants/{slug}/teams`
- `POST /api/identity/tenants/{slug}/teams`
- `PATCH /api/identity/tenants/{slug}/teams/{teamPublicId}`
- `GET /api/identity/tenants/{slug}/teams/{teamPublicId}/members`
- `POST /api/identity/tenants/{slug}/teams/{teamPublicId}/members`
- `DELETE /api/identity/tenants/{slug}/teams/{teamPublicId}/members/{userPublicId}`
- `GET /api/identity/tenants/{slug}/roles`
- `POST /api/identity/tenants/{slug}/invites`
- `GET /api/identity/tenants/{slug}/invites`
- `GET /api/identity/tenants/{slug}/invites/{invitePublicId}`
- `POST /api/identity/invites/{inviteToken}/accept`
- `POST /api/identity/sessions/login`
- `POST /api/identity/sessions/refresh`
- `POST /api/identity/password-recovery`
- `POST /api/identity/password-recovery/{resetToken}/complete`
- `GET /api/identity/tenants/{slug}/access`
- `GET /api/identity/tenants/{slug}/users/{userPublicId}/mfa`
- `GET /api/identity/tenants/{slug}/users/{userPublicId}/sessions`
- `DELETE /api/identity/tenants/{slug}/sessions/{sessionPublicId}`
- `DELETE /api/identity/tenants/{slug}/users/{userPublicId}/sessions`
- `GET /api/identity/tenants/{slug}/security/audit`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/identity mcr.microsoft.com/dotnet/sdk:8.0 dotnet test tests/Identity.UnitTests/Identity.UnitTests.csproj`
- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/identity mcr.microsoft.com/dotnet/sdk:8.0 dotnet test tests/Identity.IntegrationTests/Identity.IntegrationTests.csproj`
- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/identity mcr.microsoft.com/dotnet/sdk:8.0 dotnet test tests/Identity.ContractTests/Identity.ContractTests.csproj`
- `bash scripts/test.sh smoke`
