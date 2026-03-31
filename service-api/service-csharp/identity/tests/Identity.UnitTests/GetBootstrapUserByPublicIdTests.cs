// Estes testes cobrem a leitura individual de user por publicId.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class GetBootstrapUserByPublicIdTests
{
  [Fact]
  public void ExecuteShouldReturnUserForKnownPublicId()
  {
    var tenantRepository = new InMemoryTenantRepository();
    var userRepository = new InMemoryUserRepository();
    var useCase = new GetBootstrapUserByPublicId(tenantRepository, userRepository);
    var user = userRepository.ListByTenantId(1).First();

    var response = useCase.Execute("bootstrap-ops", user.PublicId);

    Assert.NotNull(response);
    Assert.Equal(user.PublicId, response!.PublicId);
    Assert.Equal(user.Email, response.Email);
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownUser()
  {
    var useCase = new GetBootstrapUserByPublicId(
      new InMemoryTenantRepository(),
      new InMemoryUserRepository());

    var response = useCase.Execute("bootstrap-ops", Guid.NewGuid());

    Assert.Null(response);
  }
}
