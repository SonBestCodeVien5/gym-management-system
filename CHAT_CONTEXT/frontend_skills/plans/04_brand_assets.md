# FE Plan - 04 Brand Asset Integration

Status: Completed

Created: 2026-06-01

## Goal

Integrate the official Iron Forge brand asset kit into the existing React/Vite staff console without
changing backend behavior or building new business screens.

FE04 should replace the current hand-built text/logo surfaces with selected runtime assets, reconcile
CSS tokens with `frontend/iron-forge-brand-assets/02_colors/`, and prepare consistent favicon,
loading, not-found, and app brand presentation before the resource modules start.

This cycle is frontend-only. It should keep the current dense staff-operations UI direction and avoid
turning the app into a marketing page.

## Current Baseline

- Stack: React 18 + Vite 8.
- Current app routes are already owned by FE03 route matching:
  - `/`
  - `/login`
  - `/app` -> `/app/dashboard`
  - protected module placeholders under `/app/*`
- Current brand surfaces:
  - `frontend/public/favicon.svg` is a hand-created `IF` favicon.
  - `LoginView` and `AppShell` render a text `IRON FORGE` wordmark through `.brand-wordmark`.
  - `StatusMessage` uses text-only `Iron Forge Staff` branding.
  - `StateBlock` renders generic text states for forbidden/not-found/planned/loading/error.
  - `index.html` has only a favicon link and title; no app/OG metadata.
- Current CSS tokens already mostly match the asset kit:
  - `#0d0d0d`, `#ff4614`, `#f0ece4`, `#161616`
  - muted text currently uses `#a8a39a` and `#706b64`, while asset kit includes `#888880`.
- Asset kit source:
  - `frontend/iron-forge-brand-assets/01_logo/`
  - `frontend/iron-forge-brand-assets/02_colors/`
  - `frontend/iron-forge-brand-assets/03_fonts/`
  - `frontend/iron-forge-brand-assets/06_website/`
  - `frontend/iron-forge-brand-assets/07_guidelines/README_Asset_Kit.txt`

## Screens And Routes

No new route is needed.

| Route | Access | FE04 change |
|---|---|---|
| `/login` | Public | Use official logo/wordmark asset while keeping the existing login form and compact panel. |
| `/app/dashboard` | Protected | Use the same brand component in the shell/sidebar; dashboard data remains static from FE02. |
| `/app/*` placeholders | Protected | Use official not-found/planned/status asset treatment where useful without changing route logic. |
| Unknown `/app/*` | Protected | Optionally include the official `404-illustration.svg` in the not-found state. |
| `/` | Public redirect | No visual change beyond metadata/favicon served by the page shell. |

## Asset Selection Plan

Use only assets that are served or imported by the runtime app.

| Source asset | Target use | Decision |
|---|---|---|
| `01_logo/iron-forge-logo-horizontal-light.svg` | Brand mark on dark login/sidebar surfaces | Preferred primary app logo if it renders crisply at sidebar/login sizes. |
| `01_logo/iron-forge-logo-horizontal-mono-white.svg` | Fallback on dark surfaces | Use only if the light logo is visually noisy or low-contrast. |
| `01_logo/iron-forge-logo-transparent.svg` | Neutral brand source | Keep as reference unless needed. |
| `01_logo/iron-forge-favicon-if.svg` or `06_website/iron-forge-favicon-if.svg` | `/favicon.svg` replacement | Replace the current hand-built favicon with the official SVG. |
| `01_logo/iron-forge-favicon.ico` or `06_website/iron-forge-favicon.ico` | Browser/icon fallback | Copy to `public/favicon.ico` only if `index.html` links it. |
| `01_logo/iron-forge-icon-180.png` | Apple touch icon | Add to `public/` if metadata is added. |
| `01_logo/iron-forge-icon-192.png` and `01_logo/iron-forge-icon-512.png` | PWA/app icons | Optional; add only if creating a manifest. |
| `06_website/og-image-1200x630.jpg` | Social preview metadata | Add to `public/` only if `index.html` gets OG tags. |
| `06_website/loading-spinner-icon.svg` | Auth/session loading state | Optional small visual in `StatusMessage`; keep text fallback. |
| `06_website/404-illustration.svg` | App not-found state | Optional image inside `StateBlock` for `tone="notFound"`. |
| `06_website/service-icon-*.svg` | Future marketing/service surfaces | Do not import in FE04; no current staff-console need. |
| `04_social/`, `05_print/`, `08_mockups/`, `previews/` | Reference/export material | Do not copy into runtime/public assets. |

Implementation should prefer imports from `src/assets/brand/` for component assets and `public/` for
favicon/metadata assets. Do not import directly from the full `iron-forge-brand-assets` kit in app
components if that makes the source tree look like runtime code depends on the entire kit.

## Component Plan

Add or update these files under `frontend/src/` and `frontend/public/`:

