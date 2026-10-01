# ADR-19: Dashboard Read Model and CQRS

## Context

Ledger owns income and expense records and must preserve transactional edits. Dashboard queries aggregate those records by category, time range, type, and bank provider. ADR-12 selects ClickHouse for a derived dashboard store, but a separate decision is needed for source-of-truth ownership, synchronization, failure behavior, and rebuilds.

Options considered:
- **Query Ledger PostgreSQL directly**: simplest and fast enough at the expected scale, but couples dashboard query load and release changes to the transactional service.
- **Move old rows from PostgreSQL into ClickHouse**: reduces PostgreSQL size, but makes common date ranges span two stores and makes old edits and outages unsafe.
- **Keep PostgreSQL authoritative and project events into ClickHouse**: duplicates data and introduces eventual consistency, but isolates dashboard reads and makes the analytical store disposable.

## Decision

Use **CQRS with a derived Dashboard read model**.

Ledger PostgreSQL retains every financial record and remains the only system of record. Ledger publishes created, updated, and deleted record events through the transactional outbox defined by ADR-14. Dashboard Service consumes those events idempotently and projects them into ClickHouse tables and rollups.

ClickHouse data is eventually consistent, disposable, and rebuildable from authoritative Ledger data and retained events. No age-based process moves or deletes records from PostgreSQL. If ClickHouse is unavailable or its projection is not ready, Dashboard Service falls back to `Ledger.GetRecords` over gRPC and aggregates the result in memory for the supported date range.

S3 or Iceberg archival and hot/cold financial-record tiering are out of scope at the expected data volume.

## Status

Accepted. Refines the ClickHouse role selected in ADR-12 and relies on ADR-14 for event delivery; it supersedes neither ADR.

## Consequences

**Positive**
- Financial records remain transactionally editable and recoverable in one authoritative store.
- Dashboard query load and analytical schemas evolve independently from Ledger's write model.
- ClickHouse corruption or loss can be repaired by rebuilding instead of restoring it as primary data.
- The fallback path provides degraded dashboard service while ClickHouse is unavailable.

**Negative**
- Dashboard results may briefly lag behind a Ledger write.
- Projection code, replay tooling, schema evolution, and reconciliation tests are required.
- The fallback path must be kept compatible and tested rather than treated as unused code.
- ClickHouse may be more operational complexity than the project's data volume alone requires.
