// Este adapter fornece empresas basicas por tenant durante o bootstrap do servico.
// A persistencia real pode substituir esta implementacao sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryCompanyRepository : ICompanyRepository
{
  private readonly object _sync = new();
  private readonly Dictionary<long, List<Company>> _companiesByTenantId = [];

  public InMemoryCompanyRepository()
  {
    SeedDefaults(new Tenant(
      1,
      Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000001"),
      "bootstrap-ops",
      "Bootstrap Ops",
      "active"));

    SeedDefaults(new Tenant(
      2,
      Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000002"),
      "northwind-group",
      "Northwind Group",
      "active"));
  }

  public IReadOnlyCollection<Company> ListByTenantId(long tenantId)
  {
    lock (_sync)
    {
      if (_companiesByTenantId.TryGetValue(tenantId, out var companies))
      {
        return companies.ToArray();
      }

      return [];
    }
  }

  public Company? FindByTenantIdAndDisplayName(long tenantId, string displayName)
  {
    lock (_sync)
    {
      if (!_companiesByTenantId.TryGetValue(tenantId, out var companies))
      {
        return null;
      }

      return companies.FirstOrDefault(
        company => company.DisplayName.Equals(displayName, StringComparison.OrdinalIgnoreCase));
    }
  }

  public Company? FindByTenantIdAndPublicId(long tenantId, Guid publicId)
  {
    lock (_sync)
    {
      if (!_companiesByTenantId.TryGetValue(tenantId, out var companies))
      {
        return null;
      }

      return companies.FirstOrDefault(company => company.PublicId == publicId);
    }
  }

  public Company Add(Company company)
  {
    lock (_sync)
    {
      if (!_companiesByTenantId.TryGetValue(company.TenantId, out var companies))
      {
        companies = [];
        _companiesByTenantId[company.TenantId] = companies;
      }

      companies.Add(company);

      return company;
    }
  }

  public Company Update(Company company)
  {
    lock (_sync)
    {
      if (!_companiesByTenantId.TryGetValue(company.TenantId, out var companies))
      {
        companies = [];
        _companiesByTenantId[company.TenantId] = companies;
      }

      for (var index = 0; index < companies.Count; index++)
      {
        if (companies[index].PublicId != company.PublicId)
        {
          continue;
        }

        companies[index] = company;
        return company;
      }

      companies.Add(company);
      return company;
    }
  }

  public long NextId()
  {
    lock (_sync)
    {
      return NextIdUnsafe();
    }
  }

  public IReadOnlyCollection<Company> SeedDefaults(Tenant tenant)
  {
    lock (_sync)
    {
      if (_companiesByTenantId.TryGetValue(tenant.Id, out var existingCompanies) && existingCompanies.Count > 0)
      {
        return existingCompanies.ToArray();
      }

      var seededCompanies = DefaultCompanies(tenant)
        .Select(companySeed => new Company(
          NextIdUnsafe(),
          tenant.Id,
          PublicIds.NewUuidV7(),
          companySeed.DisplayName,
          companySeed.LegalName,
          companySeed.TaxId,
          "active"))
        .ToList();

      _companiesByTenantId[tenant.Id] = seededCompanies;

      return seededCompanies.ToArray();
    }
  }

  private long NextIdUnsafe()
  {
    return _companiesByTenantId
      .SelectMany(entry => entry.Value)
      .Select(company => company.Id)
      .DefaultIfEmpty(0)
      .Max() + 1;
  }

  private static IReadOnlyCollection<(string DisplayName, string? LegalName, string? TaxId)> DefaultCompanies(Tenant tenant)
  {
    return
    [
      (tenant.DisplayName, tenant.DisplayName, null)
    ];
  }
}
