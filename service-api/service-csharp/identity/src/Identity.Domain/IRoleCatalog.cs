// Este contrato define a leitura de papeis de acesso durante o bootstrap do servico.
namespace Identity.Domain;

public interface IRoleCatalog
{
  IReadOnlyCollection<Role> ListByTenantId(long tenantId);

  IReadOnlyCollection<Role> ListByTenantSlug(string tenantSlug);

  Role? FindByTenantIdAndId(long tenantId, long roleId);

  Role? FindByTenantIdAndCode(long tenantId, string code);
}
