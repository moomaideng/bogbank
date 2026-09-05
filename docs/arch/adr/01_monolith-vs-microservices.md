# ADR-01: Adopt Microservices Architecture over Monolith

## Context

This is a Software Architecture course term project. The syllabus and graded deliverables (Service–Operations–Collaborators table, per-service design and APIs, DDD-based decomposition per the `02_02_Decomposition_Strategies` and `04_02_DDD_2026` course material) require the team to practice designing and building a microservices system, not just a monolith.

The actual product (a personal expense tracker for college students, project name TBD, built by group หมูไม่เด้ง) is functionally small and cohesive at v1 scope: auth, bank-hook receipt ingestion, income/expense/category CRUD, and a dashboard. On business merits alone, a monolith would be simpler to build and operate for a 5-person team shipping in one semester.

Team size is 5, which is small for maintaining multiple independently deployable services (build, deploy, and observability overhead scale with service count).

## Decision

Adopt a microservices architecture, decomposing the system by business capability using DDD-derived bounded contexts (see `docs/arch/microservices-design/` and `docs/ddd/` once written), primarily to satisfy the course's architectural practice requirements.

Mitigate the team-size/overhead mismatch by keeping the decomposition coarse-grained (avoid over-fragmenting services beyond what business capability boundaries justify) and by pairing this decision with a monorepo (ADR-02) and a simple v1 technology stack — single REST style, single database technology, no message broker — as the course's own guidance explicitly allows for a first architecture iteration.

## Status

Accepted

## Consequences

**Positive**
- Satisfies the graded deliverable requiring service decomposition, service APIs, and collaboration diagrams.
- Forces explicit ownership of business capabilities and data per service, which is good practice regardless of product scale.

**Negative**
- Added operational complexity (inter-service calls, per-service deploys, distributed data ownership) that is disproportionate to the actual v1 business scale.
- Slower local development loop than a single monolith app.
- Some of this complexity is deliberately deferred (reverse proxy, dev/prod environment, internal IPC style, message broker, observability stack) until service boundaries are finalized in the Microservice Design step.
