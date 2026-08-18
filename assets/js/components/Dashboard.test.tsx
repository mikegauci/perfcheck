import { describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { Dashboard } from './Dashboard';
import type { ApiClient } from '../api/client';

describe('Dashboard', () => {
  it('shows empty state', async () => {
    const client = {
      listAudits: vi.fn().mockResolvedValue({ items: [] }),
    } as unknown as ApiClient;
    render(<Dashboard client={client} />);
    expect(await screen.findByText(/no audits yet/i)).toBeInTheDocument();
  });

  it('lists audits', async () => {
    const client = {
      listAudits: vi.fn().mockResolvedValue({
        items: [
          {
            id: '1',
            url: 'https://example.com',
            createdAt: '2026-08-18T12:00:00Z',
            engine: 'mock',
            scores: { overall: 80, performance: 70, seo: 90, accessibility: 80 },
            recommendations: [],
          },
        ],
      }),
    } as unknown as ApiClient;
    render(<Dashboard client={client} />);
    await waitFor(() => {
      expect(screen.getByText('https://example.com')).toBeInTheDocument();
    });
  });
});
