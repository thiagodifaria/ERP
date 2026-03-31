// User representa o usuario interno do tenant dentro do contexto de identidade.
// Credenciais externas e seguranca aplicada entram depois sem quebrar o agregado.
namespace Identity.Domain;

public sealed class User
{
  public User(long id, long tenantId, Guid publicId, string email, string displayName, string status)
  {
    Id = id;
    TenantId = tenantId;
    PublicId = publicId;
    Email = email;
    DisplayName = displayName;
    Status = status;
  }

  public long Id { get; }

  public long TenantId { get; }

  public Guid PublicId { get; }

  public string Email { get; }

  public string DisplayName { get; }

  public string Status { get; }
}
