import { describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Dashboard } from './Dashboard';
import { ApiError } from '../api/types';
import type { ApiClient } from '../api/client';

function client(partial: Partial<ApiClient> = {}): ApiClient {
  return {
    session: vi.fn().mockResolvedValue({ authRequired: false, authenticated: true }),
    logout: vi.fn().mockResolvedValue(null),
    ...partial,
  } as ApiClient;
}

describe('Dashboard', () => {
  it('shows empty state', async () => {
    render(<Dashboard client={client({ listAudits: vi.fn().mockResolvedValue({ items: [] }) })} />);
    expect(await screen.findByText(/no audits yet/i)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /sign out/i })).not.toBeInTheDocument();
  });

  it('lists audits', async () => {
    render(
      <Dashboard
        client={client({
          listAudits: vi.fn().mockResolvedValue({
            items: [
              {
                id: '1',
                url: 'https://example.com',
                createdAt: '2026-08-18T12:00:00Z',
                engine: 'fetch',
                scores: { overall: 80, performance: 70, seo: 90, accessibility: 80 },
                recommendations: [],
              },
            ],
          }),
        })}
      />,
    );
    await waitFor(() => {
      expect(screen.getByText('https://example.com')).toBeInTheDocument();
    });
  });

  it('asks for a password on 401', async () => {
    render(
      <Dashboard
        client={client({
          listAudits: vi.fn().mockRejectedValue(new ApiError('unauthorized', 'Sign in', 401)),
        })}
      />,
    );
    expect(await screen.findByLabelText(/dashboard password/i)).toBeInTheDocument();
  });

  it('signs out when the dashboard is protected', async () => {
    const user = userEvent.setup();
    const logout = vi.fn().mockResolvedValue(null);
    const listAudits = vi
      .fn()
      .mockResolvedValueOnce({
        items: [
          {
            id: '1',
            url: 'https://example.com',
            createdAt: '2026-08-18T12:00:00Z',
            engine: 'fetch',
            scores: { overall: 80, performance: 70, seo: 90, accessibility: 80 },
            recommendations: [],
          },
        ],
      })
      .mockRejectedValueOnce(new ApiError('unauthorized', 'Sign in', 401));

    render(
      <Dashboard
        client={client({
          listAudits,
          logout,
          session: vi.fn().mockResolvedValue({ authRequired: true, authenticated: true }),
        })}
      />,
    );

    expect(await screen.findByText('https://example.com')).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: /sign out/i }));
    expect(logout).toHaveBeenCalled();
    expect(await screen.findByLabelText(/dashboard password/i)).toBeInTheDocument();
  });
});
