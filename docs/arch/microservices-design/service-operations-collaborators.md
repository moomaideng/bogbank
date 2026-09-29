# Service–Operations–Collaborators (v2)

| Service | Operations | Collaborators |
| :---- | :---- | :---- |
| Receipt Service | • UploadReceipt()<br>• GetReceipt()<br>• AddBankHook()<br>• EditBankHook()<br>• RemoveBankHook()<br>• ApplyExtractionResult()<br>• ApplySuggestion()<br>• MarkConfirmed() | Object Storage<br>Redpanda<br>PostgreSQL (receipt) |
| Receipt Reader Service | • ExtractMetadata() | Object Storage<br>Receipt OCR / VLM API<br>Redpanda<br>Redis |
| Suggestion Service | • SuggestCategory() | Ledger Service<br>• GetCategories()<br>Redpanda<br>Redis |
| Ledger Service | • AddIncome()<br>• EditIncome()<br>• RemoveIncome()<br>• AddExpense()<br>• EditExpense()<br>• RemoveExpense()<br>• CreateExpenseFromReceipt()<br>• AddCategory()<br>• EditCategory()<br>• RemoveCategory()<br>• GetCategories()<br>• GetRecords() | Receipt Service<br>• GetReceipt()<br>Redpanda<br>PostgreSQL (ledger) |
| Dashboard Service | • GetDonutChartData()<br>• GetBarChartData() | ClickHouse<br>Redpanda |
