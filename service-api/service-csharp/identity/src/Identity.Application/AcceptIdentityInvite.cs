using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class AcceptIdentityInvite
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IUserRepository _userRepository;
  private readonly IRoleCatalog _roleCatalog;
  private readonly IUserRoleRepository _userRoleRepository;
  private readonly ITeamCatalog _teamCatalog;
  private readonly ITeamMembershipRepository _teamMembershipRepository;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly IExternalIdentityProvider _identityProvider;
  private readonly SecurityAuditWriter _auditWriter;
  private readonly TenantAccessCoordinator _tenantAccessCoordinator;

  public AcceptIdentityInvite(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IUserRepository userRepository,
    IRoleCatalog roleCatalog,
    IUserRoleRepository userRoleRepository,
    ITeamCatalog teamCatalog,
    ITeamMembershipRepository teamMembershipRepository,
    IIdentitySecurityStore securityStore,
    IExternalIdentityProvider identityProvider,
    SecurityAuditWriter auditWriter,
    TenantAccessCoordinator tenantAccessCoordinator)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _userRepository = userRepository;
    _roleCatalog = roleCatalog;
    _userRoleRepository = userRoleRepository;
    _teamCatalog = teamCatalog;
    _teamMembershipRepository = teamMembershipRepository;
    _securityStore = securityStore;
    _identityProvider = identityProvider;
    _auditWriter = auditWriter;
    _tenantAccessCoordinator = tenantAccessCoordinator;
  }

  public OperationResult<AcceptInviteResponse> Execute(string inviteToken, AcceptInviteRequest request)
  {
    var invite = _securityStore.FindInviteByToken(inviteToken);
    if (invite is null)
    {
      return OperationResult<AcceptInviteResponse>.NotFound(
        new ErrorResponse("invite_not_found", "Invite was not found."));
    }

    if (invite.Status != "pending")
    {
      return OperationResult<AcceptInviteResponse>.Conflict(
        new ErrorResponse("invite_not_pending", "Invite can no longer be accepted."));
    }

    if (invite.IsExpired(DateTimeOffset.UtcNow))
    {
      _securityStore.UpdateInvite(invite.Expire());
      return OperationResult<AcceptInviteResponse>.BadRequest(
        new ErrorResponse("invite_expired", "Invite has expired."));
    }

    var tenant = _tenantCatalog.FindBySlug(invite.TenantSlug);
    if (tenant is null)
    {
      return OperationResult<AcceptInviteResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, invite.UserId);
    if (user is null)
    {
      return OperationResult<AcceptInviteResponse>.NotFound(
        new ErrorResponse("user_not_found", "Invited user was not found."));
    }

    var displayName = string.IsNullOrWhiteSpace(request.DisplayName)
      ? invite.DisplayName ?? user.DisplayName
      : request.DisplayName.Trim();

    if (string.IsNullOrWhiteSpace(displayName))
    {
      return OperationResult<AcceptInviteResponse>.BadRequest(
        new ErrorResponse("invalid_display_name", "Display name is required."));
    }

    if (!PasswordStrength.IsStrong(request.Password))
    {
      return OperationResult<AcceptInviteResponse>.BadRequest(
        new ErrorResponse("invalid_password", "Password must be at least 10 characters and include upper, lower and number."));
    }

    var updatedUser = _userRepository.Update(
      user
        .ReviseProfile(user.Email, displayName, NormalizeOptional(request.GivenName), NormalizeOptional(request.FamilyName))
        .ReviseStatus("active"));

    var profile = _securityStore.GetOrCreateProfile(updatedUser.Id);
    var externalUser = _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
      updatedUser.Email,
      updatedUser.GivenName,
      updatedUser.FamilyName,
      updatedUser.DisplayName,
      true,
      profile.IdentityProviderSubject,
      request.Password));
    _securityStore.SaveProfile(profile.AttachIdentityProviderSubject(externalUser.SubjectId));

    var roles = _roleCatalog.ListByTenantId(tenant.Id).ToDictionary(role => role.Code, StringComparer.OrdinalIgnoreCase);
    foreach (var roleCode in invite.RoleCodes)
    {
      if (!roles.TryGetValue(roleCode, out var role))
      {
        continue;
      }

      if (_userRoleRepository.FindByTenantIdAndUserIdAndRoleId(tenant.Id, updatedUser.Id, role.Id) is not null)
      {
        continue;
      }

      _userRoleRepository.Add(new UserRole(
        _userRoleRepository.NextId(),
        tenant.Id,
        updatedUser.Id,
        role.Id,
        DateTimeOffset.UtcNow));
    }

    var teams = _teamCatalog.ListByTenantId(tenant.Id).ToDictionary(team => team.PublicId);
    foreach (var teamPublicId in invite.TeamPublicIds)
    {
      if (!teams.TryGetValue(teamPublicId, out var team))
      {
        continue;
      }

      if (_teamMembershipRepository.FindByTenantIdAndTeamIdAndUserId(tenant.Id, team.Id, updatedUser.Id) is not null)
      {
        continue;
      }

      _teamMembershipRepository.Add(new TeamMembership(
        _teamMembershipRepository.NextId(),
        tenant.Id,
        team.Id,
        updatedUser.Id,
        DateTimeOffset.UtcNow));
    }

    var acceptedInvite = _securityStore.UpdateInvite(invite.Accept(DateTimeOffset.UtcNow));
    _tenantAccessCoordinator.SyncAndListRoleCodes(tenant, updatedUser);
    _auditWriter.Record(
      tenant.Id,
      updatedUser.PublicId,
      updatedUser.PublicId,
      "invite_accepted",
      "info",
      $"Invite accepted for {updatedUser.Email}.");

    return OperationResult<AcceptInviteResponse>.Success(new AcceptInviteResponse(
      tenant.Slug,
      acceptedInvite.PublicId,
      acceptedInvite.Status,
      updatedUser.ToResponse()));
  }

  private static string? NormalizeOptional(string? value)
  {
    return string.IsNullOrWhiteSpace(value)
      ? null
      : value.Trim();
  }
}
