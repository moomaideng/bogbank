# ADR-15: API Gateway and Reverse Proxy

## Context

The Expo client needs one HTTPS entry point while BogBank has several independently deployed services. The edge must route requests, terminate TLS, apply common limits, and reject unauthenticated requests without duplicating those concerns in every service.

Options considered:
- **Expose every service directly**: minimal infrastructure, but leaks internal topology to the client and duplicates cross-cutting controls.
- **Nginx with static configuration**: mature and capable, but route changes must be maintained separately from Compose or Kubernetes service metadata.
- **Kong or a custom Go gateway**: provides more policy and composition features, but adds plugins, storage, or custom code that v1 does not need.
- **Traefik**: integrates with Docker and Kubernetes discovery and provides routing, TLS, middleware, and ForwardAuth without a custom gateway service.

## Decision

Use **Traefik v3** as the API gateway and reverse proxy.

Traefik is the public entry point for client REST traffic. It owns TLS termination, path or host routing, request limits, rate limiting, and ForwardAuth integration with Auth Service. In development it discovers routes from Docker labels. If a Kubernetes environment is adopted later, it may use the Kubernetes CRD provider.

Traefik must not contain business logic, aggregate domain responses, or replace service authorization checks. A service remains responsible for authorization decisions concerning its own resources.

## Status

Accepted

## Consequences

**Positive**
- The mobile client uses one stable system endpoint.
- Common edge policies are configured once.
- Docker label discovery removes a separate static route registry in development.

**Negative**
- Traefik becomes critical infrastructure and must have health checks and observable configuration.
- Incorrect middleware or route labels can block every public API.
- Advanced API composition would require a separate BFF or application service rather than gateway configuration.
