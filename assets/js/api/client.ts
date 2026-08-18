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
  if (init.signal) {
    if (init.signal.aborted) {
      controller.abort();
    } else {
      init.signal.addEventListener('abort', () => controller.abort(), { once: true });
    }
  }

  try {
    const response = await fetchImpl(`${baseUrl}${path}`, {
      ...init,
      signal: controller.signal,
      credentials: 'include',
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
      if (init.signal?.aborted) {
        throw new ApiError('aborted', 'Audit cancelled.', 0);
      }
      throw new ApiError('timeout', 'The request timed out. Please try again.', 0);
    }
    throw new ApiError('network', 'Could not reach the API. Is the server running?', 0);
  } finally {
    clearTimeout(timer);
  }
}

export function createClient(options: ClientOptions = {}) {
  return {
    createAudit(url: string, init: { signal?: AbortSignal } = {}) {
      return request<Audit>(
        '/api/v1/audits',
        { method: 'POST', body: JSON.stringify({ url }), signal: init.signal },
        { ...options, timeoutMs: options.timeoutMs ?? 75000 },
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
    login(password: string) {
      return request<{ ok: boolean }>(
        '/api/v1/session',
        { method: 'POST', body: JSON.stringify({ password }) },
        options,
      );
    },
    logout() {
      return request<null>('/api/v1/session', { method: 'DELETE' }, options);
    },
    session() {
      return request<{ authRequired: boolean; authenticated: boolean }>(
        '/api/v1/session',
        { method: 'GET' },
        options,
      );
    },
  };
}

export type ApiClient = ReturnType<typeof createClient>;
