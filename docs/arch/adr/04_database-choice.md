# ADR-04: Database Technology

## Context

Core domain data is relational. Income and expense records point to a category. They can also point to a bank hook. Those are clear foreign-key relationships. Dashboard queries group and sum by category, time range, and bank provider. Current data is not a document shape that needs a schema-less store. Query volume is not at OLAP scale.

Options considered:
- **PostgreSQL**: a relational SQL database with foreign keys and group-by queries that fits income, expense, and category data.
- **MongoDB**: a document store that is easy to change later, but it is weak for the joins and dashboard sums we need.
- **ClickHouse**: a column store built for large analytics that is faster for huge dashboards and more than we need for a student expense tracker.

## Decision

Use PostgreSQL as the database technology for every service's data store in this iteration.

## Status

Superseded by ADR-12.

## Consequences

**Positive**
- Mature relational and foreign-key support fits the domain model directly.
- Strong SQL aggregate support for dashboard group-by queries.
- One technology for the whole team to run and learn.
- Wide and mature driver support.

**Negative**
- We give up MongoDB's schema flexibility. We do not need that here. We also give up ClickHouse's OLAP performance. We do not need that at this scale.
- If dashboard query volume grows a lot later, we may need a read-optimized store. That is not a v1 concern.
