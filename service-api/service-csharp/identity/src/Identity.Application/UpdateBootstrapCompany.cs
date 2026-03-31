// Este caso de uso atualiza empresas basicas ja existentes dentro de um tenant conhecido.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class UpdateBootstrapCompany
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ICompanyRepository _companyRepository;

  public UpdateBootstrapCompany(ITenantCatalog tenantCatalog, ICompanyRepository companyRepository)
  {
    _tenantCatalog = tenantCatalog;
    _companyRepository = companyRepository;
  }

  public UpdateBootstrapCompanyResult Execute(string tenantSlug, Guid companyPublicId, UpdateCompanyRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return UpdateBootstrapCompanyResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var company = _companyRepository.FindByTenantIdAndPublicId(tenant.Id, companyPublicId);

    if (company is null)
    {
      return UpdateBootstrapCompanyResult.NotFound(
        new ErrorResponse("company_not_found", "Company was not found."));
    }

    var displayName = request.DisplayName.Trim();
    var legalName = NormalizeOptional(request.LegalName);
    var taxId = NormalizeOptional(request.TaxId);

    if (string.IsNullOrWhiteSpace(displayName))
    {
      return UpdateBootstrapCompanyResult.BadRequest(
        new ErrorResponse("invalid_display_name", "Display name is required."));
    }

    var existingCompany = _companyRepository.FindByTenantIdAndDisplayName(tenant.Id, displayName);
    if (existingCompany is not null && existingCompany.PublicId != company.PublicId)
    {
      return UpdateBootstrapCompanyResult.Conflict(
        new ErrorResponse("company_display_name_conflict", "Company display name already exists for tenant."));
    }

    var updatedCompany = _companyRepository.Update(company.ReviseProfile(displayName, legalName, taxId));

    return UpdateBootstrapCompanyResult.Success(new CompanyResponse(
      updatedCompany.Id,
      updatedCompany.PublicId,
      updatedCompany.TenantId,
      updatedCompany.DisplayName,
      updatedCompany.LegalName,
      updatedCompany.TaxId,
      updatedCompany.Status));
  }

  private static string? NormalizeOptional(string? value)
  {
    return string.IsNullOrWhiteSpace(value)
      ? null
      : value.Trim();
  }
}
