// Este adapter fornece papeis basicos por tenant durante o bootstrap do servico.
// A persistencia real pode substituir esta implementacao sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryRoleRepository : IRoleRepository
{
  private readonly object _sync = new();
  private readonly Dictionary<string, List<Role>> _rolesByTenantSlug =
    new(StringComparer.OrdinalIgnoreCase);

  public InMemoryRoleRepository()
  {
    SeedDefaults(new Tenant(
      1,
      Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000001"),
      "bootstrap-ops",
      "Bootstrap Ops",
      "active"));

    SeedDefaults(new Tenant(
      2,
      Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000002"),
      "northwind-group",
      "Northwind Group",
      "active"));
  }

  public IReadOnlyCollection<Role> ListByTenantSlug(string tenantSlug)
  {
    lock (_sync)
    {
      if (_rolesByTenantSlug.TryGetValue(tenantSlug, out var roles))
      {
        return roles.ToArray();
      }

      return [];
    }
  }

  public IReadOnlyCollection<Role> SeedDefaults(Tenant tenant)
  {
    lock (_sync)
    {
      if (_rolesByTenantSlug.TryGetValue(tenant.Slug, out var existingRoles) && existingRoles.Count > 0)
      {
        return existingRoles.ToArray();
      }

      var nextId = _rolesByTenantSlug
        .SelectMany(entry => entry.Value)
        .Select(role => role.Id)
        .DefaultIfEmpty(0)
        .Max() + 1;

      var seededRoles = DefaultRoles()
        .Select(roleSeed => new Role(
          nextId++,
          tenant.Id,
          PublicIds.NewUuidV7(),
          roleSeed.Code,
          roleSeed.DisplayName,
          "active"))
        .ToList();

      _rolesByTenantSlug[tenant.Slug] = seededRoles;

      return seededRoles.ToArray();
    }
  }

  private static IReadOnlyCollection<(string Code, string DisplayName)> DefaultRoles()
  {
    return
    [
      ("owner", "Owner"),
      ("admin", "Administrator"),
      ("manager", "Manager"),
      ("operator", "Operator"),
      ("viewer", "Viewer")
    ];
  }
}
