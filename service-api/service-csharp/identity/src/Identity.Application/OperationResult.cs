// OperationResult reduz repeticao em casos de uso novos sem recorrer a exceptions.
using Identity.Contracts;

namespace Identity.Application;

public sealed class OperationResult<T>
{
  private OperationResult(
    T? payload,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict,
    bool isForbidden,
    bool isNotFound,
    bool isUnauthorized)
  {
    Payload = payload;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
    IsForbidden = isForbidden;
    IsNotFound = isNotFound;
    IsUnauthorized = isUnauthorized;
  }

  public T? Payload { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsForbidden { get; }

  public bool IsNotFound { get; }

  public bool IsUnauthorized { get; }

  public bool IsSuccess => Payload is not null;

  public static OperationResult<T> Success(T payload)
  {
    return new OperationResult<T>(payload, null, false, false, false, false, false);
  }

  public static OperationResult<T> BadRequest(ErrorResponse error)
  {
    return new OperationResult<T>(default, error, true, false, false, false, false);
  }

  public static OperationResult<T> Conflict(ErrorResponse error)
  {
    return new OperationResult<T>(default, error, false, true, false, false, false);
  }

  public static OperationResult<T> Forbidden(ErrorResponse error)
  {
    return new OperationResult<T>(default, error, false, false, true, false, false);
  }

  public static OperationResult<T> NotFound(ErrorResponse error)
  {
    return new OperationResult<T>(default, error, false, false, false, true, false);
  }

  public static OperationResult<T> Unauthorized(ErrorResponse error)
  {
    return new OperationResult<T>(default, error, false, false, false, false, true);
  }
}
