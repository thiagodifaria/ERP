// Este caso de uso expoe a leitura minima de times por tenant em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapTeams
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ITeamCatalog _teamCatalog;

  public ListBootstrapTeams(ITenantCatalog tenantCatalog, ITeamCatalog teamCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _teamCatalog = teamCatalog;
  }

  public IReadOnlyCollection<TeamResponse>? Execute(string tenantSlug)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    return _teamCatalog
      .ListByTenantId(tenant.Id)
      .Select(team => new TeamResponse(
        team.Id,
        team.PublicId,
        team.TenantId,
        team.CompanyId,
        team.Name,
        team.Status))
      .ToArray();
  }
}
