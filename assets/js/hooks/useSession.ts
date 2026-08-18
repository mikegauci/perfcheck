import { useCallback, useState } from 'react';
import type { ApiClient } from '../api/client';

export function useSession(client: ApiClient) {
  const [authRequired, setAuthRequired] = useState(false);
  const [authenticated, setAuthenticated] = useState(false);

  const refresh = useCallback(async () => {
    const status = await client.session();
    setAuthRequired(status.authRequired);
    setAuthenticated(status.authenticated);
    return status;
  }, [client]);

  const signOut = useCallback(async () => {
    await client.logout();
    await refresh();
  }, [client, refresh]);

  return { authRequired, authenticated, refresh, signOut };
}
