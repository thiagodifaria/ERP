// Este contrato define a escrita e leitura minima de tenants no contexto de identidade.
namespace Identity.Domain;

public interface ITenantRepository : ITenantCatalog
{
  Tenant? FindBySlug(string slug);

  Tenant Add(Tenant tenant);

  long NextId();
}
