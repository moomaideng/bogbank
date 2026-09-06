# ADR-05: Client Delivery

## Context

Target users are college students who mostly use phones (NFR5). The product needs dashboard charts, forms, and receipt intake (FR5). FR5 requires watching the device media gallery for new bank receipt images after the user grants permission at hook setup. A PWA or mobile browser cannot access the photo library in the background, so receipt auto-ingest needs a native mobile client.

Options considered:
- **Responsive website**: works in the mobile browser with no install step, but cannot watch the media gallery or run background receipt detection.
- **Mobile-first web + PWA**: installable from the browser, but still blocked from gallery access on mobile, only manual file pick or Web Share Target (explicit user action each time).
- **Native mobile app** (React Native, Flutter, or platform-native): can request media-library permission and watch for new receipt images (e.g. bank app folders or the Screenshots album). Required for the gallery-watch intake model.

## Decision

Deliver v1 as a **native mobile app**. Receipt intake uses **media-gallery watch**: when the customer configures a bank hook, the app requests media-library permission and monitors for newly added receipt images tied to that provider. Detected images upload to the backend for extraction and category suggestion. Keep backend REST APIs stable so other clients could be added later.

## Status

Accepted

## Consequences

**Positive**
- Matches FR5 gallery-watch intake and NFR5 mobile-first on real devices.
- Native push notifications for expense confirmation (FR7) work better than on a PWA.

**Negative**
- More client work than web-only: store releases, native permissions, and iOS/Android testing.
- Bank behavior varies. Some apps block screenshots (`FLAG_SECURE`). Others save receipts to their own album. Manual entry (UC-02) stays the fallback (NFR7).
- iOS background limits. Reliable detection may require the app in the foreground; 24/7 background polling is not guaranteed.
