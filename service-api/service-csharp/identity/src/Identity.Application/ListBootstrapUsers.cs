// Este caso de uso expoe a leitura minima de usuarios por tenant em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapUsers
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;

  public ListBootstrapUsers(ITenantCatalog tenantCatalog, IUserCatalog userCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
  }

  public IReadOnlyCollection<UserResponse>? Execute(string tenantSlug)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    return _userCatalog
      .ListByTenantId(tenant.Id)
      .Select(user => new UserResponse(
        user.Id,
        user.PublicId,
        user.TenantId,
        user.CompanyId,
        user.Email,
        user.DisplayName,
        user.GivenName,
        user.FamilyName,
        user.Status))
      .ToArray();
  }
}
