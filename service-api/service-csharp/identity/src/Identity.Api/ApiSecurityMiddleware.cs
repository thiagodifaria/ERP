using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using System.Net.Http.Json;

namespace Identity.Api;

public sealed class ApiSecurityMiddleware
{
  private static readonly JsonSerializerOptions JsonOptions = new(JsonSerializerDefaults.Web);
  private readonly RequestDelegate _next;
  private readonly string _serviceName;

  public ApiSecurityMiddleware(RequestDelegate next, string serviceName)
  {
    _next = next;
    _serviceName = serviceName;
  }

  public async Task InvokeAsync(HttpContext context)
  {
    if (!RequiresEnforcement() || IsPublicRequest(context.Request))
    {
      await _next(context);
      return;
    }

    if (!HasCorrelation(context.Request))
    {
      await WriteError(context, StatusCodes.Status400BadRequest, "correlation_id_required", "Mutation requests require X-Correlation-Id.");
      return;
    }

    var auth = await Authenticate(context.Request);
    if (!auth.IsAuthenticated)
    {
      await WriteError(context, StatusCodes.Status401Unauthorized, "unauthorized", auth.Error ?? "Bearer token is invalid or missing.");
      return;
    }

    context.Request.Headers["X-ERP-Auth-Subject"] = auth.Subject;
    context.Request.Headers["X-ERP-Auth-Tenant"] = auth.TenantSlug;
    context.Request.Headers["X-ERP-Auth-Scopes"] = string.Join(' ', auth.Scopes);

    var openFga = await AuthorizeWithOpenFga(context.Request, auth);
    if (!openFga.IsAuthorized)
    {
      await WriteError(context, openFga.StatusCode, openFga.Code, openFga.Message);
      return;
    }

    await _next(context);
  }

  private static bool RequiresEnforcement()
  {
    var mode = Environment.GetEnvironmentVariable("ERP_AUTH_ENFORCEMENT")?.Trim().ToLowerInvariant();
    if (mode is "disabled" or "off" or "false")
    {
      return false;
    }

    if (mode is "enforced" or "strict" or "true")
    {
      return true;
    }

    var environment =
      Environment.GetEnvironmentVariable("ERP_ENV") ??
      Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT") ??
      Environment.GetEnvironmentVariable("DOTNET_ENVIRONMENT") ??
      "local";

    return !IsLocalEnvironment(environment);
  }

  private static bool IsLocalEnvironment(string environment)
  {
    var normalized = environment.Trim().ToLowerInvariant();
    return normalized is "local" or "dev" or "development" or "test" or "testing";
  }

  private static bool IsPublicRequest(HttpRequest request)
  {
    if (request.Path.StartsWithSegments("/health"))
    {
      return true;
    }

    if (!request.Path.StartsWithSegments("/api/identity"))
    {
      return false;
    }

    var path = request.Path.Value?.ToLowerInvariant() ?? "";
    return path is "/api/identity/sessions/login" or "/api/identity/sessions/refresh" or "/api/identity/password-recovery"
      || (path.StartsWith("/api/identity/password-recovery/") && path.EndsWith("/complete"))
      || (path.StartsWith("/api/identity/invites/") && path.EndsWith("/accept"));
  }

  private static bool HasCorrelation(HttpRequest request)
  {
    if (HttpMethods.IsGet(request.Method) || HttpMethods.IsHead(request.Method) || HttpMethods.IsOptions(request.Method))
    {
      return true;
    }

    return request.Headers.TryGetValue("X-Correlation-Id", out var value) && !string.IsNullOrWhiteSpace(value);
  }

