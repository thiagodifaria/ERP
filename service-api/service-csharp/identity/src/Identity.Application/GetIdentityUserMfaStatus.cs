using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class GetIdentityUserMfaStatus
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;

  public GetIdentityUserMfaStatus(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
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

    var profile = _securityStore.GetOrCreateProfile(user.Id);
    return OperationResult<MfaStatusResponse>.Success(new MfaStatusResponse(user.PublicId, profile.MfaEnabled));
  }
}
