# ADR-06: RustFS (Local) / S3-Compatible (Prod) Object Storage

## Context

Receipt screenshots shared by users (UC-01) need to be persisted as objects that the backend can store and reference from expense records. Local dev options compared:

- **MinIO** — most mature/battle-tested S3-compatible server, huge community and docs.
- **LocalStack** — broader AWS-service emulation, heavier than needed since only the S3 API is used here.
- **RustFS** — newer (public releases since late 2025), full S3 API compatibility, single ~18 MB static binary with no runtime dependencies, faster than MinIO for small objects (receipt images are small/medium), runs via Docker Compose, and is a drop-in swap with MinIO since both speak the same S3 API.

Because only local dev is at stake here, the blast radius of trying RustFS is small: swapping back to MinIO later is a Docker Compose change, not an application code change, since the app only ever talks to the S3 API.

## Decision

Use RustFS as the local/dev S3-compatible object store. Use a real S3-compatible cloud provider in production (exact provider — AWS S3, Cloudflare R2, or GCS via interop — deferred to the production-environment ADR, Tier B). Application code depends only on the S3 API, never a vendor-specific SDK feature, so the backend stays swappable in either environment.

## Status

Accepted

## Consequences

**Positive**
- Fast, lightweight local dev loop: single binary, Docker Compose.
- Zero code lock-in to RustFS specifically, since it's S3-API-only.
- Matches the team's own earlier brainstormed idea (see `todo.md`).

**Negative**
- RustFS is young — smaller community and thinner troubleshooting resources than MinIO.
- If it causes friction, fallback to MinIO is cheap (same API) but is still a context-switch the team should budget time for if it happens.
