// Este resultado evita exceptions para validacoes de revogacao de papeis no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class RevokeBootstrapUserRoleResult
{
  private RevokeBootstrapUserRoleResult(
    UserRoleResponse? userRole,
    ErrorResponse? error,
    bool isBadRequest,
    bool isNotFound)
  {
    UserRole = userRole;
    Error = error;
    IsBadRequest = isBadRequest;
    IsNotFound = isNotFound;
  }

  public UserRoleResponse? UserRole { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => UserRole is not null;

  public static RevokeBootstrapUserRoleResult Success(UserRoleResponse userRole)
  {
    return new RevokeBootstrapUserRoleResult(userRole, null, false, false);
  }

  public static RevokeBootstrapUserRoleResult BadRequest(ErrorResponse error)
  {
    return new RevokeBootstrapUserRoleResult(null, error, true, false);
  }

  public static RevokeBootstrapUserRoleResult NotFound(ErrorResponse error)
  {
    return new RevokeBootstrapUserRoleResult(null, error, false, true);
  }
}
