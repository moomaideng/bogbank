# Bogbank

Personal expense tracker. Course monorepo: Go backend services, shared libraries, and (later) the mobile client.

## Prerequisites

- [Docker](https://www.docker.com/) + Compose
- [Task](https://taskfile.dev/installation/)

[Go 1.27](https://go.dev/) and `protoc` are optional (IDE use). The documented loop runs everything in Compose with pinned tooling.

## Get started

From the repo root:

```bash
task dev:ledger
```

When you're done:

```bash
task down
```

`task dev:ledger` starts infra (e.g., Postgres and RustFS) and ledger service with hot reload via [air](https://github.com/air-verse/air). Migrations apply on serve. Other presets:

```bash
task dev SVC=ledger   # same as dev:ledger
task dev:full         # infra + every app service defined in Compose
task proto            # regenerate gRPC stubs (pinned container; no host protoc)
```

Proto stubs are committed; run `task proto` only after changing `.proto` files.

- HTTP: `http://localhost:8080`
- gRPC: `localhost:9090`
- Health: `GET /livez`, `GET /readyz`
- S3 API (RustFS): `http://localhost:9000` (console `:9001`; keys `admin` / `admin1234`)

Config overrides use the `LEDGER_` prefix (`LEDGER_HTTP_ADDRESS`, `LEDGER_DATABASE_DSN`, …).

### API testing

[`bruno/`](bruno/) contains a [Bruno](https://www.usebruno.com/) collection for testing REST and gRPC requests manually — open it in the Bruno app (env `local`) against a running local server.

## Repository

| Path | What |
|---|---|
| [`bruno/`](bruno/) | Bruno collection (REST + gRPC) |
| [`docker/`](docker/) | Dev images (Go+air, pinned protoc); add runtime-specific Dockerfiles here later |
| [`services/template/`](services/template/README.md) | Copy this when adding a service |
| [`docs/arch/adr/`](docs/arch/adr/README.md) | Architecture decisions |
| [`Taskfile.yml`](Taskfile.yml) | Local developer commands |
| [`.air.toml`](.air.toml) | Shared hot-reload config for Go services (`SERVICE` selects which) |
