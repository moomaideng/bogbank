# Service–Operations–Collaborators (v1)

| Service | Operations | Collaborators |
| :---- | :---- | :---- |
| Receipt Ingestion Service | • UploadReceipt()<br>• ExtractMetadata()<br>• SuggestCategory()<br>• SubmitExpenseCandidate()<br>• AddBankHook()<br>• EditBankHook()<br>• RemoveBankHook() | Ledger Service<br>• GetCategories()<br>• CreateExpenseFromCandidate() |
| Ledger Service | • AddIncome()<br>• EditIncome()<br>• RemoveIncome()<br>• AddExpense()<br>• EditExpense()<br>• RemoveExpense()<br>• CreateExpenseFromCandidate()<br>• AddCategory()<br>• EditCategory()<br>• RemoveCategory()<br>• GetCategories()<br>• GetRecords() | — |
| Dashboard Service | • GetDonutChartData()<br>• GetBarChartData() | Ledger Service<br>• GetRecords() |
