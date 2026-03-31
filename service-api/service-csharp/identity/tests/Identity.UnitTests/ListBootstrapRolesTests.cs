// Estes testes cobrem a leitura de papeis basicos por tenant no bootstrap.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapRolesTests
{
  [Fact]
  public void ExecuteShouldReturnRolesForKnownTenant()
  {
    var useCase = new ListBootstrapRoles(new InMemoryRoleRepository());

    var response = useCase.Execute("bootstrap-ops");

    Assert.Equal(5, response.Count);
    Assert.Contains(response, role => role.Code == "owner");
    Assert.Contains(response, role => role.Code == "viewer");
  }

  [Fact]
  public void ExecuteShouldReturnEmptyForUnknownTenant()
  {
    var useCase = new ListBootstrapRoles(new InMemoryRoleRepository());

    var response = useCase.Execute("missing-tenant");

    Assert.Empty(response);
  }
}
