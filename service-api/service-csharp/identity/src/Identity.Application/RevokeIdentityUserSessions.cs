using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class RevokeIdentityUserSessions
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public RevokeIdentityUserSessions(
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

  public OperationResult<IReadOnlyCollection<UserSessionResponse>> Execute(string tenantSlug, Guid userPublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<IReadOnlyCollection<UserSessionResponse>>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);
    if (user is null)
    {
      return OperationResult<IReadOnlyCollection<UserSessionResponse>>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    _securityStore.RevokeSessionsByUserId(tenant.Id, user.Id, DateTimeOffset.UtcNow);
    var sessions = _securityStore.ListSessionsByTenantIdAndUserId(tenant.Id, user.Id)
      .Select(session => session.ToResponse(user.PublicId))
      .ToArray();

    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "sessions_revoked", "warning", $"All active sessions revoked for {user.Email}.");

    return OperationResult<IReadOnlyCollection<UserSessionResponse>>.Success(sessions);
  }
}
