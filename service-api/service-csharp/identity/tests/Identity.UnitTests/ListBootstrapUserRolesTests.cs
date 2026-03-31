// Estes testes cobrem a leitura de papeis atribuidos a usuarios no bootstrap.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapUserRolesTests
{
  [Fact]
  public void ExecuteShouldReturnRolesForKnownUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var user = userRepository.ListByTenantId(1).First();
    var useCase = new ListBootstrapUserRoles(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var response = useCase.Execute("bootstrap-ops", user.PublicId);

    Assert.NotNull(response);
    Assert.NotEmpty(response);
    Assert.Contains(response!, userRole => userRole.RoleCode == "owner");
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new ListBootstrapUserRoles(
      tenantRepository,
      userRepository,
      roleRepository,
      userRoleRepository);

    var response = useCase.Execute("bootstrap-ops", Guid.NewGuid());

    Assert.Null(response);
  }
}
