// Estes testes validam o contrato minimo do agregado Tenant.
using Identity.Domain;
using Xunit;

namespace Identity.UnitTests;

public sealed class TenantTests
{
  [Fact]
  public void ConstructorShouldKeepIdentityFields()
  {
    var publicId = Guid.NewGuid();
    var tenant = new Tenant(10, publicId, "acme", "Acme", "active");

    Assert.Equal(10, tenant.Id);
    Assert.Equal(publicId, tenant.PublicId);
    Assert.Equal("acme", tenant.Slug);
    Assert.Equal("Acme", tenant.DisplayName);
    Assert.Equal("active", tenant.Status);
  }
}
