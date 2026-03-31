// Estes testes cobrem a criacao minima de tenants no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class CreateBootstrapTenantTests
{
  [Fact]
  public void ExecuteShouldCreateTenantForValidPayload()
  {
    var repository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    var userRepository = new InMemoryUserRepository();
    var teamRepository = new InMemoryTeamRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      repository,
      teamRepository,
      userRepository);
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      repository,
      userRepository,
      roleRepository);
    var useCase = new CreateBootstrapTenant(
      repository,
      companyRepository,
      userRepository,
      teamRepository,
      teamMembershipRepository,
      roleRepository,
      userRoleRepository);

    var result = useCase.Execute(new CreateTenantRequest("tenant-lab", "Tenant Lab"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Tenant);
    Assert.Equal("tenant-lab", result.Tenant!.Slug);
    Assert.Equal("Tenant Lab", result.Tenant.DisplayName);
    Assert.Single(companyRepository.ListByTenantId(result.Tenant.Id));
    Assert.Single(userRepository.ListByTenantId(result.Tenant.Id));
    Assert.Single(teamRepository.ListByTenantId(result.Tenant.Id));
    Assert.Single(teamMembershipRepository.ListByTenantIdAndTeamId(
      result.Tenant.Id,
      teamRepository.ListByTenantId(result.Tenant.Id).First().Id));
    Assert.Single(userRoleRepository.ListByTenantIdAndUserId(
      result.Tenant.Id,
      userRepository.ListByTenantId(result.Tenant.Id).First().Id));
    Assert.Equal(5, roleRepository.ListByTenantSlug("tenant-lab").Count);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateSlug()
  {
    var repository = new InMemoryTenantRepository();
    var companyRepository = new InMemoryCompanyRepository();
    var userRepository = new InMemoryUserRepository();
    var teamRepository = new InMemoryTeamRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      repository,
      teamRepository,
      userRepository);
    var roleRepository = new InMemoryRoleRepository();
    var userRoleRepository = new InMemoryUserRoleRepository(
      repository,
      userRepository,
      roleRepository);
    var useCase = new CreateBootstrapTenant(
      repository,
      companyRepository,
      userRepository,
      teamRepository,
      teamMembershipRepository,
      roleRepository,
      userRoleRepository);

    var result = useCase.Execute(new CreateTenantRequest("bootstrap-ops", "Another Name"));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("tenant_slug_conflict", result.Error!.Code);
  }
}
