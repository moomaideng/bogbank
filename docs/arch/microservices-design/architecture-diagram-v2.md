# Architecture Diagram (v2)

Supersedes [architecture-diagram.md](./architecture-diagram.md) (v1).

v1 was a single flowchart. The graded deliverable is the **C4 model across four levels plus sequence
diagrams** (Presentation Guidelines, Part 1), so v2 is structured as C4 instead.

| Level | Scope | Status |
| :---- | :---- | :---- |
| L1 Context | BogBank vs. its users and external systems | This document |
| L2 Container | Services, datastores, infrastructure | This document |
| L3 Component | Internals of Ledger SVC (Hexagonal architecture) | This document |
| L4 Code | Hexagonal port/adapter layout (ADR-10, Ledger SVC) | This document |
| Sequence | UC-01 ingest, UC-02 manual, UC-03 dashboard | This document |

---

## Communication style rules

Three rules, applied without exception. They are what makes the style legible in the presentation
and they satisfy three separate course requirements at once.

| Boundary | Style | Why |
| :---- | :---- | :---- |
| Mobile client → system | **REST** (JSON over HTTPS) | Public contract, OpenAPI-documented (ADR-09) |
| Service → service, synchronous | **gRPC** | Typed contracts, internal only, never leaves the cluster |
| Service → service, asynchronous | **Redpanda** (Kafka API) | Choreography, decoupled fan-out, replay |

---

## L1 — System Context

```mermaid
flowchart TB
    Customer(["Customer<br/>(college student)"])

    subgraph BogBank["BogBank"]
        System["Personal expense tracker<br/>Receipt ingestion, ledger, dashboard"]
    end

    Google[["Google OAuth 2.0<br/>Identity provider"]]
    VLM[["Receipt OCR / VLM API<br/>(Gemini Flash or Typhoon OCR)"]]
    Push[["Expo Push Notification Service"]]
    S3[["S3-compatible object storage<br/>(RustFS dev / cloud prod, ADR-07)"]]

    Customer -- "REST: manage hooks, records,<br/>confirm receipts, view dashboard" --> System
    Customer -. "OAuth 2.0 + PKCE sign-in" .-> Google
    System -- "verify ID token" --> Google
    System -- "extract receipt fields" --> VLM
    System -- "push notification" --> Push
    System -- "store / fetch receipt image" --> S3
    Push -. "notification" .-> Customer
```

---

## L2 — Container Diagram

Split into two views. Every container appears in both; the split is by *communication style*, which
keeps each diagram readable and maps to one slide each.

### L2a — Synchronous view (REST at the edge, gRPC internally)

```mermaid
flowchart TB
    Client["Mobile App<br/>Expo / React Native<br/>ADR-05, ADR-06"]

    Traefik["<b>Traefik</b><br/>API Gateway / Ingress<br/>TLS · routing · rate limit<br/>ForwardAuth · discovery"]

    subgraph Services["Application services — Go, hexagonal (ADR-08, ADR-10)"]
        direction LR
        Auth["<b>Auth SVC</b><br/>Google OAuth PKCE<br/>issues BogBank JWT"]
        Receipt["<b>Receipt SVC</b><br/>Bank hooks · upload · dedupe<br/>receipt lifecycle · push"]
        Ledger["<b>Ledger SVC</b><br/>Income · expense · category<br/><i>system of record</i>"]
        Dash["<b>Dashboard SVC</b><br/>Read model + projector"]
        Suggest["<b>Suggestion SVC</b><br/>Merchant to category"]
    end

    subgraph Data["Private datastores"]
        direction LR
        AuthDB[("PostgreSQL<br/>auth")]
        ReceiptDB[("PostgreSQL<br/>receipt<br/>+ outbox")]
        LedgerDB[("PostgreSQL<br/>ledger<br/>+ outbox")]
        CH[("ClickHouse<br/>read model")]
        Redis[("Redis<br/>idempotency<br/>cache · sessions")]
        Blob[("S3 / RustFS<br/>receipt images")]
    end

    Client == "REST" ==> Traefik
    Traefik -- REST --> Auth
    Traefik -- REST --> Receipt
    Traefik -- REST --> Ledger
    Traefik -- REST --> Dash
    Traefik -. "ForwardAuth<br/>validate JWT" .-> Auth

    Ledger -- "gRPC GetReceipt<br/>validate + bank provider" --> Receipt
    Suggest -- "gRPC GetCategories" --> Ledger
    Dash -- "gRPC GetRecords<br/>Sprint 1 + CH-down fallback" --> Ledger

    Auth --> AuthDB
    Receipt --> ReceiptDB
    Receipt --> Blob
    Ledger --> LedgerDB
    Dash --> CH
    Suggest --> Redis
    Auth --> Redis
```

Receipt Reader SVC is absent here: it is stateless and purely event-driven, with no synchronous
inbound API. Every consumer also writes idempotency keys to Redis — omitted as edges to keep the
diagram legible, shown in L2b's note.

### L2b — Asynchronous view (Redpanda choreography)

```mermaid
flowchart LR
    Receipt["<b>Receipt SVC</b>"]
    Reader["<b>Receipt Reader SVC</b><br/>OCR / VLM · stateless"]
    Suggest["<b>Suggestion SVC</b>"]
    Ledger["<b>Ledger SVC</b>"]
    Dash["<b>Dashboard SVC</b><br/>projector"]
    VLM[["OCR / VLM API"]]
    Blob[("S3 / RustFS")]
    CH[("ClickHouse")]
    DLQ["*.dlq"]

    Receipt -- "① receipt.uploaded" --> Reader
    Reader -- "get image" --> Blob
    Reader -- "extract" --> VLM
    Reader -- "② receipt.extracted" --> Receipt
    Reader -- "② receipt.extracted" --> Suggest
    Reader -. "receipt.extraction_failed" .-> Receipt
    Reader -. "poison message" .-> DLQ
    Suggest -- "③ receipt.suggested" --> Receipt
    Ledger -- "④ ledger.expense.created" --> Receipt
    Ledger -- "⑤ ledger.*" --> Dash
    Dash -- "project" --> CH
```

