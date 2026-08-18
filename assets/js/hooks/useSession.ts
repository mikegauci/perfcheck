import { useCallback, useState } from 'react';
import type { ApiClient } from '../api/client';

export function useSession(client: ApiClient) {
  const [authed, setAuthed] = useState(false);

  const markAuthed = useCallback(() => setAuthed(true), []);
  const signOut = useCallback(async () => {
    await client.logout();
    setAuthed(false);
  }, [client]);

  return { authed, markAuthed, signOut };
}
