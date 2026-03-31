// Este contrato define a leitura minima de usuarios por tenant.
namespace Identity.Domain;

public interface IUserCatalog
{
  IReadOnlyCollection<User> ListByTenantId(long tenantId);

  User? FindByTenantIdAndId(long tenantId, long userId);

  User? FindByTenantIdAndPublicId(long tenantId, Guid publicId);

  User? FindByTenantIdAndEmail(long tenantId, string email);
}
