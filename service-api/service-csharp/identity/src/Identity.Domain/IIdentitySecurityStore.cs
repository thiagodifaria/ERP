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

  Session AddSession(Session session);

  Session UpdateSession(Session session);

  int RevokeSessionsByUserId(long tenantId, long userId, DateTimeOffset revokedAt);

  IReadOnlyCollection<SecurityAuditEvent> ListSecurityAuditByTenantId(long tenantId, int limit);

  SecurityAuditEvent AddSecurityAuditEvent(SecurityAuditEvent auditEvent);

  long NextInviteId();

  long NextSessionId();

  long NextSecurityAuditEventId();
}
