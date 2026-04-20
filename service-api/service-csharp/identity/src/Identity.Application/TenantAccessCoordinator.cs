// TenantAccessCoordinator centraliza a sincronizacao de acesso com o grafo externo.
using Identity.Domain;

namespace Identity.Application;

public sealed class TenantAccessCoordinator
{
  private readonly IRoleCatalog _roleCatalog;
  private readonly IUserRoleCatalog _userRoleCatalog;
  private readonly IAuthorizationGraph _authorizationGraph;

  public TenantAccessCoordinator(
    IRoleCatalog roleCatalog,
    IUserRoleCatalog userRoleCatalog,
    IAuthorizationGraph authorizationGraph)
  {
    _roleCatalog = roleCatalog;
    _userRoleCatalog = userRoleCatalog;
    _authorizationGraph = authorizationGraph;
  }

  public IReadOnlyCollection<string> SyncAndListRoleCodes(Tenant tenant, User user)
  {
    var roleLookup = _roleCatalog.ListByTenantId(tenant.Id).ToDictionary(role => role.Id, role => role.Code);
    var roleCodes = _userRoleCatalog.ListByTenantIdAndUserId(tenant.Id, user.Id)
      .Select(userRole => roleLookup.TryGetValue(userRole.RoleId, out var roleCode) ? roleCode : null)
      .Where(roleCode => !string.IsNullOrWhiteSpace(roleCode))
      .Select(roleCode => roleCode!)
      .Distinct(StringComparer.OrdinalIgnoreCase)
      .OrderBy(roleCode => roleCode, StringComparer.OrdinalIgnoreCase)
      .ToArray();

    _authorizationGraph.SyncTenantAccess(tenant.Slug, user.PublicId, roleCodes, user.Status == "active");

    return roleCodes;
  }

  public bool CanAccessTenant(Tenant tenant, User user)
  {
    return _authorizationGraph.CanAccessTenant(tenant.Slug, user.PublicId);
  }
}