  private static async Task<AuthResult> Authenticate(HttpRequest request)
  {
    var authorization = request.Headers.Authorization.ToString();
    if (!authorization.StartsWith("Bearer ", StringComparison.OrdinalIgnoreCase))
    {
      return AuthResult.Fail("Authorization header must contain a Bearer token.");
    }

    var token = authorization["Bearer ".Length..].Trim();
    var internalToken = Environment.GetEnvironmentVariable("ERP_INTERNAL_SERVICE_TOKEN");
    if (!string.IsNullOrWhiteSpace(internalToken) && FixedTimeEquals(token, internalToken))
    {
      return AuthResult.Success("service:internal", ResolveTenant(request), ["service"]);
    }

    var secret = Environment.GetEnvironmentVariable("ERP_JWT_HS256_SECRET");
    if (string.IsNullOrWhiteSpace(secret))
    {
      return AuthResult.Fail("ERP_JWT_HS256_SECRET is required when API auth enforcement is enabled.");
    }

    var parts = token.Split('.');
    if (parts.Length != 3)
    {
      return AuthResult.Fail("JWT must have header, payload and signature.");
    }

    var signedContent = $"{parts[0]}.{parts[1]}";
    var expectedSignature = Base64UrlEncode(HMACSHA256.HashData(Encoding.UTF8.GetBytes(secret), Encoding.ASCII.GetBytes(signedContent)));
    if (!FixedTimeEquals(parts[2], expectedSignature))
    {
      return AuthResult.Fail("JWT signature is invalid.");
    }

    using var header = JsonDocument.Parse(Base64UrlDecode(parts[0]));
    if (!header.RootElement.TryGetProperty("alg", out var alg) || alg.GetString() != "HS256")
    {
      return AuthResult.Fail("JWT alg must be HS256.");
    }

    using var payload = JsonDocument.Parse(Base64UrlDecode(parts[1]));
    if (payload.RootElement.TryGetProperty("exp", out var exp) && exp.TryGetInt64(out var expSeconds))
    {
      var expiresAt = DateTimeOffset.FromUnixTimeSeconds(expSeconds);
      if (expiresAt <= DateTimeOffset.UtcNow)
      {
        return AuthResult.Fail("JWT is expired.");
      }
    }

    var audience = Environment.GetEnvironmentVariable("ERP_JWT_AUDIENCE");
    if (!string.IsNullOrWhiteSpace(audience) && payload.RootElement.TryGetProperty("aud", out var aud) && aud.GetString() != audience)
    {
      return AuthResult.Fail("JWT audience is invalid.");
    }

    var subject = ReadClaim(payload.RootElement, "sub") ?? ReadClaim(payload.RootElement, "user_public_id") ?? "unknown";
    var tenant = ReadClaim(payload.RootElement, "tenant_slug") ?? ReadClaim(payload.RootElement, "tenant") ?? ResolveTenant(request);
    var scopes = ReadScopes(payload.RootElement);
    await Task.CompletedTask;
    return AuthResult.Success(subject, tenant, scopes);
  }

  private async Task<OpenFgaResult> AuthorizeWithOpenFga(HttpRequest request, AuthResult auth)
  {
    if (!string.Equals(Environment.GetEnvironmentVariable("ERP_OPENFGA_ENFORCEMENT"), "true", StringComparison.OrdinalIgnoreCase))
    {
      return OpenFgaResult.Allow();
    }

    var baseUrl = Environment.GetEnvironmentVariable("OPENFGA_BASE_URL")?.TrimEnd('/');
    var storeId = Environment.GetEnvironmentVariable("OPENFGA_STORE_ID");
    if (string.IsNullOrWhiteSpace(baseUrl) || string.IsNullOrWhiteSpace(storeId))
    {
      return OpenFgaResult.Deny(StatusCodes.Status503ServiceUnavailable, "openfga_not_configured", "OpenFGA enforcement is enabled but OPENFGA_BASE_URL or OPENFGA_STORE_ID is missing.");
    }

    var relation = HttpMethods.IsGet(request.Method) || HttpMethods.IsHead(request.Method) ? "read" : "write";
    var targetObject = string.IsNullOrWhiteSpace(auth.TenantSlug) ? $"service:{NormalizeObject(_serviceName)}" : $"tenant:{NormalizeObject(auth.TenantSlug)}";
    var modelId = Environment.GetEnvironmentVariable("OPENFGA_AUTHORIZATION_MODEL_ID");
    var payload = new Dictionary<string, object?>
    {
      ["tuple_key"] = new Dictionary<string, string>
      {
        ["user"] = auth.Subject.StartsWith("service:", StringComparison.Ordinal) ? auth.Subject : $"user:{auth.Subject}",
        ["relation"] = relation,
        ["object"] = targetObject
      }
    };

    if (!string.IsNullOrWhiteSpace(modelId))
    {
      payload["authorization_model_id"] = modelId;
    }

    using var client = new HttpClient { Timeout = TimeSpan.FromSeconds(2) };
    using var response = await client.PostAsJsonAsync($"{baseUrl}/stores/{storeId}/check", payload, JsonOptions);
    if (!response.IsSuccessStatusCode)
    {
      return OpenFgaResult.Deny(StatusCodes.Status403Forbidden, "openfga_denied", "OpenFGA denied the request.");
    }

    using var body = JsonDocument.Parse(await response.Content.ReadAsStringAsync());
    var allowed = body.RootElement.TryGetProperty("allowed", out var allowedProperty) && allowedProperty.GetBoolean();
    return allowed
      ? OpenFgaResult.Allow()
      : OpenFgaResult.Deny(StatusCodes.Status403Forbidden, "openfga_denied", "OpenFGA denied the request.");
  }

