# ADR-17: Runtime Environments and Service Discovery

## Context

BogBank has multiple services and supporting datastores, but the team has five members and one semester. Local development must support running one vertical slice without starting the whole system. The guaranteed demonstration environment must also be reproducible and must not depend on optional infrastructure work.

Options considered:
- **Docker Compose for development and demonstration**: simple and reproducible on one machine, with built-in DNS, but not representative of a multi-node production platform.
- **Kubernetes for every environment**: consistent orchestration and discovery, but adds substantial setup, manifests, and troubleshooting to the critical path.
- **Compose as the baseline with k3s as an optional production exercise**: protects delivery while leaving a path to demonstrate Kubernetes deployment later.

## Decision

Use **Docker Compose** as the required development, CI integration, and end-to-end demonstration environment. Compose profiles group infrastructure and each vertical service slice. Makefile targets provide stable commands for starting one slice or the full system.

Services discover each other by Compose service name through embedded DNS. Traefik discovers public routes from Docker labels. One development PostgreSQL container may host private per-service schemas as allowed by ADR-12.

A single-node **k3s** environment is a deferred stretch goal. It must not be required for CI or the guaranteed demonstration. If adopted, services use Kubernetes Service DNS and Traefik's Kubernetes provider. A later ADR must accept the production deployment topology before k3s becomes a committed environment.

## Status

Accepted for the Compose baseline. Production topology and k3s remain deferred.

## Consequences

**Positive**
- Every team member can run a predictable environment with standard container tooling.
- Profiles reduce resource use and make vertical-slice development practical.
- Built-in DNS satisfies service discovery without operating Eureka or Consul.
- Optional Kubernetes work cannot prevent the required demo from shipping.

**Negative**
- Compose does not demonstrate multi-node scheduling, autoscaling, or production high availability.
- Environment parity is limited if k3s is adopted later.
- Profiles and health dependencies require ongoing maintenance as services are added.
