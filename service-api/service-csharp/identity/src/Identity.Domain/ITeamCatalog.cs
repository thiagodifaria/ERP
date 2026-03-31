// Este contrato define a leitura minima de times por tenant.
namespace Identity.Domain;

public interface ITeamCatalog
{
  IReadOnlyCollection<Team> ListByTenantId(long tenantId);

  Team? FindByTenantIdAndName(long tenantId, string name);

  Team? FindByTenantIdAndPublicId(long tenantId, Guid publicId);
}
