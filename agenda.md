# Agenda

- Who is the PO?
- Presenting Architecture Diagram V2
- Linear or Github Issue?
- List out User Stoires
- List out Tasks for Sprint 1, 2 & 3
- How to Git & Github


## Who is the PO?

I just want someone who leads the meeting, keep track of user stories and tasks.

## Presenting Architecture Diagram V2

See [architecture-diagram-v2.md](docs/arch/microservices-design/architecture-diagram-v2.md)

Decisions to ratify in the meeting (everything else in v2 is already settled):

1. **Redis as the NoSQL store.** Syllabus demands it twice, verbatim: §17 "at least 2 types of
   databases (**both RDBMS and No-SQL**)" and Presentation Guidelines "one relational (RDBMS) and
   **one NoSQL**". ClickHouse self-describes as a SQL DBMS. Redis is needed anyway for consumer
   idempotency, so this costs one container.
2. **Notification stays inside Receipt SVC** (Expo Push) rather than becoming an another service.
3. **Suggestion SVC is rule-based**, not an LLM. Reader SVC's VLM already carries the AI narrative
   for the week 13–14 lectures.
4. **ADR-04 must be superseded** — it says "PostgreSQL for every service" and explicitly rejects
   ClickHouse. Plus ADR-11..18 to close the deferred Tier B list.

Rejected after review, with reasons in the v2 doc: Authelia (cannot federate to Google, stores
password hashes — contradicts ADR-03/NFR1), hot/cold storage tiering (CQRS instead), Iceberg/S3
tiering, service mesh.

## Linear or Github Issue?

Context

- Linear has better UX imo, but you guys already know Github Issue.

## List out User Stoires

- US1: Login (& Register)
- US2: Logout
- US3: Manage Bank Expense Hook
- US4: Manage Income & Expense Record
- US5: Bank Hook Ingestion
- US6: View Income & Expense Dashboard

**Dependency check**: dashboard depends on **US4 only**, not on US5.

Reclassifying accordingly:

**Sprint 1** shall be US1, US2, US3, US4, US6 — everything that only needs Auth + Ledger.
**Sprint 2** shall be **US5 alone** — the ingestion pipeline, isolated so it doesn't share a sprint
with anything else and the team's first Kafka work gets the full 2-week window.
**Sprint 3** shall be Non-functional requirements.

## List out Tasks for Sprint 1, 2 & 3

**Everyone owns backend services. Nobody is the frontend person** — each person builds the screens
that consume their own service, so the frontend is spread across all five. Slices touch disjoint
services and disjoint route folders, so they merge without conflicts.

Full backlog, ready for GitHub Issues / Linear import: [backlog.md](docs/backlog.md).
Ownership rationale and IPC detail: [architecture-diagram-v2.md](docs/arch/microservices-design/architecture-diagram-v2.md), section "Proposed work split".

**Story points**: Fibonacci (1/2/3/5/8/13), sized for a team learning gRPC/Kafka/Traefik/hexagonal
Go/Expo alongside the lectures teaching them. Capacity assumption: **~28 pts/person/sprint**
`[estimate, not measured]`.

| Owner | Story | Service |
| :---- | :---- | :---- |
| P1 | US1, US2 | Auth SVC |
| P2 | US4 | Ledger SVC |
| P3 | US3, US5 | Receipt SVC |
| P4 | US6 | Dashboard SVC |
| P5 | US5 | Reader SVC, Suggestion SVC, Platform |

---

### Sprint 1 — US1, US2, US3, US4, US6 (target: **progress video 1, due 29 Sep – 3 Oct**)

#### Sprint 0 — shared, done together before anyone branches (half a day) · **24 pts**

- Monorepo layout, Go workspace, service template (chi + Huma + hexagonal, ADR-09/ADR-10) — 5
- **Protobuf contracts**: `GetReceipt`, `GetCategories`, `GetRecords` + the event envelope schema — 3
- Shared `pkg/outbox` + `pkg/eventbus` (producer, consumer, Redis idempotency) — 8
  — moved here from Auth SVC's task list: everyone needs to understand it to write their own
  producer/consumer in Sprint 2, and it was overloading one person
- Expo app scaffold + route-folder convention (`app/(tabs)/<slice>/`) — 3
- Docker Compose `infra` profile (postgres, redis, traefik, rustfs) + Makefile targets — 3
- CI: build, vet, test, lint on PR — 2

#### P1 — Auth SVC + edge (US1, US2) · **31 pts**

