import { useMemo, useState } from 'react'
import { useAuth } from '../context/AuthContext.jsx'
import { hasAnyRole } from '../lib/permissions.js'
import { APP_NAV_ITEMS, NAV_GROUPS } from '../routes/routeConfig.js'

function AppShell({ navigate, activeRoute, children }) {
  const { employee, logout } = useAuth()
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const [isSidebarOpen, setIsSidebarOpen] = useState(false)
  const roles = employee?.role || []

  const navItems = useMemo(
    () => APP_NAV_ITEMS.map((item) => ({
      ...item,
      visible: hasAnyRole(roles, item.roles),
    })),
    [roles],
  )
  async function handleLogout() {
    setIsLoggingOut(true)
    await logout()
    navigate('/login', { replace: true })
  }

  function handleNavSelect(path) {
    navigate(path)
    setIsSidebarOpen(false)
  }

  function isRouteActive(item) {
    return activeRoute?.navKey === item.navKey
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
                    className={isRouteActive(item) ? 'nav-item nav-item--active' : 'nav-item'}
                    aria-current={isRouteActive(item) ? 'page' : undefined}
                    onClick={() => handleNavSelect(item.path)}
                  >
                    <span>{item.label}</span>
                    {item.status !== 'ready' ? <em>{item.status === 'blocked' ? 'Later' : 'Next'}</em> : null}
                  </button>
                ))}
              </div>
            )
          })}
        </nav>
      </aside>

      <section className="workspace">
        <header className="workspace-topbar">
          <button
            className="mobile-sidebar-toggle"
            type="button"
            aria-expanded={isSidebarOpen}
            aria-controls="mobile-sidebar"
            onClick={() => setIsSidebarOpen((current) => !current)}
          >
            Menu
          </button>

          <div className="topbar-title">
            <p>Staff workspace</p>
            <h1>{activeRoute?.label || 'Workspace'}</h1>
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
            <div className="staff-identity">
              <strong>{employee.full_name}</strong>
              <span>{employee.email}</span>
            </div>
            <button className="btn-outline" type="button" onClick={handleLogout} disabled={isLoggingOut}>
              {isLoggingOut ? 'Dang thoat' : 'Dang xuat'}
            </button>
          </div>
        </header>

        {isSidebarOpen ? (
          <nav className="mobile-sidebar-panel" id="mobile-sidebar" aria-label="Dieu huong mobile">
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
                      className={isRouteActive(item) ? 'mobile-nav-item mobile-nav-item--active' : 'mobile-nav-item'}
                      aria-current={isRouteActive(item) ? 'page' : undefined}
                      onClick={() => handleNavSelect(item.path)}
                    >
                      <span>{item.label}</span>
                      {item.status !== 'ready' ? <em>{item.status === 'blocked' ? 'Later' : 'Next'}</em> : null}
                    </button>
                  ))}
                </div>
              )
            })}
          </nav>
        ) : null}

        {children}
      </section>
    </main>
  )
}

export default AppShell
