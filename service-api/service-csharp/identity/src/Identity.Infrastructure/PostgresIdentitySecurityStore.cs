// Este adapter persiste convites, sessoes, MFA e auditoria no PostgreSQL.
using Identity.Domain;
using Npgsql;

namespace Identity.Infrastructure;

public sealed class PostgresIdentitySecurityStore : IIdentitySecurityStore
{
  private readonly string _connectionString;

  public PostgresIdentitySecurityStore(string connectionString)
  {
    _connectionString = connectionString;
  }

  public UserSecurityProfile GetOrCreateProfile(long userId)
  {
    var existing = FindProfile(userId);
    if (existing is not null)
    {
      return existing;
    }

    const string sql = """
      INSERT INTO identity.user_security_profiles (user_id, identity_provider_subject, mfa_enabled, mfa_secret, last_login_at)
      VALUES (@user_id, NULL, false, NULL, NULL)
      RETURNING user_id, identity_provider_subject, mfa_enabled, mfa_secret, last_login_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("user_id", userId);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapProfile(reader);
  }

  public UserSecurityProfile SaveProfile(UserSecurityProfile profile)
  {
    const string sql = """
      INSERT INTO identity.user_security_profiles (user_id, identity_provider_subject, mfa_enabled, mfa_secret, last_login_at)
      VALUES (@user_id, @identity_provider_subject, @mfa_enabled, @mfa_secret, @last_login_at)
      ON CONFLICT (user_id) DO UPDATE
      SET identity_provider_subject = EXCLUDED.identity_provider_subject,
          mfa_enabled = EXCLUDED.mfa_enabled,
          mfa_secret = EXCLUDED.mfa_secret,
          last_login_at = EXCLUDED.last_login_at
      RETURNING user_id, identity_provider_subject, mfa_enabled, mfa_secret, last_login_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("user_id", profile.UserId);
    command.Parameters.AddWithValue("identity_provider_subject", (object?)profile.IdentityProviderSubject ?? DBNull.Value);
    command.Parameters.AddWithValue("mfa_enabled", profile.MfaEnabled);
    command.Parameters.AddWithValue("mfa_secret", (object?)profile.MfaSecret ?? DBNull.Value);
    command.Parameters.AddWithValue("last_login_at", (object?)profile.LastLoginAt ?? DBNull.Value);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapProfile(reader);
  }

  public Invite? FindInviteByToken(string inviteToken)
  {
    const string sql = """
      SELECT id, tenant_id, tenant_slug, user_id, public_id, invite_token, email::text, display_name, role_codes, team_public_ids, status, expires_at, accepted_at, created_at
      FROM identity.invites
      WHERE invite_token = @invite_token
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("invite_token", inviteToken);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapInvite(reader) : null;
  }

  public Invite? FindPendingInviteByTenantIdAndEmail(long tenantId, string email)
  {
    const string sql = """
      SELECT id, tenant_id, tenant_slug, user_id, public_id, invite_token, email::text, display_name, role_codes, team_public_ids, status, expires_at, accepted_at, created_at
      FROM identity.invites
      WHERE tenant_id = @tenant_id
        AND lower(email::text) = lower(@email)
        AND status = 'pending'
      ORDER BY created_at DESC
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("email", email);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapInvite(reader) : null;
  }

  public IReadOnlyCollection<Invite> ListInvitesByTenantId(long tenantId)
  {
    const string sql = """
      SELECT id, tenant_id, tenant_slug, user_id, public_id, invite_token, email::text, display_name, role_codes, team_public_ids, status, expires_at, accepted_at, created_at
      FROM identity.invites
      WHERE tenant_id = @tenant_id
      ORDER BY created_at DESC;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    using var reader = command.ExecuteReader();

    var invites = new List<Invite>();
    while (reader.Read())
    {
      invites.Add(MapInvite(reader));
    }

    return invites;
  }

  public Invite AddInvite(Invite invite)
  {
    const string sql = """
      INSERT INTO identity.invites (tenant_id, tenant_slug, user_id, public_id, invite_token, email, display_name, role_codes, team_public_ids, status, expires_at, accepted_at, created_at)
      VALUES (@tenant_id, @tenant_slug, @user_id, @public_id, @invite_token, @email, @display_name, @role_codes, @team_public_ids, @status, @expires_at, @accepted_at, @created_at)
      RETURNING id, tenant_id, tenant_slug, user_id, public_id, invite_token, email::text, display_name, role_codes, team_public_ids, status, expires_at, accepted_at, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    BindInvite(command, invite);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapInvite(reader);
  }

