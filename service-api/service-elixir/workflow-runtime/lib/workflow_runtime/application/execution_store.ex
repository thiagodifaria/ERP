defmodule WorkflowRuntime.Application.ExecutionStore do
  @moduledoc false

  use Agent

  def start_link(_opts) do
    Agent.start_link(fn -> [] end, name: __MODULE__)
  end

  def list do
    Agent.get(__MODULE__, &Enum.reverse/1)
  end

  def find(public_id) do
    Agent.get(__MODULE__, fn executions ->
      Enum.find(executions, &(&1["publicId"] == public_id))
    end)
  end

  def create(attributes) do
    execution = %{
      "publicId" => generate_public_id(),
      "workflowDefinitionKey" => String.trim(attributes["workflowDefinitionKey"] || ""),
      "subjectType" => String.trim(attributes["subjectType"] || ""),
      "subjectPublicId" => String.trim(attributes["subjectPublicId"] || ""),
      "initiatedBy" => String.trim(attributes["initiatedBy"] || ""),
      "status" => "pending",
      "createdAt" => now(),
      "startedAt" => nil,
      "completedAt" => nil,
      "failedAt" => nil
    }

    Agent.update(__MODULE__, fn executions -> [execution | executions] end)
    execution
  end

  def transition(public_id, next_status, allowed_statuses) do
    Agent.get_and_update(__MODULE__, fn executions ->
      case Enum.split_with(executions, &(&1["publicId"] != public_id)) do
        {_, []} ->
          {{:error, :not_found}, executions}

        {left, [execution | right]} ->
          if execution["status"] in allowed_statuses do
            updated_execution =
              execution
              |> Map.put("status", next_status)
              |> maybe_put_timestamp(next_status)

            {{:ok, updated_execution}, left ++ [updated_execution | right]}
          else
            {{:error, :invalid_transition}, executions}
          end
      end
    end)
  end

  def reset do
    Agent.update(__MODULE__, fn _executions -> [] end)
  end

  def summary do
    executions = list()

    %{
      total: length(executions),
      pending: Enum.count(executions, &(&1["status"] == "pending")),
      running: Enum.count(executions, &(&1["status"] == "running")),
      completed: Enum.count(executions, &(&1["status"] == "completed")),
      failed: Enum.count(executions, &(&1["status"] == "failed"))
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
  defp maybe_put_timestamp(execution, _next_status), do: execution

  defp now do
    DateTime.utc_now() |> DateTime.truncate(:second) |> DateTime.to_iso8601()
  end
end