Every labelled edge is a Redpanda topic, never a direct call. Producers publish via the
**transactional outbox**; consumers dedupe on `event_id` in **Redis** because delivery is
at-least-once. Solid = happy path, dotted = failure path.

### Container responsibilities

| Container | Tech | Owns | Talks |
| :---- | :---- | :---- | :---- |
| Traefik | Traefik v3 | Routing, TLS, rate limit, ForwardAuth | — |
| Auth SVC | Go | Users, refresh tokens, JWT signing keys | REST in, gRPC in (token introspection) |
| Receipt SVC | Go | Bank hooks, receipts, extraction results, receipt state machine | REST in, **gRPC in**, Kafka in/out |
| Receipt Reader SVC | Go | Nothing (stateless) | Kafka in/out, HTTP out to VLM |
| Suggestion SVC | Go | Merchant→category rules | Kafka in/out, **gRPC out** to Ledger |
| Ledger SVC | Go | Income, expense, category — **system of record** | REST in, **gRPC in/out**, Kafka out |
| Dashboard SVC | Go | Derived read model only | REST in, Kafka in, **gRPC out** (Sprint 1 + fallback) |

---

## Data ownership and storage

One rule: **a service's datastore is private**. Cross-service reads go through gRPC or events, never
through a shared connection string.

| Service | Store | Shape | Note |
| :---- | :---- | :---- | :---- |
| Auth | PostgreSQL | `users`, `refresh_tokens` | No passwords (NFR1, ADR-03) |
| Receipt | PostgreSQL | `bank_hooks`, `receipts`, `receipt_events`, `outbox` | See below |
| Ledger | PostgreSQL | `categories`, `records`, `outbox` | **System of record, retained forever** |
| Dashboard | ClickHouse | `records_mv` + `AggregatingMergeTree` rollups | **Derived, rebuildable** |
| Suggestion / all | Redis | Consumer idempotency keys, category cache | Required by at-least-once delivery |
| Receipt | S3 / RustFS | Receipt images, SSE-encrypted (NFR3) | ADR-07 |

### Ledger storage: CQRS, not hot/cold tiering

**Postgres keeps every record forever.** ClickHouse is a *projection* of the ledger event stream, not
a second home for old rows. Nothing is ever moved out of Postgres.

Rejected alternative — age-based tiering (PG < 3 months, ClickHouse older):

- FR12's default filters ("current year", "last year") straddle the 3-month boundary, so the common
  path becomes two queries in two SQL dialects merged in application code.
- FR10 allows editing any record. Editing a 4-month-old row means a ClickHouse `ALTER UPDATE` —
  async, non-transactional — instead of a one-line Postgres `UPDATE`.
- The mover job becomes the only copy path: a crash mid-move destroys financial records permanently.
- ClickHouse down means data older than 3 months is *unreachable*, not merely slower.
- The premise doesn't hold: ~3.6k rows/user/year. Postgres with a btree on `(user_id, occurred_at)`
  serves this in single-digit milliseconds. There is no OLTP pressure to relieve.

Under CQRS each of those inverts: one query path, transactional edits, ClickHouse is disposable and
**replayable from Postgres**, and ClickHouse being down degrades the dashboard instead of hiding data.

Explicitly **out of scope**: S3/Iceberg tiering. No justification at this data volume.

### Receipt storage: mutable row + append-only audit trail

Answering "should it be append-only?" — a receipt is a short state machine, so full event sourcing is
overkill. Use both shapes:

```
receipts         (id, user_id, bank_hook_id, image_key, image_sha256,
                  status, extracted jsonb, suggested_category_id,
                  ledger_record_id, created_at, updated_at)
                  UNIQUE (user_id, image_sha256)   -- UC-01 duplicate rule

receipt_events   (id, receipt_id, type, payload jsonb, occurred_at)   -- INSERT only
```

- `status`: `uploaded → extracted → suggested → confirmed`, plus terminal `failed` and `duplicate`.
- `extracted jsonb` absorbs per-bank field variation (KBank and SCB receipts differ) without migrations.
- `receipt_events` is append-only: full audit trail, and it is what you show when asked how the
  async pipeline is debugged.
- `UNIQUE (user_id, image_sha256)` enforces the UC-01 duplicate rule at the database, not in code.

---

## Reliability patterns

These are not optional decoration — they are what makes the event chain correct, and three of them
come straight from the course textbook (Richardson, *Microservices Patterns*).

| Pattern | Where | Problem solved |
| :---- | :---- | :---- |
| **Transactional outbox** | Receipt SVC, Ledger SVC | DB write and event publish are not atomic; a crash between them loses the event. Write both in one transaction, relay the outbox to Redpanda. |
| **Consumer idempotency** | Every consumer | Kafka is at-least-once. Dedupe on `(consumer_group, event_id)` in Redis before handling. |
| **Content-hash dedupe** | Receipt SVC upload | UC-01 duplicate receipt rule, enforced by unique index. |
| **Business-key idempotency** | Ledger `CreateExpense` | `UNIQUE (receipt_id)` — double-tapping the notification cannot create two expenses. |
| **Dead letter topic** | Reader, Suggestion | Extraction failures go to `*.dlq`; receipt is marked `failed` and the user falls back to manual entry (NFR7). |
| **Read-model replay** | Dashboard | ClickHouse rebuilt from Postgres/Kafka; corruption is a routine fix, not an incident. |

