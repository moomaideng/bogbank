# ADR-18: Observability and Telemetry Stack

## Context

BogBank requests cross REST, gRPC, and Redpanda boundaries. Receipt processing also continues after the original HTTP request has returned. Per-container console logs are not sufficient to explain latency, duplicate processing, failed extraction, or consumer backlog across that flow.

Options considered:
- **Local console logs only**: smallest setup, but no cross-service search, metrics history, or event correlation.
- **A self-hosted Grafana stack**: keeps the development environment reproducible and covers metrics and centralized logs with widely used tools.
- **A hosted observability service**: reduces operations but introduces an external account, cost, and data-governance decision that the project does not need.

## Decision

Instrument every Go service with **OpenTelemetry-compatible context propagation**. REST, gRPC, and event handlers propagate `trace_id`; event processing also records `event_id`. Logs are structured and include at least timestamp, level, service, environment, operation, and relevant correlation identifiers without logging tokens, receipt images, or extracted financial payloads.

Use **Prometheus** for metrics, **Loki** for centralized logs, and **Grafana** for dashboards. Minimum metrics include request latency and errors, gRPC latency and errors, Redpanda consumer lag, dead-letter counts, outbox backlog, and receipt pipeline duration.

**Tempo and full distributed-trace storage are deferred** until metrics, logs, and context propagation work reliably. Instrumentation must not require Tempo to be running.

## Status

Accepted. Tempo remains deferred.

## Consequences

**Positive**
- Operators can correlate an upload request with later asynchronous processing.
- Consumer lag, outbox backlog, and dead letters become visible before they cause silent data staleness.
- The same metrics support the project's measurable quality attribute and load tests.

**Negative**
- Instrumentation and dashboards add work to every service.
- Prometheus, Loki, and Grafana increase local resource use.
- Poor label selection can create high-cardinality metrics and excessive storage.
