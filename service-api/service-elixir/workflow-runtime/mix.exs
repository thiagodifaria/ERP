defmodule WorkflowRuntime.MixProject do
  use Mix.Project

  def project do
    [
      app: :workflow_runtime,
      version: "0.1.0",
      elixir: "~> 1.17",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  def application do
    [
      extra_applications: [:logger],
      mod: {WorkflowRuntime.Application, []}
    ]
  end

  defp deps do
    [
      {:jason, "~> 1.4"},
      {:plug_cowboy, "~> 2.7"},
      {:postgrex, "~> 0.20"}
    ]
  end
end
