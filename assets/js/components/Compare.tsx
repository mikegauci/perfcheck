import { FormEvent, useState } from 'react';
import type { ApiClient } from '../api/client';
import type { Audit } from '../api/types';
import { Alert } from './Alert';
import { EngineBadge } from './EngineBadge';
import { ScoreGauge } from './ScoreGauge';

type CompareProps = {
  client: ApiClient;
};

type Side = {
  url: string;
  status: 'idle' | 'loading' | 'error' | 'ready';
  audit?: Audit;
  message?: string;
};

export function Compare({ client }: CompareProps) {
  const [leftUrl, setLeftUrl] = useState('https://');
  const [rightUrl, setRightUrl] = useState('https://');
  const [left, setLeft] = useState<Side>({ url: '', status: 'idle' });
  const [right, setRight] = useState<Side>({ url: '', status: 'idle' });

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setLeft({ url: leftUrl, status: 'loading' });
    setRight({ url: rightUrl, status: 'loading' });
    const [a, b] = await Promise.allSettled([
      client.createAudit(leftUrl.trim()),
      client.createAudit(rightUrl.trim()),
    ]);
    setLeft(fromSettled(leftUrl, a));
    setRight(fromSettled(rightUrl, b));
  }

  return (
    <div>
      <form className="perfcheck-form" onSubmit={onSubmit}>
        <div>
          <label className="perfcheck-form__label" htmlFor="compare-a">
            First URL
          </label>
          <input
            id="compare-a"
            className="perfcheck-form__input"
            type="url"
            value={leftUrl}
            onChange={(e) => setLeftUrl(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="perfcheck-form__label" htmlFor="compare-b">
            Second URL
          </label>
          <input
            id="compare-b"
            className="perfcheck-form__input"
            type="url"
            value={rightUrl}
            onChange={(e) => setRightUrl(e.target.value)}
            required
          />
        </div>
        <div className="perfcheck-form__actions">
          <button className="perfcheck-button perfcheck-button--primary" type="submit">
            Compare
          </button>
        </div>
      </form>
      <div className="perfcheck-compare">
        <CompareColumn side={left} other={right} />
        <CompareColumn side={right} other={left} />
      </div>
    </div>
  );
}

function fromSettled(url: string, result: PromiseSettledResult<Audit>): Side {
  if (result.status === 'fulfilled') {
    return { url, status: 'ready', audit: result.value };
  }
  const message = result.reason instanceof Error ? result.reason.message : 'Failed';
  return { url, status: 'error', message };
}

function CompareColumn({ side, other }: { side: Side; other: Side }) {
  if (side.status === 'idle') {
    return <div className="perfcheck-card"><p className="perfcheck-card__body">Waiting…</p></div>;
  }
  if (side.status === 'loading') {
    return <div className="perfcheck-card"><p className="perfcheck-card__body">Running audit…</p></div>;
  }
  if (side.status === 'error' || !side.audit) {
    return <Alert title="Could not audit this URL" variant="error">{side.message ?? 'Unknown error'}</Alert>;
  }
  const delta = other.status === 'ready' && other.audit
    ? side.audit.scores.overall - other.audit.scores.overall
    : null;
  return (
    <article className="perfcheck-card">
      <h2 className="perfcheck-card__title">{side.audit.url}</h2>
      <EngineBadge engine={side.audit.engine} />
      <ScoreGauge label="Overall" value={side.audit.scores.overall} />
      <p className="perfcheck-card__meta">
        Perf {side.audit.scores.performance} · SEO {side.audit.scores.seo} · A11y {side.audit.scores.accessibility}
        {delta != null ? ` · ${delta > 0 ? '+' : ''}${delta} vs other` : ''}
      </p>
    </article>
  );
}
