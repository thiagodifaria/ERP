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
