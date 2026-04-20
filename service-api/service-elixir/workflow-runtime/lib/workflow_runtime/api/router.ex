defmodule WorkflowRuntime.Api.Router do
  @moduledoc false

  use Plug.Router

  alias WorkflowRuntime.Application.ExecutionStore

  plug :match
  plug Plug.Parsers, parsers: [:json], pass: ["application/json"], json_decoder: Jason
  plug :dispatch

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
