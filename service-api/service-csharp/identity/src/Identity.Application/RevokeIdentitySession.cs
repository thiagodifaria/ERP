using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class RevokeIdentitySession
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public RevokeIdentitySession(
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

  public OperationResult<UserSessionResponse> Execute(string tenantSlug, Guid sessionPublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<UserSessionResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var session = _securityStore.FindSessionByPublicId(sessionPublicId);
    if (session is null || session.TenantId != tenant.Id)
    {
      return OperationResult<UserSessionResponse>.NotFound(
        new ErrorResponse("session_not_found", "Session was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, session.UserId);
    if (user is null)
    {
      return OperationResult<UserSessionResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var revokedSession = session.Status == "active"
      ? _securityStore.UpdateSession(session.Revoke(DateTimeOffset.UtcNow))
      : session;

    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "session_revoked", "warning", $"Session revoked for {user.Email}.");

    return OperationResult<UserSessionResponse>.Success(revokedSession.ToResponse(user.PublicId));
  }
}
