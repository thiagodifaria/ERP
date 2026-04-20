// MfaEnrollmentResponse descreve o estado inicial de configuracao de MFA.
namespace Identity.Contracts;

public sealed record MfaEnrollmentResponse(
  Guid UserPublicId,
  bool Enabled,
  string Secret,
  string OtpAuthUri);
