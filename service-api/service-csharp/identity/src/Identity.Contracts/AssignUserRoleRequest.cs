// AssignUserRoleRequest descreve a entrada publica minima para atribuicao de papel a usuario.
namespace Identity.Contracts;

public sealed record AssignUserRoleRequest(
  string RoleCode);
