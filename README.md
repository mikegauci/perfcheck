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

## Prerequisites

- [Hugo Extended](https://gohugo.io/installation/) ≥ 0.146
- [Go](https://go.dev/dl/) ≥ 1.22
- [Dart Sass](https://sass-lang.com/install/) on `PATH` (`brew install sass/sass/sass`)
- Node.js ≥ 20 (`nvm use` if you use nvm)

## Quick start

```bash
npm install
make api     # terminal 1 — API on :8080
make hugo    # terminal 2 — site on :1313
```

Open http://localhost:1313/.

## Tests

```bash
make test
```

## Project intent

Built as a portfolio piece for a senior frontend role that lists Hugo/Go as desirable experience. The site is intentionally small enough to understand completely: two React islands, one Go API, deterministic mock audits (no real Lighthouse integration).

See the decisions log in later documentation commits for what was deliberately left out and why.
