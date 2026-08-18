import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AuditForm } from './AuditForm';
import { ApiError, type Audit } from '../api/types';
import type { ApiClient } from '../api/client';

const sample: Audit = {
  id: '1',
  url: 'https://example.com',
  createdAt: '2026-08-18T12:00:00Z',
  scores: { overall: 80, performance: 72, seo: 88, accessibility: 74 },
  recommendations: [
    {
      id: 'img-format',
      category: 'performance',
      severity: 'high',
      title: 'Serve images in modern formats',
      detail: 'Use AVIF or WebP.',
    },
  ],
};

describe('AuditForm', () => {
  it('validates empty-ish URLs client-side', async () => {
    const user = userEvent.setup();
    const client = { createAudit: vi.fn() } as unknown as ApiClient;
    render(<AuditForm client={client} />);

    await user.clear(screen.getByLabelText(/website url/i));
    await user.type(screen.getByLabelText(/website url/i), 'ftp://x');
    await user.click(screen.getByRole('button', { name: /run audit/i }));

    expect(await screen.findByRole('alert')).toHaveTextContent(/full URL/i);
    expect(client.createAudit).not.toHaveBeenCalled();
  });

  it('shows scores on success', async () => {
    const user = userEvent.setup();
    const client = {
      createAudit: vi.fn().mockResolvedValue(sample),
    } as unknown as ApiClient;
    render(<AuditForm client={client} />);

    await user.clear(screen.getByLabelText(/website url/i));
    await user.type(screen.getByLabelText(/website url/i), 'https://example.com');
    await user.click(screen.getByRole('button', { name: /run audit/i }));

    expect(await screen.findByRole('heading', { name: /results for/i })).toBeInTheDocument();
    expect(screen.getByText(/recommendations/i)).toBeInTheDocument();
  });

  it('shows server errors', async () => {
    const user = userEvent.setup();
    const client = {
      createAudit: vi.fn().mockRejectedValue(new ApiError('internal', 'boom', 500)),
    } as unknown as ApiClient;
    render(<AuditForm client={client} />);

    await user.clear(screen.getByLabelText(/website url/i));
    await user.type(screen.getByLabelText(/website url/i), 'https://example.com');
    await user.click(screen.getByRole('button', { name: /run audit/i }));

    expect(await screen.findByRole('alert')).toHaveTextContent(/boom/i);
  });
});
