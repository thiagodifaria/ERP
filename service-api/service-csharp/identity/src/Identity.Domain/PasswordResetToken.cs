// PasswordResetToken representa o fluxo local de recuperacao de senha do tenant.
namespace Identity.Domain;

public sealed class PasswordResetToken
{
  public PasswordResetToken(
    long id,
    long tenantId,
    long userId,
    Guid publicId,
    string resetToken,
    string status,
    DateTimeOffset expiresAt,
    DateTimeOffset? consumedAt,
    DateTimeOffset createdAt)
  {
    Id = id;
    TenantId = tenantId;
    UserId = userId;
    PublicId = publicId;
    ResetToken = resetToken;
    Status = status;
    ExpiresAt = expiresAt;
    ConsumedAt = consumedAt;
    CreatedAt = createdAt;
  }

  public long Id { get; }

  public long TenantId { get; }

  public long UserId { get; }

  public Guid PublicId { get; }

  public string ResetToken { get; }

  public string Status { get; }

  public DateTimeOffset ExpiresAt { get; }

  public DateTimeOffset? ConsumedAt { get; }

  public DateTimeOffset CreatedAt { get; }

  public bool IsExpired(DateTimeOffset now)
  {
    return ExpiresAt < now;
  }

  public PasswordResetToken Consume(DateTimeOffset consumedAt)
  {
    return new PasswordResetToken(
      Id,
      TenantId,
      UserId,
      PublicId,
      ResetToken,
      "consumed",
      ExpiresAt,
      consumedAt,
      CreatedAt);
  }

  public PasswordResetToken Expire()
  {
    return new PasswordResetToken(
      Id,
      TenantId,
      UserId,
      PublicId,
      ResetToken,
      "expired",
      ExpiresAt,
      ConsumedAt,
      CreatedAt);
  }
}
