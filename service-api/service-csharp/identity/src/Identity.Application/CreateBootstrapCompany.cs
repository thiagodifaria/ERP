// Este caso de uso cria empresas basicas vinculadas a um tenant existente.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CreateBootstrapCompany
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ICompanyRepository _companyRepository;

  public CreateBootstrapCompany(ITenantCatalog tenantCatalog, ICompanyRepository companyRepository)
  {
    _tenantCatalog = tenantCatalog;
    _companyRepository = companyRepository;
  }

  public CreateBootstrapCompanyResult Execute(string tenantSlug, CreateCompanyRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return CreateBootstrapCompanyResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var displayName = request.DisplayName.Trim();
    var legalName = NormalizeOptional(request.LegalName);
    var taxId = NormalizeOptional(request.TaxId);

    if (string.IsNullOrWhiteSpace(displayName))
    {
      return CreateBootstrapCompanyResult.BadRequest(
        new ErrorResponse("invalid_display_name", "Display name is required."));
    }

    if (_companyRepository.FindByTenantIdAndDisplayName(tenant.Id, displayName) is not null)
    {
      return CreateBootstrapCompanyResult.Conflict(
        new ErrorResponse("company_display_name_conflict", "Company display name already exists for tenant."));
    }

    var company = new Company(
      _companyRepository.NextId(),
      tenant.Id,
      PublicIds.NewUuidV7(),
      displayName,
      legalName,
      taxId,
      "active");

    var createdCompany = _companyRepository.Add(company);

    return CreateBootstrapCompanyResult.Success(new CompanyResponse(
      createdCompany.Id,
      createdCompany.PublicId,
      createdCompany.TenantId,
      createdCompany.DisplayName,
      createdCompany.LegalName,
      createdCompany.TaxId,
      createdCompany.Status));
  }

  private static string? NormalizeOptional(string? value)
  {
    return string.IsNullOrWhiteSpace(value)
      ? null
      : value.Trim();
  }
}
