// AccessResolutionResponse entrega o contexto minimo de autorizacao para a borda.
namespace Identity.Contracts;

public sealed record AccessResolutionResponse(
  string TenantSlug,
  Guid SessionPublicId,
  Guid UserPublicId,
  string Email,
  string DisplayName,
  IReadOnlyCollection<string> RoleCodes,
  bool MfaEnabled,
  bool Authorized,
  string Status);
