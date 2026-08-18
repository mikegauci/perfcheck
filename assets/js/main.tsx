import { createRoot } from 'react-dom/client';
import * as params from '@params';
import { createClient } from './api/client';
import { AuditForm } from './components/AuditForm';
import { Compare } from './components/Compare';
import { Dashboard } from './components/Dashboard';
import { ErrorBoundary } from './components/ErrorBoundary';
import { ResultView } from './components/ResultView';

const client = createClient({ baseUrl: params.apiBaseUrl ?? '' });

function mountIslands() {
  const nodes = document.querySelectorAll<HTMLElement>('[data-react-island]');
  nodes.forEach((node) => {
    const component = node.dataset.component;
    const rootEl = node.querySelector<HTMLElement>('[data-island-root]');
    if (!rootEl) {
      return;
    }

    const root = createRoot(rootEl);
    if (component === 'audit-form') {
      root.render(
        <ErrorBoundary>
          <AuditForm client={client} />
        </ErrorBoundary>,
      );
    } else if (component === 'dashboard') {
      root.render(
        <ErrorBoundary>
          <Dashboard client={client} />
        </ErrorBoundary>,
      );
    } else if (component === 'result-view') {
      root.render(
        <ErrorBoundary>
          <ResultView client={client} />
        </ErrorBoundary>,
      );
    } else if (component === 'compare') {
      root.render(
        <ErrorBoundary>
          <Compare client={client} />
        </ErrorBoundary>,
      );
    } else {
      return;
    }

    node.removeAttribute('aria-busy');
  });
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', mountIslands);
} else {
  mountIslands();
}
