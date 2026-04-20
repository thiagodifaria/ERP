using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListIdentityUserSessions
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;

  public ListIdentityUserSessions(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
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

    var sessions = _securityStore.ListSessionsByTenantIdAndUserId(tenant.Id, user.Id)
      .Select(session => session.ToResponse(user.PublicId))
      .ToArray();

    return OperationResult<IReadOnlyCollection<UserSessionResponse>>.Success(sessions);
  }
}