Backend
- Auth SVC: Google OAuth 2.0 + PKCE, deep-link redirect for Expo — 8
- `users` + `refresh_tokens`; first sign-in auto-creates the account (FR2) — 3
- Issue/verify BogBank JWT, refresh, logout with revocation in Redis — 5
- gRPC `Introspect` for internal callers — 5
- Traefik: routing, TLS, ForwardAuth middleware wired to Auth SVC — 5

Frontend
- Sign-in / sign-out screens, secure token storage, refresh interceptor in the API client — 5

#### P2 — Ledger SVC (US4) · **25 pts**

Backend
- `categories`, `records`, `outbox` schema; `bank_provider` column; `UNIQUE (receipt_id)` — 3
- REST: category CRUD (FR11) — 3
- REST: income CRUD (FR9) — 3
- REST: expense CRUD (FR10) — 3
- gRPC server: `GetCategories`, `GetRecords` — 5

Frontend
- Add / edit / remove income and expense screens; category management (NFR7 manual fallback) — 8

#### P3 — Receipt SVC (US3) · **19 pts** (+5 stretch)

Backend
- `bank_hooks`, `receipts`, `receipt_events`, `outbox` schema — 3
- REST: bank hook CRUD (FR4), supported-provider list — 5
- gRPC server: `GetReceipt` — 3
- *Stretch, only if hook CRUD lands early:* upload endpoint (sha256 dedupe, S3 put,
  `UNIQUE (user_id, image_sha256)`) — 5. Not required for US3, but it's pure REST with no Kafka
  dependency, so it can ship into its own outbox ahead of Sprint 2's consumers existing.

Frontend
- Bank hook management screen, media-library permission request (NFR2) — 8

#### P4 — Dashboard SVC (US6) · **21 pts**

Backend
- Dashboard SVC reading Ledger over **gRPC** — sufficient for US6 end-to-end since UC-03 only
  requires manual records (UC-02); this is also the permanent CH-down fallback, not a temporary
  hack — 3
- Donut endpoint (FR12) with type / range / bank-provider filters — 5
- Bar endpoint (FR13) with interval / bank-provider filters — 5

Frontend
- Donut chart, bar chart, filter controls — 8

ClickHouse projection is **not** Sprint 2 — it moved to Sprint 3 (see below).

#### P5 — Platform + app shell (unblocks everyone) · **26 pts**

