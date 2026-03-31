// Este caso de uso remove atribuicoes diretas de papeis por usuario dentro do tenant.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class RevokeBootstrapUserRole
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IRoleCatalog _roleCatalog;
  private readonly IUserRoleRepository _userRoleRepository;

  public RevokeBootstrapUserRole(
    ITenantCatalog tenantCatalog,
    IUserCatalog userCatalog,
    IRoleCatalog roleCatalog,
    IUserRoleRepository userRoleRepository)
  {
    _tenantCatalog = tenantCatalog;
    _userCatalog = userCatalog;
    _roleCatalog = roleCatalog;
    _userRoleRepository = userRoleRepository;
  }

  public RevokeBootstrapUserRoleResult Execute(string tenantSlug, Guid userPublicId, string roleCode)
  {
    var normalizedRoleCode = roleCode.Trim().ToLowerInvariant();

    if (string.IsNullOrWhiteSpace(normalizedRoleCode))
    {
      return RevokeBootstrapUserRoleResult.BadRequest(
        new ErrorResponse("invalid_role_code", "Role code is required."));
    }

    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return RevokeBootstrapUserRoleResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);

    if (user is null)
    {
      return RevokeBootstrapUserRoleResult.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var role = _roleCatalog.FindByTenantIdAndCode(tenant.Id, normalizedRoleCode);

    if (role is null)
    {
      return RevokeBootstrapUserRoleResult.NotFound(
        new ErrorResponse("role_not_found", "Role was not found."));
    }

    var userRole = _userRoleRepository.FindByTenantIdAndUserIdAndRoleId(tenant.Id, user.Id, role.Id);

    if (userRole is null)
    {
      return RevokeBootstrapUserRoleResult.NotFound(
        new ErrorResponse("user_role_not_found", "Role assignment was not found."));
    }

    _userRoleRepository.RemoveByTenantIdAndUserIdAndRoleId(tenant.Id, user.Id, role.Id);

    return RevokeBootstrapUserRoleResult.Success(new UserRoleResponse(
      userRole.Id,
      userRole.TenantId,
      user.PublicId,
      role.PublicId,
      role.Code,
      role.DisplayName,
      userRole.CreatedAt));
  }
}