---

## L3 — Component Diagram (Ledger Service)

The **Ledger Service** is the financial system of record. Consistent with **ADR-10 (Backend Coding Style: Hexagonal Architecture)** and **ADR-09 (Backend HTTP & API Stack: chi + Huma)**, the service strictly isolates its core business logic from delivery protocols, databases, and remote service clients via ports and adapters.

```mermaid
flowchart TB
    subgraph External["Callers & Remote Collaborators"]
        Traefik["<b>Traefik API Gateway</b><br/>Routes /records, /categories"]
        Dash["<b>Dashboard SVC</b><br/>Direct gRPC client"]
        KafkaIn["<b>Redpanda</b><br/>Topic: receipt.suggested"]
        Postgres[("<b>PostgreSQL</b><br/>ledger schema + outbox")]
        ReceiptSVC["<b>Receipt SVC</b><br/>gRPC server"]
        KafkaOut["<b>Redpanda</b><br/>Topic: ledger.*"]
    end

    subgraph LedgerSVC["Ledger Service (Go) — Hexagonal Architecture (ADR-10)"]
        subgraph InboundAdapters["Inbound / Driving Adapters"]
            HTTPAdapter["<b>HTTP Handler Adapter</b><br/>chi router + Huma v2 (ADR-09)<br/>OpenAPI 3.1 validation"]
            GRPCAdapter["<b>gRPC Server Adapter</b><br/>GetRecords, GetCategories"]
            EventConsumer["<b>Kafka Consumer Adapter</b><br/>receipt.suggested consumer"]
        end

        subgraph InboundPorts["Inbound Ports (Interfaces)"]
            RecordUC["<b>RecordUseCase</b><br/>CreateExpense()<br/>CreateIncome()<br/>UpdateRecord()<br/>DeleteRecord()"]
            CategoryUC["<b>CategoryUseCase</b><br/>ListCategories()<br/>CreateCategory()"]
            QueryUC["<b>RecordQueryUseCase</b><br/>GetRecordsByFilter()"]
        end

        subgraph DomainCore["Domain Core (Hexagon)"]
            RecordService["<b>RecordService</b><br/>Core financial rules<br/>Business-key dedupe"]
            CategoryService["<b>CategoryService</b><br/>Category rules & defaults"]
            DomainEntities["<b>Domain Entities</b><br/>Record, Category, Money"]
        end

        subgraph OutboundPorts["Outbound Ports (Interfaces)"]
            RecordRepoPort["<b>RecordRepository</b><br/>Save(), FindByID(), Query()"]
            CategoryRepoPort["<b>CategoryRepository</b><br/>FindAll(), FindByID()"]
            ReceiptClientPort["<b>ReceiptClient</b><br/>GetReceipt(receipt_id)"]
            OutboxPort["<b>OutboxPublisher</b><br/>SaveEvent(tx, event)"]
        end

        subgraph OutboundAdapters["Outbound / Driven Adapters"]
            PGRepoAdapter["<b>Postgres Repository Adapter</b><br/>Bun ORM (PostgreSQL)"]
            ReceiptGRPCAdapter["<b>Receipt gRPC Client Adapter</b><br/>Typed protobuf client"]
            OutboxAdapter["<b>Transactional Outbox Adapter</b><br/>Writes to ledger.outbox"]
        end
    end

    Traefik -- "REST" --> HTTPAdapter
    Dash -- "gRPC: GetRecords" --> GRPCAdapter
    KafkaIn -- "Consume event" --> EventConsumer

    HTTPAdapter --> RecordUC
    HTTPAdapter --> CategoryUC
    GRPCAdapter --> QueryUC
    GRPCAdapter --> CategoryUC
    EventConsumer --> RecordUC

    RecordUC -. "implements" .-> RecordService
    CategoryUC -. "implements" .-> CategoryService
    QueryUC -. "implements" .-> RecordService

    RecordService --> DomainEntities
    CategoryService --> DomainEntities

    RecordService --> RecordRepoPort
    RecordService --> ReceiptClientPort
    RecordService --> OutboxPort
    CategoryService --> CategoryRepoPort

    RecordRepoPort -. "implements" .-> PGRepoAdapter
    CategoryRepoPort -. "implements" .-> PGRepoAdapter
    ReceiptClientPort -. "implements" .-> ReceiptGRPCAdapter
    OutboxPort -. "implements" .-> OutboxAdapter

    PGRepoAdapter --> Postgres
    OutboxAdapter --> Postgres
    ReceiptGRPCAdapter -- "gRPC: GetReceipt" --> ReceiptSVC
    OutboxAdapter -. "Outbox relayed to" .-> KafkaOut
```

### Component responsibilities (Ledger Service)

