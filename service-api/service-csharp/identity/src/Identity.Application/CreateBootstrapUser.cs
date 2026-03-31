// Este caso de uso cria usuarios basicos vinculados a um tenant existente.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CreateBootstrapUser
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserRepository _userRepository;

  public CreateBootstrapUser(ITenantCatalog tenantCatalog, IUserRepository userRepository)
  {
    _tenantCatalog = tenantCatalog;
    _userRepository = userRepository;
  }

  public CreateBootstrapUserResult Execute(string tenantSlug, CreateUserRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return CreateBootstrapUserResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var email = NormalizeEmail(request.Email);
    var displayName = request.DisplayName.Trim();
    var givenName = NormalizeOptional(request.GivenName);
    var familyName = NormalizeOptional(request.FamilyName);

    if (!IsValidEmail(email))
    {
      return CreateBootstrapUserResult.BadRequest(
        new ErrorResponse("invalid_email", "Email is invalid."));
    }

    if (string.IsNullOrWhiteSpace(displayName))
    {
      return CreateBootstrapUserResult.BadRequest(
        new ErrorResponse("invalid_display_name", "Display name is required."));
    }

    if (_userRepository.FindByTenantIdAndEmail(tenant.Id, email) is not null)
    {
      return CreateBootstrapUserResult.Conflict(
        new ErrorResponse("user_email_conflict", "User email already exists for tenant."));
    }

    var user = new User(
      _userRepository.NextId(),
      tenant.Id,
      null,
      PublicIds.NewUuidV7(),
      email,
      displayName,
      givenName,
      familyName,
      "active");

    var createdUser = _userRepository.Add(user);

    return CreateBootstrapUserResult.Success(new UserResponse(
      createdUser.Id,
      createdUser.PublicId,
      createdUser.TenantId,
      createdUser.CompanyId,
      createdUser.Email,
      createdUser.DisplayName,
      createdUser.GivenName,
      createdUser.FamilyName,
      createdUser.Status));
  }

  private static string NormalizeEmail(string email)
  {
    return email.Trim().ToLowerInvariant();
  }

  private static string? NormalizeOptional(string? value)
  {
    return string.IsNullOrWhiteSpace(value)
      ? null
      : value.Trim();
  }

  private static bool IsValidEmail(string email)
  {
    if (string.IsNullOrWhiteSpace(email))
    {
      return false;
    }

    if (email.Contains(' '))
    {
      return false;
    }

    var separatorIndex = email.IndexOf('@');

    return separatorIndex > 0
      && separatorIndex < email.Length - 1
      && separatorIndex == email.LastIndexOf('@');
  }
}
