// IdentityResponseMapper centraliza a montagem dos contratos publicos do contexto.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

internal static class IdentityResponseMapper
{
  public static UserResponse ToResponse(this User user)
  {
    return new UserResponse(
      user.Id,
      user.PublicId,
      user.TenantId,
      user.CompanyId,
      user.Email,
      user.DisplayName,
      user.GivenName,
      user.FamilyName,
      user.Status);
  }

  public static InviteResponse ToResponse(this Invite invite, Guid userPublicId)
  {
    return new InviteResponse(
      invite.PublicId,
      userPublicId,
      invite.Email,
      invite.DisplayName,
      invite.RoleCodes,
      invite.TeamPublicIds,
      invite.Status,
      invite.InviteToken,
      invite.ExpiresAt);
  }

  public static SecurityAuditEventResponse ToResponse(this SecurityAuditEvent auditEvent)
  {
    return new SecurityAuditEventResponse(
      auditEvent.PublicId,
      auditEvent.ActorUserPublicId,
      auditEvent.SubjectUserPublicId,
      auditEvent.EventCode,
      auditEvent.Severity,
      auditEvent.Summary,
      auditEvent.CreatedAt);
  }

  public static PasswordRecoveryResponse ToResponse(this PasswordResetToken passwordResetToken, string tenantSlug, User user)
  {
    return new PasswordRecoveryResponse(
      passwordResetToken.PublicId,
      tenantSlug,
      user.PublicId,
      user.Email,
      passwordResetToken.Status,
      passwordResetToken.ResetToken,
      passwordResetToken.ExpiresAt);
  }

  public static UserSessionResponse ToResponse(this Session session, Guid userPublicId)
  {
    return new UserSessionResponse(
      session.PublicId,
      userPublicId,
      session.Status,
      session.ExpiresAt,
      session.RefreshExpiresAt,
      session.CreatedAt,
      session.LastUsedAt,
      session.RevokedAt);
  }
}
