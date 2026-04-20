// Invite representa um convite de acesso ainda pendente ou historico.
// O token publico funciona como chave de aceite enquanto nao houver canal de envio real.
namespace Identity.Domain;

public sealed class Invite
{
  public Invite(
    long id,
    long tenantId,
    string tenantSlug,
    long userId,
    Guid publicId,
    string inviteToken,
    string email,
    string? displayName,
    IReadOnlyCollection<string> roleCodes,
    IReadOnlyCollection<Guid> teamPublicIds,
    string status,
    DateTimeOffset expiresAt,
    DateTimeOffset? acceptedAt,
    DateTimeOffset createdAt)
  {
    Id = id;
    TenantId = tenantId;
    TenantSlug = tenantSlug;
    UserId = userId;
    PublicId = publicId;
    InviteToken = inviteToken;
    Email = email;
    DisplayName = displayName;
    RoleCodes = roleCodes;
    TeamPublicIds = teamPublicIds;
    Status = status;
    ExpiresAt = expiresAt;
    AcceptedAt = acceptedAt;
    CreatedAt = createdAt;
  }

  public long Id { get; }

  public long TenantId { get; }

  public string TenantSlug { get; }

  public long UserId { get; }

  public Guid PublicId { get; }

  public string InviteToken { get; }

  public string Email { get; }

  public string? DisplayName { get; }

  public IReadOnlyCollection<string> RoleCodes { get; }

  public IReadOnlyCollection<Guid> TeamPublicIds { get; }

  public string Status { get; }

  public DateTimeOffset ExpiresAt { get; }

  public DateTimeOffset? AcceptedAt { get; }

  public DateTimeOffset CreatedAt { get; }

  public bool IsExpired(DateTimeOffset now)
  {
    return ExpiresAt < now;
  }

  public Invite Accept(DateTimeOffset acceptedAt)
  {
    return new Invite(
      Id,
      TenantId,
      TenantSlug,
      UserId,
      PublicId,
      InviteToken,
      Email,
      DisplayName,
      RoleCodes,
      TeamPublicIds,
      "accepted",
      ExpiresAt,
      acceptedAt,
      CreatedAt);
  }

  public Invite Expire()
  {
    return new Invite(
      Id,
      TenantId,
      TenantSlug,
      UserId,
      PublicId,
      InviteToken,
      Email,
      DisplayName,
      RoleCodes,
      TeamPublicIds,
      "expired",
      ExpiresAt,
      AcceptedAt,
      CreatedAt);
  }
}
