# spent

Bank statement analyzer in Go. Takes a list of transactions, categorizes
spending, computes monthly totals, detects subscriptions and anomalies, and
produces a report.

## Task lifecycle

```
Pending ──► Processing ──► Completed
                      └──► Failed
```

| From         | To           | Allowed |
|--------------|--------------|---------|
| Pending      | Processing   | yes     |
| Processing   | Completed    | yes     |
| Processing   | Failed       | yes     |
| anything else |             | no      |

`Completed` and `Failed` are terminal.

## Layout

```
cmd/spent/              CLI entry point, wiring only
internal/statement/     Transaction, Statement, validation, Parser, CSVParser
internal/analyzer/      categorization, totals, subscriptions, anomalies, Report
internal/task/          Task, Status, Store, Service
internal/pkg/logger/    slog logger constructor
testdata/               synthetic fixtures only
```

## CSV format

Amounts and balances are integers in tiyn (1/100 KZT). Negative = expense.

```
opening_balance,100000000
closing_balance,233650000
date,amount,currency,description,mcc
2026-01-05,-1500000,KZT,Magnum,5411
```

## Processing semantics

Validation rejects input at creation; no task is created. Balance
reconciliation runs during processing: a mismatch moves the task to `Failed`
with the reason in `FailReason`, and `Process` returns the task with a nil
error. `Process` returns an error only for an unknown ID or an invalid
transition.

## Development

```bash
make check   # gofmt, go vet, go test
make run     # go run ./cmd/spent testdata/sample.csv
```

Logs go to stderr; stdout is reserved for the report.
