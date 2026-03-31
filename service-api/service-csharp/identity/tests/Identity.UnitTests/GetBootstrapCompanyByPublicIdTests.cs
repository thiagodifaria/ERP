// Estes testes cobrem a leitura individual de company por publicId.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class GetBootstrapCompanyByPublicIdTests
{
  [Fact]
  public void ExecuteShouldReturnCompanyForKnownPublicId()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    var useCase = new GetBootstrapCompanyByPublicId(tenantRepository, companyRepository);
    var company = companyRepository.ListByTenantId(1).First();

    var response = useCase.Execute("bootstrap-ops", company.PublicId);

    Assert.NotNull(response);
    Assert.Equal(company.PublicId, response!.PublicId);
    Assert.Equal(company.DisplayName, response.DisplayName);
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownCompany()
  {
    var useCase = new GetBootstrapCompanyByPublicId(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var response = useCase.Execute("bootstrap-ops", Guid.NewGuid());

    Assert.Null(response);
  }
}
