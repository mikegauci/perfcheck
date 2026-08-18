import { describe, expect, it, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useAudit } from './useAudit';
import { ApiError, type Audit } from '../api/types';
import type { ApiClient } from '../api/client';

const sample: Audit = {
  id: '1',
  url: 'https://example.com',
  createdAt: '2026-08-18T12:00:00Z',
  scores: { overall: 80, performance: 70, seo: 90, accessibility: 80 },
  recommendations: [],
};

describe('useAudit', () => {
  it('transitions idle -> submitting -> success', async () => {
    const client = {
      createAudit: vi.fn().mockResolvedValue(sample),
    } as unknown as ApiClient;

    const { result } = renderHook(() => useAudit(client));
    expect(result.current.state.status).toBe('idle');

    await act(async () => {
      await result.current.submit('https://example.com');
    });

    expect(result.current.state.status).toBe('success');
    if (result.current.state.status === 'success') {
      expect(result.current.state.audit.id).toBe('1');
    }
  });

  it('transitions to error on ApiError', async () => {
    const client = {
      createAudit: vi.fn().mockRejectedValue(new ApiError('invalid_url', 'bad', 400, 'url')),
    } as unknown as ApiClient;

    const { result } = renderHook(() => useAudit(client));
    await act(async () => {
      await result.current.submit('bad');
    });
    expect(result.current.state.status).toBe('error');
  });
});
