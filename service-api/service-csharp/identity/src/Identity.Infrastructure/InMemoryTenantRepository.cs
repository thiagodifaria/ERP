// Este adapter fornece bootstrap de tenants enquanto a persistencia real nao esta conectada.
// Ele permite leitura e escrita minima para validar os contratos publicos do servico.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryTenantRepository : ITenantRepository
{
  private readonly object _sync = new();
  private readonly List<Tenant> _tenants =
  [
    new Tenant(1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000001"), "bootstrap-ops", "Bootstrap Ops", "active"),
    new Tenant(2, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000002"), "northwind-group", "Northwind Group", "active")
  ];

  public IReadOnlyCollection<Tenant> List()
  {
    lock (_sync)
    {
      return _tenants.ToArray();
    }
  }

  public Tenant? FindBySlug(string slug)
  {
    lock (_sync)
    {
      return _tenants.FirstOrDefault(
        tenant => tenant.Slug.Equals(slug, StringComparison.OrdinalIgnoreCase));
    }
  }

  public Tenant Add(Tenant tenant)
  {
    lock (_sync)
    {
      _tenants.Add(tenant);
      return tenant;
    }
  }

  public long NextId()
  {
    lock (_sync)
    {
      return _tenants.Count == 0
        ? 1
        : _tenants.Max(tenant => tenant.Id) + 1;
    }
  }
}
