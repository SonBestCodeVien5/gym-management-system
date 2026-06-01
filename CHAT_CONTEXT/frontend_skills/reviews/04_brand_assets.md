# Review - 04 Brand Asset Integration

## Status

- Status: reviewed
- Feature: Brand Asset Integration
- Plan file: `CHAT_CONTEXT/frontend_skills/plans/04_brand_assets.md`
- Implementation file: `CHAT_CONTEXT/frontend_skills/implementations/04_brand_assets.md`
- Reviewed at: 2026-06-01

## Review summary

- Result: pass, no blocking findings.
- Build status: pass with `npm run build`.
- Test status: browser smoke covered `/login`, mocked authenticated `/app/dashboard`, and mocked
  authenticated `/app/not-real`.

## Checklist

- [x] UI matches intended design/style.
- [x] Routes and components are scoped cleanly.
- [x] API base URL/env handling is correct.
- [x] Auth/session state is not hardcoded or leaked.
- [x] Loading, empty, and error states are handled where relevant.
- [x] Responsive layout works on mobile and desktop for reviewed FE04 surfaces.
- [x] Accessibility basics are covered.
- [x] Docs/context are aligned.

## Issues found

| Severity | File | Issue | Fix |
|---|---|---|---|
| none | n/a | No blocking review findings. | n/a |

## Review evidence

- `frontend/src/components/BrandMark.jsx` keeps brand rendering isolated and exposes a usable image
  `alt` label for full and compact variants.
- `frontend/src/components/AppShell.jsx` and `frontend/src/components/LoginView.jsx` replace the old
  text wordmark without changing auth, navigation, or form submit logic.
- `frontend/src/components/StatusMessage.jsx` uses the loading mark only in existing session-check
  status surfaces.
- `frontend/src/components/StateBlock.jsx` renders the 404 illustration as decorative
  `alt="" aria-hidden="true"` while retaining the text title/message.
- `frontend/index.html` references public assets that exist:
  - `/favicon.svg`
  - `/favicon.ico`
  - `/apple-touch-icon.png`
  - `/og-image.jpg`
- Runtime assets under `frontend/src/assets/brand/` are limited to selected logo/status assets; no
  social, print, mockup, preview, or service-icon folders were copied into runtime output.

## Commands run

```bash
cd frontend
npm run build
```

Result: pass. Vite built 41 modules.

```bash
git diff --check
```

Result: pass.

Browser checks:

- `/login` at 1280px: official logo, Staff Portal label, login form, and API footer rendered.
- `/app/dashboard` at 1280px with mocked `GET /api/v1/auth/me`: sidebar official logo and admin
  metadata rendered in the authenticated shell.
- `/app/dashboard` at 320px with mocked auth: mobile shell rendered without FE04 brand overflow.
- `/app/not-real` with mocked auth: official 404 illustration rendered with readable not-found text.
- Browser console: no warnings or errors at warning level or higher.

## Fixes applied during review

- None.

## Remaining risks

- Authenticated browser checks used Playwright route interception for `GET /api/v1/auth/me`; they did
  not validate against a live backend or real credentials.
- Broader dashboard/mobile responsive concerns remain out of FE04 scope and belong to FE12.
- The official SVG wordmark still depends on browser font fallback for embedded logo text.

## Handoff to test

- Use `$gym-fe-test` with `CHAT_CONTEXT/frontend_skills/reviews/04_brand_assets.md`.
- Re-run build and browser checks.
- If a backend with credentials is available, smoke login/restore/logout once to confirm the visual
  changes did not disturb auth flow.
