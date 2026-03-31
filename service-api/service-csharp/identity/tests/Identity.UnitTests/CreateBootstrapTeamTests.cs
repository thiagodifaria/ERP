// Estes testes cobrem a criacao minima de times no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class CreateBootstrapTeamTests
{
  [Fact]
  public void ExecuteShouldCreateTeamForKnownTenant()
  {
    var useCase = new CreateBootstrapTeam(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateTeamRequest("Field Ops"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.Team);
    Assert.Equal("Field Ops", result.Team!.Name);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownTenant()
  {
    var useCase = new CreateBootstrapTeam(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var result = useCase.Execute(
      "missing-tenant",
      new CreateTeamRequest("Field Ops"));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("tenant_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateNamePerTenant()
  {
    var useCase = new CreateBootstrapTeam(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateTeamRequest("Core"));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("team_name_conflict", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectBlankName()
  {
    var useCase = new CreateBootstrapTeam(
      new InMemoryTenantRepository(),
      new InMemoryTeamRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateTeamRequest("   "));

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_team_name", result.Error!.Code);
  }
}
