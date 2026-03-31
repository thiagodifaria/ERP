// Este contrato define a escrita minima de atribuicoes de papeis durante o bootstrap.
namespace Identity.Domain;

public interface IUserRoleRepository : IUserRoleCatalog
{
  UserRole Add(UserRole userRole);

  long NextId();
}
