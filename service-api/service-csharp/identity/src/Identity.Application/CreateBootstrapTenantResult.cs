// Este resultado evita exceptions para validacoes de contrato no bootstrap do servico.
using Identity.Contracts;

namespace Identity.Application;

public sealed class CreateBootstrapTenantResult
{
  private CreateBootstrapTenantResult(
    TenantResponse? tenant,
    ErrorResponse? error,
    bool isBadRequest,
    bool isConflict)
  {
    Tenant = tenant;
    Error = error;
    IsBadRequest = isBadRequest;
    IsConflict = isConflict;
  }

  public TenantResponse? Tenant { get; }

  public ErrorResponse? Error { get; }

  public bool IsBadRequest { get; }

  public bool IsConflict { get; }

  public bool IsSuccess => Tenant is not null;

  public static CreateBootstrapTenantResult Success(TenantResponse tenant)
  {
    return new CreateBootstrapTenantResult(tenant, null, false, false);
  }

  public static CreateBootstrapTenantResult BadRequest(ErrorResponse error)
  {
    return new CreateBootstrapTenantResult(null, error, true, false);
  }

  public static CreateBootstrapTenantResult Conflict(ErrorResponse error)
  {
    return new CreateBootstrapTenantResult(null, error, false, true);
  }
}
