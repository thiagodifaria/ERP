// Este caso de uso remove memberships entre times e usuarios dentro do tenant.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class RemoveBootstrapTeamMember
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ITeamCatalog _teamCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly ITeamMembershipRepository _teamMembershipRepository;

  public RemoveBootstrapTeamMember(
    ITenantCatalog tenantCatalog,
    ITeamCatalog teamCatalog,
    IUserCatalog userCatalog,
    ITeamMembershipRepository teamMembershipRepository)
  {
    _tenantCatalog = tenantCatalog;
    _teamCatalog = teamCatalog;
    _userCatalog = userCatalog;
    _teamMembershipRepository = teamMembershipRepository;
  }

  public RemoveBootstrapTeamMemberResult Execute(string tenantSlug, Guid teamPublicId, Guid userPublicId)
  {
    if (userPublicId == Guid.Empty)
    {
      return RemoveBootstrapTeamMemberResult.BadRequest(
        new ErrorResponse("invalid_user_public_id", "User public id is required."));
    }

    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return RemoveBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var team = _teamCatalog.FindByTenantIdAndPublicId(tenant.Id, teamPublicId);

    if (team is null)
    {
      return RemoveBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("team_not_found", "Team was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);

    if (user is null)
    {
      return RemoveBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var membership = _teamMembershipRepository.FindByTenantIdAndTeamIdAndUserId(tenant.Id, team.Id, user.Id);

    if (membership is null)
    {
      return RemoveBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("team_membership_not_found", "Team membership was not found."));
    }

    _teamMembershipRepository.RemoveByTenantIdAndTeamIdAndUserId(tenant.Id, team.Id, user.Id);

    return RemoveBootstrapTeamMemberResult.Success(new TeamMembershipResponse(
      membership.Id,
      membership.TenantId,
      team.PublicId,
      user.PublicId,
      user.Email,
      user.DisplayName,
      membership.CreatedAt));
  }
}
