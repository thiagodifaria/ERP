// UpdateUserAccessRequest descreve o bloqueio ou reativacao publica de acesso.
namespace Identity.Contracts;

public sealed record UpdateUserAccessRequest(
  string Status);
