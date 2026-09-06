# Business Use Cases

Per course convention, Sign In and Sign Out are not counted as Business Use Cases. These are the 3 primary Business Use Cases for the project.

## UC-01: Ingest Expense from Bank Receipt (Gallery Watch)

**Primary Actor**: Customer (college student)

**Goal**: Record an expense with minimal manual typing by automatically picking up new bank receipt images from the device gallery instead of entering each transaction by hand.

**Preconditions**: Customer is signed in; customer has at least one bank expense hook configured and has granted media-library access for that hook (see `Manage Bank Expense Hook`).

**Main Flow**:
1. Customer completes a bank payment; the bank app saves a receipt image to the device gallery (where the OS and bank app allow).
2. System detects a newly added receipt image associated with the configured bank hook (media-library listener or polling via the native mobile app).
3. System uploads the image to object storage.
4. System extracts metadata (amount, merchant/counterparty, date) from the stored image.
5. System creates a pending expense record and notifies the customer with a suggested category.
6. Customer confirms the suggested category or picks a different one (`<<include>> Manage Category`).
7. System finalizes the expense record under the confirmed category.

**Postcondition**: A categorized expense record exists and is reflected in the dashboard.

**Outcome the customer wants**: Their spending gets tracked accurately without manually re-typing every bank transaction.

## UC-02: Track Income and Expenses Manually

**Primary Actor**: Customer (college student)

**Goal**: Record or correct income/expense transactions that aren't covered by gallery-watch receipt intake — cash transactions, other banks, blocked screenshots, a missed detection, or fixing a mistake.

**Preconditions**: Customer is signed in.

**Main Flow**:
1. Customer opens "Add income" or "Add expense".
2. Customer enters amount, date, category (`<<include>> Manage Category`), and an optional note.
3. Customer saves the record.
4. Customer may later edit or remove an existing income/expense record from either flow (including ones created via UC-01).

**Postcondition**: The customer's income/expense record set is complete and accurate.

**Outcome the customer wants**: Nothing falls through the cracks even when the automatic gallery-watch flow doesn't apply.

## UC-03: Review Income/Expense Dashboard

**Primary Actor**: Customer (college student)

**Goal**: Understand spending and income patterns over time and by bank provider, to make better financial decisions — this is the core problem the product solves (expenses are tedious to track and categorize across many channels).

**Preconditions**: Customer is signed in; customer has at least one income or expense record.

**Main Flow**:
1. Customer opens the dashboard.
2. Customer views a donut chart of income or expense broken down by category, filterable by time range (current/last month, current/last year, or a custom range) and bank provider.
3. Customer switches to a bar chart comparing income vs. expense across a fixed set of time intervals (week/month/3-months/year), filterable by bank provider.

**Postcondition**: None (read-only use case).

**Outcome the customer wants**: An at-a-glance, filterable view of their financial activity without manual spreadsheet work.

## Supporting Use Cases (not counted among the 3)

- **Manage Bank Expense Hook** — Customer: add/edit/remove a bank provider hook and grant media-library permission for receipt detection. Prerequisite for UC-01.
- **Manage Category** — Customer: add/edit/remove categories used to classify income/expense records. Included by UC-01 and UC-02.
