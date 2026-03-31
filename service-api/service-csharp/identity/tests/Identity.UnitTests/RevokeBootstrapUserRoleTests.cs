// Estes testes cobrem a revogacao minima de papeis a usuarios no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class RevokeBootstrapUserRoleTests
{
  [Fact]
  public void ExecuteShouldRevokeAssignedRoleForKnownUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var assignUseCase = new AssignBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var revokeUseCase = new RevokeBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.Add(new Identity.Domain.User(
      userRepository.NextId(),
      1,
      null,
      Identity.Domain.PublicIds.NewUuidV7(),
      "revoke.user@bootstrap-ops.local",
      "Revoke User",
      null,
      null,
      "active"));

    assignUseCase.Execute("bootstrap-ops", user.PublicId, new AssignUserRoleRequest("viewer"));
    var result = revokeUseCase.Execute("bootstrap-ops", user.PublicId, "viewer");

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.UserRole);
    Assert.Equal("viewer", result.UserRole!.RoleCode);
    Assert.Empty(userRoleRepository.ListByTenantIdAndUserId(1, user.Id));
  }

  [Fact]
  public void ExecuteShouldRejectBlankRoleCode()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new RevokeBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", user.PublicId, "   ");

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_role_code", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownRole()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new RevokeBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", user.PublicId, "missing-role");

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("role_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectMissingAssignment()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new RevokeBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", user.PublicId, "viewer");

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("user_role_not_found", result.Error!.Code);
  }
}
