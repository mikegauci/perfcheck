import { useEffect, useState } from 'react';
import type { ApiClient } from '../api/client';
import { ApiError } from '../api/types';
import { Alert } from './Alert';
import { AuditReport } from './AuditReport';

type ResultViewProps = {
  client: ApiClient;
};

export function ResultView({ client }: ResultViewProps) {
  const [id] = useState(() => new URLSearchParams(window.location.search).get('id') ?? '');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [audit, setAudit] = useState<Awaited<ReturnType<ApiClient['getAudit']>> | null>(null);

  useEffect(() => {
    if (!id) {
      setLoading(false);
      setError('Missing audit id.');
      return;
    }
    let cancelled = false;
    (async () => {
      try {
        const result = await client.getAudit(id);
        if (!cancelled) {
          setAudit(result);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof ApiError ? err.message : 'Could not load this audit.');
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [client, id]);

  if (loading) {
    return <p>Loading result…</p>;
  }
  if (error || !audit) {
    return (
      <Alert title="Result unavailable" variant="error">
        {error ?? 'Unknown error'}
      </Alert>
    );
  }

  return <AuditReport audit={audit} />;
}