| Component | Layer | Purpose | Technology / Pattern |
| :---- | :---- | :---- | :---- |
| **HTTP Handler Adapter** | Inbound Adapter | Exposes REST endpoints (`/records`, `/categories`), deserializes JSON, validates inputs, generates OpenAPI 3.1 | chi + Huma v2 (`humachi` adapter, ADR-09) |
| **gRPC Server Adapter** | Inbound Adapter | Handles synchronous internal RPC calls (`GetRecords`, `GetCategories`) from Dashboard and Suggestion services | gRPC, Go protobuf |
| **Kafka Consumer Adapter** | Inbound Adapter | Consumes asynchronous `receipt.suggested` events, deduplicates in Redis, triggers auto-creation flow | Segmentio/kafka-go or confluent-kafka-go |
| **RecordService** | Domain Core | Implements financial transaction invariants: positive amounts, valid categories, `UNIQUE (receipt_id)` | Pure Go (no external dependencies, ADR-10) |
| **CategoryService** | Domain Core | Manages system default categories (Food, Transport, Bills) and custom user categories | Pure Go |
| **Postgres Repository Adapter**| Outbound Adapter | Executes SQL queries and transactions using Bun ORM connection pool | `uptrace/bun`, PostgreSQL |
| **Receipt gRPC Client Adapter**| Outbound Adapter | Calls Receipt SVC's `GetReceipt` endpoint to fetch receipt details for verified ledger records | Go gRPC Client |
| **Transactional Outbox Adapter**| Outbound Adapter | Appends domain events to `ledger.outbox` inside the same SQL transaction as the record | Transactional Outbox pattern |

---

## L3b — Component Diagram (Receipt Service)

The **Receipt Service** is the asynchronous ingestion engine. It handles media gallery uploads, deduplicates receipts by SHA-256 content hash (UC-01), maintains receipt lifecycle states (`uploaded` $\rightarrow$ `extracted` $\rightarrow$ `suggested` $\rightarrow$ `confirmed`), coordinates with object storage, and triggers push notifications.

```mermaid
flowchart TB
    subgraph External["Callers & Remote Collaborators"]
        Traefik["<b>Traefik API Gateway</b><br/>Routes /receipts/*, /bank-hooks/*"]
        LedgerSVC["<b>Ledger SVC</b><br/>gRPC caller: GetReceipt"]
        KafkaIn["<b>Redpanda</b><br/>Topics: receipt.extracted, receipt.suggested, ledger.expense.created"]
        Postgres[("<b>PostgreSQL</b><br/>receipt schema + outbox")]
        ObjectStorage[("<b>S3 / RustFS</b><br/>Receipt images (ADR-07)")]
        ExpoPush["<b>Expo Push Service</b><br/>Mobile push delivery"]
        KafkaOut["<b>Redpanda</b><br/>Topic: receipt.uploaded"]
    end

    subgraph ReceiptSVC["Receipt Service (Go) — Hexagonal Architecture (ADR-10)"]
        subgraph InboundAdapters["Inbound / Driving Adapters"]
            HTTPAdapter["<b>HTTP Handler Adapter</b><br/>chi router + Huma v2 (ADR-09)<br/>Multipart upload & hook management"]
            GRPCAdapter["<b>gRPC Server Adapter</b><br/>GetReceipt RPC"]
            EventConsumer["<b>Kafka Consumer Adapter</b><br/>Consumes async pipeline events"]
        end

        subgraph InboundPorts["Inbound Ports (Interfaces)"]
            UploadUC["<b>ReceiptUploadUseCase</b><br/>UploadReceipt()<br/>DeduplicateByHash()"]
            HookUC["<b>BankHookUseCase</b><br/>ManageBankHooks()"]
            LifecycleUC["<b>ReceiptLifecycleUseCase</b><br/>HandleExtraction()<br/>HandleSuggestion()<br/>HandleConfirmation()"]
        end

        subgraph DomainCore["Domain Core (Hexagon)"]
            ReceiptService["<b>ReceiptService</b><br/>Receipt state machine<br/>SHA-256 dedupe check"]
            HookService["<b>BankHookService</b><br/>Provider validation (KBank, SCB)"]
            DomainEntities["<b>Domain Entities</b><br/>Receipt, BankHook, ReceiptEvent"]
        end

        subgraph OutboundPorts["Outbound Ports (Interfaces)"]
            ReceiptRepoPort["<b>ReceiptRepository</b><br/>Save(), FindByID(), CheckDedupe()"]
            HookRepoPort["<b>BankHookRepository</b><br/>SaveHook(), ListHooks()"]
            BlobPort["<b>ObjectStorageClient</b><br/>PutObject(), GetObject()"]
            PushPort["<b>NotificationClient</b><br/>SendPushNotification()"]
            OutboxPort["<b>OutboxPublisher</b><br/>SaveEvent(tx, event)"]
        end

        subgraph OutboundAdapters["Outbound / Driven Adapters"]
            BunRepoAdapter["<b>Postgres Repository Adapter</b><br/>Bun ORM (PostgreSQL)"]
            S3Adapter["<b>S3 Storage Adapter</b><br/>AWS S3 SDK (RustFS / Cloud S3)"]
            PushAdapter["<b>Expo Push Adapter</b><br/>Expo Push HTTP Client"]
            OutboxAdapter["<b>Transactional Outbox Adapter</b><br/>Writes to receipt.outbox"]
        end
    end

    Traefik -- "REST" --> HTTPAdapter
    LedgerSVC -- "gRPC: GetReceipt" --> GRPCAdapter
    KafkaIn -- "Consume events" --> EventConsumer

    HTTPAdapter --> UploadUC
    HTTPAdapter --> HookUC
    GRPCAdapter --> LifecycleUC
    EventConsumer --> LifecycleUC

    UploadUC -. "implements" .-> ReceiptService
    HookUC -. "implements" .-> HookService
    LifecycleUC -. "implements" .-> ReceiptService

    ReceiptService --> DomainEntities
    HookService --> DomainEntities

    ReceiptService --> ReceiptRepoPort
    ReceiptService --> BlobPort
    ReceiptService --> PushPort
    ReceiptService --> OutboxPort
    HookService --> HookRepoPort

    ReceiptRepoPort -. "implements" .-> BunRepoAdapter
    HookRepoPort -. "implements" .-> BunRepoAdapter
    BlobPort -. "implements" .-> S3Adapter
    PushPort -. "implements" .-> PushAdapter
    OutboxPort -. "implements" .-> OutboxAdapter

    BunRepoAdapter --> Postgres
    OutboxAdapter --> Postgres
    S3Adapter -- "S3 API" --> ObjectStorage
    PushAdapter -- "HTTPS" --> ExpoPush
    OutboxAdapter -. "Outbox relayed to" .-> KafkaOut
```

