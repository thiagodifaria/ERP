// Este contrato define a escrita minima de times durante o bootstrap.
namespace Identity.Domain;

public interface ITeamRepository : ITeamCatalog
{
  Team Add(Team team);

  long NextId();

  IReadOnlyCollection<Team> SeedDefaults(Tenant tenant);
}
