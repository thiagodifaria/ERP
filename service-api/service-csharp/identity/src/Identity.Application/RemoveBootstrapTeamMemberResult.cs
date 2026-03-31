// Este resultado evita exceptions para validacoes de remocao de membership no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class RemoveBootstrapTeamMemberResult
{
  private RemoveBootstrapTeamMemberResult(
    TeamMembershipResponse? membership,
    ErrorResponse? error,
    bool isBadRequest,
    bool isNotFound)
  {
    Membership = membership;
    Error = error;
    IsBadRequest = isBadRequest;
    IsNotFound = isNotFound;
  }

  public TeamMembershipResponse? Membership { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => Membership is not null;

  public static RemoveBootstrapTeamMemberResult Success(TeamMembershipResponse membership)
  {
    return new RemoveBootstrapTeamMemberResult(membership, null, false, false);
  }

  public static RemoveBootstrapTeamMemberResult BadRequest(ErrorResponse error)
  {
    return new RemoveBootstrapTeamMemberResult(null, error, true, false);
  }

  public static RemoveBootstrapTeamMemberResult NotFound(ErrorResponse error)
  {
    return new RemoveBootstrapTeamMemberResult(null, error, false, true);
  }
}
