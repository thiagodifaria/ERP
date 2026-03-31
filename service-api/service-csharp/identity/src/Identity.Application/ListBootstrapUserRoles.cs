// Este caso de uso expoe a leitura minima de papeis atribuidos a um usuario em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapUserRoles
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IRoleCatalog _roleCatalog;
  private readonly IUserRoleCatalog _userRoleCatalog;

  public ListBootstrapUserRoles(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IRoleCatalog roleCatalog,
    IUserRoleCatalog userRoleCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _roleCatalog = roleCatalog;
    _userRoleCatalog = userRoleCatalog;
  }

  public IReadOnlyCollection<UserRoleResponse>? Execute(string tenantSlug, Guid userPublicId)
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

    return _userRoleCatalog
      .ListByTenantIdAndUserId(tenant.Id, user.Id)
      .Select(userRole =>
      {
        var role = _roleCatalog.FindByTenantIdAndId(tenant.Id, userRole.RoleId);

        return role is null
          ? null
          : new UserRoleResponse(
            userRole.Id,
            userRole.TenantId,
            user.PublicId,
            role.PublicId,
            role.Code,
            role.DisplayName,
            userRole.CreatedAt);
      })
      .Where(response => response is not null)
      .Select(response => response!)
      .ToArray();
  }
}
