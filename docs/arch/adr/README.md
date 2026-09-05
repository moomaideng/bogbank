# Architecture Decision Records

This directory contains the Architecture Decision Records (ADRs). Each ADR documents one significant architectural decision using the Nygard template: **Context**, **Decision**, **Status**, **Consequences**.

ADRs are numbered sequentially in the order they are decided and are never renumbered or edited after acceptance. A reversed decision gets a new ADR that supersedes the old one (mark the old one's Status as `Superseded by ADR-NN`).

## Index

| # | Title | Status |
|---|-------|--------|
| [01](./01_monolith-vs-microservices.md) | Adopt Microservices Architecture over Monolith | Accepted |
| [02](./02_monorepo-vs-polyrepo.md) | Monorepo over Polyrepo | Accepted |
| [03](./03_auth-pattern.md) | Google OAuth as Sole Auth Provider | Accepted |
| [04](./04_database-choice.md) | PostgreSQL as Database Technology | Accepted |
| [05](./05_frontend-language-framework.md) | TypeScript + Next.js for Frontend | Accepted |
| [06](./06_object-storage.md) | RustFS (Local) / S3-Compatible (Prod) Object Storage | Accepted |
| [07](./07_backend-language-framework.md) | Go + chi + Huma for Backend | Accepted |
| [08](./08_backend-coding-style.md) | Hexagonal Architecture (Ports & Adapters) for Backend Services | Accepted |

## Deferred (Tier B)

Not decided yet. Wait until service boundaries exist (see `docs/arch/microservices-design/`): reverse proxy, dev/prod environment, internal IPC style, message broker, observability stack.
