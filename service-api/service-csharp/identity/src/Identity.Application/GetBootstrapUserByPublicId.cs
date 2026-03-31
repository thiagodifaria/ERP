// Este caso de uso expoe a leitura individual de user por publicId no bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class GetBootstrapUserByPublicId
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;

  public GetBootstrapUserByPublicId(ITenantCatalog tenantCatalog, IUserCatalog userCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
  }

  public UserResponse? Execute(string tenantSlug, Guid userPublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);

    if (user is null)
    {
      return null;
    }

    return new UserResponse(
      user.Id,
      user.PublicId,
      user.TenantId,
      user.CompanyId,
      user.Email,
      user.DisplayName,
      user.GivenName,
      user.FamilyName,
      user.Status);
  }
}
