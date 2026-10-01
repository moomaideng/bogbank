# Architecture Decision Records

This directory contains the Architecture Decision Records (ADRs). Each ADR documents one significant architectural decision using the Nygard template: **Context**, **Decision**, **Status**, **Consequences**.

ADRs are numbered in the order they are decided. They are never renumbered. They are not rewritten after acceptance. A reversed decision gets a new ADR that supersedes the old one. Mark the old one's Status as `Superseded by ADR-NN`.

## Index

| # | Title | Status |
|---|-------|--------|
| [01](./01_monolith-vs-microservices.md) | Monolith vs Microservices | Accepted |
| [02](./02_monorepo-vs-polyrepo.md) | Monorepo vs Polyrepo | Accepted |
| [03](./03_auth-pattern.md) | Authentication Pattern | Accepted |
| [04](./04_database-choice.md) | Database Technology | Accepted |
| [05](./05_client-delivery.md) | Client Delivery | Accepted |
| [06](./06_frontend-language-framework.md) | Frontend Language & Framework | Accepted |
| [07](./07_object-storage.md) | Object Storage | Accepted |
| [08](./08_backend-language.md) | Backend Language | Accepted |
| [09](./09_backend-http-api-stack.md) | Backend HTTP & API Stack | Accepted |
| [10](./10_backend-coding-style.md) | Backend Coding Style | Superseded by ADR-11 |
| [11](./11_backend-coding-style-v2.md) | Backend Coding Style v2 | Accepted |

## Deferred (Tier B)

Not decided yet. Wait until service boundaries exist (see `docs/arch/microservices-design/`). That includes reverse proxy, dev/prod environment, internal IPC style, message broker, and observability stack.
