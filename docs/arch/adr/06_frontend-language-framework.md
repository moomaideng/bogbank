# ADR-06: TypeScript + Next.js for Frontend

## Context

The client is delivered as mobile-first responsive web with PWA support (ADR-05). The frontend needs dashboard charts and forms. A small team working over one semester needs a typed component ecosystem that is fast to work in.

Options considered:
- **TypeScript + Next.js**: typed React with routing and server rendering built in, plus the largest ecosystem, with more framework rules to learn.
- **TypeScript + SvelteKit**: typed, with smaller bundles because Svelte compiles away, and fewer people on the team know Svelte.
- **Plain JavaScript**: no compile step and no type checks against backend API contracts.

## Decision

Use TypeScript with Next.js for the frontend.

## Status

Accepted

## Consequences

**Positive**
- Type safety catches mismatches against backend API contracts at compile time.
- Next.js has a large ecosystem and community. That means faster problem-solving on a one-semester timeline.
- React skills are more common on the team than Svelte.
- SSR and routing come out of the box.

**Negative**
- More framework convention and "magic" to learn than a plain SPA.
- Heavier bundle and runtime than Svelte's compiled approach.
- Next.js is a web stack. A future native app (deferred in ADR-05) would be a separate client.
