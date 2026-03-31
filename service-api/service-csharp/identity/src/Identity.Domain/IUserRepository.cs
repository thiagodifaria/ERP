// Este contrato define a escrita minima de usuarios durante o bootstrap.
namespace Identity.Domain;

public interface IUserRepository : IUserCatalog
{
  User Add(User user);

  long NextId();

  IReadOnlyCollection<User> SeedDefaults(Tenant tenant);
}
