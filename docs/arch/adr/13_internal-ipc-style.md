# ADR-13: Internal IPC Style

## Context

The mobile client needs a stable, documented public API. Services also need synchronous collaboration for operations such as Ledger retrieving receipt details, Suggestion retrieving categories, and Dashboard falling back to Ledger records. Using one protocol everywhere is simple, but public clients and internal services have different compatibility and performance needs.

Options considered:
- **REST everywhere**: one protocol and familiar debugging, but internal contracts remain loosely typed and require repeated HTTP schema handling.
- **gRPC everywhere**: strongly typed and efficient, but adds unnecessary complexity to the Expo client and public edge.
- **REST at the edge and gRPC internally**: keeps the public API conventional while giving internal calls generated, typed contracts.

## Decision

Use **REST with JSON over HTTPS** from the mobile client to BogBank through Traefik. Document client-facing APIs with the OpenAPI output from chi and Huma (ADR-09).

Use **gRPC** for synchronous service-to-service calls. Protobuf contracts live in a shared, versioned monorepo location. Internal calls must set deadlines, propagate request or trace identifiers, and use bounded retries only for idempotent operations. Services must use an API or an event rather than connect to another service's datastore.

Asynchronous service collaboration is outside this ADR and is governed by ADR-14.

## Status

Accepted. Supersedes the one-REST-style constraint in ADR-01. ADR-09 remains valid for client-facing HTTP APIs.

## Consequences

**Positive**
- The mobile client uses a broadly supported public protocol and receives generated OpenAPI documentation.
- Internal services get typed contracts and generated clients and servers.
- Protocol boundaries make public and internal compatibility policies explicit.

**Negative**
- The team must maintain both OpenAPI and protobuf contracts.
- Local debugging of gRPC is less convenient than plain JSON over HTTP.
- Incorrect retry or deadline settings can amplify failures between services.
