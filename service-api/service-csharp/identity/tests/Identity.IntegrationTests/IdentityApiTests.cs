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
}
