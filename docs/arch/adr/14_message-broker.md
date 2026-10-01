# ADR-14: Message Broker

## Context

Receipt ingestion contains slow and failure-prone work: storing an image, extracting fields through an external VLM, suggesting a category, updating receipt state, and projecting confirmed records for the dashboard. A synchronous call chain would make the upload request wait for every dependency and would couple availability across services.

Options considered:
- **No broker**: use synchronous REST or gRPC for every step. This is operationally smaller but makes latency and failures propagate through the whole chain.
- **RabbitMQ or NATS**: lighter queue-oriented systems with straightforward work distribution, but they do not match the replayable event-stream model selected for projections.
- **Apache Kafka**: mature event streaming with a large ecosystem, but heavier to operate for this project.
- **Redpanda using the Kafka API**: Kafka-compatible event streaming with a simpler single-binary local deployment and standard Kafka clients.

## Decision

Use **Redpanda through the Kafka API** for asynchronous service-to-service communication.

Delivery is at least once. Receipt and Ledger producers use a transactional outbox so a domain update and its pending event are committed in the same PostgreSQL transaction. Consumers deduplicate by `(consumer_group, event_id)` in Redis and still enforce business-level unique constraints in PostgreSQL. Poison or repeatedly failing messages go to a dead-letter topic.

Every event uses a versioned envelope containing `event_id`, `event_type`, `occurred_at`, `user_id`, `trace_id`, and `payload`. Partition by the affected entity identifier to preserve per-entity ordering. Topic ownership and producer/consumer relationships are documented in the event catalog.

## Status

Accepted. Supersedes the no-message-broker constraint in ADR-01.

## Consequences

**Positive**
- Receipt uploads return without waiting for OCR, suggestion, or dashboard projection.
- Consumers can scale independently and recover by processing retained events.
- The outbox and idempotency rules prevent the most common lost-event and duplicate-event failures.

**Negative**
- Redpanda, outbox relays, consumer groups, schema compatibility, and dead-letter handling add operational work.
- The system becomes eventually consistent across service boundaries.
- At-least-once delivery means every consumer must be designed and tested for duplicates.
