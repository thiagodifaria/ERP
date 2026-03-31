// UpdateUserRequest descreve a entrada publica minima para atualizar usuario existente.
namespace Identity.Contracts;

public sealed record UpdateUserRequest(
  string Email,
  string DisplayName,
  string? GivenName,
  string? FamilyName);
