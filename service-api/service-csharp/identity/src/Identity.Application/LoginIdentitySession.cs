using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class LoginIdentitySession
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly IExternalIdentityProvider _identityProvider;
  private readonly IIdentityAccessSettings _settings;
  private readonly ITotpService _totpService;
  private readonly TenantAccessCoordinator _tenantAccessCoordinator;
  private readonly SecurityAuditWriter _auditWriter;

  public LoginIdentitySession(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    IExternalIdentityProvider identityProvider,
    IIdentityAccessSettings settings,
    ITotpService totpService,
    TenantAccessCoordinator tenantAccessCoordinator,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _identityProvider = identityProvider;
    _settings = settings;
    _totpService = totpService;
    _tenantAccessCoordinator = tenantAccessCoordinator;
    _auditWriter = auditWriter;
  }

  public OperationResult<SessionResponse> Execute(LoginSessionRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(request.TenantSlug.Trim().ToLowerInvariant());
    if (tenant is null)
    {
      return OperationResult<SessionResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var email = request.Email.Trim().ToLowerInvariant();
    var user = _userCatalog.FindByTenantIdAndEmail(tenant.Id, email);
    if (user is null)
    {
      _auditWriter.Record(tenant.Id, null, null, "login_failed", "warning", $"Login failed for unknown email {email}.");
      return OperationResult<SessionResponse>.Unauthorized(
        new ErrorResponse("invalid_credentials", "Credentials are invalid."));
    }

    if (user.Status != "active")
    {
      _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_blocked", "warning", $"Login blocked for {email} with status {user.Status}.");
      return OperationResult<SessionResponse>.Forbidden(
        new ErrorResponse("access_blocked", "User access is not active."));
    }

    var profile = _securityStore.GetOrCreateProfile(user.Id);
    if (string.IsNullOrWhiteSpace(profile.IdentityProviderSubject))
    {
      if (!request.Password.Equals(_settings.BootstrapPassword, StringComparison.Ordinal))
      {
        _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_failed", "warning", $"Bootstrap login failed for {email}.");
        return OperationResult<SessionResponse>.Unauthorized(
          new ErrorResponse("invalid_credentials", "Credentials are invalid."));
      }

      var externalUser = _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
        user.Email,
        user.GivenName,
        user.FamilyName,
        user.DisplayName,
        true,
        null,
        request.Password));
      profile = _securityStore.SaveProfile(profile.AttachIdentityProviderSubject(externalUser.SubjectId));
    }
    else
    {
      _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
        user.Email,
        user.GivenName,
        user.FamilyName,
        user.DisplayName,
        true,
        profile.IdentityProviderSubject,
        null));
    }

    IdentityProviderTokenResult tokenResult;
    try
    {
      tokenResult = _identityProvider.PasswordGrant(user.Email, request.Password);
    }
    catch (ExternalIdentityAuthenticationException)
    {
      if (request.Password.Equals(_settings.BootstrapPassword, StringComparison.Ordinal))
      {
        var repairedExternalUser = _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
          user.Email,
          user.GivenName,
          user.FamilyName,
          user.DisplayName,
          true,
          profile.IdentityProviderSubject,
          request.Password));
        profile = _securityStore.SaveProfile(profile.AttachIdentityProviderSubject(repairedExternalUser.SubjectId));

        var granted = false;
        tokenResult = default!;
        for (var attempt = 0; attempt < 3; attempt++)
        {
          try
          {
            tokenResult = _identityProvider.PasswordGrant(user.Email, request.Password);
            granted = true;
            break;
          }
          catch (ExternalIdentityAuthenticationException)
          {
            Thread.Sleep(250);
          }
        }

        if (!granted)
        {
          _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_failed", "warning", $"Bootstrap login failed for {email} after identity sync.");
          return OperationResult<SessionResponse>.Unauthorized(
            new ErrorResponse("invalid_credentials", "Credentials are invalid."));
        }
      }
      else
      {
        _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_failed", "warning", $"Login failed for {email}.");
        return OperationResult<SessionResponse>.Unauthorized(
          new ErrorResponse("invalid_credentials", "Credentials are invalid."));
      }
    }

    profile = _securityStore.SaveProfile(
      profile
        .AttachIdentityProviderSubject(tokenResult.SubjectId)
        .RecordLogin(DateTimeOffset.UtcNow));

    if (profile.MfaEnabled)
    {
      if (string.IsNullOrWhiteSpace(request.OtpCode))
      {
        _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_mfa_required", "warning", $"Login requires MFA for {email}.");
        return OperationResult<SessionResponse>.Unauthorized(
          new ErrorResponse("mfa_required", "MFA code is required."));
      }

      if (string.IsNullOrWhiteSpace(profile.MfaSecret) || !_totpService.VerifyCode(profile.MfaSecret, request.OtpCode, DateTimeOffset.UtcNow))
      {
        _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_mfa_failed", "warning", $"Invalid MFA code for {email}.");
        return OperationResult<SessionResponse>.Unauthorized(
          new ErrorResponse("invalid_otp_code", "OTP code is invalid."));
      }
    }

    var roleCodes = _tenantAccessCoordinator.SyncAndListRoleCodes(tenant, user);
    if (!_tenantAccessCoordinator.CanAccessTenant(tenant, user))
    {
      _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_forbidden", "warning", $"Access graph denied tenant access for {email}.");
      return OperationResult<SessionResponse>.Forbidden(
        new ErrorResponse("tenant_scope_forbidden", "User does not have tenant access."));
    }

    var session = _securityStore.AddSession(new Session(
      _securityStore.NextSessionId(),
      tenant.Id,
      user.Id,
      PublicIds.NewUuidV7(),
      PublicIds.NewUuidV7().ToString(),
      PublicIds.NewUuidV7().ToString(),
      tokenResult.SubjectId,
      tokenResult.RefreshToken,
      "active",
      tokenResult.AccessExpiresAt,
      tokenResult.RefreshExpiresAt,
      DateTimeOffset.UtcNow,
      DateTimeOffset.UtcNow,
      null));

    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "login_succeeded", "info", $"Login succeeded for {email}.");

    return OperationResult<SessionResponse>.Success(new SessionResponse(
      session.PublicId,
      tenant.Slug,
      user.PublicId,
      user.Email,
      user.DisplayName,
      session.SessionToken,
      session.RefreshToken,
      session.ExpiresAt,
      session.RefreshExpiresAt,
      profile.MfaEnabled,
      roleCodes));
  }
}
