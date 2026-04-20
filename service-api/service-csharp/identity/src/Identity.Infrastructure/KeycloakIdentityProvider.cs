// KeycloakIdentityProvider integra autenticacao e diretoria do identity com o Keycloak local.
using System.Net;
using System.Net.Http.Json;
using System.Text;
using System.Text.Json;
using Identity.Application;

namespace Identity.Infrastructure;

public sealed class KeycloakIdentityProvider : IExternalIdentityProvider
{
  private readonly IdentityInfrastructureOptions _options;
  private readonly HttpClient _httpClient;

  public KeycloakIdentityProvider(IdentityInfrastructureOptions options)
  {
    _options = options;
    _httpClient = new HttpClient
    {
      Timeout = TimeSpan.FromSeconds(10)
    };
  }

  public ExternalIdentityUser EnsureUser(ExternalIdentityUpsertRequest request)
  {
    var adminToken = GetAdminToken();
    var user = FindUser(adminToken, request.SubjectId, request.Email);

    if (user is null)
    {
      user = CreateUser(adminToken, request);
    }
    else
    {
      UpdateUser(adminToken, user.Id, request);
    }

    if (!string.IsNullOrWhiteSpace(request.Password))
    {
      ResetPassword(adminToken, user.Id, request.Password!);
    }

    return new ExternalIdentityUser(user.Id, request.Email, request.Enabled);
  }

  public IdentityProviderTokenResult PasswordGrant(string email, string password)
  {
    return RequestToken([
      new("grant_type", "password"),
      new("client_id", _options.KeycloakClientId),
      new("username", email),
      new("password", password)
    ]);
  }

  public IdentityProviderTokenResult RefreshGrant(string refreshToken)
  {
    return RequestToken([
      new("grant_type", "refresh_token"),
      new("client_id", _options.KeycloakClientId),
      new("refresh_token", refreshToken)
    ]);
  }

  private IdentityProviderTokenResult RequestToken(List<KeyValuePair<string, string>> pairs)
  {
    var request = new HttpRequestMessage(
      HttpMethod.Post,
      $"{_options.KeycloakBaseUrl.TrimEnd('/')}/realms/{_options.KeycloakRealm}/protocol/openid-connect/token")
    {
      Content = new FormUrlEncodedContent(pairs)
    };

    using var response = _httpClient.Send(request);
    var payload = response.Content.ReadAsStringAsync().GetAwaiter().GetResult();

    if (response.StatusCode == HttpStatusCode.BadRequest || response.StatusCode == HttpStatusCode.Unauthorized)
    {
      throw new ExternalIdentityAuthenticationException("Invalid credentials.");
    }

    if (!response.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException($"Keycloak token request failed with status {(int)response.StatusCode}.");
    }

    using var document = JsonDocument.Parse(payload);
    var root = document.RootElement;
    var accessToken = root.GetProperty("access_token").GetString() ?? string.Empty;
    var refreshToken = root.GetProperty("refresh_token").GetString() ?? string.Empty;
    var expiresIn = root.GetProperty("expires_in").GetInt32();
    var refreshExpiresIn = root.GetProperty("refresh_expires_in").GetInt32();

    return new IdentityProviderTokenResult(
      ExtractSubjectId(accessToken),
      accessToken,
      refreshToken,
      DateTimeOffset.UtcNow.AddSeconds(expiresIn),
      DateTimeOffset.UtcNow.AddSeconds(refreshExpiresIn));
  }

  private string GetAdminToken()
  {
    var request = new HttpRequestMessage(
      HttpMethod.Post,
      $"{_options.KeycloakBaseUrl.TrimEnd('/')}/realms/master/protocol/openid-connect/token")
    {
      Content = new FormUrlEncodedContent(
      [
        new KeyValuePair<string, string>("grant_type", "password"),
        new KeyValuePair<string, string>("client_id", "admin-cli"),
        new KeyValuePair<string, string>("username", _options.KeycloakAdminUsername),
        new KeyValuePair<string, string>("password", _options.KeycloakAdminPassword)
      ])
    };

    using var response = _httpClient.Send(request);
    var payload = response.Content.ReadAsStringAsync().GetAwaiter().GetResult();

    if (!response.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException($"Keycloak admin token request failed with status {(int)response.StatusCode}.");
    }

    using var document = JsonDocument.Parse(payload);
    return document.RootElement.GetProperty("access_token").GetString()
      ?? throw new ExternalIdentityProviderException("Keycloak admin token response was missing access_token.");
  }

  private KeycloakUserRepresentation? FindUser(string adminToken, string? subjectId, string email)
  {
    if (!string.IsNullOrWhiteSpace(subjectId))
    {
      var byIdRequest = CreateAdminRequest(
        HttpMethod.Get,
        $"/admin/realms/{_options.KeycloakRealm}/users/{subjectId}",
        adminToken);
      using var byIdResponse = _httpClient.Send(byIdRequest);

      if (byIdResponse.IsSuccessStatusCode)
      {
        return byIdResponse.Content.ReadFromJsonAsync<KeycloakUserRepresentation>().GetAwaiter().GetResult();
      }
    }

    var searchRequest = CreateAdminRequest(
      HttpMethod.Get,
      $"/admin/realms/{_options.KeycloakRealm}/users?email={Uri.EscapeDataString(email)}&exact=true",
      adminToken);
    using var searchResponse = _httpClient.Send(searchRequest);

    if (!searchResponse.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException($"Keycloak user lookup failed with status {(int)searchResponse.StatusCode}.");
    }

    var users = searchResponse.Content.ReadFromJsonAsync<KeycloakUserRepresentation[]>().GetAwaiter().GetResult() ?? [];
    return users.FirstOrDefault();
  }

