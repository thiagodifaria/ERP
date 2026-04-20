using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CreateIdentityInvite
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ICompanyCatalog _companyCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IUserRepository _userRepository;
  private readonly IRoleCatalog _roleCatalog;
  private readonly ITeamCatalog _teamCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly IExternalIdentityProvider _identityProvider;
  private readonly SecurityAuditWriter _auditWriter;

  public CreateIdentityInvite(
    ITenantCatalog tenantCatalog,
    ICompanyCatalog companyCatalog,
    IUserCatalog userCatalog,
    IUserRepository userRepository,
    IRoleCatalog roleCatalog,
    ITeamCatalog teamCatalog,
    IIdentitySecurityStore securityStore,
    IExternalIdentityProvider identityProvider,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _companyCatalog = companyCatalog;
    _userCatalog = userCatalog;
    _userRepository = userRepository;
    _roleCatalog = roleCatalog;
    _teamCatalog = teamCatalog;
    _securityStore = securityStore;
    _identityProvider = identityProvider;
    _auditWriter = auditWriter;
  }

  public OperationResult<InviteResponse> Execute(string tenantSlug, CreateInviteRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<InviteResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var email = NormalizeEmail(request.Email);
    if (!IsValidEmail(email))
    {
      return OperationResult<InviteResponse>.BadRequest(
        new ErrorResponse("invalid_email", "Email is invalid."));
    }

    var pendingInvite = _securityStore.FindPendingInviteByTenantIdAndEmail(tenant.Id, email);
    if (pendingInvite is not null)
    {
      return OperationResult<InviteResponse>.Conflict(
        new ErrorResponse("invite_conflict", "A pending invite already exists for this email."));
    }

    var roleCodes = NormalizeRoleCodes(request.RoleCodes);
    var availableRoles = _roleCatalog.ListByTenantId(tenant.Id)
      .Select(role => role.Code)
      .ToHashSet(StringComparer.OrdinalIgnoreCase);

    if (roleCodes.Any(roleCode => !availableRoles.Contains(roleCode)))
    {
      return OperationResult<InviteResponse>.BadRequest(
        new ErrorResponse("invalid_role_code", "Invite contains an unknown role code."));
    }

    var teamPublicIds = NormalizeTeamPublicIds(request.TeamPublicIds);
    var availableTeamIds = _teamCatalog.ListByTenantId(tenant.Id)
      .Select(team => team.PublicId)
      .ToHashSet();

    if (teamPublicIds.Any(teamPublicId => !availableTeamIds.Contains(teamPublicId)))
    {
      return OperationResult<InviteResponse>.BadRequest(
        new ErrorResponse("invalid_team_public_id", "Invite contains an unknown team public id."));
    }

    var user = _userCatalog.FindByTenantIdAndEmail(tenant.Id, email);
    if (user is not null && user.Status == "active")
    {
      return OperationResult<InviteResponse>.Conflict(
        new ErrorResponse("user_already_active", "User already exists as an active member of the tenant."));
    }

    if (user is not null && user.Status != "invited")
    {
      return OperationResult<InviteResponse>.Conflict(
        new ErrorResponse("user_status_conflict", "User cannot receive an invite in the current status."));
    }

    if (user is null)
    {
      var displayName = NormalizeOptional(request.DisplayName) ?? DefaultDisplayName(email);
      var companyId = _companyCatalog.ListByTenantId(tenant.Id)
        .Select(company => (long?)company.Id)
        .FirstOrDefault();

      user = _userRepository.Add(new User(
        _userRepository.NextId(),
        tenant.Id,
        companyId,
        PublicIds.NewUuidV7(),
        email,
        displayName,
        null,
        null,
        "invited"));
    }

    var profile = _securityStore.GetOrCreateProfile(user.Id);
    var externalUser = _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
      user.Email,
      user.GivenName,
      user.FamilyName,
      user.DisplayName,
      false,
      profile.IdentityProviderSubject,
      null));
    _securityStore.SaveProfile(profile.AttachIdentityProviderSubject(externalUser.SubjectId));

    var invite = new Invite(
      _securityStore.NextInviteId(),
      tenant.Id,
      tenant.Slug,
      user.Id,
      PublicIds.NewUuidV7(),
      PublicIds.NewUuidV7().ToString(),
      email,
      NormalizeOptional(request.DisplayName),
      roleCodes,
      teamPublicIds,
      "pending",
      DateTimeOffset.UtcNow.AddDays(request.ExpiresInDays is > 0 and <= 30 ? request.ExpiresInDays.Value : 7),
      null,
      DateTimeOffset.UtcNow);
    var createdInvite = _securityStore.AddInvite(invite);

    _auditWriter.Record(
      tenant.Id,
      null,
      user.PublicId,
      "invite_created",
      "info",
      $"Invite created for {email}.");

    return OperationResult<InviteResponse>.Success(createdInvite.ToResponse(user.PublicId));
  }

  private static IReadOnlyCollection<string> NormalizeRoleCodes(IReadOnlyCollection<string>? roleCodes)
  {
    return roleCodes?
      .Where(roleCode => !string.IsNullOrWhiteSpace(roleCode))
      .Select(roleCode => roleCode.Trim().ToLowerInvariant())
      .Distinct(StringComparer.OrdinalIgnoreCase)
      .ToArray()
      ?? [];
  }

  private static IReadOnlyCollection<Guid> NormalizeTeamPublicIds(IReadOnlyCollection<Guid>? teamPublicIds)
  {
    return teamPublicIds?
      .Where(teamPublicId => teamPublicId != Guid.Empty)
      .Distinct()
      .ToArray()
      ?? [];
  }

  private static string NormalizeEmail(string email)
  {
    return email.Trim().ToLowerInvariant();
  }

  private static string? NormalizeOptional(string? value)
  {
    return string.IsNullOrWhiteSpace(value)
      ? null
      : value.Trim();
  }

  private static string DefaultDisplayName(string email)
  {
    var localPart = email.Split('@').FirstOrDefault() ?? email;
    return localPart.Replace('.', ' ');
  }

  private static bool IsValidEmail(string email)
  {
    if (string.IsNullOrWhiteSpace(email) || email.Contains(' '))
    {
      return false;
    }

    var separatorIndex = email.IndexOf('@');
    return separatorIndex > 0
      && separatorIndex < email.Length - 1
      && separatorIndex == email.LastIndexOf('@');
  }
}
