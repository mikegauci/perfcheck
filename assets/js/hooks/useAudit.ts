import { useCallback, useEffect, useRef, useState } from 'react';
import { ApiError, type Audit } from '../api/types';
import type { ApiClient } from '../api/client';

export type AuditState =
  | { status: 'idle' }
  | { status: 'submitting'; elapsed: number }
  | { status: 'success'; audit: Audit }
  | { status: 'error'; error: ApiError };

export function useAudit(client: ApiClient) {
  const [state, setState] = useState<AuditState>({ status: 'idle' });
  const abortRef = useRef<AbortController | null>(null);

  useEffect(() => {
    if (state.status !== 'submitting') {
      return;
    }
    const timer = window.setInterval(() => {
      setState((prev) =>
        prev.status === 'submitting' ? { status: 'submitting', elapsed: prev.elapsed + 1 } : prev,
      );
    }, 1000);
    return () => window.clearInterval(timer);
  }, [state.status]);

  const cancel = useCallback(() => {
    abortRef.current?.abort();
  }, []);

  const submit = useCallback(
    async (url: string) => {
      abortRef.current?.abort();
      const controller = new AbortController();
      abortRef.current = controller;
      setState({ status: 'submitting', elapsed: 0 });
      try {
        const audit = await client.createAudit(url, { signal: controller.signal });
        setState({ status: 'success', audit });
      } catch (err) {
        const error =
          err instanceof ApiError
            ? err
            : new ApiError('unknown', 'Something went wrong.', 0);
        setState({ status: 'error', error });
      }
    },
    [client],
  );

  const reset = useCallback(() => setState({ status: 'idle' }), []);

  return { state, submit, reset, cancel };
}
