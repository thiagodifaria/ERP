// Este caso de uso expoe a leitura minima de membros de um time em bootstrap.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListBootstrapTeamMembers
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ITeamCatalog _teamCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly ITeamMembershipCatalog _teamMembershipCatalog;

  public ListBootstrapTeamMembers(
    ITenantCatalog tenantCatalog,
    ITeamCatalog teamCatalog,
    IUserCatalog userCatalog,
    ITeamMembershipCatalog teamMembershipCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _teamCatalog = teamCatalog;
    _userCatalog = userCatalog;
    _teamMembershipCatalog = teamMembershipCatalog;
  }

  public IReadOnlyCollection<TeamMembershipResponse>? Execute(string tenantSlug, Guid teamPublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    var team = _teamCatalog.FindByTenantIdAndPublicId(tenant.Id, teamPublicId);

    if (team is null)
    {
      return null;
    }

    return _teamMembershipCatalog
      .ListByTenantIdAndTeamId(tenant.Id, team.Id)
      .Select(membership =>
      {
        var user = _userCatalog.FindByTenantIdAndId(tenant.Id, membership.UserId);

        return user is null
          ? null
          : new TeamMembershipResponse(
            membership.Id,
            membership.TenantId,
            team.PublicId,
            user.PublicId,
            user.Email,
            user.DisplayName,
            membership.CreatedAt);
      })
      .Where(response => response is not null)
      .Select(response => response!)
      .ToArray();
  }
}
