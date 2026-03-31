// Este resultado evita exceptions para validacoes de update de usuario no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class UpdateBootstrapUserResult
{
  private UpdateBootstrapUserResult(
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

  public static UpdateBootstrapUserResult Success(UserResponse user)
  {
    return new UpdateBootstrapUserResult(user, null, false, false, false);
  }

  public static UpdateBootstrapUserResult BadRequest(ErrorResponse error)
  {
    return new UpdateBootstrapUserResult(null, error, true, false, false);
  }

  public static UpdateBootstrapUserResult Conflict(ErrorResponse error)
  {
    return new UpdateBootstrapUserResult(null, error, false, true, false);
  }

  public static UpdateBootstrapUserResult NotFound(ErrorResponse error)
  {
    return new UpdateBootstrapUserResult(null, error, false, false, true);
  }
}
