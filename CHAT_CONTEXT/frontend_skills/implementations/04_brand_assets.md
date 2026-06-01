# Implementation - 04 Brand Asset Integration

## Status

- Status: implemented
- Feature: Brand Asset Integration
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/04_brand_assets.md`
- Started at: 2026-06-01
- Finished at: 2026-06-01

## Scope implemented

- [x] Routes/pages
- [x] Components
- [x] Styling/responsive states
- [x] API client/state
- [x] Docs/context

## Files changed

- `frontend/src/assets/brand/` - added selected runtime copies of official logo, favicon mark,
  loading spinner, and 404 illustration assets.
- `frontend/public/favicon.svg` - replaced the hand-built favicon with the official `IF` favicon.
- `frontend/public/favicon.ico` - added the official favicon fallback.
- `frontend/public/apple-touch-icon.png` - added the official touch icon.
- `frontend/public/og-image.jpg` - added the website OG image for metadata.
- `frontend/src/components/BrandMark.jsx` - added a shared official brand mark component with full
  and compact variants.
- `frontend/src/components/AppShell.jsx` - replaced sidebar text wordmark with `BrandMark` while
  preserving navigation and logout behavior.
- `frontend/src/components/LoginView.jsx` - replaced login text wordmark with `BrandMark` while
  preserving the auth form.
- `frontend/src/components/StatusMessage.jsx` - added compact official mark and loading spinner for
  full-page session-check states.
- `frontend/src/components/StateBlock.jsx` - added the official 404 illustration for app not-found
  states.
- `frontend/src/index.css` - added asset-kit token aliases, brand mark sizing, status spinner,
  not-found illustration styling, and mobile logo constraints.
- `frontend/index.html` - added favicon fallback, Apple touch icon, theme color, description, and OG
  metadata.
- `CHAT_CONTEXT/frontend_skills/plans/04_brand_assets.md` - marked FE04 implemented.
- `CHAT_CONTEXT/frontend_skills/worklog.md` - updated FE04 implementation handoff.

## Key decisions

- Used `iron-forge-logo-horizontal-dark.svg` for the app UI because the planned primary
  `*-light.svg` variant has a white background and is unsuitable on the current dark staff console.
- Copied only selected runtime assets. Social, print, mockup, preview, and service-icon assets remain
  reference material inside the asset kit.
- Kept all auth, route, and backend API behavior unchanged.
- Kept the app as a compact staff operations console. No marketing hero or new route was added.
- Added metadata assets through `public/`, while React-rendered logo/status assets live under
  `src/assets/brand/` for Vite imports.

## Commands run

```bash
cd frontend
npm run build
```

Result: pass. Vite built 41 modules and emitted production assets.

```bash
git diff --check
```

Result: pass.

```bash
cd frontend
npm run dev -- --host 127.0.0.1 --port 5173
```

Result: pass after sandbox escalation. Vite served `http://127.0.0.1:5173/`.

```bash
curl -I http://127.0.0.1:5173/favicon.svg
curl -I http://127.0.0.1:5173/og-image.jpg
```

Result: pass. Both public assets returned `200`.

Playwright smoke:

- `/login` at 1280px rendered the official logo, Staff Portal label, login form, and API footer.
- `/login` at 320px rendered without visible horizontal overflow in the accessibility snapshot.
- Browser console warnings/errors: none at warning level or higher.

## Known limitations

- Browser visual verification is limited to `/login` local dev-server smoke unless a full FE test
  phase runs.
- Authenticated `/app/dashboard` branding was not browser-verified because no backend session or
  seeded credentials were provided in this implementation pass.
- The official SVG logo still depends on available font fallback for its embedded text, as noted in
  the asset-kit guidelines.
- FE04 does not address the broader FE12 responsive hardening backlog.

## Handoff to review

- Review `BrandMark.jsx`, `AppShell.jsx`, `LoginView.jsx`, `StatusMessage.jsx`, and `StateBlock.jsx`
  for accessibility and whether brand assets are meaningful or decorative in the right places.
- Confirm `index.html` references only existing public assets.
- Confirm no large social/print/mockup files were copied into runtime output.
- Use `$gym-fe-review` with this implementation note.
