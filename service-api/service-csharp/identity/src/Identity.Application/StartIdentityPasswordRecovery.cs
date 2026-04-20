using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class StartIdentityPasswordRecovery
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public StartIdentityPasswordRecovery(
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

  public OperationResult<PasswordRecoveryResponse> Execute(StartPasswordRecoveryRequest request)
  {
    var tenantSlug = request.TenantSlug.Trim().ToLowerInvariant();
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<PasswordRecoveryResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var email = request.Email.Trim().ToLowerInvariant();
    var user = _userCatalog.FindByTenantIdAndEmail(tenant.Id, email);
    if (user is null)
    {
      _auditWriter.Record(tenant.Id, null, null, "password_reset_rejected", "warning", $"Password reset requested for unknown email {email}.");
      return OperationResult<PasswordRecoveryResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    if (user.Status != "active")
    {
      _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "password_reset_blocked", "warning", $"Password reset blocked for {email} with status {user.Status}.");
      return OperationResult<PasswordRecoveryResponse>.Forbidden(
        new ErrorResponse("access_blocked", "User access is not active."));
    }

    var pendingToken = _securityStore.FindPendingPasswordResetTokenByTenantIdAndUserId(tenant.Id, user.Id);
    if (pendingToken is not null)
    {
      _securityStore.UpdatePasswordResetToken(pendingToken.IsExpired(DateTimeOffset.UtcNow) ? pendingToken.Expire() : pendingToken.Expire());
    }

    var expiresInMinutes = request.ExpiresInMinutes is >= 10 and <= 240
      ? request.ExpiresInMinutes.Value
      : 30;

    var passwordResetToken = _securityStore.AddPasswordResetToken(new PasswordResetToken(
      _securityStore.NextPasswordResetTokenId(),
      tenant.Id,
      user.Id,
      PublicIds.NewUuidV7(),
      PublicIds.NewUuidV7().ToString(),
      "pending",
      DateTimeOffset.UtcNow.AddMinutes(expiresInMinutes),
      null,
      DateTimeOffset.UtcNow));

    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "password_reset_requested", "info", $"Password reset requested for {email}.");

    return OperationResult<PasswordRecoveryResponse>.Success(passwordResetToken.ToResponse(tenant.Slug, user));
  }
}
