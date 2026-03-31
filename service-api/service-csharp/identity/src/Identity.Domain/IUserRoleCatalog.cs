// Este contrato define a leitura minima de atribuicoes de papeis por usuario.
namespace Identity.Domain;

public interface IUserRoleCatalog
{
  IReadOnlyCollection<UserRole> ListByTenantIdAndUserId(long tenantId, long userId);

  UserRole? FindByTenantIdAndUserIdAndRoleId(long tenantId, long userId, long roleId);
}
