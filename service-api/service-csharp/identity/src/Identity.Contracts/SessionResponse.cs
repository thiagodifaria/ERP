// SessionResponse descreve a sessao local do ERP e seus tokens opacos.
namespace Identity.Contracts;

public sealed record SessionResponse(
  Guid PublicId,
  string TenantSlug,
  Guid UserPublicId,
  string Email,
  string DisplayName,
  string SessionToken,
  string RefreshToken,
  DateTimeOffset ExpiresAt,
  DateTimeOffset RefreshExpiresAt,
  bool MfaEnabled,
  IReadOnlyCollection<string> RoleCodes);
