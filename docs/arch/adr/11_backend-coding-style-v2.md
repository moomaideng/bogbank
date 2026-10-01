# ADR-11: Backend Coding Style v2

## Context

ADR-10 chose hexagonal ports and adapters so each service's domain core would depend only on interfaces, with HTTP, Postgres, and object storage as adapters. That isolation matches the course's DDD material. It also adds a set of port files per feature. Most Bogbank operations are small CRUD flows (for example Manage Category). The semester is short, and the extra files slow the team down more than they protect the design.

The part of ADR-10 worth keeping is the outbound boundary: a usecase should not import a driver, so tests can run without HTTP, Postgres, or object storage, and an adapter can be swapped (ADR-07) without rewriting the usecase.

CQRS stays out of the service template. A derived read model, if the dashboard needs one later, belongs to that service's own decision.

## Decision

Use a 3-tier layout in every backend service: handler, then usecase, then adapter. The usecase owns the interfaces it calls. Adapters implement those interfaces.

## Status

Accepted

## Consequences

**Positive**
- One folder layout for every service, with fewer files per feature than ADR-10.
- Usecase tests still need no database, HTTP server, or object storage.
- Swapping an adapter still leaves the usecase alone.

**Negative**
- Nothing in the compiler stops a handler from talking to an adapter directly. Review has to keep that path handler → usecase → adapter.
