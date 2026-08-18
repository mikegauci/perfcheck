# PerfCheck

Mock website performance, SEO and accessibility audits. A small portfolio project demonstrating **Hugo Extended**, **Go**, **React**, **TypeScript** and **SASS/BEM** — built to show senior frontend practices while learning Hugo and Go quickly.

## Stack

| Layer | Technology |
| --- | --- |
| Marketing site | Hugo Extended |
| Interactive islands | React + TypeScript (esbuild via Hugo `js.Build`) |
| Styles | SASS (Dart Sass) + BEM |
| API | Go (`net/http`, stdlib only) |
| Frontend tests | Vitest + React Testing Library |
| Backend tests | Go `testing` package |

## Architecture

```text
Browser
  ├─ Hugo static HTML/CSS (marketing pages, zero JS)
  ├─ React islands on /audit/ and /dashboard/ only
  └─ fetch ──► Go API (/api/v1/audits)
                  ├─ validate URL
                  ├─ deterministic mock scorer
                  ├─ recommendation rules
                  └─ in-memory Repository
```

- **Hugo owns the document** — content and metadata are in the HTML before JS runs.
- **React owns interactive regions** — islands mount from `[data-react-island]`.
- **Go owns the data** — stdlib HTTP service, no third-party modules.

## Prerequisites

- Hugo Extended ≥ 0.146 (`brew install hugo`)
- Go ≥ 1.22 (`brew install go`)
- Dart Sass on `PATH` (`brew install sass/sass/sass`)
- Node.js ≥ 20 (`nvm use`)

## Quick start

```bash
npm install
make api     # terminal 1 — API on :8080 (CORS allows Hugo)
make hugo    # terminal 2 — site on :1313
```

Open http://localhost:1313/.

Production-style single process (after `hugo --minify`):

```bash
cd api && go run ./cmd/perfcheckd -cors-origin="" -static ../public
```

## Tests

```bash
make test          # Go (-race) + Vitest
make typecheck     # tsc --noEmit (esbuild does not typecheck)
make lint
```

## API

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/api/v1/healthz` | Liveness |
| `POST` | `/api/v1/audits` | Create mock audit `{ "url": "https://…" }` |
| `GET` | `/api/v1/audits?limit=20` | List newest audits |
| `GET` | `/api/v1/audits/{id}` | Fetch one audit |

Errors use a single envelope:

```json
{ "error": { "code": "invalid_url", "message": "…", "field": "url" } }
```

## Decisions log

| Decision | Choice | Why |
| --- | --- | --- |
| Storage | In-memory behind `audit.Repository` | Zero setup; interface shows DIP; swap to SQLite later |
| Bundling | Hugo `js.Build` (esbuild) | One build pipeline; Vitest stays for tests |
| CSS | Dart Sass + BEM, shared with React | LibSass deprecated; one class vocabulary |
| Scores | Deterministic FNV hash mock | Repeatable demos and golden tests; no API keys |
| Router | Go 1.22 `ServeMux` | No framework; clear and enough |
| State | Explicit unions in hooks | Impossible states unrepresentable |
| Auth on dashboard | None | Demo of an internal tool, not a secured product |

## What was deliberately left out

No real Lighthouse/PageSpeed integration, no database, no auth, no CSS framework, no Redux/Zustand, no i18n, no blog. Each would grow the surface area without strengthening the learning story.

**What I’d add next:** SQLite behind the same `Repository`, basic auth for `/dashboard/`, and a live Lighthouse adapter behind a `Scorer` interface.

## Performance & accessibility targets

- Marketing pages: Lighthouse 95+ on all four categories (run locally after `hugo --minify`).
- Interactive pages: JS under ~50 KB gzipped; React loaded only when `island` is set.
- WCAG 2.1 AA intent: skip link, labels, live regions, focus-visible, contrast-aware tokens, `prefers-reduced-motion`.

See [docs/cross-browser.md](docs/cross-browser.md) for the manual browser matrix.

## Project layout

See the repository tree — Hugo at the root, Go module in `api/`, React source in `assets/js/`, SASS in `assets/scss/`.

## Licence

MIT — portfolio / demonstration use.
