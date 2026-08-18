import { useCallback, useState } from 'react';
import { ApiError, type Audit } from '../api/types';
import type { ApiClient } from '../api/client';

export type AuditState =
  | { status: 'idle' }
  | { status: 'submitting' }
  | { status: 'success'; audit: Audit }
  | { status: 'error'; error: ApiError };

export function useAudit(client: ApiClient) {
  const [state, setState] = useState<AuditState>({ status: 'idle' });

  const submit = useCallback(
    async (url: string) => {
      setState({ status: 'submitting' });
      try {
        const audit = await client.createAudit(url);
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

  return { state, submit, reset };
}
