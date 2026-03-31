// UserRoleResponse descreve a saida publica minima de papeis atribuidos a um usuario.
namespace Identity.Contracts;

public sealed record UserRoleResponse(
  long Id,
  long TenantId,
  Guid UserPublicId,
  Guid RolePublicId,
  string RoleCode,
  string RoleDisplayName,
  DateTimeOffset CreatedAt);
