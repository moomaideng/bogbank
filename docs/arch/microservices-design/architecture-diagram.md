# Architecture Diagram (v1)

Arrows show request/invocation flow: `A --> B` means **A calls B** (per the course convention — no response arrows drawn, no call chains implied that aren't real).

```mermaid
flowchart LR
    Customer(["Customer\n(college student)"])
    GoogleOAuth[["Google OAuth\n(external identity provider)"]]
    ObjectStorage[["Object Storage\nRustFS (dev) / S3-compatible (prod)\n— ADR-07"]]

    subgraph System["Expense Tracker (project name TBD)"]
        Ingestion["Receipt Ingestion Service\n(no DB — stateless)"]
        Ledger["Ledger Service\nowns: income, expense,\ncategory, bank hook data"]
        Dashboard["Dashboard Service\n(no DB — reads via Ledger)"]
        DB[("PostgreSQL")]
    end

    Customer -- "sign in" --> GoogleOAuth
    Customer -- "ingest receipt via gallery watch (UC-01)" --> Ingestion
    Customer -- "manage income/expense/\ncategory/bank hook (UC-02)" --> Ledger
    Customer -- "view dashboard (UC-03)" --> Dashboard

    Ingestion -- "store/fetch receipt image" --> ObjectStorage
    Ingestion -- "fetch categories,\nsubmit confirmed candidate" --> Ledger
    Dashboard -- "fetch records for aggregation" --> Ledger
    Ledger --> DB
```

## Notes

- No dedicated Auth/User service in v1 (ADR-03, course convention): each backend service validates the JWT issued after Google OAuth sign-in; the diagram shows Google OAuth as an external system the customer authenticates against directly.
- No API Gateway/reverse proxy shown yet — deferred (ADR Tier B) until service boundaries are stable.
- Receipt Ingestion Service and Dashboard Service intentionally own no data store (per convention: not every service needs a database) — both delegate persistence/source-of-truth to the Ledger Service.
- This is a first iteration; expected to be revised as the design progresses (per course guidance).
