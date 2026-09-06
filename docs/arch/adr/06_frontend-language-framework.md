# ADR-06: Frontend Language & Framework

## Context

The client is a native mobile app (ADR-05) that must watch the device media gallery for bank receipt images. The frontend needs dashboard charts, forms, and OAuth sign-in. A small team working over one semester needs a typed component ecosystem that is fast to work in and supports iOS and Android from one codebase.

Options considered:
- **TypeScript + Expo (React Native)**: typed React on native, `expo-media-library` for gallery access and album queries, large ecosystem. Expo simplifies build and permissions.
- **TypeScript + React Native (bare)**: same React model with more native wiring for media APIs.
- **Dart + Flutter**: strong cross-platform UI, but fewer people on the team know Dart. Gallery plugins exist but the team would learn a new UI stack.
- **TypeScript + Next.js**: strong web stack, but cannot implement FR5 gallery watch on mobile (ruled out by ADR-05).

## Decision

Use **TypeScript with Expo (React Native)** for the mobile client.

## Status

Accepted

## Consequences

**Positive**
- Type safety catches mismatches against backend API contracts at compile time.
- React skills transfer from common web experience. One language (TypeScript) on client and shared patterns with backend OpenAPI contracts.
- Expo provides media-library permissions, album queries, and simpler dev builds for a one-semester timeline.
- Native push notifications and OAuth deep links are well supported.

**Negative**
- Expo adds its own conventions and upgrade path. Bare React Native is an escape hatch if a native module is missing.
- Charts and navigation differ from a web stack (e.g. React Navigation, a native chart library instead of web-only options).
