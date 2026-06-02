# FE Plan - 06 Courses And Branches

Status: Planned

Created: 2026-06-02

## Goal

Build the settings UI for course/package and branch reference data.

FE06 should replace the Courses and Branches placeholders with live CRUD screens that managers/admins
can use before subscription, attendance, and session forms depend on this reference data. The cycle
should keep the UI compact and operational, reuse the existing FE03/FE05 page and state patterns, and
avoid hardcoded course/branch options in later features.

## Current Baseline

- Stack: React 18 + Vite 8, no router dependency.
- Existing protected routes:
  - `/app/settings/courses`
  - `/app/settings/branches`
- Existing route registry does not yet include:
  - `/app/settings/courses/:id`
  - `/app/settings/branches/:id`
- FE05 introduced local resource API helpers and local feature component folders.
- `apiRequest` already normalizes backend error envelopes and accepts `accessToken`.
- Course and branch APIs are manager/admin only.

## Screens And Routes

| Route | Access | Behavior |
|---|---|---|
| `/app/settings/courses` | `admin`, `manager` | Course list, create form/panel, empty state, and row actions for detail/edit/delete. |
| `/app/settings/courses/:id` | `admin`, `manager` | Course detail/edit shell for one course. |
| `/app/settings/branches` | `admin`, `manager` | Branch list, create form/panel, nearby search panel, and row actions for detail/edit/delete. |
| `/app/settings/branches/:id` | `admin`, `manager` | Branch detail/edit shell for one branch. |

Route config changes:

- Add course-detail and branch-detail static dynamic routes after their list routes.
- Keep navigation entries at `/app/settings/courses` and `/app/settings/branches`.
- Mark Courses and Branches routes `ready` only after implementation.

## Component Plan

Add or update:

| Path | Responsibility |
|---|---|
| `src/lib/coursesApi.js` | `listCourses`, `getCourse`, `createCourse`, `updateCourse`, `deleteCourse`. |
| `src/lib/branchesApi.js` | `listBranches`, `getBranch`, `createBranch`, `updateBranch`, `deleteBranch`, `nearbyBranches`. |
| `src/components/settings/CoursesPage.jsx` | Course list + create/edit surface. |
| `src/components/settings/CourseDetailView.jsx` | One-course detail/edit/delete state. |
| `src/components/settings/CourseForm.jsx` | Controlled course form. |
| `src/components/settings/BranchesPage.jsx` | Branch list + create/edit + nearby search surface. |
| `src/components/settings/BranchDetailView.jsx` | One-branch detail/edit/delete state. |
| `src/components/settings/BranchForm.jsx` | Controlled branch form with GeoJSON fields. |
| `src/components/settings/NearbyBranchesPanel.jsx` | Nearby search command panel and results. |
| `src/components/settings/settingsFormatters.js` | Local money, coordinate, ObjectID, tag, and status formatting. |
| `src/App.jsx` | Render FE06 route components. |
| `src/routes/routeConfig.js` | Add detail routes and mark implemented settings routes ready. |
| `src/index.css` | Scoped settings form/list/table/nearby responsive styles. |

Keep shared extraction minimal. Use local settings components until the same form/table pattern is
needed by later features.

## State And API Plan

Courses:

```js
listCourses(accessToken)
getCourse(accessToken, courseId)
createCourse(accessToken, payload)
updateCourse(accessToken, courseId, payload)
deleteCourse(accessToken, courseId)
```

Course payload:

- `title`
- `level`
- `allowed_tags`
- `base_price`
- `session_count`
- `description`

Branches:

```js
listBranches(accessToken)
getBranch(accessToken, branchId)
createBranch(accessToken, payload)
updateBranch(accessToken, branchId, payload)
deleteBranch(accessToken, branchId)
nearbyBranches(accessToken, { lng, lat, max_distance, limit })
```

Branch payload:

- `branch_code`
- `name`
- `address`
- `province`
- `location: { type: "Point", coordinates: [lng, lat] }`
- `manager_id` optional ObjectID string

State:

- List pages own local query state: `loading | success | empty | error`.
- Create/update/delete actions own local mutation state.
- After create/update/delete, refetch list or route to detail as appropriate.
- Nearby search state stays local to `NearbyBranchesPanel`.

Validation:

- Course:
  - title required
  - level required
  - base_price positive integer
  - session_count positive integer
  - allowed_tags parsed from comma/newline input into a string array
- Branch:
  - branch_code, name, address, province required
  - location type fixed to `Point`
  - longitude `-180..180`
  - latitude `-90..90`
  - manager_id optional but must be ObjectID when supplied
- Nearby:
  - lng/lat required and range-checked
  - max_distance positive when supplied
  - limit `1..100` when supplied

## UX States

- Initial loading for course and branch lists.
- Empty course/branch list.
- Create/update validation errors.
- `409` duplicate branch code.
- `404` course/branch detail not found.
- Delete confirmation before irreversible delete.
- Delete success with list refresh/navigation away from deleted detail.
- Nearby search loading, empty results, invalid coordinates, and API errors.
- Forbidden route state continues to be owned by route guard.

## Responsive And Accessibility Notes

- Use the existing dense staff console layout: page header, action bar, data panel, table/list.
- At 320px:
  - forms stack into one column
  - list rows convert to compact stacked records
  - nearby search controls wrap without page-level overflow
- Form labels must be visible.
- Field errors should be connected with `aria-describedby`.
- Mutation success/error messages use `role="status"` or `role="alert"` where relevant.
- Delete confirmation must be keyboard reachable before executing `DELETE`.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/06_courses_branches.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/06_courses_branches.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/06_courses_branches.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Checks:

- Course list/create/update/delete with backend or mocked API.
- Branch list/create/update/delete with backend or mocked API.
- Duplicate branch code `409` state.
- Nearby search success, empty, and invalid query.
- Mobile 320-375px and desktop route screenshots or DOM checks.

## Risks And Boundaries

- Do not implement subscription/session forms in FE06.
- Do not add map libraries; nearby search can be numeric lng/lat fields for this cycle.
- Do not rely on employee manager listing unless FE10 or backend support is available; manager_id can
  remain manual optional ObjectID.
- Branch deletion may break later references. Show a confirmation and leave business rejection to the
  backend if it exists.
- Course update endpoint expects a full course shape, not a partial patch.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/06_courses_branches.md`.
