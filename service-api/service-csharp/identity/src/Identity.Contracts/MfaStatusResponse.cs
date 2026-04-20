// MfaStatusResponse descreve o estado publico atual de MFA do usuario.
namespace Identity.Contracts;

public sealed record MfaStatusResponse(
  Guid UserPublicId,
  bool Enabled);
