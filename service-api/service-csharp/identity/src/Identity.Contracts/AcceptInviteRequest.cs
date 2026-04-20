// AcceptInviteRequest descreve a ativacao publica de um convite pendente.
namespace Identity.Contracts;

public sealed record AcceptInviteRequest(
  string DisplayName,
  string? GivenName,
  string? FamilyName,
  string Password);
