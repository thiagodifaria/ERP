// Este caso de uso expoe a leitura individual de company por publicId no bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class GetBootstrapCompanyByPublicId
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ICompanyCatalog _companyCatalog;

  public GetBootstrapCompanyByPublicId(ITenantCatalog tenantCatalog, ICompanyCatalog companyCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _companyCatalog = companyCatalog;
  }

  public CompanyResponse? Execute(string tenantSlug, Guid companyPublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    var company = _companyCatalog.FindByTenantIdAndPublicId(tenant.Id, companyPublicId);

    if (company is null)
    {
      return null;
    }

    return new CompanyResponse(
      company.Id,
      company.PublicId,
      company.TenantId,
      company.DisplayName,
      company.LegalName,
      company.TaxId,
      company.Status);
  }
}
