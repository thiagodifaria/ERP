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
    json(conn, 200, %{
      service: "workflow-runtime",
      status: "ready",
      dependencies: [
        %{name: "execution-store", status: "ready"},
        %{name: "timer-wheel", status: "ready"},
        %{name: "workflow-catalog", status: "pending-runtime-wiring"}
      ]
    })
  end

  get "/api/workflow-runtime/executions" do
    json(conn, 200, ExecutionStore.list())
  end

  get "/api/workflow-runtime/executions/summary" do
    json(conn, 200, ExecutionStore.summary())
  end

  get "/api/workflow-runtime/executions/:public_id" do
    case ExecutionStore.find(public_id) do
      nil ->
        json(conn, 404, %{code: "workflow_runtime_execution_not_found", message: "Execution was not found."})

      execution ->
        json(conn, 200, execution)
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
      execution = ExecutionStore.create(conn.body_params)
      json(conn, 201, execution)
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
end