Backend / platform
- Compose profiles for every service + `full`; hot reload via `air` — 5
- Grafana + Prometheus + Loki; OpenTelemetry SDK wired into the service template — 8
- k6 harness stub for the Sprint 3 load tests — 2
- Reader SVC and Suggestion SVC skeletons (no logic yet — that's Sprint 2) — 3

Frontend
- App shell, navigation, shared UI kit, OpenAPI → typed API client codegen — 8

**Sprint 1 total: 146 pts** (24 shared + 31 + 25 + 19 + 21 + 26).

#### Definition of done for Sprint 1

Sign in on a real device, add a bank hook, add/edit/remove an income and an expense, see them on a
dashboard chart. Everything behind Traefik, everything startable with `make dev-full`.

Everyone must have merged backend commits by the end of the sprint — per-person GIT contribution is
graded three separate times (5% + 5% + 5%).

---

### Sprint 2 — US5 alone (target: **progress video 2, due 13–17 Oct**)

US5 is isolated in its own sprint because it's the team's first Kafka work and the one place real
distributed-systems risk concentrates (outbox, at-least-once delivery, a brand-new native
gallery-watch listener). **ClickHouse moved in from Sprint 1's "later" note and back out to
Sprint 3** — see the callout below.

#### P1 — Broker + reliability hardening · **18 pts**

- Redpanda topics, partitioning, ACLs / prod config — 5
- Outbox relay hardening: retry/backoff, crash-safety — 5
- Dead-letter topic tooling and admin visibility — 3
- `trace_id` / OTel span propagation across the event envelope — 5

#### P2 — Ledger event production + confirm flow · **18 pts**

- `ledger.expense.*` / `ledger.income.*` event emission via outbox — 5
- Confirm flow: `CreateExpense` calls `GetReceipt` over gRPC, transactional insert — 5
- Frontend: confirm / category-picker screen reached from the push notification — 8

#### P3 — Receipt SVC state machine · **28 pts**

- Upload endpoint w/ dedupe, if not already shipped as Sprint 1 stretch — 5
- Receipt state machine (`uploaded → extracted → suggested → confirmed`, `failed`) — 5
- 3 Kafka consumers (`receipt.extracted`, `receipt.suggested`, `ledger.expense.created`) with
  Redis event-id idempotency — 8
- Push notification integration (Expo Push SDK, device token management) — 5
- Frontend: receipt status / notification list UI — 5

#### P4 — Gallery watch (reassigned to help P3) · **15 pts**

With ClickHouse deferred, P4's Sprint 1 dashboard already works end to end — Sprint 2 capacity goes
to the hardest single item in the backlog instead of sitting idle.

- Frontend: gallery-watch native listener — device background limits, iOS/Android divergence — 13
- Verify the Ledger-gRPC fallback still holds once real ingestion traffic exists — 2

#### P5 — Reader SVC + Suggestion SVC (new services) · **19 pts**

- Reader SVC: fetch from S3, VLM extraction call, retry/backoff, DLQ on failure — 8
- Reader SVC: consumer wiring for `receipt.uploaded` — 3
- Suggestion SVC: `GetCategories` gRPC client, merchant-to-category rule matching, event emission — 5
- Suggestion SVC: consumer wiring for `receipt.extracted` — 3

**Sprint 2 total: 98 pts** (18 + 18 + 28 + 15 + 19).

#### Definition of done for Sprint 2

Save a real bank receipt to the gallery, the app detects it, extracts amount/merchant/date, suggests
a category, user confirms with one tap, the expense appears in the ledger and on the dashboard — no
manual entry anywhere in that path. This is the **progress video 2** deliverable.

---

### Sprint 3 — Non-functional requirements (target: **progress video 3, due 27–31 Oct**; presentation **11–15 Nov**)

Lighter, longer window — mostly prep, not new distributed-systems surface. **k3s is the first thing
to cut** if the schedule slips; it has no functional payoff and was always scoped as a stretch item.

| Item | Suggested owner | Points |
| :---- | :---- | :---- |
| ClickHouse projection (moved from Sprint 2) — schema, Kafka consumer, incremental MV | P4 | 13 |
| Load testing, 2 types (constant-arrival + spike) + writeup | P5 | 8 |
| Quality-attribute measurement: scale Reader replicas, kill ClickHouse mid-demo, measure p95 | P1 + P5 | 8 |
| Security hardening: TLS, secrets, PDPA review, rate limits | P1 | 8 |
| Risk matrix, formalized to the graded format | Whole team | 3 |
| Observability polish: Grafana dashboards, alerts, Tempo (stretch) | P5 | 5 |
| C4 L3 component diagrams (Ledger + Receipt internals) | P2 + P3 | 5 |
| C4 L4 code-level diagram, one service | Whole team | 3 |
| ADR-11..18 (closes the deferred Tier B list) | Whole team | 8 |
| Slide deck + demo rehearsal | Whole team | 8 |
| Stabilization buffer | Whole team | 8 |
| **k3s single-node deploy (optional, cut first)** | TBD | 13 |

**Sprint 3 total: 77 pts** (+13 optional).

#### Definition of done for Sprint 3

Load test results with interpretation, risk matrix with ≥3 issues, one measured quality attribute
with before/after evidence, C4 L1–L4 + sequence diagrams, all ADRs current, rehearsed 15-minute demo.
This is the **progress video 3** deliverable and the direct input to the final presentation.

---

### Story point summary

| Sprint | Points | Window |
| :---- | :---- | :---- |
| 1 | 146 | today → 29 Sep – 3 Oct |
| 2 | 98 | → 13–17 Oct |
| 3 | 77 (+13 optional) | → 27–31 Oct, present 11–15 Nov |
| **Total** | **~321** (+13 optional) | |

~28 pts/person/sprint is the assumed capacity band. Sprint 1's Auth owner (P1, 31 pts) and Sprint 2's
Receipt owner (P3, 28 pts) are the tightest — if either slips, that's where to send help first.


## How to Git & Github

You should know

- Commit convention `type(appname|scope): commit name here`
- PR title convention `type(appname|scope): Pull request title`
- Squash merge
- Rebase & rebase onto

Scenario

1. You develop on `branch-a`, already open PR
2. You develop on `branch-b` based on `branch-a`, and also open another PR, as a [Github Stack PRs](https://docs.github.com/en/pull-requests/how-tos/stacked-pull-requests)
3. There's requested change on `branch-a`, you commit it
  - You then should rebase `[branch-b]$ git fetch origin && git rebase branch-a`
4. There's another requested change on `branch-a`, you commit it, reviewer **Squash merge it**
  - Now the squash commit is not exactly the same as `branch-a` changes in the `branch-b`
  - So you do rebase onto `[branch-b]$ git fetch origin && git rebase --onto branch-a a21b8935a89c` where `a21b8935a89c` is the cutoff point. Git will grab everything *after* this commit up to your current `HEAD`
5. Then your branch-b gets merged safely
