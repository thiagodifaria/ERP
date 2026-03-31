// Este resultado evita exceptions para validacoes de update de empresa no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class UpdateBootstrapCompanyResult
{
  private UpdateBootstrapCompanyResult(
    CompanyResponse? company,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict,
    bool isNotFound)
  {
    Company = company;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
    IsNotFound = isNotFound;
  }

  public CompanyResponse? Company { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsNotFound { get; }

  public bool IsSuccess => Company is not null;

  public static UpdateBootstrapCompanyResult Success(CompanyResponse company)
  {
    return new UpdateBootstrapCompanyResult(company, null, false, false, false);
  }

  public static UpdateBootstrapCompanyResult BadRequest(ErrorResponse error)
  {
    return new UpdateBootstrapCompanyResult(null, error, true, false, false);
  }

  public static UpdateBootstrapCompanyResult Conflict(ErrorResponse error)
  {
    return new UpdateBootstrapCompanyResult(null, error, false, true, false);
  }

  public static UpdateBootstrapCompanyResult NotFound(ErrorResponse error)
  {
    return new UpdateBootstrapCompanyResult(null, error, false, false, true);
  }
}
