// Estes testes cobrem a atualizacao minima de times no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class UpdateBootstrapTeamTests
{
  [Fact]
  public void ExecuteShouldUpdateTeamForKnownTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var useCase = new UpdateBootstrapTeam(tenantRepository, teamRepository);
    var team = teamRepository.ListByTenantId(1).First();

    var result = useCase.Execute(
      "bootstrap-ops",
      team.PublicId,
      new UpdateTeamRequest("Core Prime"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Team);
    Assert.Equal("Core Prime", result.Team!.Name);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownTeam()
  {
    var useCase = new UpdateBootstrapTeam(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      Guid.NewGuid(),
      new UpdateTeamRequest("Core Prime"));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("team_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectBlankName()
  {
    var teamRepository = new InMemoryTeamRepository();
    var useCase = new UpdateBootstrapTeam(new InMemoryTenantRepository(), teamRepository);
    var team = teamRepository.ListByTenantId(1).First();

    var result = useCase.Execute(
      "bootstrap-ops",
      team.PublicId,
      new UpdateTeamRequest("   "));

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_team_name", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateName()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var teamRepository = new InMemoryTeamRepository();
    var useCase = new UpdateBootstrapTeam(tenantRepository, teamRepository);

    var existingTeam = teamRepository.ListByTenantId(1).First();
    var secondTeam = teamRepository.Add(new Identity.Domain.Team(
      teamRepository.NextId(),
      1,
      null,
      Identity.Domain.PublicIds.NewUuidV7(),
      "Field Ops",
      "active"));

    var result = useCase.Execute(
      "bootstrap-ops",
      secondTeam.PublicId,
      new UpdateTeamRequest(existingTeam.Name));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("team_name_conflict", result.Error!.Code);
  }
}
