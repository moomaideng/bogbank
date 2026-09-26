# Bogbank

Personal expense tracker. Course monorepo: Go backend services, shared libraries, and (later) the mobile client.

## Prerequisites

- [Go](https://go.dev/) (see `go.mod`)
- [Docker](https://www.docker.com/) + Compose
- [make](https://www.gnu.org/software/make/)
- `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc` on `PATH` (only needed when regenerating stubs)

## Get started

From the repo root:

```bash
make dev-ledger
```

When you're done:

```bash
make down
```

`make dev-ledger` also applies pending migrations before it listens. Proto stubs are already committed; run `make proto` only after changing `.proto` files.

- HTTP: `http://localhost:8080`
- gRPC: `localhost:9090`
- Health: `GET /livez`, `GET /readyz`

Config overrides use the `LEDGER_` prefix (`LEDGER_HTTP_ADDRESS`, `LEDGER_DATABASE_DSN`, …).

### API testing

[`bruno/`](bruno/) contains a [Bruno](https://www.usebruno.com/) collection for testing REST and gRPC requests manually — open it in the Bruno app (env `local`) against a running local server.

## Repository

| Path | What |
|---|---|
| [`bruno`](bruno/) | Bruno collection (REST + gRPC) |
| [`services/template`](services/template/README.md) | Copy this when adding a service |
| [`docs/arch/adr`](docs/arch/adr/README.md) | Architecture decisions |
