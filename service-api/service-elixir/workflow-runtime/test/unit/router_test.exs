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
    assert Enum.any?(payload["dependencies"], &(&1["name"] == "workflow-catalog" and &1["status"] == "ready"))
  end

  test "capabilities route exposes runtime foundations" do
    conn = conn(:get, "/api/workflow-runtime/capabilities") |> Router.call([])
    payload = Jason.decode!(conn.resp_body)

    assert conn.status == 200
    assert payload["service"] == "workflow-runtime"
    assert payload["timers"]["enabled"] == true
    assert payload["timers"]["supportsDelayActions"] == true
    assert payload["retries"]["enabled"] == true
    assert payload["compensations"]["mode"] == "basic"
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
             "failed" => 0,
             "cancelled" => 0
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
    assert create_payload["tenantSlug"] == "bootstrap-ops"
    assert create_payload["workflowDefinitionKey"] == "lead-follow-up"
    assert create_payload["workflowDefinitionVersionNumber"] == 1
    assert create_payload["status"] == "pending"
    assert create_payload["currentActionIndex"] == 0
    assert create_payload["retryCount"] == 0
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

  test "execution create rejects missing workflow definition from catalog" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "missing-flow",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001118",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    payload = Jason.decode!(create_conn.resp_body)

    assert create_conn.status == 404
    assert payload["code"] == "workflow_runtime_definition_not_found"
  end

  test "execution can start and complete through runtime lifecycle" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001234",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    start_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/start") |> Router.call([])
    complete_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/complete") |> Router.call([])
    summary_conn = conn(:get, "/api/workflow-runtime/executions/summary") |> Router.call([])

    start_payload = Jason.decode!(start_conn.resp_body)
    complete_payload = Jason.decode!(complete_conn.resp_body)
    summary_payload = Jason.decode!(summary_conn.resp_body)

    assert start_conn.status == 200
    assert start_payload["status"] == "running"
    assert start_payload["startedAt"] != nil

    assert complete_conn.status == 200
    assert complete_payload["status"] == "completed"
    assert complete_payload["completedAt"] != nil

    assert summary_payload["completed"] == 1
  end

  test "execution transitions ledger reflects lifecycle" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001290",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    _start_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/start") |> Router.call([])
    _complete_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/complete") |> Router.call([])

    transitions_conn = conn(:get, "/api/workflow-runtime/executions/#{public_id}/transitions") |> Router.call([])
    transitions_payload = Jason.decode!(transitions_conn.resp_body)

    assert transitions_conn.status == 200
    assert Enum.map(transitions_payload, & &1["status"]) == ["pending", "running", "completed"]
    assert Enum.all?(transitions_payload, &(&1["changedBy"] == "unit-test"))
  end

  test "execution actions route exposes the runtime action ledger" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001289",
        "initiatedBy" => "action-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    _start_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/start") |> Router.call([])
    _advance_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/advance") |> Router.call([])

    actions_conn = conn(:get, "/api/workflow-runtime/executions/#{public_id}/actions") |> Router.call([])
    actions_payload = Jason.decode!(actions_conn.resp_body)

    assert actions_conn.status == 200
    assert length(actions_payload) == 1
    assert hd(actions_payload)["actionKey"] == "task.create"
    assert hd(actions_payload)["status"] == "completed"
  end

  test "execution list filters by tenant and status" do
    first_create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001291",
        "initiatedBy" => "dispatch-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    second_create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-west",
        "workflowDefinitionKey" => "deal-follow-up",
        "subjectType" => "crm.deal",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001292",
        "initiatedBy" => "dispatch-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    first_public_id = Jason.decode!(first_create_conn.resp_body)["publicId"]
    second_public_id = Jason.decode!(second_create_conn.resp_body)["publicId"]

    _start_conn = conn(:post, "/api/workflow-runtime/executions/#{second_public_id}/start") |> Router.call([])

    filtered_conn =
      conn(:get, "/api/workflow-runtime/executions?tenantSlug=ops-west&status=running") |> Router.call([])

    filtered_payload = Jason.decode!(filtered_conn.resp_body)

    assert filtered_conn.status == 200
    assert Enum.map(filtered_payload, & &1["publicId"]) == [second_public_id]
    assert Enum.all?(filtered_payload, &(&1["tenantSlug"] == "ops-west"))
    assert first_public_id not in Enum.map(filtered_payload, & &1["publicId"])
  end

  test "execution can be cancelled from runtime lifecycle" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001293",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    cancel_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/cancel") |> Router.call([])
    summary_conn = conn(:get, "/api/workflow-runtime/executions/summary") |> Router.call([])

    cancel_payload = Jason.decode!(cancel_conn.resp_body)
    summary_payload = Jason.decode!(summary_conn.resp_body)

    assert cancel_conn.status == 200
    assert cancel_payload["status"] == "cancelled"
    assert cancel_payload["cancelledAt"] != nil
    assert summary_payload["cancelled"] == 1
  end

  test "execution advance should wait on delay actions and finish after the delay" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "quote-follow-up",
        "subjectType" => "sales.quote",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001333",
        "initiatedBy" => "delay-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    _start_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/start") |> Router.call([])
    waiting_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/advance") |> Router.call([])
    waiting_payload = Jason.decode!(waiting_conn.resp_body)
    blocked_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/advance") |> Router.call([])
    blocked_payload = Jason.decode!(blocked_conn.resp_body)

    Process.sleep(1_100)

    release_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/advance") |> Router.call([])
    complete_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/advance") |> Router.call([])
    actions_conn = conn(:get, "/api/workflow-runtime/executions/#{public_id}/actions") |> Router.call([])

    release_payload = Jason.decode!(release_conn.resp_body)
    complete_payload = Jason.decode!(complete_conn.resp_body)
    actions_payload = Jason.decode!(actions_conn.resp_body)

    assert waiting_conn.status == 200
    assert waiting_payload["waitingUntil"] != nil

    assert blocked_conn.status == 409
    assert blocked_payload["code"] == "workflow_runtime_execution_waiting"

    assert release_conn.status == 200
    assert release_payload["waitingUntil"] == nil
    assert release_payload["currentActionIndex"] == 1

    assert complete_conn.status == 200
    assert complete_payload["status"] == "completed"
    assert complete_payload["completedAt"] != nil

    assert Enum.map(actions_payload, &{&1["actionKey"], &1["status"]}) == [
             {"delay.wait", "waiting"},
             {"delay.wait", "completed"},
             {"integration.webhook", "completed"}
           ]
  end

  test "execution summary filters by tenant and workflow definition" do
    completed_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001294",
        "initiatedBy" => "summary-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    failed_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "quote-follow-up",
        "subjectType" => "sales.quote",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001295",
        "initiatedBy" => "summary-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    completed_public_id = Jason.decode!(completed_conn.resp_body)["publicId"]
    failed_public_id = Jason.decode!(failed_conn.resp_body)["publicId"]

    _start_conn = conn(:post, "/api/workflow-runtime/executions/#{completed_public_id}/start") |> Router.call([])
    _complete_conn = conn(:post, "/api/workflow-runtime/executions/#{completed_public_id}/complete") |> Router.call([])
    _fail_conn = conn(:post, "/api/workflow-runtime/executions/#{failed_public_id}/fail") |> Router.call([])

    filtered_summary_conn =
      conn(:get, "/api/workflow-runtime/executions/summary?tenantSlug=ops-east&workflowDefinitionKey=lead-follow-up")
      |> Router.call([])

    filtered_summary_payload = Jason.decode!(filtered_summary_conn.resp_body)

    assert filtered_summary_conn.status == 200
    assert filtered_summary_payload["total"] == 1
    assert filtered_summary_payload["completed"] == 1
    assert filtered_summary_payload["failed"] == 0
  end

  test "execution summary by workflow groups runtime state and retries" do
    first_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001298",
        "initiatedBy" => "board-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    second_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "quote-follow-up",
        "subjectType" => "sales.quote",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001299",
        "initiatedBy" => "board-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    first_public_id = Jason.decode!(first_conn.resp_body)["publicId"]
    second_public_id = Jason.decode!(second_conn.resp_body)["publicId"]

    _start_first_conn = conn(:post, "/api/workflow-runtime/executions/#{first_public_id}/start") |> Router.call([])
    _complete_conn = conn(:post, "/api/workflow-runtime/executions/#{first_public_id}/complete") |> Router.call([])
    _fail_conn = conn(:post, "/api/workflow-runtime/executions/#{second_public_id}/fail") |> Router.call([])
    _retry_conn = conn(:post, "/api/workflow-runtime/executions/#{second_public_id}/retry") |> Router.call([])

    grouped_summary_conn =
      conn(:get, "/api/workflow-runtime/executions/summary/by-workflow?tenantSlug=ops-east") |> Router.call([])

    grouped_summary_payload = Jason.decode!(grouped_summary_conn.resp_body)

    lead_follow_up =
      Enum.find(grouped_summary_payload, &(&1["workflowDefinitionKey"] == "lead-follow-up"))

    quote_follow_up =
      Enum.find(grouped_summary_payload, &(&1["workflowDefinitionKey"] == "quote-follow-up"))

    assert grouped_summary_conn.status == 200
    assert lead_follow_up["total"] == 1
    assert lead_follow_up["completed"] == 1
    assert lead_follow_up["retriesTotal"] == 0
    assert quote_follow_up["total"] == 1
    assert quote_follow_up["pending"] == 1
    assert quote_follow_up["failed"] == 0
    assert quote_follow_up["retriesTotal"] == 1
  end

  test "execution can be retried after failure and resume lifecycle" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "tenantSlug" => "ops-east",
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001296",
        "initiatedBy" => "retry-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    _fail_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/fail") |> Router.call([])
    retry_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/retry") |> Router.call([])
    restart_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/start") |> Router.call([])
    complete_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/complete") |> Router.call([])
    transitions_conn = conn(:get, "/api/workflow-runtime/executions/#{public_id}/transitions") |> Router.call([])

    retry_payload = Jason.decode!(retry_conn.resp_body)
    restart_payload = Jason.decode!(restart_conn.resp_body)
    complete_payload = Jason.decode!(complete_conn.resp_body)
    transitions_payload = Jason.decode!(transitions_conn.resp_body)

    assert retry_conn.status == 200
    assert retry_payload["status"] == "pending"
    assert retry_payload["retryCount"] == 1
    assert retry_payload["failedAt"] == nil
    assert retry_payload["startedAt"] == nil
    assert retry_payload["completedAt"] == nil

    assert restart_conn.status == 200
    assert restart_payload["status"] == "running"

    assert complete_conn.status == 200
    assert complete_payload["status"] == "completed"
    assert complete_payload["retryCount"] == 1

    assert Enum.map(transitions_payload, & &1["status"]) == ["pending", "failed", "pending", "running", "completed"]
  end

  test "execution fail should append compensation ledger for completed actions" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001334",
        "initiatedBy" => "comp-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    _start_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/start") |> Router.call([])
    _advance_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/advance") |> Router.call([])
    fail_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/fail") |> Router.call([])
    actions_conn = conn(:get, "/api/workflow-runtime/executions/#{public_id}/actions") |> Router.call([])

    fail_payload = Jason.decode!(fail_conn.resp_body)
    actions_payload = Jason.decode!(actions_conn.resp_body)

    assert fail_conn.status == 200
    assert fail_payload["status"] == "failed"
    assert Enum.map(actions_payload, & &1["status"]) == ["completed", "compensated"]
    assert List.last(actions_payload)["actionKey"] == "task.create"
  end

  test "execution blocks retry when lifecycle has not failed or cancelled" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001297",
        "initiatedBy" => "retry-user"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    retry_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/retry") |> Router.call([])
    retry_payload = Jason.decode!(retry_conn.resp_body)

    assert retry_conn.status == 409
    assert retry_payload["code"] == "workflow_runtime_execution_retry_invalid"
  end

  test "execution blocks invalid lifecycle transition" do
    create_conn =
      conn(:post, "/api/workflow-runtime/executions", %{
        "workflowDefinitionKey" => "lead-follow-up",
        "subjectType" => "crm.lead",
        "subjectPublicId" => "00000000-0000-0000-0000-000000001235",
        "initiatedBy" => "unit-test"
      })
      |> Plug.Conn.put_req_header("content-type", "application/json")
      |> Router.call([])

    public_id = Jason.decode!(create_conn.resp_body)["publicId"]

    invalid_complete_conn = conn(:post, "/api/workflow-runtime/executions/#{public_id}/complete") |> Router.call([])
    invalid_payload = Jason.decode!(invalid_complete_conn.resp_body)

    assert invalid_complete_conn.status == 409
    assert invalid_payload["code"] == "workflow_runtime_execution_transition_invalid"
  end
end
