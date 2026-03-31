// Este contrato define a leitura minima de memberships de time por tenant.
namespace Identity.Domain;

public interface ITeamMembershipCatalog
{
  IReadOnlyCollection<TeamMembership> ListByTenantIdAndTeamId(long tenantId, long teamId);

  TeamMembership? FindByTenantIdAndTeamIdAndUserId(long tenantId, long teamId, long userId);
}
