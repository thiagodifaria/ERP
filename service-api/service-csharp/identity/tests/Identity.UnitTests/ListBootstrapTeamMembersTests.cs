// Estes testes cobrem a leitura de memberships basicos por time no bootstrap.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class ListBootstrapTeamMembersTests
{
  [Fact]
  public void ExecuteShouldReturnMembersForKnownTeam()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var useCase = new ListBootstrapTeamMembers(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var response = useCase.Execute("bootstrap-ops", team.PublicId);

    Assert.NotNull(response);
    Assert.NotEmpty(response);
    Assert.Contains(response!, member => member.UserEmail == "owner@bootstrap-ops.local");
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownTeam()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new ListBootstrapTeamMembers(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var response = useCase.Execute("bootstrap-ops", Guid.NewGuid());

    Assert.Null(response);
  }
}
