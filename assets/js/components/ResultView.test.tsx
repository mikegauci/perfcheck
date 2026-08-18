import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ResultView } from './ResultView';
import { ApiError, type Audit } from '../api/types';
import type { ApiClient } from '../api/client';

const sample: Audit = {
  id: 'audit-1',
  url: 'https://example.com',
  createdAt: '2026-08-18T12:00:00Z',
  engine: 'fetch',
  scores: { overall: 80, performance: 72, seo: 88, accessibility: 74 },
  recommendations: [],
};

function setSearch(search: string) {
  window.history.pushState({}, '', search ? `/result/${search}` : '/result/');
}

describe('ResultView', () => {
  beforeEach(() => {
    setSearch('');
  });

  it('shows an error when the id is missing', async () => {
    const client = { getAudit: vi.fn() } as unknown as ApiClient;
    render(<ResultView client={client} />);
    expect(await screen.findByText(/missing audit id/i)).toBeInTheDocument();
    expect(client.getAudit).not.toHaveBeenCalled();
  });

  it('renders a saved audit', async () => {
    setSearch('?id=audit-1');
    const client = {
      getAudit: vi.fn().mockResolvedValue(sample),
    } as unknown as ApiClient;
    render(<ResultView client={client} />);
    expect(await screen.findByRole('heading', { name: /results for/i })).toBeInTheDocument();
    expect(client.getAudit).toHaveBeenCalledWith('audit-1');
  });

  it('shows an API error', async () => {
    setSearch('?id=missing');
    const client = {
      getAudit: vi.fn().mockRejectedValue(new ApiError('not_found', 'Not found.', 404)),
    } as unknown as ApiClient;
    render(<ResultView client={client} />);
    expect(await screen.findByRole('alert')).toHaveTextContent(/not found/i);
  });
});
