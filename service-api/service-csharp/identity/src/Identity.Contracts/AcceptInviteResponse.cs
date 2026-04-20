// AcceptInviteResponse descreve o resultado publico de aceite de convite.
namespace Identity.Contracts;

public sealed record AcceptInviteResponse(
  string TenantSlug,
  Guid InvitePublicId,
  string InviteStatus,
  UserResponse User);
