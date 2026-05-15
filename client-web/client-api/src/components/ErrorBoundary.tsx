import { Component, type ErrorInfo, type ReactNode } from "react";

type Props = {
  children: ReactNode;
};

type State = {
  error: Error | null;
};

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    console.error("ERP Control Console render failure", error, errorInfo);
  }

  render(): ReactNode {
    if (!this.state.error) {
      return this.props.children;
    }

    return (
      <main className="min-h-screen bg-white px-6 py-12 text-gray-900">
        <div className="mx-auto max-w-2xl rounded-xl border border-red-200 bg-red-50 p-6">
          <p className="text-sm font-semibold uppercase tracking-wide text-red-700">Falha na interface</p>
          <h1 className="mt-2 text-2xl font-semibold">O console encontrou um erro inesperado.</h1>
          <p className="mt-3 text-sm text-red-900">{this.state.error.message}</p>
          <button
            className="mt-6 rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
            onClick={() => this.setState({ error: null })}
            type="button"
          >
            Tentar novamente
          </button>
        </div>
      </main>
    );
  }
}
