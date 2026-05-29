import { useMemo, useState } from 'react'
import { useAuth } from '../context/AuthContext.jsx'
import DashboardHome from './DashboardHome.jsx'

const NAV_ITEMS = [
  { key: 'dashboard', label: 'Dashboard', group: 'Tong quan', roles: ['admin', 'manager', 'trainer', 'receptionist'], ready: true },
  { key: 'members', label: 'Members', group: 'Tong quan', roles: ['admin', 'manager', 'receptionist'], ready: false },
  { key: 'attendance', label: 'Attendance', group: 'Tong quan', roles: ['admin', 'manager', 'receptionist'], ready: false },
  { key: 'sessions', label: 'Sessions', group: 'Tong quan', roles: ['admin', 'manager', 'trainer'], ready: false },
  { key: 'reports', label: 'Reports', group: 'Tong quan', roles: ['admin', 'manager'], ready: false },
  { key: 'subscriptions', label: 'Subscriptions', group: 'Quan ly', roles: ['admin', 'manager', 'receptionist'], ready: false },
  { key: 'employees', label: 'Employees', group: 'Quan ly', roles: ['admin'], ready: false },
  { key: 'courses', label: 'Courses', group: 'Quan ly', roles: ['admin', 'manager'], ready: false },
  { key: 'branches', label: 'Branches', group: 'Quan ly', roles: ['admin', 'manager'], ready: false },
  { key: 'payments', label: 'Payments', group: 'Quan ly', roles: ['admin', 'manager', 'receptionist'], ready: false },
]

const NAV_GROUPS = ['Tong quan', 'Quan ly']

function hasAllowedRole(employeeRoles, allowedRoles) {
  return allowedRoles.some((role) => employeeRoles.includes(role))
}

function AppShell({ navigate }) {
  const { employee, logout } = useAuth()
  const [activeItem, setActiveItem] = useState('dashboard')
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const roles = employee?.role || []

  const navItems = useMemo(
    () => NAV_ITEMS.map((item) => ({
      ...item,
      available: item.ready && hasAllowedRole(roles, item.roles),
      visible: hasAllowedRole(roles, item.roles),
    })),
    [roles],
  )

  async function handleLogout() {
    setIsLoggingOut(true)
    await logout()
    navigate('/login', { replace: true })
  }

  return (
    <main className="app-shell">
      <aside className="app-sidebar" aria-label="Dieu huong staff portal">
        <div className="sidebar-header">
          <div className="sidebar-brand brand-wordmark" aria-label="Iron Forge">
            <span>IRON</span>
            <strong>FORGE</strong>
          </div>
          <p>Admin Panel</p>
        </div>

        <nav className="sidebar-nav" aria-label="Dieu huong desktop">
          {NAV_GROUPS.map((group) => {
            const groupItems = navItems.filter((item) => item.visible && item.group === group)

            if (!groupItems.length) {
              return null
            }

            return (
              <div className="nav-group" key={group}>
                <p>{group}</p>
                {groupItems.map((item) => (
                  <button
                    key={item.key}
                    type="button"
                    className={item.key === activeItem ? 'nav-item nav-item--active' : 'nav-item'}
                    disabled={!item.available}
                    onClick={() => setActiveItem(item.key)}
                  >
                    <span>{item.label}</span>
                    {!item.ready ? <em>Next</em> : null}
                  </button>
                ))}
              </div>
            )
          })}
        </nav>
      </aside>

      <section className="workspace">
        <header className="workspace-topbar">
          <div>
            <p>Staff workspace</p>
            <h1>Dashboard</h1>
          </div>

          <div className="staff-summary">
            <div className="topbar-actions" aria-label="Dashboard tools">
              <button type="button" title="Notifications coming soon" aria-label="Notifications coming soon" disabled>
                <span className="tool-mark tool-mark--alert" aria-hidden="true" />
              </button>
              <button type="button" title="Search coming soon" aria-label="Search coming soon" disabled>
                <span className="tool-mark tool-mark--search" aria-hidden="true" />
              </button>
            </div>
            <span className="staff-avatar" aria-hidden="true">
              {employee.full_name
                .split(' ')
                .map((part) => part[0])
                .join('')
                .slice(0, 2)
                .toUpperCase()}
            </span>
            <div>
              <strong>{employee.full_name}</strong>
              <span>{employee.email}</span>
            </div>
            <button className="btn-outline" type="button" onClick={handleLogout} disabled={isLoggingOut}>
              {isLoggingOut ? 'Dang thoat' : 'Dang xuat'}
            </button>
          </div>
        </header>

        <nav className="mobile-nav" aria-label="Dieu huong nhanh">
          {navItems.filter((item) => item.visible).map((item) => (
            <button
              key={item.key}
              type="button"
              className={item.key === activeItem ? 'mobile-nav-item mobile-nav-item--active' : 'mobile-nav-item'}
              disabled={!item.available}
              onClick={() => setActiveItem(item.key)}
            >
              {item.label}
            </button>
          ))}
        </nav>

        <DashboardHome employee={employee} navItems={navItems} activeItem={activeItem} />
      </section>
    </main>
  )
}

export default AppShell
