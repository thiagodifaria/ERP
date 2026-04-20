using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class VerifyIdentityUserMfa
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly ITotpService _totpService;
  private readonly SecurityAuditWriter _auditWriter;

  public VerifyIdentityUserMfa(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    ITotpService totpService,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _totpService = totpService;
    _auditWriter = auditWriter;
  }

  public OperationResult<MfaStatusResponse> Execute(string tenantSlug, Guid userPublicId, VerifyMfaRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<MfaStatusResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);
    if (user is null)
    {
      return OperationResult<MfaStatusResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var profile = _securityStore.GetOrCreateProfile(user.Id);
    if (string.IsNullOrWhiteSpace(profile.MfaSecret))
    {
      return OperationResult<MfaStatusResponse>.BadRequest(
        new ErrorResponse("mfa_not_initialized", "MFA enrollment has not been started."));
    }

    if (!_totpService.VerifyCode(profile.MfaSecret, request.OtpCode, DateTimeOffset.UtcNow))
    {
      return OperationResult<MfaStatusResponse>.BadRequest(
        new ErrorResponse("invalid_otp_code", "OTP code is invalid."));
    }

    _securityStore.SaveProfile(profile.EnableMfa());
    _auditWriter.Record(
      tenant.Id,
      user.PublicId,
      user.PublicId,
      "mfa_enabled",
      "info",
      $"MFA enabled for {user.Email}.");

    return OperationResult<MfaStatusResponse>.Success(new MfaStatusResponse(user.PublicId, true));
  }
}
