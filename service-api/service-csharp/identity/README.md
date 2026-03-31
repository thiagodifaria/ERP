# identity

The identity service owns tenants, companies, users, teams, roles and access foundations.

Current scope:

- minimal API bootstrap by layer
- health and readiness endpoints
- tenant bootstrap with default structure provisioning
- company, user and team management in bootstrap mode
- team memberships and direct user-role assignments
- consolidated tenant access snapshot
- selectable repository driver with PostgreSQL-backed runtime support
- company update path started for existing tenant structure records
- user update path started for existing tenant records
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
- `POST /api/identity/tenants/{slug}/companies`
- `PATCH /api/identity/tenants/{slug}/companies/{companyPublicId}`
- `GET /api/identity/tenants/{slug}/users`
- `POST /api/identity/tenants/{slug}/users`
- `PATCH /api/identity/tenants/{slug}/users/{userPublicId}`
- `GET /api/identity/tenants/{slug}/users/{userPublicId}/roles`
- `POST /api/identity/tenants/{slug}/users/{userPublicId}/roles`
- `DELETE /api/identity/tenants/{slug}/users/{userPublicId}/roles/{roleCode}`
- `GET /api/identity/tenants/{slug}/teams`
- `POST /api/identity/tenants/{slug}/teams`
- `PATCH /api/identity/tenants/{slug}/teams/{teamPublicId}`
- `GET /api/identity/tenants/{slug}/teams/{teamPublicId}/members`
- `POST /api/identity/tenants/{slug}/teams/{teamPublicId}/members`
- `DELETE /api/identity/tenants/{slug}/teams/{teamPublicId}/members/{userPublicId}`
- `GET /api/identity/tenants/{slug}/roles`

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/identity mcr.microsoft.com/dotnet/sdk:8.0 dotnet test tests/Identity.UnitTests/Identity.UnitTests.csproj`
- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/identity mcr.microsoft.com/dotnet/sdk:8.0 dotnet test tests/Identity.IntegrationTests/Identity.IntegrationTests.csproj`
- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-csharp/identity mcr.microsoft.com/dotnet/sdk:8.0 dotnet test tests/Identity.ContractTests/Identity.ContractTests.csproj`
