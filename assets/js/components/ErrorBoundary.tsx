import { Component, type ErrorInfo, type ReactNode } from 'react';
import { Alert } from './Alert';

type ErrorBoundaryProps = {
  children: ReactNode;
};

type ErrorBoundaryState = {
  error: Error | null;
};

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error(error, info.componentStack);
  }

  render() {
    if (this.state.error) {
      return (
        <Alert title="This tool failed to load" variant="error">
          {this.state.error.message}
        </Alert>
      );
    }
    return this.props.children;
  }
}
