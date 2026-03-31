// Este caso de uso atribui papeis existentes a usuarios dentro do tenant.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class AssignBootstrapUserRole
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserCatalog _userCatalog;
  private readonly IRoleCatalog _roleCatalog;
  private readonly IUserRoleRepository _userRoleRepository;

  public AssignBootstrapUserRole(
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

  public AssignBootstrapUserRoleResult Execute(
    string tenantSlug,
    Guid userPublicId,
    AssignUserRoleRequest request)
  {
    var roleCode = request.RoleCode.Trim().ToLowerInvariant();

    if (string.IsNullOrWhiteSpace(roleCode))
    {
      return AssignBootstrapUserRoleResult.BadRequest(
        new ErrorResponse("invalid_role_code", "Role code is required."));
    }

    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return AssignBootstrapUserRoleResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userCatalog.FindByTenantIdAndPublicId(tenant.Id, userPublicId);

    if (user is null)
    {
      return AssignBootstrapUserRoleResult.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var role = _roleCatalog.FindByTenantIdAndCode(tenant.Id, roleCode);

    if (role is null)
    {
      return AssignBootstrapUserRoleResult.NotFound(
        new ErrorResponse("role_not_found", "Role was not found."));
    }

    if (_userRoleRepository.FindByTenantIdAndUserIdAndRoleId(tenant.Id, user.Id, role.Id) is not null)
    {
      return AssignBootstrapUserRoleResult.Conflict(
        new ErrorResponse("user_role_conflict", "Role is already assigned to this user."));
    }

    var userRole = new UserRole(
      _userRoleRepository.NextId(),
      tenant.Id,
      user.Id,
      role.Id,
      DateTimeOffset.UtcNow);

    var createdUserRole = _userRoleRepository.Add(userRole);

    return AssignBootstrapUserRoleResult.Success(new UserRoleResponse(
      createdUserRole.Id,
      createdUserRole.TenantId,
      user.PublicId,
      role.PublicId,
      role.Code,
      role.DisplayName,
      createdUserRole.CreatedAt));
  }
}
