# Bogbank mobile

Expo (React Native) client. TypeScript, Expo Router, pnpm. See [ADR-05](../docs/arch/adr/05_client-delivery.md) and [ADR-06](../docs/arch/adr/06_frontend-language-framework.md) for why.

## Route convention

```text
src/
  app/                          # ONLY routes and layouts: every file here becomes a screen
    _layout.tsx                 # root Stack: (auth), (tabs). Platform-owned
    index.tsx                   # redirect to /dashboard
    +not-found.tsx
    (auth)/
      _layout.tsx
      sign-in.tsx
    (tabs)/
      _layout.tsx               # Tabs: registers every slice once. Platform-owned
      dashboard/  _layout.tsx  index.tsx
      ledger/     _layout.tsx  index.tsx
      receipts/   _layout.tsx  index.tsx
      bank-hooks/ _layout.tsx  index.tsx
      account/    _layout.tsx  index.tsx
  modules/
    <slice>/                    # auth, dashboard, ledger, receipts, bank-hooks, account
      index.tsx                 # public surface: the screen(s) the routes render
      components/               # add as the slice grows
      hooks/
      api.ts
  components/ui/                # shared kit
  lib/api/                      # fetch wrapper + base URL, later replaced by a generated OpenAPI client
```

### Slice map

| Slice                          | Route               | Module               |
| ------------------------------ | ------------------- | -------------------- |
| Sign-in                        | `(auth)/sign-in`    | `modules/auth`       |
| Dashboard                      | `(tabs)/dashboard`  | `modules/dashboard`  |
| Ledger (manual income/expense) | `(tabs)/ledger`     | `modules/ledger`     |
| Bank hooks                     | `(tabs)/bank-hooks` | `modules/bank-hooks` |
| Receipts                       | `(tabs)/receipts`   | `modules/receipts`   |
| Account                        | `(tabs)/account`    | `modules/account`    |

### Rules

- **Owned surface.** Touch only `src/app/(tabs)/<slice>/**` and `src/modules/<slice>/**` for your slice, plus `src/app/(auth)/**` and `src/modules/auth/**` for sign-in.
- **Shared files.** `src/app/_layout.tsx` and `src/app/(tabs)/_layout.tsx` change only via whoever owns the app shell. Every tab is already registered there, so adding a screen inside your slice never touches these files.
- **Slice layouts.** Each slice's `_layout.tsx` is a `Stack`. Detail screens (`[id].tsx`, `new.tsx`) go inside that folder and push onto this stack.
- **URLs.** Groups like `(tabs)` and `(auth)` vanish from the URL: `(tabs)/ledger/[id].tsx` is `/ledger/123`. Each slice owns its URL prefix; no two slices define the same path.
- **Route files stay thin.** A route file default-exports a screen imported from `modules/<slice>`. Never put components, hooks, or helpers under `src/app`; Expo Router treats every file there as a route.
- **Cross-slice boundaries.** A slice never imports from another slice's `modules/`. Shared code goes in `components/ui` or `lib`. Navigate between slices with typed `router.push('/ledger/...')`, not a direct component import.
- **File naming.** kebab-case route files; `[param].tsx` for dynamic segments.
- **Conflict hotspots.**
  - `pnpm-lock.yaml`: resolve by re-running `pnpm install`, never by hand-merging.
  - `app.json` plugins and permissions (media-library, notifications): coordinate with whoever owns the app shell before changing these.

## Developing on your machine

