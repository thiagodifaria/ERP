// Estes testes cobrem a criacao minima de memberships de time no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class AddBootstrapTeamMemberTests
{
  [Fact]
  public void ExecuteShouldCreateMembershipForKnownTeamAndUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new AddBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var user = userRepository.Add(new Identity.Domain.User(
      userRepository.NextId(),
      1,
      null,
      Identity.Domain.PublicIds.NewUuidV7(),
      "new.member@bootstrap-ops.local",
      "New Member",
      null,
      null,
      "active"));

    var team = teamRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, new AddTeamMemberRequest(user.PublicId));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Membership);
    Assert.Equal(user.PublicId, result.Membership!.UserPublicId);
  }

  [Fact]
  public void ExecuteShouldRejectMissingUserPublicId()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new AddBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, new AddTeamMemberRequest(Guid.Empty));

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_user_public_id", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownTeam()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new AddBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", Guid.NewGuid(), new AddTeamMemberRequest(user.PublicId));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("team_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownUser()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new AddBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, new AddTeamMemberRequest(Guid.NewGuid()));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("user_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateMembership()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var userRepository = new InMemoryUserRepository();
    var membershipRepository = new InMemoryTeamMembershipRepository(
      tenantRepository,
      teamRepository,
      userRepository);

    var useCase = new AddBootstrapTeamMember(
      tenantRepository,
      teamRepository,
      userRepository,
      membershipRepository);

    var team = teamRepository.ListByTenantId(1).First();
    var user = userRepository.ListByTenantId(1).First();
    var result = useCase.Execute("bootstrap-ops", team.PublicId, new AddTeamMemberRequest(user.PublicId));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("team_membership_conflict", result.Error!.Code);
  }
}
