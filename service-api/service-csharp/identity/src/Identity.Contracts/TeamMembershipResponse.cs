// TeamMembershipResponse descreve a saida publica minima de membros por time.
namespace Identity.Contracts;

public sealed record TeamMembershipResponse(
  long Id,
  long TenantId,
  Guid TeamPublicId,
  Guid UserPublicId,
  string UserEmail,
  string UserDisplayName,
  DateTimeOffset CreatedAt);