### Component responsibilities (Receipt Service)

| Component | Layer | Purpose | Technology / Pattern |
| :---- | :---- | :---- | :---- |
| **HTTP Handler Adapter** | Inbound Adapter | Exposes `POST /receipts` (multipart image upload) and bank hook management routes | chi + Huma v2 (`humachi` adapter, ADR-09) |
| **gRPC Server Adapter** | Inbound Adapter | Serves `GetReceipt` for Ledger Service to verify receipt metadata upon user confirmation | gRPC, Go protobuf |
| **Kafka Consumer Adapter** | Inbound Adapter | Listens for `receipt.extracted`, `receipt.suggested`, and `ledger.expense.created` | Segmentio/kafka-go or confluent-kafka-go |
| **ReceiptService** | Domain Core | Manages receipt state machine (`uploaded` $\rightarrow$ `extracted` $\rightarrow$ `suggested` $\rightarrow$ `confirmed`), enforces SHA-256 deduplication | Pure Go (ADR-10) |
| **BankHookService** | Domain Core | Validates configured bank providers (KBank, SCB, etc.) and user gallery permissions | Pure Go (ADR-10) |
| **Postgres Repository Adapter**| Outbound Adapter | Persists receipts, bank hooks, append-only `receipt_events`, and outbox table | `uptrace/bun`, PostgreSQL |
| **S3 Storage Adapter** | Outbound Adapter | Uploads raw encrypted receipt images to S3-compatible storage | AWS S3 Go SDK v2 (RustFS dev / S3 prod, ADR-07) |
| **Expo Push Adapter** | Outbound Adapter | Dispatches push notifications to customer mobile device when category is suggested | Expo Push HTTP API |
| **Transactional Outbox Adapter**| Outbound Adapter | Appends `receipt.uploaded` to `receipt.outbox` atomically with the receipt insertion | Transactional Outbox pattern |

---

## L4 — Code Diagram (Ledger Service: Record Domain Hexagonal Layout)

Zooming into the code structure of the `Record` bounded context within `Ledger SVC`. Per **ADR-10**, the domain core depends strictly on interfaces (ports), while HTTP handlers, DB drivers, and gRPC clients are pluggable adapters.

```mermaid
classDiagram
    class RecordHandler {
        -recordUC RecordUseCase
        +PostExpense(ctx context.Context, input PostExpenseInput) (PostExpenseOutput, error)
        +PostIncome(ctx context.Context, input PostIncomeInput) (PostIncomeOutput, error)
        +RegisterRoutes(api huma.API)
    }

    class RecordUseCase {
        <<interface>>
        +CreateExpenseFromReceipt(ctx context.Context, cmd CreateExpenseCmd) (*Record, error)
        +CreateManualExpense(ctx context.Context, cmd CreateManualExpenseCmd) (*Record, error)
        +GetRecords(ctx context.Context, filter RecordFilter) ([]Record, error)
    }

    class RecordService {
        -repo RecordRepository
        -receiptClient ReceiptClient
        -outbox OutboxPublisher
        +CreateExpenseFromReceipt(ctx context.Context, cmd CreateExpenseCmd) (*Record, error)
        +CreateManualExpense(ctx context.Context, cmd CreateManualExpenseCmd) (*Record, error)
        +GetRecords(ctx context.Context, filter RecordFilter) ([]Record, error)
    }

    class RecordRepository {
        <<interface>>
        +Save(ctx context.Context, tx DBTx, record *Record) error
        +FindByID(ctx context.Context, id uuid.UUID) (*Record, error)
        +FindByFilter(ctx context.Context, filter RecordFilter) ([]Record, error)
    }

    class PostgresRecordRepo {
        -db *bun.DB
        +Save(ctx context.Context, tx DBTx, record *Record) error
        +FindByID(ctx context.Context, id uuid.UUID) (*Record, error)
        +FindByFilter(ctx context.Context, filter RecordFilter) ([]Record, error)
    }

    class ReceiptClient {
        <<interface>>
        +GetReceipt(ctx context.Context, id uuid.UUID) (*ReceiptDetails, error)
    }

    class GRPCReceiptClient {
        -client pb.ReceiptServiceClient
        +GetReceipt(ctx context.Context, id uuid.UUID) (*ReceiptDetails, error)
    }

    class OutboxPublisher {
        <<interface>>
        +Publish(ctx context.Context, tx DBTx, event OutboxEvent) error
    }

    class PostgresOutboxPublisher {
        -db *bun.DB
        +Publish(ctx context.Context, tx DBTx, event OutboxEvent) error
    }

    class Record {
        +ID uuid.UUID
        +UserID uuid.UUID
        +Amount decimal.Decimal
        +Type RecordType
        +CategoryID uuid.UUID
        +ReceiptID *uuid.UUID
        +BankProvider *string
        +OccurredAt time.Time
        +Validate() error
    }

    RecordHandler ..> RecordUseCase : drives (HTTP)
    RecordService ..|> RecordUseCase : implements port
    RecordService --> RecordRepository : driven port
    RecordService --> ReceiptClient : driven port
    RecordService --> OutboxPublisher : driven port
    RecordService ..> Record : constructs & validates
    PostgresRecordRepo ..|> RecordRepository : adapter implements
    GRPCReceiptClient ..|> ReceiptClient : adapter implements
    PostgresOutboxPublisher ..|> OutboxPublisher : adapter implements
```

