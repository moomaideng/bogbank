# ADR-05: TypeScript + Next.js for Frontend

## Context

The frontend needs a responsive, mobile-first UI (dashboard charts, forms) and benefits from a productive, typed component ecosystem for a small team iterating over one semester. The backend is already committed to a statically typed language (Go). Options considered: TypeScript + Next.js, TypeScript + SvelteKit, plain JavaScript.

## Decision

Use TypeScript with Next.js for the frontend.

## Status

Accepted

## Consequences

**Positive**
- Type safety catches mismatches against backend API contracts at compile time.
- Next.js's large ecosystem and community mean faster problem-solving on a one-semester timeline.
- React skills are more commonly already known across the team than Svelte.
- SSR/routing out of the box reduces boilerplate.

**Negative**
- More framework convention/"magic" to learn than a plain SPA.
- Heavier bundle/runtime than Svelte's compiled-away approach.
