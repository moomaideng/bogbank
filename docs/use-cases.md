# Business Use Cases

Per course convention, Sign In / Sign Out are infrastructure, not counted as Business Use Cases (see `supplimentary-materials/convention-for-work.md`). These are the 3 primary Business Use Cases for the project (name TBD, built by group หมูไม่เด้ง).

Supporting use cases (`Manage Bank Expense Hook`, `Manage Category`) are CRUD utilities that the primary use cases depend on — they show up as `<<include>>` in the use case diagram, not as one of the 3.

## UC-01: Ingest Expense from Bank Receipt (Share-to-App)

**Primary Actor**: Customer (college student)

**Goal**: Record an expense with minimal manual typing by sharing a bank payment/transfer confirmation screenshot instead of entering the transaction by hand.

**Preconditions**: Customer is signed in; customer has at least one bank expense hook configured (see `Manage Bank Expense Hook`).

**Main Flow**:
1. Customer receives or takes a screenshot of a bank payment confirmation on their phone.
2. Customer shares the screenshot to the app via the OS share sheet, selecting the relevant bank hook.
3. System extracts metadata (amount, merchant/counterparty, date) from the image.
4. System creates a pending expense record and notifies the customer with a suggested category.
5. Customer confirms the suggested category or picks a different one (`<<include>> Manage Category`).
6. System finalizes the expense record under the confirmed category.

**Postcondition**: A categorized expense record exists and is reflected in the dashboard.

**Outcome the customer wants**: Their spending gets tracked accurately without manually re-typing every bank transaction.

## UC-02: Track Income and Expenses Manually

**Primary Actor**: Customer

**Goal**: Record or correct income/expense transactions that aren't covered by a bank-hook receipt share — cash transactions, other banks, a missed share, or fixing a mistake.

**Preconditions**: Customer is signed in.

**Main Flow**:
1. Customer opens "Add income" or "Add expense".
2. Customer enters amount, date, category (`<<include>> Manage Category`), and an optional note.
3. Customer saves the record.
4. Customer may later edit or remove an existing income/expense record from either flow (including ones created via UC-01).

**Postcondition**: The customer's income/expense record set is complete and accurate.

**Outcome the customer wants**: Nothing falls through the cracks even when the automatic receipt-share flow doesn't apply.

## UC-03: Review Income/Expense Dashboard

**Primary Actor**: Customer

**Goal**: Understand spending and income patterns over time and by bank provider, to make better financial decisions — this is the core problem the product solves (expenses are tedious to track and categorize across many channels).

**Preconditions**: Customer is signed in; customer has at least one income or expense record.

**Main Flow**:
1. Customer opens the dashboard.
2. Customer views a donut chart of income or expense broken down by category, filterable by time range (current/last month, current/last year, or a custom range) and bank provider.
3. Customer switches to a bar chart comparing income vs. expense across a fixed set of time intervals (week/month/3-months/year), filterable by bank provider.

**Postcondition**: None (read-only use case).

**Outcome the customer wants**: An at-a-glance, filterable view of their financial activity without manual spreadsheet work.

## Supporting Use Cases (not counted among the 3)

- **Manage Bank Expense Hook** — Customer: add/edit/remove a bank provider hook. Prerequisite for UC-01.
- **Manage Category** — Customer: add/edit/remove categories used to classify income/expense records. Included by UC-01 and UC-02.