  public Invite UpdateInvite(Invite invite)
  {
    const string sql = """
      UPDATE identity.invites
      SET display_name = @display_name,
          role_codes = @role_codes,
          team_public_ids = @team_public_ids,
          status = @status,
          expires_at = @expires_at,
          accepted_at = @accepted_at
      WHERE id = @id
      RETURNING id, tenant_id, tenant_slug, user_id, public_id, invite_token, email::text, display_name, role_codes, team_public_ids, status, expires_at, accepted_at, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    BindInvite(command, invite);
    command.Parameters.AddWithValue("id", invite.Id);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapInvite(reader);
  }

  public Session? FindSessionBySessionToken(string sessionToken)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at
      FROM identity.sessions
      WHERE session_token = @session_token
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("session_token", sessionToken);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapSession(reader) : null;
  }

  public Session? FindSessionByRefreshToken(string refreshToken)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at
      FROM identity.sessions
      WHERE refresh_token = @refresh_token
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("refresh_token", refreshToken);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapSession(reader) : null;
  }

  public Session? FindSessionByPublicId(Guid sessionPublicId)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at
      FROM identity.sessions
      WHERE public_id = @public_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("public_id", sessionPublicId);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapSession(reader) : null;
  }

  public IReadOnlyCollection<Session> ListSessionsByTenantIdAndUserId(long tenantId, long userId)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at
      FROM identity.sessions
      WHERE tenant_id = @tenant_id
        AND user_id = @user_id
      ORDER BY created_at DESC;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("user_id", userId);
    using var reader = command.ExecuteReader();

    var sessions = new List<Session>();
    while (reader.Read())
    {
      sessions.Add(MapSession(reader));
    }

    return sessions;
  }

  public Session AddSession(Session session)
  {
    const string sql = """
      INSERT INTO identity.sessions (tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at)
      VALUES (@tenant_id, @user_id, @public_id, @session_token, @refresh_token, @identity_provider_subject, @identity_provider_refresh_token, @status, @expires_at, @refresh_expires_at, @created_at, @last_used_at, @revoked_at)
      RETURNING id, tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    BindSession(command, session);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapSession(reader);
  }

  public Session UpdateSession(Session session)
  {
    const string sql = """
      UPDATE identity.sessions
      SET session_token = @session_token,
          refresh_token = @refresh_token,
          identity_provider_subject = @identity_provider_subject,
          identity_provider_refresh_token = @identity_provider_refresh_token,
          status = @status,
          expires_at = @expires_at,
          refresh_expires_at = @refresh_expires_at,
          last_used_at = @last_used_at,
          revoked_at = @revoked_at
      WHERE id = @id
      RETURNING id, tenant_id, user_id, public_id, session_token, refresh_token, identity_provider_subject, identity_provider_refresh_token, status, expires_at, refresh_expires_at, created_at, last_used_at, revoked_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    BindSession(command, session);
    command.Parameters.AddWithValue("id", session.Id);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapSession(reader);
  }

  public int RevokeSessionsByUserId(long tenantId, long userId, DateTimeOffset revokedAt)
  {
    const string sql = """
      UPDATE identity.sessions
      SET status = 'revoked',
          revoked_at = @revoked_at
      WHERE tenant_id = @tenant_id
        AND user_id = @user_id
        AND status = 'active';
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("user_id", userId);
    command.Parameters.AddWithValue("revoked_at", revokedAt);
    return command.ExecuteNonQuery();
  }

