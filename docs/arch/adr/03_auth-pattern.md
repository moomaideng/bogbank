# ADR-03: Google OAuth as Sole Auth Provider

## Context

`docs/project-description.md` FR1/FR2 already commit to Google OAuth as the sign-in method, with first-sign-in auto-provisioning an account (no separate registration form). NFR1 already commits to not storing user passwords.

Target users are college students, who near-universally already hold Google accounts (many via their university's Google Workspace), making OAuth lower-friction than building and securing an email/password flow.

## Decision

Use Google OAuth 2.0 as the sole identity provider. After a successful OAuth sign-in, the system issues its own short-lived session token (e.g. JWT) for subsequent API calls, so internal services don't need to call Google on every request.

No email/password authentication path, and no other OAuth providers, for v1.

## Status

Accepted

## Consequences

**Positive**
- No password storage, reset flows, or credential-stuffing surface to build or secure.
- Fast, familiar onboarding for the target demographic.
- Delegates credential security and MFA to Google.

**Negative**
- Users without a Google account cannot use the app.
- Hard dependency on Google OAuth's availability; a Google-side outage blocks all sign-ins.
- Must handle edge cases: user revokes the app's Google access, user's Google email changes, session token refresh/expiry.
