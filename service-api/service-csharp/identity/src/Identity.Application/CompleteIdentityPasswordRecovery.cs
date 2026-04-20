using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CompleteIdentityPasswordRecovery
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly IExternalIdentityProvider _identityProvider;
  private readonly SecurityAuditWriter _auditWriter;

  public CompleteIdentityPasswordRecovery(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    IExternalIdentityProvider identityProvider,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _identityProvider = identityProvider;
    _auditWriter = auditWriter;
  }

  public OperationResult<PasswordRecoveryResponse> Execute(string resetToken, ResetPasswordRequest request)
  {
    var passwordResetToken = _securityStore.FindPasswordResetTokenByResetToken(resetToken.Trim());
    if (passwordResetToken is null)
    {
      return OperationResult<PasswordRecoveryResponse>.NotFound(
        new ErrorResponse("password_reset_not_found", "Password reset token was not found."));
    }

    if (passwordResetToken.Status != "pending")
    {
      return OperationResult<PasswordRecoveryResponse>.Conflict(
        new ErrorResponse("password_reset_not_pending", "Password reset token can no longer be used."));
    }

    if (passwordResetToken.IsExpired(DateTimeOffset.UtcNow))
    {
      _securityStore.UpdatePasswordResetToken(passwordResetToken.Expire());
      return OperationResult<PasswordRecoveryResponse>.BadRequest(
        new ErrorResponse("password_reset_expired", "Password reset token has expired."));
    }

    if (!PasswordStrength.IsStrong(request.Password))
    {
      return OperationResult<PasswordRecoveryResponse>.BadRequest(
        new ErrorResponse("invalid_password", "Password must be at least 10 characters and include upper, lower and number."));
    }

    var tenant = _tenantCatalog.List().FirstOrDefault(candidate => candidate.Id == passwordResetToken.TenantId);
    if (tenant is null)
    {
      return OperationResult<PasswordRecoveryResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, passwordResetToken.UserId);
    if (user is null)
    {
      return OperationResult<PasswordRecoveryResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var profile = _securityStore.GetOrCreateProfile(user.Id);
    var externalUser = _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
      user.Email,
      user.GivenName,
      user.FamilyName,
      user.DisplayName,
      true,
      profile.IdentityProviderSubject,
      request.Password));
    _securityStore.SaveProfile(profile.AttachIdentityProviderSubject(externalUser.SubjectId));
    _securityStore.RevokeSessionsByUserId(tenant.Id, user.Id, DateTimeOffset.UtcNow);

    var consumedToken = _securityStore.UpdatePasswordResetToken(passwordResetToken.Consume(DateTimeOffset.UtcNow));
    _auditWriter.Record(tenant.Id, user.PublicId, user.PublicId, "password_reset_completed", "info", $"Password reset completed for {user.Email}.");

    return OperationResult<PasswordRecoveryResponse>.Success(consumedToken.ToResponse(tenant.Slug, user));
  }
}
