# ADR-07: Go + chi + Huma for Backend

## Context

**Language**: Options considered were Go, Rust, and TypeScript (Node). Go offers fast build/run cycles, simple concurrency, and is the only language the course's own example framework list (Huma/chi/fiber/echo/gin) targets.

**Router**: Among chi, gin, fiber, and echo, the team wants a lightweight, `net/http`-compatible router — chi fits that directly, unlike fiber (built on `fasthttp`, not `net/http`-compatible) and gin (heavier feature set than needed here).

**API framework**: The course explicitly calls out "Defining a public API contract (e.g., OpenAPI spec)" as an architecturally-significant interface decision, and the Microservice Design deliverable requires documented service APIs. Two OpenAPI-code-first Go options were compared:

- **Fuego** — a complete, standalone framework built directly on `net/http` with its own routing; generates OpenAPI 3.0.
- **Huma** — layers onto an existing router via an adapter interface (chi, gin, echo, fiber, stdlib, etc.) rather than owning routing itself; generates OpenAPI 3.1 and JSON Schema from Go types via generics.

Since chi is already the router choice, Huma's official `humachi` adapter (`github.com/danielgtaylor/huma/v2/adapters/humachi`) lets the team keep chi's routing and add OpenAPI-first request/response typing on top. Picking Fuego instead would mean replacing chi's routing model entirely with Fuego's own — an unnecessary swap once chi is already chosen.

## Decision

Use Go for all backend services. Use chi as the HTTP router, with Huma layered on top via the official `humachi` adapter for OpenAPI-first request/response definitions, validation, and generated API docs.

## Status

Accepted

## Consequences

**Positive**
- `net/http`-compatible router keeps middleware/ecosystem compatibility broad.
- Huma generates OpenAPI 3.1 specs directly from Go types, which doubles as the service-API documentation the Microservice Design deliverable needs.
- Smaller footprint than gin/fiber.
- Official adapter means no unsupported glue code between chi and Huma.

**Negative**
- Huma has a smaller community than gin or Fuego — fewer tutorials/answers if the team gets stuck.
- Two libraries (chi + Huma) to learn instead of one all-in-one framework.
- Team must follow Huma's generics-based type-to-schema conventions — a small learning curve.
