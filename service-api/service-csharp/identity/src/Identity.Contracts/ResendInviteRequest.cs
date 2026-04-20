namespace Identity.Contracts;

public sealed record ResendInviteRequest(
  int? ExpiresInDays);
