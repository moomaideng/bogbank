# ADR-08: Hexagonal Architecture (Ports & Adapters) for Backend Services

## Context

The course's DDD material teaches bounded contexts and domain modeling, which pairs naturally with keeping each service's domain core isolated from infrastructure (HTTP framework, DB driver, object storage client) behind explicit interfaces. Options considered:

- **3-tier layered** — simplest, fastest for small CRUD-heavy services, but doesn't enforce domain isolation; business logic tends to leak into handlers or the DB layer.
- **Clean Architecture** — similar isolation goals to hexagonal, but with a stricter, more layered dependency rule and more boilerplate than needed for services this small.
- **Hexagonal (ports & adapters)** — domain core depends only on interfaces (ports); frameworks and infrastructure are adapters implementing those ports.

## Decision

Use hexagonal architecture (ports & adapters) for each backend service: a domain core (entities + use-case/application services) depends only on ports (interfaces); HTTP handlers (Huma/chi), the Postgres repository, and the object storage client (ADR-06) are adapters implementing those ports.

## Status

Accepted

## Consequences

**Positive**
- Domain logic stays testable without spinning up HTTP, a database, or object storage.
- Swapping an adapter (e.g. the object storage provider per ADR-06) doesn't touch domain code.
- Aligns with the DDD concepts the course is grading.

**Negative**
- More files/interfaces per feature than a 3-tier service — real ceremony overhead for genuinely simple CRUD operations (e.g. Manage Category).
- The team needs to agree on and consistently follow one port/adapter folder layout across all services to actually get the benefit.
