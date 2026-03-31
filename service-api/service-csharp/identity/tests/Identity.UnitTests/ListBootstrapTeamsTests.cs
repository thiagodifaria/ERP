// Estes testes cobrem a leitura de times basicos por tenant no bootstrap.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapTeamsTests
{
  [Fact]
  public void ExecuteShouldReturnTeamsForKnownTenant()
  {
    var useCase = new ListBootstrapTeams(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var response = useCase.Execute("bootstrap-ops");

    Assert.NotNull(response);
    Assert.NotEmpty(response);
    Assert.Contains(response!, team => team.Name == "Core");
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownTenant()
  {
    var useCase = new ListBootstrapTeams(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var response = useCase.Execute("missing-tenant");

    Assert.Null(response);
  }
}
