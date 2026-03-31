// Estes testes cobrem o update minimo de usuarios no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class UpdateBootstrapUserTests
{
  [Fact]
  public void ExecuteShouldUpdateUserForKnownTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var existingUser = userRepository.ListByTenantId(1).Single();
    var useCase = new UpdateBootstrapUser(tenantRepository, userRepository);

    var result = useCase.Execute(
      "bootstrap-ops",
      existingUser.PublicId,
      new UpdateUserRequest("owner.prime@bootstrap-ops.local", "Bootstrap Ops Prime Owner", "Bootstrap", "Prime"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.User);
    Assert.Equal("owner.prime@bootstrap-ops.local", result.User!.Email);
    Assert.Equal("Bootstrap Ops Prime Owner", result.User.DisplayName);
    Assert.Equal("Bootstrap", result.User.GivenName);
    Assert.Equal("Prime", result.User.FamilyName);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownUser()
  {
    var useCase = new UpdateBootstrapUser(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      Guid.NewGuid(),
      new UpdateUserRequest("owner.prime@bootstrap-ops.local", "Bootstrap Ops Prime Owner", null, null));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("user_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateEmailPerTenant()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    userRepository.Add(new Identity.Domain.User(
      userRepository.NextId(),
      1,
      null,
      Identity.Domain.PublicIds.NewUuidV7(),
      "operator@bootstrap-ops.local",
      "Bootstrap Operator",
      null,
      null,
      "active"));
    var existingUser = userRepository.ListByTenantId(1).Single(user => user.Email == "owner@bootstrap-ops.local");
    var useCase = new UpdateBootstrapUser(tenantRepository, userRepository);

    var result = useCase.Execute(
      "bootstrap-ops",
      existingUser.PublicId,
      new UpdateUserRequest("operator@bootstrap-ops.local", "Bootstrap Ops Prime Owner", null, null));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("user_email_conflict", result.Error!.Code);
  }
}
