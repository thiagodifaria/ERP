// OpenFgaAuthorizationGraph sincroniza o acesso minimo de tenant no runtime local.
using System.Net;
using System.Net.Http.Json;
using System.Text;
using System.Text.Json;
using Identity.Application;

namespace Identity.Infrastructure;

public sealed class OpenFgaAuthorizationGraph : IAuthorizationGraph
{
  private readonly IdentityInfrastructureOptions _options;
  private readonly HttpClient _httpClient;
  private readonly object _syncRoot = new();
  private string? _storeId;
  private string? _authorizationModelId;

  public OpenFgaAuthorizationGraph(IdentityInfrastructureOptions options)
  {
    _options = options;
    _httpClient = new HttpClient
    {
      Timeout = TimeSpan.FromSeconds(10)
    };
  }

  public void SyncTenantAccess(string tenantSlug, Guid userPublicId, IReadOnlyCollection<string> roleCodes, bool active)
  {
    EnsureReady();

    var user = $"user:{userPublicId}";
    var obj = $"tenant:{tenantSlug}";
    var desiredRelations = new HashSet<string>(StringComparer.OrdinalIgnoreCase);

    if (active && roleCodes.Count > 0)
    {
      desiredRelations.Add("member");
    }

    if (active && roleCodes.Any(roleCode => roleCode.Equals("admin", StringComparison.OrdinalIgnoreCase) || roleCode.Equals("owner", StringComparison.OrdinalIgnoreCase)))
    {
      desiredRelations.Add("admin");
    }

    if (active && roleCodes.Any(roleCode => roleCode.Equals("owner", StringComparison.OrdinalIgnoreCase)))
    {
      desiredRelations.Add("owner");
    }

    var tuples = new[] { "member", "admin", "owner" };
    var existingRelations = ReadRelations(user, obj);
    var writes = tuples
      .Where(relation => desiredRelations.Contains(relation) && !existingRelations.Contains(relation))
      .Select(relation => new { user, relation, @object = obj })
      .ToArray();
    var deletes = tuples
      .Where(relation => !desiredRelations.Contains(relation) && existingRelations.Contains(relation))
      .Select(relation => new { user, relation, @object = obj })
      .ToArray();

    if (writes.Length == 0 && deletes.Length == 0)
    {
      return;
    }

    var payload = new Dictionary<string, object?>
    {
      ["authorization_model_id"] = _authorizationModelId
    };

    if (writes.Length > 0)
    {
      payload["writes"] = new
      {
        tuple_keys = writes
      };
    }

    if (deletes.Length > 0)
    {
      payload["deletes"] = new
      {
        tuple_keys = deletes
      };
    }

    var request = new HttpRequestMessage(
      HttpMethod.Post,
      $"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores/{_storeId}/write")
    {
      Content = JsonContent.Create(payload)
    };

    using var response = _httpClient.Send(request);
    if (!response.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException(
        $"OpenFGA write failed with status {(int)response.StatusCode}: {ReadResponseBody(response)}");
    }
  }

  public bool CanAccessTenant(string tenantSlug, Guid userPublicId)
  {
    EnsureReady();

    var request = new HttpRequestMessage(
      HttpMethod.Post,
      $"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores/{_storeId}/check")
    {
      Content = JsonContent.Create(new
      {
        tuple_key = new
        {
          user = $"user:{userPublicId}",
          relation = "can_access",
          @object = $"tenant:{tenantSlug}"
        },
        authorization_model_id = _authorizationModelId
      })
    };

    using var response = _httpClient.Send(request);
    if (!response.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException(
        $"OpenFGA check failed with status {(int)response.StatusCode}: {ReadResponseBody(response)}");
    }

    using var document = JsonDocument.Parse(response.Content.ReadAsStringAsync().GetAwaiter().GetResult());
    return document.RootElement.TryGetProperty("allowed", out var allowed) && allowed.GetBoolean();
  }

  private void EnsureReady()
  {
    if (!string.IsNullOrWhiteSpace(_storeId) && !string.IsNullOrWhiteSpace(_authorizationModelId))
    {
      return;
    }

    lock (_syncRoot)
    {
      if (string.IsNullOrWhiteSpace(_storeId))
      {
        _storeId = ResolveStoreId();
      }

      if (string.IsNullOrWhiteSpace(_authorizationModelId))
      {
        _authorizationModelId = ResolveAuthorizationModelId(_storeId!);
      }
    }
  }

