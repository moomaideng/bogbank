# ADR-05: Mobile-First Web + PWA for Client Delivery

## Context

Target users are college students who mostly use phones (NFR5). The product needs dashboard charts, forms, and receipt intake (FR5). FR5 asks for sharing a receipt screenshot via the OS share sheet. A 5-person team has one semester, and the course grades architecture (services, APIs) more than shipping a native app.

Options considered:
- **Responsive website**: works in the mobile browser with no install step, but weak OS share-sheet integration (usually file upload instead).
- **Mobile-first web + PWA**: installable from the browser, can register as a share target on Android; iOS share-to-PWA support is still limited.
- **Native mobile app** (React Native, Flutter, or platform-native): best share sheet and notifications on iOS and Android, but the most work for a semester focused on microservices.

## Decision

Deliver v1 as a mobile-first responsive web app, PWA-capable where it helps (web app manifest, add-to-home-screen). Receipt intake uses in-app upload in v1. Add Android Web Share Target API if time allows. Defer a native app to a later iteration. Keep backend REST APIs stable so a future native client can reuse them.

## Status

Accepted

## Consequences

**Positive**
- Matches NFR5 mobile-first without App Store deployment overhead.
- One codebase for mobile and desktop browsers.
- PWA install and optional Android share target improve phone UX without a native stack.

**Negative**
- FR5 is only fully met on all platforms with in-app upload in v1. Perfect share-sheet UX on iOS may wait for native or better PWA support.
- Push notifications for expense confirmation (FR7) are limited on iOS PWAs compared with native.
- If the product later requires native-only UX, the UI layer will be rebuilt. Backend APIs should stay stable.
