// Este caso de uso cria tenants de bootstrap com validacao minima de contrato.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CreateBootstrapTenant
{
  private readonly ITenantRepository _tenantRepository;
  private readonly ICompanyRepository _companyRepository;
  private readonly IUserRepository _userRepository;
  private readonly IRoleRepository _roleRepository;

  public CreateBootstrapTenant(
    ITenantRepository tenantRepository,
    ICompanyRepository companyRepository,
    IUserRepository userRepository,
    IRoleRepository roleRepository)
  {
    _tenantRepository = tenantRepository;
    _companyRepository = companyRepository;
    _userRepository = userRepository;
    _roleRepository = roleRepository;
  }

  public CreateBootstrapTenantResult Execute(CreateTenantRequest request)
  {
    var slug = NormalizeSlug(request.Slug);
    var displayName = request.DisplayName.Trim();

    if (string.IsNullOrWhiteSpace(slug))
    {
      return CreateBootstrapTenantResult.BadRequest(
        new ErrorResponse("invalid_slug", "Slug is required."));
    }

    if (string.IsNullOrWhiteSpace(displayName))
    {
      return CreateBootstrapTenantResult.BadRequest(
        new ErrorResponse("invalid_display_name", "Display name is required."));
    }

    if (!IsValidSlug(slug))
    {
      return CreateBootstrapTenantResult.BadRequest(
        new ErrorResponse("invalid_slug", "Slug must use lowercase letters, numbers or hyphens."));
    }

    if (_tenantRepository.FindBySlug(slug) is not null)
    {
      return CreateBootstrapTenantResult.Conflict(
        new ErrorResponse("tenant_slug_conflict", "Tenant slug already exists."));
    }

    var tenant = new Tenant(
      _tenantRepository.NextId(),
      PublicIds.NewUuidV7(),
      slug,
      displayName,
      "active");

    var createdTenant = _tenantRepository.Add(tenant);
    _companyRepository.SeedDefaults(createdTenant);
    _userRepository.SeedDefaults(createdTenant);
    _roleRepository.SeedDefaults(createdTenant);

    return CreateBootstrapTenantResult.Success(new TenantResponse(
      createdTenant.Id,
      createdTenant.PublicId,
      createdTenant.Slug,
      createdTenant.DisplayName,
      createdTenant.Status));
  }

  private static string NormalizeSlug(string slug)
  {
    return slug.Trim().ToLowerInvariant();
  }

  private static bool IsValidSlug(string slug)
  {
    foreach (var character in slug)
    {
      if (char.IsLower(character) || char.IsDigit(character) || character == '-')
      {
        continue;
      }

      return false;
    }

    return true;
  }
}
