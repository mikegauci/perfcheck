import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Compare } from './Compare';
import type { Audit } from '../api/types';
import type { ApiClient } from '../api/client';

const audit = (url: string, overall: number): Audit => ({
  id: overall.toString(),
  url,
  createdAt: '2026-08-18T12:00:00Z',
  engine: 'mock',
  scores: { overall, performance: overall, seo: overall, accessibility: overall },
  recommendations: [],
});

describe('Compare', () => {
  it('shows both results', async () => {
    const user = userEvent.setup();
    const client = {
      createAudit: vi
        .fn()
        .mockResolvedValueOnce(audit('https://a.example', 80))
        .mockResolvedValueOnce(audit('https://b.example', 60)),
    } as unknown as ApiClient;
    render(<Compare client={client} />);
    await user.clear(screen.getByLabelText(/first url/i));
    await user.type(screen.getByLabelText(/first url/i), 'https://a.example');
    await user.clear(screen.getByLabelText(/second url/i));
    await user.type(screen.getByLabelText(/second url/i), 'https://b.example');
    await user.click(screen.getByRole('button', { name: /compare/i }));
    expect(await screen.findByText('https://a.example')).toBeInTheDocument();
    expect(screen.getByText('https://b.example')).toBeInTheDocument();
  });
});
