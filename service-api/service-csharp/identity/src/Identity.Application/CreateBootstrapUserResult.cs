// Este resultado evita exceptions para validacoes de criacao de usuario no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class CreateBootstrapUserResult
{
  private CreateBootstrapUserResult(
    UserResponse? user,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict,
    bool isNotFound)
  {
    User = user;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
    IsNotFound = isNotFound;
  }

  public UserResponse? User { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => User is not null;

  public static CreateBootstrapUserResult Success(UserResponse user)
  {
    return new CreateBootstrapUserResult(user, null, false, false, false);
  }

  public static CreateBootstrapUserResult BadRequest(ErrorResponse error)
  {
    return new CreateBootstrapUserResult(null, error, true, false, false);
  }

  public static CreateBootstrapUserResult Conflict(ErrorResponse error)
  {
    return new CreateBootstrapUserResult(null, error, false, true, false);
  }

  public static CreateBootstrapUserResult NotFound(ErrorResponse error)
  {
    return new CreateBootstrapUserResult(null, error, false, false, true);
  }
}
