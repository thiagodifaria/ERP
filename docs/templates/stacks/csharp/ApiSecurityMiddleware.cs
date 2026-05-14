using System.Security.Claims;

namespace StackStarter;

public sealed class ApiSecurityMiddleware(RequestDelegate next, string serviceName)
{
  public async Task Invoke(HttpContext context)
  {
    if (context.Request.Path.StartsWithSegments("/health"))
    {
      await next(context);
      return;
    }

    var auth = Authenticate(context);
    if (auth is null)
    {
      await WriteError(context, StatusCodes.Status401Unauthorized, "unauthorized", "Bearer token is invalid or missing.");
      return;
    }
    if (!HttpMethods.IsGet(context.Request.Method) && string.IsNullOrWhiteSpace(context.Request.Headers["X-Correlation-Id"]))
    {
      await WriteError(context, StatusCodes.Status400BadRequest, "correlation_id_required", "Mutation requests require X-Correlation-Id.");
      return;
    }
    if (!Authorize(context, auth))
    {
      await WriteError(context, StatusCodes.Status403Forbidden, "forbidden", "Request is not authorized.");
      return;
    }

    context.Request.Headers["X-ERP-Auth-Subject"] = auth.FindFirstValue(ClaimTypes.NameIdentifier) ?? "service";
    await next(context);
  }

  private ClaimsPrincipal? Authenticate(HttpContext context) => null;

  private bool Authorize(HttpContext context, ClaimsPrincipal auth) => !string.IsNullOrWhiteSpace(serviceName);

  private static async Task WriteError(HttpContext context, int status, string code, string message)
  {
    context.Response.StatusCode = status;
    context.Response.ContentType = "application/json";
    await context.Response.WriteAsJsonAsync(new { code, message });
  }
}
