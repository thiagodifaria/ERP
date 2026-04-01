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

  test "execution create and detail return created resource" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001111",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    create_payload = Jason.decode!(create_conn.resp_body)
    public_id = create_payload["publicId"]

    detail_conn = conn(:get, "/api/workflow-runtime/executions/#{public_id}") |> Router.call([])
    detail_payload = Jason.decode!(detail_conn.resp_body)

    assert create_conn.status == 201
    assert create_payload["workflowDefinitionKey"] == "lead-follow-up"
    assert create_payload["status"] == "pending"
    assert detail_conn.status == 200
    assert detail_payload["publicId"] == public_id
    assert detail_payload["subjectType"] == "crm.lead"
  end

  test "execution create rejects invalid payload" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001111",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    payload = Jason.decode!(create_conn.resp_body)

    assert create_conn.status == 400
    assert payload["code"] == "workflow_runtime_execution_payload_invalid"
  end
end
