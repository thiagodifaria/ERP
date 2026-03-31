// CreateUserRequest descreve a entrada publica minima para criacao de usuario.
namespace Identity.Contracts;

public sealed record CreateUserRequest(
  string Email,
  string DisplayName,
  string? GivenName,
  string? FamilyName);
