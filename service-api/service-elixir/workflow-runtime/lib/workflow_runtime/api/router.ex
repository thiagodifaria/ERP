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
