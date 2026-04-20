using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ResendIdentityInvite
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public ResendIdentityInvite(
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

  public OperationResult<InviteResponse> Execute(string tenantSlug, Guid invitePublicId, ResendInviteRequest request)
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

    if (invite.Status.Equals("accepted", StringComparison.OrdinalIgnoreCase))
    {
      return OperationResult<InviteResponse>.Conflict(
        new ErrorResponse("invite_already_accepted", "Invite was already accepted."));
    }

    var resentInvite = _securityStore.UpdateInvite(invite.Reissue(
      PublicIds.NewUuidV7().ToString(),
      DateTimeOffset.UtcNow.AddDays(request.ExpiresInDays is > 0 and <= 30 ? request.ExpiresInDays.Value : 7)));
    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, invite.UserId);
    var userPublicId = user?.PublicId ?? Guid.Empty;

    _auditWriter.Record(tenant.Id, null, userPublicId == Guid.Empty ? null : userPublicId, "invite_resent", "info", $"Invite reissued for {invite.Email}.");

    return OperationResult<InviteResponse>.Success(resentInvite.ToResponse(userPublicId));
  }
}
