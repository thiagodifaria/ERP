// Estes testes validam o comportamento HTTP minimo da API de identidade.
using System.Net;
using System.Net.Http.Json;
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
  public async Task TenantCompaniesEndpointShouldReturnCompaniesForKnownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/companies");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.Contains(payload, company => company.DisplayName == "Bootstrap Ops");
  }

  [Fact]
  public async Task TenantCompaniesEndpointShouldReturnNotFoundForUnknownTenant()
  {
    var response = await _client.GetAsync("/api/identity/tenants/missing-tenant/companies");

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
}
