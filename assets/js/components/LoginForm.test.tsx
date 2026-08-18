import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LoginForm } from './LoginForm';
import { ApiError } from '../api/types';
import type { ApiClient } from '../api/client';

describe('LoginForm', () => {
  it('signs in successfully', async () => {
    const user = userEvent.setup();
    const onSuccess = vi.fn();
    const client = {
      login: vi.fn().mockResolvedValue({ ok: true }),
    } as unknown as ApiClient;
    render(<LoginForm client={client} onSuccess={onSuccess} />);
    await user.type(screen.getByLabelText(/dashboard password/i), 'secret');
    await user.click(screen.getByRole('button', { name: /sign in/i }));
    expect(client.login).toHaveBeenCalledWith('secret');
    expect(onSuccess).toHaveBeenCalled();
  });

  it('shows an error on failure', async () => {
    const user = userEvent.setup();
    const client = {
      login: vi.fn().mockRejectedValue(new ApiError('unauthorized', 'Incorrect password.', 401, 'password')),
    } as unknown as ApiClient;
    render(<LoginForm client={client} onSuccess={vi.fn()} />);
    await user.type(screen.getByLabelText(/dashboard password/i), 'nope');
    await user.click(screen.getByRole('button', { name: /sign in/i }));
    expect(await screen.findByRole('alert')).toHaveTextContent(/incorrect password/i);
  });
});
