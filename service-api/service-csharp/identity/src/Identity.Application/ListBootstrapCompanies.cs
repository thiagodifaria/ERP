// Este caso de uso expoe a leitura minima de empresas por tenant em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapCompanies
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ICompanyCatalog _companyCatalog;

  public ListBootstrapCompanies(ITenantCatalog tenantCatalog, ICompanyCatalog companyCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _companyCatalog = companyCatalog;
  }

  public IReadOnlyCollection<CompanyResponse>? Execute(string tenantSlug)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    return _companyCatalog
      .ListByTenantId(tenant.Id)
      .Select(company => new CompanyResponse(
        company.Id,
        company.PublicId,
        company.TenantId,
        company.DisplayName,
        company.LegalName,
        company.TaxId,
        company.Status))
      .ToArray();
  }
}
