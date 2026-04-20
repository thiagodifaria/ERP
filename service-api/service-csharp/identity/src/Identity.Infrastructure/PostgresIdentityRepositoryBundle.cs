// Este adapter conecta o bootstrap de identidade ao PostgreSQL.
// O objetivo e manter os contratos de dominio estaveis enquanto o runtime deixa a memoria.
using Identity.Domain;
using Npgsql;

namespace Identity.Infrastructure;

public sealed class PostgresIdentityRepositoryBundle :
  ITenantRepository,
  ICompanyRepository,
  IUserRepository,
  ITeamRepository,
  IRoleRepository,
  ITeamMembershipRepository,
  IUserRoleRepository
{
  private static readonly IReadOnlyCollection<(string Code, string DisplayName)> DefaultRoles =
  [
    ("owner", "Owner"),
    ("admin", "Administrator"),
    ("manager", "Manager"),
    ("operator", "Operator"),
    ("viewer", "Viewer")
  ];

  private readonly string _connectionString;

  public PostgresIdentityRepositoryBundle(string connectionString)
  {
    _connectionString = connectionString;
  }

  IReadOnlyCollection<Tenant> ITenantCatalog.List()
  {
    const string sql = """
      SELECT id, public_id, slug, display_name, status
      FROM identity.tenants
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    using var reader = command.ExecuteReader();

    var tenants = new List<Tenant>();
    while (reader.Read())
    {
      tenants.Add(MapTenant(reader));
    }

    return tenants;
  }

  Tenant? ITenantCatalog.FindBySlug(string slug)
  {
    const string sql = """
      SELECT id, public_id, slug, display_name, status
      FROM identity.tenants
      WHERE lower(slug) = lower(@slug)
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("slug", slug);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapTenant(reader)
      : null;
  }

  Tenant ITenantRepository.Add(Tenant tenant)
  {
    const string sql = """
      INSERT INTO identity.tenants (public_id, slug, display_name, status)
      VALUES (@public_id, @slug, @display_name, @status)
      RETURNING id, public_id, slug, display_name, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("public_id", tenant.PublicId);
    command.Parameters.AddWithValue("slug", tenant.Slug);
    command.Parameters.AddWithValue("display_name", tenant.DisplayName);
    command.Parameters.AddWithValue("status", tenant.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapTenant(reader);
  }

  long ITenantRepository.NextId()
  {
    return NextId("identity.tenants");
  }

  IReadOnlyCollection<Company> ICompanyCatalog.ListByTenantId(long tenantId)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, display_name, legal_name, tax_id, status
      FROM identity.companies
      WHERE tenant_id = @tenant_id
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    using var reader = command.ExecuteReader();

    var companies = new List<Company>();
    while (reader.Read())
    {
      companies.Add(MapCompany(reader));
    }

    return companies;
  }

  Company? ICompanyCatalog.FindByTenantIdAndDisplayName(long tenantId, string displayName)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, display_name, legal_name, tax_id, status
      FROM identity.companies
      WHERE tenant_id = @tenant_id
        AND lower(display_name) = lower(@display_name)
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("display_name", displayName);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapCompany(reader)
      : null;
  }

  Company? ICompanyCatalog.FindByTenantIdAndPublicId(long tenantId, Guid publicId)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, display_name, legal_name, tax_id, status
      FROM identity.companies
      WHERE tenant_id = @tenant_id
        AND public_id = @public_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("public_id", publicId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapCompany(reader)
      : null;
  }

  Company ICompanyRepository.Add(Company company)
  {
    const string sql = """
      INSERT INTO identity.companies (tenant_id, public_id, display_name, legal_name, tax_id, status)
      VALUES (@tenant_id, @public_id, @display_name, @legal_name, @tax_id, @status)
      RETURNING id, tenant_id, public_id, display_name, legal_name, tax_id, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", company.TenantId);
    command.Parameters.AddWithValue("public_id", company.PublicId);
    command.Parameters.AddWithValue("display_name", company.DisplayName);
    command.Parameters.AddWithValue("legal_name", (object?)company.LegalName ?? DBNull.Value);
    command.Parameters.AddWithValue("tax_id", (object?)company.TaxId ?? DBNull.Value);
    command.Parameters.AddWithValue("status", company.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapCompany(reader);
  }

  Company ICompanyRepository.Update(Company company)
  {
    const string sql = """
      UPDATE identity.companies
      SET display_name = @display_name,
          legal_name = @legal_name,
          tax_id = @tax_id,
          status = @status
      WHERE tenant_id = @tenant_id
        AND public_id = @public_id
      RETURNING id, tenant_id, public_id, display_name, legal_name, tax_id, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", company.TenantId);
    command.Parameters.AddWithValue("public_id", company.PublicId);
    command.Parameters.AddWithValue("display_name", company.DisplayName);
    command.Parameters.AddWithValue("legal_name", (object?)company.LegalName ?? DBNull.Value);
    command.Parameters.AddWithValue("tax_id", (object?)company.TaxId ?? DBNull.Value);
    command.Parameters.AddWithValue("status", company.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapCompany(reader);
  }

  long ICompanyRepository.NextId()
  {
    return NextId("identity.companies");
  }

  IReadOnlyCollection<Company> ICompanyRepository.SeedDefaults(Tenant tenant)
  {
    var companyCatalog = (ICompanyCatalog)this;
    var companyRepository = (ICompanyRepository)this;

    if (companyCatalog.FindByTenantIdAndDisplayName(tenant.Id, tenant.DisplayName) is null)
    {
      companyRepository.Add(new Company(
        0,
        tenant.Id,
        PublicIds.NewUuidV7(),
        tenant.DisplayName,
        tenant.DisplayName,
        null,
        "active"));
    }

    return companyCatalog.ListByTenantId(tenant.Id);
  }

  IReadOnlyCollection<User> IUserCatalog.ListByTenantId(long tenantId)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, email::text, display_name, given_name, family_name, status
      FROM identity.users
      WHERE tenant_id = @tenant_id
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    using var reader = command.ExecuteReader();

    var users = new List<User>();
    while (reader.Read())
    {
      users.Add(MapUser(reader));
    }

    return users;
  }

  User? IUserCatalog.FindByTenantIdAndId(long tenantId, long userId)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, email::text, display_name, given_name, family_name, status
      FROM identity.users
      WHERE tenant_id = @tenant_id
        AND id = @id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("id", userId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapUser(reader)
      : null;
  }

  User? IUserCatalog.FindByTenantIdAndPublicId(long tenantId, Guid publicId)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, email::text, display_name, given_name, family_name, status
      FROM identity.users
      WHERE tenant_id = @tenant_id
        AND public_id = @public_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("public_id", publicId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapUser(reader)
      : null;
  }

  User? IUserCatalog.FindByTenantIdAndEmail(long tenantId, string email)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, email::text, display_name, given_name, family_name, status
      FROM identity.users
      WHERE tenant_id = @tenant_id
        AND email = @email
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("email", email);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapUser(reader)
      : null;
  }

  User IUserRepository.Add(User user)
  {
    const string sql = """
      INSERT INTO identity.users (tenant_id, company_id, public_id, email, display_name, given_name, family_name, status)
      VALUES (@tenant_id, @company_id, @public_id, @email, @display_name, @given_name, @family_name, @status)
      RETURNING id, tenant_id, company_id, public_id, email::text, display_name, given_name, family_name, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", user.TenantId);
    command.Parameters.AddWithValue("company_id", (object?)user.CompanyId ?? DBNull.Value);
    command.Parameters.AddWithValue("public_id", user.PublicId);
    command.Parameters.AddWithValue("email", user.Email);
    command.Parameters.AddWithValue("display_name", user.DisplayName);
    command.Parameters.AddWithValue("given_name", (object?)user.GivenName ?? DBNull.Value);
    command.Parameters.AddWithValue("family_name", (object?)user.FamilyName ?? DBNull.Value);
    command.Parameters.AddWithValue("status", user.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapUser(reader);
  }

  User IUserRepository.Update(User user)
  {
    const string sql = """
      UPDATE identity.users
      SET email = @email,
          display_name = @display_name,
          given_name = @given_name,
          family_name = @family_name,
          status = @status
      WHERE tenant_id = @tenant_id
        AND public_id = @public_id
      RETURNING id, tenant_id, company_id, public_id, email::text, display_name, given_name, family_name, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", user.TenantId);
    command.Parameters.AddWithValue("public_id", user.PublicId);
    command.Parameters.AddWithValue("email", user.Email);
    command.Parameters.AddWithValue("display_name", user.DisplayName);
    command.Parameters.AddWithValue("given_name", (object?)user.GivenName ?? DBNull.Value);
    command.Parameters.AddWithValue("family_name", (object?)user.FamilyName ?? DBNull.Value);
    command.Parameters.AddWithValue("status", user.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapUser(reader);
  }

  long IUserRepository.NextId()
  {
    return NextId("identity.users");
  }

  IReadOnlyCollection<User> IUserRepository.SeedDefaults(Tenant tenant)
  {
    var userCatalog = (IUserCatalog)this;
    var userRepository = (IUserRepository)this;
    var companyCatalog = (ICompanyCatalog)this;
    var ownerEmail = $"owner@{tenant.Slug}.local";

    if (userCatalog.FindByTenantIdAndEmail(tenant.Id, ownerEmail) is null)
    {
      var companyId = companyCatalog
        .ListByTenantId(tenant.Id)
        .Select(company => (long?)company.Id)
        .FirstOrDefault();

      userRepository.Add(new User(
        0,
        tenant.Id,
        companyId,
        PublicIds.NewUuidV7(),
        ownerEmail,
        $"{tenant.DisplayName} Owner",
        tenant.DisplayName,
        "Owner",
        "active"));
    }

    return userCatalog.ListByTenantId(tenant.Id);
  }

  IReadOnlyCollection<Team> ITeamCatalog.ListByTenantId(long tenantId)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, name, status
      FROM identity.teams
      WHERE tenant_id = @tenant_id
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    using var reader = command.ExecuteReader();

    var teams = new List<Team>();
    while (reader.Read())
    {
      teams.Add(MapTeam(reader));
    }

    return teams;
  }

  Team? ITeamCatalog.FindByTenantIdAndName(long tenantId, string name)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, name, status
      FROM identity.teams
      WHERE tenant_id = @tenant_id
        AND lower(name) = lower(@name)
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("name", name);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapTeam(reader)
      : null;
  }

  Team? ITeamCatalog.FindByTenantIdAndPublicId(long tenantId, Guid publicId)
  {
    const string sql = """
      SELECT id, tenant_id, company_id, public_id, name, status
      FROM identity.teams
      WHERE tenant_id = @tenant_id
        AND public_id = @public_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("public_id", publicId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapTeam(reader)
      : null;
  }

  Team ITeamRepository.Add(Team team)
  {
    const string sql = """
      INSERT INTO identity.teams (tenant_id, company_id, public_id, name, status)
      VALUES (@tenant_id, @company_id, @public_id, @name, @status)
      RETURNING id, tenant_id, company_id, public_id, name, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", team.TenantId);
    command.Parameters.AddWithValue("company_id", (object?)team.CompanyId ?? DBNull.Value);
    command.Parameters.AddWithValue("public_id", team.PublicId);
    command.Parameters.AddWithValue("name", team.Name);
    command.Parameters.AddWithValue("status", team.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapTeam(reader);
  }

  Team ITeamRepository.Update(Team team)
  {
    const string sql = """
      UPDATE identity.teams
      SET name = @name,
          status = @status
      WHERE tenant_id = @tenant_id
        AND public_id = @public_id
      RETURNING id, tenant_id, company_id, public_id, name, status;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", team.TenantId);
    command.Parameters.AddWithValue("public_id", team.PublicId);
    command.Parameters.AddWithValue("name", team.Name);
    command.Parameters.AddWithValue("status", team.Status);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapTeam(reader);
  }

  long ITeamRepository.NextId()
  {
    return NextId("identity.teams");
  }

  IReadOnlyCollection<Team> ITeamRepository.SeedDefaults(Tenant tenant)
  {
    var teamCatalog = (ITeamCatalog)this;
    var teamRepository = (ITeamRepository)this;
    var companyCatalog = (ICompanyCatalog)this;

    if (teamCatalog.FindByTenantIdAndName(tenant.Id, "Core") is null)
    {
      var companyId = companyCatalog
        .ListByTenantId(tenant.Id)
        .Select(company => (long?)company.Id)
        .FirstOrDefault();

      teamRepository.Add(new Team(
        0,
        tenant.Id,
        companyId,
        PublicIds.NewUuidV7(),
        "Core",
        "active"));
    }

    return teamCatalog.ListByTenantId(tenant.Id);
  }

  IReadOnlyCollection<Role> IRoleCatalog.ListByTenantId(long tenantId)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, code, display_name, status
      FROM identity.roles
      WHERE tenant_id = @tenant_id
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    using var reader = command.ExecuteReader();

    var roles = new List<Role>();
    while (reader.Read())
    {
      roles.Add(MapRole(reader));
    }

    return roles;
  }

  IReadOnlyCollection<Role> IRoleCatalog.ListByTenantSlug(string tenantSlug)
  {
    const string sql = """
      SELECT role.id, role.tenant_id, role.public_id, role.code, role.display_name, role.status
      FROM identity.roles AS role
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = role.tenant_id
      WHERE lower(tenant.slug) = lower(@tenant_slug)
      ORDER BY role.id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_slug", tenantSlug);
    using var reader = command.ExecuteReader();

    var roles = new List<Role>();
    while (reader.Read())
    {
      roles.Add(MapRole(reader));
    }

    return roles;
  }

  Role? IRoleCatalog.FindByTenantIdAndId(long tenantId, long roleId)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, code, display_name, status
      FROM identity.roles
      WHERE tenant_id = @tenant_id
        AND id = @id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("id", roleId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapRole(reader)
      : null;
  }

  Role? IRoleCatalog.FindByTenantIdAndCode(long tenantId, string code)
  {
    const string sql = """
      SELECT id, tenant_id, public_id, code, display_name, status
      FROM identity.roles
      WHERE tenant_id = @tenant_id
        AND lower(code) = lower(@code)
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("code", code);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapRole(reader)
      : null;
  }

  IReadOnlyCollection<Role> IRoleRepository.SeedDefaults(Tenant tenant)
  {
    var roleCatalog = (IRoleCatalog)this;

    foreach (var roleSeed in DefaultRoles)
    {
      if (roleCatalog.FindByTenantIdAndCode(tenant.Id, roleSeed.Code) is not null)
      {
        continue;
      }

      const string sql = """
        INSERT INTO identity.roles (tenant_id, public_id, code, display_name, status)
        VALUES (@tenant_id, @public_id, @code, @display_name, @status);
        """;

      using var connection = OpenConnection();
      using var command = new NpgsqlCommand(sql, connection);
      command.Parameters.AddWithValue("tenant_id", tenant.Id);
      command.Parameters.AddWithValue("public_id", PublicIds.NewUuidV7());
      command.Parameters.AddWithValue("code", roleSeed.Code);
      command.Parameters.AddWithValue("display_name", roleSeed.DisplayName);
      command.Parameters.AddWithValue("status", "active");
      command.ExecuteNonQuery();
    }

    return roleCatalog.ListByTenantId(tenant.Id);
  }

  IReadOnlyCollection<TeamMembership> ITeamMembershipCatalog.ListByTenantIdAndTeamId(long tenantId, long teamId)
  {
    const string sql = """
      SELECT id, tenant_id, team_id, user_id, created_at
      FROM identity.team_memberships
      WHERE tenant_id = @tenant_id
        AND team_id = @team_id
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("team_id", teamId);
    using var reader = command.ExecuteReader();

    var memberships = new List<TeamMembership>();
    while (reader.Read())
    {
      memberships.Add(MapTeamMembership(reader));
    }

    return memberships;
  }

  TeamMembership? ITeamMembershipCatalog.FindByTenantIdAndTeamIdAndUserId(long tenantId, long teamId, long userId)
  {
    const string sql = """
      SELECT id, tenant_id, team_id, user_id, created_at
      FROM identity.team_memberships
      WHERE tenant_id = @tenant_id
        AND team_id = @team_id
        AND user_id = @user_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("team_id", teamId);
    command.Parameters.AddWithValue("user_id", userId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapTeamMembership(reader)
      : null;
  }

  TeamMembership ITeamMembershipRepository.Add(TeamMembership membership)
  {
    const string sql = """
      INSERT INTO identity.team_memberships (tenant_id, team_id, user_id, created_at)
      VALUES (@tenant_id, @team_id, @user_id, @created_at)
      RETURNING id, tenant_id, team_id, user_id, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", membership.TenantId);
    command.Parameters.AddWithValue("team_id", membership.TeamId);
    command.Parameters.AddWithValue("user_id", membership.UserId);
    command.Parameters.AddWithValue("created_at", membership.CreatedAt);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapTeamMembership(reader);
  }

  bool ITeamMembershipRepository.RemoveByTenantIdAndTeamIdAndUserId(long tenantId, long teamId, long userId)
  {
    const string sql = """
      DELETE FROM identity.team_memberships
      WHERE tenant_id = @tenant_id
        AND team_id = @team_id
        AND user_id = @user_id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("team_id", teamId);
    command.Parameters.AddWithValue("user_id", userId);

    return command.ExecuteNonQuery() > 0;
  }

  long ITeamMembershipRepository.NextId()
  {
    return NextId("identity.team_memberships");
  }

  IReadOnlyCollection<UserRole> IUserRoleCatalog.ListByTenantIdAndUserId(long tenantId, long userId)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, role_id, created_at
      FROM identity.user_roles
      WHERE tenant_id = @tenant_id
        AND user_id = @user_id
      ORDER BY id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("user_id", userId);
    using var reader = command.ExecuteReader();

    var userRoles = new List<UserRole>();
    while (reader.Read())
    {
      userRoles.Add(MapUserRole(reader));
    }

    return userRoles;
  }

  UserRole? IUserRoleCatalog.FindByTenantIdAndUserIdAndRoleId(long tenantId, long userId, long roleId)
  {
    const string sql = """
      SELECT id, tenant_id, user_id, role_id, created_at
      FROM identity.user_roles
      WHERE tenant_id = @tenant_id
        AND user_id = @user_id
        AND role_id = @role_id
      LIMIT 1;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("user_id", userId);
    command.Parameters.AddWithValue("role_id", roleId);
    using var reader = command.ExecuteReader();

    return reader.Read()
      ? MapUserRole(reader)
      : null;
  }

  UserRole IUserRoleRepository.Add(UserRole userRole)
  {
    const string sql = """
      INSERT INTO identity.user_roles (tenant_id, user_id, role_id, created_at)
      VALUES (@tenant_id, @user_id, @role_id, @created_at)
      RETURNING id, tenant_id, user_id, role_id, created_at;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", userRole.TenantId);
    command.Parameters.AddWithValue("user_id", userRole.UserId);
    command.Parameters.AddWithValue("role_id", userRole.RoleId);
    command.Parameters.AddWithValue("created_at", userRole.CreatedAt);
    using var reader = command.ExecuteReader();
    reader.Read();

    return MapUserRole(reader);
  }

  bool IUserRoleRepository.RemoveByTenantIdAndUserIdAndRoleId(long tenantId, long userId, long roleId)
  {
    const string sql = """
      DELETE FROM identity.user_roles
      WHERE tenant_id = @tenant_id
        AND user_id = @user_id
        AND role_id = @role_id;
      """;

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue("tenant_id", tenantId);
    command.Parameters.AddWithValue("user_id", userId);
    command.Parameters.AddWithValue("role_id", roleId);

    return command.ExecuteNonQuery() > 0;
  }

  long IUserRoleRepository.NextId()
  {
    return NextId("identity.user_roles");
  }

  private long NextId(string tableName)
  {
    var sql = $"SELECT COALESCE(MAX(id), 0) + 1 FROM {tableName};";

    using var connection = OpenConnection();
    using var command = new NpgsqlCommand(sql, connection);

    return Convert.ToInt64(command.ExecuteScalar()!);
  }

  private NpgsqlConnection OpenConnection()
  {
    var connection = new NpgsqlConnection(_connectionString);
    connection.Open();

    return connection;
  }

  private static Tenant MapTenant(NpgsqlDataReader reader)
  {
    return new Tenant(
      reader.GetInt64(0),
      reader.GetGuid(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4));
  }

  private static Company MapCompany(NpgsqlDataReader reader)
  {
    return new Company(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetGuid(2),
      reader.GetString(3),
      reader.IsDBNull(4) ? null : reader.GetString(4),
      reader.IsDBNull(5) ? null : reader.GetString(5),
      reader.GetString(6));
  }

  private static User MapUser(NpgsqlDataReader reader)
  {
    return new User(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.IsDBNull(2) ? null : reader.GetInt64(2),
      reader.GetGuid(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.IsDBNull(6) ? null : reader.GetString(6),
      reader.IsDBNull(7) ? null : reader.GetString(7),
      reader.GetString(8));
  }

  private static Team MapTeam(NpgsqlDataReader reader)
  {
    return new Team(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.IsDBNull(2) ? null : reader.GetInt64(2),
      reader.GetGuid(3),
      reader.GetString(4),
      reader.GetString(5));
  }

  private static Role MapRole(NpgsqlDataReader reader)
  {
    return new Role(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetGuid(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5));
  }

  private static TeamMembership MapTeamMembership(NpgsqlDataReader reader)
  {
    return new TeamMembership(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetInt64(2),
      reader.GetInt64(3),
      reader.GetFieldValue<DateTimeOffset>(4));
  }

  private static UserRole MapUserRole(NpgsqlDataReader reader)
  {
    return new UserRole(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetInt64(2),
      reader.GetInt64(3),
      reader.GetFieldValue<DateTimeOffset>(4));
  }
}
