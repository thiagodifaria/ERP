// Este contrato agrupa os dados de seguranca e sessao do contexto de identidade.
// O objetivo e manter convites, MFA, auditoria e sessoes consistentes em um unico adapter.
namespace Identity.Domain;

public interface IIdentitySecurityStore
{
  UserSecurityProfile GetOrCreateProfile(long userId);

  UserSecurityProfile SaveProfile(UserSecurityProfile profile);

  Invite? FindInviteByToken(string inviteToken);

  Invite? FindPendingInviteByTenantIdAndEmail(long tenantId, string email);

  IReadOnlyCollection<Invite> ListInvitesByTenantId(long tenantId);

  Invite AddInvite(Invite invite);

  Invite UpdateInvite(Invite invite);

  Session? FindSessionBySessionToken(string sessionToken);

  Session? FindSessionByRefreshToken(string refreshToken);

  Session? FindSessionByPublicId(Guid sessionPublicId);

  IReadOnlyCollection<Session> ListSessionsByTenantIdAndUserId(long tenantId, long userId);

  Session AddSession(Session session);

  Session UpdateSession(Session session);

  int RevokeSessionsByUserId(long tenantId, long userId, DateTimeOffset revokedAt);

  PasswordResetToken? FindPasswordResetTokenByResetToken(string resetToken);

  PasswordResetToken? FindPendingPasswordResetTokenByTenantIdAndUserId(long tenantId, long userId);

  PasswordResetToken AddPasswordResetToken(PasswordResetToken passwordResetToken);

  PasswordResetToken UpdatePasswordResetToken(PasswordResetToken passwordResetToken);

  IReadOnlyCollection<SecurityAuditEvent> ListSecurityAuditByTenantId(long tenantId, int limit);

  SecurityAuditEvent AddSecurityAuditEvent(SecurityAuditEvent auditEvent);

  long NextInviteId();

  long NextSessionId();

  long NextPasswordResetTokenId();

  long NextSecurityAuditEventId();
}
