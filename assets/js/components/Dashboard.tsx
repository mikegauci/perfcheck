import { useCallback, useEffect, useMemo, useState } from 'react';
import type { ApiClient } from '../api/client';
import { ApiError, type Audit } from '../api/types';
import { useSession } from '../hooks/useSession';
import { Alert } from './Alert';
import { EngineBadge } from './EngineBadge';
import { LoginForm } from './LoginForm';

type DashboardProps = {
  client: ApiClient;
};

type DashState =
  | { status: 'loading' }
  | { status: 'auth' }
  | { status: 'empty' }
  | { status: 'ready'; items: Audit[] }
  | { status: 'error'; message: string };

type SortKey = 'newest' | 'oldest' | 'score-desc' | 'score-asc';

export function Dashboard({ client }: DashboardProps) {
  const [state, setState] = useState<DashState>({ status: 'loading' });
  const [sort, setSort] = useState<SortKey>('newest');
  const [query, setQuery] = useState('');
  const [visible, setVisible] = useState(10);
  const [reload, setReload] = useState(0);
  const { authRequired, authenticated, refresh, signOut } = useSession(client);
  const canSignOut = authRequired && authenticated;

  const load = useCallback(async () => {
    setState({ status: 'loading' });
    try {
      const res = await client.listAudits(50);
      await refresh().catch(() => undefined);
      if (res.items.length === 0) {
        setState({ status: 'empty' });
      } else {
        setState({ status: 'ready', items: res.items });
      }
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        setState({ status: 'auth' });
        return;
      }
      const message = err instanceof ApiError ? err.message : 'Could not load audits.';
      setState({ status: 'error', message });
    }
  }, [client, refresh]);

  useEffect(() => {
    void load();
  }, [load, reload]);

  async function onSignOut() {
    await signOut();
    setReload((n) => n + 1);
  }

  const filtered = useMemo(() => {
    if (state.status !== 'ready') return [];
    const q = query.trim().toLowerCase();
    let items = state.items;
    if (q) {
      items = items.filter((a) => a.url.toLowerCase().includes(q));
    }
    const sorted = [...items];
    sorted.sort((a, b) => {
      switch (sort) {
        case 'oldest':
          return a.createdAt.localeCompare(b.createdAt);
        case 'score-desc':
          return b.scores.overall - a.scores.overall;
        case 'score-asc':
          return a.scores.overall - b.scores.overall;
        case 'newest':
        default:
          return b.createdAt.localeCompare(a.createdAt);
      }
    });
    return sorted;
  }, [state, sort, query]);

  if (state.status === 'loading') {
    return (
      <div className="perfcheck-dashboard" aria-busy="true">
        <div className="perfcheck-dashboard__skeleton" />
        <div className="perfcheck-dashboard__skeleton" />
        <p className="perfcheck-sr-only">Loading audits…</p>
      </div>
    );
  }

  if (state.status === 'auth') {
    return (
      <div>
        <Alert title="Sign in required" variant="info">
          This dashboard is protected. Enter the password configured on the API.
        </Alert>
        <LoginForm client={client} onSuccess={() => setReload((n) => n + 1)} />
      </div>
    );
  }

  if (state.status === 'error') {
    return (
      <div>
        <Alert title="Could not load dashboard" variant="error">{state.message}</Alert>
        <button type="button" className="perfcheck-button perfcheck-button--primary" onClick={() => void load()}>
          Retry
        </button>
      </div>
    );
  }

  if (state.status === 'empty') {
    return (
      <div className="perfcheck-dashboard">
        {canSignOut ? (
          <div className="perfcheck-dashboard__toolbar">
            <button
              type="button"
              className="perfcheck-button perfcheck-button--ghost"
              onClick={() => void onSignOut()}
            >
              Sign out
            </button>
          </div>
        ) : null}
        <Alert title="No audits yet" variant="info">
          Run an audit on the <a href="/audit/">audit page</a> to populate this list.
        </Alert>
      </div>
    );
  }

  return (
    <div className="perfcheck-dashboard">
      <div className="perfcheck-dashboard__toolbar">
        <div>
          <label className="perfcheck-form__label" htmlFor="dash-filter">
            Filter by URL
          </label>
          <input
            id="dash-filter"
            className="perfcheck-form__input"
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <div>
          <label className="perfcheck-form__label" htmlFor="dash-sort">
            Sort
          </label>
          <select
            id="dash-sort"
            className="perfcheck-form__input"
            value={sort}
            onChange={(e) => setSort(e.target.value as SortKey)}
          >
            <option value="newest">Newest first</option>
            <option value="oldest">Oldest first</option>
            <option value="score-desc">Highest score</option>
            <option value="score-asc">Lowest score</option>
          </select>
        </div>
        {canSignOut ? (
          <button
            type="button"
            className="perfcheck-button perfcheck-button--ghost"
            onClick={() => void onSignOut()}
          >
            Sign out
          </button>
        ) : null}
      </div>

      {filtered.length === 0 ? (
        <Alert title="No matches" variant="info">
          No audits match that filter.
        </Alert>
      ) : (
        <>
        <ul className="perfcheck-dashboard__list">
          {filtered.slice(0, visible).map((item) => (
            <li key={item.id} className="perfcheck-dashboard__item">
              <div className="perfcheck-dashboard__url">
                <a href={`/result/?id=${item.id}`}>{item.url}</a>
              </div>
              <EngineBadge engine={item.engine} />
              <div className="perfcheck-dashboard__scores">
                <span>Overall {item.scores.overall}</span>
                <span>Perf {item.scores.performance}</span>
                <span>SEO {item.scores.seo}</span>
                <span>A11y {item.scores.accessibility}</span>
                <span>{new Date(item.createdAt).toLocaleString()}</span>
              </div>
            </li>
          ))}
        </ul>
        {visible < filtered.length ? (
          <button
            type="button"
            className="perfcheck-button perfcheck-button--ghost"
            onClick={() => setVisible((n) => n + 10)}
          >
            Load more
          </button>
        ) : null}
        </>
      )}
    </div>
  );
}
