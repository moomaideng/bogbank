# ADR-12: Database Technology v2

## Context

ADR-04 selected PostgreSQL for every service. The service boundaries are now defined, and they introduce three different data-access needs: transactional financial records, a rebuildable dashboard read model, and short-lived coordination data for asynchronous consumers.

Ledger and receipt data need relational constraints and transactions. Dashboard queries repeatedly aggregate by user, category, date, and bank provider. Redpanda (ADR-13) consumers also need a fast store for idempotency keys so at-least-once delivery does not apply the same event twice. The course additionally requires both an RDBMS and a NoSQL database.

Options considered:
- **PostgreSQL only**: the simplest operational choice and sufficient for the expected data volume, but it does not provide the intended independent dashboard projection or a suitable shared idempotency store.
- **PostgreSQL and Redis**: satisfies the RDBMS/NoSQL requirement and covers transactional data plus idempotency, but dashboard traffic remains coupled to Ledger's operational store.
- **PostgreSQL, ClickHouse, and Redis**: separates transactional ownership, analytical projections, and ephemeral coordination data, at the cost of operating three database technologies.

## Decision

Use each database for one explicit role:

- **PostgreSQL** is the authoritative store for service-owned transactional data. Auth, Receipt, and Ledger own private schemas or databases. Other services never read those stores directly.
- **ClickHouse** stores only the Dashboard Service's derived, rebuildable read model. It is not a system of record and is governed by ADR-19.
- **Redis** stores consumer idempotency keys and bounded caches or session data. It must not be the only copy of financial or receipt data.

Development may run one PostgreSQL instance with a separate schema and credentials per service to reduce laptop resource use without weakening logical ownership.

## Status

Accepted. Supersedes ADR-04 and the one-database-technology constraint in ADR-01.

## Consequences

**Positive**
- Transactional data keeps PostgreSQL's constraints and atomic updates.
- Dashboard projections can be rebuilt without risking the financial system of record.
- Redis gives event consumers a low-latency idempotency store and satisfies the required NoSQL category.
- Each datastore has a narrow responsibility instead of becoming a shared integration database.

**Negative**
- The team must operate, back up, monitor, and learn three database technologies.
- Data projected into ClickHouse is eventually consistent with Ledger.
- Local development and CI require more containers and health checks.
- Redis loss can cause events to be processed again, so business-level unique constraints are still required.
