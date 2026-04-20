using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CancelIdentityInvite
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly SecurityAuditWriter _auditWriter;

  public CancelIdentityInvite(
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

    if (!invite.Status.Equals("pending", StringComparison.OrdinalIgnoreCase))
    {
      return OperationResult<InviteResponse>.Conflict(
        new ErrorResponse("invite_not_pending", "Invite can no longer be updated."));
    }

    var revokedInvite = _securityStore.UpdateInvite(invite.Revoke());
    var user = _userCatalog.FindByTenantIdAndId(tenant.Id, invite.UserId);
    var userPublicId = user?.PublicId ?? Guid.Empty;

    _auditWriter.Record(tenant.Id, null, userPublicId == Guid.Empty ? null : userPublicId, "invite_revoked", "warning", $"Invite revoked for {invite.Email}.");

    return OperationResult<InviteResponse>.Success(revokedInvite.ToResponse(userPublicId));
  }
}
