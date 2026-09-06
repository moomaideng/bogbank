# ADR-08: Backend Language

## Context

The system has several backend microservices (ADR-01). They expose REST APIs to the frontend (ADR-06). The team has five people and one semester. The backend language should be quick to build in and easy for the whole team to learn.

Options considered:
- **Go**: fast compile and run cycles, simple concurrency, and a strong fit for small HTTP services.
- **Rust**: very safe and fast at runtime, but slower to learn and compile for this timeline.
- **TypeScript (Node)**: same language as the frontend (ADR-06), but a separate runtime from the browser stack and less natural for compiled microservices on the team.

## Decision

Use Go for all backend services.

## Status

Accepted

## Consequences

**Positive**
- Fast feedback loop when building and running many small services locally.
- Simple concurrency model for services that handle parallel HTTP requests.
- One backend language for the whole team to standardize on.

**Negative**
- Since it's a different language from the TypeScript frontend, API contracts must stay explicit (OpenAPI or shared types at the boundary).
