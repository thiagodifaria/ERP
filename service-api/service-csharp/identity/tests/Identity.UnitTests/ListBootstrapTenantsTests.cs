// Estes testes cobrem o fluxo basico de leitura de tenants no bootstrap do servico.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapTenantsTests
{
  [Fact]
  public void ExecuteShouldReturnBootstrapTenants()
  {
    var useCase = new ListBootstrapTenants(new InMemoryTenantRepository());

    var response = useCase.Execute();

    Assert.NotEmpty(response);
    Assert.Contains(response, tenant => tenant.Slug == "bootstrap-ops");
    Assert.Contains(response, tenant => tenant.Slug == "northwind-group");
  }
}
