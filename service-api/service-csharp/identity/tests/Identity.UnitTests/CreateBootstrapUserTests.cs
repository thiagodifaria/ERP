// Estes testes cobrem a criacao minima de usuarios no bootstrap.
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class CreateBootstrapUserTests
{
  [Fact]
  public void ExecuteShouldCreateUserForKnownTenant()
  {
    var useCase = new CreateBootstrapUser(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateUserRequest("new.user@bootstrap-ops.local", "New User", "New", "User"));

    Assert.True(result.IsSuccess);
    Assert.NotNull(result.User);
    Assert.Equal("new.user@bootstrap-ops.local", result.User!.Email);
    Assert.Equal("New User", result.User.DisplayName);
  }

  [Fact]
  public void ExecuteShouldRejectUnknownTenant()
  {
    var useCase = new CreateBootstrapUser(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var result = useCase.Execute(
      "missing-tenant",
      new CreateUserRequest("new.user@missing.local", "New User", null, null));

    Assert.True(result.IsNotFound);
    Assert.NotNull(result.Error);
    Assert.Equal("tenant_not_found", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectDuplicateEmailPerTenant()
  {
    var useCase = new CreateBootstrapUser(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateUserRequest("owner@bootstrap-ops.local", "Duplicate User", null, null));

    Assert.True(result.IsConflict);
    Assert.NotNull(result.Error);
    Assert.Equal("user_email_conflict", result.Error!.Code);
  }

  [Fact]
  public void ExecuteShouldRejectInvalidEmail()
  {
    var useCase = new CreateBootstrapUser(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var result = useCase.Execute(
      "bootstrap-ops",
      new CreateUserRequest("invalid-email", "New User", null, null));

    Assert.True(result.IsBadRequest);
    Assert.NotNull(result.Error);
    Assert.Equal("invalid_email", result.Error!.Code);
  }
}
