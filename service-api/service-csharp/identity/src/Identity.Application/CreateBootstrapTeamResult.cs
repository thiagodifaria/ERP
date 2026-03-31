// Este resultado evita exceptions para validacoes de criacao de time no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class CreateBootstrapTeamResult
{
  private CreateBootstrapTeamResult(
    TeamResponse? team,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict,
    bool isNotFound)
  {
    Team = team;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
    IsNotFound = isNotFound;
  }

  public TeamResponse? Team { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => Team is not null;

  public static CreateBootstrapTeamResult Success(TeamResponse team)
  {
    return new CreateBootstrapTeamResult(team, null, false, false, false);
  }

  public static CreateBootstrapTeamResult BadRequest(ErrorResponse error)
  {
    return new CreateBootstrapTeamResult(null, error, true, false, false);
  }

  public static CreateBootstrapTeamResult Conflict(ErrorResponse error)
  {
    return new CreateBootstrapTeamResult(null, error, false, true, false);
  }

  public static CreateBootstrapTeamResult NotFound(ErrorResponse error)
  {
    return new CreateBootstrapTeamResult(null, error, false, false, true);
  }
}
