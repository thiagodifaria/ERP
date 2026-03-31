// Este adapter fornece tenants de bootstrap enquanto a persistencia real nao esta conectada.
// Ele existe para validar contratos e fluxo da API sem esconder regra de negocio.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryTenantCatalog : ITenantCatalog
{
  private static readonly IReadOnlyCollection<Tenant> Tenants =
  [
    new Tenant(1, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000001"), "bootstrap-ops", "Bootstrap Ops", "active"),
    new Tenant(2, Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000002"), "northwind-group", "Northwind Group", "active")
  ];

  public IReadOnlyCollection<Tenant> List()
  {
    return Tenants;
  }
}
