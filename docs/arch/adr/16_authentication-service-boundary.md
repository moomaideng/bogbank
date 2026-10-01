# ADR-16: Authentication Service Boundary and Session Management

## Context

ADR-03 selects Google OAuth 2.0 as the only external identity provider and requires BogBank to issue a short-lived session token after sign-in. The microservice design must now decide where Google verification, local account creation, token issuance, refresh, revocation, and signing-key ownership live.

Options considered:
- **Every service verifies Google tokens directly**: avoids a dedicated authentication service, but couples every service to Google and duplicates identity mapping and failure handling.
- **A small BogBank Auth Service**: centralizes Google federation and BogBank sessions while keeping the implementation scoped to the project's needs.
- **Keycloak or Zitadel**: mature identity platforms, but their deployment and administration surface is large for Google-only sign-in in a semester project.
- **Authelia**: useful for protecting web applications, but it does not replace the required Google-to-BogBank account and mobile-session flow without additional integration.

## Decision

Keep **Google OAuth 2.0 with PKCE** as the external sign-in mechanism. Implement a dedicated **Go Auth Service** that verifies Google's identity response, creates or resolves the BogBank user, and issues BogBank access and refresh tokens.

Access tokens are short-lived and signed by Auth Service. Refresh tokens are revocable, stored as hashes in Auth Service's PostgreSQL store, and rotated when used. Auth Service owns signing-key rotation and publishes or exposes the verification material required by trusted services. Traefik may use Auth Service through ForwardAuth for coarse authentication, while each domain service still enforces resource authorization.

BogBank stores no user passwords. Internal services do not call Google for each request.

## Status

Accepted. Refines ADR-03; it does not supersede it. Google remains the only external identity provider.

## Consequences

**Positive**
- Google-specific logic and session lifecycle rules have one owner.
- Domain services can validate BogBank credentials without depending on Google availability for every request.
- Revocation, refresh-token rotation, and account mapping are consistent across the system.

**Negative**
- The team owns security-sensitive token, key-rotation, and refresh-token code.
- Auth Service and its signing keys become critical security infrastructure.
- A Google outage still blocks new sign-ins even though existing unexpired sessions can continue.
