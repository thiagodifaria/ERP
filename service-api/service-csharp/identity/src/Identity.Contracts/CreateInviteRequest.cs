// CreateInviteRequest descreve o convite minimo de onboarding por tenant.
namespace Identity.Contracts;

public sealed record CreateInviteRequest(
  string Email,
  string? DisplayName,
  IReadOnlyCollection<string>? RoleCodes,
  IReadOnlyCollection<Guid>? TeamPublicIds,
  int? ExpiresInDays);
