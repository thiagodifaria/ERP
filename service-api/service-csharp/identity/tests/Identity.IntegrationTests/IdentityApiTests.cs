// Estes testes validam o comportamento HTTP minimo da API de identidade.
using System.Net;
using System.Net.Http.Json;
using System.Security.Cryptography;
using System.Text;
using Identity.Contracts;
using Microsoft.AspNetCore.Mvc.Testing;
using Xunit;

namespace Identity.IntegrationTests;

public sealed class IdentityApiTests : IClassFixture<WebApplicationFactory<Program>>
{
  private readonly HttpClient _client;

  public IdentityApiTests(WebApplicationFactory<Program> factory)
  {
    _client = factory.CreateClient();
  }

  [Fact]
  public async Task HealthDetailsShouldReturnReadyDependencies()
  {
    var response = await _client.GetAsync("/health/details");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ReadinessResponse>();

    Assert.NotNull(payload);
    Assert.Equal("identity", payload.Service);
    Assert.Equal("ready", payload.Status);
    Assert.Contains(payload.Dependencies, dependency => dependency.Name == "tenant-catalog");
  }

  [Fact]
  public async Task TenantsEndpointShouldReturnBootstrapTenants()
  {
    var response = await _client.GetAsync("/api/identity/tenants");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, tenant => tenant.Slug == "bootstrap-ops");
  }

  [Fact]
  public async Task TenantBySlugEndpointShouldReturnTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(payload);
    Assert.Equal("bootstrap-ops", payload.Slug);
  }

  [Fact]
  public async Task TenantBySlugEndpointShouldReturnNotFoundForUnknownSlug()
  {
    var response = await _client.GetAsync("/api/identity/tenants/missing-tenant");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task TenantSnapshotEndpointShouldReturnSnapshotForKnownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/snapshot");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantAccessSnapshotResponse>();

    Assert.NotNull(payload);
    Assert.Equal("bootstrap-ops", payload!.Tenant.Slug);
    Assert.Equal(1, payload.Counts.Companies);
    Assert.Equal(1, payload.Counts.Users);
    Assert.Equal(1, payload.Counts.Teams);
    Assert.Equal(5, payload.Counts.Roles);
    Assert.Equal(1, payload.Counts.TeamMemberships);
    Assert.Equal(1, payload.Counts.UserRoles);
    Assert.Single(payload.Users);
    Assert.Single(payload.Users.First().Roles);
    Assert.Single(payload.Teams);
    Assert.Single(payload.Teams.First().Members);
  }

  [Fact]
  public async Task TenantSnapshotEndpointShouldReturnNotFoundForUnknownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/missing-tenant/snapshot");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task TenantCompaniesEndpointShouldReturnCompaniesForKnownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/companies");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, company =>
      company.TenantId == 1 &&
      !string.IsNullOrWhiteSpace(company.DisplayName));
  }

  [Fact]
  public async Task TenantCompaniesEndpointShouldReturnNotFoundForUnknownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/missing-tenant/companies");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task CompanyByPublicIdEndpointShouldReturnCompanyForKnownTenant()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Company Lookup");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var companiesResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/companies");
    var companiesPayload = await companiesResponse.Content.ReadFromJsonAsync<CompanyResponse[]>();

    Assert.NotNull(companiesPayload);
    Assert.NotEmpty(companiesPayload);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/companies/{companiesPayload![0].PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse>();

    Assert.NotNull(payload);
    Assert.Equal(companiesPayload[0].PublicId, payload!.PublicId);
    Assert.Equal(tenantPayload.DisplayName, payload.DisplayName);
  }

  [Fact]
  public async Task CompanyByPublicIdEndpointShouldReturnNotFoundForUnknownCompany()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Missing Company");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/companies/{Guid.NewGuid()}");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task TenantUsersEndpointShouldReturnUsersForKnownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/users");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, user => user.Email == "owner@bootstrap-ops.local");
  }

  [Fact]
  public async Task TenantUsersEndpointShouldReturnNotFoundForUnknownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/missing-tenant/users");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task UserByPublicIdEndpointShouldReturnUserForKnownTenant()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With User Lookup");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"lookup.user@{tenantPayload.Slug}.local",
        "Lookup User",
        "Lookup",
        "User"));

    var createdUser = await createResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.Equal(createdUser.PublicId, payload!.PublicId);
    Assert.Equal($"lookup.user@{tenantPayload.Slug}.local", payload.Email);
  }

  [Fact]
  public async Task UserByPublicIdEndpointShouldReturnNotFoundForUnknownUser()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Missing User");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users/{Guid.NewGuid()}");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task UserRolesEndpointShouldReturnRolesForKnownUser()
  {
    var usersPayload = await _client.GetFromJsonAsync<UserResponse[]>("/api/identity/tenants/bootstrap-ops/users");

    Assert.NotNull(usersPayload);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/bootstrap-ops/users/{usersPayload![0].PublicId}/roles");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserRoleResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, userRole => userRole.RoleCode == "owner");
  }

  [Fact]
  public async Task UserRolesEndpointShouldReturnNotFoundForUnknownUser()
  {
    var response = await _client.GetAsync(
      $"/api/identity/tenants/bootstrap-ops/users/{Guid.NewGuid()}/roles");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task TenantTeamsEndpointShouldReturnTeamsForKnownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/teams");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, team => team.Name == "Core");
  }

  [Fact]
  public async Task TenantTeamsEndpointShouldReturnNotFoundForUnknownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/missing-tenant/teams");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task TeamMembersEndpointShouldReturnMembersForKnownTeam()
  {
    var teamsResponse = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/teams");
    var teamsPayload = await teamsResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamsPayload);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/bootstrap-ops/teams/{teamsPayload![0].PublicId}/members");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamMembershipResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, member => member.UserEmail == "owner@bootstrap-ops.local");
  }

  [Fact]
  public async Task TeamMembersEndpointShouldReturnNotFoundForUnknownTeam()
  {
    var response = await _client.GetAsync(
      $"/api/identity/tenants/bootstrap-ops/teams/{Guid.NewGuid()}/members");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task CreateTenantShouldReturnCreatedTenant()
  {
    var request = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Created Tenant");

    var response = await _client.PostAsJsonAsync("/api/identity/tenants", request);

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(payload);
    Assert.Equal(request.Slug, payload.Slug);
    Assert.Equal(request.DisplayName, payload.DisplayName);

    var rolesResponse = await _client.GetAsync($"/api/identity/tenants/{payload.Slug}/roles");

    Assert.Equal(HttpStatusCode.OK, rolesResponse.StatusCode);

    var rolesPayload = await rolesResponse.Content.ReadFromJsonAsync<RoleResponse[]>();

    Assert.NotNull(rolesPayload);
    Assert.Equal(5, rolesPayload.Length);
    Assert.Contains(rolesPayload, role => role.Code == "owner");

    var companiesResponse = await _client.GetAsync($"/api/identity/tenants/{payload.Slug}/companies");

    Assert.Equal(HttpStatusCode.OK, companiesResponse.StatusCode);

    var companiesPayload = await companiesResponse.Content.ReadFromJsonAsync<CompanyResponse[]>();

    Assert.NotNull(companiesPayload);
    Assert.Single(companiesPayload);
    Assert.Equal(request.DisplayName, companiesPayload[0].DisplayName);

    var usersResponse = await _client.GetAsync($"/api/identity/tenants/{payload.Slug}/users");

    Assert.Equal(HttpStatusCode.OK, usersResponse.StatusCode);

    var usersPayload = await usersResponse.Content.ReadFromJsonAsync<UserResponse[]>();

    Assert.NotNull(usersPayload);
    Assert.Single(usersPayload);
    Assert.Equal($"owner@{payload.Slug}.local", usersPayload[0].Email);

    var userRolesResponse = await _client.GetAsync(
      $"/api/identity/tenants/{payload.Slug}/users/{usersPayload[0].PublicId}/roles");

    Assert.Equal(HttpStatusCode.OK, userRolesResponse.StatusCode);

    var userRolesPayload = await userRolesResponse.Content.ReadFromJsonAsync<UserRoleResponse[]>();

    Assert.NotNull(userRolesPayload);
    Assert.Single(userRolesPayload);
    Assert.Equal("owner", userRolesPayload[0].RoleCode);

    var teamsResponse = await _client.GetAsync($"/api/identity/tenants/{payload.Slug}/teams");

    Assert.Equal(HttpStatusCode.OK, teamsResponse.StatusCode);

    var teamsPayload = await teamsResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamsPayload);
    Assert.Single(teamsPayload);
    Assert.Equal("Core", teamsPayload[0].Name);

    var membersResponse = await _client.GetAsync(
      $"/api/identity/tenants/{payload.Slug}/teams/{teamsPayload[0].PublicId}/members");

    Assert.Equal(HttpStatusCode.OK, membersResponse.StatusCode);

    var membersPayload = await membersResponse.Content.ReadFromJsonAsync<TeamMembershipResponse[]>();

    Assert.NotNull(membersPayload);
    Assert.Single(membersPayload);
    Assert.Equal($"owner@{payload.Slug}.local", membersPayload[0].UserEmail);

    var snapshotResponse = await _client.GetAsync($"/api/identity/tenants/{payload.Slug}/snapshot");

    Assert.Equal(HttpStatusCode.OK, snapshotResponse.StatusCode);

    var snapshotPayload = await snapshotResponse.Content.ReadFromJsonAsync<TenantAccessSnapshotResponse>();

    Assert.NotNull(snapshotPayload);
    Assert.Equal(1, snapshotPayload!.Counts.Companies);
    Assert.Equal(1, snapshotPayload.Counts.Users);
    Assert.Equal(1, snapshotPayload.Counts.Teams);
    Assert.Equal(5, snapshotPayload.Counts.Roles);
    Assert.Equal(1, snapshotPayload.Counts.TeamMemberships);
    Assert.Equal(1, snapshotPayload.Counts.UserRoles);
  }

  [Fact]
  public async Task CreateTenantShouldReturnConflictForDuplicateSlug()
  {
    var request = new CreateTenantRequest("bootstrap-ops", "Bootstrap Ops Duplicate");

    var response = await _client.PostAsJsonAsync("/api/identity/tenants", request);

    Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("tenant_slug_conflict", payload.Code);
  }

  [Fact]
  public async Task CreateTenantShouldReturnBadRequestForInvalidSlug()
  {
    var request = new CreateTenantRequest("invalid slug", "Invalid Slug");

    var response = await _client.PostAsJsonAsync("/api/identity/tenants", request);

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("invalid_slug", payload.Code);
  }

  [Fact]
  public async Task TenantRolesEndpointShouldReturnRolesForKnownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/roles");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<RoleResponse[]>();

    Assert.NotNull(payload);
    Assert.Equal(5, payload.Length);
    Assert.Contains(payload, role => role.Code == "owner");
  }

  [Fact]
  public async Task TenantRolesEndpointShouldReturnNotFoundForUnknownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/unknown/roles");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
  }

  [Fact]
  public async Task CreateCompanyShouldReturnCreatedCompany()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Branch");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var request = new CreateCompanyRequest(
      "Tenant With Branch Filial",
      "Tenant With Branch Filial LTDA",
      "12345678901234");

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/companies",
      request);

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse>();

    Assert.NotNull(payload);
    Assert.Equal(request.DisplayName, payload.DisplayName);
    Assert.Equal(request.LegalName, payload.LegalName);
    Assert.Equal(request.TaxId, payload.TaxId);
  }

  [Fact]
  public async Task UpdateCompanyShouldReturnUpdatedCompany()
  {
    var companiesResponse = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/companies");
    var companiesPayload = await companiesResponse.Content.ReadFromJsonAsync<CompanyResponse[]>();

    Assert.NotNull(companiesPayload);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/companies/{companiesPayload![0].PublicId}",
      new UpdateCompanyRequest("Bootstrap Ops Prime", "Bootstrap Ops Prime LTDA", "12345678901234"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse>();

    Assert.NotNull(payload);
    Assert.Equal("Bootstrap Ops Prime", payload!.DisplayName);
    Assert.Equal("Bootstrap Ops Prime LTDA", payload.LegalName);
    Assert.Equal("12345678901234", payload.TaxId);
  }

  [Fact]
  public async Task UpdateCompanyShouldReturnNotFoundForUnknownCompany()
  {
    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/companies/{Guid.NewGuid()}",
      new UpdateCompanyRequest("Bootstrap Ops Prime", null, null));

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("company_not_found", payload.Code);
  }

  [Fact]
  public async Task CreateCompanyShouldReturnConflictForDuplicateDisplayName()
  {
    var request = new CreateCompanyRequest("Bootstrap Ops", null, null);

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/companies",
      request);

    Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("company_display_name_conflict", payload.Code);
  }

  [Fact]
  public async Task CreateCompanyShouldReturnNotFoundForUnknownTenant()
  {
    var request = new CreateCompanyRequest("Missing Tenant Company", null, null);

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/missing-tenant/companies",
      request);

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("tenant_not_found", payload.Code);
  }

  [Fact]
  public async Task CreateCompanyShouldReturnBadRequestForBlankDisplayName()
  {
    var request = new CreateCompanyRequest("   ", null, null);

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/companies",
      request);

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("invalid_display_name", payload.Code);
  }

  [Fact]
  public async Task CreateUserShouldReturnCreatedUser()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With User");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var request = new CreateUserRequest(
      "tenant.user@tenant-with-user.local",
      "Tenant User",
      "Tenant",
      "User");

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      request);

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.Equal(request.Email, payload.Email);
    Assert.Equal(request.DisplayName, payload.DisplayName);
  }

  [Fact]
  public async Task UpdateUserShouldReturnUpdatedUser()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With User Update");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"tenant.user@{tenantPayload.Slug}.local",
        "Tenant User",
        "Tenant",
        "User"));

    var createdUser = await createResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}",
      new UpdateUserRequest(
        $"tenant.user.prime@{tenantPayload.Slug}.local",
        "Tenant User Prime",
        "Tenant",
        "Prime"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.Equal($"tenant.user.prime@{tenantPayload.Slug}.local", payload!.Email);
    Assert.Equal("Tenant User Prime", payload.DisplayName);
    Assert.Equal("Prime", payload.FamilyName);
  }

  [Fact]
  public async Task UpdateUserShouldReturnNotFoundForUnknownUser()
  {
    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/users/{Guid.NewGuid()}",
      new UpdateUserRequest("owner.prime@bootstrap-ops.local", "Bootstrap Owner Prime", "Bootstrap", "Prime"));

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("user_not_found", payload.Code);
  }

  [Fact]
  public async Task CreateUserShouldReturnConflictForDuplicateEmail()
  {
    var request = new CreateUserRequest("owner@bootstrap-ops.local", "Duplicate Owner", null, null);

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/users",
      request);

    Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("user_email_conflict", payload.Code);
  }

  [Fact]
  public async Task CreateUserShouldReturnBadRequestForInvalidEmail()
  {
    var request = new CreateUserRequest("invalid-email", "Invalid User", null, null);

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/users",
      request);

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("invalid_email", payload.Code);
  }

  [Fact]
  public async Task CreateTeamShouldReturnCreatedTeam()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Team");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var request = new CreateTeamRequest("Field Ops");

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/teams",
      request);

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamResponse>();

    Assert.NotNull(payload);
    Assert.Equal(request.Name, payload.Name);
  }

  [Fact]
  public async Task UpdateTeamShouldReturnUpdatedTeam()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Team Update");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teamsPayload = await _client.GetFromJsonAsync<TeamResponse[]>(
      $"/api/identity/tenants/{tenantPayload!.Slug}/teams");

    Assert.NotNull(teamsPayload);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamsPayload![0].PublicId}",
      new UpdateTeamRequest("Core Prime"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamResponse>();

    Assert.NotNull(payload);
    Assert.Equal("Core Prime", payload!.Name);
  }

  [Fact]
  public async Task UpdateTeamShouldReturnNotFoundForUnknownTeam()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Missing Team");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/teams/{Guid.NewGuid()}",
      new UpdateTeamRequest("Core Prime"));

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("team_not_found", payload.Code);
  }

  [Fact]
  public async Task CreateTeamShouldReturnConflictForDuplicateName()
  {
    var request = new CreateTeamRequest("Core");

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/teams",
      request);

    Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("team_name_conflict", payload.Code);
  }

  [Fact]
  public async Task CreateTeamShouldReturnNotFoundForUnknownTenant()
  {
    var request = new CreateTeamRequest("Field Ops");

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/missing-tenant/teams",
      request);

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("tenant_not_found", payload.Code);
  }

  [Fact]
  public async Task CreateTeamShouldReturnBadRequestForBlankName()
  {
    var request = new CreateTeamRequest("   ");

    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/teams",
      request);

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("invalid_team_name", payload.Code);
  }

  [Fact]
  public async Task AddTeamMemberShouldReturnCreatedMembership()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Membership");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teamsResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/teams");
    var teamsPayload = await teamsResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamsPayload);

    var userResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users",
      new CreateUserRequest("member@tenant-membership.local", "Member User", null, null));

    var userPayload = await userResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(userPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamsPayload![0].PublicId}/members",
      new AddTeamMemberRequest(userPayload!.PublicId));

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamMembershipResponse>();

    Assert.NotNull(payload);
    Assert.Equal(userPayload.PublicId, payload.UserPublicId);
  }

  [Fact]
  public async Task RemoveTeamMemberShouldReturnRemovedMembership()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Membership Removal");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teamResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/teams");
    var teamPayload = await teamResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamPayload);

    var userResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users",
      new CreateUserRequest("member.remove@tenant-membership-removal.local", "Member Remove", null, null));

    var userPayload = await userResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(userPayload);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamPayload![0].PublicId}/members",
      new AddTeamMemberRequest(userPayload!.PublicId));

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamPayload[0].PublicId}/members/{userPayload.PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamMembershipResponse>();

    Assert.NotNull(payload);
    Assert.Equal(userPayload.PublicId, payload!.UserPublicId);

    var membersResponse = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamPayload[0].PublicId}/members");

    var membersPayload = await membersResponse.Content.ReadFromJsonAsync<TeamMembershipResponse[]>();

    Assert.NotNull(membersPayload);
    Assert.DoesNotContain(membersPayload, member => member.UserPublicId == userPayload.PublicId);
  }

  [Fact]
  public async Task RemoveTeamMemberShouldReturnNotFoundForMissingMembership()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Missing Membership");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teamResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/teams");
    var teamPayload = await teamResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamPayload);

    var userResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users",
      new CreateUserRequest("missing.membership@tenant-membership-removal.local", "Missing Membership", null, null));

    var userPayload = await userResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(userPayload);

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamPayload![0].PublicId}/members/{userPayload!.PublicId}");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("team_membership_not_found", payload.Code);
  }

  [Fact]
  public async Task AddTeamMemberShouldReturnConflictForDuplicateMembership()
  {
    var teamsPayload = await _client.GetFromJsonAsync<TeamResponse[]>("/api/identity/tenants/bootstrap-ops/teams");
    var usersPayload = await _client.GetFromJsonAsync<UserResponse[]>("/api/identity/tenants/bootstrap-ops/users");

    Assert.NotNull(teamsPayload);
    Assert.NotNull(usersPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/teams/{teamsPayload![0].PublicId}/members",
      new AddTeamMemberRequest(usersPayload![0].PublicId));

    Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("team_membership_conflict", payload.Code);
  }

  [Fact]
  public async Task AddTeamMemberShouldReturnBadRequestForEmptyUserPublicId()
  {
    var teamsPayload = await _client.GetFromJsonAsync<TeamResponse[]>("/api/identity/tenants/bootstrap-ops/teams");

    Assert.NotNull(teamsPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/teams/{teamsPayload![0].PublicId}/members",
      new AddTeamMemberRequest(Guid.Empty));

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("invalid_user_public_id", payload.Code);
  }

  [Fact]
  public async Task AssignUserRoleShouldReturnCreatedRoleAssignment()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Extra Role");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var userResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest("role.user@tenant-extra-role.local", "Role User", null, null));

    var userPayload = await userResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(userPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{userPayload!.PublicId}/roles",
      new AssignUserRoleRequest("viewer"));

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserRoleResponse>();

    Assert.NotNull(payload);
    Assert.Equal("viewer", payload.RoleCode);
  }

  [Fact]
  public async Task RevokeUserRoleShouldReturnRevokedRoleAssignment()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Tenant With Role Revoke");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var userResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest("revoke.role@tenant-role-revoke.local", "Revoke Role User", null, null));

    var userPayload = await userResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(userPayload);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{userPayload!.PublicId}/roles",
      new AssignUserRoleRequest("viewer"));

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{userPayload.PublicId}/roles/viewer");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserRoleResponse>();

    Assert.NotNull(payload);
    Assert.Equal("viewer", payload!.RoleCode);

    var rolesResponse = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{userPayload.PublicId}/roles");

    var rolesPayload = await rolesResponse.Content.ReadFromJsonAsync<UserRoleResponse[]>();

    Assert.NotNull(rolesPayload);
    Assert.Empty(rolesPayload);
  }

  [Fact]
  public async Task RevokeUserRoleShouldReturnNotFoundForMissingAssignment()
  {
    var usersPayload = await _client.GetFromJsonAsync<UserResponse[]>("/api/identity/tenants/bootstrap-ops/users");

    Assert.NotNull(usersPayload);

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/bootstrap-ops/users/{usersPayload![0].PublicId}/roles/viewer");

    Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("user_role_not_found", payload.Code);
  }

  [Fact]
  public async Task AssignUserRoleShouldReturnConflictForDuplicateRole()
  {
    var usersPayload = await _client.GetFromJsonAsync<UserResponse[]>("/api/identity/tenants/bootstrap-ops/users");

    Assert.NotNull(usersPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/users/{usersPayload![0].PublicId}/roles",
      new AssignUserRoleRequest("owner"));

    Assert.Equal(HttpStatusCode.Conflict, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("user_role_conflict", payload.Code);
  }

  [Fact]
  public async Task AssignUserRoleShouldReturnBadRequestForBlankRoleCode()
  {
    var usersPayload = await _client.GetFromJsonAsync<UserResponse[]>("/api/identity/tenants/bootstrap-ops/users");

    Assert.NotNull(usersPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/users/{usersPayload![0].PublicId}/roles",
      new AssignUserRoleRequest("   "));

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.Equal("invalid_role_code", payload.Code);
  }

  [Fact]
  public async Task InviteAcceptLoginRefreshAndAccessResolveShouldWorkEndToEnd()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity Access Flow");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"invite.flow@{tenantPayload.Slug}.local",
        "Invite Flow User",
        ["viewer"],
        null,
        7));

    Assert.Equal(HttpStatusCode.Created, inviteResponse.StatusCode);

    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);
    Assert.False(string.IsNullOrWhiteSpace(invitePayload!.InviteToken));
    Assert.Contains("viewer", invitePayload.RoleCodes);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload.InviteToken}/accept",
      new AcceptInviteRequest("Invite Flow User", "Invite", "Flow", "PhaseTwo123"));

    Assert.Equal(HttpStatusCode.OK, acceptResponse.StatusCode);

    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);
    Assert.Equal("accepted", acceptedPayload!.InviteStatus);
    Assert.Equal($"invite.flow@{tenantPayload.Slug}.local", acceptedPayload.User.Email);

    var loginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        $"invite.flow@{tenantPayload.Slug}.local",
        "PhaseTwo123",
        null));

    Assert.Equal(HttpStatusCode.OK, loginResponse.StatusCode);

    var sessionPayload = await loginResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(sessionPayload);
    Assert.False(string.IsNullOrWhiteSpace(sessionPayload!.SessionToken));
    Assert.False(string.IsNullOrWhiteSpace(sessionPayload.RefreshToken));
    Assert.Contains("viewer", sessionPayload.RoleCodes);

    var accessRequest = new HttpRequestMessage(HttpMethod.Get, $"/api/identity/tenants/{tenantPayload.Slug}/access");
    accessRequest.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", sessionPayload.SessionToken);
    var accessResponse = await _client.SendAsync(accessRequest);

    Assert.Equal(HttpStatusCode.OK, accessResponse.StatusCode);

    var accessPayload = await accessResponse.Content.ReadFromJsonAsync<AccessResolutionResponse>();

    Assert.NotNull(accessPayload);
    Assert.True(accessPayload!.Authorized);
    Assert.Equal(sessionPayload.UserPublicId, accessPayload.UserPublicId);
    Assert.Contains("viewer", accessPayload.RoleCodes);

    var refreshResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/refresh",
      new RefreshSessionRequest(sessionPayload.RefreshToken));

    Assert.Equal(HttpStatusCode.OK, refreshResponse.StatusCode);

    var refreshedPayload = await refreshResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(refreshedPayload);
    Assert.NotEqual(sessionPayload.RefreshToken, refreshedPayload!.RefreshToken);
  }

  [Fact]
  public async Task MfaShouldRequireOtpAfterEnrollment()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity MFA Flow");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"mfa.flow@{tenantPayload.Slug}.local",
        "MFA Flow User",
        ["viewer"],
        null,
        7));

    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("MFA Flow User", "MFA", "Flow", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    var enrollResponse = await _client.PostAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{acceptedPayload!.User.PublicId}/mfa/enroll",
      content: null);

    Assert.Equal(HttpStatusCode.OK, enrollResponse.StatusCode);

    var enrollmentPayload = await enrollResponse.Content.ReadFromJsonAsync<MfaEnrollmentResponse>();

    Assert.NotNull(enrollmentPayload);
    Assert.False(enrollmentPayload!.Enabled);

    var otpCode = ComputeTotp(enrollmentPayload.Secret);
    var verifyResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{acceptedPayload.User.PublicId}/mfa/verify",
      new VerifyMfaRequest(otpCode));

    Assert.Equal(HttpStatusCode.OK, verifyResponse.StatusCode);

    var loginWithoutOtpResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        "PhaseTwo123",
        null));

    Assert.Equal(HttpStatusCode.Unauthorized, loginWithoutOtpResponse.StatusCode);

    var loginWithoutOtpPayload = await loginWithoutOtpResponse.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(loginWithoutOtpPayload);
    Assert.Equal("mfa_required", loginWithoutOtpPayload!.Code);

    var loginWithOtpResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        "PhaseTwo123",
        ComputeTotp(enrollmentPayload.Secret)));

    Assert.Equal(HttpStatusCode.OK, loginWithOtpResponse.StatusCode);
  }

  [Fact]
  public async Task BlockedUserShouldLoseSessionAccess()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity Block Flow");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"blocked.flow@{tenantPayload.Slug}.local",
        "Blocked Flow User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Blocked Flow User", "Blocked", "Flow", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    var loginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload!.User.Email,
        "PhaseTwo123",
        null));
    var sessionPayload = await loginResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(sessionPayload);

    var blockResponse = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{acceptedPayload.User.PublicId}/access",
      new UpdateUserAccessRequest("suspended"));

    Assert.Equal(HttpStatusCode.OK, blockResponse.StatusCode);

    var loginAfterBlockResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        "PhaseTwo123",
        null));

    Assert.Equal(HttpStatusCode.Forbidden, loginAfterBlockResponse.StatusCode);

    var accessRequest = new HttpRequestMessage(HttpMethod.Get, $"/api/identity/tenants/{tenantPayload.Slug}/access");
    accessRequest.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", sessionPayload!.SessionToken);
    var accessResponse = await _client.SendAsync(accessRequest);

    Assert.Equal(HttpStatusCode.Unauthorized, accessResponse.StatusCode);

    var auditResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload.Slug}/security/audit");
    var auditPayload = await auditResponse.Content.ReadFromJsonAsync<SecurityAuditEventResponse[]>();

    Assert.NotNull(auditPayload);
    Assert.Contains(auditPayload!, auditEvent => auditEvent.EventCode == "access_blocked");
  }

  [Fact]
  public async Task PasswordRecoveryShouldRotateCredentialsAndRevokeSessions()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity Password Recovery");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"password.flow@{tenantPayload.Slug}.local",
        "Password Flow User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Password Flow User", "Password", "Flow", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    var loginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload!.User.Email,
        "PhaseTwo123",
        null));
    var sessionPayload = await loginResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(sessionPayload);

    var startRecoveryResponse = await _client.PostAsJsonAsync(
      "/api/identity/password-recovery",
      new StartPasswordRecoveryRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        45));

    Assert.Equal(HttpStatusCode.OK, startRecoveryResponse.StatusCode);

    var recoveryPayload = await startRecoveryResponse.Content.ReadFromJsonAsync<PasswordRecoveryResponse>();

    Assert.NotNull(recoveryPayload);
    Assert.False(string.IsNullOrWhiteSpace(recoveryPayload!.ResetToken));

    var completeRecoveryResponse = await _client.PostAsJsonAsync(
      $"/api/identity/password-recovery/{recoveryPayload.ResetToken}/complete",
      new ResetPasswordRequest("PhaseThree123"));

    Assert.Equal(HttpStatusCode.OK, completeRecoveryResponse.StatusCode);

    var accessRequest = new HttpRequestMessage(HttpMethod.Get, $"/api/identity/tenants/{tenantPayload.Slug}/access");
    accessRequest.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", sessionPayload!.SessionToken);
    var staleAccessResponse = await _client.SendAsync(accessRequest);

    Assert.Equal(HttpStatusCode.Unauthorized, staleAccessResponse.StatusCode);

    var staleLoginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        "PhaseTwo123",
        null));

    Assert.Equal(HttpStatusCode.Unauthorized, staleLoginResponse.StatusCode);

    var freshLoginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        "PhaseThree123",
        null));

    Assert.Equal(HttpStatusCode.OK, freshLoginResponse.StatusCode);
  }

  [Fact]
  public async Task UserSessionsEndpointsShouldListAndRevokeLiveSessions()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity Session Governance");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"session.flow@{tenantPayload.Slug}.local",
        "Session Flow User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Session Flow User", "Session", "Flow", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload!.User.Email,
        "PhaseTwo123",
        null));
    var secondLoginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload.User.Email,
        "PhaseTwo123",
        null));
    var secondSession = await secondLoginResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(secondSession);

    var listResponse = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{acceptedPayload.User.PublicId}/sessions");

    Assert.Equal(HttpStatusCode.OK, listResponse.StatusCode);

    var listedSessions = await listResponse.Content.ReadFromJsonAsync<UserSessionResponse[]>();

    Assert.NotNull(listedSessions);
    Assert.Equal(2, listedSessions!.Length);
    Assert.Contains(listedSessions, session => session.Status == "active");

    var revokeResponse = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/sessions/{secondSession!.PublicId}");

    Assert.Equal(HttpStatusCode.OK, revokeResponse.StatusCode);

    var revokeAllResponse = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{acceptedPayload.User.PublicId}/sessions");

    Assert.Equal(HttpStatusCode.OK, revokeAllResponse.StatusCode);

    var revokedSessions = await revokeAllResponse.Content.ReadFromJsonAsync<UserSessionResponse[]>();

    Assert.NotNull(revokedSessions);
    Assert.All(revokedSessions!, session => Assert.Equal("revoked", session.Status));
  }

  [Fact]
  public async Task InviteGovernanceEndpointsShouldResendAndCancelPendingInvite()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity Invite Governance");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"invite.governance@{tenantPayload.Slug}.local",
        "Invite Governance User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var resendResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/invites/{invitePayload!.PublicId}/resend",
      new ResendInviteRequest(14));

    Assert.Equal(HttpStatusCode.OK, resendResponse.StatusCode);

    var resentInvite = await resendResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(resentInvite);
    Assert.Equal(invitePayload.PublicId, resentInvite!.PublicId);
    Assert.Equal("pending", resentInvite.Status);
    Assert.NotEqual(invitePayload.InviteToken, resentInvite.InviteToken);

    var oldAcceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload.InviteToken}/accept",
      new AcceptInviteRequest("Invite Governance User", "Invite", "Governance", "PhaseTwo123"));

    Assert.Equal(HttpStatusCode.NotFound, oldAcceptResponse.StatusCode);

    var cancelResponse = await _client.PostAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/invites/{invitePayload.PublicId}/cancel",
      content: null);

    Assert.Equal(HttpStatusCode.OK, cancelResponse.StatusCode);

    var cancelledInvite = await cancelResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(cancelledInvite);
    Assert.Equal("revoked", cancelledInvite!.Status);

    var acceptAfterCancelResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{resentInvite.InviteToken}/accept",
      new AcceptInviteRequest("Invite Governance User", "Invite", "Governance", "PhaseTwo123"));

    Assert.Equal(HttpStatusCode.Conflict, acceptAfterCancelResponse.StatusCode);
  }

  [Fact]
  public async Task LogoutEndpointShouldRevokeCurrentSession()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Identity Logout");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"logout.flow@{tenantPayload.Slug}.local",
        "Logout Flow User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Logout Flow User", "Logout", "Flow", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    var loginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload!.User.Email,
        "PhaseTwo123",
        null));
    var loginPayload = await loginResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(loginPayload);

    using var logoutRequest = new HttpRequestMessage(HttpMethod.Post, "/api/identity/sessions/logout");
    logoutRequest.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", loginPayload!.SessionToken);

    var logoutResponse = await _client.SendAsync(logoutRequest);

    Assert.Equal(HttpStatusCode.OK, logoutResponse.StatusCode);

    var logoutPayload = await logoutResponse.Content.ReadFromJsonAsync<UserSessionResponse>();

    Assert.NotNull(logoutPayload);
    Assert.Equal("revoked", logoutPayload!.Status);
    Assert.NotNull(logoutPayload.RevokedAt);

    using var accessRequest = new HttpRequestMessage(HttpMethod.Get, $"/api/identity/tenants/{tenantPayload.Slug}/access");
    accessRequest.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", loginPayload.SessionToken);

    var accessResponse = await _client.SendAsync(accessRequest);

    Assert.Equal(HttpStatusCode.Unauthorized, accessResponse.StatusCode);

    var errorPayload = await accessResponse.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(errorPayload);
    Assert.Equal("invalid_session", errorPayload!.Code);
  }

  private static string ComputeTotp(string secret)
  {
    var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";
    var normalized = secret.Trim().Replace("=", string.Empty, StringComparison.Ordinal).ToUpperInvariant();
    var buffer = 0;
    var bitsLeft = 0;
    var bytes = new List<byte>();

    foreach (var character in normalized)
    {
      var index = alphabet.IndexOf(character);
      if (index < 0)
      {
        continue;
      }

      buffer = (buffer << 5) | index;
      bitsLeft += 5;

      if (bitsLeft < 8)
      {
        continue;
      }

      bytes.Add((byte)((buffer >> (bitsLeft - 8)) & 0xff));
      bitsLeft -= 8;
    }

    var counter = DateTimeOffset.UtcNow.ToUnixTimeSeconds() / 30;
    Span<byte> counterBytes = stackalloc byte[8];
    for (var index = 7; index >= 0; index--)
    {
      counterBytes[index] = (byte)(counter & 0xff);
      counter >>= 8;
    }

    using var hmac = new HMACSHA1(bytes.ToArray());
    var hash = hmac.ComputeHash(counterBytes.ToArray());
    var offset = hash[^1] & 0x0f;
    var binaryCode =
      ((hash[offset] & 0x7f) << 24)
      | (hash[offset + 1] << 16)
      | (hash[offset + 2] << 8)
      | hash[offset + 3];

    return (binaryCode % 1_000_000).ToString("D6");
  }
}
