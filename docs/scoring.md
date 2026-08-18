# Scoring engines

PerfCheck never pretends a number is more real than it is. Every audit records an `engine` field so the UI can label provenance.

## `psi` — PageSpeed Insights

Calls `https://www.googleapis.com/pagespeedonline/v5/runPagespeed` with `strategy=mobile` and the Performance, SEO and Accessibility categories.

- Scores are Lighthouse category scores scaled 0–100.
- Optional `PSI_API_KEY` raises the anonymous quota.
- Typical latency is 10–30 seconds. Results are cached in process for 15 minutes.
- **Limits:** Google’s lab environment, not your users’ devices. Rate-limited without a key. Does not crawl authenticated pages.

## `fetch` — live HTML GET

One HTTP GET (10s timeout, 2 MiB cap, 5 redirects) plus `golang.org/x/net/html` inspection.

Deductions from 100 per category are table-driven in `scoreFromSignals` (HTTPS, TTFB, bytes, title, meta description, H1 count, lang, alt text, labels, viewport). Overall is 40% performance, 30% SEO, 30% accessibility.

- **Limits:** A single document, no JS execution, no real Core Web Vitals. Useful when PSI is unavailable.

## `mock` — deterministic hash

FNV-1a of the normalised URL seeds `math/rand`. Same URL always yields the same scores. No network.

- **Limits:** The numbers are synthetic. Kept as the last fallback so a demo never hard-fails.

## Chain (`-scorer=auto`)

`PSIScorer` → `FetchScorer` → `MockScorer`. The first success wins. The dashboard badge tells you which one answered.
