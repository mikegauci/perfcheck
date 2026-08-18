---
title: "About"
description: "Why PerfCheck exists — a portfolio project for learning Hugo and Go with senior frontend practices."
layout: "page"
---

PerfCheck is a portfolio project. It exists to show that I can pick up **Hugo** and **Go** quickly while applying the frontend practices I already use day to day: responsive layout, TypeScript, SASS with BEM, accessibility, SEO, testing and a maintainable architecture.

## Design goals

- **Small enough to understand completely** — no framework soup.
- **Customer-facing pages** for marketing, plus a small **internal dashboard**.
- **Separation of concerns** — Hugo owns documents, React owns interactive regions, Go owns data.
- **Git-friendly incremental development** — conventional commits, phase branches, green `main` at every step.

## What it is not

It is not a production auditing product. Scores come from one live HTML fetch and table-driven signals — not Lighthouse or real-user monitoring. The dashboard can sit behind an optional password. History lives in memory or SQLite. Those limits are deliberate — each is listed under “what I’d add next” in the project documentation.
