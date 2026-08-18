# PerfCheck

Website performance, SEO and accessibility audits. A portfolio project demonstrating **Hugo Extended**, **Go**, **React**, **TypeScript** and **SASS/BEM**.

Scores come from a pluggable `Scorer` chain: Google PageSpeed Insights first, a live HTML fetch second, and a deterministic mock last. Every result is labelled with the engine that produced it.

## Stack

| Layer | Technology |
| --- | --- |
| Marketing site | Hugo Extended |
| Interactive islands | React + TypeScript (esbuild via Hugo `js.Build`) |
| Styles | SASS (Dart Sass) + BEM |
| API | Go `net/http` |
| Persistence | SQLite (`modernc.org/sqlite`) or in-memory |
| Frontend tests | Vitest + React Testing Library |
| Backend tests | Go `testing` package |

## Scoring engines

| Engine | Data source | Network | When it runs |
| --- | --- | --- | --- |
| `psi` | Google PageSpeed Insights (Lighthouse lab) | PSI API | First choice in `auto` |
| `fetch` | One HTML GET + `x/net/html` signals | Target URL | PSI disabled or error |
| `mock` | Deterministic FNV hash of the URL | None | Last fallback so demos never hard-fail |

Every audit response includes `engine` so the UI can show provenance. Details: [docs/scoring.md](docs/scoring.md).

## Architecture

```text
Browser
  ├─ Hugo static HTML/CSS (marketing pages, zero JS)
  ├─ React islands on /audit/, /compare/, /result/ and /dashboard/
  └─ fetch ──► Go API (/api/v1/audits)
                  ├─ validate URL
                  ├─ TTL cache
                  ├─ ChainScorer: PSI → fetch → mock
                  ├─ recommendation rules
                  └─ Repository (memory or SQLite)
```

- **Hugo owns the document** — content and metadata are in the HTML before JS runs.
- **React owns interactive regions** — islands mount from `[data-react-island]`.
- **Go owns the data** — scoring, persistence and auth.

## Prerequisites

- Hugo Extended ≥ 0.146
- Go ≥ 1.25
- Dart Sass on `PATH`
- Node.js ≥ 20

Optional: `PSI_API_KEY` for a higher PageSpeed Insights quota.

## Quick start

```bash
npm install
make api     # terminal 1 — API on :8080
make hugo    # terminal 2 — site on :1313
```

Open http://localhost:1313/.

Useful flags:

```bash
cd api && go run ./cmd/perfcheckd \
  -scorer=auto \
  -db ../data/perfcheck.db \
  -dashboard-password=secret \
  -cors-origin=http://localhost:1313
```

`-scorer` is `auto` (default), `psi`, `fetch` or `mock`.

Production-style single process (after `hugo --minify`):

```bash
cd api && go run ./cmd/perfcheckd -cors-origin="" -static ../public -db ../data/perfcheck.db
```

## Tests

```bash
make test
make typecheck
```

Network is never required for tests. PSI and HTML fetch use `httptest` fixtures.

## API

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/api/v1/healthz` | Liveness |
| `POST` | `/api/v1/audits` | Create audit `{ "url": "https://…" }` |
| `GET` | `/api/v1/audits?limit=20` | List history (optional session) |
| `GET` | `/api/v1/audits/{id}` | Fetch one audit |
| `POST` | `/api/v1/session` | Sign in `{ "password": "…" }` |
| `DELETE` | `/api/v1/session` | Sign out |
| `GET` | `/api/v1/session` | Session status |

Errors use a single envelope:

```json
{ "error": { "code": "invalid_url", "message": "…", "field": "url" } }
```

## Dependency budget

Frontend runtime: `react`, `react-dom`.

Go direct dependencies:

| Module | Why |
| --- | --- |
| `golang.org/x/net/html` | Correct HTML parsing for the fetch scorer |
| `modernc.org/sqlite` | Pure-Go SQLite so `CGO_ENABLED=0` still builds |

No HTTP framework, no ORM, no CSS framework, no client state library.

## Decisions log

| Decision | Choice | Why |
| --- | --- | --- |
| Scoring | `Scorer` interface + chain | DIP; mock/fetch/PSI are drop-in |
| Provenance | `engine` on every audit | Fake numbers without a label would be dishonest |
| Persistence | SQLite behind `Repository` | Original interface paid off; empty `-db` stays in-memory |
| Bundling | Hugo `js.Build` (esbuild) | One build pipeline |
| CSS | Dart Sass + BEM, shared with React | One class vocabulary |
| Router | Go 1.22 `ServeMux` | No framework |
| Dashboard auth | Optional HMAC cookie | Off by default for local demos |
| Real lab data | PageSpeed Insights API | No Chrome in the image |

See [docs/scoring.md](docs/scoring.md) and [docs/sqlite.md](docs/sqlite.md).

## Performance & accessibility targets

- Marketing pages: Lighthouse 95+ on all four categories.
- React loads only on island pages.
- WCAG 2.1 AA intent: skip link, labels, live regions, focus-visible, `prefers-reduced-motion`.

See [docs/cross-browser.md](docs/cross-browser.md).

## Licence

MIT — portfolio / demonstration use.
