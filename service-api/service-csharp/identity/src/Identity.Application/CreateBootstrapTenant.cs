// Este caso de uso cria tenants de bootstrap com validacao minima de contrato.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CreateBootstrapTenant
{
  private readonly ITenantRepository _tenantRepository;
  private readonly ICompanyRepository _companyRepository;
  private readonly IUserRepository _userRepository;
  private readonly ITeamRepository _teamRepository;
  private readonly ITeamMembershipRepository _teamMembershipRepository;
  private readonly IRoleRepository _roleRepository;

  public CreateBootstrapTenant(
    ITenantRepository tenantRepository,
    ICompanyRepository companyRepository,
    IUserRepository userRepository,
    ITeamRepository teamRepository,
    ITeamMembershipRepository teamMembershipRepository,
    IRoleRepository roleRepository)
  {
    _tenantRepository = tenantRepository;
    _companyRepository = companyRepository;
    _userRepository = userRepository;
    _teamRepository = teamRepository;
    _teamMembershipRepository = teamMembershipRepository;
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
    var defaultUsers = _userRepository.SeedDefaults(createdTenant);
    var defaultTeams = _teamRepository.SeedDefaults(createdTenant);
    _roleRepository.SeedDefaults(createdTenant);

    var defaultUser = defaultUsers.FirstOrDefault();
    var defaultTeam = defaultTeams.FirstOrDefault();

    if (defaultUser is not null
      && defaultTeam is not null
      && _teamMembershipRepository.FindByTenantIdAndTeamIdAndUserId(
        createdTenant.Id,
        defaultTeam.Id,
        defaultUser.Id) is null)
    {
      _teamMembershipRepository.Add(new TeamMembership(
        _teamMembershipRepository.NextId(),
        createdTenant.Id,
        defaultTeam.Id,
        defaultUser.Id,
        DateTimeOffset.UtcNow));
    }

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
