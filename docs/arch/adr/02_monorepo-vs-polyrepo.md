# ADR-02: Monorepo over Polyrepo

## Context

The project spans a frontend, multiple backend services (per ADR-01), infra config, and architecture documentation (ADRs, microservice design, DDD docs) that frequently need to change together — e.g. adding a field can touch a service's API, its docs, and the frontend that consumes it in the same change.

The team is 5 people working on one semester-long project. The coordination overhead of separate repos (independent versioning, cross-repo PRs, release orchestration) outweighs the isolation benefits a polyrepo would give a larger or longer-lived organization.

Project conventions already assume a single-repo layout: `docs/arch/adr/`, `docs/arch/microservices-design/`, and `docs/ddd/` living alongside the frontend, backend, and infra code.

## Decision

Use a single monorepo for all services, the frontend, infra config, and architecture/DDD documentation.

## Status

Accepted

## Consequences

**Positive**
- Atomic commits/PRs across service, frontend, and doc changes.
- One CI/CD pipeline and one place to browse the whole system; easy for a small team to stay in sync.
- Architecture docs stay versioned next to the code that motivated them.

**Negative**
- The repo will grow large across all services; CI will eventually need path-based filtering to avoid rebuilding/testing everything on every change.
- No repo-boundary enforcement of service isolation — the team must rely on code review and CI checks (e.g. import-boundary linting) instead, once services exist.
