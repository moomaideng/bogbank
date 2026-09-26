# Agent notes

## Layout

One module, `github.com/moomaideng/bogbank`. No `go.work`.

- Repo-root `internal/` is shared: config, database, baserepo, httpserver, objectstorage.
- A service lives at `services/<name>/`. Its `internal/` is private to that service.
- Inside a service: `handler` → `usecase` → `adapter/<technology>`.
- The usecase owns outbound interfaces. Adapters implement them.
- `services/<name>/main.go` is the cobra entrypoint. `serve --config` loads YAML. Env prefix overrides nested keys (`TEMPLATE_HTTP_ADDRESS` → `http.address`).

Coding style is [ADR-11](docs/arch/adr/11_backend-coding-style-v2.md), which supersedes [ADR-10](docs/arch/adr/10_backend-coding-style.md). HTTP is chi + Huma ([ADR-09](docs/arch/adr/09_backend-http-api-stack.md)).

## Comments

Comment only when the code cannot say it. Prefer HOW: a constraint, an ordering, a reason the next line is shaped that way.

Explain WHAT only when the code cannot be made clear. A comment is the last resort, not a caption.

Skip comments that restate the type, the function name, or the architecture. Package comments are one sentence of ownership.
