// Tenant representa o agregado raiz minimo para a fundacao multi-tenant.
// Regras de consistencia entram aqui conforme os casos de uso nascerem.
namespace Identity.Domain;

public sealed class Tenant
{
  public Tenant(long id, Guid publicId, string name)
  {
    Id = id;
    PublicId = publicId;
    Name = name;
  }

  public long Id { get; }

  public Guid PublicId { get; }

  public string Name { get; }
}
