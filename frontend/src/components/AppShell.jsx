import { useMemo, useState } from 'react'
import { useAuth } from '../context/AuthContext.jsx'
import DashboardHome from './DashboardHome.jsx'

const NAV_ITEMS = [
  { key: 'dashboard', label: 'Dashboard', roles: ['admin', 'manager', 'trainer', 'receptionist'], ready: true },
  { key: 'members', label: 'Members', roles: ['admin', 'manager', 'receptionist'], ready: false },
  { key: 'subscriptions', label: 'Subscriptions', roles: ['admin', 'manager', 'receptionist'], ready: false },
  { key: 'attendance', label: 'Attendance', roles: ['admin', 'manager', 'receptionist'], ready: false },
  { key: 'sessions', label: 'Sessions', roles: ['admin', 'manager', 'trainer'], ready: false },
  { key: 'employees', label: 'Employees', roles: ['admin'], ready: false },
]

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
        <div className="sidebar-brand brand-wordmark" aria-label="Iron Forge">
          <span>IRON</span>
          <strong>FORGE</strong>
        </div>

        <nav className="sidebar-nav">
          {navItems.filter((item) => item.visible).map((item) => (
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
        </nav>
      </aside>

      <section className="workspace">
        <header className="workspace-topbar">
          <div>
            <p>Staff workspace</p>
            <h1>Dashboard</h1>
          </div>

          <div className="staff-summary">
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
