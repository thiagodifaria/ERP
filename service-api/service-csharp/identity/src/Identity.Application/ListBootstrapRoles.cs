// Este caso de uso expõe a leitura minima de papeis por tenant em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapRoles
{
  private readonly IRoleCatalog _roleCatalog;

  public ListBootstrapRoles(IRoleCatalog roleCatalog)
  {
    _roleCatalog = roleCatalog;
  }

  public IReadOnlyCollection<RoleResponse> Execute(string tenantSlug)
  {
    return _roleCatalog
      .ListByTenantSlug(tenantSlug)
      .Select(role => new RoleResponse(
        role.Id,
        role.PublicId,
        role.Code,
        role.DisplayName,
        role.Status))
      .ToArray();
  }
}
