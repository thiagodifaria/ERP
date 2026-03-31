// Este contrato define a leitura minima de usuarios por tenant.
namespace Identity.Domain;

public interface IUserCatalog
{
  IReadOnlyCollection<User> ListByTenantId(long tenantId);

  User? FindByTenantIdAndEmail(long tenantId, string email);
}
