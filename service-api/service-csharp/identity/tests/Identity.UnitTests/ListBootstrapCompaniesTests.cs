// Estes testes cobrem a leitura de empresas basicas por tenant no bootstrap.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapCompaniesTests
{
  [Fact]
  public void ExecuteShouldReturnCompaniesForKnownTenant()
  {
    var useCase = new ListBootstrapCompanies(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var response = useCase.Execute("bootstrap-ops");

    Assert.NotNull(response);
    Assert.NotEmpty(response);
    Assert.Contains(response!, company => company.DisplayName == "Bootstrap Ops");
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownTenant()
  {
    var useCase = new ListBootstrapCompanies(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var response = useCase.Execute("missing-tenant");

    Assert.Null(response);
  }
}
