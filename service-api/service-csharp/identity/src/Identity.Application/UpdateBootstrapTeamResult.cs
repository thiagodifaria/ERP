// Este resultado evita exceptions para validacoes de update de time no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class UpdateBootstrapTeamResult
{
  private UpdateBootstrapTeamResult(
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

  public static UpdateBootstrapTeamResult Success(TeamResponse team)
  {
    return new UpdateBootstrapTeamResult(team, null, false, false, false);
  }

  public static UpdateBootstrapTeamResult BadRequest(ErrorResponse error)
  {
    return new UpdateBootstrapTeamResult(null, error, true, false, false);
  }

  public static UpdateBootstrapTeamResult Conflict(ErrorResponse error)
  {
    return new UpdateBootstrapTeamResult(null, error, false, true, false);
  }

  public static UpdateBootstrapTeamResult NotFound(ErrorResponse error)
  {
    return new UpdateBootstrapTeamResult(null, error, false, false, true);
  }
}
