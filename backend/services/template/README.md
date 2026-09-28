# Service template

Copy this directory when adding a service. Rename the module path under `services/<name>/` and the `TEMPLATE` environment prefix in `internal/config`.

One Go module, `github.com/moomaideng/bogbank`, covers the repo. Shared clients live in the repo-root `internal/` tree. This service's code lives in `internal/` here, which other services cannot import.

## Layout

```text
services/template/
  config.yaml
  main.go                       # cobra entrypoint
  internal/
    cli/                        # serve --config
    config/                     # this service's YAML struct
    handler/                    # Huma operations
    usecase/                    # rules and outbound interfaces
    adapter/postgres/           # bun repositories
    adapter/objectstorage/      # S3 calls for this service
```

A broker adapter, when one exists, goes in `internal/adapter/broker/`. ClickHouse, when a service needs it, is another constructor beside `internal/database.NewPostgres`, not a second ORM.

Repo-root packages:

| Package | Role |
|---|---|
| `internal/config` | YAML file plus `PREFIX_NESTED_KEY` env overlay |
| `internal/database` | bun + pgx Postgres |
| `internal/migrator` | goose `migrate` cobra command |
| `internal/baserepo` | CRUD, transactioner, cursor page types |
| `internal/httpserver` | chi, Huma, `GET /livez`, `GET /readyz` |
| `internal/objectstorage` | S3 API client (path-style) |

## Where a feature goes

Handler → usecase → adapter. The usecase owns the interfaces. Adapters implement them and call the shared clients.

```text
internal/handler/expense.go
internal/usecase/expense.go              # interface the adapter implements
internal/adapter/postgres/expense.go     # bun via internal/baserepo
internal/adapter/objectstorage/receipt.go
```

`internal/baserepo.CursorInput` is the list-page shape (`First`/`After` forward, `Last`/`Before` backward). The keyset query arrives with the first list endpoint.

## Run

From `backend/`, with no database:

```bash
go run ./services/template serve --config services/template/config.yaml
```

`GET /livez` is always 200. `GET /readyz` pings Postgres and the S3 bucket when those blocks are present. Omitting a block leaves that pointer nil.

Environment variables override the file. `TEMPLATE_HTTP_ADDRESS` replaces `http.address`. The same rule applies to every other key (`TEMPLATE_DATABASE_DSN`, `TEMPLATE_S3_BUCKET`, and so on).
