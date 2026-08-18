import { FormEvent, useId, useState } from 'react';
import type { ApiClient } from '../api/client';
import { rememberUrl, readRecentUrls } from '../hooks/useRecentUrls';
import { useAudit } from '../hooks/useAudit';
import { Alert } from './Alert';
import { EngineBadge } from './EngineBadge';
import { RecommendationList } from './RecommendationList';
import { ScoreCard } from './ScoreCard';
import { SignalList } from './SignalList';

type AuditFormProps = {
  client: ApiClient;
};

export function AuditForm({ client }: AuditFormProps) {
  const inputId = useId();
  const hintId = useId();
  const errorId = useId();
  const listId = useId();
  const [url, setUrl] = useState('https://');
  const [recent, setRecent] = useState<string[]>(() =>
    typeof window === 'undefined' ? [] : readRecentUrls(),
  );
  const [clientError, setClientError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const { state, submit, cancel } = useAudit(client);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmed = url.trim();
    if (!/^https?:\/\//i.test(trimmed) || trimmed.length < 12) {
      setClientError('Enter a full URL including https://');
      return;
    }
    setClientError(null);
    setCopied(false);
    setRecent(rememberUrl(trimmed));
    await submit(trimmed);
  }

  async function copyLink(id: string) {
    const share = `${window.location.origin}/result/?id=${id}`;
    try {
      await navigator.clipboard.writeText(share);
      setCopied(true);
    } catch {
      setCopied(false);
    }
  }

  const fieldError =
    clientError ??
    (state.status === 'error' && state.error.field === 'url' ? state.error.message : null);
  const describedBy = [hintId, fieldError ? errorId : null].filter(Boolean).join(' ') || undefined;

  return (
    <div className="perfcheck-audit">
      <form className="perfcheck-form" onSubmit={onSubmit} noValidate>
        <div>
          <label className="perfcheck-form__label" htmlFor={inputId}>
            Website URL
          </label>
          <input
            id={inputId}
            className="perfcheck-form__input"
            type="url"
            name="url"
            inputMode="url"
            autoComplete="url"
            list={listId}
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            aria-invalid={fieldError ? true : undefined}
            aria-describedby={describedBy}
            aria-busy={state.status === 'submitting' || undefined}
            required
          />
          <datalist id={listId}>
            {recent.map((item) => (
              <option key={item} value={item} />
            ))}
          </datalist>
          <p id={hintId} className="perfcheck-form__hint">
            Example: https://example.com
          </p>
          {fieldError ? (
            <p id={errorId} className="perfcheck-form__error" role="alert">
              {fieldError}
            </p>
          ) : null}
        </div>
        <div className="perfcheck-form__actions">
          <button
            className="perfcheck-button perfcheck-button--primary"
            type="submit"
            disabled={state.status === 'submitting'}
          >
            {state.status === 'submitting' ? 'Running audit…' : 'Run audit'}
          </button>
          {state.status === 'submitting' ? (
            <button className="perfcheck-button perfcheck-button--ghost" type="button" onClick={cancel}>
              Cancel
            </button>
          ) : null}
        </div>
      </form>

      <div role="status" aria-live="polite" className="perfcheck-audit__live">
        {state.status === 'submitting' ? (
          <p>
            Running audit… {state.elapsed}s elapsed. Live engines can take up to a minute.
          </p>
        ) : null}
        {state.status === 'error' && state.error.field !== 'url' ? (
          <Alert title="Audit failed" variant="error">
            {state.error.message}
          </Alert>
        ) : null}
        {state.status === 'success' ? (
          <>
            <p>
              <EngineBadge engine={state.audit.engine} />
            </p>
            <ScoreCard scores={state.audit.scores} url={state.audit.url} />
            <SignalList signals={state.audit.signals} />
            <RecommendationList items={state.audit.recommendations} />
            <p>
              <button
                type="button"
                className="perfcheck-button perfcheck-button--ghost"
                onClick={() => void copyLink(state.audit.id)}
              >
                Copy link
              </button>
              {copied ? <span className="perfcheck-form__hint"> Link copied.</span> : null}
            </p>
          </>
        ) : null}
      </div>
    </div>
  );
}
