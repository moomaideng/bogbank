# Backlog Export

Source of truth for the sprint plan agreed in `agenda.md`. This file is structured for a future
agent to parse and create issues in **GitHub Issues** or **Linear** — do not hand-edit issue content
here without also updating `agenda.md`, and vice versa.

## Format

Every issue is an `##`-level heading followed by a fenced `yaml` metadata block and a body. An
importer should split the document on `## ` headings, parse the `yaml` block, and treat everything
after it (until the next heading) as the issue body.

```yaml
id: BL-000          # stable id, referenced by other issues' `parent`/`blocks`
type: epic | task    # epic = one per user story (or platform); task = actual backlog item
epic: BL-000         # parent epic id (tasks only)
sprint: 0 | 1 | 2 | 3 # 0 = shared/foundation, done before Sprint 1 branches
points: 5             # Fibonacci story points (see agenda.md for scale/capacity assumptions)
service: auth-svc     # owning service, or "platform" / "mobile" / "cross-cutting"
owner: P1             # placeholder slot from agenda.md — replace with a real handle at import time
labels: [backend, grpc]
stretch: true         # optional; omit when false
```

### Field mapping for import

| Field | GitHub Issues | Linear |
| :---- | :---- | :---- |
| `## Title` | issue title | issue title |
| body (below yaml) | issue body | issue description |
| `labels` | `labels` | `labels` |
| `sprint` | milestone (`Sprint N`) | cycle |
| `epic` | task list / `Epic: BL-00N` in body, or GitHub sub-issue parent | parent issue |
| `points` | not native — add a `points:N` label or a project custom field | `estimate` |
| `service` | label `svc:<service>` | label `svc:<service>` |
| `owner` | assignee (map `P1`.. to GitHub usernames) | assignee (map `P1`.. to Linear user) |

Owner slots (`P1`–`P5`) map to team members per the meeting; substitute real usernames at import
time. Group members: Ashira Aungsumal, Kittichet Arayasujin, Kittichon Chaonawig, Pasin Thanyakasikol,
Peeravas Piboolvorakul (`docs/project-description.md`).

---

## Epics

### Login & Register

```yaml
id: BL-EPIC-US1
type: epic
sprint: 1
service: auth-svc
owner: P1
labels: [epic, us1]
```

FR1, FR2, FR3, NFR1. Google OAuth 2.0 sign-in; first successful sign-in auto-creates the account.
No passwords are ever stored.

### Logout

```yaml
id: BL-EPIC-US2
type: epic
sprint: 1
service: auth-svc
owner: P1
labels: [epic, us2]
```

FR3. Ends the current session; refresh token revoked.

### Manage Bank Expense Hook

```yaml
id: BL-EPIC-US3
type: epic
sprint: 1
service: receipt-svc
owner: P3
labels: [epic, us3]
```

FR4. Add / edit / remove a bank expense hook; grant media-library permission. CRUD only — the
gallery-watch ingestion pipeline itself is US5.

### Manage Income & Expense Record

```yaml
id: BL-EPIC-US4
type: epic
sprint: 1
service: ledger-svc
owner: P2
labels: [epic, us4]
```

FR9, FR10, FR11, NFR7. Manual add/edit/remove of income, expense, and category records — always
available as the fallback path.

### Bank Hook Ingestion

```yaml
id: BL-EPIC-US5
type: epic
sprint: 2
service: receipt-svc
owner: P3
labels: [epic, us5]
```

FR5, FR6, FR7, FR8. Detect a new receipt image, upload it, extract metadata, suggest a category,
let the user confirm with one tap, create the expense. The full async choreography
(Receipt → Reader → Suggestion → Receipt → Ledger over Redpanda).

### View Income & Expense Dashboard

```yaml
id: BL-EPIC-US6
type: epic
sprint: 1
service: dashboard-svc
owner: P4
labels: [epic, us6]
```

FR12, FR13, NFR8. Donut chart and bar chart, filterable by type/range/bank provider. Functionally
complete via Ledger gRPC alone (UC-03 depends only on UC-02); ClickHouse is a Sprint 3 performance
upgrade, not a functional prerequisite.

### Project Setup & Platform

```yaml
id: BL-EPIC-0
type: epic
sprint: 0
service: platform
owner: P5
labels: [epic, platform]
```

Not a user story. Shared foundation (monorepo, contracts, CI, Compose) plus cross-cutting platform
work (observability, deployment) that every other epic depends on.

