// Este contrato define a escrita minima de papeis basicos durante o bootstrap.
namespace Identity.Domain;

public interface IRoleRepository : IRoleCatalog
{
  IReadOnlyCollection<Role> SeedDefaults(Tenant tenant);
}
