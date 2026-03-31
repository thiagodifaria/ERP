// Estes testes cobrem a remocao minima de memberships de time no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class RemoveBootstrapTeamMemberTests
{
  [Fact]
  public void ExecuteShouldRemoveExistingMembership()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new RemoveBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      teamMembershipRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, user.PublicId);

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Membership);
    Assert.Equal(user.PublicId, result.Membership!.UserPublicId);
    Assert.Empty(teamMembershipRepository.ListByTenantIdAndTeamId(1, team.Id));
  }

  [Fact]
  public void ExecuteShouldRejectEmptyUserPublicId()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new RemoveBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      teamMembershipRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, Guid.Empty);

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_user_public_id", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new RemoveBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      teamMembershipRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, Guid.NewGuid());

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("user_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectMissingMembership()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var teamMembershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new RemoveBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      teamMembershipRepository);

    var team = teamRepository.Add(new Identity.Domain.Team(
      teamRepository.NextId(),
      1,
      null,
      Identity.Domain.PublicIds.NewUuidV7(),
      "Field Ops",
      "active"));
    var user = userRepository.ListByTenantId(1).First();

    var result = useCase.Execute("bootstrap-ops", team.PublicId, user.PublicId);

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("team_membership_not_found", result.Error!.Code);
  }
}
