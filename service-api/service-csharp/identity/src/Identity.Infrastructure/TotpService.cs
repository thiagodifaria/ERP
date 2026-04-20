// TotpService implementa TOTP simples para MFA local do contexto de identidade.
using System.Security.Cryptography;
using System.Text;
using Identity.Application;

namespace Identity.Infrastructure;

public sealed class TotpService : ITotpService
{
  private const string Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";

  public string GenerateSecret()
  {
    var bytes = RandomNumberGenerator.GetBytes(20);
    return Base32Encode(bytes);
  }

  public bool VerifyCode(string secret, string otpCode, DateTimeOffset now)
  {
    if (string.IsNullOrWhiteSpace(secret) || string.IsNullOrWhiteSpace(otpCode))
    {
      return false;
    }

    var normalizedCode = otpCode.Trim();
    if (normalizedCode.Length != 6 || normalizedCode.Any(character => !char.IsDigit(character)))
    {
      return false;
    }

    var secretBytes = Base32Decode(secret);
    foreach (var offset in new[] { -1L, 0L, 1L })
    {
      var code = ComputeCode(secretBytes, now.ToUnixTimeSeconds() / 30 + offset);
      if (code.Equals(normalizedCode, StringComparison.Ordinal))
      {
        return true;
      }
    }

    return false;
  }

  public string BuildOtpAuthUri(string issuer, string accountName, string secret)
  {
    var encodedIssuer = Uri.EscapeDataString(issuer);
    var encodedAccountName = Uri.EscapeDataString(accountName);
    return $"otpauth://totp/{encodedIssuer}:{encodedAccountName}?secret={secret}&issuer={encodedIssuer}&digits=6&period=30";
  }

  private static string ComputeCode(byte[] secret, long counter)
  {
    Span<byte> counterBytes = stackalloc byte[8];
    for (var index = 7; index >= 0; index--)
    {
      counterBytes[index] = (byte)(counter & 0xff);
      counter >>= 8;
    }

    using var hmac = new HMACSHA1(secret);
    var hash = hmac.ComputeHash(counterBytes.ToArray());
    var offset = hash[^1] & 0x0f;
    var binaryCode =
      ((hash[offset] & 0x7f) << 24)
      | (hash[offset + 1] << 16)
      | (hash[offset + 2] << 8)
      | hash[offset + 3];

    return (binaryCode % 1_000_000).ToString("D6");
  }

  private static string Base32Encode(byte[] data)
  {
    var output = new StringBuilder((data.Length + 4) / 5 * 8);
    var buffer = 0;
    var bitsLeft = 0;

    foreach (var value in data)
    {
      buffer = (buffer << 8) | value;
      bitsLeft += 8;

      while (bitsLeft >= 5)
      {
        output.Append(Alphabet[(buffer >> (bitsLeft - 5)) & 0x1f]);
        bitsLeft -= 5;
      }
    }

    if (bitsLeft > 0)
    {
      output.Append(Alphabet[(buffer << (5 - bitsLeft)) & 0x1f]);
    }

    return output.ToString();
  }

  private static byte[] Base32Decode(string secret)
  {
    var normalized = secret.Trim().Replace("=", string.Empty, StringComparison.Ordinal).ToUpperInvariant();
    var output = new List<byte>((normalized.Length * 5) / 8);
    var buffer = 0;
    var bitsLeft = 0;

    foreach (var character in normalized)
    {
      var index = Alphabet.IndexOf(character);
      if (index < 0)
      {
        continue;
      }

      buffer = (buffer << 5) | index;
      bitsLeft += 5;

      if (bitsLeft < 8)
      {
        continue;
      }

      output.Add((byte)((buffer >> (bitsLeft - 8)) & 0xff));
      bitsLeft -= 8;
    }

    return output.ToArray();
  }
}
