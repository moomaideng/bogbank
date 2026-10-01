# ADR-17: Runtime Environments and Service Discovery

## Context

BogBank has multiple services and supporting datastores, but the team has five members and one semester. Environments must be reproducible for CI and demonstration. Local work must still allow running one vertical slice without starting the whole system. The team must choose an orchestration model that fits that delivery constraint without unnecessary operational overhead.

Options considered (orchestration platform):
- **Docker Compose on a single host**: simple and familiar for laptops and CI, with embedded DNS between services, but no built-in replicated scheduling or rolling updates across nodes.
- **Docker Swarm**: multi-node-capable orchestration using Compose-compatible stack definitions, overlay networking, and built-in service discovery, with less moving parts than Kubernetes.
- **Kubernetes (including lightweight distributions such as k3s)**: industry-standard orchestration and discovery, but adds substantial cluster setup, manifests or operators, and troubleshooting to the critical path for a small team.

## Decision

Use **Docker Swarm** as the orchestration platform for the integrated BogBank stack in development integration, CI end-to-end runs, and demonstration. Service definitions live in **Compose-format stack files** deployed with `docker stack deploy` (single-node Swarm is sufficient for the course).

**Local vertical slices.** Developers may still use plain `docker compose` with profiles or Makefile targets to start infrastructure plus one service slice on a laptop when a full stack deploy is unnecessary. Those flows must not contradict the Swarm stack definitions used for the full system.

**Service discovery.** Services on the Swarm overlay network resolve each other by **service name** through Swarm's embedded DNS. Traefik is the public entry point and discovers routes from **Docker Swarm** service labels (ADR-15). One development PostgreSQL container may host private per-service schemas as allowed by ADR-12.

Kubernetes is out of scope for the required demonstration path. A later ADR would be needed before adopting Kubernetes in place of Swarm.

## Status

Accepted.

## Consequences

**Positive**
- One orchestration model covers integration, CI, and demo without maintaining parallel Kubernetes manifests.
- Swarm provides rolling updates, restarts, and overlay DNS closer to production than Compose-only on a single machine.
- Stack files stay close to Compose, which lowers the learning curve relative to Kubernetes.
- Built-in DNS avoids operating a separate registry such as Eureka or Consul.

**Negative**
- Swarm is less common in industry than Kubernetes; skills transfer is partial.
- Compose profiles and `docker stack deploy` do not align one-to-one; slice workflows need separate compose overrides or documented commands.
- Multi-node Swarm is optional for the course but adds host and networking setup if exercised.
- Health dependencies and stack files require maintenance as services are added.
