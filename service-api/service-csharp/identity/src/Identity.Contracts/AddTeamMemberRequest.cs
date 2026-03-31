// AddTeamMemberRequest descreve a entrada publica minima para adicionar um membro ao time.
namespace Identity.Contracts;

public sealed record AddTeamMemberRequest(
  Guid UserPublicId);
