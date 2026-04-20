using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListIdentityInvites
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;

  public ListIdentityInvites(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
  }

  public OperationResult<IReadOnlyCollection<InviteResponse>> Execute(string tenantSlug)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<IReadOnlyCollection<InviteResponse>>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var userLookup = _userCatalog.ListByTenantId(tenant.Id).ToDictionary(user => user.Id, user => user.PublicId);
    var invites = _securityStore.ListInvitesByTenantId(tenant.Id)
      .Select(invite => invite.ToResponse(userLookup.TryGetValue(invite.UserId, out var userPublicId) ? userPublicId : Guid.Empty))
      .ToArray();

    return OperationResult<IReadOnlyCollection<InviteResponse>>.Success(invites);
  }
}
