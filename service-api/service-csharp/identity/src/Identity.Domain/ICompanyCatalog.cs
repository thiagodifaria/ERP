// Este contrato define a leitura minima de empresas por tenant.
// Consultas de bootstrap devem depender desta abstracao.
namespace Identity.Domain;

public interface ICompanyCatalog
{
  IReadOnlyCollection<Company> ListByTenantId(long tenantId);

  Company? FindByTenantIdAndDisplayName(long tenantId, string displayName);
}
