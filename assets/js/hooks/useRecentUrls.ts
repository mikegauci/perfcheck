const KEY = 'perfcheck.recentUrls';
const MAX = 5;

export function readRecentUrls(): string[] {
  try {
    const raw = window.localStorage.getItem(KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw) as unknown;
    return Array.isArray(parsed) ? parsed.filter((v) => typeof v === 'string') : [];
  } catch {
    return [];
  }
}

export function rememberUrl(url: string): string[] {
  const next = [url, ...readRecentUrls().filter((item) => item !== url)].slice(0, MAX);
  try {
    window.localStorage.setItem(KEY, JSON.stringify(next));
  } catch {
    // Safari private mode may throw.
  }
  return next;
}
