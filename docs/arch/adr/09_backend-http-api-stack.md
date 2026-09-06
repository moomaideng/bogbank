# ADR-09: chi + Huma for HTTP and API Contracts

## Context

Backend services use Go (ADR-08) and REST (ADR-01). Each service needs a documented HTTP API for the Service–Operations–Collaborators deliverable. The team wants to keep the v1 stack simple while avoiding hand-maintained specs that drift from the code.

Options considered:
- **chi + hand-written OpenAPI**: lightweight routing on `net/http` with a spec maintained separately; simplest dependency set, but easy for code and docs to diverge.
- **chi + Huma**: chi for routing on `net/http`, Huma for typed handlers, validation, and OpenAPI 3.1 generated from Go types via the official `humachi` adapter.
- **Fuego**: a complete framework on `net/http` with its own routing and OpenAPI 3.0 generation; replaces chi rather than layering on it.
- **gin**: popular and full-featured, but heavier than needed and a different ecosystem from Huma's `humachi` path.

## Decision

Use chi as the HTTP router. Layer Huma on top through the official [`humachi`](github.com/danielgtaylor/huma/v2/adapters/humachi) adapter for OpenAPI-first request and response definitions, validation, and generated API docs.

## Status

Accepted

## Consequences

**Positive**
- A `net/http`-compatible router keeps middleware and ecosystem compatibility broad.
- Huma generates OpenAPI 3.1 specs from Go types. Those specs serve as the service API docs the Microservice Design deliverable needs.
- Smaller footprint than gin.
- The official adapter means no unsupported glue code between chi and Huma.

**Negative**
- Huma has a smaller community than gin or Fuego. There are fewer tutorials and answers if the team gets stuck.
- Two libraries (chi and Huma) to learn instead of one all-in-one framework.
- The team must follow Huma's generics-based type-to-schema conventions. That is a small learning curve.
