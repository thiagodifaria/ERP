// Este store em memoria sustenta sessoes, convites, MFA e auditoria em testes locais.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryIdentitySecurityStore : IIdentitySecurityStore
{
  private readonly Dictionary<long, UserSecurityProfile> _profilesByUserId = [];
  private readonly Dictionary<string, Invite> _invitesByToken = new(StringComparer.OrdinalIgnoreCase);
  private readonly Dictionary<Guid, Invite> _invitesByPublicId = [];
  private readonly Dictionary<long, Invite> _invitesById = [];
  private readonly Dictionary<string, Session> _sessionsByToken = new(StringComparer.Ordinal);
  private readonly Dictionary<string, Session> _sessionsByRefreshToken = new(StringComparer.Ordinal);
  private readonly Dictionary<Guid, Session> _sessionsByPublicId = [];
  private readonly Dictionary<long, Session> _sessionsById = [];
  private readonly Dictionary<string, PasswordResetToken> _passwordResetTokensByToken = new(StringComparer.Ordinal);
  private readonly Dictionary<long, PasswordResetToken> _passwordResetTokensById = [];
  private readonly List<SecurityAuditEvent> _auditEvents = [];
  private long _inviteSequence = 1;
  private long _sessionSequence = 1;
  private long _passwordResetTokenSequence = 1;
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

  public Invite? FindInviteByPublicId(Guid invitePublicId)
  {
    return _invitesByPublicId.TryGetValue(invitePublicId, out var invite)
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
    _invitesByPublicId[invite.PublicId] = invite;
    _invitesByToken[invite.InviteToken] = invite;
    return invite;
  }

  public Invite UpdateInvite(Invite invite)
  {
    if (_invitesById.TryGetValue(invite.Id, out var currentInvite))
    {
      _invitesByToken.Remove(currentInvite.InviteToken);
    }

    _invitesById[invite.Id] = invite;
    _invitesByPublicId[invite.PublicId] = invite;
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

  public Session? FindSessionByPublicId(Guid sessionPublicId)
  {
    return _sessionsByPublicId.TryGetValue(sessionPublicId, out var session)
      ? session
      : null;
  }

  public IReadOnlyCollection<Session> ListSessionsByTenantIdAndUserId(long tenantId, long userId)
  {
    return _sessionsById.Values
      .Where(session => session.TenantId == tenantId && session.UserId == userId)
      .OrderByDescending(session => session.CreatedAt)
      .ToArray();
  }

  public Session AddSession(Session session)
  {
    _sessionsById[session.Id] = session;
    _sessionsByPublicId[session.PublicId] = session;
    _sessionsByToken[session.SessionToken] = session;
    _sessionsByRefreshToken[session.RefreshToken] = session;
    return session;
  }

  public Session UpdateSession(Session session)
  {
    if (_sessionsById.TryGetValue(session.Id, out var current))
    {
      _sessionsByPublicId.Remove(current.PublicId);
      _sessionsByToken.Remove(current.SessionToken);
      _sessionsByRefreshToken.Remove(current.RefreshToken);
    }

    _sessionsById[session.Id] = session;
    _sessionsByPublicId[session.PublicId] = session;
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

  public PasswordResetToken? FindPasswordResetTokenByResetToken(string resetToken)
  {
    return _passwordResetTokensByToken.TryGetValue(resetToken, out var passwordResetToken)
      ? passwordResetToken
      : null;
  }

  public PasswordResetToken? FindPendingPasswordResetTokenByTenantIdAndUserId(long tenantId, long userId)
  {
    return _passwordResetTokensById.Values
      .Where(passwordResetToken => passwordResetToken.TenantId == tenantId && passwordResetToken.UserId == userId && passwordResetToken.Status == "pending")
      .OrderByDescending(passwordResetToken => passwordResetToken.CreatedAt)
      .FirstOrDefault();
  }

  public PasswordResetToken AddPasswordResetToken(PasswordResetToken passwordResetToken)
  {
    _passwordResetTokensById[passwordResetToken.Id] = passwordResetToken;
    _passwordResetTokensByToken[passwordResetToken.ResetToken] = passwordResetToken;
    return passwordResetToken;
  }

  public PasswordResetToken UpdatePasswordResetToken(PasswordResetToken passwordResetToken)
  {
    _passwordResetTokensById[passwordResetToken.Id] = passwordResetToken;
    _passwordResetTokensByToken[passwordResetToken.ResetToken] = passwordResetToken;
    return passwordResetToken;
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

  public long NextPasswordResetTokenId()
  {
    return _passwordResetTokenSequence++;
  }

  public long NextSecurityAuditEventId()
  {
    return _auditSequence++;
  }
}