  private KeycloakUserRepresentation CreateUser(string adminToken, ExternalIdentityUpsertRequest request)
  {
    var profile = ResolveProfile(request);
    var createRequest = CreateAdminRequest(
      HttpMethod.Post,
      $"/admin/realms/{_options.KeycloakRealm}/users",
      adminToken);
    createRequest.Content = JsonContent.Create(new
    {
      username = request.Email,
      email = request.Email,
      firstName = profile.GivenName,
      lastName = profile.FamilyName,
      enabled = request.Enabled,
      emailVerified = true,
      attributes = new Dictionary<string, string[]>
      {
        ["display_name"] = [profile.DisplayName]
      }
    });

    using var createResponse = _httpClient.Send(createRequest);
    if (createResponse.StatusCode != HttpStatusCode.Created)
    {
      throw new ExternalIdentityProviderException($"Keycloak user create failed with status {(int)createResponse.StatusCode}.");
    }

    var location = createResponse.Headers.Location?.ToString();
    if (!string.IsNullOrWhiteSpace(location))
    {
      var userId = location.Split('/').Last();
      return new KeycloakUserRepresentation { Id = userId, Email = request.Email };
    }

    return FindUser(adminToken, null, request.Email)
      ?? throw new ExternalIdentityProviderException("Keycloak user was created but could not be loaded again.");
  }

  private void UpdateUser(string adminToken, string userId, ExternalIdentityUpsertRequest request)
  {
    var profile = ResolveProfile(request);
    var updateRequest = CreateAdminRequest(
      HttpMethod.Put,
      $"/admin/realms/{_options.KeycloakRealm}/users/{userId}",
      adminToken);
    updateRequest.Content = JsonContent.Create(new
    {
      username = request.Email,
      email = request.Email,
      firstName = profile.GivenName,
      lastName = profile.FamilyName,
      enabled = request.Enabled,
      emailVerified = true,
      attributes = new Dictionary<string, string[]>
      {
        ["display_name"] = [profile.DisplayName]
      }
    });

    using var updateResponse = _httpClient.Send(updateRequest);
    if (updateResponse.StatusCode != HttpStatusCode.NoContent)
    {
      throw new ExternalIdentityProviderException($"Keycloak user update failed with status {(int)updateResponse.StatusCode}.");
    }
  }

  private void ResetPassword(string adminToken, string userId, string password)
  {
    var resetRequest = CreateAdminRequest(
      HttpMethod.Put,
      $"/admin/realms/{_options.KeycloakRealm}/users/{userId}/reset-password",
      adminToken);
    resetRequest.Content = JsonContent.Create(new
    {
      type = "password",
      value = password,
      temporary = false
    });

    using var response = _httpClient.Send(resetRequest);
    if (response.StatusCode != HttpStatusCode.NoContent)
    {
      var payload = response.Content.ReadAsStringAsync().GetAwaiter().GetResult();
      throw new ExternalIdentityProviderException($"Keycloak password reset failed with status {(int)response.StatusCode}: {payload}");
    }
  }

  private HttpRequestMessage CreateAdminRequest(HttpMethod method, string path, string adminToken)
  {
    var request = new HttpRequestMessage(method, $"{_options.KeycloakBaseUrl.TrimEnd('/')}{path}");
    request.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", adminToken);
    return request;
  }

  private static string ExtractSubjectId(string accessToken)
  {
    var parts = accessToken.Split('.');
    if (parts.Length < 2)
    {
      return string.Empty;
    }

    var payload = parts[1]
      .Replace('-', '+')
      .Replace('_', '/');

    var padding = payload.Length % 4;
    if (padding > 0)
    {
      payload = payload.PadRight(payload.Length + (4 - padding), '=');
    }

    var bytes = Convert.FromBase64String(payload);
    using var document = JsonDocument.Parse(Encoding.UTF8.GetString(bytes));
    return document.RootElement.TryGetProperty("sub", out var sub)
      ? sub.GetString() ?? string.Empty
      : string.Empty;
  }

  private static NormalizedProfile ResolveProfile(ExternalIdentityUpsertRequest request)
  {
    var displayName = string.IsNullOrWhiteSpace(request.DisplayName)
      ? request.Email.Trim()
      : request.DisplayName.Trim();

    var fallbackTokens = displayName
      .Replace('.', ' ')
      .Replace('_', ' ')
      .Replace('-', ' ')
      .Split(' ', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

    var givenName = string.IsNullOrWhiteSpace(request.GivenName)
      ? null
      : request.GivenName.Trim();
    var familyName = string.IsNullOrWhiteSpace(request.FamilyName)
      ? null
      : request.FamilyName.Trim();

    if (string.IsNullOrWhiteSpace(givenName))
    {
      givenName = fallbackTokens.Length switch
      {
        > 1 => string.Join(' ', fallbackTokens[..^1]),
        1 => fallbackTokens[0],
        _ => "ERP"
      };
    }

    if (string.IsNullOrWhiteSpace(familyName))
    {
      familyName = fallbackTokens.Length switch
      {
        > 1 => fallbackTokens[^1],
        _ => "User"
      };
    }

    return new NormalizedProfile(givenName, familyName, displayName);
  }

  private sealed record NormalizedProfile(string GivenName, string FamilyName, string DisplayName);

  private sealed class KeycloakUserRepresentation
  {
    public string Id { get; set; } = string.Empty;

    public string Email { get; set; } = string.Empty;
  }
}
