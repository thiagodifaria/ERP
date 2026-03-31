// Estes testes cobrem a atribuicao minima de papeis a usuarios no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class AssignBootstrapUserRoleTests
{
  [Fact]
  public void ExecuteShouldAssignRoleForKnownUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new AssignBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.Add(new Identity.Domain.User(
      userRepository.NextId(),
      1,
      null,
      Identity.Domain.PublicIds.NewUuidV7(),
      "role.user@bootstrap-ops.local",
      "Role User",
      null,
      null,
      "active"));

    var result = useCase.Execute("bootstrap-ops", user.PublicId, new AssignUserRoleRequest("viewer"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.UserRole);
    Assert.Equal("viewer", result.UserRole!.RoleCode);
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

    var useCase = new AssignBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", user.PublicId, new AssignUserRoleRequest("   "));

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

    var useCase = new AssignBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", user.PublicId, new AssignUserRoleRequest("missing-role"));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("role_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateAssignment()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new AssignBootstrapUserRole(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", user.PublicId, new AssignUserRoleRequest("owner"));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("user_role_conflict", result.Error!.Code);
  }
}
