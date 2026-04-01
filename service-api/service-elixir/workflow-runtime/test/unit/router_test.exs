defmodule WorkflowRuntime.Api.RouterTest do
  use ExUnit.Case, async: false
  import Plug.Test

  alias WorkflowRuntime.Api.Router

  setup do
    WorkflowRuntime.Application.ExecutionStore.reset()
    :ok
  end

  test "live route returns ok" do
    conn = conn(:get, "/health/live") |> Router.call([])

    assert conn.status == 200
    assert Jason.decode!(conn.resp_body)["service"] == "workflow-runtime"
  end

  test "details route exposes runtime dependencies" do
    conn = conn(:get, "/health/details") |> Router.call([])
    payload = Jason.decode!(conn.resp_body)

    assert conn.status == 200
    assert payload["status"] == "ready"
    assert Enum.any?(payload["dependencies"], &(&1["name"] == "execution-store"))
  end

  test "execution list and summary start empty" do
    list_conn = conn(:get, "/api/workflow-runtime/executions") |> Router.call([])
    summary_conn = conn(:get, "/api/workflow-runtime/executions/summary") |> Router.call([])

    assert list_conn.status == 200
    assert Jason.decode!(list_conn.resp_body) == []

    assert summary_conn.status == 200
    assert Jason.decode!(summary_conn.resp_body) == %{
             "total" => 0,
             "pending" => 0,
             "running" => 0,
             "completed" => 0,
             "failed" => 0
           }
  end
end
