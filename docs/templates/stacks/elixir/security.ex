defmodule StackStarter.Security do
  import Plug.Conn

  def call(%Plug.Conn{request_path: "/health/" <> _rest} = conn, _opts), do: conn

  def call(conn, opts) do
    service_name = Keyword.fetch!(opts, :service_name)

    with {:ok, auth} <- authenticate(conn),
         :ok <- ensure_correlation(conn),
         :ok <- authorize(service_name, conn, auth) do
      conn
      |> put_req_header("x-erp-auth-subject", auth.subject)
      |> put_req_header("x-erp-auth-tenant", auth.tenant_slug)
    else
      {:error, :auth} -> error(conn, 401, "unauthorized", "Bearer token is invalid or missing.")
      {:error, :correlation} -> error(conn, 400, "correlation_id_required", "Mutation requests require X-Correlation-Id.")
      {:error, :forbidden} -> error(conn, 403, "forbidden", "Request is not authorized.")
    end
  end

  defp authenticate(_conn), do: {:error, :auth}

  defp ensure_correlation(%Plug.Conn{method: method}) when method in ["GET", "HEAD", "OPTIONS"], do: :ok
  defp ensure_correlation(conn), do: if(get_req_header(conn, "x-correlation-id") == [], do: {:error, :correlation}, else: :ok)

  defp authorize(service_name, _conn, auth) do
    if service_name != "" and auth.subject != "", do: :ok, else: {:error, :forbidden}
  end

  defp error(conn, status, code, message) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(status, Jason.encode!(%{code: code, message: message}))
    |> halt()
  end
end
