import { describe, expect, it, vi } from 'vitest';
import { ApiError } from './types';
import { createClient } from './client';

function jsonResponse(status: number, body: unknown) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

describe('createClient', () => {
  it('creates an audit on success', async () => {
    const fetchImpl = vi.fn().mockResolvedValue(
      jsonResponse(201, {
        id: 'abc',
        url: 'https://example.com',
        createdAt: '2026-08-18T12:00:00Z',
        engine: 'mock',
        scores: { overall: 80, performance: 70, seo: 90, accessibility: 80 },
        recommendations: [],
      }),
    );
    const client = createClient({ baseUrl: 'http://api', fetchImpl });
    const audit = await client.createAudit('https://example.com');
    expect(audit.id).toBe('abc');
    expect(fetchImpl).toHaveBeenCalledWith(
      'http://api/api/v1/audits',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('maps error envelopes to ApiError', async () => {
    const fetchImpl = vi.fn().mockResolvedValue(
      jsonResponse(400, {
        error: { code: 'invalid_url', message: 'bad url', field: 'url' },
      }),
    );
    const client = createClient({ fetchImpl });
    await expect(client.createAudit('nope')).rejects.toMatchObject({
      code: 'invalid_url',
      field: 'url',
      status: 400,
    } satisfies Partial<ApiError>);
  });

  it('maps network failures', async () => {
    const fetchImpl = vi.fn().mockRejectedValue(new TypeError('Failed to fetch'));
    const client = createClient({ fetchImpl });
    await expect(client.listAudits()).rejects.toMatchObject({ code: 'network' });
  });
});
