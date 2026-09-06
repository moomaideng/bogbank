# Service–Operations–Collaborators (v1)

First-iteration decomposition, derived from the 3 Business Use Cases in `docs/use-cases.md`. Follows the course's recommended path: Actor → Business Use Cases → Operations → Service Responsibilities → Collaborators → Data Ownership → Diagram (see `supplimentary-materials/convention-for-work.md`).

Per the convention notes for this round: no dedicated Auth/User service (auth handled by ADR-03's Google OAuth + JWT, validated per-service), single REST style everywhere, single database technology (ADR-04), no message broker.

## Services

### Receipt Ingestion Service
**Business capability**: Turn a newly detected bank receipt image into a structured expense candidate with a suggested category. Owns UC-01.

| Operation | Description |
|---|---|
| `UploadReceipt(image)` | Accepts a receipt image uploaded by the mobile client after gallery-watch detection, stores it in object storage. |
| `ExtractMetadata(imageRef)` | Extracts amount, merchant/counterparty, and date from the stored image. |
| `SuggestCategory(merchantHint, customerId)` | Proposes a category for the extracted transaction, using the customer's existing categories. |
| `SubmitExpenseCandidate(customerId, extractedData, confirmedCategoryId)` | Hands off the confirmed candidate to the Ledger Service to create the final expense record. |

**Collaborators**: Object Storage (store/fetch the receipt image, ADR-07) · Ledger Service (fetch the customer's categories for suggestion; create the final expense record once confirmed).

**Data ownership**: None. Stateless processing pipeline — persistence of the final record is delegated to the Ledger Service.

### Ledger Service
**Business capability**: Own and maintain every income record, expense record, category, and bank hook — the single source of truth for a customer's financial data. Owns UC-02, plus the `Manage Bank Expense Hook` and `Manage Category` supporting use cases, plus finalizing records handed off from Receipt Ingestion.

| Operation | Description |
|---|---|
| `AddIncome` / `EditIncome` / `RemoveIncome` | Manual income record CRUD. |
| `AddExpense` / `EditExpense` / `RemoveExpense` | Manual expense record CRUD (also used as the fallback path in UC-02). |
| `CreateExpenseFromCandidate(customerId, extractedData, categoryId)` | Finalizes an expense record from Receipt Ingestion's confirmed candidate. |
| `AddCategory` / `EditCategory` / `RemoveCategory` | Category CRUD. |
| `AddBankHook` / `EditBankHook` / `RemoveBankHook` | Bank expense hook CRUD. |
| `GetRecords(customerId, filters)` | Returns income/expense records filtered by type, time range, and bank provider — consumed by the Dashboard Service. |

**Collaborators**: PostgreSQL (owns, ADR-04) · called by Receipt Ingestion Service and Dashboard Service.

**Data ownership**: Owns the `income_records`, `expense_records`, `categories`, and `bank_hooks` tables in PostgreSQL. No other service accesses this data directly — only through these operations.

### Dashboard Service
**Business capability**: Aggregate and shape a customer's financial data for visualization. Owns UC-03.

| Operation | Description |
|---|---|
| `GetDonutChartData(customerId, type, timeRange, bankProvider)` | Category breakdown for the donut chart. |
| `GetBarChartData(customerId, interval, bankProvider)` | Income-vs-expense series across the fixed time-interval presets (ADR-driven, see FR13). |

**Collaborators**: Ledger Service (fetches the underlying records via `GetRecords`; no direct database access).

**Data ownership**: None. Read-only aggregation over data owned by the Ledger Service.

## Actor → Use Case → Service Trace

| Actor | Use Case | Operation | Responsible Service | Collaborators | Data Owner |
|---|---|---|---|---|---|
| Customer | UC-01 Ingest Expense from Bank Receipt | `UploadReceipt` → `ExtractMetadata` → `SuggestCategory` → `SubmitExpenseCandidate` | Receipt Ingestion Service | Object Storage, Ledger Service | Ledger Service (final record) |
| Customer | UC-02 Track Income/Expenses Manually | `AddIncome`/`AddExpense`/etc. | Ledger Service | PostgreSQL | Ledger Service |
| Customer | UC-03 Review Dashboard | `GetDonutChartData` / `GetBarChartData` | Dashboard Service | Ledger Service | — (read-only) |
| Customer | Manage Bank Expense Hook (supporting) | `AddBankHook`/etc. | Ledger Service | PostgreSQL | Ledger Service |
| Customer | Manage Category (supporting) | `AddCategory`/etc. | Ledger Service | PostgreSQL | Ledger Service |
| Customer | Sign In / Sign Out | — | *(not a dedicated service in v1 — see ADR-03 and convention rule on Auth)* | Google OAuth | — |

## Sanity check (per course convention)

Removing the service names, each box still explains itself: one service turns a raw receipt image into a categorized transaction; one service is the system of record for all financial data; one service turns that data into charts. All 3 Business Use Cases are covered end-to-end.
