// TeamAccessSnapshotResponse agrupa dados do time com seus membros atuais.
namespace Identity.Contracts;

public sealed record TeamAccessSnapshotResponse(
  TeamResponse Team,
  IReadOnlyCollection<TeamMembershipResponse> Members);
