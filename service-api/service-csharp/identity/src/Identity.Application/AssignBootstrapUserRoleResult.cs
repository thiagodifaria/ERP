// Este resultado evita exceptions para validacoes de atribuicao de papeis no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class AssignBootstrapUserRoleResult
{
  private AssignBootstrapUserRoleResult(
    UserRoleResponse? userRole,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict,
    bool isNotFound)
  {
    UserRole = userRole;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
    IsNotFound = isNotFound;
  }

  public UserRoleResponse? UserRole { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => UserRole is not null;

  public static AssignBootstrapUserRoleResult Success(UserRoleResponse userRole)
  {
    return new AssignBootstrapUserRoleResult(userRole, null, false, false, false);
  }

  public static AssignBootstrapUserRoleResult BadRequest(ErrorResponse error)
  {
    return new AssignBootstrapUserRoleResult(null, error, true, false, false);
  }

  public static AssignBootstrapUserRoleResult Conflict(ErrorResponse error)
  {
    return new AssignBootstrapUserRoleResult(null, error, false, true, false);
  }

  public static AssignBootstrapUserRoleResult NotFound(ErrorResponse error)
  {
    return new AssignBootstrapUserRoleResult(null, error, false, false, true);
  }
}
