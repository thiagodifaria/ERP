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
}
