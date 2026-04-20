// ITotpService define a verificacao minima de MFA por TOTP.
namespace Identity.Application;

public interface ITotpService
{
  string GenerateSecret();

  bool VerifyCode(string secret, string otpCode, DateTimeOffset now);

  string BuildOtpAuthUri(string issuer, string accountName, string secret);
}
