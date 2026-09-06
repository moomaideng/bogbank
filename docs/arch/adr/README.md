# Architecture Decision Records

This directory contains the Architecture Decision Records (ADRs). Each ADR documents one significant architectural decision using the Nygard template: **Context**, **Decision**, **Status**, **Consequences**.

ADRs are numbered in the order they are decided. They are never renumbered. They are not rewritten after acceptance. A reversed decision gets a new ADR that supersedes the old one. Mark the old one's Status as `Superseded by ADR-NN`.

## Index

| # | Title | Status |
|---|-------|--------|
| [01](./01_monolith-vs-microservices.md) | Adopt Microservices Architecture over Monolith | Accepted |
| [02](./02_monorepo-vs-polyrepo.md) | Monorepo over Polyrepo | Accepted |
| [03](./03_auth-pattern.md) | Google OAuth as Sole Auth Provider | Accepted |
| [04](./04_database-choice.md) | PostgreSQL as Database Technology | Accepted |
| [05](./05_client-delivery.md) | Mobile-First Web + PWA for Client Delivery | Accepted |
| [06](./06_frontend-language-framework.md) | TypeScript + Next.js for Frontend | Accepted |
| [07](./07_object-storage.md) | RustFS (Local) / S3-Compatible (Prod) Object Storage | Accepted |
| [08](./08_backend-language.md) | Go for Backend Language | Accepted |
| [09](./09_backend-http-api-stack.md) | chi + Huma for HTTP and API Contracts | Accepted |
| [10](./10_backend-coding-style.md) | Hexagonal Architecture (Ports & Adapters) for Backend Services | Accepted |

## Deferred (Tier B)

Not decided yet. Wait until service boundaries exist (see `docs/arch/microservices-design/`). That includes reverse proxy, dev/prod environment, internal IPC style, message broker, and observability stack.
