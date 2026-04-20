using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class StartIdentityUserMfaEnrollment
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly ITotpService _totpService;
  private readonly IIdentityAccessSettings _settings;
  private readonly SecurityAuditWriter _auditWriter;

  public StartIdentityUserMfaEnrollment(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    ITotpService totpService,
    IIdentityAccessSettings settings,
    SecurityAuditWriter auditWriter)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _totpService = totpService;
    _settings = settings;
    _auditWriter = auditWriter;
  }

  public OperationResult<MfaEnrollmentResponse> Execute(string tenantSlug, Guid userPublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<MfaEnrollmentResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);
    if (user is null)
    {
      return OperationResult<MfaEnrollmentResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var secret = _totpService.GenerateSecret();
    _securityStore.SaveProfile(_securityStore.GetOrCreateProfile(user.Id).StartMfaEnrollment(secret));

    _auditWriter.Record(
      tenant.Id,
      user.PublicId,
      user.PublicId,
      "mfa_enrollment_started",
      "info",
      $"MFA enrollment started for {user.Email}.");

    return OperationResult<MfaEnrollmentResponse>.Success(new MfaEnrollmentResponse(
      user.PublicId,
      false,
      secret,
      _totpService.BuildOtpAuthUri(_settings.MfaIssuer, user.Email, secret)));
  }
}
