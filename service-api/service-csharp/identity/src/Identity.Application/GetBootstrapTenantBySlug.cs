// Este caso de uso expõe a leitura individual de tenant por slug no bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class GetBootstrapTenantBySlug
{
  private readonly ITenantRepository _tenantRepository;

  public GetBootstrapTenantBySlug(ITenantRepository tenantRepository)
  {
    _tenantRepository = tenantRepository;
  }

  public TenantResponse? Execute(string slug)
  {
    var tenant = _tenantRepository.FindBySlug(slug);
    if (tenant is null)
    {
      return null;
    }

    return new TenantResponse(
      tenant.Id,
      tenant.PublicId,
      tenant.Slug,
      tenant.DisplayName,
      tenant.Status);
  }
}
