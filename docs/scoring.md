# Scoring

PerfCheck scores a URL from one live HTML fetch. Every audit records `engine: "fetch"` so the UI can label provenance.

## `fetch` — live HTML GET

One HTTP GET (10s timeout, 2 MiB cap, 5 redirects) plus `golang.org/x/net/html` inspection.

Deductions from 100 per category are table-driven in `scoreFromSignals` (HTTPS, TTFB, bytes, title, meta description, H1 count, lang, alt text, labels, viewport). Overall is 40% performance, 30% SEO, 30% accessibility.

Results are cached in process for 15 minutes by default (`-cache-ttl`).

- **Limits:** A single document, no JS execution, no real Core Web Vitals. Authenticated or heavily client-rendered pages will look incomplete. A failed fetch surfaces as an API error.
