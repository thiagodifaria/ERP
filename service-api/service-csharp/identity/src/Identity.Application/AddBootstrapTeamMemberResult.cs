// Este resultado evita exceptions para validacoes de memberships de time no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class AddBootstrapTeamMemberResult
{
  private AddBootstrapTeamMemberResult(
    TeamMembershipResponse? membership,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict,
    bool isNotFound)
  {
    Membership = membership;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
    IsNotFound = isNotFound;
  }

  public TeamMembershipResponse? Membership { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => Membership is not null;

  public static AddBootstrapTeamMemberResult Success(TeamMembershipResponse membership)
  {
    return new AddBootstrapTeamMemberResult(membership, null, false, false, false);
  }

  public static AddBootstrapTeamMemberResult BadRequest(ErrorResponse error)
  {
    return new AddBootstrapTeamMemberResult(null, error, true, false, false);
  }

  public static AddBootstrapTeamMemberResult Conflict(ErrorResponse error)
  {
    return new AddBootstrapTeamMemberResult(null, error, false, true, false);
  }

  public static AddBootstrapTeamMemberResult NotFound(ErrorResponse error)
  {
    return new AddBootstrapTeamMemberResult(null, error, false, false, true);
  }
}
