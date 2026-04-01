defmodule WorkflowRuntime.Application.ExecutionStore do
  @moduledoc false

  use Agent

  def start_link(_opts) do
    Agent.start_link(fn -> %{executions: [], transitions: %{}} end, name: __MODULE__)
  end

  def list(filters \\ %{}) do
    filters = normalize_filters(filters)

    Agent.get(__MODULE__, fn state ->
      state.executions
      |> Enum.reverse()
      |> Enum.filter(&matches_filters?(&1, filters))
    end)
  end

  def find(public_id) do
    Agent.get(__MODULE__, fn state ->
      Enum.find(state.executions, &(&1["publicId"] == public_id))
    end)
  end

  def transitions(public_id) do
    Agent.get(__MODULE__, fn state ->
      Map.get(state.transitions, public_id, [])
    end)
  end

  def create(attributes) do
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

  def transition(public_id, next_status, allowed_statuses) do
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

  def reset do
    Agent.update(__MODULE__, fn _state -> %{executions: [], transitions: %{}} end)
  end

  def summary do
    executions = list()

    %{
      total: length(executions),
      pending: Enum.count(executions, &(&1["status"] == "pending")),
      running: Enum.count(executions, &(&1["status"] == "running")),
      completed: Enum.count(executions, &(&1["status"] == "completed")),
      failed: Enum.count(executions, &(&1["status"] == "failed")),
      cancelled: Enum.count(executions, &(&1["status"] == "cancelled"))
    }
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
