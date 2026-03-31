// Estes testes cobrem a leitura individual de tenant por slug.
using Identity.Application;
using Identity.Infrastructure;
using Xunit;

namespace Identity.UnitTests;

public sealed class GetBootstrapTenantBySlugTests
{
  [Fact]
  public void ExecuteShouldReturnTenantForKnownSlug()
  {
    var useCase = new GetBootstrapTenantBySlug(new InMemoryTenantRepository());

    var response = useCase.Execute("bootstrap-ops");

    Assert.NotNull(response);
    Assert.Equal("bootstrap-ops", response!.Slug);
  }

  [Fact]
  public void ExecuteShouldReturnNullForUnknownSlug()
  {
    var useCase = new GetBootstrapTenantBySlug(new InMemoryTenantRepository());

    var response = useCase.Execute("missing-tenant");

    Assert.Null(response);
  }
}
