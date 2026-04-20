// Este store em memoria sustenta sessoes, convites, MFA e auditoria em testes locais.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryIdentitySecurityStore : IIdentitySecurityStore
{
  private readonly Dictionary<long, UserSecurityProfile> _profilesByUserId = [];
  private readonly Dictionary<string, Invite> _invitesByToken = new(StringComparer.OrdinalIgnoreCase);
  private readonly Dictionary<long, Invite> _invitesById = [];
  private readonly Dictionary<string, Session> _sessionsByToken = new(StringComparer.Ordinal);
  private readonly Dictionary<string, Session> _sessionsByRefreshToken = new(StringComparer.Ordinal);
  private readonly Dictionary<long, Session> _sessionsById = [];
  private readonly List<SecurityAuditEvent> _auditEvents = [];
  private long _inviteSequence = 1;
  private long _sessionSequence = 1;
  private long _auditSequence = 1;

  public UserSecurityProfile GetOrCreateProfile(long userId)
  {
    if (_profilesByUserId.TryGetValue(userId, out var profile))
    {
      return profile;
    }

    profile = new UserSecurityProfile(userId, null, false, null, null);
    _profilesByUserId[userId] = profile;
    return profile;
  }

  public UserSecurityProfile SaveProfile(UserSecurityProfile profile)
  {
    _profilesByUserId[profile.UserId] = profile;
    return profile;
  }

  public Invite? FindInviteByToken(string inviteToken)
  {
    return _invitesByToken.TryGetValue(inviteToken, out var invite)
      ? invite
      : null;
  }

  public Invite? FindPendingInviteByTenantIdAndEmail(long tenantId, string email)
  {
    return _invitesById.Values
      .Where(invite => invite.TenantId == tenantId && invite.Status == "pending")
      .FirstOrDefault(invite => invite.Email.Equals(email, StringComparison.OrdinalIgnoreCase));
  }

  public IReadOnlyCollection<Invite> ListInvitesByTenantId(long tenantId)
  {
    return _invitesById.Values
      .Where(invite => invite.TenantId == tenantId)
      .OrderByDescending(invite => invite.CreatedAt)
      .ToArray();
  }

  public Invite AddInvite(Invite invite)
  {
    _invitesById[invite.Id] = invite;
    _invitesByToken[invite.InviteToken] = invite;
    return invite;
  }

  public Invite UpdateInvite(Invite invite)
  {
    _invitesById[invite.Id] = invite;
    _invitesByToken[invite.InviteToken] = invite;
    return invite;
  }

  public Session? FindSessionBySessionToken(string sessionToken)
  {
    return _sessionsByToken.TryGetValue(sessionToken, out var session)
      ? session
      : null;
  }

  public Session? FindSessionByRefreshToken(string refreshToken)
  {
    return _sessionsByRefreshToken.TryGetValue(refreshToken, out var session)
      ? session
      : null;
  }

  public Session AddSession(Session session)
  {
    _sessionsById[session.Id] = session;
    _sessionsByToken[session.SessionToken] = session;
    _sessionsByRefreshToken[session.RefreshToken] = session;
    return session;
  }

  public Session UpdateSession(Session session)
  {
    if (_sessionsById.TryGetValue(session.Id, out var current))
    {
      _sessionsByToken.Remove(current.SessionToken);
      _sessionsByRefreshToken.Remove(current.RefreshToken);
    }

    _sessionsById[session.Id] = session;
    _sessionsByToken[session.SessionToken] = session;
    _sessionsByRefreshToken[session.RefreshToken] = session;
    return session;
  }

  public int RevokeSessionsByUserId(long tenantId, long userId, DateTimeOffset revokedAt)
  {
    var sessions = _sessionsById.Values
      .Where(session => session.TenantId == tenantId && session.UserId == userId && session.Status == "active")
      .ToArray();

    foreach (var session in sessions)
    {
      UpdateSession(session.Revoke(revokedAt));
    }

    return sessions.Length;
  }

  public IReadOnlyCollection<SecurityAuditEvent> ListSecurityAuditByTenantId(long tenantId, int limit)
  {
    return _auditEvents
      .Where(auditEvent => auditEvent.TenantId == tenantId)
      .OrderByDescending(auditEvent => auditEvent.CreatedAt)
      .Take(limit)
      .ToArray();
  }

  public SecurityAuditEvent AddSecurityAuditEvent(SecurityAuditEvent auditEvent)
  {
    _auditEvents.Add(auditEvent);
    return auditEvent;
  }

  public long NextInviteId()
  {
    return _inviteSequence++;
  }

  public long NextSessionId()
  {
    return _sessionSequence++;
  }

  public long NextSecurityAuditEventId()
  {
    return _auditSequence++;
  }
}
