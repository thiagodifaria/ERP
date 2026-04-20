namespace Identity.Contracts;

public sealed record StartPasswordRecoveryRequest(
  string TenantSlug,
  string Email,
  int? ExpiresInMinutes);
