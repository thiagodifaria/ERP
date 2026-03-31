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
}
