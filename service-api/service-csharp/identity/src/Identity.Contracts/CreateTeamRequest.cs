// CreateTeamRequest descreve a entrada publica minima para criacao de time.
namespace Identity.Contracts;

public sealed record CreateTeamRequest(
  string Name);
