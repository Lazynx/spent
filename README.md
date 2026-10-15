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
internal/pkg/logger/    slog logger constructor
testdata/               synthetic fixtures only
```

Domain packages (`statement`, `task`, `analyzer`) are added next.

## Development

```bash
make check   # gofmt, go vet, go test
make run     # go run ./cmd/spent --log-level=debug --log-format=json
```

Logs go to stderr; stdout is reserved for the report.
