// Este contrato define a leitura minima de tenants durante o bootstrap do servico.
// Persistencia real pode evoluir sem alterar a API publica basica.
namespace Identity.Domain;

public interface ITenantCatalog
{
  IReadOnlyCollection<Tenant> List();
}