| Path | Responsibility |
|---|---|
| `src/assets/brand/` | Runtime copies of selected logo/status SVG assets only. |
| `public/favicon.svg` | Official `IF` favicon from the asset kit. |
| `public/favicon.ico` | Optional official favicon fallback if linked from `index.html`. |
| `public/apple-touch-icon.png` | Optional touch icon from `iron-forge-icon-180.png`. |
| `public/og-image.jpg` | Optional OG image if metadata is added. |
| `src/components/BrandMark.jsx` | Shared logo/wordmark component with compact/full variants and accessible label handling. |
| `src/components/AppShell.jsx` | Replace sidebar text wordmark with `BrandMark`; preserve navigation/logout behavior. |
| `src/components/LoginView.jsx` | Replace login text wordmark with `BrandMark`; preserve auth form behavior. |
| `src/components/StatusMessage.jsx` | Optionally render small brand/loading mark for full-page auth checks. |
| `src/components/StateBlock.jsx` | Optionally support an illustration prop or not-found illustration for route not-found states. |
| `src/App.jsx` | Pass not-found illustration only if `StateBlock` API needs explicit details. |
| `src/index.css` | Reconcile brand tokens, logo sizing, image layout, and state illustration styles. |
| `index.html` | Update favicon links, title/description/theme color, and optional OG metadata. |

Keep `BrandMark` small. It should not own navigation or layout; it only renders the selected asset and
optional supporting text such as `Staff Portal` or `Admin Panel`.

## State And API Plan

State:

- No new React state is required.
- Auth/session/route state remains unchanged.
- Loading, forbidden, and not-found status logic remains owned by the current route/auth components.

API:

- No backend API calls.
- No API contract changes.
- No environment variable changes.

Asset path rules:

- Use Vite imports for assets referenced by React components.
- Use `public/` paths for favicon, Apple touch icon, and OG image links in `index.html`.
- Keep copied asset filenames stable and descriptive so future docs/tests can reference them.

## UX States

FE04 should verify brand behavior across current app states:

- login page idle, validation error, backend login error, and submitting states
- full-page auth checking status
- authenticated dashboard shell/sidebar
- module placeholder state
- app forbidden state
- app not-found state
- logout submitting state remains unchanged

Do not add decorative brand artwork to dense data panels unless it improves state recognition, such as
a not-found illustration in an empty route state.

## Responsive And Accessibility Notes

- Logo/wordmark must fit at 320px without pushing the login panel or sidebar/mobile header wider than
  the viewport.
- Brand assets should have stable dimensions through CSS, for example fixed `height` with `max-width:
  100%`.
- Decorative images should use empty `alt=""`; meaningful logo images need an accessible name through
  `aria-label`, `alt`, or visible text.
- Do not rely on the SVG text inside a logo as the only accessible name.
- Preserve visible labels and `aria-live` behavior in existing forms/status messages.
- Ensure not-found/loading illustrations do not hide or replace the text explanation.
- If `index.html` metadata adds theme colors, use asset-kit tokens:
  - `#0D0D0D` for background/theme
  - `#FF4614` for accent
  - `#F0ECE4` for text surfaces

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/04_brand_assets.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/04_brand_assets.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/04_brand_assets.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Static/manual checks:

- Build succeeds and emitted asset names are reasonable.
- `index.html` references only assets that exist in `public/`.
- `/login` renders official branding without layout overflow at 320px and desktop width.
- `/app/dashboard` renders official branding in the sidebar/mobile shell after auth.
- Unknown `/app/*` route renders a readable not-found state; illustration, if added, is visible but
  not required for comprehension.
- Favicon resolves from Vite dev server.
- No large social/print/mockup assets are copied into runtime/public output.

Backend checks:

- Not required for FE04 unless implementation accidentally touches auth behavior.
- If a backend is already running, a quick login/restore/logout smoke can be run to confirm visual
  changes did not disturb auth flow.

## Risks And Boundaries

- Some logo SVGs contain live text and depend on font fallback. Verify rendering before using them in
  critical UI.
- Do not bulk-copy the asset kit into `public/` or `src/assets`; only copy selected runtime assets.
- Do not introduce an icon library or router dependency in this cycle.
- Do not redesign the dashboard or continue FE12 responsive hardening here.
- Do not add marketing hero/website sections; this is a staff operations console.
- Do not change backend API docs because FE04 has no API contract.

## Completion Summary

FE04 is complete as a frontend-only brand asset cycle.

- Official Iron Forge runtime assets are integrated for favicon, metadata, login/sidebar/status
  branding, loading state, and not-found state.
- No backend API contract changed.
- Build, review, and browser verification passed.
- Live backend login/restore/logout remains a residual test gap because no backend credentials or
  seeded local backend session were available during the FE04 test pass.

## Next Action

Use `$gym-git` to review/commit the FE04 changes, or use `$gym-fe-plan` for FE06 Courses And
Branches / FE05 Members depending on the next frontend priority.
