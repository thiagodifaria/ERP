// Estes testes cobrem a leitura de usuarios basicos por tenant no bootstrap.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapUsersTests
{
  [Fact]
  public void ExecuteShouldReturnUsersForKnownTenant()
  {
    var useCase = new ListBootstrapUsers(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var response = useCase.Execute("bootstrap-ops");

    Assert.NotNull(response);
    Assert.NotEmpty(response);
    Assert.Contains(response!, user => user.Email == "owner@bootstrap-ops.local");
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownTenant()
  {
    var useCase = new ListBootstrapUsers(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var response = useCase.Execute("missing-tenant");

    Assert.Null(response);
  }
}
