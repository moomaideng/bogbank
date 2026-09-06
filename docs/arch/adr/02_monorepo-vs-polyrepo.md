# ADR-02: Monorepo vs Polyrepo

## Context

The project has a frontend, several backend services (per ADR-01), infra config, and architecture docs (ADRs, microservice design, usecase diagram). These often need to change together. For example, adding a field can touch a service API, its docs, and the frontend in the same change.

Options considered:
- **Monorepo**: all services, the frontend, infra, and docs live in one repo, so one PR can change them together, though CI may later need path-based filtering.
- **Polyrepo**: one repo per service for stronger isolation. That pays off for bigger or longer-running organizations, but extra versioning, cross-repo PRs, and coordinated releases outweigh that benefit for a 5-person semester project.

## Decision

Use one monorepo for all services, the frontend, infra config, and architecture/DDD docs.

## Status

Accepted

## Consequences

**Positive**
- One commit or PR can change a service, the frontend, and the docs together.
- One CI/CD pipeline and one place to browse the whole system. That is easy for a small team to stay in sync.
- Architecture docs stay versioned next to the code.

**Negative**
- The repo will grow large across all services. CI will later need path-based filtering so we do not rebuild and test everything on every change.
- The repo boundary does not enforce service isolation. The team must rely on code review and CI checks instead, once services exist.
