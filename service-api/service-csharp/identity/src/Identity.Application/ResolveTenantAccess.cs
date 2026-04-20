using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ResolveTenantAccess
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly TenantAccessCoordinator _tenantAccessCoordinator;

  public ResolveTenantAccess(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore,
    TenantAccessCoordinator tenantAccessCoordinator)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
    _tenantAccessCoordinator = tenantAccessCoordinator;
  }

  public OperationResult<AccessResolutionResponse> Execute(string tenantSlug, string sessionToken)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<AccessResolutionResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var session = _securityStore.FindSessionBySessionToken(sessionToken);
    if (session is null || !session.IsActive(DateTimeOffset.UtcNow))
    {
      return OperationResult<AccessResolutionResponse>.Unauthorized(
        new ErrorResponse("invalid_session", "Session is invalid."));
    }

    if (session.TenantId != tenant.Id)
    {
      return OperationResult<AccessResolutionResponse>.Forbidden(
        new ErrorResponse("tenant_scope_forbidden", "Session does not belong to the requested tenant."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, session.UserId);
    if (user is null)
    {
      return OperationResult<AccessResolutionResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    if (user.Status != "active")
    {
      return OperationResult<AccessResolutionResponse>.Forbidden(
        new ErrorResponse("access_blocked", "User access is not active."));
    }

    var roleCodes = _tenantAccessCoordinator.SyncAndListRoleCodes(tenant, user);
    if (!_tenantAccessCoordinator.CanAccessTenant(tenant, user))
    {
      return OperationResult<AccessResolutionResponse>.Forbidden(
        new ErrorResponse("tenant_scope_forbidden", "User does not have tenant access."));
    }

    var profile = _securityStore.GetOrCreateProfile(user.Id);
    _securityStore.UpdateSession(session.Touch(DateTimeOffset.UtcNow));

    return OperationResult<AccessResolutionResponse>.Success(new AccessResolutionResponse(
      tenant.Slug,
      session.PublicId,
      user.PublicId,
      user.Email,
      user.DisplayName,
      roleCodes,
      profile.MfaEnabled,
      true,
      user.Status));
  }
}
