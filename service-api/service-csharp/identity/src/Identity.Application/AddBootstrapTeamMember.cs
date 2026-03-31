// Este caso de uso cria memberships entre times e usuarios dentro do tenant.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class AddBootstrapTeamMember
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ITeamCatalog _teamCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly ITeamMembershipRepository _teamMembershipRepository;

  public AddBootstrapTeamMember(
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

  public AddBootstrapTeamMemberResult Execute(
    string tenantSlug,
    Guid teamPublicId,
    AddTeamMemberRequest request)
  {
    if (request.UserPublicId == Guid.Empty)
    {
      return AddBootstrapTeamMemberResult.BadRequest(
        new ErrorResponse("invalid_user_public_id", "User public id is required."));
    }

    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return AddBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var team = _teamCatalog.FindByTenantIdAndPublicId(tenant.Id, teamPublicId);

    if (team is null)
    {
      return AddBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("team_not_found", "Team was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, request.UserPublicId);

    if (user is null)
    {
      return AddBootstrapTeamMemberResult.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    if (_teamMembershipRepository.FindByTenantIdAndTeamIdAndUserId(tenant.Id, team.Id, user.Id) is not null)
    {
      return AddBootstrapTeamMemberResult.Conflict(
        new ErrorResponse("team_membership_conflict", "User is already a member of this team."));
    }

    var membership = new TeamMembership(
      _teamMembershipRepository.NextId(),
      tenant.Id,
      team.Id,
      user.Id,
      DateTimeOffset.UtcNow);

    var createdMembership = _teamMembershipRepository.Add(membership);

    return AddBootstrapTeamMemberResult.Success(new TeamMembershipResponse(
      createdMembership.Id,
      createdMembership.TenantId,
      team.PublicId,
      user.PublicId,
      user.Email,
      user.DisplayName,
      createdMembership.CreatedAt));
  }
}
