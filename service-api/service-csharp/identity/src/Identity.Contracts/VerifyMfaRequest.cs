// VerifyMfaRequest descreve a confirmacao publica do codigo TOTP.
namespace Identity.Contracts;

public sealed record VerifyMfaRequest(
  string OtpCode);