---

## Sequence — UC-01: Ingest expense from bank receipt

The main demo path. Note that every arrow is labelled with its style (REST / gRPC / event), which is
the requirement-coverage evidence.

```mermaid
sequenceDiagram
    actor U as Customer
    participant App as Mobile App
    participant TR as Traefik
    participant RC as Receipt SVC
    participant S3 as Object Storage
    participant K as Redpanda
    participant RD as Reader SVC
    participant SG as Suggestion SVC
    participant LG as Ledger SVC
    participant DS as Dashboard SVC

    U->>App: bank app saves receipt to gallery
    App->>App: media-library watcher detects new image
    App->>TR: POST /receipts (multipart) [REST]
    TR->>RC: forward (JWT validated via ForwardAuth)
    RC->>RC: sha256(image), reject if duplicate
    RC->>S3: PutObject (encrypted)
    RC->>RC: TX: INSERT receipt(status=uploaded) + outbox
    RC-->>App: 202 Accepted {receipt_id}
    RC->>K: receipt.uploaded

    K->>RD: receipt.uploaded
    RD->>S3: GetObject
    RD->>RD: OCR / VLM extract (amount, merchant, date)
    alt extraction succeeded
        RD->>K: receipt.extracted
    else failed
        RD->>K: receipt.extraction_failed
    end

    par Receipt persists the facts
        K->>RC: receipt.extracted
        RC->>RC: UPDATE status=extracted, extracted=jsonb
        RC->>App: push "New expense detected"
    and Suggestion classifies
        K->>SG: receipt.extracted
        SG->>LG: GetCategories(user_id) [gRPC]
        LG-->>SG: categories
        SG->>SG: match merchant to category
        SG->>K: receipt.suggested
    end

    K->>RC: receipt.suggested
    RC->>RC: UPDATE suggested_category_id, status=suggested
    opt user has not confirmed yet
        RC->>App: push "Suggested: Food - confirm?"
    end

    U->>App: tap category (confirm or change)
    App->>TR: POST /records {receipt_id, category_id} [REST]
    TR->>LG: forward
    LG->>RC: GetReceipt(receipt_id) [gRPC]
    RC-->>LG: {amount, merchant, occurred_at, bank_provider}
    LG->>LG: TX: INSERT record (UNIQUE receipt_id) + outbox
    LG-->>App: 201 Created
    LG->>K: ledger.expense.created

    par
        K->>RC: ledger.expense.created
        RC->>RC: UPDATE status=confirmed, ledger_record_id
    and
        K->>DS: ledger.expense.created
        DS->>DS: project into ClickHouse read model
    end
```

**Why the confirm call goes to Ledger, not Receipt:** Ledger is the authority for expense records, so
the user gets a synchronous `201` confirming the money was recorded. The receipt's own state converges
asynchronously — it is metadata about the ingest pipeline, not the thing the user is waiting on.
`UNIQUE (receipt_id)` makes a double tap a no-op.

## Sequence — UC-03: Dashboard

```mermaid
sequenceDiagram
    actor U as Customer
    participant App as Mobile App
    participant TR as Traefik
    participant DS as Dashboard SVC
    participant CH as ClickHouse
    participant LG as Ledger SVC

    U->>App: open dashboard (type, range, bank filter)
    App->>TR: GET /dashboard/donut?... [REST]
    TR->>DS: forward
    alt Sprint 2+ (read model live)
        DS->>CH: SELECT from AggregatingMergeTree rollup
        CH-->>DS: aggregated buckets
    else Sprint 1 / ClickHouse unavailable
        DS->>LG: GetRecords(filter) [gRPC]
        LG-->>DS: rows
        DS->>DS: aggregate in memory
    end
    DS-->>App: chart data
```

The `else` branch is not a hack — it is the graceful-degradation path that makes ClickHouse
disposable, and it is how Sprint 1 ships before the read model exists.

---

## Event catalog

| Topic | Producer | Consumers | Key |
| :---- | :---- | :---- | :---- |
| `receipt.uploaded` | Receipt | Reader | `receipt_id` |
| `receipt.extracted` | Reader | Receipt, Suggestion | `receipt_id` |
| `receipt.extraction_failed` | Reader | Receipt | `receipt_id` |
| `receipt.suggested` | Suggestion | Receipt | `receipt_id` |
| `ledger.expense.created` | Ledger | Receipt, Dashboard | `record_id` |
| `ledger.expense.updated` / `.deleted` | Ledger | Dashboard | `record_id` |
| `ledger.income.created` / `.updated` / `.deleted` | Ledger | Dashboard | `record_id` |

Envelope on every event: `event_id` (UUID, the idempotency key), `event_type`, `occurred_at`,
`user_id`, `trace_id` (propagated for observability), `payload`.

Partition key is the entity id so that per-entity ordering holds; `user_id` is available for
future repartitioning.

---

## Service discovery

A required deliverable, answered per environment.

