// Este caso de uso atualiza usuarios basicos ja existentes dentro de um tenant conhecido.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class UpdateBootstrapUser
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserRepository _userRepository;

  public UpdateBootstrapUser(ITenantCatalog tenantCatalog, IUserRepository userRepository)
  {
    _tenantCatalog = tenantCatalog;
    _userRepository = userRepository;
  }

  public UpdateBootstrapUserResult Execute(string tenantSlug, Guid userPublicId, UpdateUserRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return UpdateBootstrapUserResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userRepository.FindByTenantIdAndPublicId(tenant.Id, userPublicId);

    if (user is null)
    {
      return UpdateBootstrapUserResult.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var email = NormalizeEmail(request.Email);
    var displayName = request.DisplayName.Trim();
    var givenName = NormalizeOptional(request.GivenName);
    var familyName = NormalizeOptional(request.FamilyName);

    if (!IsValidEmail(email))
    {
      return UpdateBootstrapUserResult.BadRequest(
        new ErrorResponse("invalid_email", "Email is invalid."));
    }

    if (string.IsNullOrWhiteSpace(displayName))
    {
      return UpdateBootstrapUserResult.BadRequest(
        new ErrorResponse("invalid_display_name", "Display name is required."));
    }

    var existingUser = _userRepository.FindByTenantIdAndEmail(tenant.Id, email);
    if (existingUser is not null && existingUser.PublicId != user.PublicId)
    {
      return UpdateBootstrapUserResult.Conflict(
        new ErrorResponse("user_email_conflict", "User email already exists for tenant."));
    }

    var updatedUser = _userRepository.Update(user.ReviseProfile(email, displayName, givenName, familyName));

    return UpdateBootstrapUserResult.Success(new UserResponse(
      updatedUser.Id,
      updatedUser.PublicId,
      updatedUser.TenantId,
      updatedUser.CompanyId,
      updatedUser.Email,
      updatedUser.DisplayName,
      updatedUser.GivenName,
      updatedUser.FamilyName,
      updatedUser.Status));
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