---

## Sprint 0 — shared foundation

### Monorepo layout, Go workspace, service template

```yaml
id: BL-001
type: task
epic: BL-EPIC-0
sprint: 0
points: 5
service: platform
owner: whole-team
labels: [platform, setup]
```

chi + Huma HTTP stack, hexagonal port/adapter folder layout (ADR-09, ADR-10) as the template every
service copies.

**Acceptance criteria**
- [ ] `go.work` wires all service modules
- [ ] One reference service compiles and serves a health-check endpoint
- [ ] Folder layout documented in the template's README

### Protobuf contracts + event envelope schema

```yaml
id: BL-002
type: task
epic: BL-EPIC-0
sprint: 0
points: 3
service: platform
owner: whole-team
labels: [platform, grpc, kafka]
```

`GetReceipt`, `GetCategories`, `GetRecords` service definitions, plus the Kafka event envelope
(`event_id`, `event_type`, `occurred_at`, `user_id`, `trace_id`, `payload`). The only cross-slice
coupling — agreed once, in the room.

**Acceptance criteria**
- [ ] `.proto` files committed, generate Go stubs via one Makefile target
- [ ] Event envelope schema documented in `architecture-diagram-v2.md`'s event catalog

### Shared outbox + eventbus library

```yaml
id: BL-003
type: task
epic: BL-EPIC-0
sprint: 0
points: 8
service: platform
owner: whole-team
labels: [platform, kafka, reliability]
```

`pkg/outbox` (transactional outbox write + relay) and `pkg/eventbus` (producer, consumer, Redis
idempotency dedupe). Moved out of Auth SVC's task list into the shared session — every Sprint 2
producer/consumer depends on it.

**Acceptance criteria**
- [ ] Outbox table schema + relay process reusable by any service via one import
- [ ] Consumer wrapper rejects already-processed `event_id` via Redis

### Expo app scaffold + route convention

```yaml
id: BL-004
type: task
epic: BL-EPIC-0
sprint: 0
points: 3
service: mobile
owner: P5
labels: [platform, mobile]
```

`app/(tabs)/<slice>/` route-folder convention so five people can build screens without merge
conflicts.

**Acceptance criteria**
- [ ] Expo app builds and runs on a real device
- [ ] Route convention documented for the other four owners

### Compose infra profile + Makefile targets

```yaml
id: BL-005
type: task
epic: BL-EPIC-0
sprint: 0
points: 3
service: platform
owner: P5
labels: [platform, docker]
```

`infra` profile (postgres, redis, traefik, rustfs) always on; per-service profiles layer on top.

**Acceptance criteria**
- [ ] `make dev SVC=<name>` starts one service and its dependencies only
- [ ] `make dev-full` starts everything

### CI pipeline

```yaml
id: BL-006
type: task
epic: BL-EPIC-0
sprint: 0
points: 2
service: platform
owner: P5
labels: [platform, ci]
```

**Acceptance criteria**
- [ ] PR triggers build, `go vet`, unit tests, lint
- [ ] Red CI blocks merge

---

## Sprint 1 — US1, US2, US3, US4, US6

### P1 — Auth SVC + edge

```yaml
id: BL-101
type: task
epic: BL-EPIC-US1
sprint: 1
points: 8
service: auth-svc
owner: P1
labels: [backend, auth, oauth]
```

Google OAuth 2.0 + PKCE, deep-link redirect back into Expo.

**Acceptance criteria**
- [ ] Sign-in completes end-to-end on a real device
- [ ] Deep link returns control to the app with the auth code

```yaml
id: BL-102
type: task
epic: BL-EPIC-US1
sprint: 1
points: 3
service: auth-svc
owner: P1
labels: [backend, database]
```

`users` + `refresh_tokens` schema; first sign-in auto-creates the account (FR2).

**Acceptance criteria**
- [ ] First-time Google sign-in creates a row with no separate registration step
- [ ] Second sign-in reuses the existing account

```yaml
id: BL-103
type: task
epic: BL-EPIC-US2
sprint: 1
points: 5
service: auth-svc
owner: P1
labels: [backend, auth]
```

Issue/verify BogBank JWT, refresh, logout with revocation in Redis.

**Acceptance criteria**
- [ ] Expired access token is rejected; refresh token issues a new one
- [ ] Logout revokes the refresh token; a revoked token is rejected