This app uses a **local development build**, not Expo Go: planned screens need push notifications and full media-library access, which Expo Go doesn't support. You build your own copy of the app once; after that the daily loop reloads JS instantly, same as Expo Go. See Expo's [development builds introduction](https://docs.expo.dev/develop/development-builds/introduction/) for the concept and its own prerequisites reference.

Two independent choices, pick both:

- **Platform**: Android or iOS. iOS needs a Mac; nobody on the team has verified that path yet, so treat the iOS steps below as untested.
- **Target**: an emulator/simulator running on your own machine, or your physical phone. Either works with the exact same build; only day-to-day networking (stage 3) differs.

### 1. One-time setup

Common to every combination: [pnpm](https://pnpm.io/) and `cp mobile/.env.example mobile/.env` once, so `EXPO_PUBLIC_API_URL` points at the backend.

#### Android

Follow Expo's [Android Studio Emulator guide](https://docs.expo.dev/workflow/android-studio-emulator/) for JDK 17 + Android Studio + `ANDROID_HOME`. Then pick a target:

- **Emulator**: create a virtual device in Android Studio's Device Manager (the guide above covers this) and leave it running.
- **Physical device**: Settings → About phone → tap "Build number" 7 times → Developer Options → enable USB debugging. Plug in, accept the "Allow USB debugging?" prompt, confirm `adb devices` lists it.

#### iOS (macOS only)

Follow Expo's [iOS Simulator guide](https://docs.expo.dev/workflow/ios-simulator/) for Xcode + Command Line Tools. Then pick a target:

- **Simulator**: included with Xcode once you install a runtime from Xcode → Settings → Components.
- **Physical iPhone**: needs a paid Apple Developer account for code signing, and hasn't been set up by anyone on the team yet.

### 2. First build

```bash
task mobile:android   # or: task mobile:ios
```

Either command runs Expo's Continuous Native Generation to generate a gitignored `android/` or `ios/` from `app.json`, compiles it, installs it on whichever target is currently visible (a booted emulator/simulator, or a connected device), and starts Metro. Expect 5–15 minutes the first time; if more than one target is available it prompts you to pick.

### 3. Daily loop

```bash
task mobile:start
```

Reopen the already-installed app; it reconnects to Metro automatically. No rebuild, no reinstall; that's what Fast Refresh is for.

**Reaching the backend** depends on your target:

| Target                                               | What you need to do                                                                                           |
| ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| Android Emulator                                     | Nothing: its virtual network already reaches your laptop's LAN IP.                                            |
| iOS Simulator                                        | Nothing: it shares your Mac's network stack directly; `localhost` on the simulator _is_ your laptop.          |
| Physical device, same Wi-Fi                          | Nothing, as long as your network doesn't isolate clients from each other (common on university/public Wi-Fi). |
| Physical device, different network or isolated Wi-Fi | USB tunnel, see below.                                                                                        |

For a physical device that can't reach your laptop over Wi-Fi:

```bash
task mobile:adb-reverse       # tunnels 8081 (Metro) and 8080 (ledger) over the USB cable
task mobile:start -- --localhost
```

`--localhost` tells Metro to advertise itself as `localhost:8081` instead of your LAN IP. On the phone, "localhost" normally means the phone itself; `adb reverse` intercepts that one port and redirects it down the cable to your laptop. The two only work together: `--localhost` picks the address, `adb reverse` makes that address reach your laptop. If the app shows "Failed to connect to /\<ip\>:8081", this is why. Run the two commands above, then reload the app (shake → Reload, or the Reload button on the red error screen).

**When to redo step 2:** after pulling a change that adds/updates a native dependency, or changes `app.json` (icon, permissions, plugins, package name). A pure JS/TSX change never needs a rebuild. If you want the native project regenerated without also building and launching, `task mobile:prebuild:clean` does just the CNG step.

### Troubleshooting

- **"Failed to connect to /\<ip\>:8081" on a physical device**: Wi-Fi mismatch or client isolation. Use the USB path above.
- **App crashes with a native-module-not-found-style error right after `git pull`**: someone added a native dependency. Redo step 2.
- **`adb devices` shows nothing**: check the USB cable (some are charge-only), re-accept the debugging prompt, or `adb kill-server && adb start-server`.
- **iOS anything**: untested by the team so far (see the setup section above). Ask before assuming a given iOS path works.

## Quality checks

```bash
task mobile:check   # run before pushing: format:check + eslint-config-expo + tsc --noEmit + expo-doctor
task mobile:format  # prettier --write, fixes what mobile:check's format:check flags
```
