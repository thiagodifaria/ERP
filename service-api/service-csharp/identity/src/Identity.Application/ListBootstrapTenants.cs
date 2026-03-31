// Este caso de uso expõe a leitura minima de tenants para smoke e bootstrap do servico.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapTenants
{
  private readonly ITenantRepository _tenantRepository;

  public ListBootstrapTenants(ITenantRepository tenantRepository)
  {
    _tenantRepository = tenantRepository;
  }

  public IReadOnlyCollection<TenantResponse> Execute()
  {
    return _tenantRepository
      .List()
      .Select(tenant => new TenantResponse(
        tenant.Id,
        tenant.PublicId,
        tenant.Slug,
        tenant.DisplayName,
        tenant.Status))
      .ToArray();
  }
}