```yaml
id: BL-104
type: task
epic: BL-EPIC-US1
sprint: 1
points: 5
service: auth-svc
owner: P1
labels: [backend, grpc]
```

gRPC `Introspect` for internal callers.

**Acceptance criteria**
- [ ] Another service can validate a JWT via gRPC without hitting Google

```yaml
id: BL-105
type: task
epic: BL-EPIC-US1
sprint: 1
points: 5
service: platform
owner: P1
labels: [backend, gateway, traefik]
```

Traefik: routing, TLS, ForwardAuth middleware wired to Auth SVC.

**Acceptance criteria**
- [ ] Unauthenticated request to a protected route is rejected by Traefik before reaching the service
- [ ] Routes discovered from Docker labels, no static route config

```yaml
id: BL-106
type: task
epic: BL-EPIC-US1
sprint: 1
points: 5
service: mobile
owner: P1
labels: [frontend, auth]
```

Sign-in / sign-out screens, secure token storage, refresh interceptor in the API client.

**Acceptance criteria**
- [ ] Token persists across app restart
- [ ] Expired token triggers silent refresh, not a forced re-login

### P2 — Ledger SVC

```yaml
id: BL-201
type: task
epic: BL-EPIC-US4
sprint: 1
points: 3
service: ledger-svc
owner: P2
labels: [backend, database]
```

`categories`, `records`, `outbox` schema; `bank_provider` column; `UNIQUE (receipt_id)`.

