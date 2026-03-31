// Este caso de uso consolida a estrutura basica do tenant em uma unica leitura.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class GetBootstrapTenantAccessSnapshot
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ICompanyCatalog _companyCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly ITeamCatalog _teamCatalog;
  private readonly IRoleCatalog _roleCatalog;
  private readonly ITeamMembershipCatalog _teamMembershipCatalog;
  private readonly IUserRoleCatalog _userRoleCatalog;

  public GetBootstrapTenantAccessSnapshot(
    ITenantCatalog tenantCatalog,
    ICompanyCatalog companyCatalog,
    IUserCatalog userCatalog,
    ITeamCatalog teamCatalog,
    IRoleCatalog roleCatalog,
    ITeamMembershipCatalog teamMembershipCatalog,
    IUserRoleCatalog userRoleCatalog)
  {
    _tenantCatalog = tenantCatalog;
    _companyCatalog = companyCatalog;
    _userCatalog = userCatalog;
    _teamCatalog = teamCatalog;
    _roleCatalog = roleCatalog;
    _teamMembershipCatalog = teamMembershipCatalog;
    _userRoleCatalog = userRoleCatalog;
  }

  public TenantAccessSnapshotResponse? Execute(string tenantSlug)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return null;
    }

    var companies = _companyCatalog
      .ListByTenantId(tenant.Id)
      .Select(company => new CompanyResponse(
        company.Id,
        company.PublicId,
        company.TenantId,
        company.DisplayName,
        company.LegalName,
        company.TaxId,
        company.Status))
      .ToArray();

    var roles = _roleCatalog
      .ListByTenantId(tenant.Id)
      .Select(role => new RoleResponse(
        role.Id,
        role.PublicId,
        role.Code,
        role.DisplayName,
        role.Status))
      .ToArray();

    var users = _userCatalog
      .ListByTenantId(tenant.Id)
      .Select(user =>
      {
        var roleResponses = _userRoleCatalog
          .ListByTenantIdAndUserId(tenant.Id, user.Id)
          .Select(userRole =>
          {
            var role = _roleCatalog.FindByTenantIdAndId(tenant.Id, userRole.RoleId);

            return role is null
              ? null
              : new UserRoleResponse(
                userRole.Id,
                userRole.TenantId,
                user.PublicId,
                role.PublicId,
                role.Code,
                role.DisplayName,
                userRole.CreatedAt);
          })
          .Where(userRole => userRole is not null)
          .Select(userRole => userRole!)
          .ToArray();

        return new UserAccessSnapshotResponse(
          new UserResponse(
            user.Id,
            user.PublicId,
            user.TenantId,
            user.CompanyId,
            user.Email,
            user.DisplayName,
            user.GivenName,
            user.FamilyName,
            user.Status),
          roleResponses);
      })
      .ToArray();

    var teams = _teamCatalog
      .ListByTenantId(tenant.Id)
      .Select(team =>
      {
        var memberResponses = _teamMembershipCatalog
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
          .Where(member => member is not null)
          .Select(member => member!)
          .ToArray();

        return new TeamAccessSnapshotResponse(
          new TeamResponse(
            team.Id,
            team.PublicId,
            team.TenantId,
            team.CompanyId,
            team.Name,
            team.Status),
          memberResponses);
      })
      .ToArray();

    var teamMemberships = teams.Sum(team => team.Members.Count);
    var userRoles = users.Sum(user => user.Roles.Count);

    return new TenantAccessSnapshotResponse(
      new TenantResponse(
        tenant.Id,
        tenant.PublicId,
        tenant.Slug,
        tenant.DisplayName,
        tenant.Status),
      new TenantStructureCountsResponse(
        companies.Length,
        users.Length,
        teams.Length,
        roles.Length,
        teamMemberships,
        userRoles),
      companies,
      users,
      teams,
      roles);
  }
}
