namespace Identity.Contracts;

public sealed record PasswordRecoveryResponse(
  Guid PublicId,
  string TenantSlug,
  Guid UserPublicId,
  string Email,
  string Status,
  string ResetToken,
  DateTimeOffset ExpiresAt);
