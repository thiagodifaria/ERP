// Tenant representa o agregado raiz minimo para a fundacao multi-tenant.
// Regras de consistencia entram aqui conforme os casos de uso nascerem.
namespace Identity.Domain;

public sealed class Tenant
{
  public Tenant(long id, Guid publicId, string slug, string displayName, string status)
  {
    Id = id;
    PublicId = publicId;
    Slug = slug;
    DisplayName = displayName;
    Status = status;
  }

  public long Id { get; }

  public Guid PublicId { get; }

  public string Slug { get; }

  public string DisplayName { get; }

  public string Status { get; }
}
