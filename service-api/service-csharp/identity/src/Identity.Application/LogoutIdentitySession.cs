using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class LogoutIdentitySession
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public LogoutIdentitySession(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _auditWriter = auditWriter;
  }

  public OperationResult<UserSessionResponse> Execute(string sessionToken)
  {
    var session = _securityStore.FindSessionBySessionToken(sessionToken.Trim());
    if (session is null || !session.IsActive(DateTimeOffset.UtcNow))
    {
      return OperationResult<UserSessionResponse>.Unauthorized(
        new ErrorResponse("invalid_session", "Session is invalid."));
    }

    var tenant = _tenantCatalog.List().FirstOrDefault(candidate => candidate.Id == session.TenantId);
    if (tenant is null)
    {
      return OperationResult<UserSessionResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, session.UserId);
    if (user is null)
    {
      return OperationResult<UserSessionResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var revokedSession = _securityStore.UpdateSession(session.Revoke(DateTimeOffset.UtcNow));
    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "logout_succeeded", "info", $"Session logged out for {user.Email}.");

    return OperationResult<UserSessionResponse>.Success(revokedSession.ToResponse(user.PublicId));
  }
}