  private HashSet<string> ReadRelations(string user, string obj)
  {
    var request = new HttpRequestMessage(
      HttpMethod.Post,
      $"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores/{_storeId}/read")
    {
      Content = JsonContent.Create(new
      {
        tuple_key = new
        {
          user,
          @object = obj
        }
      })
    };

    using var response = _httpClient.Send(request);
    if (!response.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException(
        $"OpenFGA read failed with status {(int)response.StatusCode}: {ReadResponseBody(response)}");
    }

    using var document = JsonDocument.Parse(response.Content.ReadAsStringAsync().GetAwaiter().GetResult());
    var relations = new HashSet<string>(StringComparer.OrdinalIgnoreCase);

    if (!document.RootElement.TryGetProperty("tuples", out var tuples))
    {
      return relations;
    }

    foreach (var tuple in tuples.EnumerateArray())
    {
      if (!tuple.TryGetProperty("key", out var key))
      {
        continue;
      }

      if (!key.TryGetProperty("relation", out var relation))
      {
        continue;
      }

      var relationValue = relation.GetString();
      if (!string.IsNullOrWhiteSpace(relationValue))
      {
        relations.Add(relationValue);
      }
    }

    return relations;
  }

  private string ResolveStoreId()
  {
    using var listResponse = _httpClient.GetAsync($"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores").GetAwaiter().GetResult();
    if (!listResponse.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException($"OpenFGA store list failed with status {(int)listResponse.StatusCode}.");
    }

    using var listDocument = JsonDocument.Parse(listResponse.Content.ReadAsStringAsync().GetAwaiter().GetResult());
    if (listDocument.RootElement.TryGetProperty("stores", out var stores))
    {
      foreach (var store in stores.EnumerateArray())
      {
        if (store.TryGetProperty("name", out var name)
          && name.GetString() == _options.OpenFgaStoreName
          && store.TryGetProperty("id", out var id))
        {
          return id.GetString() ?? string.Empty;
        }
      }
    }

    var createRequest = new HttpRequestMessage(HttpMethod.Post, $"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores")
    {
      Content = JsonContent.Create(new { name = _options.OpenFgaStoreName })
    };
    using var createResponse = _httpClient.Send(createRequest);
    if (createResponse.StatusCode != HttpStatusCode.Created)
    {
      throw new ExternalIdentityProviderException($"OpenFGA store create failed with status {(int)createResponse.StatusCode}.");
    }

    using var createDocument = JsonDocument.Parse(createResponse.Content.ReadAsStringAsync().GetAwaiter().GetResult());
    return createDocument.RootElement.GetProperty("id").GetString()
      ?? throw new ExternalIdentityProviderException("OpenFGA store id was missing.");
  }

  private string ResolveAuthorizationModelId(string storeId)
  {
    using var listResponse = _httpClient.GetAsync($"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores/{storeId}/authorization-models").GetAwaiter().GetResult();
    if (!listResponse.IsSuccessStatusCode)
    {
      throw new ExternalIdentityProviderException($"OpenFGA authorization model list failed with status {(int)listResponse.StatusCode}.");
    }

    using var listDocument = JsonDocument.Parse(listResponse.Content.ReadAsStringAsync().GetAwaiter().GetResult());
    if (listDocument.RootElement.TryGetProperty("authorization_models", out var models))
    {
      var latest = models.EnumerateArray().FirstOrDefault();
      if (latest.ValueKind != JsonValueKind.Undefined && latest.TryGetProperty("id", out var id))
      {
        return id.GetString() ?? string.Empty;
      }
    }

    var writeRequest = new HttpRequestMessage(HttpMethod.Post, $"{_options.OpenFgaBaseUrl.TrimEnd('/')}/stores/{storeId}/authorization-models")
    {
      Content = new StringContent(BuildAuthorizationModelJson(), Encoding.UTF8, "application/json")
    };
    using var writeResponse = _httpClient.Send(writeRequest);
    if (writeResponse.StatusCode != HttpStatusCode.Created)
    {
      throw new ExternalIdentityProviderException($"OpenFGA authorization model write failed with status {(int)writeResponse.StatusCode}.");
    }

    using var writeDocument = JsonDocument.Parse(writeResponse.Content.ReadAsStringAsync().GetAwaiter().GetResult());
    return writeDocument.RootElement.GetProperty("authorization_model_id").GetString()
      ?? throw new ExternalIdentityProviderException("OpenFGA authorization model id was missing.");
  }

  private static string BuildAuthorizationModelJson()
  {
    return """
      {
        "schema_version": "1.1",
        "type_definitions": [
          {
            "type": "user"
          },
          {
            "type": "tenant",
            "relations": {
              "member": { "this": {} },
              "admin": { "this": {} },
              "owner": { "this": {} },
              "can_access": {
                "union": {
                  "child": [
                    { "computedUserset": { "relation": "member" } },
                    { "computedUserset": { "relation": "admin" } },
                    { "computedUserset": { "relation": "owner" } }
                  ]
                }
              }
            },
            "metadata": {
              "relations": {
                "member": {
                  "directly_related_user_types": [
                    { "type": "user" }
                  ]
                },
                "admin": {
                  "directly_related_user_types": [
                    { "type": "user" }
                  ]
                },
                "owner": {
                  "directly_related_user_types": [
                    { "type": "user" }
                  ]
                }
              }
            }
          }
        ]
      }
      """;
  }

  private static string ReadResponseBody(HttpResponseMessage response)
  {
    return response.Content.ReadAsStringAsync().GetAwaiter().GetResult();
  }
}
