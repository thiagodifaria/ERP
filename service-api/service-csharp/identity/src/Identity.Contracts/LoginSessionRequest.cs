// LoginSessionRequest descreve a autenticacao publica com tenant explicito.
namespace Identity.Contracts;

public sealed record LoginSessionRequest(
  string TenantSlug,
  string Email,
  string Password,
  string? OtpCode);
