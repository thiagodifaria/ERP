// UserAccessSnapshotResponse agrupa dados do usuario com seus papeis atribuidos.
namespace Identity.Contracts;

public sealed record UserAccessSnapshotResponse(
  UserResponse User,
  IReadOnlyCollection<UserRoleResponse> Roles);
