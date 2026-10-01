# ADR-10: Backend Coding Style

## Context

The course's DDD material teaches bounded contexts and domain modeling. That pairs well with keeping each service's domain core away from infrastructure. Infrastructure includes the HTTP framework, the DB driver, and the object storage client. Options considered:

- **3-tier layered**: simplest and fastest for small CRUD-heavy services, but business logic tends to leak into handlers or the DB layer.
- **Clean Architecture**: similar isolation to hexagonal, with a stricter layer rule and more boilerplate than we need.
- **Hexagonal (ports and adapters)**: the domain core depends only on interfaces (ports), frameworks are adapters, and there are more files per feature than 3-tier.

## Decision

Use hexagonal architecture (ports and adapters) for each backend service. A domain core depends only on ports. HTTP handlers (ADR-09), the Postgres repository, and the object storage client (ADR-07) are adapters that implement those ports.

## Status

Superseded by ADR-11

## Consequences

**Positive**
- Domain logic stays testable without starting HTTP, a database, or object storage.
- Swapping an adapter (for example the object storage provider per ADR-07) does not touch domain code.
- Aligns with the DDD concepts the course is grading.

**Negative**
- More files and interfaces per feature than a 3-tier service. That extra structure is real overhead for simple CRUD operations (for example Manage Category).
- The team needs to agree on one port/adapter folder layout and use it across all services. Otherwise the benefit is lost.
