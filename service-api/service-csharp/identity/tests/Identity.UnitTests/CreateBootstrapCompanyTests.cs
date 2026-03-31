// Estes testes cobrem a criacao minima de empresas no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class CreateBootstrapCompanyTests
{
  [Fact]
  public void ExecuteShouldCreateCompanyForKnownTenant()
  {
    var useCase = new CreateBootstrapCompany(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateCompanyRequest("Bootstrap Ops Filial", "Bootstrap Ops Filial LTDA", "12345678901234"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Company);
    Assert.Equal("Bootstrap Ops Filial", result.Company!.DisplayName);
    Assert.Equal("Bootstrap Ops Filial LTDA", result.Company.LegalName);
    Assert.Equal("12345678901234", result.Company.TaxId);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownTenant()
  {
    var useCase = new CreateBootstrapCompany(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var result = useCase.Execute(
      "missing-tenant",
      new CreateCompanyRequest("New Company", null, null));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("tenant_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateDisplayNamePerTenant()
  {
    var useCase = new CreateBootstrapCompany(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateCompanyRequest("Bootstrap Ops", null, null));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("company_display_name_conflict", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectBlankDisplayName()
  {
    var useCase = new CreateBootstrapCompany(
      new InMemoryTenantRepository(),
      new InMemoryCompanyRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateCompanyRequest("   ", null, null));

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_display_name", result.Error!.Code);
  }
}