  private static string ResolveTenant(HttpRequest request)
  {
    if (request.Headers.TryGetValue("X-Tenant-Slug", out var headerTenant) && !string.IsNullOrWhiteSpace(headerTenant))
    {
      return headerTenant.ToString();
    }

    if (request.Headers.TryGetValue("X-ERP-Tenant-Slug", out var erpTenant) && !string.IsNullOrWhiteSpace(erpTenant))
    {
      return erpTenant.ToString();
    }

    return request.Query.TryGetValue("tenant_slug", out var queryTenant) ? queryTenant.ToString() : "";
  }

  private static string? ReadClaim(JsonElement payload, string name)
  {
    return payload.TryGetProperty(name, out var property) && property.ValueKind == JsonValueKind.String
      ? property.GetString()
      : null;
  }

  private static string[] ReadScopes(JsonElement payload)
  {
    if (!payload.TryGetProperty("scope", out var scope))
    {
      return [];
    }

    if (scope.ValueKind == JsonValueKind.String)
    {
      return scope.GetString()?.Split(' ', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries) ?? [];
    }

    return scope.ValueKind == JsonValueKind.Array
      ? scope.EnumerateArray().Where(item => item.ValueKind == JsonValueKind.String).Select(item => item.GetString()!).ToArray()
      : [];
  }

  private static string NormalizeObject(string value)
  {
    return value.Trim().ToLowerInvariant().Replace(' ', '-');
  }

  private static byte[] Base64UrlDecode(string value)
  {
    var padded = value.Replace('-', '+').Replace('_', '/');
    padded = padded.PadRight(padded.Length + (4 - padded.Length % 4) % 4, '=');
    return Convert.FromBase64String(padded);
  }

  private static string Base64UrlEncode(byte[] value)
  {
    return Convert.ToBase64String(value).TrimEnd('=').Replace('+', '-').Replace('/', '_');
  }

  private static bool FixedTimeEquals(string left, string right)
  {
    return CryptographicOperations.FixedTimeEquals(Encoding.UTF8.GetBytes(left), Encoding.UTF8.GetBytes(right));
  }

  private static async Task WriteError(HttpContext context, int statusCode, string code, string message)
  {
    context.Response.StatusCode = statusCode;
    context.Response.ContentType = "application/json";
    await context.Response.WriteAsJsonAsync(new { code, message }, JsonOptions);
  }

  private sealed record AuthResult(bool IsAuthenticated, string Subject, string TenantSlug, string[] Scopes, string? Error)
  {
    public static AuthResult Success(string subject, string tenantSlug, string[] scopes) => new(true, subject, tenantSlug, scopes, null);
    public static AuthResult Fail(string error) => new(false, "", "", [], error);
  }

  private sealed record OpenFgaResult(bool IsAuthorized, int StatusCode, string Code, string Message)
  {
    public static OpenFgaResult Allow() => new(true, StatusCodes.Status200OK, "", "");
    public static OpenFgaResult Deny(int statusCode, string code, string message) => new(false, statusCode, code, message);
  }
}
