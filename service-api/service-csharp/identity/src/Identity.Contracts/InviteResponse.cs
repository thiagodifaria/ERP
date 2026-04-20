// InviteResponse descreve o convite atual e seu token publico de aceite.
namespace Identity.Contracts;

public sealed record InviteResponse(
  Guid PublicId,
  Guid UserPublicId,
  string Email,
  string? DisplayName,
  IReadOnlyCollection<string> RoleCodes,
  IReadOnlyCollection<Guid> TeamPublicIds,
  string Status,
  string InviteToken,
  DateTimeOffset ExpiresAt);
