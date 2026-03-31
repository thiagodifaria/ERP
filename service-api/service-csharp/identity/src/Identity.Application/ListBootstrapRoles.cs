// Este caso de uso expõe a leitura minima de papeis por tenant em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapRoles
{
  private readonly IRoleRepository _roleRepository;

  public ListBootstrapRoles(IRoleRepository roleRepository)
  {
    _roleRepository = roleRepository;
  }

  public IReadOnlyCollection<RoleResponse> Execute(string tenantSlug)
  {
    return _roleRepository
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
