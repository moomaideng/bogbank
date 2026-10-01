# ADR-18: Observability and Telemetry Stack

## Context

BogBank requests cross REST, gRPC, and Redpanda boundaries. Receipt processing also continues after the original HTTP request has returned. Per-container console logs are not sufficient to explain latency, duplicate processing, failed extraction, or consumer backlog across that flow.

Options considered (where to run centralized observability):
- **No centralized platform**: services emit logs to stdout only. Smallest setup, but no cross-service log search, metrics history, or retained distributed traces.
- **Self-hosted open-source stack**: run the Grafana ecosystem locally and in CI. Higher local resource use; no vendor account or egress of operational data.
- **Managed observability vendor**: reduces operations but adds cost, an external account, and data-governance choices the project does not want for local development.

## Decision

Choose the **self-hosted open-source stack**. Run **Prometheus** (metrics), **Loki** (logs), **Grafana** (dashboards and exploration), and **Tempo** (traces) in local and CI environments.

**Instrumentation.** Every Go service uses **OpenTelemetry-compatible context propagation**. REST, gRPC, and event handlers propagate `trace_id`; event processing also records `event_id`. Logs are structured and include at least timestamp, level, service, environment, operation, and relevant correlation identifiers without logging tokens, receipt images, or extracted financial payloads.

**Metrics.** Expose request rate, latency, and errors for public HTTP and internal gRPC. Where Redpanda consumers, transactional outbox, or the receipt pipeline exist, also expose messaging and async-pipeline health (for example lag, dead-letter volume, and backlog). Prefer low-cardinality labels; detailed metric names and dashboards are defined during implementation.

**Traces.** Services may export spans to Tempo via OTLP as instrumentation is added. A service must start and run correctly when Tempo is unavailable or not yet configured. Until spans are exported, correlate cross-service work using `trace_id` and `event_id` in Loki and supporting metrics.

## Status

Accepted.

## Consequences

**Positive**
- One operator surface in Grafana and one propagation model across metrics, logs, and traces.
- Shared identifiers still correlate an upload request with later asynchronous processing before every hop emits spans.
- Messaging and pipeline metrics surface consumer lag, outbox backlog, and dead letters before they cause silent data staleness.
- The same metrics support the project's measurable quality attribute and load tests.

**Negative**
- Instrumentation and dashboards add work to every service.
- Prometheus, Loki, Grafana, and Tempo increase local resource use.
- Full trace coverage may lag metrics and logs; multi-hop latency stays harder to inspect until spans exist end to end.
- Poor label selection can create high-cardinality metrics and excessive storage.