| Environment | Mechanism |
| :---- | :---- | 
| Development | Docker Compose embedded DNS; service name = hostname. Traefik discovers routes from the **Docker provider** via container labels — no static route config. |
| Production (k3s) | Kubernetes `Service` DNS (`<svc>.<ns>.svc.cluster.local`); Traefik discovers routes from the **Kubernetes CRD provider** (`IngressRoute`). |
| Async | Broker-based: producers and consumers only know a topic name, never a peer address. Consumer-group rebalance handles instance membership. |

Client-side discovery and a dedicated registry (Eureka/Consul) were not adopted: the platform already
provides DNS-based server-side discovery in both environments, so a registry would be a component
with no responsibility.

---

## Environments

### Development — Docker Compose profiles

Compose `profiles:` gives exactly the "enable service A and its dependencies" behaviour we want.

```
infra      always on: postgres, redpanda, redis, rustfs, traefik
auth       auth-svc            (+ infra)
receipt    receipt-svc         (+ infra)
reader     reader-svc          (+ infra)
suggest    suggest-svc         (+ ledger, infra)
ledger     ledger-svc          (+ infra)
dash       dashboard-svc       (+ clickhouse, infra)
obs        grafana, prometheus, loki
full       everything
```

Driven by Makefile targets (`make dev SVC=receipt`, `make dev-full`). Go services run under `air`
for hot reload. One Postgres instance, one schema per service — separate schemas preserve the
"private datastore" rule while keeping the laptop footprint to a single container.

### Production — k3s, as a Sprint 3 stretch

Docker Compose on a single VM is the **guaranteed demo path** and must be working end to end first.
k3s is additive:

- Single-node k3s (Traefik ships built in, so the edge config carries over).
- **Kustomize only** — no Helm charts authored by us.
- **One dedicated owner.** CI stays Compose-based so the critical path never depends on k3s.
- Explicitly out of scope: service mesh, multi-node, autoscaling.

Kubernetes is not required by the syllabus — "describe how to manage Service Discovery" is, and
Compose DNS answers it. k3s is taken for the week-12 deployment story, and it is droppable.

### Observability — Grafana stack

Prometheus (metrics) + Loki (logs) + Grafana (dashboards), with OpenTelemetry SDK in every Go service.
Tempo (traces) added only if Sprint 3 allows — distributed tracing across the async chain is the most
compelling observability demo, so it is the first stretch item.

---

## Quality attribute (graded: "1 measurable quality attribute")

**Primary: Performance / Scalability of the ingest pipeline.**
Measured as end-to-end latency from `POST /receipts` to `receipt.suggested`, at p50/p95/p99, under
k6 load. The measurable claim: horizontally scaling Reader SVC consumer instances reduces p95 in
proportion to partition count, until partitions are saturated.

**Secondary: Resilience.** Kill ClickHouse mid-demo; dashboard degrades to the Ledger gRPC path
instead of erroring. Kill Reader SVC mid-demo; restart it and Redpanda replays the backlog with no
lost receipts.

Load tests (two types required):
1. **Constant-arrival-rate load test** on `GET /dashboard/*` — steady-state latency.
2. **Spike/stress test** on `POST /receipts` — queue depth, consumer lag, recovery time.

---

## Requirement traceability

| Course requirement | Satisfied by |
| :---- | :---- |
| ≥3 business use cases + authentication | UC-01, UC-02, UC-03 + Auth SVC |
| Front-end UI live demo | Expo mobile app (ADR-05, ADR-06) |
| ≥1 REST service | All client-facing services (Auth, Receipt, Ledger, Dashboard) |
| ≥1 gRPC service | Ledger→Receipt, Suggestion→Ledger, Dashboard→Ledger |
| ≥1 message broker service | Reader, Suggestion, Receipt, Dashboard (Redpanda) |
| API Gateway | Traefik |
| Service discovery described | Compose DNS / k3s Service DNS / consumer groups (above) |
| ≥2 database types | PostgreSQL (RDBMS) + Redis (NoSQL KV) + ClickHouse (OLAP columnar) |
| Load test, 2 types | Constant-arrival-rate + spike/stress (k6) |
| Risk matrix ≥3 issues | Below |
| 1 measurable quality attribute | Ingest-pipeline performance; resilience secondary |
| C4 four levels + sequence | This document (L1, L2, L3, L4 complete + sequences) |

---

## Risk matrix (draft — expand to the graded format later)

| # | Risk | Likelihood | Impact | Mitigation |
| :---- | :---- | :---- | :---- | :---- |
| R1 | Thai receipt OCR accuracy is poor | High | High | Use a VLM API (Gemini Flash / Typhoon OCR) not Tesseract; manual entry always available (NFR7); low-confidence extraction leaves the receipt `failed` rather than guessing |
| R2 | Bank apps block screenshots (`FLAG_SECURE`) or use private albums | High | Medium | Documented per-provider support matrix; UC-02 manual fallback; demo with a provider known to work |
| R3 | k3s deployment consumes Sprint 3 and no demo ships | Medium | High | Compose-on-VM is the guaranteed path; k3s is additive with one owner; hard cutoff date |
| R4 | Team unfamiliar with Kafka semantics → lost or duplicated records | Medium | High | Outbox + Redis idempotency from day one; Redpanda to reduce operational load; the week-10 tutorial lands just before Sprint 2 |
| R5 | 7 services for 5 people in 9 weeks | Medium | High | Vertical slice ownership (below); Suggestion stays rule-based; no service mesh, no Iceberg, no tiering |
| R6 | iOS background execution limits break gallery watch | Medium | Medium | Foreground-triggered scan on app open; documented as a known platform constraint |

---

## Proposed work split (5 people, full-stack vertical slices)

Two rules drive this split:

