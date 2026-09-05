# ADR-04: PostgreSQL as Database Technology

## Context

Core domain data is relational: income/expense records reference a category and, optionally, a bank hook, with clear foreign-key relationships. Dashboard queries (UC-03) are aggregate-style — group and sum by category, time range, and bank provider. Nothing in the current scope is genuinely document-shaped/schema-less, and query volume is nowhere near OLAP scale.

Each microservice still owns its own schema/database instance (ADR-01) — this ADR is about which single database *technology* the team standardizes on for v1, not about sharing one database across services. The course explicitly allows a single database technology in the first architecture iteration.

Options considered: PostgreSQL, MongoDB, ClickHouse.

## Decision

Use PostgreSQL as the database technology for every service's data store in this iteration.

## Status

Accepted

## Consequences

**Positive**
- Mature relational/foreign-key support fits the domain model directly.
- Strong SQL aggregate support for dashboard group-by queries.
- One technology for the whole team to operate and learn.
- Wide, mature Go driver support (e.g. pgx), pairing well with the backend language choice.

**Negative**
- Gives up MongoDB's schema flexibility (not needed here) and ClickHouse's OLAP performance (not needed at this scale).
- If dashboard query volume grows significantly later, a read-optimized store may be needed — deferred, not a v1 concern.
