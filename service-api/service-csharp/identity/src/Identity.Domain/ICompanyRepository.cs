// Este contrato define a escrita minima de empresas durante o bootstrap.
namespace Identity.Domain;

public interface ICompanyRepository : ICompanyCatalog
{
  Company Add(Company company);

  Company Update(Company company);

  long NextId();

  IReadOnlyCollection<Company> SeedDefaults(Tenant tenant);
}
