// Estes testes protegem os contratos publicos mais importantes da API de identidade.
using System.Net;
using System.Net.Http.Json;
using Identity.Contracts;
using Microsoft.AspNetCore.Mvc.Testing;
using Xunit;

namespace Identity.ContractTests;

public sealed class IdentityContractsTests : IClassFixture<WebApplicationFactory<Program>>
{
  private readonly HttpClient _client;

  public IdentityContractsTests(WebApplicationFactory<Program> factory)
  {
    _client = factory.CreateClient();
  }

  [Fact]
  public async Task TenantListContractShouldExposePublicFields()
  {
    var response = await _client.GetAsync("/api/identity/tenants");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.All(payload, tenant =>
    {
      Assert.NotEqual(0, tenant.Id);
      Assert.NotEqual(Guid.Empty, tenant.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(tenant.Slug));
      Assert.False(string.IsNullOrWhiteSpace(tenant.DisplayName));
      Assert.False(string.IsNullOrWhiteSpace(tenant.Status));
    });
  }

  [Fact]
  public async Task CreateTenantContractShouldReturnCreatedResource()
  {
    var request = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Tenant");

    var response = await _client.PostAsJsonAsync("/api/identity/tenants", request);

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(payload);
    Assert.Equal(request.Slug, payload!.Slug);
    Assert.Equal(request.DisplayName, payload.DisplayName);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task SnapshotContractShouldExposeConsistentNestedStructure()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/snapshot");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantAccessSnapshotResponse>();

    Assert.NotNull(payload);
    Assert.Equal("bootstrap-ops", payload!.Tenant.Slug);
    Assert.True(payload.Counts.Companies >= 1);
    Assert.True(payload.Counts.Users >= 1);
    Assert.True(payload.Counts.Teams >= 1);
    Assert.True(payload.Counts.Roles >= 1);
    Assert.All(payload.Users, user =>
    {
      Assert.NotEqual(Guid.Empty, user.User.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(user.User.Email));
      Assert.All(user.Roles, role => Assert.False(string.IsNullOrWhiteSpace(role.RoleCode)));
    });
    Assert.All(payload.Teams, team =>
    {
      Assert.NotEqual(Guid.Empty, team.Team.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(team.Team.Name));
      Assert.All(team.Members, member => Assert.NotEqual(Guid.Empty, member.UserPublicId));
    });
  }

  [Fact]
  public async Task UpdateCompanyContractShouldReturnUpdatedResourceShape()
  {
    var companies = await _client.GetFromJsonAsync<CompanyResponse[]>("/api/identity/tenants/bootstrap-ops/companies");

    Assert.NotNull(companies);
    Assert.NotEmpty(companies);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/companies/{companies![0].PublicId}",
      new UpdateCompanyRequest("Bootstrap Ops Contract", "Bootstrap Ops Contract LTDA", "12345678901234"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal("Bootstrap Ops Contract", payload.DisplayName);
    Assert.Equal("Bootstrap Ops Contract LTDA", payload.LegalName);
    Assert.Equal("12345678901234", payload.TaxId);
  }

  [Fact]
  public async Task CompanyDetailContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Company Detail");

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
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal(tenantPayload.DisplayName, payload.DisplayName);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task UpdateUserContractShouldReturnUpdatedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract User Update");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"contract.user@{tenantPayload.Slug}.local",
        "Contract User",
        "Contract",
        "User"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}",
      new UpdateUserRequest(
        $"contract.user.prime@{tenantPayload.Slug}.local",
        "Contract User Prime",
        "Contract",
        "Prime"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal($"contract.user.prime@{tenantPayload.Slug}.local", payload.Email);
    Assert.Equal("Contract User Prime", payload.DisplayName);
    Assert.Equal("Prime", payload.FamilyName);
  }

  [Fact]
  public async Task UserDetailContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract User Detail");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"contract.user.detail@{tenantPayload.Slug}.local",
        "Contract User Detail",
        "Contract",
        "Detail"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal($"contract.user.detail@{tenantPayload.Slug}.local", payload.Email);
    Assert.Equal("Contract User Detail", payload.DisplayName);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task RevokeUserRoleContractShouldReturnRevokedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Role Revoke");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"contract.revoke@{tenantPayload.Slug}.local",
        "Contract Revoke User",
        "Contract",
        "Revoke"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}/roles",
      new AssignUserRoleRequest("viewer"));

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser.PublicId}/roles/viewer");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserRoleResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.UserPublicId);
    Assert.NotEqual(Guid.Empty, payload.RolePublicId);
    Assert.Equal("viewer", payload.RoleCode);
    Assert.False(string.IsNullOrWhiteSpace(payload.RoleDisplayName));
  }

  [Fact]
  public async Task UpdateTeamContractShouldReturnUpdatedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Team Update");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teams = await _client.GetFromJsonAsync<TeamResponse[]>(
      $"/api/identity/tenants/{tenantPayload!.Slug}/teams");

    Assert.NotNull(teams);
    Assert.NotEmpty(teams);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teams![0].PublicId}",
      new UpdateTeamRequest("Core Contract Prime"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal("Core Contract Prime", payload.Name);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task RemoveTeamMemberContractShouldReturnRemovedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Membership Removal");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teamsResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/teams");
    var teamsPayload = await teamsResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamsPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users",
      new CreateUserRequest(
        $"contract.membership.remove@{tenantPayload.Slug}.local",
        "Contract Membership Remove",
        "Contract",
        "Remove"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamsPayload![0].PublicId}/members",
      new AddTeamMemberRequest(createdUser!.PublicId));

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamsPayload[0].PublicId}/members/{createdUser.PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamMembershipResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.TeamPublicId);
    Assert.NotEqual(Guid.Empty, payload.UserPublicId);
    Assert.False(string.IsNullOrWhiteSpace(payload.UserEmail));
    Assert.False(string.IsNullOrWhiteSpace(payload.UserDisplayName));
  }

  [Fact]
  public async Task ErrorContractShouldExposeCodeAndMessage()
  {
    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/users",
      new CreateUserRequest("invalid-email", "Invalid User", null, null));

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.False(string.IsNullOrWhiteSpace(payload!.Code));
    Assert.False(string.IsNullOrWhiteSpace(payload.Message));
  }
}
