// Este contrato define a leitura de papeis de acesso durante o bootstrap do servico.
namespace Identity.Domain;

public interface IRoleCatalog
{
  IReadOnlyCollection<Role> ListByTenantSlug(string tenantSlug);
}
