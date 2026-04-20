defmodule WorkflowRuntime.Application.ExecutionStore do
  @moduledoc false

  use Agent

  alias WorkflowRuntime.Infrastructure.Postgres

  def start_link(_opts) do
    if postgres_driver?() do
      :ignore
    else
      Agent.start_link(fn -> %{executions: [], transitions: %{}, actions: %{}, plans: %{}} end, name: __MODULE__)
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

  def actions(public_id) do
    if postgres_driver?() do
      actions_postgres(public_id)
    else
      actions_memory(public_id)
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

  def advance(public_id) do
    if postgres_driver?() do
      advance_postgres(public_id)
    else
      advance_memory(public_id)
    end
  end

  def reset do
    unless postgres_driver?() do
      Agent.update(__MODULE__, fn _state -> %{executions: [], transitions: %{}, actions: %{}, plans: %{}} end)
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

  def catalog_status do
    if postgres_driver?() do
      catalog_status_postgres()
    else
      :ready
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

  defp actions_memory(public_id) do
    Agent.get(__MODULE__, fn state ->
      Map.get(state.actions, public_id, [])
    end)
  end

  defp create_memory(attributes) do
    initiated_by = String.trim(attributes["initiatedBy"] || "")
    workflow_definition_key = String.trim(attributes["workflowDefinitionKey"] || "")

    case fetch_catalog_plan_memory(workflow_definition_key) do
      {:ok, workflow_definition_version_number, plan} ->
        public_id = generate_public_id()

        execution = %{
          "publicId" => public_id,
          "tenantSlug" => normalize_tenant_slug(attributes["tenantSlug"]),
          "workflowDefinitionKey" => workflow_definition_key,
          "workflowDefinitionVersionNumber" => workflow_definition_version_number,
          "subjectType" => String.trim(attributes["subjectType"] || ""),
          "subjectPublicId" => String.trim(attributes["subjectPublicId"] || ""),
          "initiatedBy" => initiated_by,
          "status" => "pending",
          "currentActionIndex" => 0,
          "waitingUntil" => nil,
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
              transitions: Map.put(state.transitions, public_id, [build_transition("pending", initiated_by)]),
              actions: Map.put(state.actions, public_id, []),
              plans: Map.put(state.plans, public_id, plan)
          }
        end)

        execution

      error ->
        {:error, error}
    end
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

            updated_actions =
              if next_status in ["failed", "cancelled"] do
                append_compensations_memory(
                  Map.get(state.actions, public_id, []),
                  execution["initiatedBy"]
                )
              else
                Map.get(state.actions, public_id, [])
              end

            updated_state = %{
              state
              | executions:
                  Enum.map(state.executions, fn current_execution ->
                    if current_execution["publicId"] == public_id, do: updated_execution, else: current_execution
                  end),
                transitions: updated_transitions,
                actions: Map.put(state.actions, public_id, updated_actions)
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

  defp advance_memory(public_id) do
    Agent.get_and_update(__MODULE__, fn state ->
      case Enum.find(state.executions, &(&1["publicId"] == public_id)) do
        nil ->
          {{:error, :not_found}, state}

        execution ->
          if execution["status"] != "running" do
            {{:error, :invalid_transition}, state}
          else
            plan = Map.get(state.plans, public_id, [])
            current_action_index = execution["currentActionIndex"]

            if current_action_index >= length(plan) do
              completed_execution =
                execution
                |> Map.put("status", "completed")
                |> maybe_put_timestamp("completed")

              updated_transitions =
                Map.update(
                  state.transitions,
                  public_id,
                  [build_transition("completed", execution["initiatedBy"])],
                  fn existing ->
                    existing ++ [build_transition("completed", execution["initiatedBy"])]
                  end
                )

              updated_state = %{
                state
                | executions:
                    Enum.map(state.executions, fn current_execution ->
                      if current_execution["publicId"] == public_id, do: completed_execution, else: current_execution
                    end),
                  transitions: updated_transitions
              }

              {{:ok, completed_execution}, updated_state}
            else
              action = Enum.at(plan, current_action_index)

              case advance_execution_action(execution, action, length(plan)) do
                {:waiting, updated_execution, action_entry} ->
                  updated_state = %{
                    state
                    | executions:
                        Enum.map(state.executions, fn current_execution ->
                          if current_execution["publicId"] == public_id, do: updated_execution, else: current_execution
                        end),
                      actions: Map.update(state.actions, public_id, [action_entry], fn existing -> existing ++ [action_entry] end)
                  }

                  {{:ok, updated_execution}, updated_state}

                {:completed, updated_execution, action_entry, transition_status} ->
                  updated_transitions =
                    if transition_status == nil do
                      state.transitions
                    else
                      Map.update(
                        state.transitions,
                        public_id,
                        [build_transition(transition_status, execution["initiatedBy"])],
                        fn existing ->
                          existing ++ [build_transition(transition_status, execution["initiatedBy"])]
                        end
                      )
                    end

                  updated_state = %{
                    state
                    | executions:
                        Enum.map(state.executions, fn current_execution ->
                          if current_execution["publicId"] == public_id, do: updated_execution, else: current_execution
                        end),
                      transitions: updated_transitions,
                      actions: Map.update(state.actions, public_id, [action_entry], fn existing -> existing ++ [action_entry] end)
                  }

                  {{:ok, updated_execution}, updated_state}

                {:error, :waiting} ->
                  {{:error, :waiting}, state}
              end
            end
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
      execution.workflow_definition_version_number,
      execution.subject_type,
      execution.subject_public_id::text,
      execution.initiated_by,
      execution.status,
      execution.current_action_index,
      execution.waiting_until,
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
      execution.workflow_definition_version_number,
      execution.subject_type,
      execution.subject_public_id::text,
      execution.initiated_by,
      execution.status,
      execution.current_action_index,
      execution.waiting_until,
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

  defp actions_postgres(public_id) do
    statement = """
    SELECT
      action.public_id::text,
      action.step_id,
      action.action_key,
      action.label,
      action.status,
      action.delay_seconds,
      action.compensation_action_key,
      action.created_at
    FROM workflow_runtime.execution_actions AS action
    JOIN workflow_runtime.executions AS execution ON execution.id = action.execution_id
    WHERE execution.public_id = $1::uuid
    ORDER BY action.id ASC
    """

    Postgres.query!(statement, [uuid_bytes(public_id)]).rows
    |> Enum.map(&map_action_row/1)
  end

  defp create_postgres(attributes) do
    tenant_slug = normalize_tenant_slug(attributes["tenantSlug"])
    initiated_by = String.trim(attributes["initiatedBy"] || "")
    workflow_definition_key = String.trim(attributes["workflowDefinitionKey"] || "")

    case fetch_tenant_id(tenant_slug) do
      nil ->
        {:error, :tenant_not_found}

      tenant_id ->
        case fetch_catalog_plan_postgres(tenant_id, workflow_definition_key) do
          {:ok, workflow_definition_version_number, _plan} ->
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
                    workflow_definition_version_number,
                    subject_type,
                    subject_public_id,
                    initiated_by,
                    status,
                    current_action_index,
                    retry_count
                  )
                  VALUES ($1, $2::uuid, $3, $4, $5, $6::uuid, $7, 'pending', 0, 0)
                  """,
                  [
                    tenant_id,
                    uuid_bytes(public_id),
                    workflow_definition_key,
                    workflow_definition_version_number,
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

          error ->
            {:error, error}
        end
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

              if next_status in ["failed", "cancelled"] do
                Postgrex.query!(
                  connection,
                  """
                  INSERT INTO workflow_runtime.execution_actions (
                    public_id,
                    execution_id,
                    step_id,
                    action_key,
                    label,
                    status,
                    delay_seconds,
                    compensation_action_key
                  )
                  SELECT
                    gen_random_uuid(),
                    runtime_execution.id,
                    completed_action.step_id || ':compensation',
                    completed_action.compensation_action_key,
                    'Compensate: ' || completed_action.label,
                    'compensated',
                    NULL,
                    NULL
                  FROM workflow_runtime.executions AS runtime_execution
                  JOIN workflow_runtime.execution_actions AS completed_action ON completed_action.execution_id = runtime_execution.id
                  WHERE runtime_execution.public_id = $1::uuid
                    AND completed_action.status = 'completed'
                    AND completed_action.compensation_action_key IS NOT NULL
                  ORDER BY completed_action.id DESC
                  """,
                  [uuid_bytes(public_id)]
                )
              end

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
                  current_action_index = 0,
                  waiting_until = NULL,
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

  defp advance_postgres(public_id) do
    case find_postgres(public_id) do
      nil ->
        {:error, :not_found}

      execution ->
        if execution["status"] != "running" do
          {:error, :invalid_transition}
        else
          case fetch_execution_plan_postgres(execution) do
            {:error, reason} ->
              {:error, reason}

            {:ok, plan} ->
              current_action_index = execution["currentActionIndex"]

              if current_action_index >= length(plan) do
                transition_postgres(public_id, "completed", ["running"])
              else
                action = Enum.at(plan, current_action_index)

                case advance_execution_action(execution, action, length(plan)) do
                  {:waiting, updated_execution, action_entry} ->
                    {:ok, :updated} =
                      Postgres.transaction(fn connection ->
                        Postgrex.query!(
                          connection,
                          """
                          UPDATE workflow_runtime.executions
                          SET
                            waiting_until = $2::timestamptz,
                            updated_at = timezone('utc', now())
                          WHERE public_id = $1::uuid
                          """,
                          [uuid_bytes(public_id), to_datetime_param(updated_execution["waitingUntil"])]
                        )

                        Postgrex.query!(
                          connection,
                          """
                          INSERT INTO workflow_runtime.execution_actions (
                            public_id,
                            execution_id,
                            step_id,
                            action_key,
                            label,
                            status,
                            delay_seconds,
                            compensation_action_key
                          )
                          SELECT $1::uuid, runtime_execution.id, $3, $4, $5, $6, $7, $8
                          FROM workflow_runtime.executions AS runtime_execution
                          WHERE runtime_execution.public_id = $2::uuid
                          """,
                          [
                            uuid_bytes(action_entry["publicId"]),
                            uuid_bytes(public_id),
                            action_entry["stepId"],
                            action_entry["actionKey"],
                            action_entry["label"],
                            action_entry["status"],
                            action_entry["delaySeconds"],
                            action_entry["compensationActionKey"]
                          ]
                        )

                        :updated
                      end)

                    {:ok, find_postgres(public_id)}

                  {:completed, updated_execution, action_entry, transition_status} ->
                    {:ok, :updated} =
                      Postgres.transaction(fn connection ->
                        Postgrex.query!(
                          connection,
                          """
                          UPDATE workflow_runtime.executions
                          SET
                            current_action_index = $2,
                            waiting_until = $3::timestamptz,
                            status = $4::varchar,
                            completed_at = CASE WHEN $4::varchar = 'completed' THEN timezone('utc', now()) ELSE completed_at END,
                            updated_at = timezone('utc', now())
                          WHERE public_id = $1::uuid
                          """,
                          [
                            uuid_bytes(public_id),
                            updated_execution["currentActionIndex"],
                            to_datetime_param(updated_execution["waitingUntil"]),
                            updated_execution["status"]
                          ]
                        )

                        Postgrex.query!(
                          connection,
                          """
                          INSERT INTO workflow_runtime.execution_actions (
                            public_id,
                            execution_id,
                            step_id,
                            action_key,
                            label,
                            status,
                            delay_seconds,
                            compensation_action_key
                          )
                          SELECT $1::uuid, runtime_execution.id, $3, $4, $5, $6, $7, $8
                          FROM workflow_runtime.executions AS runtime_execution
                          WHERE runtime_execution.public_id = $2::uuid
                          """,
                          [
                            uuid_bytes(action_entry["publicId"]),
                            uuid_bytes(public_id),
                            action_entry["stepId"],
                            action_entry["actionKey"],
                            action_entry["label"],
                            action_entry["status"],
                            action_entry["delaySeconds"],
                            action_entry["compensationActionKey"]
                          ]
                        )

                        if transition_status != nil do
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
                            [uuid_bytes(generate_public_id()), uuid_bytes(public_id), transition_status]
                          )
                        end

                        :updated
                      end)

                    {:ok, find_postgres(public_id)}

                  {:error, :waiting} ->
                    {:error, :waiting}
                end
              end
          end
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

  defp catalog_status_postgres do
    case Postgres.query!(
           """
           SELECT count(*)
           FROM workflow_control.workflow_definitions AS definition
           JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
           WHERE tenant.slug = $1
           """,
           [default_tenant_slug()]
         ).rows do
      [[count]] when count >= 0 -> :ready
      _ -> :not_ready
    end
  rescue
    _ -> :not_ready
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

  defp fetch_catalog_plan_memory(workflow_definition_key) do
    catalog = %{
      "lead-follow-up" => {1,
        [
          %{"stepId" => "create-task", "actionKey" => "task.create", "label" => "Criar tarefa comercial inicial", "delaySeconds" => nil, "compensationActionKey" => "task.create"},
          %{"stepId" => "cooldown", "actionKey" => "delay.wait", "label" => "Aguardar janela curta de acompanhamento", "delaySeconds" => 1, "compensationActionKey" => nil},
          %{"stepId" => "notify-webhook", "actionKey" => "integration.webhook", "label" => "Emitir webhook operacional", "delaySeconds" => nil, "compensationActionKey" => "integration.webhook"}
        ]},
      "deal-follow-up" => {1,
        [
          %{"stepId" => "advance-stage", "actionKey" => "sales.stage.advance", "label" => "Avancar etapa comercial", "delaySeconds" => nil, "compensationActionKey" => "sales.stage.advance"}
        ]},
      "quote-follow-up" => {1,
        [
          %{"stepId" => "wait-touchpoint", "actionKey" => "delay.wait", "label" => "Aguardar janela de retorno", "delaySeconds" => 1, "compensationActionKey" => nil},
          %{"stepId" => "notify-quote", "actionKey" => "integration.webhook", "label" => "Notificar fluxo comercial", "delaySeconds" => nil, "compensationActionKey" => "integration.webhook"}
        ]}
    }

    case Map.fetch(catalog, workflow_definition_key) do
      {:ok, {workflow_definition_version_number, plan}} -> {:ok, workflow_definition_version_number, plan}
      :error -> :workflow_definition_not_found
    end
  end

  defp fetch_catalog_plan_postgres(tenant_id, workflow_definition_key) do
    case Postgres.query!(
           """
           SELECT
             definition.status,
             version.version_number,
             version.snapshot_status,
             version.snapshot_actions
           FROM workflow_control.workflow_definitions AS definition
           LEFT JOIN LATERAL (
             SELECT
               workflow_version.version_number,
               workflow_version.snapshot_status,
               workflow_version.snapshot_actions
             FROM workflow_control.workflow_definition_versions AS workflow_version
             WHERE workflow_version.workflow_definition_id = definition.id
             ORDER BY workflow_version.version_number DESC
             LIMIT 1
           ) AS version ON true
           WHERE definition.tenant_id = $1
             AND definition.key = $2
           LIMIT 1
           """,
           [tenant_id, workflow_definition_key]
         ).rows do
      [] ->
        :workflow_definition_not_found

      [[status, nil, _latest_snapshot_status, _snapshot_actions]] when status in ["draft", "active", "archived"] ->
        :workflow_definition_version_not_found

      [[status, _version_number, _latest_snapshot_status, _snapshot_actions]] when status != "active" ->
        :workflow_definition_inactive

      [[_status, _version_number, latest_snapshot_status, _snapshot_actions]] when latest_snapshot_status != "active" ->
        :workflow_definition_inactive

      [[_status, version_number, _latest_snapshot_status, snapshot_actions]] ->
        {:ok, version_number, Enum.map(snapshot_actions || [], &normalize_action_map/1)}
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
      "workflowDefinitionVersionNumber" => Enum.at(row, 3),
      "subjectType" => Enum.at(row, 4),
      "subjectPublicId" => Enum.at(row, 5),
      "initiatedBy" => Enum.at(row, 6),
      "status" => Enum.at(row, 7),
      "currentActionIndex" => Enum.at(row, 8),
      "waitingUntil" => encode_datetime(Enum.at(row, 9)),
      "retryCount" => Enum.at(row, 10),
      "createdAt" => encode_datetime(Enum.at(row, 11)),
      "startedAt" => encode_datetime(Enum.at(row, 12)),
      "completedAt" => encode_datetime(Enum.at(row, 13)),
      "failedAt" => encode_datetime(Enum.at(row, 14)),
      "cancelledAt" => encode_datetime(Enum.at(row, 15))
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

  defp map_action_row(row) do
    %{
      "publicId" => Enum.at(row, 0),
      "stepId" => Enum.at(row, 1),
      "actionKey" => Enum.at(row, 2),
      "label" => Enum.at(row, 3),
      "status" => Enum.at(row, 4),
      "delaySeconds" => Enum.at(row, 5),
      "compensationActionKey" => Enum.at(row, 6),
      "createdAt" => encode_datetime(Enum.at(row, 7))
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
    |> Map.put("currentActionIndex", 0)
    |> Map.put("waitingUntil", nil)
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

  defp build_action_entry(action, status) do
    %{
      "publicId" => generate_public_id(),
      "stepId" => action["stepId"],
      "actionKey" => action["actionKey"],
      "label" => action["label"],
      "status" => status,
      "delaySeconds" => action["delaySeconds"],
      "compensationActionKey" => action["compensationActionKey"],
      "createdAt" => now()
    }
  end

  defp append_compensations_memory(action_entries, initiated_by) do
    compensation_entries =
      action_entries
      |> Enum.filter(&(&1["status"] == "completed" && !is_nil(&1["compensationActionKey"])))
      |> Enum.reverse()
      |> Enum.map(fn action_entry ->
        %{
          "publicId" => generate_public_id(),
          "stepId" => action_entry["stepId"] <> ":compensation",
          "actionKey" => action_entry["compensationActionKey"],
          "label" => "Compensate: " <> action_entry["label"],
          "status" => "compensated",
          "delaySeconds" => nil,
          "compensationActionKey" => nil,
          "createdAt" => now(),
          "changedBy" => initiated_by
        }
      end)

    action_entries ++ compensation_entries
  end

  defp advance_execution_action(execution, action, plan_length) do
    case action["actionKey"] do
      "delay.wait" ->
        advance_delay_action(execution, action, plan_length)

      _action_key ->
        updated_execution =
          execution
          |> Map.put("currentActionIndex", execution["currentActionIndex"] + 1)
          |> Map.put("waitingUntil", nil)

        action_entry = build_action_entry(action, "completed")
        transition_status = if updated_execution["currentActionIndex"] >= plan_length, do: "completed", else: nil

        finalized_execution =
          if transition_status == "completed" do
            updated_execution
            |> Map.put("status", "completed")
            |> maybe_put_timestamp("completed")
          else
            updated_execution
          end

        {:completed, finalized_execution, action_entry, transition_status}
    end
  end

  defp advance_delay_action(execution, action, plan_length) do
    waiting_until = execution["waitingUntil"]

    cond do
      is_nil(waiting_until) ->
        updated_execution =
          execution
          |> Map.put("waitingUntil", shift_seconds(action["delaySeconds"]))

        {:waiting, updated_execution, build_action_entry(action, "waiting")}

      wait_due?(waiting_until) ->
        updated_execution =
          execution
          |> Map.put("currentActionIndex", execution["currentActionIndex"] + 1)
          |> Map.put("waitingUntil", nil)

        transition_status = if updated_execution["currentActionIndex"] >= plan_length, do: "completed", else: nil

        finalized_execution =
          if transition_status == "completed" do
            updated_execution
            |> Map.put("status", "completed")
            |> maybe_put_timestamp("completed")
          else
            updated_execution
          end

        {:completed, finalized_execution, build_action_entry(action, "completed"), transition_status}

      true ->
        {:error, :waiting}
    end
  end

  defp fetch_execution_plan_postgres(execution) do
    case Postgres.query!(
           """
           SELECT version.snapshot_actions
           FROM workflow_control.workflow_definition_versions AS version
           JOIN workflow_control.workflow_definitions AS definition ON definition.id = version.workflow_definition_id
           JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
           WHERE tenant.slug = $1
             AND definition.key = $2
             AND version.version_number = $3
           LIMIT 1
           """,
           [
             execution["tenantSlug"],
             execution["workflowDefinitionKey"],
             execution["workflowDefinitionVersionNumber"]
           ]
         ).rows do
      [[snapshot_actions]] ->
        {:ok, Enum.map(snapshot_actions || [], &normalize_action_map/1)}

      _ ->
        {:error, :workflow_definition_version_not_found}
    end
  end

  defp normalize_action_map(action) do
    %{
      "stepId" => action["stepId"] || action[:stepId],
      "actionKey" => action["actionKey"] || action[:actionKey],
      "label" => action["label"] || action[:label],
      "delaySeconds" => action["delaySeconds"] || action[:delaySeconds],
      "compensationActionKey" => action["compensationActionKey"] || action[:compensationActionKey]
    }
  end

  defp wait_due?(waiting_until) do
    {:ok, waiting_at, _offset} = DateTime.from_iso8601(waiting_until)
    DateTime.compare(DateTime.utc_now(), waiting_at) != :lt
  end

  defp to_datetime_param(nil), do: nil

  defp to_datetime_param(value) when is_binary(value) do
    {:ok, datetime, _offset} = DateTime.from_iso8601(value)
    datetime
  end

  defp shift_seconds(nil), do: nil

  defp shift_seconds(seconds) do
    DateTime.utc_now()
    |> DateTime.add(seconds, :second)
    |> DateTime.truncate(:second)
    |> DateTime.to_iso8601()
  end

  defp now do
    DateTime.utc_now() |> DateTime.truncate(:second) |> DateTime.to_iso8601()
  end
end
