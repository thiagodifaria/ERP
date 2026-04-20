using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class DisableIdentityUserMfa
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public DisableIdentityUserMfa(
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

  public OperationResult<MfaStatusResponse> Execute(string tenantSlug, Guid userPublicId)
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

    _securityStore.SaveProfile(_securityStore.GetOrCreateProfile(user.Id).DisableMfa());
    _auditWriter.Record(
      tenant.Id,
      user.PublicId,
      user.PublicId,
      "mfa_disabled",
      "warning",
      $"MFA disabled for {user.Email}.");

    return OperationResult<MfaStatusResponse>.Success(new MfaStatusResponse(user.PublicId, false));
  }
}
