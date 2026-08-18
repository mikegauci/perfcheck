import { FormEvent, useState } from 'react';
import type { ApiClient } from '../api/client';
import { ApiError } from '../api/types';

type LoginFormProps = {
  client: ApiClient;
  onSuccess: () => void;
};

export function LoginForm({ client, onSuccess }: LoginFormProps) {
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await client.login(password);
      onSuccess();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not sign in.');
    } finally {
      setBusy(false);
    }
  }

  return (
    <form className="perfcheck-form" onSubmit={onSubmit}>
      <div>
        <label className="perfcheck-form__label" htmlFor="dashboard-password">
          Dashboard password
        </label>
        <input
          id="dashboard-password"
          className="perfcheck-form__input"
          type="password"
          name="password"
          autoComplete="current-password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        {error ? (
          <p className="perfcheck-form__error" role="alert">
            {error}
          </p>
        ) : null}
      </div>
      <div className="perfcheck-form__actions">
        <button className="perfcheck-button perfcheck-button--primary" type="submit" disabled={busy}>
          {busy ? 'Signing in…' : 'Sign in'}
        </button>
      </div>
    </form>
  );
}
