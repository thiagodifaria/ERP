using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class RefreshIdentitySession
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly IExternalIdentityProvider _identityProvider;
  private readonly TenantAccessCoordinator _tenantAccessCoordinator;
  private readonly SecurityAuditWriter _auditWriter;

  public RefreshIdentitySession(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    IExternalIdentityProvider identityProvider,
    TenantAccessCoordinator tenantAccessCoordinator,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _identityProvider = identityProvider;
    _tenantAccessCoordinator = tenantAccessCoordinator;
    _auditWriter = auditWriter;
  }

  public OperationResult<SessionResponse> Execute(RefreshSessionRequest request)
  {
    var session = _securityStore.FindSessionByRefreshToken(request.RefreshToken.Trim());
    if (session is null || !session.CanRefresh(DateTimeOffset.UtcNow))
    {
      return OperationResult<SessionResponse>.Unauthorized(
        new ErrorResponse("invalid_refresh_token", "Refresh token is invalid."));
    }

    var tenant = _tenantCatalog.List().FirstOrDefault(candidate => candidate.Id == session.TenantId);
    if (tenant is null)
    {
      return OperationResult<SessionResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, session.UserId);
    if (user is null)
    {
      return OperationResult<SessionResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    if (user.Status != "active")
    {
      _securityStore.UpdateSession(session.Revoke(DateTimeOffset.UtcNow));
      return OperationResult<SessionResponse>.Forbidden(
        new ErrorResponse("access_blocked", "User access is not active."));
    }

    IdentityProviderTokenResult tokenResult;
    try
    {
      tokenResult = _identityProvider.RefreshGrant(session.IdentityProviderRefreshToken ?? string.Empty);
    }
    catch (ExternalIdentityAuthenticationException)
    {
      _securityStore.UpdateSession(session.Revoke(DateTimeOffset.UtcNow));
      return OperationResult<SessionResponse>.Unauthorized(
        new ErrorResponse("invalid_refresh_token", "Refresh token is invalid."));
    }

    var roleCodes = _tenantAccessCoordinator.SyncAndListRoleCodes(tenant, user);
    if (!_tenantAccessCoordinator.CanAccessTenant(tenant, user))
    {
      _securityStore.UpdateSession(session.Revoke(DateTimeOffset.UtcNow));
      return OperationResult<SessionResponse>.Forbidden(
        new ErrorResponse("tenant_scope_forbidden", "User does not have tenant access."));
    }

    var refreshedSession = _securityStore.UpdateSession(session.Refresh(
      PublicIds.NewUuidV7().ToString(),
      tokenResult.RefreshToken,
      tokenResult.AccessExpiresAt,
      tokenResult.RefreshExpiresAt,
      DateTimeOffset.UtcNow));
    var profile = _securityStore.GetOrCreateProfile(user.Id);

    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "session_refreshed", "info", $"Session refreshed for {user.Email}.");

    return OperationResult<SessionResponse>.Success(new SessionResponse(
      refreshedSession.PublicId,
      tenant.Slug,
      user.PublicId,
      user.Email,
      user.DisplayName,
      refreshedSession.SessionToken,
      refreshedSession.RefreshToken,
      refreshedSession.ExpiresAt,
      refreshedSession.RefreshExpiresAt,
      profile.MfaEnabled,
      roleCodes));
  }
}
