import { FormEvent, useId, useState } from 'react';
import type { ApiClient } from '../api/client';
import { useAudit } from '../hooks/useAudit';
import { Alert } from './Alert';
import { RecommendationList } from './RecommendationList';
import { ScoreCard } from './ScoreCard';

type AuditFormProps = {
  client: ApiClient;
};

export function AuditForm({ client }: AuditFormProps) {
  const inputId = useId();
  const hintId = useId();
  const errorId = useId();
  const [url, setUrl] = useState('https://');
  const [clientError, setClientError] = useState<string | null>(null);
  const { state, submit } = useAudit(client);

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmed = url.trim();
    if (!/^https?:\/\//i.test(trimmed) || trimmed.length < 12) {
      setClientError('Enter a full URL including https://');
      return;
    }
    setClientError(null);
    await submit(trimmed);
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
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            aria-invalid={fieldError ? true : undefined}
            aria-describedby={describedBy}
            aria-busy={state.status === 'submitting' || undefined}
            required
          />
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
        </div>
      </form>

      <div role="status" aria-live="polite" className="perfcheck-audit__live">
        {state.status === 'submitting' ? <p>Running mock audit…</p> : null}
        {state.status === 'error' && state.error.field !== 'url' ? (
          <Alert title="Audit failed" variant="error">
            {state.error.message}
          </Alert>
        ) : null}
        {state.status === 'success' ? (
          <>
            <ScoreCard scores={state.audit.scores} url={state.audit.url} />
            <RecommendationList items={state.audit.recommendations} />
          </>
        ) : null}
      </div>
    </div>
  );
}
