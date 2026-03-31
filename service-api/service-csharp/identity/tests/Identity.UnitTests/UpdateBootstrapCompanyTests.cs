// Estes testes cobrem o update minimo de empresas no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Domain;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class UpdateBootstrapCompanyTests
{
  [Fact]
  public void ExecuteShouldUpdateCompanyForKnownTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    var existingCompany = companyRepository.ListByTenantId(1).Single();
    var useCase = new UpdateBootstrapCompany(tenantRepository, companyRepository);

    var result = useCase.Execute(
      "bootstrap-ops",
      existingCompany.PublicId,
      new UpdateCompanyRequest("Bootstrap Ops Prime", "Bootstrap Ops Prime LTDA", "12345678901234"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Company);
    Assert.Equal("Bootstrap Ops Prime", result.Company!.DisplayName);
    Assert.Equal("Bootstrap Ops Prime LTDA", result.Company.LegalName);
    Assert.Equal("12345678901234", result.Company.TaxId);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownCompany()
  {
    var useCase = new UpdateBootstrapCompany(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      Guid.NewGuid(),
      new UpdateCompanyRequest("Bootstrap Ops Prime", null, null));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("company_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateDisplayNamePerTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    companyRepository.Add(new Company(
      companyRepository.NextId(),
      1,
      PublicIds.NewUuidV7(),
      "Bootstrap Ops Filial",
      null,
      null,
      "active"));
    var existingCompany = companyRepository.ListByTenantId(1).Single(company => company.DisplayName == "Bootstrap Ops");
    var useCase = new UpdateBootstrapCompany(tenantRepository, companyRepository);

    var result = useCase.Execute(
      "bootstrap-ops",
      existingCompany.PublicId,
      new UpdateCompanyRequest("Bootstrap Ops Filial", null, null));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("company_display_name_conflict", result.Error!.Code);
  }
}
