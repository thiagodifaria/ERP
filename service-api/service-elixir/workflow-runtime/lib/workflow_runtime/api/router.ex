defmodule WorkflowRuntime.Api.Router do
  @moduledoc false

  use Plug.Router

  alias WorkflowRuntime.Application.ExecutionStore

  plug :match
  plug Plug.Parsers, parsers: [:json], pass: ["application/json"], json_decoder: Jason
  plug :security
  plug :dispatch

  def security(%Plug.Conn{request_path: "/health/" <> _rest} = conn, _opts), do: conn

  def security(conn, _opts) do
    if security_enforced?() do
      with :ok <- ensure_correlation(conn),
           {:ok, auth} <- authenticate(conn),
           :ok <- authorize_openfga(conn, auth) do
        conn
        |> Plug.Conn.put_req_header("x-erp-auth-subject", auth.subject)
        |> Plug.Conn.put_req_header("x-erp-auth-tenant", auth.tenant_slug)
        |> Plug.Conn.put_req_header("x-erp-auth-scopes", Enum.join(auth.scopes, " "))
      else
        {:error, :correlation} ->
          conn |> json(400, %{code: "correlation_id_required", message: "Mutation requests require X-Correlation-Id."}) |> halt()

        {:error, :auth} ->
          conn |> json(401, %{code: "unauthorized", message: "Bearer token is invalid or missing."}) |> halt()

        {:error, :openfga} ->
          conn |> json(403, %{code: "openfga_denied", message: "OpenFGA denied the request."}) |> halt()
      end
    else
      conn
    end
  end

  defp security_enforced? do
    mode = System.get_env("ERP_AUTH_ENFORCEMENT", "") |> String.trim() |> String.downcase()

    cond do
      mode in ["disabled", "off", "false"] -> false
      mode in ["enforced", "strict", "true"] -> true
      true ->
        environment = System.get_env("ERP_ENV", "local") |> String.trim() |> String.downcase()
        environment not in ["", "local", "dev", "development", "test", "testing"]
    end
  end

  defp ensure_correlation(%Plug.Conn{method: method}) when method in ["GET", "HEAD", "OPTIONS"], do: :ok

  defp ensure_correlation(conn) do
    case Plug.Conn.get_req_header(conn, "x-correlation-id") do
      [value | _] when value != "" -> :ok
      _ -> {:error, :correlation}
    end
  end

  defp authenticate(conn) do
    with ["Bearer " <> token | _] <- Plug.Conn.get_req_header(conn, "authorization") do
      internal_token = System.get_env("ERP_INTERNAL_SERVICE_TOKEN", "")

      cond do
        internal_token != "" and secure_compare(token, internal_token) ->
          {:ok, %{subject: "service:internal", tenant_slug: resolve_tenant(conn), scopes: ["service"]}}

        true ->
          case verify_jwt(token) do
            {:ok, claims} ->
              subject = claims["sub"] || claims["user_public_id"]
              tenant_slug = claims["tenant_slug"] || claims["tenant"] || resolve_tenant(conn)
              scopes = parse_scopes(claims["scope"])
              if subject, do: {:ok, %{subject: subject, tenant_slug: tenant_slug, scopes: scopes}}, else: {:error, :auth}

            _ ->
              {:error, :auth}
          end
      end
    else
      _ -> {:error, :auth}
    end
  end

  defp verify_jwt(token) do
    secret = System.get_env("ERP_JWT_HS256_SECRET", "")

    with true <- secret != "",
         [encoded_header, encoded_payload, signature] <- String.split(token, "."),
         {:ok, header_json} <- Base.url_decode64(encoded_header, padding: false),
         {:ok, %{"alg" => "HS256"}} <- Jason.decode(header_json),
         expected <- :crypto.mac(:hmac, :sha256, secret, "#{encoded_header}.#{encoded_payload}") |> Base.url_encode64(padding: false),
         true <- secure_compare(signature, expected),
         {:ok, payload_json} <- Base.url_decode64(encoded_payload, padding: false),
         {:ok, claims} <- Jason.decode(payload_json),
         true <- valid_expiration?(claims) do
      {:ok, claims}
    else
      _ -> {:error, :auth}
    end
  end

  defp authorize_openfga(conn, auth) do
    if String.downcase(System.get_env("ERP_OPENFGA_ENFORCEMENT", "")) == "true" do
      do_authorize_openfga(conn, auth)
    else
      :ok
    end
  end

  defp do_authorize_openfga(conn, auth) do
    base_url = System.get_env("OPENFGA_BASE_URL", "") |> String.trim_trailing("/")
    store_id = System.get_env("OPENFGA_STORE_ID", "")

    if base_url == "" or store_id == "" do
      {:error, :openfga}
    else
      {:ok, _apps} = Application.ensure_all_started(:inets)
      relation = if conn.method in ["GET", "HEAD", "OPTIONS"], do: "read", else: "write"
      object = if auth.tenant_slug == "", do: "service:workflow-runtime", else: "tenant:#{normalize_object(auth.tenant_slug)}"
      user = if String.starts_with?(auth.subject, "service:"), do: auth.subject, else: "user:#{auth.subject}"

      payload =
        %{tuple_key: %{user: user, relation: relation, object: object}}
        |> maybe_put_authorization_model()
        |> Jason.encode!()

      url = String.to_charlist("#{base_url}/stores/#{store_id}/check")

      case :httpc.request(:post, {url, [{~c"content-type", ~c"application/json"}], ~c"application/json", payload}, [{:timeout, 2_000}], []) do
        {:ok, {{_, status, _}, _headers, body}} when status in 200..299 ->
          case Jason.decode(to_string(body)) do
            {:ok, %{"allowed" => true}} -> :ok
            _ -> {:error, :openfga}
          end

        _ ->
          {:error, :openfga}
      end
    end
  end

  defp maybe_put_authorization_model(payload) do
    case System.get_env("OPENFGA_AUTHORIZATION_MODEL_ID", "") do
      "" -> payload
      model_id -> Map.put(payload, :authorization_model_id, model_id)
    end
  end

  defp valid_expiration?(%{"exp" => exp}) when is_integer(exp), do: exp > DateTime.to_unix(DateTime.utc_now())
  defp valid_expiration?(_claims), do: true

  defp parse_scopes(scope) when is_binary(scope), do: String.split(scope, " ", trim: true)
  defp parse_scopes(scope) when is_list(scope), do: Enum.map(scope, &to_string/1)
  defp parse_scopes(_scope), do: []

  defp resolve_tenant(conn) do
    conn = Plug.Conn.fetch_query_params(conn)

    Plug.Conn.get_req_header(conn, "x-tenant-slug")
    |> Kernel.++(Plug.Conn.get_req_header(conn, "x-erp-tenant-slug"))
    |> List.first()
    |> Kernel.||(conn.query_params["tenant_slug"] || "")
  end

  defp secure_compare(left, right) when byte_size(left) == byte_size(right), do: Plug.Crypto.secure_compare(left, right)
  defp secure_compare(_left, _right), do: false

  defp normalize_object(value), do: value |> String.trim() |> String.downcase() |> String.replace(" ", "-")

  get "/health/live" do
    json(conn, 200, %{service: "workflow-runtime", status: "live"})
  end

  get "/health/ready" do
    json(conn, 200, %{service: "workflow-runtime", status: "ready"})
  end

  get "/health/details" do
    repository_dependency =
      if WorkflowRuntime.Infrastructure.Postgres.enabled?() do
        %{name: "postgresql", status: if(WorkflowRuntime.Infrastructure.Postgres.ready?(), do: "ready", else: "not_ready")}
      else
        %{name: "execution-store", status: "ready"}
      end

    workflow_catalog_status =
      case ExecutionStore.catalog_status() do
        :ready -> "ready"
        _ -> "not_ready"
      end

    json(conn, 200, %{
      service: "workflow-runtime",
      status: "ready",
      dependencies: [
        repository_dependency,
        %{name: "timer-wheel", status: "ready"},
        %{name: "workflow-catalog", status: workflow_catalog_status}
      ]
    })
  end

  get "/api/workflow-runtime/capabilities" do
    json(conn, 200, %{
      service: "workflow-runtime",
      timers: %{enabled: true, mode: "timer-wheel", supportsDelayActions: true, supportsSchedules: false},
      retries: %{enabled: true, maxAttempts: 5, supportsManualRetry: true},
      compensations: %{enabled: true, mode: "basic"},
      transports: %{catalogValidation: true, postgres: WorkflowRuntime.Infrastructure.Postgres.enabled?()}
    })
  end

  get "/api/workflow-runtime/executions" do
    conn = Plug.Conn.fetch_query_params(conn)

    json(conn, 200, ExecutionStore.list(conn.query_params))
  end

  get "/api/workflow-runtime/executions/summary" do
    conn = Plug.Conn.fetch_query_params(conn)
    json(conn, 200, ExecutionStore.summary(conn.query_params))
  end

  get "/api/workflow-runtime/executions/summary/by-workflow" do
    conn = Plug.Conn.fetch_query_params(conn)
    json(conn, 200, ExecutionStore.summary_by_workflow(conn.query_params))
  end

  get "/api/workflow-runtime/executions/:public_id" do
    case ExecutionStore.find(public_id) do
      nil ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      execution ->
        json(conn, 200, execution)
    end
  end

  get "/api/workflow-runtime/executions/:public_id/transitions" do
    case ExecutionStore.find(public_id) do
      nil ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      _execution ->
        json(conn, 200, ExecutionStore.transitions(public_id))
    end
  end

  get "/api/workflow-runtime/executions/:public_id/actions" do
    case ExecutionStore.find(public_id) do
      nil ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      _execution ->
        json(conn, 200, ExecutionStore.actions(public_id))
    end
  end

  get "/api/workflow-runtime/executions/:public_id/timeline" do
    case ExecutionStore.find(public_id) do
      nil ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      _execution ->
        json(conn, 200, ExecutionStore.timeline(public_id))
    end
  end

  post "/api/workflow-runtime/executions/:public_id/start" do
    transition_execution(conn, public_id, "running", ["pending"])
  end

  post "/api/workflow-runtime/executions/:public_id/complete" do
    transition_execution(conn, public_id, "completed", ["running"])
  end

  post "/api/workflow-runtime/executions/:public_id/fail" do
    transition_execution(conn, public_id, "failed", ["pending", "running"])
  end

  post "/api/workflow-runtime/executions/:public_id/cancel" do
    transition_execution(conn, public_id, "cancelled", ["pending", "running"])
  end

  post "/api/workflow-runtime/executions/:public_id/retry" do
    case ExecutionStore.retry(public_id, ["failed", "cancelled"]) do
      {:ok, execution} ->
        json(conn, 200, execution)

      {:error, :not_found} ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      {:error, :invalid_transition} ->
        json(conn, 409, %{
          code: "workflow_runtime_execution_retry_invalid",
          message: "Execution cannot be retried from the current status."
        })
    end
  end

  post "/api/workflow-runtime/executions/:public_id/advance" do
    case ExecutionStore.advance(public_id) do
      {:ok, execution} ->
        json(conn, 200, execution)

      {:error, :not_found} ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      {:error, :waiting} ->
        json(conn, 409, %{
          code: "workflow_runtime_execution_waiting",
          message: "Execution is waiting for the configured delay to elapse."
        })

      {:error, :invalid_transition} ->
        json(conn, 409, %{
          code: "workflow_runtime_execution_transition_invalid",
          message: "Execution cannot transition from the current status."
        })

      {:error, :workflow_definition_version_not_found} ->
        json(conn, 409, %{
          code: "workflow_runtime_definition_version_not_found",
          message: "Workflow definition has no published version available for runtime."
        })
    end
  end

  post "/api/workflow-runtime/executions" do
    required_fields = ["workflowDefinitionKey", "subjectType", "subjectPublicId", "initiatedBy"]

    if Enum.any?(required_fields, fn field ->
         conn.body_params[field] |> to_string() |> String.trim() == ""
       end) do
      json(conn, 400, %{
        code: "workflow_runtime_execution_payload_invalid",
        message: "Execution payload is invalid."
      })
    else
      case ExecutionStore.create(conn.body_params) do
        {:error, :tenant_not_found} ->
          json(conn, 404, %{
            code: "workflow_runtime_tenant_not_found",
            message: "Workflow runtime tenant was not found."
          })

        {:error, :workflow_definition_not_found} ->
          json(conn, 404, %{
            code: "workflow_runtime_definition_not_found",
            message: "Workflow definition was not found in the runtime catalog."
          })

        {:error, :workflow_definition_version_not_found} ->
          json(conn, 409, %{
            code: "workflow_runtime_definition_version_not_found",
            message: "Workflow definition has no published version available for runtime."
          })

        {:error, :workflow_definition_inactive} ->
          json(conn, 409, %{
            code: "workflow_runtime_definition_inactive",
            message: "Workflow definition is not active for runtime execution."
          })

        execution ->
          json(conn, 201, execution)
      end
    end
  end

  match _ do
    json(conn, 404, %{code: "workflow_runtime_route_not_found", message: "Route was not found."})
  end

  defp json(conn, status, payload) do
    body = Jason.encode!(payload)

    conn
    |> Plug.Conn.put_resp_content_type("application/json")
    |> Plug.Conn.send_resp(status, body)
  end

  defp transition_execution(conn, public_id, next_status, allowed_statuses) do
    case ExecutionStore.transition(public_id, next_status, allowed_statuses) do
      {:ok, execution} ->
        json(conn, 200, execution)

      {:error, :not_found} ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      {:error, :invalid_transition} ->
        json(conn, 409, %{
          code: "workflow_runtime_execution_transition_invalid",
          message: "Execution cannot transition from the current status."
        })
    end
  end
end