  public PasswordResetToken? FindPasswordResetTokenByResetToken(string resetToken)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, public_id, reset_token, status, expires_at, consumed_at, created_at
      FROM identity.password_reset_tokens
      WHERE reset_token = @reset_token
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("reset_token", resetToken);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapPasswordResetToken(reader) : null;
  }

  public PasswordResetToken? FindPendingPasswordResetTokenByTenantIdAndUserId(long tenantId, long userId)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, public_id, reset_token, status, expires_at, consumed_at, created_at
      FROM identity.password_reset_tokens
      WHERE tenant_id = @tenant_id
        AND user_id = @user_id
        AND status = 'pending'
      ORDER BY created_at DESC
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("user_id", userId);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapPasswordResetToken(reader) : null;
  }

  public PasswordResetToken AddPasswordResetToken(PasswordResetToken passwordResetToken)
  {
    const string sql = """
      INSERT INTO identity.password_reset_tokens (tenant_id, user_id, public_id, reset_token, status, expires_at, consumed_at, created_at)
      VALUES (@tenant_id, @user_id, @public_id, @reset_token, @status, @expires_at, @consumed_at, @created_at)
      RETURNING id, tenant_id, user_id, public_id, reset_token, status, expires_at, consumed_at, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    BindPasswordResetToken(command, passwordResetToken);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapPasswordResetToken(reader);
  }

  public PasswordResetToken UpdatePasswordResetToken(PasswordResetToken passwordResetToken)
  {
    const string sql = """
      UPDATE identity.password_reset_tokens
      SET status = @status,
          expires_at = @expires_at,
          consumed_at = @consumed_at
      WHERE id = @id
      RETURNING id, tenant_id, user_id, public_id, reset_token, status, expires_at, consumed_at, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    BindPasswordResetToken(command, passwordResetToken);
    command.Parameters.AddWithValue("id", passwordResetToken.Id);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapPasswordResetToken(reader);
  }

  public IReadOnlyCollection<SecurityAuditEvent> ListSecurityAuditByTenantId(long tenantId, int limit)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, actor_user_public_id, subject_user_public_id, event_code, severity, summary, created_at
      FROM identity.security_audit_events
      WHERE tenant_id = @tenant_id
      ORDER BY created_at DESC
      LIMIT @limit;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("limit", limit);
    using var reader = command.ExecuteReader();

    var events = new List<SecurityAuditEvent>();
    while (reader.Read())
    {
      events.Add(MapSecurityAuditEvent(reader));
    }

    return events;
  }

  public SecurityAuditEvent AddSecurityAuditEvent(SecurityAuditEvent auditEvent)
  {
    const string sql = """
      INSERT INTO identity.security_audit_events (tenant_id, public_id, actor_user_public_id, subject_user_public_id, event_code, severity, summary, created_at)
      VALUES (@tenant_id, @public_id, @actor_user_public_id, @subject_user_public_id, @event_code, @severity, @summary, @created_at)
      RETURNING id, tenant_id, public_id, actor_user_public_id, subject_user_public_id, event_code, severity, summary, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", auditEvent.TenantId);
    command.Parameters.AddWithValue("public_id", auditEvent.PublicId);
    command.Parameters.AddWithValue("actor_user_public_id", (object?)auditEvent.ActorUserPublicId ?? DBNull.Value);
    command.Parameters.AddWithValue("subject_user_public_id", (object?)auditEvent.SubjectUserPublicId ?? DBNull.Value);
    command.Parameters.AddWithValue("event_code", auditEvent.EventCode);
    command.Parameters.AddWithValue("severity", auditEvent.Severity);
    command.Parameters.AddWithValue("summary", auditEvent.Summary);
    command.Parameters.AddWithValue("created_at", auditEvent.CreatedAt);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapSecurityAuditEvent(reader);
  }

  public long NextInviteId()
  {
    return NextId("identity.invites");
  }

  public long NextSessionId()
  {
    return NextId("identity.sessions");
  }

  public long NextPasswordResetTokenId()
  {
    return NextId("identity.password_reset_tokens");
  }

  public long NextSecurityAuditEventId()
  {
    return NextId("identity.security_audit_events");
  }

  private UserSecurityProfile? FindProfile(long userId)
  {
    const string sql = """
      SELECT user_id, identity_provider_subject, mfa_enabled, mfa_secret, last_login_at
      FROM identity.user_security_profiles
      WHERE user_id = @user_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("user_id", userId);
    using var reader = command.ExecuteReader();

    return reader.Read() ? MapProfile(reader) : null;
  }

  private static void BindInvite(NpgsqlCommand command, Invite invite)
  {
    command.Parameters.AddWithValue("tenant_id", invite.TenantId);
    command.Parameters.AddWithValue("tenant_slug", invite.TenantSlug);
    command.Parameters.AddWithValue("user_id", invite.UserId);
    command.Parameters.AddWithValue("public_id", invite.PublicId);
    command.Parameters.AddWithValue("invite_token", invite.InviteToken);
    command.Parameters.AddWithValue("email", invite.Email);
    command.Parameters.AddWithValue("display_name", (object?)invite.DisplayName ?? DBNull.Value);
    command.Parameters.AddWithValue("role_codes", invite.RoleCodes.ToArray());
    command.Parameters.AddWithValue("team_public_ids", invite.TeamPublicIds.ToArray());
    command.Parameters.AddWithValue("status", invite.Status);
    command.Parameters.AddWithValue("expires_at", invite.ExpiresAt);
    command.Parameters.AddWithValue("accepted_at", (object?)invite.AcceptedAt ?? DBNull.Value);
    command.Parameters.AddWithValue("created_at", invite.CreatedAt);
  }

  private static void BindSession(NpgsqlCommand command, Session session)
  {
    command.Parameters.AddWithValue("tenant_id", session.TenantId);
    command.Parameters.AddWithValue("user_id", session.UserId);
    command.Parameters.AddWithValue("public_id", session.PublicId);
    command.Parameters.AddWithValue("session_token", session.SessionToken);
    command.Parameters.AddWithValue("refresh_token", session.RefreshToken);
    command.Parameters.AddWithValue("identity_provider_subject", (object?)session.IdentityProviderSubject ?? DBNull.Value);
    command.Parameters.AddWithValue("identity_provider_refresh_token", (object?)session.IdentityProviderRefreshToken ?? DBNull.Value);
    command.Parameters.AddWithValue("status", session.Status);
    command.Parameters.AddWithValue("expires_at", session.ExpiresAt);
    command.Parameters.AddWithValue("refresh_expires_at", session.RefreshExpiresAt);
    command.Parameters.AddWithValue("created_at", session.CreatedAt);
    command.Parameters.AddWithValue("last_used_at", (object?)session.LastUsedAt ?? DBNull.Value);
    command.Parameters.AddWithValue("revoked_at", (object?)session.RevokedAt ?? DBNull.Value);
  }

  private static void BindPasswordResetToken(NpgsqlCommand command, PasswordResetToken passwordResetToken)
  {
    command.Parameters.AddWithValue("tenant_id", passwordResetToken.TenantId);
    command.Parameters.AddWithValue("user_id", passwordResetToken.UserId);
    command.Parameters.AddWithValue("public_id", passwordResetToken.PublicId);
    command.Parameters.AddWithValue("reset_token", passwordResetToken.ResetToken);
    command.Parameters.AddWithValue("status", passwordResetToken.Status);
    command.Parameters.AddWithValue("expires_at", passwordResetToken.ExpiresAt);
    command.Parameters.AddWithValue("consumed_at", (object?)passwordResetToken.ConsumedAt ?? DBNull.Value);
    command.Parameters.AddWithValue("created_at", passwordResetToken.CreatedAt);
  }

  private long NextId(string tableName)
  {
    using var connection = OpenConnection();
    using var command = new NpgsqlCommand($"SELECT nextval(pg_get_serial_sequence('{tableName}', 'id'));", connection);
    return Convert.ToInt64(command.ExecuteScalar());
  }

  private NpgsqlConnection OpenConnection()
  {
    var connection = new NpgsqlConnection(_connectionString);
    connection.Open();
    return connection;
  }

  private static UserSecurityProfile MapProfile(NpgsqlDataReader reader)
  {
    return new UserSecurityProfile(
      reader.GetInt64(0),
      reader.IsDBNull(1) ? null : reader.GetString(1),
      reader.GetBoolean(2),
      reader.IsDBNull(3) ? null : reader.GetString(3),
      reader.IsDBNull(4) ? null : reader.GetFieldValue<DateTimeOffset>(4));
  }

  private static Invite MapInvite(NpgsqlDataReader reader)
  {
    return new Invite(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetString(2),
      reader.GetInt64(3),
      reader.GetGuid(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.IsDBNull(7) ? null : reader.GetString(7),
      reader.GetFieldValue<string[]>(8),
      reader.GetFieldValue<Guid[]>(9),
      reader.GetString(10),
      reader.GetFieldValue<DateTimeOffset>(11),
      reader.IsDBNull(12) ? null : reader.GetFieldValue<DateTimeOffset>(12),
      reader.GetFieldValue<DateTimeOffset>(13));
  }

  private static Session MapSession(NpgsqlDataReader reader)
  {
    return new Session(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetInt64(2),
      reader.GetGuid(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.IsDBNull(6) ? null : reader.GetString(6),
      reader.IsDBNull(7) ? null : reader.GetString(7),
      reader.GetString(8),
      reader.GetFieldValue<DateTimeOffset>(9),
      reader.GetFieldValue<DateTimeOffset>(10),
      reader.GetFieldValue<DateTimeOffset>(11),
      reader.IsDBNull(12) ? null : reader.GetFieldValue<DateTimeOffset>(12),
      reader.IsDBNull(13) ? null : reader.GetFieldValue<DateTimeOffset>(13));
  }

  private static PasswordResetToken MapPasswordResetToken(NpgsqlDataReader reader)
  {
    return new PasswordResetToken(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetInt64(2),
      reader.GetGuid(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetFieldValue<DateTimeOffset>(6),
      reader.IsDBNull(7) ? null : reader.GetFieldValue<DateTimeOffset>(7),
      reader.GetFieldValue<DateTimeOffset>(8));
  }

  private static SecurityAuditEvent MapSecurityAuditEvent(NpgsqlDataReader reader)
  {
    return new SecurityAuditEvent(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetGuid(2),
      reader.IsDBNull(3) ? null : reader.GetGuid(3),
      reader.IsDBNull(4) ? null : reader.GetGuid(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetFieldValue<DateTimeOffset>(8));
  }
}