**Acceptance criteria**
- [ ] Schema migration applies cleanly
- [ ] `UNIQUE (receipt_id)` present (needed by Sprint 2's confirm flow)

```yaml
id: BL-202
type: task
epic: BL-EPIC-US4
sprint: 1
points: 3
service: ledger-svc
owner: P2
labels: [backend, rest]
```

REST: category CRUD (FR11).

**Acceptance criteria**
- [ ] Add/edit/remove category via REST, validated against OpenAPI spec

```yaml
id: BL-203
type: task
epic: BL-EPIC-US4
sprint: 1
points: 3
service: ledger-svc
owner: P2
labels: [backend, rest]
```

REST: income CRUD (FR9).

**Acceptance criteria**
- [ ] Add/edit/remove income record via REST

```yaml
id: BL-204
type: task
epic: BL-EPIC-US4
sprint: 1
points: 3
service: ledger-svc
owner: P2
labels: [backend, rest]
```

REST: expense CRUD (FR10).

**Acceptance criteria**
- [ ] Add/edit/remove expense record via REST (manual fallback path, NFR7)

```yaml
id: BL-205
type: task
epic: BL-EPIC-US4
sprint: 1
points: 5
service: ledger-svc
owner: P2
labels: [backend, grpc]
```

gRPC server: `GetCategories`, `GetRecords`.

**Acceptance criteria**
- [ ] Suggestion SVC (Sprint 2) and Dashboard SVC can both call these today

```yaml
id: BL-206
type: task
epic: BL-EPIC-US4
sprint: 1
points: 8
service: mobile
owner: P2
labels: [frontend, ledger]
```

Add / edit / remove income and expense screens; category management (NFR7 manual fallback path).

**Acceptance criteria**
- [ ] Full manual CRUD flow usable without any receipt ever being ingested

### P3 — Receipt SVC

```yaml
id: BL-301
type: task
epic: BL-EPIC-US3
sprint: 1
points: 3
service: receipt-svc
owner: P3
labels: [backend, database]
```

`bank_hooks`, `receipts`, `receipt_events`, `outbox` schema.

**Acceptance criteria**
- [ ] Schema migration applies cleanly
- [ ] `UNIQUE (user_id, image_sha256)` present (needed by the upload dedupe rule)

```yaml
id: BL-302
type: task
epic: BL-EPIC-US3
sprint: 1
points: 5
service: receipt-svc
owner: P3
labels: [backend, rest]
```

REST: bank hook CRUD (FR4), supported-provider list.

**Acceptance criteria**
- [ ] Add/edit/remove a bank hook selecting from supported providers (KBank, SCB, ...)

```yaml
id: BL-303
type: task
epic: BL-EPIC-US3
sprint: 1
points: 3
service: receipt-svc
owner: P3
labels: [backend, grpc]
```

gRPC server: `GetReceipt`.

**Acceptance criteria**
- [ ] Ledger SVC (Sprint 2 confirm flow) can fetch a receipt's amount/merchant/bank_provider today

```yaml
id: BL-304
type: task
epic: BL-EPIC-US5
sprint: 1
points: 5
service: receipt-svc
owner: P3
stretch: true
labels: [backend, rest, stretch]
```

Stretch, only if hook CRUD lands early: upload endpoint (sha256 dedupe, S3 put,
`UNIQUE (user_id, image_sha256)`). Not required for US3; pure REST with no Kafka dependency, so it
can ship into its own outbox ahead of Sprint 2's consumers existing.

**Acceptance criteria**
- [ ] Duplicate image (same user, same hash) is rejected, not double-stored

```yaml
id: BL-305
type: task
epic: BL-EPIC-US3
sprint: 1
points: 8
service: mobile
owner: P3
labels: [frontend, receipt]
```

Bank hook management screen, media-library permission request (NFR2).

**Acceptance criteria**
- [ ] Permission is requested only when configuring a hook, not on app launch (NFR2)

### P4 — Dashboard SVC

```yaml
id: BL-401
type: task
epic: BL-EPIC-US6
sprint: 1
points: 3
service: dashboard-svc
owner: P4
labels: [backend, grpc]
```

Dashboard SVC reading Ledger over gRPC — sufficient for US6 end-to-end since UC-03 only requires
manual records (UC-02). This is the permanent ClickHouse-down fallback, not a temporary hack.

**Acceptance criteria**
- [ ] Dashboard SVC has zero dependency on ClickHouse or Redpanda in Sprint 1

```yaml
id: BL-402
type: task
epic: BL-EPIC-US6
sprint: 1
points: 5
service: dashboard-svc
owner: P4
labels: [backend, rest]
```

Donut endpoint (FR12) with type / range / bank-provider filters.

**Acceptance criteria**
- [ ] All filter combinations from FR12 return correct aggregates against seeded data

```yaml
id: BL-403
type: task
epic: BL-EPIC-US6
sprint: 1
points: 5
service: dashboard-svc
owner: P4
labels: [backend, rest]
```

Bar endpoint (FR13) with interval / bank-provider filters.

**Acceptance criteria**
- [ ] Week/month/3-month/year intervals each return two series (income, expense)

```yaml
id: BL-404
type: task
epic: BL-EPIC-US6
sprint: 1
points: 8
service: mobile
owner: P4
labels: [frontend, dashboard]
```

Donut chart, bar chart, filter controls.

**Acceptance criteria**
- [ ] Empty state renders a friendly message, not a blank/error chart (UC-03 alternate flow)

### P5 — Platform + app shell

```yaml
id: BL-501
type: task
epic: BL-EPIC-0
sprint: 1
points: 5
service: platform
owner: P5
labels: [platform, docker]
```

Compose profiles for every service + `full`; hot reload via `air`.

**Acceptance criteria**
- [ ] Each of the 5 owners can run `make dev SVC=<theirs>` independently

```yaml
id: BL-502
type: task
epic: BL-EPIC-0
sprint: 1
points: 8
service: platform
owner: P5
labels: [platform, observability]
```

Grafana + Prometheus + Loki; OpenTelemetry SDK wired into the service template.

**Acceptance criteria**
- [ ] Any service's logs and a basic request-rate metric are visible in Grafana out of the box

```yaml
id: BL-503
type: task
epic: BL-EPIC-0
sprint: 1
points: 2
service: platform
owner: P5
labels: [platform, testing]
```

k6 harness stub for the Sprint 3 load tests.

**Acceptance criteria**
- [ ] `make loadtest` runs a trivial smoke scenario against a running stack

```yaml
id: BL-504
type: task
epic: BL-EPIC-US5
sprint: 1
points: 3
service: platform
owner: P5
labels: [backend, scaffold]
```

Reader SVC and Suggestion SVC skeletons (no logic yet — that's Sprint 2).

**Acceptance criteria**
- [ ] Both services build, start under Compose, and pass a health check

```yaml
id: BL-505
type: task
epic: BL-EPIC-0
sprint: 1
points: 8
service: mobile
owner: P5
labels: [frontend, platform]
```

App shell, navigation, shared UI kit, OpenAPI → typed API client codegen.

**Acceptance criteria**
- [ ] A new screen under any owner's route folder gets navigation and theming for free
- [ ] API client regenerates from any service's OpenAPI spec via one command

---

## Sprint 2 — US5 alone

### P1 — Broker + reliability hardening

```yaml
id: BL-601
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: platform
owner: P1
labels: [backend, kafka]
```

Redpanda topics, partitioning, ACLs / prod config.

```yaml
id: BL-602
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: platform
owner: P1
labels: [backend, reliability]
```

Outbox relay hardening: retry/backoff, crash-safety.

**Acceptance criteria**
- [ ] Killing the relay mid-publish does not lose or duplicate an event on restart

```yaml
id: BL-603
type: task
epic: BL-EPIC-US5
sprint: 2
points: 3
service: platform
owner: P1
labels: [backend, reliability]
```

Dead-letter topic tooling and admin visibility.

```yaml
id: BL-604
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: platform
owner: P1
labels: [backend, observability]
```

`trace_id` / OTel span propagation across the event envelope.

**Acceptance criteria**
- [ ] A single receipt's trace is visible end-to-end across Receipt → Reader → Suggestion → Ledger

### P2 — Ledger event production + confirm flow

```yaml
id: BL-701
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: ledger-svc
owner: P2
labels: [backend, kafka]
```

`ledger.expense.*` / `ledger.income.*` event emission via outbox.

```yaml
id: BL-702
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: ledger-svc
owner: P2
labels: [backend, grpc]
```

Confirm flow: `CreateExpense` calls `GetReceipt` over gRPC, transactional insert.

**Acceptance criteria**
- [ ] Double-tapping "confirm" does not create a second expense (`UNIQUE (receipt_id)`)

```yaml
id: BL-703
type: task
epic: BL-EPIC-US5
sprint: 2
points: 8
service: mobile
owner: P2
labels: [frontend, ledger]
```

Confirm / category-picker screen reached from the push notification.

### P3 — Receipt SVC state machine

```yaml
id: BL-801
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: receipt-svc
owner: P3
labels: [backend, rest]
```

Upload endpoint w/ dedupe, if not already shipped as Sprint 1 stretch (`BL-304`).

```yaml
id: BL-802
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: receipt-svc
owner: P3
labels: [backend]
```

Receipt state machine (`uploaded → extracted → suggested → confirmed`, `failed`).

```yaml
id: BL-803
type: task
epic: BL-EPIC-US5
sprint: 2
points: 8
service: receipt-svc
owner: P3
labels: [backend, kafka]
```

3 Kafka consumers (`receipt.extracted`, `receipt.suggested`, `ledger.expense.created`) with Redis
event-id idempotency.

**Acceptance criteria**
- [ ] Redelivering the same event does not double-apply its effect

```yaml
id: BL-804
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: receipt-svc
owner: P3
labels: [backend, push]
```

Push notification integration (Expo Push SDK, device token management).

```yaml
id: BL-805
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: mobile
owner: P3
labels: [frontend, receipt]
```

Receipt status / notification list UI.

### P4 — Gallery watch (reassigned to help P3)

```yaml
id: BL-901
type: task
epic: BL-EPIC-US5
sprint: 2
points: 13
service: mobile
owner: P4
labels: [frontend, receipt, mobile-native]
```

Gallery-watch native listener — device background limits, iOS/Android divergence. The single
hardest item in the backlog; reassigned here because Sprint 1's ClickHouse deferral freed P4's
Sprint 2 capacity.

**Acceptance criteria**
- [ ] A new receipt image saved to the gallery is detected without the user opening the app
      (foreground-triggered scan is an acceptable documented fallback for iOS background limits)

```yaml
id: BL-902
type: task
epic: BL-EPIC-US6
sprint: 2
points: 2
service: dashboard-svc
owner: P4
labels: [backend, verification]
```

Verify the Ledger-gRPC fallback still holds once real ingestion traffic exists.

### P5 — Reader SVC + Suggestion SVC

```yaml
id: BL-1001
type: task
epic: BL-EPIC-US5
sprint: 2
points: 8
service: reader-svc
owner: P5
labels: [backend, ocr]
```

Reader SVC: fetch from S3, VLM extraction call, retry/backoff, DLQ on failure.

```yaml
id: BL-1002
type: task
epic: BL-EPIC-US5
sprint: 2
points: 3
service: reader-svc
owner: P5
labels: [backend, kafka]
```

Reader SVC: consumer wiring for `receipt.uploaded`.

```yaml
id: BL-1003
type: task
epic: BL-EPIC-US5
sprint: 2
points: 5
service: suggestion-svc
owner: P5
labels: [backend, grpc]
```

Suggestion SVC: `GetCategories` gRPC client, merchant-to-category rule matching, event emission.

```yaml
id: BL-1004
type: task
epic: BL-EPIC-US5
sprint: 2
points: 3
service: suggestion-svc
owner: P5
labels: [backend, kafka]
```

Suggestion SVC: consumer wiring for `receipt.extracted`.

---

## Sprint 3 — Non-functional requirements

```yaml
id: BL-1101
type: task
epic: BL-EPIC-US6
sprint: 3
points: 13
service: dashboard-svc
owner: P4
labels: [backend, clickhouse, performance]
```

ClickHouse projection (moved from Sprint 2): schema, `AggregatingMergeTree` rollups, Kafka
consumer/projector, incremental materialized view.

**Acceptance criteria**
- [ ] Before/after p95 latency comparison vs. the Sprint 1 Ledger-gRPC path is captured for the demo

```yaml
id: BL-1102
type: task
epic: BL-EPIC-0
sprint: 3
points: 8
service: platform
owner: P5
labels: [testing, load-test]
```

Load testing, 2 types (constant-arrival-rate + spike/stress) + interpretation writeup.

```yaml
id: BL-1103
type: task
epic: BL-EPIC-0
sprint: 3
points: 8
service: platform
owner: P1+P5
labels: [testing, quality-attribute]
```

Quality-attribute measurement: scale Reader SVC consumer replicas, kill ClickHouse mid-demo,
measure p95 ingest latency.

```yaml
id: BL-1104
type: task
epic: BL-EPIC-0
sprint: 3
points: 8
service: auth-svc
owner: P1
labels: [security]
```

Security hardening: TLS everywhere, secrets management, PDPA compliance review (NFR3, NFR4),
rate limits.

```yaml
id: BL-1105
type: task
epic: BL-EPIC-0
sprint: 3
points: 3
service: cross-cutting
owner: whole-team
labels: [docs]
```

Risk matrix, formalized to the graded format (≥3 issues, likelihood/impact/mitigation).

```yaml
id: BL-1106
type: task
epic: BL-EPIC-0
sprint: 3
points: 5
service: platform
owner: P5
labels: [observability]
```

Observability polish: Grafana dashboards, alert rules, Tempo tracing (stretch).

```yaml
id: BL-1107
type: task
epic: BL-EPIC-0
sprint: 3
points: 5
service: cross-cutting
owner: P2+P3
labels: [docs, c4]
```

C4 L3 component diagrams for Ledger SVC and Receipt SVC internals.

```yaml
id: BL-1108
type: task
epic: BL-EPIC-0
sprint: 3
points: 3
service: cross-cutting
owner: whole-team
labels: [docs, c4]
```

C4 L4 code-level hexagonal diagram, one service.

```yaml
id: BL-1109
type: task
epic: BL-EPIC-0
sprint: 3
points: 8
service: cross-cutting
owner: whole-team
labels: [docs, adr]
```

ADR-11 through ADR-18 — closes the deferred Tier B list (`docs/arch/adr/README.md`).

```yaml
id: BL-1110
type: task
epic: BL-EPIC-0
sprint: 3
points: 8
service: cross-cutting
owner: whole-team
labels: [presentation]
```

Slide deck (C4 model, per presentation guidelines) + demo rehearsal.

```yaml
id: BL-1111
type: task
epic: BL-EPIC-0
sprint: 3
points: 8
service: cross-cutting
owner: whole-team
labels: [stabilization]
```

Stabilization buffer.

```yaml
id: BL-1112
type: task
epic: BL-EPIC-0
sprint: 3
points: 13
service: platform
owner: TBD
stretch: true
labels: [platform, k8s, stretch]
```

k3s single-node deployment, Kustomize manifests, Traefik IngressRoute. **Optional — cut first if
Sprint 3 runs over.** Compose-on-VM is the guaranteed demo path; this is additive.

---

## Totals (cross-check against `agenda.md`)

| Sprint | Points | Matches `agenda.md` as |
| :---- | :---- | :---- |
| 0 (foundation) | 24 | folded into "Sprint 1 total" |
| 1 | 122 (+5 stretch) | 24 + 122 = **146** |
| 2 | 98 | **98** |
| 3 | 77 (+13 optional) | **77** (+13 optional) |
| **Total** | **321** (+18 stretch/optional) | |
