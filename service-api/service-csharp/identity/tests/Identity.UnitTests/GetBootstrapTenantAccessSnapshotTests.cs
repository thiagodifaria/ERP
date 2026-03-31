// Estes testes cobrem a leitura consolidada da estrutura de acesso do tenant.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class GetBootstrapTenantAccessSnapshotTests
{
  [Fact]
  public void ExecuteShouldReturnSnapshotForKnownTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    var userRepository = new InMemoryUserRepository();
    var teamRepository = new InMemoryTeamRepository();
    var roleRepository = new InMemoryRoleRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new GetBootstrapTenantAccessSnapshot(
      tenantRepository,
      companyRepository,
      userRepository,
      teamRepository,
      roleRepository,
      teamMembershipRepository,
      userRoleRepository);

    var response = useCase.Execute("bootstrap-ops");

    Assert.NotNull(response);
    Assert.Equal("bootstrap-ops", response!.Tenant.Slug);
    Assert.Equal(1, response.Counts.Companies);
    Assert.Equal(1, response.Counts.Users);
    Assert.Equal(1, response.Counts.Teams);
    Assert.Equal(5, response.Counts.Roles);
    Assert.Equal(1, response.Counts.TeamMemberships);
    Assert.Equal(1, response.Counts.UserRoles);
    Assert.Single(response.Users);
    Assert.Single(response.Users.First().Roles);
    Assert.Single(response.Teams);
    Assert.Single(response.Teams.First().Members);
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    var userRepository = new InMemoryUserRepository();
    var teamRepository = new InMemoryTeamRepository();
    var roleRepository = new InMemoryRoleRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);
    var userRoleRepository = new InMemoryUserRoleRepository(
      tenantRepository,
      userRepository,
      roleRepository);

    var useCase = new GetBootstrapTenantAccessSnapshot(
      tenantRepository,
      companyRepository,
      userRepository,
      teamRepository,
      roleRepository,
      teamMembershipRepository,
      userRoleRepository);

    var response = useCase.Execute("missing-tenant");

    Assert.Null(response);
  }
}
