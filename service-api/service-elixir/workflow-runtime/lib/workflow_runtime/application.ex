defmodule WorkflowRuntime.Application do
  @moduledoc false

  use Application

  def start(_type, _args) do
    validate_runtime_posture!()

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

  defp validate_runtime_posture! do
    environment =
      System.get_env("ERP_ENV", "local")
      |> String.trim()
      |> String.downcase()

    if environment not in ["", "local", "dev", "development", "test", "testing"] do
      repository_driver = Application.fetch_env!(:workflow_runtime, :repository_driver)
      postgres_password = Application.fetch_env!(:workflow_runtime, :postgres_password)

      if repository_driver != "postgres" do
        raise "workflow_runtime_requires_postgres_outside_local"
      end

      if postgres_password in ["", "erp", "admin"] or String.starts_with?(postgres_password, "change-me-unsafe-local-only") do
        raise "workflow_runtime_requires_non_local_postgres_password"
      end
    end
  end
end
