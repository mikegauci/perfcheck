import { ApiError, type ApiErrorBody, type Audit, type AuditListResponse } from './types';

export type ClientOptions = {
  baseUrl?: string;
  fetchImpl?: typeof fetch;
  timeoutMs?: number;
};

async function request<T>(
  path: string,
  init: RequestInit,
  options: ClientOptions = {},
): Promise<T> {
  const baseUrl = options.baseUrl ?? '';
  const fetchImpl = options.fetchImpl ?? fetch;
  const timeoutMs = options.timeoutMs ?? 10000;
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await fetchImpl(`${baseUrl}${path}`, {
      ...init,
      signal: controller.signal,
      headers: {
        Accept: 'application/json',
        ...(init.body ? { 'Content-Type': 'application/json' } : {}),
        ...init.headers,
      },
    });

    const text = await response.text();
    let data: unknown = null;
    if (text) {
      try {
        data = JSON.parse(text);
      } catch {
        throw new ApiError('unknown', 'Unexpected response from the server.', response.status);
      }
    }

    if (!response.ok) {
      const body = data as ApiErrorBody | null;
      throw new ApiError(
        body?.error?.code ?? 'unknown',
        body?.error?.message ?? 'Request failed.',
        response.status,
        body?.error?.field,
      );
    }

    return data as T;
  } catch (err) {
    if (err instanceof ApiError) {
      throw err;
    }
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw new ApiError('timeout', 'The request timed out. Please try again.', 0);
    }
    throw new ApiError('network', 'Could not reach the API. Is the server running?', 0);
  } finally {
    clearTimeout(timer);
  }
}

export function createClient(options: ClientOptions = {}) {
  return {
    createAudit(url: string) {
      return request<Audit>(
        '/api/v1/audits',
        { method: 'POST', body: JSON.stringify({ url }) },
        options,
      );
    },
    listAudits(limit = 20) {
      return request<AuditListResponse>(
        `/api/v1/audits?limit=${limit}`,
        { method: 'GET' },
        options,
      );
    },
    getAudit(id: string) {
      return request<Audit>(`/api/v1/audits/${id}`, { method: 'GET' }, options);
    },
  };
}

export type ApiClient = ReturnType<typeof createClient>;
