// RoleResponse descreve a leitura publica minima de um papel de acesso.
namespace Identity.Contracts;

public sealed record RoleResponse(
  long Id,
  Guid PublicId,
  string Code,
  string DisplayName,
  string Status);
