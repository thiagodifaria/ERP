namespace Identity.Contracts;

public sealed record UserSessionResponse(
  Guid PublicId,
  Guid UserPublicId,
  string Status,
  DateTimeOffset ExpiresAt,
  DateTimeOffset RefreshExpiresAt,
  DateTimeOffset CreatedAt,
  DateTimeOffset? LastUsedAt,
  DateTimeOffset? RevokedAt);
