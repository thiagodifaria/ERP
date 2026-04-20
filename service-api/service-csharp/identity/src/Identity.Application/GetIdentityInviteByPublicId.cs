using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class GetIdentityInviteByPublicId
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;

  public GetIdentityInviteByPublicId(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IIdentitySecurityStore securityStore)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _securityStore = securityStore;
  }

  public OperationResult<InviteResponse> Execute(string tenantSlug, Guid invitePublicId)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<InviteResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var invite = _securityStore.FindInviteByPublicId(invitePublicId);
    if (invite is null || invite.TenantId != tenant.Id)
    {
      return OperationResult<InviteResponse>.NotFound(
        new ErrorResponse("invite_not_found", "Invite was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, invite.UserId);
    return OperationResult<InviteResponse>.Success(invite.ToResponse(user?.PublicId ?? Guid.Empty));
  }
}
