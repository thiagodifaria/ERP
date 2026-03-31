// User representa o usuario interno do tenant dentro do contexto de identidade.
// Credenciais externas e seguranca aplicada entram depois sem quebrar o agregado.
namespace Identity.Domain;

public sealed class User
{
  public User(
    long id,
    long tenantId,
    long? companyId,
    Guid publicId,
    string email,
    string displayName,
    string? givenName,
    string? familyName,
    string status)
  {
    Id = id;
    TenantId = tenantId;
    CompanyId = companyId;
    PublicId = publicId;
    Email = email;
    DisplayName = displayName;
    GivenName = givenName;
    FamilyName = familyName;
    Status = status;
  }

  public long Id { get; }

  public long TenantId { get; }

  public long? CompanyId { get; }

  public Guid PublicId { get; }

  public string Email { get; }

  public string DisplayName { get; }

  public string? GivenName { get; }

  public string? FamilyName { get; }

  public string Status { get; }
}
