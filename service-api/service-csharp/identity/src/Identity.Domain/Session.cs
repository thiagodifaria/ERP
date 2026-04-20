// Session representa a sessao publica do ERP, isolando o token local do provider externo.
// Isso permite refresh, revogacao e enforcement de tenant sem expor detalhes do Keycloak.
namespace Identity.Domain;

public sealed class Session
{
  public Session(
    long id,
    long tenantId,
    long userId,
    Guid publicId,
    string sessionToken,
    string refreshToken,
    string? identityProviderSubject,
    string? identityProviderRefreshToken,
    string status,
    DateTimeOffset expiresAt,
    DateTimeOffset refreshExpiresAt,
    DateTimeOffset createdAt,
    DateTimeOffset? lastUsedAt,
    DateTimeOffset? revokedAt)
  {
    Id = id;
    TenantId = tenantId;
    UserId = userId;
    PublicId = publicId;
    SessionToken = sessionToken;
    RefreshToken = refreshToken;
    IdentityProviderSubject = identityProviderSubject;
    IdentityProviderRefreshToken = identityProviderRefreshToken;
    Status = status;
    ExpiresAt = expiresAt;
    RefreshExpiresAt = refreshExpiresAt;
    CreatedAt = createdAt;
    LastUsedAt = lastUsedAt;
    RevokedAt = revokedAt;
  }

  public long Id { get; }

  public long TenantId { get; }

  public long UserId { get; }

  public Guid PublicId { get; }

  public string SessionToken { get; }

  public string RefreshToken { get; }

  public string? IdentityProviderSubject { get; }

  public string? IdentityProviderRefreshToken { get; }

  public string Status { get; }

  public DateTimeOffset ExpiresAt { get; }

  public DateTimeOffset RefreshExpiresAt { get; }

  public DateTimeOffset CreatedAt { get; }

  public DateTimeOffset? LastUsedAt { get; }

  public DateTimeOffset? RevokedAt { get; }

  public bool IsActive(DateTimeOffset now)
  {
    return Status == "active" && RevokedAt is null && ExpiresAt >= now;
  }

  public bool CanRefresh(DateTimeOffset now)
  {
    return Status == "active" && RevokedAt is null && RefreshExpiresAt >= now;
  }

  public Session Refresh(
    string refreshToken,
    string? identityProviderRefreshToken,
    DateTimeOffset expiresAt,
    DateTimeOffset refreshExpiresAt,
    DateTimeOffset lastUsedAt)
  {
    return new Session(
      Id,
      TenantId,
      UserId,
      PublicId,
      SessionToken,
      refreshToken,
      IdentityProviderSubject,
      identityProviderRefreshToken,
      Status,
      expiresAt,
      refreshExpiresAt,
      CreatedAt,
      lastUsedAt,
      RevokedAt);
  }

  public Session Touch(DateTimeOffset lastUsedAt)
  {
    return new Session(
      Id,
      TenantId,
      UserId,
      PublicId,
      SessionToken,
      RefreshToken,
      IdentityProviderSubject,
      IdentityProviderRefreshToken,
      Status,
      ExpiresAt,
      RefreshExpiresAt,
      CreatedAt,
      lastUsedAt,
      RevokedAt);
  }

  public Session Revoke(DateTimeOffset revokedAt)
  {
    return new Session(
      Id,
      TenantId,
      UserId,
      PublicId,
      SessionToken,
      RefreshToken,
      IdentityProviderSubject,
      IdentityProviderRefreshToken,
      "revoked",
      ExpiresAt,
      RefreshExpiresAt,
      CreatedAt,
      LastUsedAt,
      revokedAt);
  }
}
