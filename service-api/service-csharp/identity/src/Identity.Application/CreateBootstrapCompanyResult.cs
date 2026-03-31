// Este resultado evita exceptions para validacoes de criacao de empresa no bootstrap.
using Identity.Contracts;

namespace Identity.Application;

public sealed class CreateBootstrapCompanyResult
{
  private CreateBootstrapCompanyResult(
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

  public static CreateBootstrapCompanyResult Success(CompanyResponse company)
  {
    return new CreateBootstrapCompanyResult(company, null, false, false, false);
  }

  public static CreateBootstrapCompanyResult BadRequest(ErrorResponse error)
  {
    return new CreateBootstrapCompanyResult(null, error, true, false, false);
  }

  public static CreateBootstrapCompanyResult Conflict(ErrorResponse error)
  {
    return new CreateBootstrapCompanyResult(null, error, false, true, false);
  }

  public static CreateBootstrapCompanyResult NotFound(ErrorResponse error)
  {
    return new CreateBootstrapCompanyResult(null, error, false, false, true);
  }
}
