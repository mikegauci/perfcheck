---
title: "How it works"
description: "How PerfCheck produces mock performance, SEO and accessibility scores from a URL."
layout: "page"
---

PerfCheck is intentionally small. You submit a URL, the Go API validates it, and a deterministic mock scorer returns scores plus recommendations. There is no real Lighthouse or PageSpeed call — the point is the architecture, not synthetic third-party latency.

## The flow

1. **You enter a URL** on the audit page. The form uses real HTML semantics and accessible validation.
2. **Hugo serves the page** as static HTML. React mounts only on the audit island.
3. **The Go API** validates the URL, hashes it, and derives repeatable scores.
4. **Recommendations** come from a rule catalogue — adding a rule means appending a condition, not editing a switch.
5. **The dashboard** lists previous audits from an in-memory store behind a Repository interface.

## Why mock scores?

A portfolio project that depends on Google APIs, rate limits and API keys is harder to run and review. Deterministic mocks keep demos repeatable and tests golden. The README documents this trade-off explicitly.

## What this demonstrates

- Static HTML first (SEO and performance)
- React islands instead of an SPA
- BEM components shared between Hugo partials and React
- A tiny Go REST API with clear package boundaries
