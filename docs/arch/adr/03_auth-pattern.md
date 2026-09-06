# ADR-03: Google OAuth as Sole Auth Provider

## Context

Target users are college students. Most of them already have Google accounts. Many use their university's Google Workspace. OAuth is easier than building and securing an email/password flow.

Options considered:
- **Google OAuth**: students mostly already have Google accounts and we store no passwords, but sign-in depends on Google being up.
- **Email and password**: works without a Google account, but we would store and reset passwords.
- **Other OAuth providers** (Apple, GitHub): more sign-in choices, which is extra work for little gain because most target users already use Google.
- **OAuth plus email/password**: the most flexible path, and it doubles the auth work we have to build and secure.

## Decision

Use Google OAuth 2.0 as the only identity provider. After a successful OAuth sign-in, the system issues its own short-lived session token (for example a JWT) for later API calls. Internal services then do not need to call Google on every request.

## Status

Accepted

## Consequences

**Positive**
- No password storage, password reset flows, or credential-stuffing attacks to build or secure.
- Fast, familiar sign-in for the target customer.
- Google handles credential security and MFA (Multi-factor authentication).

**Negative**
- Users without a Google account cannot use the app.
- We depend on Google OAuth. A Google OAuth outage blocks all sign-ins, which is very rare.
- We must handle edge cases. A user may revoke the app's Google access. A user's Google email may change. Session tokens must refresh and expire.
