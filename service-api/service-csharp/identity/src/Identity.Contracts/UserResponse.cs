// UserResponse descreve a saida publica minima para leitura de usuarios.
namespace Identity.Contracts;

public sealed record UserResponse(
  long Id,
  Guid PublicId,
  long TenantId,
  long? CompanyId,
  string Email,
  string DisplayName,
  string? GivenName,
  string? FamilyName,
  string Status);
