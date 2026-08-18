import { createRoot } from 'react-dom/client';
import * as params from '@params';
import { createClient } from './api/client';
import { AuditForm } from './components/AuditForm';
import { Compare } from './components/Compare';
import { Dashboard } from './components/Dashboard';
import { ResultView } from './components/ResultView';

const client = createClient({ baseUrl: params.apiBaseUrl ?? '' });

function mountIslands() {
  const nodes = document.querySelectorAll<HTMLElement>('[data-react-island]');
  nodes.forEach((node) => {
    const component = node.dataset.component;
    const fallback = node.querySelector('[data-island-fallback]');
    if (fallback) {
      fallback.remove();
    }

    const root = createRoot(node);
    if (component === 'audit-form') {
      root.render(<AuditForm client={client} />);
      return;
    }
    if (component === 'dashboard') {
      root.render(<Dashboard client={client} />);
      return;
    }
    if (component === 'result-view') {
      root.render(<ResultView client={client} />);
      return;
    }
    if (component === 'compare') {
      root.render(<Compare client={client} />);
      return;
    }
  });
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', mountIslands);
} else {
  mountIslands();
}
