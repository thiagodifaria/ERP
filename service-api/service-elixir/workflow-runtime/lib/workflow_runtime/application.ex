defmodule WorkflowRuntime.Application do
  @moduledoc false

  use Application

  def start(_type, _args) do
    children =
      if WorkflowRuntime.Infrastructure.Postgres.enabled?() do
        [
          WorkflowRuntime.Infrastructure.Postgres,
          WorkflowRuntime.Server
        ]
      else
        [
          WorkflowRuntime.Application.ExecutionStore,
          WorkflowRuntime.Server
        ]
      end

    opts = [strategy: :one_for_one, name: WorkflowRuntime.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
