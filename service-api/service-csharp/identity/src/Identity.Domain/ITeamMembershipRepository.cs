// Este contrato define a escrita minima de memberships de time durante o bootstrap.
namespace Identity.Domain;

public interface ITeamMembershipRepository : ITeamMembershipCatalog
{
  TeamMembership Add(TeamMembership membership);

  long NextId();
}
