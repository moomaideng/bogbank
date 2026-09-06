# ADR-07: RustFS (Local) / S3-Compatible (Prod) Object Storage

## Context

Receipt screenshots shared by users (UC-01) need to be stored as objects. The backend stores them and links them from expense records. This choice only affects local dev. Trying RustFS is low risk. Switching to another S3-compatible server later is a Docker Compose change, not an application code change. The app only talks to the S3 API.

Local dev options compared:
- **MinIO**: was the most mature S3-compatible server, with a large community and docs, but [minio/minio](https://github.com/minio/minio) is archived and no longer actively maintained (read-only since April 2026).
- **LocalStack**: emulates more AWS services, and is heavier than we need because we only use the S3 API.
- **RustFS**: a small S3-compatible binary that is fast for receipt-sized objects, newer, with a smaller community than MinIO.

Production options (deferred to the production-environment ADR):
- **AWS S3**: the original S3 API with the most tooling, billed by usage.
- **Cloudflare R2**: S3-compatible storage that does not charge extra when files leave, with a smaller ecosystem than AWS S3.
- **Google Cloud Storage**: S3-compatible if we later host on Google Cloud, with extra mapping compared with native AWS S3.

## Decision

Use RustFS as the local/dev S3-compatible object store. Use a real S3-compatible cloud provider in production. The exact provider is deferred to the production-environment ADR. Application code depends only on the S3 API. It never depends on a vendor-specific SDK feature. The backend stays easy to swap in either environment.

## Status

Accepted

## Consequences

**Positive**
- Fast, lightweight local dev loop: single binary, Docker Compose.
- No code lock-in to RustFS. The app uses the S3 API only.

**Negative**
- RustFS is young. It has a smaller community and fewer troubleshooting docs than MinIO had.
- If it causes friction, falling back to another S3-compatible server is cheap because the API is the same. MinIO is a poor fallback now that its open-source repo is archived. The team should still budget time for that switch if it happens.
