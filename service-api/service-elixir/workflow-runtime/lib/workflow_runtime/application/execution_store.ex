defmodule WorkflowRuntime.Application.ExecutionStore do
  @moduledoc false

  use Agent

  alias WorkflowRuntime.Infrastructure.Postgres

  def start_link(_opts) do
    if postgres_driver?() do
      :ignore
    else
      Agent.start_link(fn -> %{executions: [], transitions: %{}} end, name: __MODULE__)
    end
  end

  def list(filters \\ %{}) do
    if postgres_driver?() do
      list_postgres(filters)
    else
      list_memory(filters)
    end
  end

  def find(public_id) do
    if postgres_driver?() do
      find_postgres(public_id)
    else
      find_memory(public_id)
    end
  end

  def transitions(public_id) do
    if postgres_driver?() do
      transitions_postgres(public_id)
    else
      transitions_memory(public_id)
    end
  end

  def create(attributes) do
    if postgres_driver?() do
      create_postgres(attributes)
    else
      create_memory(attributes)
    end
  end

  def transition(public_id, next_status, allowed_statuses) do
    if postgres_driver?() do
      transition_postgres(public_id, next_status, allowed_statuses)
    else
      transition_memory(public_id, next_status, allowed_statuses)
    end
  end

  def retry(public_id, allowed_statuses) do
    if postgres_driver?() do
      retry_postgres(public_id, allowed_statuses)
    else
      retry_memory(public_id, allowed_statuses)
    end
  end

  def reset do
    unless postgres_driver?() do
      Agent.update(__MODULE__, fn _state -> %{executions: [], transitions: %{}} end)
    end
  end

  def summary(filters \\ %{}) do
    if postgres_driver?() do
      summary_postgres(filters)
    else
      summary_memory(filters)
    end
  end

  def summary_by_workflow(filters \\ %{}) do
    if postgres_driver?() do
      summary_by_workflow_postgres(filters)
    else
      summary_by_workflow_memory(filters)
    end
  end

  defp list_memory(filters) do
    filters = normalize_filters(filters)

    Agent.get(__MODULE__, fn state ->
      state.executions
      |> Enum.reverse()
      |> Enum.filter(&matches_filters?(&1, filters))
    end)
  end

  defp find_memory(public_id) do
    Agent.get(__MODULE__, fn state ->
      Enum.find(state.executions, &(&1["publicId"] == public_id))
    end)
  end

  defp transitions_memory(public_id) do
    Agent.get(__MODULE__, fn state ->
      Map.get(state.transitions, public_id, [])
    end)
  end

  defp create_memory(attributes) do
    initiated_by = String.trim(attributes["initiatedBy"] || "")

    execution = %{
      "publicId" => generate_public_id(),
      "tenantSlug" => normalize_tenant_slug(attributes["tenantSlug"]),
      "workflowDefinitionKey" => String.trim(attributes["workflowDefinitionKey"] || ""),
      "subjectType" => String.trim(attributes["subjectType"] || ""),
      "subjectPublicId" => String.trim(attributes["subjectPublicId"] || ""),
      "initiatedBy" => initiated_by,
      "status" => "pending",
      "retryCount" => 0,
      "createdAt" => now(),
      "startedAt" => nil,
      "completedAt" => nil,
      "failedAt" => nil,
      "cancelledAt" => nil
    }

    Agent.update(__MODULE__, fn state ->
      %{
        state
        | executions: [execution | state.executions],
          transitions: Map.put(state.transitions, execution["publicId"], [build_transition("pending", initiated_by)])
      }
    end)

    execution
  end

  defp transition_memory(public_id, next_status, allowed_statuses) do
    Agent.get_and_update(__MODULE__, fn state ->
      case Enum.find(state.executions, &(&1["publicId"] == public_id)) do
        nil ->
          {{:error, :not_found}, state}

        execution ->
          if execution["status"] in allowed_statuses do
            updated_execution =
              execution
              |> Map.put("status", next_status)
              |> maybe_put_timestamp(next_status)

            updated_transitions =
              Map.update(
                state.transitions,
                public_id,
                [build_transition(next_status, execution["initiatedBy"])],
                fn existing ->
                  existing ++ [build_transition(next_status, execution["initiatedBy"])]
                end
              )

            updated_state = %{
              state
              | executions:
                  Enum.map(state.executions, fn current_execution ->
                    if current_execution["publicId"] == public_id, do: updated_execution, else: current_execution
                  end),
                transitions: updated_transitions
            }

            {{:ok, updated_execution}, updated_state}
          else
            {{:error, :invalid_transition}, state}
          end
      end
    end)
  end

  defp summary_memory(filters) do
    executions = list(filters)

    %{
      total: length(executions),
      pending: Enum.count(executions, &(&1["status"] == "pending")),
      running: Enum.count(executions, &(&1["status"] == "running")),
      completed: Enum.count(executions, &(&1["status"] == "completed")),
      failed: Enum.count(executions, &(&1["status"] == "failed")),
      cancelled: Enum.count(executions, &(&1["status"] == "cancelled"))
    }
  end

  defp retry_memory(public_id, allowed_statuses) do
    Agent.get_and_update(__MODULE__, fn state ->
      case Enum.find(state.executions, &(&1["publicId"] == public_id)) do
        nil ->
          {{:error, :not_found}, state}

        execution ->
          if execution["status"] in allowed_statuses do
            updated_execution =
              execution
              |> reset_for_retry()

            updated_transitions =
              Map.update(
                state.transitions,
                public_id,
                [build_transition("pending", execution["initiatedBy"])],
                fn existing ->
                  existing ++ [build_transition("pending", execution["initiatedBy"])]
                end
              )

            updated_state = %{
              state
              | executions:
                  Enum.map(state.executions, fn current_execution ->
                    if current_execution["publicId"] == public_id, do: updated_execution, else: current_execution
                  end),
                transitions: updated_transitions
            }

            {{:ok, updated_execution}, updated_state}
          else
            {{:error, :invalid_transition}, state}
          end
      end
    end)
  end

  defp summary_by_workflow_memory(filters) do
    filters
    |> list()
    |> Enum.group_by(& &1["workflowDefinitionKey"])
    |> Enum.sort_by(fn {workflow_definition_key, _executions} -> workflow_definition_key end)
    |> Enum.map(fn {workflow_definition_key, executions} ->
      %{
        "workflowDefinitionKey" => workflow_definition_key,
        "total" => length(executions),
        "pending" => Enum.count(executions, &(&1["status"] == "pending")),
        "running" => Enum.count(executions, &(&1["status"] == "running")),
        "completed" => Enum.count(executions, &(&1["status"] == "completed")),
        "failed" => Enum.count(executions, &(&1["status"] == "failed")),
        "cancelled" => Enum.count(executions, &(&1["status"] == "cancelled")),
        "retriesTotal" => Enum.reduce(executions, 0, &(&1["retryCount"] + &2))
      }
    end)
  end

  defp list_postgres(filters) do
    filters = normalize_filters(filters)
    {conditions, params} = postgres_filter_conditions(filters)

    statement = """
    SELECT
      execution.public_id::text,
      tenant.slug,
      execution.workflow_definition_key,
      execution.subject_type,
      execution.subject_public_id::text,
      execution.initiated_by,
      execution.status,
      execution.retry_count,
      execution.created_at,
      execution.started_at,
      execution.completed_at,
      execution.failed_at,
      execution.cancelled_at
    FROM workflow_runtime.executions AS execution
    JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
    #{build_where_clause(conditions)}
    ORDER BY execution.id ASC
    """

    Postgres.query!(statement, params).rows
    |> Enum.map(&map_execution_row/1)
  end

  defp find_postgres(public_id) do
    statement = """
    SELECT
      execution.public_id::text,
      tenant.slug,
      execution.workflow_definition_key,
      execution.subject_type,
      execution.subject_public_id::text,
      execution.initiated_by,
      execution.status,
      execution.retry_count,
      execution.created_at,
      execution.started_at,
      execution.completed_at,
      execution.failed_at,
      execution.cancelled_at
    FROM workflow_runtime.executions AS execution
    JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
    WHERE execution.public_id = $1::uuid
    LIMIT 1
    """

    case Postgres.query!(statement, [uuid_bytes(public_id)]).rows do
      [row] -> map_execution_row(row)
      _ -> nil
    end
  end

  defp transitions_postgres(public_id) do
    statement = """
    SELECT
      transition.public_id::text,
      transition.status,
      transition.changed_by,
      transition.created_at
    FROM workflow_runtime.execution_transitions AS transition
    JOIN workflow_runtime.executions AS execution ON execution.id = transition.execution_id
    WHERE execution.public_id = $1::uuid
    ORDER BY transition.id ASC
    """

    Postgres.query!(statement, [uuid_bytes(public_id)]).rows
    |> Enum.map(&map_transition_row/1)
  end

  defp create_postgres(attributes) do
    tenant_slug = normalize_tenant_slug(attributes["tenantSlug"])
    initiated_by = String.trim(attributes["initiatedBy"] || "")

    case fetch_tenant_id(tenant_slug) do
      nil ->
        {:error, :tenant_not_found}

      tenant_id ->
        public_id = generate_public_id()
        transition_public_id = generate_public_id()

        {:ok, :created} =
          Postgres.transaction(fn connection ->
            Postgrex.query!(
              connection,
              """
              INSERT INTO workflow_runtime.executions (
                tenant_id,
                public_id,
                workflow_definition_key,
                subject_type,
                subject_public_id,
                initiated_by,
                status,
                retry_count
              )
              VALUES ($1, $2::uuid, $3, $4, $5::uuid, $6, 'pending', 0)
              """,
              [
                tenant_id,
                uuid_bytes(public_id),
                String.trim(attributes["workflowDefinitionKey"] || ""),
                String.trim(attributes["subjectType"] || ""),
                uuid_bytes(String.trim(attributes["subjectPublicId"] || "")),
                initiated_by
              ]
            )

            Postgrex.query!(
              connection,
              """
              INSERT INTO workflow_runtime.execution_transitions (
                public_id,
                execution_id,
                status,
                changed_by
              )
              SELECT $1::uuid, execution.id, 'pending', $3
              FROM workflow_runtime.executions AS execution
              WHERE execution.public_id = $2::uuid
              """,
              [uuid_bytes(transition_public_id), uuid_bytes(public_id), initiated_by]
            )

            :created
          end)

        find_postgres(public_id)
    end
  end

  defp transition_postgres(public_id, next_status, allowed_statuses) do
    case find_postgres(public_id) do
      nil ->
        {:error, :not_found}

      execution ->
        if execution["status"] in allowed_statuses do
          transition_public_id = generate_public_id()
          timestamp_update = timestamp_update_clause(next_status)

          {:ok, :updated} =
            Postgres.transaction(fn connection ->
              Postgrex.query!(
                connection,
                """
                UPDATE workflow_runtime.executions
                SET
                  status = $2,
                  updated_at = timezone('utc', now())
                  #{timestamp_update}
                WHERE public_id = $1::uuid
                """,
                [uuid_bytes(public_id), next_status]
              )

              Postgrex.query!(
                connection,
                """
                INSERT INTO workflow_runtime.execution_transitions (
                  public_id,
                  execution_id,
                  status,
                  changed_by
                )
                SELECT $1::uuid, runtime_execution.id, $3, runtime_execution.initiated_by
                FROM workflow_runtime.executions AS runtime_execution
                WHERE runtime_execution.public_id = $2::uuid
                """,
                [uuid_bytes(transition_public_id), uuid_bytes(public_id), next_status]
              )

              :updated
            end)

          {:ok, find_postgres(public_id)}
        else
          {:error, :invalid_transition}
        end
    end
  end

  defp summary_postgres(filters) do
    filters = normalize_filters(filters)
    {conditions, params} = postgres_filter_conditions(filters)

    statement = """
    SELECT
      count(*) AS total,
      count(*) FILTER (WHERE execution.status = 'pending') AS pending,
      count(*) FILTER (WHERE execution.status = 'running') AS running,
      count(*) FILTER (WHERE execution.status = 'completed') AS completed,
      count(*) FILTER (WHERE execution.status = 'failed') AS failed,
      count(*) FILTER (WHERE execution.status = 'cancelled') AS cancelled
    FROM workflow_runtime.executions AS execution
    JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
    #{build_where_clause(conditions)}
    """

    [row] = Postgres.query!(statement, params).rows

    %{
      total: Enum.at(row, 0),
      pending: Enum.at(row, 1),
      running: Enum.at(row, 2),
      completed: Enum.at(row, 3),
      failed: Enum.at(row, 4),
      cancelled: Enum.at(row, 5)
    }
  end

  defp retry_postgres(public_id, allowed_statuses) do
    case find_postgres(public_id) do
      nil ->
        {:error, :not_found}

      execution ->
        if execution["status"] in allowed_statuses do
          transition_public_id = generate_public_id()

          {:ok, :updated} =
            Postgres.transaction(fn connection ->
              Postgrex.query!(
                connection,
                """
                UPDATE workflow_runtime.executions
                SET
                  status = 'pending',
                  retry_count = retry_count + 1,
                  started_at = NULL,
                  completed_at = NULL,
                  failed_at = NULL,
                  cancelled_at = NULL,
                  updated_at = timezone('utc', now())
                WHERE public_id = $1::uuid
                """,
                [uuid_bytes(public_id)]
              )

              Postgrex.query!(
                connection,
                """
                INSERT INTO workflow_runtime.execution_transitions (
                  public_id,
                  execution_id,
                  status,
                  changed_by
                )
                SELECT $1::uuid, runtime_execution.id, 'pending', runtime_execution.initiated_by
                FROM workflow_runtime.executions AS runtime_execution
                WHERE runtime_execution.public_id = $2::uuid
                """,
                [uuid_bytes(transition_public_id), uuid_bytes(public_id)]
              )

              :updated
            end)

          {:ok, find_postgres(public_id)}
        else
          {:error, :invalid_transition}
        end
    end
  end

  defp summary_by_workflow_postgres(filters) do
    filters = normalize_filters(filters)
    {conditions, params} = postgres_filter_conditions(filters)

    statement = """
    SELECT
      execution.workflow_definition_key,
      count(*) AS total,
      count(*) FILTER (WHERE execution.status = 'pending') AS pending,
      count(*) FILTER (WHERE execution.status = 'running') AS running,
      count(*) FILTER (WHERE execution.status = 'completed') AS completed,
      count(*) FILTER (WHERE execution.status = 'failed') AS failed,
      count(*) FILTER (WHERE execution.status = 'cancelled') AS cancelled,
      coalesce(sum(execution.retry_count), 0) AS retries_total
    FROM workflow_runtime.executions AS execution
    JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
    #{build_where_clause(conditions)}
    GROUP BY execution.workflow_definition_key
    ORDER BY execution.workflow_definition_key ASC
    """

    Postgres.query!(statement, params).rows
    |> Enum.map(&map_workflow_summary_row/1)
  end

  defp postgres_driver? do
    Postgres.enabled?()
  end

  defp fetch_tenant_id(tenant_slug) do
    case Postgres.query!(
           "SELECT id FROM identity.tenants WHERE slug = $1 LIMIT 1",
           [tenant_slug]
         ).rows do
      [[tenant_id]] -> tenant_id
      _ -> nil
    end
  end

  defp postgres_filter_conditions(filters) do
    filters
    |> Enum.reduce({[], []}, fn
      {"tenantSlug", nil}, acc -> acc
      {"status", nil}, acc -> acc
      {"workflowDefinitionKey", nil}, acc -> acc
      {"subjectType", nil}, acc -> acc
      {"subjectPublicId", nil}, acc -> acc
      {"initiatedBy", nil}, acc -> acc
      {"tenantSlug", value}, {conditions, params} -> {[conditions | ["tenant.slug = $" <> Integer.to_string(length(params) + 1)]], params ++ [value]}
      {"status", value}, {conditions, params} -> {[conditions | ["execution.status = $" <> Integer.to_string(length(params) + 1)]], params ++ [value]}
      {"workflowDefinitionKey", value}, {conditions, params} -> {[conditions | ["execution.workflow_definition_key = $" <> Integer.to_string(length(params) + 1)]], params ++ [value]}
      {"subjectType", value}, {conditions, params} -> {[conditions | ["execution.subject_type = $" <> Integer.to_string(length(params) + 1)]], params ++ [value]}
      {"subjectPublicId", value}, {conditions, params} -> {[conditions | ["execution.subject_public_id = $" <> Integer.to_string(length(params) + 1)]], params ++ [uuid_bytes(value)]}
      {"initiatedBy", value}, {conditions, params} -> {[conditions | ["execution.initiated_by = $" <> Integer.to_string(length(params) + 1)]], params ++ [value]}
    end)
    |> then(fn {conditions, params} -> {List.flatten(conditions), params} end)
  end

  defp build_where_clause([]), do: ""
  defp build_where_clause(conditions), do: "WHERE " <> Enum.join(conditions, " AND ")

  defp timestamp_update_clause("running"), do: ", started_at = timezone('utc', now())"
  defp timestamp_update_clause("completed"), do: ", completed_at = timezone('utc', now())"
  defp timestamp_update_clause("failed"), do: ", failed_at = timezone('utc', now())"
  defp timestamp_update_clause("cancelled"), do: ", cancelled_at = timezone('utc', now())"
  defp timestamp_update_clause(_next_status), do: ""

  defp map_execution_row(row, tenant_slug_override \\ nil) do
    %{
      "publicId" => Enum.at(row, 0),
      "tenantSlug" => tenant_slug_override || Enum.at(row, 1),
      "workflowDefinitionKey" => Enum.at(row, 2),
      "subjectType" => Enum.at(row, 3),
      "subjectPublicId" => Enum.at(row, 4),
      "initiatedBy" => Enum.at(row, 5),
      "status" => Enum.at(row, 6),
      "retryCount" => Enum.at(row, 7),
      "createdAt" => encode_datetime(Enum.at(row, 8)),
      "startedAt" => encode_datetime(Enum.at(row, 9)),
      "completedAt" => encode_datetime(Enum.at(row, 10)),
      "failedAt" => encode_datetime(Enum.at(row, 11)),
      "cancelledAt" => encode_datetime(Enum.at(row, 12))
    }
  end

  defp map_transition_row(row) do
    %{
      "publicId" => Enum.at(row, 0),
      "status" => Enum.at(row, 1),
      "changedBy" => Enum.at(row, 2),
      "createdAt" => encode_datetime(Enum.at(row, 3))
    }
  end

  defp map_workflow_summary_row(row) do
    %{
      "workflowDefinitionKey" => Enum.at(row, 0),
      "total" => Enum.at(row, 1),
      "pending" => Enum.at(row, 2),
      "running" => Enum.at(row, 3),
      "completed" => Enum.at(row, 4),
      "failed" => Enum.at(row, 5),
      "cancelled" => Enum.at(row, 6),
      "retriesTotal" => Enum.at(row, 7)
    }
  end

  defp encode_datetime(nil), do: nil
  defp encode_datetime(value), do: DateTime.to_iso8601(value)

  defp uuid_bytes(value) do
    value
    |> String.replace("-", "")
    |> Base.decode16!(case: :mixed)
  end

  defp generate_public_id do
    <<part1::binary-size(4), part2::binary-size(2), part3::binary-size(2), part4::binary-size(2), part5::binary-size(6)>> =
      :crypto.strong_rand_bytes(16)

    [
      Base.encode16(part1, case: :lower),
      Base.encode16(part2, case: :lower),
      Base.encode16(part3, case: :lower),
      Base.encode16(part4, case: :lower),
      Base.encode16(part5, case: :lower)
    ]
    |> Enum.join("-")
  end

  defp maybe_put_timestamp(execution, "running"), do: Map.put(execution, "startedAt", now())
  defp maybe_put_timestamp(execution, "completed"), do: Map.put(execution, "completedAt", now())
  defp maybe_put_timestamp(execution, "failed"), do: Map.put(execution, "failedAt", now())
  defp maybe_put_timestamp(execution, "cancelled"), do: Map.put(execution, "cancelledAt", now())
  defp maybe_put_timestamp(execution, _next_status), do: execution

  defp reset_for_retry(execution) do
    execution
    |> Map.put("status", "pending")
    |> Map.update!("retryCount", &(&1 + 1))
    |> Map.put("startedAt", nil)
    |> Map.put("completedAt", nil)
    |> Map.put("failedAt", nil)
    |> Map.put("cancelledAt", nil)
  end

  defp normalize_filters(filters) do
    %{
      "tenantSlug" => normalize_filter(filters["tenantSlug"]),
      "status" => normalize_filter(filters["status"]),
      "workflowDefinitionKey" => normalize_filter(filters["workflowDefinitionKey"]),
      "subjectType" => normalize_filter(filters["subjectType"]),
      "subjectPublicId" => normalize_filter(filters["subjectPublicId"]),
      "initiatedBy" => normalize_filter(filters["initiatedBy"])
    }
  end

  defp matches_filters?(execution, filters) do
    Enum.all?(filters, fn
      {_key, nil} -> true
      {key, value} -> execution[key] == value
    end)
  end

  defp normalize_filter(nil), do: nil
  defp normalize_filter(""), do: nil
  defp normalize_filter(value), do: String.trim(value)

  defp normalize_tenant_slug(nil), do: default_tenant_slug()

  defp normalize_tenant_slug(value) do
    case String.trim(value) do
      "" -> default_tenant_slug()
      tenant_slug -> tenant_slug
    end
  end

  defp default_tenant_slug do
    Application.fetch_env!(:workflow_runtime, :bootstrap_tenant_slug)
  end

  defp build_transition(status, changed_by) do
    %{
      "status" => status,
      "changedBy" => changed_by,
      "createdAt" => now()
    }
  end

  defp now do
    DateTime.utc_now() |> DateTime.truncate(:second) |> DateTime.to_iso8601()
  end
end
