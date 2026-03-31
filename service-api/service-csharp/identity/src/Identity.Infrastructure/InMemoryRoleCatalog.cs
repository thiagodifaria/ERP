// Este adapter fornece papeis basicos por tenant durante o bootstrap do servico.
// A leitura definitiva pode migrar para banco sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryRoleCatalog : IRoleCatalog
{
  private static readonly IReadOnlyDictionary<string, IReadOnlyCollection<Role>> RolesByTenantSlug =
    new Dictionary<string, IReadOnlyCollection<Role>>(StringComparer.OrdinalIgnoreCase)
    {
      ["bootstrap-ops"] =
      [
        new Role(1, 1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000011"), "owner", "Owner", "active"),
        new Role(2, 1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000012"), "admin", "Administrator", "active"),
        new Role(3, 1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000013"), "manager", "Manager", "active"),
        new Role(4, 1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000014"), "operator", "Operator", "active"),
        new Role(5, 1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000015"), "viewer", "Viewer", "active")
      ],
      ["northwind-group"] =
      [
        new Role(6, 2, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000021"), "admin", "Administrator", "active"),
        new Role(7, 2, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000022"), "manager", "Manager", "active"),
        new Role(8, 2, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000023"), "viewer", "Viewer", "active")
      ]
    };

  public IReadOnlyCollection<Role> ListByTenantSlug(string tenantSlug)
  {
    if (RolesByTenantSlug.TryGetValue(tenantSlug, out var roles))
    {
      return roles;
    }

    return [];
  }
}
