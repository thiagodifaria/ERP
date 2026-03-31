// PublicIds centraliza a geracao de identificadores publicos do contexto.
// O formato textual segue o padrao UUIDv7.
using System.Security.Cryptography;

namespace Identity.Domain;

public static class PublicIds
{
  public static Guid NewUuidV7()
  {
    var timestamp = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds().ToString("x12");
    var randomA = RandomNumberGenerator.GetHexString(3).ToLowerInvariant();
    var randomB = RandomNumberGenerator.GetHexString(3).ToLowerInvariant();
    var randomC = RandomNumberGenerator.GetHexString(12).ToLowerInvariant();
    var variant = "89ab"[RandomNumberGenerator.GetInt32(0, 4)];

    return Guid.Parse(
      $"{timestamp[..8]}-{timestamp[8..12]}-7{randomA}-{variant}{randomB}-{randomC}");
  }
}
