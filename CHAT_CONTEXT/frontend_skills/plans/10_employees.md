# FE Plan - 10 Employees

Status: Planned

Created: 2026-06-02

## Goal

Build the admin-only employee management workspace for staff listing, creation, profile updates,
password reset, and deactivation through status updates.

FE10 should expose the existing backend employee management APIs without leaking password hashes or
stored secrets. The UI must clearly communicate that employee management is admin-only and that
deactivation/password reset revoke active refresh tokens.

## Current Baseline

- Existing route:
  - `/app/employees`
- Route registry does not yet include:
  - `/app/employees/new`
  - `/app/employees/:id`
- Employee management APIs are admin-only.
- Branch list from FE06 can provide branch assignment options.
- Backend never returns `password_hash` or `normalized_email`.

## Screens And Routes

| Route | Access | Behavior |
|---|---|---|
| `/app/employees` | `admin` | Staff list with role/status/branch filters and create action. |
| `/app/employees/new` | `admin` | Create employee account form with initial password. |
| `/app/employees/:id` | `admin` | Employee detail/edit, status update/deactivation, branch assignment, password reset. |

Route config changes:

- Add `/app/employees/new` before `/app/employees/:id`.
- Add `/app/employees/:id`.
- Keep Employees navigation at `/app/employees`.
- Mark employee routes ready after implementation.

## Component Plan

Add or update:

| Path | Responsibility |
|---|---|
| `src/lib/employeesApi.js` | List/get/create/update/reset password helpers. |
| `src/components/employees/EmployeesPage.jsx` | Staff list and filters. |
| `src/components/employees/EmployeeCreateView.jsx` | Create account form. |
| `src/components/employees/EmployeeDetailView.jsx` | Detail/edit/password reset state. |
| `src/components/employees/EmployeeFilters.jsx` | Role/status/branch filter controls. |
| `src/components/employees/EmployeeForm.jsx` | Shared create/update fields. |
| `src/components/employees/PasswordResetPanel.jsx` | Password reset form and confirmation. |
| `src/components/employees/RoleSelector.jsx` | Multi-role selector with trainer level handling. |
| `src/components/employees/employeeFormatters.js` | Role/status/branch/date/ObjectID helpers. |
| `src/App.jsx` | Render FE10 route components. |
| `src/routes/routeConfig.js` | Add employee create/detail routes and mark ready. |
| `src/index.css` | Scoped employee list/form/detail responsive styles. |

## State And API Plan

API helpers:

```js
listEmployees(accessToken, { role, status, branch_id })
getEmployee(accessToken, employeeId)
createEmployee(accessToken, payload)
updateEmployee(accessToken, employeeId, payload)
resetEmployeePassword(accessToken, employeeId, password)
```

Create payload:

- `employee_id`
- `full_name`
- `email`
- `password`
- `role`
- `level`
- `phone`
- `branch_id`
- `status`

Update payload:

- Partial mutable fields:
  - `employee_id`
  - `full_name`
  - `email`
  - `role`
  - `level`
  - `phone`
  - `branch_id`
  - `status`

Password payload:

- `password`

State:

- List page owns filters and list query state.
- Create page owns form and submit state; navigate to detail on success.
- Detail page owns employee query state, update mutation state, and password reset mutation state.
- Branch options load from FE06 branch API when available; fallback to manual branch IDs or disabled
  branch assignment with clear error.

Validation:

- Create:
  - employee_id required
  - full_name required
  - email required
  - password minimum 8 characters
  - at least one role
  - trainer role requires level
  - branch IDs must be ObjectIDs
  - status active/inactive
- Update:
  - require at least one changed mutable field
  - same role/status/level/branch validation
- Password reset:
  - password minimum 8 characters
  - confirmation field should match if included

## UX States

- List loading, empty, error.
- Role/status/branch filter loading and invalid branch ID.
- Create duplicate email/employee ID `409`.
- Detail invalid ID, not found, success.
- Update `409` self-deactivation/self-admin-removal conflicts.
- Deactivation confirmation before setting `status=inactive`.
- Password reset confirmation and success notice; never display or store generated/entered password
  beyond the form.
- Forbidden route is handled by existing route guard for non-admin users.

## Responsive And Accessibility Notes

- Staff list becomes stacked records at 320px.
- Create/edit forms use one-column mobile and two-column desktop.
- Role selector must be keyboard accessible and expose selected roles.
- Password inputs should use visible labels and avoid storing examples/secrets in docs.
- Destructive status changes and password reset require clear confirmation.
- Success/error messages use `aria-live`.

## Docs And Test Plan

Implementation note:

- `CHAT_CONTEXT/frontend_skills/implementations/10_employees.md`

Review note:

- `CHAT_CONTEXT/frontend_skills/reviews/10_employees.md`

Test note:

- `CHAT_CONTEXT/frontend_skills/tests/10_employees.md`

Verification:

```sh
cd frontend
npm run build
npm run dev -- --host 127.0.0.1 --port 5173
```

Checks:

- Admin route access and non-admin forbidden route.
- List filters.
- Create employee success/validation/duplicate conflict.
- Detail update success/conflict.
- Password reset success/error.
- Deactivation confirmation.
- Mobile and desktop layout.

## Backend Contract Gaps

- No hard delete employee endpoint; deactivation is `PATCH /api/v1/employees/:id` with
  `status: "inactive"`.
- Branch options depend on FE06 branch list.
- Employee list is admin-only; do not reuse it as trainer lookup for FE09 for non-admin roles.

## Risks And Boundaries

- Do not store passwords in localStorage, docs, screenshots, or implementation notes.
- Do not show backend-hidden fields (`password_hash`, `normalized_email`).
- Be careful with self-deactivation/admin-role conflicts; surface backend `409` clearly.
- Do not implement audit logs or token revocation views; backend handles token revoke side effects.

## Next Action

Use `$gym-fe-implement` with `CHAT_CONTEXT/frontend_skills/plans/10_employees.md`.
