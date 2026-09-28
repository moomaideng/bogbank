# Agent notes

## Layout

Backend and mobile client live in one repo, in separate top-level trees. No `go.work`; the Go module is `github.com/moomaideng/bogbank` under `backend/`.

- `backend/internal/` is shared: config, database, migrator, baserepo, httpserver, objectstorage.
- A service lives at `backend/services/<name>/`. Its `internal/` is private to that service.
- Inside a service: `handler` → `usecase` → `adapter/<technology>`. HTTP and gRPC edges split into `handler/rest` and `handler/grpc`.
- The usecase owns outbound interfaces. Adapters implement them.
- `backend/services/<name>/main.go` is the cobra entrypoint. `serve --config` loads YAML. Env prefix overrides nested keys (`TEMPLATE_HTTP_ADDRESS` → `http.address`).
- `mobile/` is the Expo (React Native) client. See `mobile/README.md` once it exists.

Coding style is [ADR-11](docs/arch/adr/11_backend-coding-style-v2.md), which supersedes [ADR-10](docs/arch/adr/10_backend-coding-style.md). HTTP is chi + Huma ([ADR-09](docs/arch/adr/09_backend-http-api-stack.md)).

## Comments

Comment only when the code cannot say it. Prefer HOW: a constraint, an ordering, a reason the next line is shaped that way.

Explain WHAT only when the code cannot be made clear. A comment is the last resort, not a caption.

Skip comments that restate the type, the function name, or the architecture. Package comments are one sentence of ownership.
