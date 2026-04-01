defmodule WorkflowRuntime.Application.ExecutionStore do
  @moduledoc false

  use Agent

  def start_link(_opts) do
    Agent.start_link(fn -> [] end, name: __MODULE__)
  end

  def list do
    Agent.get(__MODULE__, &Enum.reverse/1)
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
end