1. **Everyone owns backend microservices.** Nobody is "the frontend person". Project progress is
   graded on per-person GIT contribution three times (5% + 5% + 5%), and every member has to be able
   to answer architecture questions in the 5-minute Q&A about code they actually wrote.
2. **Frontend is sliced, not owned.** Each person builds the screens that consume their own service.
   The Expo app is one codebase but the slices touch disjoint route folders, so they merge cleanly.

| # | Services owned (backend) | Frontend slice | IPC styles they personally implement |
| :---- | :---- | :---- | :---- |
| P1 | **Auth SVC** + Traefik edge config + shared `pkg/outbox`, `pkg/eventbus` | Sign-in / sign-out, token storage, refresh interceptor | REST server, gRPC server (introspection), **Kafka producer/consumer library everyone reuses** |
| P2 | **Ledger SVC** | Income / expense / category CRUD screens | REST server, gRPC server + gRPC client (`GetReceipt`), Kafka producer |
| P3 | **Receipt SVC** | Bank hook management, media-library permission, upload, notification confirm sheet | REST server, gRPC server, Kafka producer **and** consumer |
| P4 | **Dashboard SVC** + ClickHouse projection | Donut chart, bar chart, filter controls | REST server, gRPC client, Kafka consumer |
| P5 | **Receipt Reader SVC** + **Suggestion SVC** + platform (Compose, CI, Grafana, k6) | App shell, navigation, shared UI kit, generated API client | Kafka consumer/producer ×2, gRPC client (`GetCategories`), REST admin API |

Every person ends up having written at least one REST handler, one gRPC endpoint or client, and one
Kafka producer or consumer. That is deliberate: it is the same list the syllabus grades, and it means
any of the five can field any Q&A question.

**Load balance check.** P5's two services are the smallest in the system (Reader is stateless glue
around a VLM call; Suggestion is rule matching) and both land in Sprint 2 — so P5 carries platform
and app shell in Sprint 1, when the others are blocked on scaffolding anyway. P1's Auth SVC is small
after Sprint 1, which is why P1 also owns the shared outbox/eventbus library that Sprint 2 depends on.

**Frontend merge convention.** Each slice owns `app/(tabs)/<slice>/` and registers its routes; the
shell, navigation, theme, and `components/ui/` belong to P5 and change only by PR. The generated
API client is produced from each service's OpenAPI spec (ADR-09), so nobody hand-writes another
person's request types.

**Cross-slice contracts** — protobuf service definitions and the event envelope schema — are agreed
**before** Sprint 1 coding and live in a shared monorepo package. They are the only coupling between
slices, and changing one requires the downstream owner to review the PR.

---

## ADR consequences

| ADR | Status after v2 | Action |
| :---- | :---- | :---- |
| ADR-03 Authentication Pattern | Still valid | Own Auth SVC is exactly what it specifies. Note in the new IdP ADR **why Authelia was rejected**: it implements the OIDC *Provider* role only, cannot federate to Google, and its auth backends store password hashes — contradicting NFR1. |
| ADR-04 Database Technology | **Must be superseded** | It states "PostgreSQL for every service's data store" and explicitly rejects ClickHouse. v2 adds ClickHouse and Redis. Write ADR-11 superseding it. |
| ADR-08 Backend Language | Still valid | Reader SVC stays Go by calling a hosted VLM over HTTP rather than embedding a Python OCR stack. |
| ADR-10 Backend Coding Style | Reinforced | Dashboard's PG→ClickHouse switch is a pure adapter swap behind a port — the payoff this ADR promised. |

New ADRs to write (these close the "Deferred (Tier B)" list in `docs/arch/adr/README.md`):

| # | Title |
| :---- | :---- |
| 11 | Database Technology (supersedes ADR-04) — PostgreSQL + ClickHouse + Redis |
| 12 | Internal IPC Style — gRPC internal, REST at the edge |
| 13 | Message Broker — Redpanda (Kafka API) |
| 14 | API Gateway / Reverse Proxy — Traefik |
| 15 | Identity Provider — own Go Auth SVC (rejects Authelia, Keycloak, Zitadel) |
| 16 | Dev & Prod Environments — Compose profiles, k3s stretch |
| 17 | Observability Stack — Grafana + Prometheus + Loki (+ Tempo) |
| 18 | Ledger Read Model — CQRS, rejects hot/cold tiering |

---

## Open items for the meeting

1. **Redis as the NoSQL store.** The syllabus requires it twice, in those words: §17 "Conduct at
   least 2 types of databases (both RDBMS and No-SQL)" and Presentation Guidelines "one relational
   (RDBMS) and one NoSQL". ClickHouse self-describes as a SQL DBMS, so it is a taxonomy argument we
   would have to win live. Redis is needed anyway for consumer idempotency, so the marginal cost is
   one container. Decide: keep Redis, or accept the risk and drop it.
2. **Notification delivery** currently lives inside Receipt SVC (Expo Push SDK). A separate
   Notification SVC is cleaner but is an eighth service. Recommend: keep it in Receipt SVC, extract
   later only if in-app/email channels appear.
3. **Suggestion SVC algorithm.** Rule-based merchant-string matching is the scoped plan. An LLM call
   would connect to the week 13–14 AI-architecture lectures — decide whether that is worth the
   latency and cost, or whether Reader SVC's VLM already covers the AI narrative.
4. **API composition.** With Traefik-only, Dashboard SVC is the service that performs API
   composition (aggregating Ledger data into chart responses) — that is the answer to give in week 11.
   If the grader wants a distinct composition layer, a thin Go BFF can be added behind Traefik
   without changing any service contract.
